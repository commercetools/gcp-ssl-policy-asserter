package sslpolicy

import (
	"strings"

	compute "google.golang.org/api/compute/v1"
)

// SelectProxies removes black listed proxies from our list of
// targets.
func SelectProxies(proxies *compute.TargetHttpsProxyList, blacklist map[string]struct{}) *compute.TargetHttpsProxyList {
	var tmp []*compute.TargetHttpsProxy
	for _, proxy := range proxies.Items {
		splitURLMapURL := strings.Split(proxy.UrlMap, "/")
		urlMapName := splitURLMapURL[len(splitURLMapURL)-1]
		_, ok := blacklist[urlMapName]
		if !ok {
			tmp = append(tmp, proxy)
		}
	}
	proxies.Items = tmp
	return proxies
}
