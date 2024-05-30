package client

import (
	"context"
	"errors"
	"fmt"
	"log"
)

type Selected_Reg struct {
	CSP   string   `json:"csp,omitempty"`
	Names []string `json:"names,omitempty"`
}
type Regions_route struct {
	All      bool           `json:"all,omitempty"`
	Selected []Selected_Reg `json:"selected,omitempty"`
}

type Route_entry_region struct {
	ID             string         `json:"id,omitempty"`
	Status         string         `json:"status,omitempty"`
	Prefixes       []string       `json:"prefixes,omitempty"`
	Regions        *Regions_route `json:"regions,omitempty"`
	Type           string         `json:"type,omitempty"`
	Enabled        bool           `json:"enabled,omitempty"`
	OverWriteRoute bool           `json:"overwriteRoute,omitempty"`
	TeamID         string         `json:"teamID,omitempty"`
	CreatedTime    string         `json:"createdTime,omitempty"`
	UpdatedTime    string         `json:"updatedTime,omitempty"`
}

type Route_entry_region_res struct {
	Data *Route_entry_region `json:"data,omitempty"`
}

type Route_entry_region_Data_res struct {
	Data *Route_entry_region_Data `json:"data,omitempty"`
}

type Route_entry_region_Data struct {
	Records    []*Route_entry_region `json:"records,omitempty"`
	TotalCount int                   `json:"totalCount,omitempty"`
}

func (prosimoClient *ProsimoClient) CreateRouteEntry(ctx context.Context, routeEntry *Route_entry_region) (*Route_entry_region_res, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", CreateRouteEntryEndpoint, routeEntry)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &Route_entry_region_res{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) UpdateRouteEntry(ctx context.Context, routeEntry *Route_entry_region) error {

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

func (prosimoClient *ProsimoClient) GetRouteEntry(ctx context.Context) ([]*Route_entry_region, error) {

	SearchInput := Route_entry_region{}
	req, err := prosimoClient.api_client.NewRequest("POST", GetRouteEntryEndpoint, SearchInput)
	if err != nil {
		return nil, err
	}

	resData := &Route_entry_region_Data_res{}
	_, err = prosimoClient.api_client.Do(ctx, req, resData)
	if err != nil {
		return nil, err
	}
	log.Println("resData", resData)

	return resData.Data.Records, nil

}

func (prosimoClient *ProsimoClient) GetRouteEntryByID(ctx context.Context, id string) (*Route_entry_region, error) {

	rpList, err := prosimoClient.GetRouteEntry(ctx)
	if err != nil {
		return nil, err
	}
	log.Println("rpList", rpList)
	log.Println("id", id)
	var routeEntry *Route_entry_region
	for _, returnedrp := range rpList {
		log.Println("eturnedrp.ID", returnedrp.ID)
		if returnedrp.ID == id {
			routeEntry = returnedrp
			break
		}
	}

	if routeEntry == nil {
		return nil, errors.New("routing prefix doesn't exists")
	}

	return routeEntry, nil

}

func (prosimoClient *ProsimoClient) DeleteRouteEntry(ctx context.Context, route_entry_id string) error {

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
