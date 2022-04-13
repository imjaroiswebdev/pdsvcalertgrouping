package main

import (
	"fmt"
	"log"
	"os"

	"github.com/imjaroiswebdev/pdsvcalertgrouping"
)

func main() {
	if err := validateEnvVars(); err != nil {
		log.Fatal(err)
	}

	const NO_BASE_URL = ""
	svc, err := pdsvcalertgrouping.CreateServiceWithAlertGrouping(os.Getenv("PAGERDUTY_TOKEN"), os.Getenv("PD_USER_EMAIL"), NO_BASE_URL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nDetails of Service with alert grouping...\n\n%+v\n", svc)
}

func validateEnvVars() error {
	envVarsList := []string{"PAGERDUTY_TOKEN", "PD_USER_EMAIL"}

	for _, k := range envVarsList {
		_, ok := os.LookupEnv(k)
		if !ok {
			return fmt.Errorf("Error Environment Variable %s not supplied", k)
		}
	}
	return nil
}
