package client

import (
	"context"
	"errors"
	"fmt"
)

type PL_Source struct {
	Name         string          `json:"name,omitempty"`
	Region       string          `json:"region,omitempty"`
	ID           string          `json:"id,omitempty"`
	Deleted      bool            `json:"deleted,omitempty"`
	InUse        bool            `json:"inUse,omitempty"`
	TotalCount   int             `json:"totalCount,omitempty"`
	Status       string          `json:"status,omitempty"`
	CloudCredsID string          `json:"cloudCredsID,omitempty"`
	Credentials  *Credentials    `json:"credentials,omitempty"`
	CloudSources *[]Cloud_Source `json:"cloudSources,omitempty"`
	Ports        *[]PL_Port      `json:"ports,omitempty"`
	Edge         *PL_Edge        `json:"edge,omitempty"`
	Policies     *[]PL_Policy    `json:"policies,omitempty"`
	CreatedTime  string          `json:"createdTime,omitempty"`
	UpdatedTime  string          `json:"updatedTime,omitempty"`
}

type Credentials struct {
	ID          string `json:"id,omitempty"`
	Cloud       string `json:"cloud,omitempty"`
	Type        string `json:"type,omitempty"`
	Credentials string `json:"credentials,omitempty"`
}

type Cloud_Source struct {
	ID           string                 `json:"id,omitempty"`
	CloudNetwork *Cloud_Network         `json:"cloudNetwork,omitempty"`
	Deleted      bool                   `json:"deleted,omitempty"`
	Subnets      *[]Subnet              `json:"subnets,omitempty"`
	HostedZones  *[]HostedZone          `json:"hostedzones,omitempty"`
	Records      *[]Cloud_Source_Record `json:"records,omitempty"`
	CreatedTime  string                 `json:"createdTime,omitempty"`
	UpdatedTime  string                 `json:"updatedTime,omitempty"`
}

type PL_Port struct {
	ID          string `json:"id,omitempty"`
	TeamID      string `json:"teamID,omitempty"`
	EdgeID      string `json:"edgeID,omitempty"`
	Port        int    `json:"port,omitempty"`
	Status      string `json:"status,omitempty"`
	PlsID       string `json:"plsID,omitempty"`
	DomainID    string `json:"domainID,omitempty"`
	AppID       string `json:"appID,omitempty"`
	CretedTime  string `json:"createdTime,omitempty"`
	UpdatedTime string `json:"udpdatedTime,omitempty"`
}

type PL_Edge struct {
	ID           string          `json:"id,omitempty"`
	Domain       string          `json:"domain,omitempty"`
	Cloud        string          `json:"cloud,omitempty"`
	CloudCredsID string          `json:"cloudCredsID,omitempty"`
	Name         string          `json:"name,omitempty"`
	NetworkInfo  *PL_NetworkInfo `json:"networkInfo,omitempty"`
	Region       string          `json:"region,omitempty"`
}

type PL_NetworkInfo struct {
	VpcID  string `json:"vpcId,omitempty"`
	IlbDns string `json:"ilbDns,omitempty"`
}

type PL_Policy struct {
	ID     string `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
	Target string `json:"target,omitempty"`
}

type Subnet struct {
	ID                 string         `json:"id,omitempty"`
	Cidr               string         `json:"cidr,omitempty"`
	OnboardedASNetwork bool           `json:"onboardedAsNetwork,omitempty"`
	Endpoints          *[]PL_Endpoint `json:"endpoints,omitempty"`
	Deleted            bool           `json:"deleted,omitempty"`
	CreatedTime        string         `json:"createdTime,omitempty"`
	UpdatedTime        string         `json:"updatedTime,omitempty"`
}

type PL_Endpoint struct {
	ID         string                   `json:"id,omitempty"`
	Endpoint   string                   `json:"endpoint,omitempty"`
	Name       string                   `json:"name,omitempty"`
	DomainID   string                   `json:"domainID,omitempty"`
	Domain     string                   `json:"domain,omitempty"`
	PolicyID   string                   `json:"policyID,omitempty"`
	AppID      string                   `json:"appID,omitempty"`
	AppName    string                   `json:"appName,omitempty"`
	Status     string                   `json:"status,omitempty"`
	ProtoPorts *[]PL_Endpoint_ProtoPort `json:"protoPorts,omitempty"`
}

type PL_Endpoint_ProtoPort struct {
	Protocol         string   `json:"protocol,omitempty"`
	Port             int      `json:"port,omitempty"`
	WebSocketEnabled bool     `json:"webSocketEnabled,omitempty"`
	Paths            []string `json:"paths,omitempty"`
	PortList         []string `json:"portList,omitempty"`
}

type HostedZone struct {
	ID                 string `json:"id,omitempty"`
	TeamID             string `json:"teamID,omitempty"`
	CloudCredentialsID string `json:"cloudCredentialsID,omitempty"`
	HostedZoneID       string `json:"hostedZoneID,omitempty"`
	Name               string `json:"name,omitempty"`
	VpcID              string `json:"vpcID,omitempty"`
	Region             string `json:"region,omitempty"`
	ProsimoManaged     bool   `json:"prosimoManaged,omitempty"`
	Status             string `json:"status,omitempty"`
	CreatedTime        string `json:"createdTime,omitempty"`
	UpdatedTime        string `json:"updatedTime,omitempty"`
}

type Cloud_Source_Record struct {
	ID           string `json:"id,omitempty"`
	TeamID       string `json:"teamID,omitempty"`
	Record       string `json:"record,omitempty"`
	HostedZoneID string `json:"hostedZoneID,omitempty"`
	Type         string `json:"type,omitempty"`
	Target       string `json:"target,omitempty"`
	Ttl          int    `json:"ttl,omitempty"`
	Status       string `json:"status,omitempty"`
	CreatedTime  string `json:"createdTime,omitempty"`
	UpdatedTime  string `json:"updatedTime,omitempty"`
}

type Subnet_res struct {
	Subnets []Subnet `json:"data,omitempty"`
}
type PL_Source_ResID struct {
	ID string `json:"id,omitempty"`
}

type Cloud_Network struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Cloud_Network_res struct {
	CloudNetworks []Cloud_Network `json:"data,omitempty"`
}
type PL_Source_Response struct {
	PLSource_ResID *PL_Source_ResID `json:"data,omitempty"`
}

type PLSListData struct {
	Records    []*PL_Source `json:"records,omitempty"`
	TotalCount int          `json:"totalCount,omitempty"`
}

type PLSListDataResponse struct {
	Data *PLSListData `json:"data,omitempty"`
}
type Search_Input struct {
	ID           string `json:"id,omitempty"`
	CloudCredsID string `json:"cloudCredsID,omitempty"`
	Region       string `json:"region,omitempty"`
}

func (prosimoClient *ProsimoClient) CreatePrivateLinkSource(ctx context.Context, privateLinkSourceInput *PL_Source) (*PL_Source_Response, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", PrivateLinkSourceEndpoint, privateLinkSourceInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &PL_Source_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) UpdatePrivateLinkSource(ctx context.Context, privateLinkSourceInput *PL_Source) (*PL_Source_Response, error) {

	updatePrivateLinkSourceEndpoint := fmt.Sprintf("%s/%s", PrivateLinkSourceEndpoint, privateLinkSourceInput.ID)

	req, err := prosimoClient.api_client.NewRequest("PUT", updatePrivateLinkSourceEndpoint, privateLinkSourceInput)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &PL_Source_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) DeletePrivateLinkSource(ctx context.Context, privateLinkSourceID string) (*PL_Source_Response, error) {

	updatePrivateLinkSourceEndpoint := fmt.Sprintf("%s/%s", PrivateLinkSourceEndpoint, privateLinkSourceID)

	req, err := prosimoClient.api_client.NewRequest("DELETE", updatePrivateLinkSourceEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &PL_Source_Response{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) GetPrivateLinkSource(ctx context.Context) ([]*PL_Source, error) {

	PLSSearchInput := PL_Source{}
	req, err := prosimoClient.api_client.NewRequest("POST", GetPrivateLinkSourceEndpoint, PLSSearchInput)
	if err != nil {
		return nil, err
	}

	plsListData := &PLSListDataResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, plsListData)
	if err != nil {
		return nil, err
	}

	return plsListData.Data.Records, nil

}

func (prosimoClient *ProsimoClient) GetPrivateLinkSourceByID(ctx context.Context, id string) (*PL_Source, error) {

	plsList, err := prosimoClient.GetPrivateLinkSource(ctx)
	if err != nil {
		return nil, err
	}
	var privateLinkSource *PL_Source
	for _, returnedPLS := range plsList {
		if returnedPLS.ID == id {
			privateLinkSource = returnedPLS
		}
	}

	if privateLinkSource == nil {
		return nil, errors.New("private Link Source doesn't exists")
	}

	return privateLinkSource, nil

}
func (prosimoClient *ProsimoClient) GetPrivateLinkSourceByName(ctx context.Context, name string) (*PL_Source, error) {

	plsList, err := prosimoClient.GetPrivateLinkSource(ctx)
	if err != nil {
		return nil, err
	}
	var privateLinkSource *PL_Source
	for _, returnedPLS := range plsList {
		if returnedPLS.Name == name {
			privateLinkSource = returnedPLS
		}
	}

	if privateLinkSource == nil {
		return nil, errors.New("private Link Source doesn't exists")
	}

	return privateLinkSource, nil

}

func (prosimoClient *ProsimoClient) DiscoverPVSNetworks(ctx context.Context, discoverNetworkInput *Search_Input) (*Cloud_Network_res, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", DiscoverPVSNetworkEndpoint, discoverNetworkInput)
	if err != nil {
		return nil, err
	}

	discoveredNetwork := &Cloud_Network_res{}
	_, err = prosimoClient.api_client.Do(ctx, req, discoveredNetwork)
	if err != nil {
		return nil, err
	}

	return discoveredNetwork, nil

}

func (prosimoClient *ProsimoClient) DiscoverPVSSubnets(ctx context.Context, discoverSubnetInput *Search_Input) (*Subnet_res, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", DiscoverPVSSubnetEndpoint, discoverSubnetInput)
	if err != nil {
		return nil, err
	}

	discoveredSubnets := &Subnet_res{}
	_, err = prosimoClient.api_client.Do(ctx, req, discoveredSubnets)
	if err != nil {
		return nil, err
	}

	return discoveredSubnets, nil

}
