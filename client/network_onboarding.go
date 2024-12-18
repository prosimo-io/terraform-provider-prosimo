package client

import (
	"context"
	"fmt"
	"net/url"
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
	ID            string        `json:"id,omitempty"`
	Name          string        `json:"name,omitempty"`
	Exportable    bool          `json:"exportable,omitempty"`
	NamespaceID   string        `json:"namespaceID,omitempty"`
	NamespaceNID  int           `json:"namespaceNID,omitempty"`
	TeamID        string        `json:"teamID,omitempty"`
	PamCname      string        `json:"pamCname,omitempty"`
	Deployed      bool          `json:"deployed,omitempty"`
	Status        string        `json:"status,omitempty"`
	CreatedTime   string        `json:"createdTime,omitempty"`
	UpdatedTime   string        `json:"updatedTime,omitempty"`
	PolicyUpdated bool          `json:"policyUpdated,omitempty"`
	NamespaceName string        `json:"namespaceName,omitempty"`
	PublicCloud   *PublicCloud  `json:"publicCloud,omitempty"`
	PrivateCloud  *PrivateCloud `json:"privateCloud,omitempty"`
	Security      *Security     `json:"security,omitempty"`
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
	Bandwidth      string          `json:"bandwidth,omitempty"`
	BandwidthName  string          `json:"bandwidthName,omitempty"`
	CloudNetworkID string          `json:"cloudNetworkId,omitempty"`
	UpdateStatus   string          `json:"updateStatus,omitempty"`
	InstanceType   string          `json:"instanceType,omitempty"`
	Subnets        []string        `json:"subnets,omitempty"`
	BandwidthRange *BandwidthRange `json:"bandwidthRange,omitempty"`
}

type BandwidthRange struct {
	Min int `json:"min,omitempty"`
	Max int `json:"max,omitempty"`
}

type ServiceSubnets struct {
	Mode string `json:"mode,omitempty"`
}

type CloudNetworkops struct {
	Id                 string             `json:"id,omitempty"`
	CloudNetworkID     string             `json:"cloudNetworkID,omitempty"`
	ConnectorGrpID     string             `json:"connectorGroupID,omitempty"`
	EdgeConnectivityID string             `json:"edgeConnectivityID,omitempty"`
	Subnets            []InputSubnet      `json:"subnets,omitempty"`
	HubID              string             `json:"hubID,omitempty"`
	ConnectivityType   string             `json:"connectivityType,omitempty"`
	ConnectorPlacement string             `json:"connectorPlacement,omitempty"`
	Servicesubnets     *ServiceSubnets    `json:"serviceSubnets,omitempty"`
	Connectorsettings  *ConnectorSettings `json:"connectorSettings,omitempty"`
}

type InputSubnet struct {
	Subnet        string `json:"subnet,omitempty"`
	VirtualSubnet string `json:"virtualSubnet,omitempty"`
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
	Policies               []Policyops `json:"policies,omitempty"`
	InternetEgressControls []Policyops `json:"internetEgressControls,omitempty"`
}

type NetworkSecurityInput struct {
	Security *Security `json:"security,omitempty"`
}

type NetworkDeploymentresops struct {
	TaskID  string `json:"taskID,omitempty"`
	AuditID string `json:"auditID,omitempty"`
	ID      string `json:"id,omitempty"`
}

type SubnetRead struct {
	ID   string `json:"id,omitempty"`
	CIDR string `json:"cidr,omitempty"`
	Name string `json:"name,omitempty"`
}
type VPC struct {
	ID           string       `json:"id,omitempty"`
	Network      string       `json:"network,omitempty"`
	Name         string       `json:"name,omitempty"`
	Subnet_Count int          `json:"subnet_count,omitempty"`
	Subnets      []SubnetRead `json:"subnets,omitempty"`
}

type VPCLIST struct {
	VpcList []VPC `json:"data,omitempty"`
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


func (prosimoClient *ProsimoClient) OnboardNetworkDeploymentV2(ctx context.Context, networkOnboardsops *NetworkOnboardoptns, paramValue string) (*NetworkDeploymentres, error) {
	// log.Println("Entering reboard block", networkOnboardsops)

	OnboardNetworkDeploymentEndpoint := fmt.Sprintf("%s?%s=%s", NetworkOnboardEndpointNew, ParamName, paramValue)

	// onboardrequest := &NetworkDeploymentres{}
	req, err := prosimoClient.api_client.NewRequest("POST", OnboardNetworkDeploymentEndpoint, networkOnboardsops)
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

func (prosimoClient *ProsimoClient) ForceOffboardNetworkDeployment(ctx context.Context, appID string) (*NetworkDeploymentres, error) {
	PostOnboardNetworkDeploymentEndpoint := fmt.Sprintf(OnboardNetworkDeploymentEndpoint, appID)
	u, err := url.Parse(PostOnboardNetworkDeploymentEndpoint)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return nil, err
	}

	// Add query parameters
	q := u.Query()
	q.Add("force", "true")
	u.RawQuery = q.Encode()

	req, err := prosimoClient.api_client.NewRequest("DELETE", u.String(), nil)
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

func (prosimoClient *ProsimoClient) GetNetworkList(ctx context.Context, cloudKey string, region string) ([]VPC, error) {
	updateDiscoveredNetworkVpc := fmt.Sprintf(DiscoveredNetworkVpc, cloudKey, region)

	req, err := prosimoClient.api_client.NewRequest("GET", updateDiscoveredNetworkVpc, nil)
	if err != nil {
		return nil, err
	}
	readresponse := &VPCLIST{}
	_, err = prosimoClient.api_client.Do(ctx, req, readresponse)
	if err != nil {
		return nil, err
	}

	return readresponse.VpcList, nil

}
