package client

import (
	"context"
)

type DiscoveredNetworksList struct {
	StatusCode string                `json:"id"`
	Message    string                `json:"message"`
	Data       []*DiscoveredNetworks `json:"data"`
}

type DiscoveredNetworks struct {
	ID          string    `json:"id,omitempty"`
	CloudType   string    `json:"cloud_type"`
	Name        string    `json:"name"`
	AccountName string    `json:"account_name"`
	Regions     []Regions `json:"regions"`
}

type Regions struct {
	ID           string `json:"id"`
	CloudCredsID string `json:"cloud_creds_id"`
	Name         string `json:"name"`
	VpcCount     int    `json:"vpc_count"`
	Vpcs         []VPCs `json:"vpcs"`
}

type VPCs struct {
	ID          string `json:"id"`
	RegionID    string `json:"region_id"`
	Cidr        string `json:"cidr"`
	Network     string `json:"network"`
	Name        string `json:"name"`
	SubnetCount int    `json:"subnet_count"`
}

func (prosimoClient *ProsimoClient) GetDiscoveredNetworks(ctx context.Context) ([]*DiscoveredNetworks, error) {
	req, err := prosimoClient.api_client.NewRequest("GET", DiscoveredNetworksapi, nil)

	if err != nil {
		return nil, err
	}
	discoveredNetworkList := &DiscoveredNetworksList{}
	_, err = prosimoClient.api_client.Do(ctx, req, discoveredNetworkList)

	if err != nil {
		return nil, err
	}

	return discoveredNetworkList.Data, nil
}
