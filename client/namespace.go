package client

import (
	"context"
	"errors"
	"fmt"
)

type Namespace struct {
	ID     string `json:"id,omitempty"`
	TaskID string `json:"taskID,omitempty"`
	Name   string `json:"name,omitempty"`
}

type NetActNamespace struct {
	ID         string   `json:"id,omitempty"`
	NetworkID  string   `json:"networkID,omitempty"`
	Namespaces []string `json:"namespaces,omitempty"`
}

type Namespace_Response struct {
	NamespaceResponse *Namespace `json:"data,omitempty"`
}
type NSListData struct {
	Namespaces []*Namespace `json:"namespaces,omitempty"`
	TotalCount int          `json:"totalCount,omitempty"`
}

type NSListDataResponse struct {
	Data *NSListData `json:"data,omitempty"`
}

func (prosimoClient *ProsimoClient) CreateNamespace(ctx context.Context, namespaceInput *Namespace) (*Namespace_Response, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", NameSpaceEndpoint, namespaceInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &Namespace_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) UpdateNamespace(ctx context.Context, namespaceInput *Namespace) (*Namespace_Response, error) {

	updateNamespaceEndpoint := fmt.Sprintf("%s/%s", NameSpaceEndpoint, namespaceInput.ID)

	req, err := prosimoClient.api_client.NewRequest("PUT", updateNamespaceEndpoint, namespaceInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &Namespace_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) AssignNetworkToNamespace(ctx context.Context, assignNetInput *[]NetActNamespace, namespaceid string) (*Namespace_Response, error) {

	updateNetworkAssignEndpoint := fmt.Sprintf("%s/%s", AssignNetworkEndpoint, namespaceid)

	req, err := prosimoClient.api_client.NewRequest("PUT", updateNetworkAssignEndpoint, assignNetInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &Namespace_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) ExportNetworkToNamespace(ctx context.Context, exportNetInput *[]NetActNamespace, namespaceid string) (*Namespace_Response, error) {

	updateNetworkExportEndpoint := fmt.Sprintf("%s/%s", ExportNetworkEndpoint, namespaceid)

	req, err := prosimoClient.api_client.NewRequest("PUT", updateNetworkExportEndpoint, exportNetInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &Namespace_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) WithdrawNetworkToNamespace(ctx context.Context, withdrawNetInput *[]NetActNamespace, namespaceid string) (*Namespace_Response, error) {

	updateNetworkWithdrawEndpoint := fmt.Sprintf("%s/%s", WithdrawNetworkEndpoint, namespaceid)

	req, err := prosimoClient.api_client.NewRequest("POST", updateNetworkWithdrawEndpoint, withdrawNetInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &Namespace_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) DeleteNamespace(ctx context.Context, namespaceID string) (*Namespace_Response, error) {

	updateNameSpaceEndpoint := fmt.Sprintf("%s/%s", NameSpaceEndpoint, namespaceID)

	req, err := prosimoClient.api_client.NewRequest("DELETE", updateNameSpaceEndpoint, nil)
	if err != nil {
		return nil,err
	}

	resourcePostResponseData := &Namespace_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil,err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) GetNamespace(ctx context.Context) ([]*Namespace, error) {

	NSSearchInput := Namespace{}
	req, err := prosimoClient.api_client.NewRequest("POST", GetNameSpaceEndpoint, NSSearchInput)
	if err != nil {
		return nil, err
	}

	nsListData := &NSListDataResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, nsListData)
	if err != nil {
		return nil, err
	}

	return nsListData.Data.Namespaces, nil

}

func (prosimoClient *ProsimoClient) GetNamespaceByID(ctx context.Context, id string) (*Namespace, error) {

	nsList, err := prosimoClient.GetNamespace(ctx)
	if err != nil {
		return nil, err
	}
	var nameSpace *Namespace
	for _, returnedNS := range nsList {
		if returnedNS.ID == id {
			nameSpace = returnedNS
			break
		}
	}

	if nameSpace == nil {
		return nil, errors.New("Namespace doesn't exists")
	}

	return nameSpace, nil

}

func (prosimoClient *ProsimoClient) GetNamespaceByName(ctx context.Context, name string) (*Namespace, error) {
	nsList, err := prosimoClient.GetNamespace(ctx)
	if err != nil {
		return nil, err
	}
	var nameSpace *Namespace
	for _, returnedNS := range nsList {
		if returnedNS.Name == name {
			nameSpace = returnedNS
			break
		}
	}

	if nameSpace == nil {
		return nil, errors.New("Namespace doesn't exists")
	}

	return nameSpace, nil

}
