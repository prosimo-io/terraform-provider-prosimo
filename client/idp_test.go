package client

import (
	"context"
	"fmt"
	"testing"
)

func TestCreateIDP(t *testing.T) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient("https://myeveroot1611940242341.dashboard.psonar.us/", "NKcXJrqknhfivyP214TY9f5uOLvcF46wMEV1SgXwdNM=", true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}

	idpDetails := &IDPDetails{
		APIToken: "00zxRi1hF1vwwHkJxlTU15jyDu11M_SDfbbeosZX_2",
	}

	idp := &IDP{
		Auth_Type:   "oidc",
		IDPName:     "okta",
		AccountURL:  "https://dev-142456.okta.com",
		Select_Type: "primary",
		Details:     *idpDetails,
	}

	respString, err := prosimo_client.CreateIDP(ctx, idp)

	fmt.Printf("respString2 - %s", respString.ResourceData.ID)

	if err != nil {
		fmt.Printf("err - %s", err)
	}

}

func TestGetIDP(t *testing.T) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient("https://myeveroot1611940242341.dashboard.psonar.us/", "NKcXJrqknhfivyP214TY9f5uOLvcF46wMEV1SgXwdNM=", true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}

	getIDP, err := prosimo_client.GetIDP(ctx)

	fmt.Printf("respString2 - %d       ", len(getIDP.IDPs))

	if err != nil {
		fmt.Printf("err - %s", err)
	}

}
