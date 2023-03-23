package client

import (
	"context"
	"fmt"
)

type IPPoolList struct {
	IPPools []*IPPool `json:"data,omitempty"`
}

type IPPoolData struct {
	IPPool *IPPool `json:"data,omitempty"`
}

type IPPool struct {
	Cidr         string `json:"cidr,omitempty"`
	CloudType    string `json:"cloudType,omitempty"`
	Name         string `json:"name,omitempty"`
	ID           string `json:"id,omitempty"`
	TotalSubnets int    `json:"totalSubnets,omitempty"`
	SubnetsInUse int    `json:"subnetsInUse,omitempty"`
}

func (ipPool IPPool) String() string {
	return fmt.Sprintf("{Cidr:%s, CloudType:%s, ID:%s}", ipPool.Cidr, ipPool.CloudType, ipPool.ID)
}

func (prosimoClient *ProsimoClient) CreateIPPool(ctx context.Context, ipPoolOpts *IPPool) (*IPPoolData, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", IPPoolEndpoint, ipPoolOpts)
	if err != nil {
		return nil, err
	}

	ipPool := &IPPoolData{}
	_, err = prosimoClient.api_client.Do(ctx, req, ipPool)
	if err != nil {
		return nil, err
	}

	return ipPool, nil

}

func (prosimoClient *ProsimoClient) GetIPPool(ctx context.Context) (*IPPoolList, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", IPPoolEndpoint, nil)
	if err != nil {
		return nil, err
	}

	ipPoolList := &IPPoolList{}
	_, err = prosimoClient.api_client.Do(ctx, req, ipPoolList)
	if err != nil {
		return nil, err
	}

	return ipPoolList, nil

}

func (prosimoClient *ProsimoClient) GetIPPoolFiltered(ctx context.Context, id string) (*IPPool, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", IPPoolEndpoint, nil)
	if err != nil {
		return nil, err
	}

	ipPoolList := &IPPoolList{}
	_, err = prosimoClient.api_client.Do(ctx, req, ipPoolList)
	if err != nil {
		return nil, err
	}
	for _, v := range ipPoolList.IPPools {
		if v.ID == id {
			return v, nil
		}
	}
	return nil, err

}

func (prosimoClient *ProsimoClient) DeleteIPPool(ctx context.Context, ipPoolID string) error {

	deleteIPPoolEndpt := fmt.Sprintf("%s/%s", IPPoolEndpoint, ipPoolID)

	req, err := prosimoClient.api_client.NewRequest("DELETE", deleteIPPoolEndpt, nil)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}
