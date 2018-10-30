package main

import (
	"context"
	"log"
	"os"

	sslpolicy "github.com/commercetools/gcp-ssl-policy-asserter/pkg"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

func main() {
	var confPath string
	if len(os.Args) > 1 {
		confPath = os.Args[1]
	}
	config, err := sslpolicy.LoadConfig(confPath)

	if err != nil {
		log.Fatalf("Could not read config %s: %v", confPath, err)
	}

	ctx := context.Background()
	client, err := google.DefaultClient(ctx, compute.ComputeScope)

	if err != nil {
		log.Fatal(err)
	}
	svc, err := compute.New(client)
	if err != nil {
		log.Fatal(err)
	}
	if config.Project() == "" {
		log.Fatalf("GOOGLE_PROJECT environment variable not set.")
	}

	sslPolicy, err := sslpolicy.AssertPolicy(config, svc)
	if err != nil {
		log.Fatal(err)
	}

	prxySvc := compute.NewTargetHttpsProxiesService(svc)
	listCall := prxySvc.List(config.Project())
	prxyList, err := listCall.Do()

	if err != nil {
		log.Fatal(err)
	}

	for k := range config.IgnoreProxies {
		log.Printf("Ignoring HTTPSTargetProxies using URLMap %s", k)
	}
	prxyList = sslpolicy.SelectProxies(prxyList, config.IgnoreProxies)
	log.Printf("Found %d HttpsProxies in %s project", len(prxyList.Items), config.Project())

	sslPolicyReference := compute.SslPolicyReference{
		SslPolicy: sslPolicy.SelfLink,
	}

	for _, prxy := range prxyList.Items {
		if prxy.SslPolicy != sslPolicy.SelfLink {
			setPolicyCall := prxySvc.SetSslPolicy(config.Project(), prxy.Name, &sslPolicyReference)
			_, err := setPolicyCall.Do()
			if err != nil {
				log.Printf("WARNING: %s HttpsProxy errored: %v", prxy.Name, err)
			} else {
				log.Printf("%s SSLPolicy set to %s \n", prxy.Name, sslPolicy.SelfLink)
			}
		}
	}
}
