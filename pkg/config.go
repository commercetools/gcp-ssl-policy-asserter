package sslpolicy

/*
* This file should contain all of the functions that
* return configuration values.
 */
import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// config defines the structure of the YAML file
// that should be unmarshalled from disk.
type config struct {
	IgnoreProxies []string `yaml:"ignoreProxies,omitempty"`
}

// Config contains the configuration values that the rest of
// the program will leverage.
// Fields from the YAML are manipulated to become more convenient
// for the internal program.
type Config struct {
	IgnoreProxies map[string]struct{}
}

// PolicyName will be used by the service to fetch the current TLS Policy
// by this name and if not found will create it.
// The Version string at the end is very important! It will be used
// to upgrade TLS policies in the future.
func (*Config) PolicyName() string {
	return os.Getenv("SSL_POLICY_NAME")
}

// Project returns the ID (not the display name) of the Google Cloud
// Project to work under.
func (*Config) Project() string {
	return os.Getenv("GOOGLE_PROJECT")
}

// TlsVersion will be used current TLS Policy by this value.
// The acceptable value MUST BE either`TLS_1_1`/`TLS_1_2`/`TLS_1_3`
// Default is "TLS_1_2"
func (*Config) TlsVersion() string {
	tlsVersions := []string{"TLS_1_1", "TLS_1_2", "TLS_1_3"}
	inputVersion := os.Getenv("MIN_TLS_VERSION")
	result := "TLS_1_2"
	for _, ver := range tlsVersions {
		if inputVersion == ver {
			result = ver
			break
		}
	}
	return result
}

// SslProfile returns the string value of pre-configured profile defined by GCP.
// The acceptable value MUST BE either `COMPATIBLE`/`MODERN`/`RESTRICTED`.
// Default is "MODERN"
func (*Config) SslProfile() string {
	sslProfiles := []string{"COMPATIBLE", "MODERN", "RESTRICTED"}
	inputProfile := os.Getenv("SSL_PROFILE")
	result := "MODERN"
	for _, profile := range sslProfiles {
		if inputProfile == profile {
			result = profile
			break
		}
	}
	return result
}

// listToContainsMap converts a list of strings
// into a map that will be used for "contains" checking.
func listToContainsMap(x []string) (result map[string]struct{}) {
	result = make(map[string]struct{})
	var empty struct{}
	for _, v := range x {
		result[v] = empty
	}
	return result
}

// LoadConfig returns an instance of the
// unmarshalled configuration file.
func LoadConfig(path string) (*Config, error) {
	var (
		data []byte
		err  error
	)
	if path != "" {
		data, err = ioutil.ReadFile(path)

		if err != nil {
			return nil, err
		}
	}
	var rawConfig config
	var result Config
	yaml.Unmarshal(data, &rawConfig)
	result.IgnoreProxies = listToContainsMap(rawConfig.IgnoreProxies)

	return &result, err
}
