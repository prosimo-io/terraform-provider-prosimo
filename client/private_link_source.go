package client

import (
	"context"
	"errors"
	"fmt"
)

type PL_Source struct {
	Name         string          `json:"name,omitempty"`
	Region       string          `json:"region,omitempty"`
	ID           string          `json:"id,omitempty"`
	Status       string          `json:"status,omitempty"`
	CloudCredsID string          `json:"cloudCredsID,omitempty"`
	CloudSources *[]Cloud_Source `json:"cloudSources,omitempty"`
}

type Cloud_Source struct {
	ID           string         `json:"id,omitempty"`
	CloudNetwork *Cloud_Network `json:"cloudNetwork,omitempty"`
	Subnets      *[]Subnet      `json:"subnets,omitempty"`
}

type Subnet struct {
	Cidr               string `json:"cidr,omitempty"`
	OnboardedASNetwork bool   `json:"onboardedAsNetwork,omitempty"`
}

type Subnet_res struct{
	Subnets []Subnet `json:"data,omitempty"`
}
type PL_Source_ResID struct {
	ID string `json:"id,omitempty"`
}

type Cloud_Network struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Cloud_Network_res struct{
	CloudNetworks []Cloud_Network `json:"data,omitempty"`
}
type PL_Source_Response struct {
	PLSource_ResID *PL_Source_ResID `json:"data,omitempty"`
}

type PLSListData struct {
	Records    []*PL_Source `json:"records,omitempty"`
	TotalCount int          `json:"totalCount,omitempty"`
}

type PLSListDataResponse struct {
	Data *PLSListData `json:"data,omitempty"`
}
type Search_Input struct {
	ID           string `json:"id,omitempty"`
	CloudCredsID string `json:"cloudCredsID,omitempty"`
	Region       string `json:"region,omitempty"`
}

func (prosimoClient *ProsimoClient) CreatePrivateLinkSource(ctx context.Context, privateLinkSourceInput *PL_Source) (*PL_Source_Response, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", PrivateLinkSourceEndpoint, privateLinkSourceInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &PL_Source_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) UpdatePrivateLinkSource(ctx context.Context, privateLinkSourceInput *PL_Source) (*PL_Source_Response, error) {

	updatePrivateLinkSourceEndpoint := fmt.Sprintf("%s/%s", PrivateLinkSourceEndpoint, privateLinkSourceInput.ID)

	req, err := prosimoClient.api_client.NewRequest("PUT", updatePrivateLinkSourceEndpoint, privateLinkSourceInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &PL_Source_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) DeletePrivateLinkSource(ctx context.Context, privateLinkSourceID string) (*PL_Source_Response, error) {

	updatePrivateLinkSourceEndpoint := fmt.Sprintf("%s/%s", PrivateLinkSourceEndpoint, privateLinkSourceID)

	req, err := prosimoClient.api_client.NewRequest("DELETE", updatePrivateLinkSourceEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &PL_Source_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) GetPrivateLinkSource(ctx context.Context) ([]*PL_Source, error) {

	PLSSearchInput := PL_Source{}
	req, err := prosimoClient.api_client.NewRequest("POST", GetPrivateLinkSourceEndpoint, PLSSearchInput)
	if err != nil {
		return nil, err
	}

	plsListData := &PLSListDataResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, plsListData)
	if err != nil {
		return nil, err
	}

	return plsListData.Data.Records, nil

}

func (prosimoClient *ProsimoClient) GetPrivateLinkSourceByID(ctx context.Context, id string) (*PL_Source, error) {

	plsList, err := prosimoClient.GetPrivateLinkSource(ctx)
	if err != nil {
		return nil, err
	}
	var privateLinkSource *PL_Source
	for _, returnedPLS := range plsList {
		if returnedPLS.ID == id {
			privateLinkSource = returnedPLS
		}
	}

	if privateLinkSource == nil {
		return nil, errors.New("private Link Source doesn't exists")
	}

	return privateLinkSource, nil

}
func (prosimoClient *ProsimoClient) GetPrivateLinkSourceByName(ctx context.Context, name string) (*PL_Source, error) {

	plsList, err := prosimoClient.GetPrivateLinkSource(ctx)
	if err != nil {
		return nil, err
	}
	var privateLinkSource *PL_Source
	for _, returnedPLS := range plsList {
		if returnedPLS.Name == name {
			privateLinkSource = returnedPLS
		}
	}

	if privateLinkSource == nil {
		return nil, errors.New("private Link Source doesn't exists")
	}

	return privateLinkSource, nil

}

func (prosimoClient *ProsimoClient) DiscoverPVSNetworks(ctx context.Context, discoverNetworkInput *Search_Input) (*Cloud_Network_res, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", DiscoverPVSNetworkEndpoint, discoverNetworkInput)
	if err != nil {
		return nil, err
	}

	discoveredNetwork := &Cloud_Network_res{}
	_, err = prosimoClient.api_client.Do(ctx, req, discoveredNetwork)
	if err != nil {
		return nil, err
	}

	return discoveredNetwork, nil

}

func (prosimoClient *ProsimoClient) DiscoverPVSSubnets(ctx context.Context, discoverSubnetInput *Search_Input) (*Subnet_res, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", DiscoverPVSSubnetEndpoint, discoverSubnetInput)
	if err != nil {
		return nil, err
	}

	discoveredSubnets := &Subnet_res{}
	_, err = prosimoClient.api_client.Do(ctx, req, discoveredSubnets)
	if err != nil {
		return nil, err
	}

	return discoveredSubnets, nil

}
