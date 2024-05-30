package client

import (
	"context"
	"errors"
	"fmt"
)

type InternetEgressListData struct {
	Records    []*InternetEgress `json:"records,omitempty"`
	TotalCount int               `json:"totalCount,omitempty"`
}
type InternetEgressDataResponse struct {
	Data *InternetEgressListData `json:"data,omitempty"`
}

type InternetEgressData struct {
	Data *InternetEgress `json:"data,omitempty"`
}

type InternetEgress struct {
	ID            string              `json:"id,omitempty"`
	Name          string              `json:"name,omitempty"`
	Action        string              `json:"action,omitempty"`
	Matches       *IE_Matches         `json:"matches,omitempty"`
	Namespaces    *[]NamespaceList    `json:"namespaces,omitempty"`
	Networks      *[]NetworkList      `json:"networks,omitempty"`
	NetworkGroups *[]NetworkGroupList `json:"networkGroups,omitempty"`
}

type IE_Matches struct {
	Fqdn *[]FqdnDetails `json:"fqdn,omitempty"`
}

type FqdnDetails struct {
	Property  string  `json:"property,omitempty"`
	Operation string  `json:"operation,omitempty"`
	Values    *Values `json:"values,omitempty"`
}

type NamespaceList struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type NetworkList struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	NamespaceName string `json:"namespaceName,omitempty"`
	Status        string `json:"status,omitempty"`
}

type NetworkGroupList struct {
	ID            string         `json:"id,omitempty"`
	Name          string         `json:"name,omitempty"`
	NamespaceName string         `json:"namespaceName,omitempty"`
	Networks      *[]NetworkList `json:"networks,omitempty"`
}

type IEMatchItem struct {
	FQDN PropertyItemList `json:"fqdn,omitempty"`
}

func (prosimoClient *ProsimoClient) ReadInternetEgressJson() IEMatchItem {

	return GetInternetEgressPolicyServerTemplate()
}

func (prosimoClient *ProsimoClient) GetInternetEgressPolicy(ctx context.Context) ([]*InternetEgress, error) {

	PolicySearchInput := InternetEgress{}
	req, err := prosimoClient.api_client.NewRequest("POST", GetInternetEgressEndpoint, PolicySearchInput)
	if err != nil {
		return nil, err
	}

	policyListData := &InternetEgressDataResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, policyListData)
	if err != nil {
		return nil, err
	}

	return policyListData.Data.Records, nil

}

func (prosimoClient *ProsimoClient) GetInternetEgressPolicyByName(ctx context.Context, policyName string) (*InternetEgress, error) {

	policyList, err := prosimoClient.GetInternetEgressPolicy(ctx)
	if err != nil {
		return nil, err
	}

	var policy *InternetEgress
	for _, returnedPolicy := range policyList {
		if returnedPolicy.Name == policyName {
			policy = returnedPolicy
			break
		}
	}

	if policy == nil {
		return nil, errors.New("Internet Egress Policy doesn't exist")
	}

	return policy, nil

}

func (prosimoClient *ProsimoClient) GetInternetEgressPolicyID(ctx context.Context, policyName string) (string, error) {

	policyList, err := prosimoClient.GetInternetEgressPolicy(ctx)
	if err != nil {
		return "", err
	}

	var policy *InternetEgress
	for _, returnedPolicy := range policyList {
		if returnedPolicy.Name == policyName {
			policy = returnedPolicy
			break
		}
	}

	if policy == nil {
		return "", errors.New("Internet Egress Policy doesn't exist")
	}

	return policy.ID, nil

}

func (prosimoClient *ProsimoClient) GetInternetEgressPolicyByID(ctx context.Context, policyID string) (*InternetEgress, error) {

	policyList, err := prosimoClient.GetInternetEgressPolicy(ctx)
	if err != nil {
		return nil, err
	}

	var policy *InternetEgress
	for _, returnedPolicy := range policyList {
		if returnedPolicy.ID == policyID {
			policy = returnedPolicy
			break
		}
	}
	if policy == nil {
		return nil, errors.New("Internet Egress Policy doesn't exist")
	}

	return policy, nil

}

func (prosimoClient *ProsimoClient) CreateInternetEgressPolicy(ctx context.Context, inputPolicy *InternetEgress) (*InternetEgressData, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", InternetEgressEndpoint, inputPolicy)
	if err != nil {
		return nil, err
	}

	policyListData := &InternetEgressData{}
	_, err = prosimoClient.api_client.Do(ctx, req, policyListData)
	if err != nil {
		return nil, err
	}

	return policyListData, nil

}

func (prosimoClient *ProsimoClient) UpdateInternetEgressPolicy(ctx context.Context, inputPolicy *InternetEgress) (*InternetEgressData, error) {

	updatePolicyEndpoint := fmt.Sprintf("%s/%s", InternetEgressEndpoint, inputPolicy.ID)
	req, err := prosimoClient.api_client.NewRequest("PUT", updatePolicyEndpoint, inputPolicy)
	if err != nil {
		return nil, err
	}
	policyListData := &InternetEgressData{}

	_, err = prosimoClient.api_client.Do(ctx, req, policyListData)
	if err != nil {
		return nil, err
	}
	return policyListData, nil
}

func (prosimoClient *ProsimoClient) DeleteInternetEgressPolicy(ctx context.Context, policyID string) error {

	DeletePolicyEndpoint := fmt.Sprintf("%s/%s", InternetEgressEndpoint, policyID)
	req, err := prosimoClient.api_client.NewRequest("DELETE", DeletePolicyEndpoint, nil)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}
