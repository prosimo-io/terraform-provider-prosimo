package client

import (
	"context"
	"fmt"
	"os"
)

func GetIPPoolfiltered(id string) (*IPPool, bool) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient(os.Getenv("PROSIMO_BASE_URL"), os.Getenv("PROSIMO_TOKEN"), true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}

	getIPPool, err := prosimo_client.GetIPPool(ctx)

	if err != nil {
		fmt.Printf("err - %s", err)
	}
	for _, v := range getIPPool.IPPools {
		if v.ID == id {
			return v, false
		}
	}
	return nil, true
}

func GetDynamicRiskfiltered(id string) (*Dynamic_Risk, bool) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient(os.Getenv("PROSIMO_BASE_URL"), os.Getenv("PROSIMO_TOKEN"), true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}

	getDynamicRisk, err := prosimo_client.GetDYRisk(ctx)

	if err != nil {
		fmt.Printf("err - %s", err)
	}
	for _, v := range getDynamicRisk.DyRisk {
		if v.Id == id {
			return v, false
		}
	}
	return nil, true
}

func GetCloudCredfiltered(id string) (*CloudCreds, bool) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient(os.Getenv("PROSIMO_BASE_URL"), os.Getenv("PROSIMO_TOKEN"), true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}

	GetCloudCreds, err := prosimo_client.GetCloudCreds(ctx)

	if err != nil {
		fmt.Printf("err - %s", err)
	}
	for _, v := range GetCloudCreds.CloudCreds {
		if v.ID == id {
			return v, false
		}
	}
	return nil, true
}

func GetAppOnboardingfiltered(id string) (*AppOnboardSettings, bool) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient(os.Getenv("PROSIMO_BASE_URL"), os.Getenv("PROSIMO_TOKEN"), true)
	if err != nil {
		fmt.Errorf("err - %s", err)
	}

	appOnboardSetting, err1 := prosimo_client.GetAppOnboardSettings(ctx, id)

	if err1 != nil {
		// fmt.Printf("err - %s", err)
		return appOnboardSetting, true
	}

	return nil, false
}

func GetNetworkOnboardingfiltered(id string) (*NetworkOnboardoptns, bool) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient(os.Getenv("PROSIMO_BASE_URL"), os.Getenv("PROSIMO_TOKEN"), true)
	if err != nil {
		fmt.Errorf("err - %s", err)
	}

	networkOnboardSetting, err1 := prosimo_client.GetNetworkSettings(ctx, id)

	if err1 != nil {
		// fmt.Printf("err - %s", err)
		return networkOnboardSetting, true
	}

	return nil, false
}

func GetSharedServicefiltered(id string) (*Shared_Service, bool) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient(os.Getenv("PROSIMO_BASE_URL"), os.Getenv("PROSIMO_TOKEN"), true)
	if err != nil {
		fmt.Errorf("err - %s", err)
	}

	sharedService, err1 := prosimo_client.GetSharedServiceByID(ctx, id)

	if err1 != nil {
		// fmt.Printf("err - %s", err)
		return sharedService, true
	}

	return nil, false
}

func GetServiceInsertionfiltered(id string) (*Service_Insertion, bool) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient(os.Getenv("PROSIMO_BASE_URL"), os.Getenv("PROSIMO_TOKEN"), true)
	if err != nil {
		fmt.Errorf("err - %s", err)
	}

    serviceInsertion, err1 := prosimo_client.GetServiceInsertionByID(ctx, id)

	if err1 != nil {
		// fmt.Printf("err - %s", err)
		return serviceInsertion, true
	}

	return nil, false
}

func GetPrivateLinkSourcefiltered(id string) (*PL_Source, bool) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient(os.Getenv("PROSIMO_BASE_URL"), os.Getenv("PROSIMO_TOKEN"), true)
	if err != nil {
		fmt.Errorf("err - %s", err)
	}

    plSource, err1 := prosimo_client.GetPrivateLinkSourceByID(ctx, id)

	if err1 != nil {
		// fmt.Printf("err - %s", err)
		return plSource, true
	}

	return nil, false
}
func GetPrivateLinkMappingfiltered(id string) (*PL_Map, bool) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient(os.Getenv("PROSIMO_BASE_URL"), os.Getenv("PROSIMO_TOKEN"), true)
	if err != nil {
		fmt.Errorf("err - %s", err)
	}

    plMap, err1 := prosimo_client.GetPrivateLinkMappingByID(ctx, id)

	if err1 != nil {
		// fmt.Printf("err - %s", err)
		return plMap, true
	}

	return nil, false
}
func GetPolicyfiltered(id string) (*Policy, bool) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient(os.Getenv("PROSIMO_BASE_URL"), os.Getenv("PROSIMO_TOKEN"), true)
	if err != nil {
		fmt.Errorf("err - %s", err)
	}

    policy, err1 := prosimo_client.GetPolicyByID(ctx, id)

	if err1 != nil {
		// fmt.Printf("err - %s", err)
		return policy, true
	}

	return nil, false
}

func GetWaffiltered(id string) (*Waf, bool) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient(os.Getenv("PROSIMO_BASE_URL"), os.Getenv("PROSIMO_TOKEN"), true)
	if err != nil {
		fmt.Errorf("err - %s", err)
	}

    waf, err1 := prosimo_client.GetWafByID(ctx, id)

	if err1 != nil {
		// fmt.Printf("err - %s", err)
		return waf, true
	}

	return nil, false
}

func GetIDPfiltered(id string) (*IDP, bool) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient(os.Getenv("PROSIMO_BASE_URL"), os.Getenv("PROSIMO_TOKEN"), true)
	if err != nil {
		fmt.Errorf("err - %s", err)
	}

    idp, err1 := prosimo_client.GetIDPByID(ctx, id)

	if err1 != nil {
		// fmt.Printf("err - %s", err)
		return idp, true
	}

	return nil, false
}