package client

import (
	"context"
	"fmt"
	"log"
)

type NetworkDiscovery struct {
	ID      string `json:"id,omitempty"`
	Region  string `json:"region,omitempty"`
	VpcID   string `json:"vpcID,omitempty"`
	VpcCIDR string `json:"vpcCIDR,omitempty"`
}

type NetworkDiscoveryResponse struct {
	THlist []string `json:"data,omitempty"`
}

type NetworkOnboardoptns struct {
	ID           string        `json:"id,omitempty"`
	Name         string        `json:"name,omitempty"`
	Exportable   bool          `json:"exportable,omitempty"`
	NamespaceID  string        `json:"namespaceID,omitempty"`
	TeamID       string        `json:"teamID,omitempty"`
	PamCname     string        `json:"pamCname,omitempty"`
	Deployed     bool          `json:"deployed,omitempty"`
	Status       string        `json:"status,omitempty"`
	PublicCloud  *PublicCloud  `json:"publicCloud,omitempty"`
	PrivateCloud *PrivateCloud `json:"privateCloud,omitempty"`
	Security     *Security     `json:"security,omitempty"`
}

type NetworkOnboardRes struct {
	NetworkOnboardResponseData *NetworkOnboardoptns `json:"data,omitempty"`
}
type NetworkOnboardResList struct {
	NetworkOnboardResponseDataList []*NetworkOnboardoptns `json:"data,omitempty"`
}
type NetworkOnboardSearchResponse struct {
	Records    []*NetworkOnboardoptns `json:"records,omitempty"`
	TotalCount int                    `json:"totalCount,omitempty"`
}
type NetworkOnboardSearchResponseData struct {
	Data *NetworkOnboardSearchResponse `json:"data,omitempty"`
}
type ConnectorSettings struct {
	Bandwidth     string   `json:"bandwidth,omitempty"`
	BandwidthName string   `json:"bandwidthName,omitempty"`
	InstanceType  string   `json:"instanceType,omitempty"`
	Subnets       []string `json:"subnets,omitempty"`
}
type ServiceSubnets struct {
	Mode string `json:"mode,omitempty"`
}

type CloudNetworkops struct {
	Id                 string             `json:"id,omitempty"`
	CloudNetworkID     string             `json:"cloudNetworkID,omitempty"`
	ConnectorGrpID     string             `json:"connectorGroupID,omitempty"`
	EdgeConnectivityID string             `json:"edgeConnectivityID,omitempty"`
	Subnets            []string           `json:"subnets,omitempty"`
	HubID              string             `json:"hubID,omitempty"`
	ConnectivityType   string             `json:"connectivityType,omitempty"`
	ConnectorPlacement string             `json:"connectorPlacement,omitempty"`
	Servicesubnets     *ServiceSubnets    `json:"serviceSubnets,omitempty"`
	Connectorsettings  *ConnectorSettings `json:"connectorSettings,omitempty"`
}

type PublicCloud struct {
	Id               string            `json:"id,omitempty"`
	Cloud            string            `json:"cloud,omitempty"`
	CloudType        string            `json:"cloudType,omitempty"`
	ConnectionOption string            `json:"connectionOption,omitempty"`
	CloudKeyID       string            `json:"cloudKeyID,omitempty"`
	CloudRegion      string            `json:"cloudRegion,omitempty"`
	CloudNetworks    []CloudNetworkops `json:"cloudNetworks,omitempty"`
	ConnectType      string            `json:"connectType,omitempty"`
}

type PrivateCloud struct {
	Id               string   `json:"id,omitempty"`
	CloudType        string   `json:"cloudType,omitempty"`
	ConnectionOption string   `json:"connectionOption,omitempty"`
	PrivateCloudID   string   `json:"privateCloudID,omitempty"`
	Subnets          []string `json:"subnets,omitempty"`
}

type Policyops struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Security struct {
	Policies []Policyops `json:"policies,omitempty"`
}

type NetworkSecurityInput struct {
	Security *Security `json:"security,omitempty"`
}

type NetworkDeploymentresops struct {
	TaskID  string `json:"taskID,omitempty"`
	AuditID string `json:"auditID,omitempty"`
}

type NetworkDeploymentres struct {
	NetworkDeploymentResops NetworkDeploymentresops `json:"data,omitempty"`
}

func (prosimoClient *ProsimoClient) GetTransitHub(ctx context.Context, networkCloudTransitHubDiscovery *NetworkDiscovery) (*NetworkDiscoveryResponse, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", CloudTransitHubDiscoveryEndpoint, networkCloudTransitHubDiscovery)
	if err != nil {
		return nil, err
	}

	cloudTHListData := &NetworkDiscoveryResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, cloudTHListData)
	if err != nil {
		return nil, err
	}

	return cloudTHListData, nil

}

func (prosimoClient *ProsimoClient) GetSubsnets(ctx context.Context, networkCloudSubnetDiscovery *NetworkDiscovery) (*NetworkDiscoveryResponse, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", CloudSubnetDiscoveryEndpoint, networkCloudSubnetDiscovery)
	if err != nil {
		return nil, err
	}

	cloudSubnetListData := &NetworkDiscoveryResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, cloudSubnetListData)
	if err != nil {
		return nil, err
	}

	return cloudSubnetListData, nil

}

func (prosimoClient *ProsimoClient) NetworkOnboard(ctx context.Context, networkOnboardsops *NetworkOnboardoptns) (*NetworkOnboardRes, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", OnboardNetworkEndpoint, networkOnboardsops)
	if err != nil {
		return nil, err
	}

	onboardresponse := &NetworkOnboardRes{}
	_, err = prosimoClient.api_client.Do(ctx, req, onboardresponse)
	if err != nil {
		return nil, err
	}

	return onboardresponse, nil

}

func (prosimoClient *ProsimoClient) NetworkOnboardCloud(ctx context.Context, networkOnboardsops *NetworkOnboardoptns) error {
	postOnboardNetworkCloudEndpoint := fmt.Sprintf(OnboardNetworkCloudEndpoint, networkOnboardsops.ID)

	req, err := prosimoClient.api_client.NewRequest("PUT", postOnboardNetworkCloudEndpoint, networkOnboardsops)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}

func (prosimoClient *ProsimoClient) NetworkOnboardSecurity(ctx context.Context, networkOnboardsecurityops *NetworkSecurityInput, appID string) error {
	postOnboardNetworkPolicyEndpoint := fmt.Sprintf(OnboardNetworkPolicyEndpoint, appID)

	req, err := prosimoClient.api_client.NewRequest("PUT", postOnboardNetworkPolicyEndpoint, networkOnboardsecurityops)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}

func (prosimoClient *ProsimoClient) OnboardNetworkDeployment(ctx context.Context, appID string) (*NetworkDeploymentres, error) {
	PostOnboardNetworkDeploymentEndpoint := fmt.Sprintf(OnboardNetworkDeploymentEndpoint, appID)

	onboardrequest := &NetworkDeploymentres{}
	req, err := prosimoClient.api_client.NewRequest("PUT", PostOnboardNetworkDeploymentEndpoint, onboardrequest)
	if err != nil {
		return nil, err
	}

	onboardresponse := &NetworkDeploymentres{}
	_, err = prosimoClient.api_client.Do(ctx, req, onboardresponse)
	if err != nil {
		return nil, err
	}

	return onboardresponse, nil

}

func (prosimoClient *ProsimoClient) OnboardNetworkDeploymentPost(ctx context.Context, networkOnboardsops *NetworkOnboardoptns) (*NetworkDeploymentres, error) {
	log.Println("Entering reboard block", networkOnboardsops)
	PostOnboardNetworkDeploymentEndpoint := fmt.Sprintf(OnboardNetworkDeploymentEndpoint, networkOnboardsops.ID)

	// onboardrequest := &NetworkDeploymentres{}
	req, err := prosimoClient.api_client.NewRequest("POST", PostOnboardNetworkDeploymentEndpoint, networkOnboardsops)
	if err != nil {
		return nil, err
	}

	onboardresponse := &NetworkDeploymentres{}
	_, err = prosimoClient.api_client.Do(ctx, req, onboardresponse)
	if err != nil {
		return nil, err
	}

	return onboardresponse, nil

}

func (prosimoClient *ProsimoClient) OffboardNetworkDeployment(ctx context.Context, appID string) (*NetworkDeploymentres, error) {
	PostOnboardNetworkDeploymentEndpoint := fmt.Sprintf(OnboardNetworkDeploymentEndpoint, appID)

	req, err := prosimoClient.api_client.NewRequest("DELETE", PostOnboardNetworkDeploymentEndpoint, nil)
	if err != nil {
		return nil, err
	}

	onboardresponse := &NetworkDeploymentres{}
	_, err = prosimoClient.api_client.Do(ctx, req, onboardresponse)
	if err != nil {
		return nil, err
	}

	return onboardresponse, nil

}

func (prosimoClient *ProsimoClient) DeleteNetworkDeployment(ctx context.Context, appID string) error {
	PostNetworkOnboardEndpoint := fmt.Sprintf(NetworkOnboardEndpoint, appID)

	req, err := prosimoClient.api_client.NewRequest("DELETE", PostNetworkOnboardEndpoint, nil)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}

func (prosimoClient *ProsimoClient) GetNetworkSettings(ctx context.Context, appID string) (*NetworkOnboardoptns, error) {
	PostNetworkOnboardEndpoint := fmt.Sprintf(NetworkOnboardEndpoint, appID)

	req, err := prosimoClient.api_client.NewRequest("GET", PostNetworkOnboardEndpoint, nil)
	if err != nil {
		return nil, err
	}
	onboardresponse := &NetworkOnboardRes{}
	_, err = prosimoClient.api_client.Do(ctx, req, onboardresponse)
	if err != nil {
		return nil, err
	}

	return onboardresponse.NetworkOnboardResponseData, nil

}

func (prosimoClient *ProsimoClient) SearchOnboardNetworks(ctx context.Context) (*NetworkOnboardSearchResponseData, error) {

	NetworkOnboardSearchInput := NetworkOnboardoptns{}
	req, err := prosimoClient.api_client.NewRequest("POST", OnboardNetworkSearchEndpoint, NetworkOnboardSearchInput)
	if err != nil {
		return nil, err
	}

	networkOnboardSettingsResponseData := &NetworkOnboardSearchResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, networkOnboardSettingsResponseData)
	if err != nil {
		return nil, err
	}

	return networkOnboardSettingsResponseData, nil

}
