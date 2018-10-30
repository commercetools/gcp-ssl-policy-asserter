package sslpolicy

import (
	"testing"

	"google.golang.org/api/compute/v1"
)

// TestSelectProxies verifies that
// SelectProxies properly removes blacklisted items
// and keeps all other items.
func TestSelectProxies(t *testing.T) {
	proxyList := compute.TargetHttpsProxyList{
		Items: []*compute.TargetHttpsProxy{
			{
				Name:   "keep1",
				UrlMap: "https://www.googleapis.compute/v1/projects/project/global/urlMaps/keeper",
			},
			{
				Name:   "remove1",
				UrlMap: "projects/project/global/urlMaps/loser",
			},
			{
				Name:   "keep2",
				UrlMap: "global/urlMaps/keeper",
			},
			{
				Name:   "remove2",
				UrlMap: "global/urlMaps/anotherLoser",
			},
		},
	}

	blacklist := map[string]struct{}{
		"loser":        {},
		"anotherLoser": {},
	}

	expectedBlackList := map[string]struct{}{
		"remove1": {},
		"remove2": {},
	}
	expected := map[string]struct{}{
		"keep1": {},
		"keep2": {},
	}

	output := SelectProxies(&proxyList, blacklist)

	for _, v := range output.Items {

		_, ok := expectedBlackList[v.Name]
		if ok {
			t.Fatalf("expected not to find %s", v.Name)
		}
	}

	for _, v := range output.Items {
		_, ok := expected[v.Name]
		if !ok {
			t.Fatalf("missing %s", v.Name)
		}
	}
}
