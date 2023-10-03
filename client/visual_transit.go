package client

import "context"

type Visual_Transit_Setup struct {
	CloudType   string      `json:"cloudType,omitempty"`
	CloudRegion string      `json:"cloudRegion,omitempty"`
	ID          string      `json:"id,omitempty"`
	TaskID      string      `json:"taskID,omitempty"`
	Status      string      `json:"status,omitempty"`
	TeamID      string      `json:"teamID,omitempty"`
	AuditID     string      `json:"auditID,omitempty"`
	HasError    bool        `json:"hasError,omitempty"`
	Edge        *EdGe       `json:"edge,omitempty"`
	Operation   *Operation  `json:"operation,omitempty"`
	Deployment  *Deployment `json:"deployment,omitempty"`
	Account     *Account    `json:"account,omitempty"`
}

type Account struct {
	CloudKeyID  string        `json:"cloudKeyID,omitempty"`
	AccountID   string        `json:"accountID,omitempty"`
	AccountName string        `json:"accountName,omitempty"`
	VPCS        []*Constructs `json:"vpcs,omitempty"`
	VNETS       []*Constructs `json:"vnets,omitempty"`
}

type EdGe struct {
	ID             string `json:"id,omitempty"`
	Status         string `json:"status,omitempty"`
	CloudNetworkID string `json:"cloudNetworkID,omitempty"`
	EdgeFqdn       string `json:"edgeFqdn,omitempty"`
	Subnet         string `json:"subnet,omitempty"`
	AccountID      string `json:"accountID,omitempty"`
	IssueFixed     bool   `json:"issueFixed,omitempty"`
	Type           string `json:"type,omitempty"`
}
type Connection struct {
	ID         string   `json:"id,omitempty"`
	Name       string   `json:"name,omitempty"`
	AccountID  string   `json:"accountID,omitempty"`
	Action     string   `json:"action,omitempty"`
	Issue      string   `json:"issue,omitempty"`
	IssueFixed bool     `json:"issueFixed,omitempty"`
	Subnets    []string `json:"subnets,omitempty"`
	Type       string   `json:"type,omitempty"`
	PathAction string   `json:"pathAction,omitempty"`
}
type Constructs struct {
	ID           string       `json:"id,omitempty"`
	Name         string       `json:"name,omitempty"`
	AccountID    string       `json:"accountID,omitempty"`
	Status       string       `json:"status,omitempty"`
	Action       string       `json:"action,omitempty"`
	Subnets      []string     `json:"subnets,omitempty"`
	Connections  []Connection `json:"connections,omitempty"`
	Type         string       `json:"type,omitempty"`
	AddressSpace string       `json:"addressSpace,omitempty"`
	VwanID       string       `json:"vwanID,omitempty"`
}
type Operation struct {
	TGWS  []*Constructs `json:"tgws,omitempty"`
	VPCS  []*Constructs `json:"vpcs,omitempty"`
	VHUBS []*Constructs `json:"vhubs,omitempty"`
	VNETS []*Constructs `json:"vnets,omitempty"`
}
type Deployment struct {
	TGWS  []Constructs `json:"tgws,omitempty"`
	VPCS  []Constructs `json:"vpcs,omitempty"`
	VHUBS []Constructs `json:"vhubs,omitempty"`
	VNETS []Constructs `json:"vnets,omitempty"`
}

type TransitSetupSearchResponse struct {
	Data []*Visual_Transit_Setup `json:"data,omitempty"`
}

type EdgeInput struct {
	CloudType   string `json:"cloudType,omitempty"`
	CloudRegion string `json:"cloudRegion,omitempty"`
}
type TransitSearchInput struct {
	Edges []EdgeInput `json:"edges,omitempty"`
}

func (prosimoClient *ProsimoClient) TransitSetupSearch(ctx context.Context, transitSearchInput TransitSearchInput) ([]*Visual_Transit_Setup, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", GetTransitSetupEndpoint, transitSearchInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &TransitSetupSearchResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData.Data, nil

}

func (prosimoClient *ProsimoClient) TransitSetupSummary(ctx context.Context, transitSearchInput TransitSearchInput) ([]*Visual_Transit_Setup, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", TransitSetupSummaryEndpoint, transitSearchInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &TransitSetupSearchResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData.Data, nil

}

func (prosimoClient *ProsimoClient) CloudNetworkSearch(ctx context.Context, networkSearchInput EdgeInput) ([]*Visual_Transit_Setup, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", GetNetworkEndpoint, networkSearchInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &TransitSetupSearchResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData.Data, nil

}

func (prosimoClient *ProsimoClient) TransitSetup(ctx context.Context, transitSetupInput []Visual_Transit_Setup) error {

	req, err := prosimoClient.api_client.NewRequest("PUT", CreateTransitSetupEndpoint, transitSetupInput)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}

func (prosimoClient *ProsimoClient) CreateTransitDeploy(ctx context.Context, setupIDInput []string) ([]*Visual_Transit_Setup, error) {

	req, err := prosimoClient.api_client.NewRequest("PUT", CreateTransitDeploymentEndpoint, setupIDInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &TransitSetupSearchResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData.Data, nil

}

func (prosimoClient *ProsimoClient) DeleteTransitSetup(ctx context.Context, setupIDInput []string) error {

	req, err := prosimoClient.api_client.NewRequest("POST", DeleteTransitSetupEndpoint, setupIDInput)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}
