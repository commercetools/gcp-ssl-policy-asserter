package main

import (
	"context"
	"log"

	sslpolicy "github.com/commercetools/gcp-ssl-policy-asserter/pkg"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

func main() {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, compute.ComputeScope)

	if err != nil {
		log.Fatal(err)
	}
	svc, err := compute.New(client)
	if err != nil {
		log.Fatal(err)
	}
	project := sslpolicy.Project()
	if project == "" {
		log.Fatalf("GOOGLE_PROJECT environment variable not set.")
	}

	sslPolicy, err := sslpolicy.AssertPolicy(project, svc)
	if err != nil {
		log.Fatal(err)
	}

	prxySvc := compute.NewTargetHttpsProxiesService(svc)
	listCall := prxySvc.List(project)
	prxyList, err := listCall.Do()

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Found %d HttpsProxies in %s project", len(prxyList.Items), project)
	sslPolicyReference := compute.SslPolicyReference{
		SslPolicy: sslPolicy.SelfLink,
	}
	for _, prxy := range prxyList.Items {
		if prxy.SslPolicy != sslPolicy.SelfLink {
			setPolicyCall := prxySvc.SetSslPolicy(project, prxy.Name, &sslPolicyReference)
			_, err := setPolicyCall.Do()
			if err != nil {
				log.Printf("WARNING: %s HttpsProxy errored: %v", prxy.Name, err)
			} else {
				log.Printf("%s SSLPolicy set to %s \n", prxy.Name, sslPolicy.SelfLink)
			}
		}

	}

}
