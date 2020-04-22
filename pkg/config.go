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

// TlsVersion will be used current TLS Policy
// by this value. The acceptable value MUST BE TLS 1.1/1.2/1.3
// Note that the TLS 1.3 supports `TLS_AES_128_GCM_SHA256`,
// `TLS_AES_256_GCM_SHA384` and `TLS_CHACHA20_POLY1305_SHA256` ciphers only.
// Default is "TLS_1_2"
func (*Config) TlsVersion() string {
	tlsVersions := []string{"TLS_1_1", "TLS_1_2", "TLS_1_3"}
	for _, ver := range tlsVersions {
		if os.Getenv("MIN_TLS_VERSION") != ver || os.Getenv("MIN_TLS_VERSION") == "" {
			return "TLS_1_2"
		}
	}
	return os.Getenv("MIN_TLS_VERSION")
}

// SslProfile returns the string value of pre-configured profile
// defined by GCP. The acceptable value MUST BE one of these values
// COMPATIBLE/MODERN/RESTRICTED. We have not suppoerted CUSTOM value yet
// Default is "MODERN"
func (*Config) SslProfile() string {
	sslProfiles := []string{"COMPATIBLE", "MODERN", "RESTRICTED"}
	for _, profile := range sslProfiles {
		if os.Getenv("SSL_PROFILE") != profile || os.Getenv("SSL_PROFILE") == "" {
			return "MODERN"
		}
	}
	return os.Getenv("SSL_PROFILE")
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
