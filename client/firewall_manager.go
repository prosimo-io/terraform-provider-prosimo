package client

import (
	"context"
	"errors"
	"fmt"
)

type LicenseDetails struct {
	LicenseMode    string `json:"licenseMode,omitempty"`
	FirewallFamily string `json:"firewallFamily,omitempty"`
	InstanceFamily string `json:"instanceFamily,omitempty"`
	LicenseType    string `json:"licenseType,omitempty"`
}

type FWMConfig struct {
	ID              string          `json:"id,omitempty"`
	IntegrationType string          `json:"integrationType,omitempty"`
	IPAddress       string          `json:"ipAddress,omitempty"`
	APIKey          string          `json:"apiKey,omitempty"`
	LicenseDetails  *LicenseDetails `json:"licenseDetails,omitempty"`
}

type FWMConfig_Res struct {
	PlMapRes *FWConfig `json:"data,omitempty"`
}

type FWMConfigReturn struct {
	ID     string `json:"id,omitempty"`
	TaskID string `json:"taskID,omitempty"`
}

type FWMConfigListDataResponse struct {
	Data []*FWMConfig `json:"data,omitempty"`
}

func (prosimoClient *ProsimoClient) CreateFirewallManager(ctx context.Context, firewallInput *FWMConfig) (*FWMConfig_Res, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", CreateFirewallManagerEndpoint, firewallInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &FWMConfig_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) UpdateFirewallManager(ctx context.Context, firewallInput *FWMConfig)  error {
	// deployFirewallEndpointUpdated := fmt.Sprintf(DeployFirewallEndpoint, networkPrefixID)

	updateFirewallManagerEndpoint := fmt.Sprintf("%s/%s", CreateFirewallManagerEndpoint, firewallInput.ID)

	req, err := prosimoClient.api_client.NewRequest("PATCH", updateFirewallManagerEndpoint, firewallInput)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}

func (prosimoClient *ProsimoClient) DeleteFirewallManager(ctx context.Context, firewallManagerID string) error {

	updateFirewallManagerEndpoint := fmt.Sprintf("%s/%s", CreateFirewallManagerEndpoint, firewallManagerID)

	req, err := prosimoClient.api_client.NewRequest("DELETE", updateFirewallManagerEndpoint, nil)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}

func (prosimoClient *ProsimoClient) GetFirewallManager(ctx context.Context) ([]*FWMConfig, error) {

	FirewallSearchInput := FWMConfig{}
	req, err := prosimoClient.api_client.NewRequest("POST", SearchFirewallManagerEndpoint, FirewallSearchInput)
	if err != nil {
		return nil, err
	}

	plsListData := &FWMConfigListDataResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, plsListData)
	if err != nil {
		return nil, err
	}

	return plsListData.Data, nil

}
func (prosimoClient *ProsimoClient) GetFirewallManagerByID(ctx context.Context, id string) (*FWMConfig, error) {

	firewallManagerList, err := prosimoClient.GetFirewallManager(ctx)
	if err != nil {
		return nil, err
	}
	var firewallManager *FWMConfig
	for _, returnedFirewallManager := range firewallManagerList {
		if returnedFirewallManager.ID == id {
			firewallManager = returnedFirewallManager
		}
	}

	if firewallManager == nil {
		return nil, errors.New("firewall manager doesn't exists")
	}

	return firewallManager, nil

}
