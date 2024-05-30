package client

import (
	"context"
	"errors"
	"fmt"
	"log"
)

type Prefix_RouteID struct {
	Prefix        string         `json:"prefix,omitempty"`
	RouteTableIDS []Route_Tables `json:"routeTables,omitempty"`
}

type Route_Tables struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// type Regions_route struct {
// 	All      bool           `json:"all,omitempty"`
// 	Selected []Selected_Reg `json:"selected,omitempty"`
// }

type Route_entry_network struct {
	ID               string            `json:"id,omitempty"`
	CloudKeyID       string            `json:"cloudKeyID,omitempty"`
	CSP              string            `json:"csp,omitempty"`
	CloudRegion      string            `json:"cloudRegion,omitempty"`
	CloudNetworkID   string            `json:"cloudNetworkID,omitempty"`
	CloudNetworkName string            `json:"cloudNetworkName,omitempty"`
	PrefixesRT       *[]Prefix_RouteID `json:"prefixRouteTables,omitempty"`
	Enabled          bool              `json:"enabled,omitempty"`
	OverwriteRoute   bool              `json:"overwriteRoute,omitempty"`
	Status           string            `json:"status,omitempty"`
	TeamID           string            `json:"teamID,omitempty"`
	AuditID          string            `json:"auditID,omitempty"`
	CreatedTime      string            `json:"createdTime,omitempty"`
	UpdatedTime      string            `json:"updatedTime,omitempty"`
}

type Route_entry_network_res struct {
	Data *Route_entry_network `json:"data,omitempty"`
}

type Route_entry_network_res_data struct {
	Data *Route_entry_network_Data `json:"data,omitempty"`
}

type Route_entry_network_Data struct {
	Records    []*Route_entry_network `json:"records,omitempty"`
	TotalCount int                    `json:"totalCount,omitempty"`
}

func (prosimoClient *ProsimoClient) CreateNetworkRouteEntry(ctx context.Context, routeEntry *Route_entry_network) (*Route_entry_network_res, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", CreateRouteEntryNetworkEndpoint, routeEntry)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &Route_entry_network_res{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) EnableNetworkRouteEntry(ctx context.Context, networkPrefixID string) (*Route_entry_network_res, error) {
	EnableRouteEntryNetworkEndpointUpdated := fmt.Sprintf(EnableRouteEntryNetworkEndpoint, networkPrefixID)
	blankInput := Route_entry_region{}
	req, err := prosimoClient.api_client.NewRequest("PUT", EnableRouteEntryNetworkEndpointUpdated, blankInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &Route_entry_network_res{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}
func (prosimoClient *ProsimoClient) DisableNetworkRouteEntry(ctx context.Context, networkPrefixID string) (*Route_entry_network_res, error) {
	DisableRouteEntryNetworkEndpointUpdated := fmt.Sprintf(DisableRouteEntryNetworkEndpoint, networkPrefixID)
	blankInput := Route_entry_region{}
	req, err := prosimoClient.api_client.NewRequest("PUT", DisableRouteEntryNetworkEndpointUpdated, blankInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &Route_entry_network_res{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) UpdateNetworkRouteEntry(ctx context.Context, routeEntry *Route_entry_network) error {

	updateRouteEntryNetworkEndpoint := fmt.Sprintf("%s/%s", CreateRouteEntryNetworkEndpoint, routeEntry.ID)

	req, err := prosimoClient.api_client.NewRequest("PUT", updateRouteEntryNetworkEndpoint, routeEntry)
	if err != nil {
		return err
	}

	// resourcePostResponseData := &Route_entry_region_res{}
	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}
func (prosimoClient *ProsimoClient) GetNetworkRouteEntry(ctx context.Context) ([]*Route_entry_network, error) {

	SearchInput := Route_entry_region{}
	req, err := prosimoClient.api_client.NewRequest("POST", GetRouteEntryNetworkEndpoint, SearchInput)
	if err != nil {
		return nil, err
	}

	resData := &Route_entry_network_res_data{}
	_, err = prosimoClient.api_client.Do(ctx, req, resData)
	if err != nil {
		return nil, err
	}

	return resData.Data.Records, nil

}

func (prosimoClient *ProsimoClient) GetNetworkRouteEntryByID(ctx context.Context, id string) (*Route_entry_network, error) {

	npList, err := prosimoClient.GetNetworkRouteEntry(ctx)
	log.Println("npList", npList)
	if err != nil {
		return nil, err
	}
	var nwEntry *Route_entry_network
	for _, returnednp := range npList {
		if returnednp.ID == id {
			nwEntry = returnednp
			break
		}
	}

	if nwEntry == nil {
		return nil, errors.New("network prefix doesn't exists")
	}

	return nwEntry, nil

}

func (prosimoClient *ProsimoClient) DeleteNetworkRouteEntry(ctx context.Context, route_entry_id string) error {

	updateRouteEntryNetworkEndpoint := fmt.Sprintf("%s/%s", CreateRouteEntryNetworkEndpoint, route_entry_id)

	req, err := prosimoClient.api_client.NewRequest("DELETE", updateRouteEntryNetworkEndpoint, nil)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}
