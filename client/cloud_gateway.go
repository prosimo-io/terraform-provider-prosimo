package client

import (
	"context"
	"errors"
	"fmt"
)

type CloudGateway struct {
	ID                 string `json:"id,omitempty"`
	EdgeConnectivityID string `json:"edgeConnectivityID,omitempty"`
	Name               string `json:"name,omitempty"`
	Attachment         string `json:"attachment,omitempty"`
	RouteTable         string `json:"routeTable,omitempty"`
}

type ConnectivityOptions struct {
	AttachPoint       string `json:"attachPoint,omitempty"`
	CloudCredsID      string `json:"cloudCredsID,omitempty"`
	CloudType         string `json:"cloudType,omitempty"`
	CloudRegion       string `json:"cloudRegion,omitempty"`
	ConnectivityType  string `json:"connectivityType,omitempty"`
	EdgeID            string `json:"edgeID,omitempty"`
	ID                string `json:"id,omitempty"`
	InternallyManaged bool   `json:"internallyManaged,omitempty"`
	Name              string `json:"name,omitempty"`
	Network           string `json:"network,omitempty"`
	Region            string `json:"region,omitempty"`
	Status            string `json:"status,omitempty"`
	TeamID            string `json:"teamID,omitempty"`
	Attachment        string `json:"attachment,omitempty"`
	RouteTable        string `json:"routeTable,omitempty"`
}

type ConnectivityOptionsResponse struct {
	Data []*ConnectivityOptions `json:"data,omitempty"`
}

type CloudGatewayResponse struct {
	Data *CloudGateway `json:"data,omitempty"`
}

type CloudGatewayReadResponse struct {
	Data []*CloudGateway `json:"data,omitempty"`
}

func (prosimoClient *ProsimoClient) CreateCloudGateway(ctx context.Context, cgInput *CloudGateway) (*CloudGatewayResponse, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", CloudGatewayEndpoint, cgInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &CloudGatewayResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) UpdateCloudGateway(ctx context.Context, cgInput *CloudGateway) error {

	updateCloudGatewayEndpoint := fmt.Sprintf("%s/%s", CloudGatewayEndpoint, cgInput.ID)

	req, err := prosimoClient.api_client.NewRequest("PUT", updateCloudGatewayEndpoint, cgInput)
	if err != nil {
		return err
	}

	resourcePostResponseData := &CloudGatewayResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return err
	}

	return nil

}

func (prosimoClient *ProsimoClient) DeleteCloudGateway(ctx context.Context, gwID string) error {

	updateCloudGatewayeEndpoint := fmt.Sprintf("%s/%s", CloudGatewayEndpoint, gwID)

	req, err := prosimoClient.api_client.NewRequest("DELETE", updateCloudGatewayeEndpoint, nil)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}

func (prosimoClient *ProsimoClient) GetCloudGateway(ctx context.Context) ([]*ConnectivityOptions, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", CloudGatewayEndpoint, nil)
	if err != nil {
		return nil, err
	}

	gwListData := &ConnectivityOptionsResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, gwListData)
	if err != nil {
		return nil, err
	}

	return gwListData.Data, nil

}

func (prosimoClient *ProsimoClient) GetCloudGatewayByID(ctx context.Context, id string) (*ConnectivityOptions, error) {

	gwList, err := prosimoClient.GetCloudGateway(ctx)
	if err != nil {
		return nil, err
	}
	var gateWay *ConnectivityOptions
	for _, returnedGW := range gwList {
		if returnedGW.ID == id {
			gateWay = returnedGW
			break
		}
	}

	if gateWay == nil {
		return nil, errors.New("cloud gateway doesn't exists")
	}

	return gateWay, nil

}

func (prosimoClient *ProsimoClient) GetCnnectivityOptions(ctx context.Context) ([]*ConnectivityOptions, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", ConnectivityOptionsEndpont, nil)
	if err != nil {
		return nil, err
	}

	connOptnListData := &ConnectivityOptionsResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, connOptnListData)
	if err != nil {
		return nil, err
	}

	return connOptnListData.Data, nil

}
