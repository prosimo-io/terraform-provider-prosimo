package client

import (
	"context"
	"fmt"
)

type Prefix_RouteID struct {
	Prefix        string   `json:"prefix,omitempty"`
	RouteTableIDS []string `json:"routeTableIDs,omitempty"`
}
// type Regions_route struct {
// 	All      bool           `json:"all,omitempty"`
// 	Selected []Selected_Reg `json:"selected,omitempty"`
// }

type Route_entry_network struct {
	ID             string            `json:"id,omitempty"`
	CloudKeyID     string            `json:"cloudKeyID,omitempty"`
	CSP            string            `json:"csp,omitempty"`
	CloudRegion    string            `json:"cloudRegion,omitempty"`
	CloudNetworkID string            `json:"loudNetworkID,omitempty"`
	PrefixesRTID   *[]Prefix_RouteID `json:"prefixRouteTableIDs,omitempty"`
	Enabled        bool              `json:"enabled,omitempty"`
}
type Route_entry_network_res struct {
	Data *Route_entry_region `json:"data,omitempty"`
}

type Route_entry_network_Data struct {
	Records    []*Route_entry_region `json:"records,omitempty"`
	TotalCount int                   `json:"totalCount,omitempty"`
}

func (prosimoClient *ProsimoClient) CreateNetworkRouteEntry(ctx context.Context, routeEntry *Route_entry_network) (*Route_entry_network_res, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", CreateRouteEntryEndpoint, routeEntry)
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

func (prosimoClient *ProsimoClient) UpdateNetworkRouteEntry(ctx context.Context, routeEntry *Route_entry_region) error {

	updateRouteEntryEndpoint := fmt.Sprintf("%s/%s", CreateRouteEntryEndpoint, routeEntry.ID)

	req, err := prosimoClient.api_client.NewRequest("PUT", updateRouteEntryEndpoint, routeEntry)
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
func (prosimoClient *ProsimoClient) GetNetworkRouteEntry(ctx context.Context) ([]*Route_entry_region, error) {

	SearchInput := Route_entry_region{}
	req, err := prosimoClient.api_client.NewRequest("POST", GetRouteEntryEndpoint, SearchInput)
	if err != nil {
		return nil, err
	}

	resData := &Route_entry_region_Data{}
	_, err = prosimoClient.api_client.Do(ctx, req, resData)
	if err != nil {
		return nil, err
	}

	return resData.Records, nil

}

func (prosimoClient *ProsimoClient) DeleteNetworkRouteEntry(ctx context.Context, route_entry_id string) error {

	updateRouteEntryEndpoint := fmt.Sprintf("%s/%s", CreateRouteEntryEndpoint, route_entry_id)

	req, err := prosimoClient.api_client.NewRequest("DELETE", updateRouteEntryEndpoint, nil)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}
