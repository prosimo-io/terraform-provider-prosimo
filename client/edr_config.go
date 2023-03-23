package client

import (
	"context"
	"fmt"
)

type EDR_Config struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"nickname,omitempty"`
	Vendor      string `json:"vendor,omitempty"`
	Status      string `json:"status,omitempty"`
	Auth        AUTH   `json:"auth,omitempty"`
	CreatedTime string `json:"createdTime,omitempty"`
	UpdatedTime string `json:"updatedTime,omitempty"`
}

type AUTH struct {
	BaseURL      string `json:"baseURL,omitempty"`
	ClientID     string `json:"clientID,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	CustomerID   string `json:"customerID,omitempty"`
	MSSP         bool   `json:"mssp,omitempty"`
}

type EDR_Config_Res struct {
	EdrConfig EDR_Config `json:"data,omitempty"`
}

type EDR_Config_ResList struct {
	EdrConfigList []*EDR_Config `json:"data,omitempty"`
}

func (auth AUTH) String() string {
	return fmt.Sprintf("{BaseURL:%s, ClientID:%s, ClientSecret:%s, CustomerID:%s, MSSP:%t}",
		auth.BaseURL, auth.ClientID, auth.ClientSecret, auth.CustomerID, auth.MSSP)
}

func (edr_config EDR_Config) String() string {
	return fmt.Sprintf("{Id:%s, CreatedTime:%s, UpdatedTime:%s, Name:%s, Vendor:%s, Status:%s, }",
		edr_config.Id, edr_config.CreatedTime, edr_config.UpdatedTime, edr_config.Name, edr_config.Vendor, edr_config.Status)
}

func (prosimoClient *ProsimoClient) GetEDRConf(ctx context.Context) (*EDR_Config_ResList, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", EDRConfigEndpoint, nil)
	if err != nil {
		return nil, err
	}
	edrres := &EDR_Config_ResList{}
	_, err = prosimoClient.api_client.Do(ctx, req, edrres)

	if err != nil {
		return nil, err
	}

	return edrres, nil

}

func (prosimoClient *ProsimoClient) GetEDRVendor(ctx context.Context, vendor string) (bool, error) { //This function is called by edr profile

	req, err := prosimoClient.api_client.NewRequest("GET", EDRConfigEndpoint, nil)
	if err != nil {
		return false, err
	}
	edrres := &EDR_Config_ResList{}
	_, err = prosimoClient.api_client.Do(ctx, req, edrres)

	if err != nil {
		return false, err
	}
	flag := false
	for _, conf := range edrres.EdrConfigList {
		if conf.Vendor == vendor {
			flag = true
		}
	}
	return flag, nil
}

func (prosimoClient *ProsimoClient) CreateEDRConf(ctx context.Context, edrconfops *EDR_Config) (*EDR_Config_Res, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", EDRConfigEndpoint, edrconfops)
	if err != nil {
		return nil, err
	}

	edrconfpostData := &EDR_Config_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, edrconfpostData)
	if err != nil {
		return nil, err
	}

	return edrconfpostData, nil

}

func (prosimoClient *ProsimoClient) UpdateEDRConf(ctx context.Context, edrconfops *EDR_Config) (*EDR_Config_Res, error) {

	updateEDRConfigEndpoint := fmt.Sprintf("%s/%s", EDRConfigEndpoint, edrconfops.Id)
	req, err := prosimoClient.api_client.NewRequest("PUT", updateEDRConfigEndpoint, edrconfops)
	if err != nil {
		return nil, err
	}

	edrconfputData := &EDR_Config_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, edrconfputData)
	if err != nil {
		return nil, err
	}

	return edrconfputData, nil

}

func (prosimoClient *ProsimoClient) DeleteEDRConf(ctx context.Context, edrID string) error {

	updateEDRConfigEndpoint := fmt.Sprintf("%s/%s", EDRConfigEndpoint, edrID)
	req, err := prosimoClient.api_client.NewRequest("DELETE", updateEDRConfigEndpoint, nil)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}
