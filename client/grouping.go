package client

import (
	"context"
	"fmt"
)

type Grp_Config struct {
	Id          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Type        string   `json:"type,omitempty"`
	SubType     string   `json:"subType,omitempty"`
	Members     []string `json:"members,omitempty"`
	Details     DeTails  `json:"details,omitempty"`
	CreatedTime string   `json:"createdTime,omitempty"`
	UpdatedTime string   `json:"updatedTime,omitempty"`
}

type DeTails struct {
	Apps   []App    `json:"apps,omitempty"`
	Names  []Name   `json:"names,omitempty"`
	Time   []Time   `json:"time,omitempty"`
	Ranges []string `json:"string,omitempty"`
}

type App struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Name struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Time struct {
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	TimeZone string `json:"timeZone,omitempty"`
}

type Grp_Config_ResData struct {
	Grp_Config_Res Grp_Config_Res `json:"data,omitempty"`
}

type Grp_Config_Res struct {
	Id         string        `json:"id,omitempty"`
	TotalCount int           `json:"totalCount,omitempty"`
	Records    []*Grp_Config `json:"records,omitempty"`
}


func (prosimoClient *ProsimoClient) GetGrpConf(ctx context.Context, grptype string) (*Grp_Config_ResData, error) {
	grpConfig := Grp_Config{
		Type: grptype,
	}
	updateGroupingEndpoint := fmt.Sprintf("%s/%s", GroupingEndpoint, "search")
	req, err := prosimoClient.api_client.NewRequest("POST", updateGroupingEndpoint, grpConfig)
	if err != nil {
		return nil, err
	}
	grpres := &Grp_Config_ResData{}
	_, err = prosimoClient.api_client.Do(ctx, req, grpres)

	if err != nil {
		return nil, err
	}

	return grpres, nil

}

func (prosimoClient *ProsimoClient) CreateGrpConf(ctx context.Context, grpconfops *Grp_Config) (*Grp_Config_ResData, error) {
	// updateGroupingEndpoint := fmt.Sprintf("%s/%s", GroupingEndpoint, "search")
	req, err := prosimoClient.api_client.NewRequest("POST", GroupingEndpoint, grpconfops)
	if err != nil {
		return nil, err
	}

	grpconfpostData := &Grp_Config_ResData{}
	_, err = prosimoClient.api_client.Do(ctx, req, grpconfpostData)
	if err != nil {
		return nil, err
	}

	return grpconfpostData, nil

}

func (prosimoClient *ProsimoClient) UpdateGrpConf(ctx context.Context, grpconfops *Grp_Config) (*Grp_Config_ResData, error) {

	updateGroupingEndpoint := fmt.Sprintf("%s/%s", GroupingEndpoint, grpconfops.Id)
	req, err := prosimoClient.api_client.NewRequest("PUT", updateGroupingEndpoint, grpconfops)
	if err != nil {
		return nil, err
	}

	grpconfputData := &Grp_Config_ResData{}
	_, err = prosimoClient.api_client.Do(ctx, req, grpconfputData)
	if err != nil {
		return nil, err
	}

	return grpconfputData, nil

}

func (prosimoClient *ProsimoClient) DeleteGrpConf(ctx context.Context, grpID string) error {

	updateGroupingEndpoint := fmt.Sprintf("%s/%s", GroupingEndpoint, grpID)
	req, err := prosimoClient.api_client.NewRequest("DELETE", updateGroupingEndpoint, nil)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}
