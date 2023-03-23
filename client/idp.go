package client

import (
	"context"
	"fmt"
)

type IDP struct {
	AccountURL         string     `json:"accountUrl,omitempty"`
	AppIDs             []string   `json:"appIDs,omitempty"`
	Auth_Type          string     `json:"authType,omitempty"`
	Api_Creds_Provided string     `json:"apiCredsProvided,omitempty"`
	Details            IDPDetails `json:"details,omitempty"`
	DomainURL          []string   `json:"domainURL,omitempty"`
	ID                 string     `json:"id,omitempty"`
	IDPName            string     `json:"idpName,omitempty"`
	Select_Type        string     `json:"selectType,omitempty"`
	// Status           string     `json:"status,omitempty"`
	TeamID        string `json:"teamID,omitempty"`
	IsFileUpdated bool   `json:"isFileUpdated,omitempty"`
}

type IDPDetails struct {
	ClientID        string `json:"clientID,omitempty"`
	ClientSecret    string `json:"clientSecret,omitempty"`
	APIToken        string `json:"apiToken,omitempty"`
	APIClientID     string `json:"apiClientID,omitempty"`
	APIClientSecret string `json:"apiClientSecret,omitempty"`
	MetadataURL     string `json:"metadataURL,omitempty"`
	Metadata        string `json:"metadata,omitempty"`
	Region          string `json:"apiRegion,omitempty"`
	EnvID           string `json:"envID,omitempty"`
	APIEmail        string `json:"apiEmail,omitempty"`
	CustomerID      string `json:"customerID,omitempty"`
	Domain          string `json:"domain,omitempty"`
	FilePath        string `json:"filepath,omitempty"`
	APIFile         string `json:"apiFile,omitempty"`
}

type IDPList struct {
	IDPs []*IDP `json:"data,omitempty"`
}

func (idpDetails IDPDetails) String() string {
	return fmt.Sprintf("{ClientID:%s, ClientSecret:%s, APIToken:%s, MetadataURL:%s, Metadata:%s}",
		idpDetails.ClientID, idpDetails.ClientSecret, idpDetails.APIToken, idpDetails.MetadataURL, idpDetails.Metadata)
}

func (idp IDP) String() string {
	return fmt.Sprintf("{IDPName:%s, ApiCredsProvided:%s, AccountURL:%s, DomainURL:%s, AppIDs:%s, AuthType:%s, SelectType:%s}",
		idp.IDPName, idp.Api_Creds_Provided, idp.AccountURL, idp.DomainURL, idp.AppIDs, idp.Auth_Type, idp.Select_Type)
}

func (prosimoClient *ProsimoClient) CreateIDP(ctx context.Context, idp *IDP) (*ResourcePostResponseData, error) {

	return prosimoClient.api_client.PostRequest(ctx, IDPEndpoint, idp)

}

func (prosimoClient *ProsimoClient) UpdateIDP(ctx context.Context, idp *IDP) (*ResourcePostResponseData, error) {
	updateIDPEndpoint := fmt.Sprintf("%s/%s", IDPEndpoint, idp.ID)

	return prosimoClient.api_client.PutRequest(ctx, updateIDPEndpoint, idp)

}

func (prosimoClient *ProsimoClient) CreateGoogleIDP(ctx context.Context, idp *IDP) (*ResourcePostResponseData, error) {
	fileUploadEdPt := fmt.Sprintf("%s/%s", IDPEndpoint, "creds/file")

	m := make(map[string]string)
	m["idpName"] = idp.IDPName
	m["authType"] = idp.Auth_Type
	m["selectType"] = idp.Select_Type
	if idp.Select_Type == "partner" {
		// count := len(idp.DomainURL)
		// for i = 0 ; i < count;
		for _, i := range idp.DomainURL {
			m["domainURL"] = i
		}
		for _, i := range idp.AppIDs {
			m["appID"] = i
		}
	}
	m["accountUrl"] = idp.AccountURL
	m["apiCredsProvided"] = idp.Api_Creds_Provided
	m["apiEmail"] = idp.Details.APIEmail
	m["customerID"] = idp.Details.CustomerID
	m["domain"] = idp.Details.Domain
	m["clientID"] = idp.Details.ClientID
	m["clientSecret"] = idp.Details.ClientSecret

	req, err := prosimoClient.api_client.ReqFileUploadIDP("POST", fileUploadEdPt, m, idp.Details.FilePath)
	if err != nil {
		return nil, err
	}

	resourceResponseData := &ResourcePostResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourceResponseData)
	if err != nil {
		return nil, err
	}

	return resourceResponseData, nil
}

func (prosimoClient *ProsimoClient) UpdateGoogleIDP(ctx context.Context, idp *IDP) (*ResourcePostResponseData, error) {
	fileUploadEdPt := fmt.Sprintf("%s/%s", IDPEndpoint, "creds/file")

	m := make(map[string]string)
	m["id"] = idp.ID
	m["idpName"] = idp.IDPName
	m["authType"] = idp.Auth_Type
	m["selectType"] = idp.Select_Type
	if idp.Select_Type == "partner" {
		for _, i := range idp.DomainURL {
			m["domainURL"] = i
		}
		for _, i := range idp.AppIDs {
			m["appID"] = i
		}
	}
	m["accountUrl"] = idp.AccountURL
	m["apiCredsProvided"] = idp.Api_Creds_Provided
	m["apiEmail"] = idp.Details.APIEmail
	m["customerID"] = idp.Details.CustomerID
	m["domain"] = idp.Details.Domain
	m["clientID"] = idp.Details.ClientID
	m["clientSecret"] = idp.Details.ClientSecret

	req, err := prosimoClient.api_client.ReqFileUploadIDP("POST", fileUploadEdPt, m, idp.Details.FilePath)
	if err != nil {
		return nil, err
	}

	resourceResponseData := &ResourcePostResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourceResponseData)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (prosimoClient *ProsimoClient) GetIDP(ctx context.Context) (*IDPList, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", IDPEndpoint, nil)
	if err != nil {
		return nil, err
	}

	idpList := &IDPList{}
	_, err = prosimoClient.api_client.Do(ctx, req, idpList)
	if err != nil {
		return nil, err
	}

	return idpList, nil

}

func (prosimoClient *ProsimoClient) DeleteIDP(ctx context.Context, idpID string) error {

	deleteIDPEndpt := fmt.Sprintf("%s/%s", IDPEndpoint, idpID)

	return prosimoClient.api_client.DeleteRequest(ctx, deleteIDPEndpt)

}
