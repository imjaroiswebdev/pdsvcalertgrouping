package pdsvcalertgrouping

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestServiceCreationWithAlertGrouping(t *testing.T) {
	setup()
	defer teardown()
	mockServerCalls(t)

	createdSvc, err := CreateServiceWithAlertGrouping("token_just_testing", "testing@email.com", baseURLTesting)
	if err != nil {
		t.Fatal(err)
	}

	want := struct {
		Name          string
		AlertCreation string
	}{
		Name:          svcName,
		AlertCreation: alertCreationConfig,
	}
	actual := struct {
		Name          string
		AlertCreation string
	}{
		Name:          createdSvc.Name,
		AlertCreation: createdSvc.AlertCreation,
	}

	if !reflect.DeepEqual(actual, want) {
		t.Errorf("returned \n\n%#v want \n\n%#v", createdSvc, want)
	}
}

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// Base URL just for testing
	baseURLTesting string

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	baseURLTesting = server.URL
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func mockServerCalls(t *testing.T) {
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"users":[{"name":"Jos\u00e9 Antonio Reyes","email":"testing@email.com","time_zone":"America\/Santiago","color":"red","avatar_url":"https:\/\/secure.gravatar.com\/avatar\/865251c0b3fdf17b56bcc276bb46c2bd.png?d=mm&r=PG","billed":true,"role":"owner","description":null,"invitation_sent":false,"job_title":null,"teams":[]}],"limit":25,"offset":0,"total":null,"more":false}`))
	})
	mux.HandleFunc("/escalation_policies", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"escalation_policies":[{"id":"BX0A6UF","type":"escalation_policy","summary":"My Escalation Policy","self":"https:\/\/api.pagerduty.com\/escalation_policies\/BX0A6UF","html_url":"https:\/\/jareyes.pagerduty.com\/escalation_policies\/BX0A6UF","name":"My Escalation Policy","escalation_rules":[{"id":"PETAAOK","escalation_delay_in_minutes":10,"targets":[{"id":"BWDWRZA","type":"user_reference","summary":"Jos\u00e9 Antonio Reyes","self":"https:\/\/api.pagerduty.com\/users\/BWDWRZA","html_url":"https:\/\/jareyes.pagerduty.com\/users\/BWDWRZA"}]}],"services":[{"id":"B93R1LL","type":"service_reference","summary":"My Service with alert grouping","self":"https:\/\/api.pagerduty.com\/services\/B93R1LL","html_url":"https:\/\/jareyes.pagerduty.com\/service-directory\/B93R1LL"}],"num_loops":2,"teams":[],"description":null,"on_call_handoff_notifications":"if_has_services","privilege":null}],"limit":25,"offset":0,"more":false,"total":null}`))
	})
	mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"services":[{"id":"B93R1LL","name":"My Service with alert grouping","description":null,"created_at":"2022-04-12T23:52:49-05:00","updated_at":"2022-04-12T23:53:48-05:00","status":"active","teams":[],"alert_creation":"create_alerts_and_incidents","addons":[],"scheduled_actions":[],"support_hours":null,"last_incident_timestamp":null,"escalation_policy":{"id":"BX0A6UF","type":"escalation_policy_reference","summary":"My Escalation Policy","self":"https://api.pagerduty.com/escalation_policies/BX0A6UF","html_url":"https://jareyes.pagerduty.com/escalation_policies/BX0A6UF"},"incident_urgency_rule":{"type":"constant","urgency":"high"},"acknowledgement_timeout":3600,"auto_resolve_timeout":3600,"alert_grouping":null,"alert_grouping_timeout":null,"alert_grouping_parameters":{"type":null,"config":null},"integrations":[],"response_play":null,"type":"service","summary":"My Service with alert grouping","self":"https://api.pagerduty.com/services/B93R1LL","html_url":"https://jareyes.pagerduty.com/service-directory/B93R1LL"}],"limit":25,"offset":0,"total":null,"more":false}`))
	})
	mux.HandleFunc("/services/B93R1LL", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Write([]byte(`{"service":{"id":"B93R1LL","name":"My Service with alert grouping","description":null,"created_at":"2022-04-12T23:52:49-05:00","updated_at":"2022-04-13T00:13:00-05:00","status":"active","teams":[],"alert_creation":"create_alerts_and_incidents","addons":[],"scheduled_actions":[],"support_hours":null,"last_incident_timestamp":null,"escalation_policy":{"id":"BX0A6UF","type":"escalation_policy_reference","summary":"My Escalation Policy","self":"https:\/\/api.pagerduty.com\/escalation_policies\/BX0A6UF","html_url":"https:\/\/jareyes.pagerduty.com\/escalation_policies\/BX0A6UF"},"incident_urgency_rule":{"type":"constant","urgency":"high"},"acknowledgement_timeout":3600,"auto_resolve_timeout":3600,"alert_grouping":null,"alert_grouping_timeout":null,"alert_grouping_parameters":{"type":null,"config":null},"integrations":[],"response_play":null,"type":"service","summary":"My Service with alert grouping","self":"https:\/\/api.pagerduty.com\/services\/B93R1LL","html_url":"https:\/\/jareyes.pagerduty.com\/service-directory\/B93R1LL"}}`))
	})
}
