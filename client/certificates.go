package client

import (
	"context"
	"fmt"
)

type ReadCertDetails struct {
	ID                 string `json:"id,omitempty"`
	TeamID             string `json:"teamID,omitempty"`
	URL                string `json:"url,omitempty"`
	CA                 string `json:"ca,omitempty"`
	Type               string `json:"type,omitempty"`
	DN                 string `json:"dn,omitempty"`
	Status             string `json:"status,omitempty"`
	ISTeamCert         bool   `json:"isTeamCert,omitempty"`
	Certificate        string `json:"certificate,omitempty"`
	Certificatehash    string `json:"certificatehash,omitempty"`
	SAN                string `json:"san,omitempty"`
	Expirytime         string `json:"expirytime,omitempty"`
	Generatedby        string `json:"generatedby,omitempty"`
	Updatedtime        string `json:"updatedtime,omitempty"`
	Issuetime          string `json:"issuetime,omitempty"`
	Signingalgorithm   string `json:"signingalgorith,omitempty"`
	Publickeyalgorithm string `json:"publickeyalgorith,omitempty"`
	Notified           bool   `json:"notified,omitempty"`
	Keysize            string `json:"keysize,omitempty"`
	Createdtime        string `json:"createdtime,omitempty"`
}

type ReadCertDetailsResponse struct {
	CertResList []ReadCertDetails `json:"data,omitempty"`
}

type UploadCert struct {
	CertPath string `json:"certificate,omitempty"`
	KeyPath  string `json:"privateKey,omitempty"`
}

func (prosimoClient *ProsimoClient) UploadCert(ctx context.Context, certfile string, keyfile string) (*ResourcePostResponseData, error) {

	req, err := prosimoClient.api_client.ReqCertUpload("POST", GenerateCertEndpoint, certfile, keyfile)
	if err != nil {
		return nil, err
	}
	resourcePostResponseData := &ResourcePostResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}
	return resourcePostResponseData, nil
}

func (prosimoClient *ProsimoClient) UploadCertUpdate(ctx context.Context, certfile string, keyfile string, certid string) (*ResourcePostResponseData, error) {
	updateGenerateCertEndpoint := fmt.Sprintf("%s/%s", GenerateCertEndpoint, certid)
	req, err := prosimoClient.api_client.ReqCertUpload("POST", updateGenerateCertEndpoint, certfile, keyfile)
	if err != nil {
		return nil, err
	}
	resourcePostResponseData := &ResourcePostResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}
	return resourcePostResponseData, nil
}

func (prosimoClient *ProsimoClient) UploadCertClient(ctx context.Context, certfile string, keyfile string) (*ResourcePostResponseData, error) {

	req, err := prosimoClient.api_client.ReqCertUpload("POST", UploadClientCertEndpoint, certfile, keyfile)
	if err != nil {
		return nil, err
	}
	resourcePostResponseData := &ResourcePostResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}
	return resourcePostResponseData, nil
}

func (prosimoClient *ProsimoClient) UploadCertClientUpdate(ctx context.Context, certfile string, keyfile string, certid string) (*ResourcePostResponseData, error) {
	updateUploadClientCertEndpoint := fmt.Sprintf("%s/%s", UploadClientCertEndpoint, certid)
	req, err := prosimoClient.api_client.ReqCertUpload("POST", updateUploadClientCertEndpoint, certfile, keyfile)
	if err != nil {
		return nil, err
	}
	resourcePostResponseData := &ResourcePostResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}
	return resourcePostResponseData, nil
}

func (prosimoClient *ProsimoClient) UploadCertCA(ctx context.Context, cacertfile string) (*ResourcePostResponseData, error) {

	req, err := prosimoClient.api_client.ReqCACertUpload("POST", UploadCACertEndpoint, cacertfile)
	if err != nil {
		return nil, err
	}
	resourcePostResponseData := &ResourcePostResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}
	return resourcePostResponseData, nil
}

func (prosimoClient *ProsimoClient) UploadCertCAUpdate(ctx context.Context, cacertfile string, certid string) (*ResourcePostResponseData, error) {
	updateUploadCACertEndpoint := fmt.Sprintf("%s/%s", UploadCACertEndpoint, certid)
	req, err := prosimoClient.api_client.ReqCACertUpload("POST", updateUploadCACertEndpoint, cacertfile)
	if err != nil {
		return nil, err
	}
	resourcePostResponseData := &ResourcePostResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}
	return resourcePostResponseData, nil
}

func (prosimoClient *ProsimoClient) GetCertDetails(ctx context.Context) ([]ReadCertDetails, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", GetCertEndpoint, nil)
	if err != nil {
		return nil, err
	}

	certData := &ReadCertDetailsResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, certData)
	if err != nil {
		return nil, err
	}

	return certData.CertResList, nil

}


func (prosimoClient *ProsimoClient) DeleteCert(ctx context.Context, certid string) error {
	DeleteCertEndpoint := fmt.Sprintf("%s/%s", GetCertEndpoint, certid)
	req, err := prosimoClient.api_client.NewRequest("DELETE", DeleteCertEndpoint, nil)
	if err != nil {
		return err
	}

	certData := &ReadCertDetailsResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, certData)
	if err != nil {
		return err
	}

	return nil

}
