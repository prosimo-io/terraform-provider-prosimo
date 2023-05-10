package client

import (
	"context"
	"errors"
	"fmt"
)

type Service_Insertion struct {
	Name               string    `json:"name,omitempty"`
	Type               string    `json:"type,omitempty"`
	Service_name       string    `json:"serviceName,omitempty"`
	RegionID           string    `json:"regionId,omitempty"`
	ID                 string    `json:"id,omitempty"`
	Source             *Source   `json:"source,omitempty"`
	Target             *Target   `json:"target,omitempty"`
	IpRules            *[]IpRule `json:"ipRules,omitempty"`
	Status             string    `json:"status,omitempty"`
	CloudType          string    `json:"cloudType,omitempty"`
	CloudRegion        string    `json:"cloudRegion,omitempty"`
	ServiceID          string    `json:"serviceId,omitempty"`
	GwLoadbalancerID   string    `json:"gwLoadbalancerID,omitempty"`
	SharedServiceCreds string    `json:"sharedServiceCreds,omitempty"`
}

type Service_Insertion_Res struct {
	ID string `json:"id,omitempty"`
}

type Service_Input struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
type Target struct {
	Networks []Service_Input `json:"networks,omitempty"`
	Apps     []Service_Input `json:"apps,omitempty"`
}

type Source struct {
	Networks []Service_Input `json:"networks,omitempty"`
}

type IpRule struct {
	SrcAddr  []string `json:"srcAddr,omitempty"`
	SrcPort  []string `json:"srcPort,omitempty"`
	DestAddr []string `json:"destAddr,omitempty"`
	DestPort []string `json:"destPort,omitempty"`
	Protocol []string `json:"protocol,omitempty"`
}

type SI_Response struct {
	Service_Insertion_Response *Service_Insertion_Res `json:"data,omitempty"`
}
type SIListData struct {
	Records    []*Service_Insertion `json:"records,omitempty"`
	TotalCount int                  `json:"totalCount,omitempty"`
}

type SIListDataResponse struct {
	Data *SIListData `json:"data,omitempty"`
}

func (prosimoClient *ProsimoClient) CreateServiceInsertion(ctx context.Context, serviceInsertionInput *Service_Insertion) (*SI_Response, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", ServiceInsertionEndpoint, serviceInsertionInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &SI_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) UpdateServiceInsertion(ctx context.Context, serviceInsertionInput *Service_Insertion) (*SI_Response, error) {

	updateServiceInsertionEndpoint := fmt.Sprintf("%s/%s", ServiceInsertionEndpoint, serviceInsertionInput.ID)

	req, err := prosimoClient.api_client.NewRequest("PUT", updateServiceInsertionEndpoint, serviceInsertionInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &SI_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}


func (prosimoClient *ProsimoClient) DeleteServiceInsertion(ctx context.Context, serviceIinsertionID string) (*SI_Response, error) {

	updateServiceInsertionEndpoint := fmt.Sprintf("%s/%s", ServiceInsertionEndpoint, serviceIinsertionID)

	req, err := prosimoClient.api_client.NewRequest("DELETE", updateServiceInsertionEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &SI_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) GetServiceInsertion(ctx context.Context) ([]*Service_Insertion, error) {

	SISearchInput := Service_Insertion{}
	req, err := prosimoClient.api_client.NewRequest("POST", GetServiceInsertionEndpoint, SISearchInput)
	if err != nil {
		return nil, err
	}

	siListData := &SIListDataResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, siListData)
	if err != nil {
		return nil, err
	}

	return siListData.Data.Records, nil

}

func (prosimoClient *ProsimoClient) GetServiceInsertionByID(ctx context.Context, id string) (*Service_Insertion, error) {

	siList, err := prosimoClient.GetServiceInsertion(ctx)
	if err != nil {
		return nil, err
	}
	var serviceInsertion *Service_Insertion
	for _, returnedSI := range siList {
		if returnedSI.ID == id {
			serviceInsertion = returnedSI
			break
		}
	}

	if serviceInsertion == nil {
		return nil, errors.New("Service Insertion doesn't exists")
	}

	return serviceInsertion, nil

}
func (prosimoClient *ProsimoClient) GetServiceInsertionByName(ctx context.Context, name string) (*Service_Insertion, error) {

	siList, err := prosimoClient.GetServiceInsertion(ctx)
	if err != nil {
		return nil, err
	}
	var serviceInsertion *Service_Insertion
	for _, returnedSI := range siList {
		if returnedSI.Name == name {
			serviceInsertion = returnedSI
			break
		}
	}

	if serviceInsertion == nil {
		return nil, errors.New("Service Insertion doesn't exists")
	}

	return serviceInsertion, nil

}
