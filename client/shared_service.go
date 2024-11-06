package client

import (
	"context"
	"errors"
	"fmt"
	"log"
)

type Shared_Service struct {
	ID            string           `json:"id,omitempty"`
	Name          string           `json:"name,omitempty"`
	Deployed      bool             `json:"deployed,omitempty"`
	Status        string           `json:"status,omitempty"`
	Type          string           `json:"type,omitempty"`
	Region        *Region          `json:"region,omitempty"`
	ServiceInsert *[]ServiceInsert `json:"serviceInsert,omitempty"`
	CreatedTime   string           `json:"createdTime,omitempty"`
	UpdatedTime   string           `json:"updatedTime,omitempty"`
	Progress      int              `json:"progress,omitempty"`
	TaskID        string           `json:"taskID,omitempty"`
	TeamID        string           `json:"teamID,omitempty"`
}

type Region struct {
	ID               string   `json:"id,omitempty"`
	CloudRegion      string   `json:"cloudRegion,omitempty"`
	GwLoadBalancerID string   `json:"gwLoadBalancerID,omitempty"`
	CloudKeyID       string   `json:"cloudKeyID,omitempty"`
	CloudType        string   `json:"cloudType,omitempty"`
	CloudZones       []string `json:"cloudZones,omitempty"`
	ResourceGrp      string   `json:"resourceGroup,omitempty"`
}

type ServiceInsert struct {
	ID           string  `json:"id,omitempty"`
	Name         string  `json:"name,omitempty"`
	Type         string  `json:"type,omitempty"`
	RegionID     string  `json:"regionid,omitempty"`
	Source       *Source `json:"source,omitempty"`
	Target       *Target `json:"target,omitempty"`
	NamespaceNID int     `json:"namespaceNID,omitempty"`
	NamespaceID  string  `json:"namespaceID,omitempty"`
}

type SS_Response struct {
	Shared_Service_Response *Shared_Service `json:"data,omitempty"`
}
type SSListData struct {
	Records    []*Shared_Service `json:"records,omitempty"`
	TotalCount int               `json:"totalCount,omitempty"`
}

type SSListDataResponse struct {
	Data *SSListData `json:"data,omitempty"`
}
type GwlbSearch struct {
	CloudCredID string `json:"cloudCreds,omitempty"`
	Region string `json:"region,omitempty"`
}
type GWLBlistDataResponse struct {
	Data []GWLBListData `json:"data,omitempty"`
}
type GWLBListData struct {
	Name string `json:"name,omitempty"`
}
type GWLBAWSlistDataResponse struct {
	Data []string `json:"data,omitempty"`
}

func (prosimoClient *ProsimoClient) CreateSharedService(ctx context.Context, sharedServiceInput *Shared_Service) (*SS_Response, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", SharedServiceEndpoint, sharedServiceInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &SS_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) UpdateSharedService(ctx context.Context, sharedServiceInput *Shared_Service) (*SS_Response, error) {

	updateSharedServiceEndpoint := fmt.Sprintf("%s/%s", SharedServiceEndpoint, sharedServiceInput.ID)

	req, err := prosimoClient.api_client.NewRequest("PUT", updateSharedServiceEndpoint, sharedServiceInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &SS_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) DecomSharedService(ctx context.Context, sharedServiceID string) (*SS_Response, error) {

	deleteSSEndpt := fmt.Sprintf("%s/%s", SharedServiceDeploymentEndpoint, sharedServiceID)

	req, err := prosimoClient.api_client.NewRequest("DELETE", deleteSSEndpt, nil)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &SS_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) OnboardSharedService(ctx context.Context, sharedServiceID string) (*SS_Response, error) {

	ssDeployEndpt := fmt.Sprintf("%s/%s", SharedServiceDeploymentEndpoint, sharedServiceID)

	emptyInterface := &Shared_Service{}
	req, err := prosimoClient.api_client.NewRequest("PUT", ssDeployEndpt, emptyInterface)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &SS_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) DeleteSharedService(ctx context.Context, sharedServiceID string) error {

	updateSharedServiceEndpoint := fmt.Sprintf("%s/%s", SharedServiceEndpoint, sharedServiceID)

	req, err := prosimoClient.api_client.NewRequest("DELETE", updateSharedServiceEndpoint, nil)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}

func (prosimoClient *ProsimoClient) GetSharedService(ctx context.Context) ([]*Shared_Service, error) {

	SSSearchInput := Shared_Service{}
	req, err := prosimoClient.api_client.NewRequest("POST", GetSharedServiceEndpoint, SSSearchInput)
	if err != nil {
		return nil, err
	}

	ssListData := &SSListDataResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, ssListData)
	if err != nil {
		return nil, err
	}

	return ssListData.Data.Records, nil

}

func (prosimoClient *ProsimoClient) GetSharedServiceByID(ctx context.Context, id string) (*Shared_Service, error) {

	ssList, err := prosimoClient.GetSharedService(ctx)
	if err != nil {
		return nil, err
	}
	var sharedService *Shared_Service
	for _, returnedSS := range ssList {
		if returnedSS.ID == id {
			sharedService = returnedSS
			break
		}
	}

	if sharedService == nil {
		return nil, errors.New("Shared Service doesn't exists")
	}

	return sharedService, nil

}

func (prosimoClient *ProsimoClient) GetSharedServiceByName(ctx context.Context, name string) (*Shared_Service, error) {

	ssList, err := prosimoClient.GetSharedService(ctx)
	if err != nil {
		return nil, err
	}
	var sharedService *Shared_Service
	for _, returnedSS := range ssList {
		if returnedSS.Name == name {
			sharedService = returnedSS
			break
		}
	}

	if sharedService == nil {
		return nil, errors.New("Shared Service doesn't exists")
	}

	return sharedService, nil

}

func (prosimoClient *ProsimoClient) CheckIfGWLBExistsAWS(ctx context.Context, gwinput GwlbSearch, gwlb string) (bool, error) {
	// Initialize the flag as false
	gwlbExists := false

	// Create a new API request for AWS endpoint
	req, err := prosimoClient.api_client.NewRequest("POST", AWSGWLBEndpoint, gwinput)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	// Initialize the response structure
	GWListData := &GWLBAWSlistDataResponse{}

	// Send the request and unmarshal the response into GWListData
	_, err = prosimoClient.api_client.Do(ctx, req, GWListData)
	if err != nil {
		return false, fmt.Errorf("failed to perform request: %w", err)
	}
	log.Println("GWListData", GWListData)
	log.Println("gwlb", gwlb)
	// Loop through the returned list of GWLB names
	for _, retgwlb := range GWListData.Data {
		if retgwlb == gwlb {
			gwlbExists = true
			break
		}
	}

	// Return whether the GWLB exists along with no error
	return gwlbExists, nil
}


func (prosimoClient *ProsimoClient) CheckIfGWLBExistsGCP(ctx context.Context, gwinput GwlbSearch, gwlb string) (bool, error) {
	// Initialize the flag as false
	gwlbExists := false

	// Create a new API request for GCP endpoint
	req, err := prosimoClient.api_client.NewRequest("POST", GCPGWLBEndpoint, gwinput)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	// Initialize the response structure
	GWListData := &GWLBlistDataResponse{}

	// Send the request and unmarshal the response into GWListData
	_, err = prosimoClient.api_client.Do(ctx, req, GWListData)
	if err != nil {
		return false, fmt.Errorf("failed to perform request: %w", err)
	}

	// Loop through the returned data to check if the desired GWLB exists
	for _, retgwlb := range GWListData.Data {
		if retgwlb.Name == gwlb {
			gwlbExists = true
			break
		}
	}

	// Return whether the GWLB exists along with no error
	return gwlbExists, nil
}


func (prosimoClient *ProsimoClient) CheckIfGWLBExistsAZURE(ctx context.Context, gwinput GwlbSearch, gwlb string, rg string) (bool, error) {
	// Initialize the flag as false
	gwlbExists := false

	// Construct the Azure GWLB API endpoint URL
	postAzureGWLBEndpoint := fmt.Sprintf(AzureGWLBEndpoint, rg)

	// Create a new API request with the given input data
	req, err := prosimoClient.api_client.NewRequest("POST", postAzureGWLBEndpoint, gwinput)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	// Initialize the response structure
	GWListData := &GWLBlistDataResponse{}

	// Send the request and unmarshal the response into GWListData
	_, err = prosimoClient.api_client.Do(ctx, req, GWListData)
	if err != nil {
		return false, fmt.Errorf("failed to perform request: %w", err)
	}

	// Loop through the returned data to check if the desired GWLB exists
	for _, retgwlb := range GWListData.Data {
		if retgwlb.Name == gwlb {
			gwlbExists = true
			break
		}
	}

	// Return whether the GWLB exists along with no error
	return gwlbExists, nil
}
