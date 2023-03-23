package client

import (
	"context"
	"fmt"
	"testing"
)

func TestCreateIPPool(t *testing.T) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient("https://myevekuldeep1610576711105.dashboard.psonar.us/", "aggS7fie4O2mulcuEh2XbGH0cjJ4R1WHIhXggsnxr-I=", true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}

	ipPool := &IPPool{
		CloudType: "GCP",
		Cidr:      "10.0.0.0/16",
	}
	respString, err := prosimo_client.CreateIPPool(ctx, ipPool)

	fmt.Printf("respString2 - %s", respString.IPPool.ID)

	if err != nil {
		fmt.Printf("err - %s", err)
	}

}

func TestGetIPPool(t *testing.T) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient("https://myevekuldeep1610576711105.dashboard.psonar.us/", "aggS7fie4O2mulcuEh2XbGH0cjJ4R1WHIhXggsnxr-I=", true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}

	getIPPool, err := prosimo_client.GetIPPool(ctx)

	fmt.Printf("respString2 - %d       ", len(getIPPool.IPPools))

	if err != nil {
		fmt.Printf("err - %s", err)
	}

}

