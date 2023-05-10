package client

import (
	"context"
	"errors"
	"fmt"
)

type PL_Map struct {
	Source      *Service_Input  `json:"source,omitempty"`
	Target      *Service_Input  `json:"target,omitempty"`
	Region      string          `json:"region,omitempty"`
	ID          string          `json:"id,omitempty"`
	HostedZones *[]Hosted_Zones `json:"hostedZones,omitempty"`
}

type Hosted_Zones struct {
	ID       string `json:"id,omitempty"`
	SourceID string `json:"sourceID,omitempty"`
	DomainID string `json:"domainID,omitempty"`
	Name     string `json:"name,omitempty"`
}

type PL_Map_Res struct {
	PlMapRes *PlMapReturn `json:"data,omitempty"`
}

type PlMapReturn struct {
	ID     string `json:"id,omitempty"`
	TaskID string `json:"taskID,omitempty"`
}

type PLMapListData struct {
	Records    []*PL_Map `json:"records,omitempty"`
	TotalCount int       `json:"totalCount,omitempty"`
}

type PLMapListDataResponse struct {
	Data *PLMapListData `json:"data,omitempty"`
}

func (prosimoClient *ProsimoClient) CreatePrivateLinkMapping(ctx context.Context, privateLinkMappingInput *PL_Map) (*PL_Map_Res, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", PrivateLinkMappingEndpoint, privateLinkMappingInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &PL_Map_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) UpdatePrivateLinkMapping(ctx context.Context, privateLinkMappingInput *PL_Map) (*PL_Map_Res, error) {

	updatePrivateLinkMappingEndpoint := fmt.Sprintf("%s/%s", PrivateLinkMappingEndpoint, privateLinkMappingInput.ID)

	req, err := prosimoClient.api_client.NewRequest("PUT", updatePrivateLinkMappingEndpoint, privateLinkMappingInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &PL_Map_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) DeletePrivateLinkMapping(ctx context.Context, privateLinkMappingID string) (*PL_Map_Res, error) {

	updatePrivateLinkMappingEndpoint := fmt.Sprintf("%s/%s", PrivateLinkMappingEndpoint, privateLinkMappingID)

	req, err := prosimoClient.api_client.NewRequest("DELETE", updatePrivateLinkMappingEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &PL_Map_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) GetPrivateLinkMapping(ctx context.Context) ([]*PL_Map, error) {

	PLMSearchInput := PL_Map{}
	req, err := prosimoClient.api_client.NewRequest("POST", GetPrivateLinkMappingEndpoint, PLMSearchInput)
	if err != nil {
		return nil, err
	}

	plsListData := &PLMapListDataResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, plsListData)
	if err != nil {
		return nil, err
	}

	return plsListData.Data.Records, nil

}
func (prosimoClient *ProsimoClient) GetPrivateLinkMappingByID(ctx context.Context, id string) (*PL_Map, error) {

	plmList, err := prosimoClient.GetPrivateLinkMapping(ctx)
	if err != nil {
		return nil, err
	}
	var privateLinkMapping *PL_Map
	for _, returnedPLM := range plmList {
		if returnedPLM.ID == id {
			privateLinkMapping = returnedPLM
		}
	}

	if privateLinkMapping == nil {
		return nil, errors.New("private Link Mapping doesn't exists")
	}

	return privateLinkMapping, nil

}
