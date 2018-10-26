package sslpolicy

/*
* This file should contain all of the functions that
* return configuration values.
 */
import "os"

// PolicyName will be used by the service to fetch the current TLS Policy
// by this name and if not found will create it.
// The Version string at the end is very important! It will be used
// to upgrade TLS policies in the future.
func PolicyName() string {
	return os.Getenv("SSL_POLICY_NAME")
}

// Project returns the ID (not the display name) of the Google Cloud
// Project to work under.
func Project() string {
	return os.Getenv("GOOGLE_PROJECT")
}
