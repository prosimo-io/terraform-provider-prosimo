package client

import (
	"context"
	"fmt"
	"testing"
)

func TestCreateCacheRule(t *testing.T) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient("https://myeveroot1623834688941.dashboard.psonar.us/", "stU_5-1SW-CtuZ81_1qlciIPMP9gwYxgL2VQ_zLTLtM=", true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}

	appdomain1 := AppDomaiN{
		AppDomainID: "a9e051ca-ab7f-4b92-bf03-42f2ed24b7bc",
		AppDomain:   "*.app-azure-eastus2-1623834688947.myeventarena.com",
	}
	var DomainList []AppDomaiN
	DomainList = append(DomainList, appdomain1)
	ttl1 := Ttl{
		Enabled:  true,
		Time:     24,
		TimeUnit: "hours",
	}

	settings1 := Settings{
		UserIDIgnored:         false,
		QueryParamaterIgnored: false,
		Type:                  "static-short-lived",
		CacheControlIgnored:   false,
		CookieIgnored:         false,
		TTL:                   ttl1,
	}

	path1 := PATH{
		Path:      "demo1",
		ByPassURI: false,
		IsDefault: false,
		Status:    "new",
		IsNewPath: true,
		Settings:  settings1,
	}
	var pathList []PATH
	pathList = append(pathList, path1)

	CacheRule1 := &CacheRule{
		ID:                  "026a3c55-cc7b-4eff-9712-6e6a0cf0ed89",
		BypassCache:         false,
		Name:                "abcd",
		CacheControlIgnored: false,
		ShareStaticContent:  false,
		DefaultCache:        false,
		Editable:            true,
		IsNew:               true,
		PathPatterns:        pathList,
		AppDomains:          DomainList,
	}
	_ = CacheRule1

	err1 := prosimo_client.DeleteCacheRule(ctx, "026a3c55-cc7b-4eff-9712-6e6a0cf0ed89")
	if err1 != nil {
		t.Fatal(err)
	}
}
	