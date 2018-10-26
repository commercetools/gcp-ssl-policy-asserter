package sslpolicy

import (
	"errors"
	"log"
	"time"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

// NewSslPolicy returns instance of the configuration options
// necessary to create our globally enforced SSL Policy.
func NewSslPolicy() compute.SslPolicy {
	return compute.SslPolicy{
		Description:   "Commercetools TLS policy: modern features and TLS 1.2 only.",
		Name:          PolicyName(),
		Profile:       "MODERN",
		MinTlsVersion: "TLS_1_2",
	}
}

func pollOperationStatus(project string, svc *compute.Service, operation *compute.Operation) error {
	operSvc := compute.NewGlobalOperationsService(svc)
	statusCall := operSvc.Get(project, operation.Name)
	operation, _ = statusCall.Do()

	for timeout := time.After(10 * time.Second); operation.Status != "DONE"; {
		select {
		case <-timeout:
			return errors.New("polling SSLPolicy creation operation timed out")
		default:
			time.Sleep(2 * time.Second)
			statusCall := operSvc.Get(project, operation.Name)
			var err error
			operation, err = statusCall.Do()

			if err != nil {
				return err
			}
			log.Printf("STATUS: create %s is %s", PolicyName(), operation.Status)
		}
	}
	return nil
}

func pollSSLPolicy(project string, policySvc *compute.SslPoliciesService) (*compute.SslPolicy, error) {
	var (
		currPolicy *compute.SslPolicy
		err        error
	)
	for timeout := time.After(10 * time.Second); ; {
		select {
		case <-timeout:
			return nil, errors.New("polling SSL Policy timed out")
		default:
			getPolicyCall := policySvc.Get(project, PolicyName())
			currPolicy, err = getPolicyCall.Do()

			if err == nil {
				log.Printf("SSLPolicy %s exists", PolicyName())
				return currPolicy, nil
			}
			log.Printf("Polling %s status: %d", PolicyName(), err.(*googleapi.Error).Code)
			time.Sleep(1 * time.Second)
		}
	}
}

// AssertPolicy ensures that a policy exists
// that matches our expectations.
func AssertPolicy(project string, svc *compute.Service) (*compute.SslPolicy, error) {
	policySvc := compute.NewSslPoliciesService(svc)

	getPolicyCall := policySvc.Get(project, PolicyName())
	currPolicy, err := getPolicyCall.Do()

	switch err.(type) {
	case *googleapi.Error:
		if err.(*googleapi.Error).Code == 404 {
			// Clean out the old error from above.
			// := causes variable shadowing on err
			err = nil
			log.Printf("SSLPolicy %s not found, creating...", PolicyName())
			sslPolicy := NewSslPolicy()
			createPolicyCall := policySvc.Insert(project, &sslPolicy)
			operation, err := createPolicyCall.Do()
			if err != nil {
				return nil, err
			}
			// Log any warning from the operation.
			for _, warning := range operation.Warnings {
				log.Printf("WARNING: %s", warning.Message)
			}

			err = pollOperationStatus(project, svc, operation)

			if err != nil {
				log.Fatalf("Create SSLPolicy operation failed: %v", err)
			}

			currPolicy, err = pollSSLPolicy(project, policySvc)

			if err != nil {
				log.Fatalf("Could not find SSLPolicy after creation: %v", err)
			}

		}
	}

	if err != nil {
		return nil, err
	}
	log.Print(currPolicy)
	return currPolicy, nil
}
