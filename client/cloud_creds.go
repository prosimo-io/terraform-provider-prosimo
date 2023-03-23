package client

import (
	"context"
	"errors"
	"fmt"
)

type CloudCredsDetails struct {
	IAMRoleArn    string `json:"roleArn,omitempty"`
	IAMExternalID string `json:"externalID,omitempty"`
	AccessKeyID   string `json:"accessKeyID,omitempty"`
	SecretKeyID   string `json:"secretKeyID,omitempty"`

	SubscriptionID string `json:"subscriptionID,omitempty"`
	TenantID       string `json:"tenantID,omitempty"`
	ClientID       string `json:"clientID,omitempty"`
	SecretID       string `json:"clientSecret,omitempty"`

	GcpType       string `json:"type,omitempty"`
	ProjectID     string `json:"project_id,omitempty"`
	PrivateKey    string `json:"private_key,omitempty"`
	PrivateKeyID  string `json:"private_key_id,omitempty"`
	ClientEmail   string `json:"client_email,omitempty"`
	GcpClientID   string `json:"client_id,omitempty"`
	AuthURI       string `json:"auth_uri,omitempty"`
	TokenURI      string `json:"token_uri,omitempty"`
	AuthCertURL   string `json:"auth_provider_x509_cert_url,omitempty"`
	ClientCertURL string `json:"client_x509_cert_url,omitempty"`
}

/*type binaryCloudDetails struct {
	ClDetails []byte(CloudCredsDetails
}*/

type CloudCreds struct {
	CloudType         string             `json:"cloudType,omitempty"`
	Nickname          string             `json:"name,omitempty"`
	KeyType           string             `json:"keyType,omitempty"`
	CloudCredsDetails *CloudCredsDetails `json:"details,omitempty"`
	ID                string             `json:"id,omitempty"`
	ConectionType     string             `json:"connectionType,omitempty"`
	//binaryCloudDetails *binaryCloudDetails
}

type CloudCredsList struct {
	CloudCreds []*CloudCreds `json:"data,omitempty"`
}

type CloudCredsData struct {
	CloudCreds *CloudCreds `json:"data,omitempty"`
}

type CloudRegionData struct {
	CloudRegionList []*CloudRegion `json:"data,omitempty"`
}

type CloudRegion struct {
	AppCount   int    `json:"appCount,omitempty"`
	Region     string `json:"regionName,omitempty"`
	LocationID string `json:"locationID,omitempty"`
}

func (cloudCreds CloudCreds) String() string {
	return fmt.Sprintf("{CloudType:%s, Nickname:%s, KeyType:%v, CloudCredsDetails:%v}", cloudCreds.CloudType, cloudCreds.Nickname,
		cloudCreds.KeyType, cloudCreds.CloudCredsDetails)
}

func (cloudCredsDetails CloudCredsDetails) String() string {
	return fmt.Sprintf("{AWS ---> IAMRoleArn:%s, IAMExternalID:%s, AccessKeyID:%s, SecretKeyID:%s, AZURE --->   SubscriptionID:%s, TenantID:%s, ClientID:%s, SecretID:%s, GCP ---> ClientEmail:%s, GcpType:%s, ProjectID:%s, PrivateKey:%s, PrivateKeyID:%s, GcpClientID:%s, AuthURI:%s, TokenURI:%s, AuthCertURL:%s, ClientCertURL:%s}",
		cloudCredsDetails.IAMRoleArn, cloudCredsDetails.IAMExternalID, cloudCredsDetails.AccessKeyID, cloudCredsDetails.SecretKeyID,
		cloudCredsDetails.SubscriptionID, cloudCredsDetails.TenantID, cloudCredsDetails.ClientID, cloudCredsDetails.SecretID,
		cloudCredsDetails.ClientEmail, cloudCredsDetails.GcpType, cloudCredsDetails.ProjectID, cloudCredsDetails.PrivateKey, cloudCredsDetails.PrivateKeyID, cloudCredsDetails.GcpClientID, cloudCredsDetails.AuthURI, cloudCredsDetails.TokenURI, cloudCredsDetails.AuthCertURL, cloudCredsDetails.ClientCertURL)
}

func (prosimoClient *ProsimoClient) CreateCloudCreds(ctx context.Context, cloudCredsOpts *CloudCreds) (*CloudCredsData, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", CloudCredsEndpoint, cloudCredsOpts)
	if err != nil {
		return nil, err
	}

	cloudCredsData := &CloudCredsData{}
	_, err = prosimoClient.api_client.Do(ctx, req, cloudCredsData)
	if err != nil {
		return nil, err
	}

	return cloudCredsData, nil

}

func (prosimoClient *ProsimoClient) UpdateCloudCreds(ctx context.Context, cloudCredsOpts *CloudCreds) (*CloudCredsData, error) {

	updateCloudCredsEndpt := fmt.Sprintf("%s/%s", CloudCredsEndpoint, cloudCredsOpts.ID)

	req, err := prosimoClient.api_client.NewRequest("PUT", updateCloudCredsEndpt, cloudCredsOpts)
	if err != nil {
		return nil, err
	}

	cloudCredsData := &CloudCredsData{}
	_, err = prosimoClient.api_client.Do(ctx, req, cloudCredsData)
	if err != nil {
		return nil, err
	}

	return cloudCredsData, nil

}

func (prosimoClient *ProsimoClient) GetCloudCreds(ctx context.Context) (*CloudCredsList, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", CloudCredsEndpoint, nil)
	if err != nil {
		return nil, err
	}

	cloudCredsList := &CloudCredsList{}
	_, err = prosimoClient.api_client.Do(ctx, req, cloudCredsList)
	if err != nil {
		return nil, err
	}

	return cloudCredsList, nil

}

func (prosimoClient *ProsimoClient) GetCloudCredsByName(ctx context.Context, cloudCredsName string) (*CloudCreds, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", CloudCredsEndpoint, nil)
	if err != nil {
		return nil, err
	}

	cloudCredsList := &CloudCredsList{}
	_, err = prosimoClient.api_client.Do(ctx, req, cloudCredsList)
	if err != nil {
		return nil, err
	}

	var cloudCreds *CloudCreds
	for _, returnedCloudCreds := range cloudCredsList.CloudCreds {
		if returnedCloudCreds.Nickname == cloudCredsName {
			cloudCreds = returnedCloudCreds
			break
		}
	}

	if cloudCreds == nil {
		return nil, errors.New("Cloudcreds doesnt exists")
	}

	return cloudCreds, nil

}

func (prosimoClient *ProsimoClient) GetCloudCredsById(ctx context.Context, cloudCredsId string) (*CloudCreds, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", CloudCredsEndpoint, nil)
	if err != nil {
		return nil, err
	}

	cloudCredsList := &CloudCredsList{}
	_, err = prosimoClient.api_client.Do(ctx, req, cloudCredsList)
	if err != nil {
		return nil, err
	}

	var cloudCreds *CloudCreds
	for _, returnedCloudCreds := range cloudCredsList.CloudCreds {
		if returnedCloudCreds.ID == cloudCredsId {
			cloudCreds = returnedCloudCreds
			break
		}
	}

	if cloudCreds == nil {
		return nil, errors.New("Cloudcreds doesnt exists")
	}

	return cloudCreds, nil

}

func (prosimoClient *ProsimoClient) DeleteCloudCreds(ctx context.Context, cloudCredsID string) error {

	deleteCloudCredsEndpt := fmt.Sprintf("%s/%s", CloudCredsEndpoint, cloudCredsID)

	req, err := prosimoClient.api_client.NewRequest("DELETE", deleteCloudCredsEndpt, nil)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}

func (prosimoClient *ProsimoClient) GetCloudRegion(ctx context.Context, cloudCredsID string) (*CloudRegionData, error) {

	getCloudRegionEndpt := fmt.Sprintf(CloudRegionEndpoint, cloudCredsID)

	req, err := prosimoClient.api_client.NewRequest("GET", getCloudRegionEndpt, nil)
	if err != nil {
		return nil, err
	}

	cloudRegionData := &CloudRegionData{}
	_, err = prosimoClient.api_client.Do(ctx, req, cloudRegionData)
	if err != nil {
		return nil, err
	}

	return cloudRegionData, nil

}

func (prosimoClient *ProsimoClient) CheckIfCloudRegionExists(ctx context.Context, cloudCredsID string, region string) (bool, error) {

	getCloudRegionEndpt := fmt.Sprintf(CloudRegionEndpoint, cloudCredsID)

	req, err := prosimoClient.api_client.NewRequest("GET", getCloudRegionEndpt, nil)
	if err != nil {
		return false, err
	}

	cloudRegionData := &CloudRegionData{}
	_, err = prosimoClient.api_client.Do(ctx, req, cloudRegionData)
	if err != nil {
		return false, err
	}

	regionExists := false
	for _, cloudRegion := range cloudRegionData.CloudRegionList {
		if cloudRegion.Region == region {
			regionExists = true
			break
		}
	}

	return regionExists, nil

}

func (prosimoClient *ProsimoClient) UploadGcpCloudCreds(ctx context.Context, cloudCredsOpts *CloudCreds, Filepath string) (*CloudCredsData, error) {

	fileUploadEdPt := fmt.Sprintf("%s/%s", CloudCredsEndpoint, "file")
	m := make(map[string]string)

	m["cloudType"] = cloudCredsOpts.CloudType
	m["keyType"] = cloudCredsOpts.KeyType
	m["name"] = cloudCredsOpts.Nickname
	req, err := prosimoClient.api_client.ReqFileUpload("POST", fileUploadEdPt, m, Filepath)
	if err != nil {
		return nil, err
	}

	cloudCredsData := &CloudCredsData{}
	_, err = prosimoClient.api_client.Do(ctx, req, cloudCredsData)
	if err != nil {
		return nil, err
	}

	return cloudCredsData, nil

}

func (prosimoClient *ProsimoClient) UpdateGcpCloudCreds(ctx context.Context, cloudCredsOpts *CloudCreds, Filepath string) (*CloudCredsData, error) {

	fileUploadEdPt := fmt.Sprintf("%s/%s", CloudCredsEndpoint, cloudCredsOpts.ID)
	updatefileUploadEdPt := fmt.Sprintf("%s/%s", fileUploadEdPt, "file")
	meta_data := make(map[string]string)

	meta_data["cloudType"] = cloudCredsOpts.CloudType
	meta_data["keyType"] = cloudCredsOpts.KeyType
	meta_data["name"] = cloudCredsOpts.Nickname
	req, err := prosimoClient.api_client.ReqFileUpload("POST", updatefileUploadEdPt, meta_data, Filepath)
	if err != nil {
		return nil, err
	}

	cloudCredsData := &CloudCredsData{}
	_, err = prosimoClient.api_client.Do(ctx, req, cloudCredsData)
	if err != nil {
		return nil, err
	}

	return cloudCredsData, nil

}
