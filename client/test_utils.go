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
