package client

import (
	"context"
)

type Connector_Placement struct {
	ConnectorPlacementAppVpc bool   `json:"connectorPlacementInAppVpc,omitempty"`
	TeamID                   string `json:"teamID,omitempty"`
	UpdatedTime              string `json:"updatedTime,omitempty"`
}

type Connector_Placement_Read_response struct {
	ConnectorPlacementAppVpcStatus Connector_Placement `json:"data,omitempty"`
}

func (prosimoClient *ProsimoClient) PutConnectorPlacement(ctx context.Context, con_pl *Connector_Placement) error {
	req, err := prosimoClient.api_client.NewRequest("PUT", ConnectorPlacementEndpoint, con_pl)
	if err != nil {
		return err
	}

	// prepListData := &Connector_Placement_Read_response{}
	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}

func (prosimoClient *ProsimoClient) GetConnectorPlacement(ctx context.Context) (*Connector_Placement_Read_response, error) {
	req, err := prosimoClient.api_client.NewRequest("GET", ConnectorPlacementEndpoint, nil)
	if err != nil {
		return nil, err
	}

	prepListData := &Connector_Placement_Read_response{}
	_, err = prosimoClient.api_client.Do(ctx, req, prepListData)
	if err != nil {
		return nil, err
	}

	return prepListData, nil

}
