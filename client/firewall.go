package client

import (
	"context"
	"errors"
	"fmt"
)

type FWConfig struct {
	ID               string            `json:"id,omitempty"`
	TaskID           string            `json:"taskID,omitempty"`
	Status           string            `json:"status,omitempty"`
	Type             string            `json:"type,omitempty"`
	Name             string            `json:"name,omitempty"`
	CloudKeyID       string            `json:"cloudKeyID,omitempty"`
	CloudRegion      string            `json:"cloudRegion,omitempty"`
	CIDR             string            `json:"cidr,omitempty"`
	Version          string            `json:"version,omitempty"`
	InstanceSize     string            `json:"instanceSize,omitempty"`
	ActivationConfig *ActivationConfig `json:"activationConfig,omitempty"`
	ScalingConfig    *ScalingConfig    `json:"scalingConfig,omitempty"`
	DeviceConfig     *DeviceConfig     `json:"deviceConfig,omitempty"`
	AccessConfig     *AccessConfig     `json:"accessConfig,omitempty"`
	Bootstrap        string            `json:"bootstrap,omitempty"`
}

type ActivationConfig struct {
	AuthKey  string `json:"authKey,omitempty"`
	AuthCode string `json:"authCode,omitempty"`
}
type ScalingConfig struct {
	DefaultCapacity int `json:"defaultCapacity,omitempty"`
	MinCapacity     int `json:"minCapacity,omitempty"`
	MaxCapacity     int `json:"maxCapacity,omitempty"`
}

type DeviceConfig struct {
	DGName       string `json:"dgName,omitempty"`
	TPLStackname string `json:"tplStackname,omitempty"`
}
type AccessConfig struct {
	PEMDetails  *PEMDetails `json:"pemDetails,omitempty"`
	AccessCreds *AccessCred `json:"accessCreds,omitempty"`
}

type AccessCred struct {
	UserName string `json:"username,omitempty"`
	PassWord string `json:"password,omitempty"`
}

type PEMDetails struct {
	// KeypairName string `json:"keypairName,omitempty"`
	KeyGenerate bool   `json:"keyGenerate,omitempty"`
	PublicKey   string `json:"publicKey,omitempty"`
}

// type Hosted_Zones struct {
// 	ID       string `json:"id,omitempty"`
// 	SourceID string `json:"sourceID,omitempty"`
// 	DomainID string `json:"domainID,omitempty"`
// 	Name     string `json:"name,omitempty"`
// }

type FWConfig_Res struct {
	PlMapRes *FWConfig `json:"data,omitempty"`
}

type FWConfigReturn struct {
	FWConfig *FWConfig `json:"data,omitempty"`
}

type FWConfigData struct {
	Records    []*FWConfig `json:"records,omitempty"`
	TotalCount int         `json:"totalCount,omitempty"`
}

type FWConfigListDataResponse struct {
	Data *FWConfigData `json:"data,omitempty"`
}

func (prosimoClient *ProsimoClient) CreateFirewall(ctx context.Context, firewallInput *FWConfig) (*FWConfig_Res, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", CreateFirewallEndpoint, firewallInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &FWConfig_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}
func (prosimoClient *ProsimoClient) DeployFirewall(ctx context.Context, firewallID string) (*FWConfigReturn, error) {
	deployFirewallEndpointUpdated := fmt.Sprintf(DeployFirewallEndpoint, firewallID)

	req, err := prosimoClient.api_client.NewRequest("PUT", deployFirewallEndpointUpdated, firewallID)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &FWConfigReturn{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}
func (prosimoClient *ProsimoClient) DecomFirewall(ctx context.Context, firewallID string) (*FWConfigReturn, error) {
	DecommissionFirewallEndpointUpdated := fmt.Sprintf(DecommissionFirewallEndpoint, firewallID)

	req, err := prosimoClient.api_client.NewRequest("DELETE", DecommissionFirewallEndpointUpdated, firewallID)
	if err != nil {
		return nil, err
	}

	resourceResponseData := &FWConfigReturn{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourceResponseData)
	if err != nil {
		return nil, err
	}

	return resourceResponseData, nil

}

func (prosimoClient *ProsimoClient) UpdateFirewall(ctx context.Context, firewallInput *FWConfig) (*FWConfig_Res, error) {
	// deployFirewallEndpointUpdated := fmt.Sprintf(DeployFirewallEndpoint, networkPrefixID)

	updateFirewallEndpoint := fmt.Sprintf("%s/%s", CreateFirewallEndpoint, firewallInput.ID)

	req, err := prosimoClient.api_client.NewRequest("PUT", updateFirewallEndpoint, firewallInput)
	if err != nil {
		return nil, err
	}

	resourceResponseData := &FWConfig_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourceResponseData)
	if err != nil {
		return nil, err
	}

	return resourceResponseData, nil

}

func (prosimoClient *ProsimoClient) DeleteFirewall(ctx context.Context, firewallID string) error {

	updateFirewallEndpoint := fmt.Sprintf("%s/%s", CreateFirewallEndpoint, firewallID)

	req, err := prosimoClient.api_client.NewRequest("DELETE", updateFirewallEndpoint, nil)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}

func (prosimoClient *ProsimoClient) GetFirewall(ctx context.Context) ([]*FWConfig, error) {

	FirewallSearchInput := FWConfig{}
	req, err := prosimoClient.api_client.NewRequest("POST", SearchFirewallEndpoint, FirewallSearchInput)
	if err != nil {
		return nil, err
	}

	plsListData := &FWConfigListDataResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, plsListData)
	if err != nil {
		return nil, err
	}

	return plsListData.Data.Records, nil

}
func (prosimoClient *ProsimoClient) GetFirewallByID(ctx context.Context, id string) (*FWConfig, error) {

	firewallList, err := prosimoClient.GetFirewall(ctx)
	if err != nil {
		return nil, err
	}
	var firewall *FWConfig
	for _, returnedFirewall := range firewallList {
		if returnedFirewall.ID == id {
			firewall = returnedFirewall
		}
	}

	if firewall == nil {
		return nil, errors.New("Firewall doesn't exists")
	}

	return firewall, nil

}
