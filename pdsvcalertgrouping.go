package pdsvcalertgrouping

import (
	"fmt"

	"github.com/heimweh/go-pagerduty/pagerduty"
)

var (
	svcName              = "My Service with alert grouping"
	escalationPolicyName = "My Escalation Policy"
	ackTimeout           = 3600
	autoResolveTimeout   = 3600
	alertCreationConfig  = "create_alerts_and_incidents"
)

type EscalationRuleTarget struct {
	Type string
	ID   string
}

type EscalationPolicyInput struct {
	ID                       string
	Name                     string
	NumLoops                 int
	EscalationDelayInMinutes int
	Targets                  []*EscalationRuleTarget
}

func CreateServiceWithAlertGrouping(apiToken string, pdUserEmail string, baseURL string) (*pagerduty.Service, error) {
	client, err := createPDClient(apiToken, baseURL)
	if err != nil {
		return nil, err
	}

	usersList, err := getUsersList(client)
	if err != nil {
		return nil, err
	}

	user := getUserByEmail(usersList, pdUserEmail)
	if user == nil {
		return nil, fmt.Errorf("Error user %s not found", pdUserEmail)
	}

	escalationPolicy, err := getOrCreateEscalationPolicy(client, escalationPolicyName, user)
	if err != nil {
		return nil, err
	}

	svc, err := getOrCreateService(client, svcName, escalationPolicy)
	if err != nil {
		return nil, err
	}

	updatedSvc, err := updateServiceAlert(client, alertCreationConfig, svc)
	if err != nil {
		return nil, err
	}

	return updatedSvc, nil
}

func createPDClient(apiToken string, baseURL string) (*pagerduty.Client, error) {
	fmt.Println("Configuring PagerDuty Client")

	clientConfig := &pagerduty.Config{
		Token: apiToken,
	}
	if baseURL != "" {
		clientConfig.BaseURL = baseURL
	}
	client, err := pagerduty.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getUserByEmail(usersList []*pagerduty.User, lookupEmail string) *pagerduty.User {
	for _, u := range usersList {
		if u.Email == lookupEmail {
			return u
		}
	}

	return nil
}

func getUsersList(c *pagerduty.Client) ([]*pagerduty.User, error) {
	fmt.Println("Requesting Users List")
	usersList, resp, err := c.Users.List(&pagerduty.ListUsersOptions{})
	if err != nil {
		return nil, err
	}

	err = resp.Response.Body.Close()
	if err != nil {
		fmt.Printf("Error while trying to close connection: %v", err)
	}
	return usersList.Users, nil
}

func getOrCreateEscalationPolicy(c *pagerduty.Client, name string, u *pagerduty.User) (*pagerduty.EscalationPolicy, error) {
	fmt.Println("Getting Escalation Policy")
	var escalationPolicy *pagerduty.EscalationPolicy
	epList, resp, err := c.EscalationPolicies.List(&pagerduty.ListEscalationPoliciesOptions{})
	if err != nil {
		return nil, err
	}
	for _, ep := range epList.EscalationPolicies {
		if ep.Name == name {
			return ep, nil
		}
	}

	fmt.Println("Escalation Policy not found. Creating...")
	escalationPolicyDef := buildEscalationPolicyStruct(EscalationPolicyInput{
		ID:                       u.ID,
		Name:                     name,
		EscalationDelayInMinutes: 10,
		NumLoops:                 2,
		Targets: []*EscalationRuleTarget{
			{
				Type: "user",
				ID:   u.ID,
			},
		},
	})
	escalationPolicy, resp, err = c.EscalationPolicies.Create(escalationPolicyDef)
	if err != nil {
		return nil, err
	}

	err = resp.Response.Body.Close()
	if err != nil {
		fmt.Printf("Error while trying to close connection: %v", err)
	}

	return escalationPolicy, nil
}

func getOrCreateService(c *pagerduty.Client, name string, ep *pagerduty.EscalationPolicy) (*pagerduty.Service, error) {
	fmt.Println("Getting Service")
	var svc *pagerduty.Service
	svcList, resp, err := c.Services.List(&pagerduty.ListServicesOptions{})
	if err != nil {
		return nil, err
	}

	err = resp.Response.Body.Close()
	if err != nil {
		fmt.Printf("Error while trying to close connection: %v", err)
	}

	for _, s := range svcList.Services {
		if s.Name == name {
			return s, nil
		}
	}

	fmt.Println("Service not found. Creating...")
	svcDef := &pagerduty.Service{
		Name:                   svcName,
		AcknowledgementTimeout: &ackTimeout,
		AutoResolveTimeout:     &autoResolveTimeout,
		EscalationPolicy: &pagerduty.EscalationPolicyReference{
			HTMLURL: ep.HTMLURL,
			ID:      ep.ID,
			Self:    ep.Self,
			Summary: ep.Summary,
			Type:    ep.Type,
		},
	}
	fmt.Printf("Service Config...\n%+v\n", *svcDef)
	svc, resp, err = c.Services.Create(svcDef)
	if err != nil {
		return nil, err
	}

	err = resp.Response.Body.Close()
	if err != nil {
		fmt.Printf("Error while trying to close connection: %v", err)
	}

	return svc, nil
}

func updateServiceAlert(c *pagerduty.Client, alertCreationConfig string, svc *pagerduty.Service) (*pagerduty.Service, error) {
	alertGroupingConfig := "intelligent"
	fmt.Printf("Updating Service's alerts creation to %q and alerts grouping to %q\n", alertCreationConfig, alertGroupingConfig)
	svc.AlertCreation = alertCreationConfig
	svc.AlertGroupingParameters = &pagerduty.AlertGroupingParameters{
		Type: &alertGroupingConfig,
	}
	updatedSvc, resp, err := c.Services.Update(svc.ID, svc)
	if err != nil {
		return nil, err
	}

	err = resp.Response.Body.Close()
	if err != nil {
		fmt.Printf("Error while trying to close connection: %v", err)
	}

	return updatedSvc, nil
}

func buildEscalationPolicyStruct(input EscalationPolicyInput) *pagerduty.EscalationPolicy {
	escalationPolicy := &pagerduty.EscalationPolicy{
		Name:     input.Name,
		NumLoops: &input.NumLoops,
		ID:       input.ID,
	}

	var escalationRules []*pagerduty.EscalationRule
	var targets []*pagerduty.EscalationTargetReference
	for _, t := range input.Targets {
		targets = append(targets, &pagerduty.EscalationTargetReference{
			ID:   t.ID,
			Type: t.Type,
		})
	}
	escalationRules = append(escalationRules, &pagerduty.EscalationRule{
		EscalationDelayInMinutes: input.EscalationDelayInMinutes,
		ID:                       input.ID,
		Targets:                  targets,
	})

	escalationPolicy.EscalationRules = escalationRules
	return escalationPolicy
}
