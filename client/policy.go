package client

import (
	"context"
	"errors"
	"fmt"
)

type PolicyListData struct {
	Records    []*Policy `json:"records,omitempty"`
	TotalCount int       `json:"totalCount,omitempty"`
}
type PolicyListDataResponse struct {
	Data *PolicyListData `json:"data,omitempty"`
}

type PolicytData struct {
	Data *Policy `json:"data,omitempty"`
}

type Policy struct {
	PolicyType                string  `json:"type,omitempty"`
	ID                        string  `json:"id,omitempty"`
	Name                      string  `json:"name,omitempty"`
	TeamID                    string  `json:"teamID,omitempty"`
	App_Access_Type           string  `json:"appAccessType,omitempty"`
	NamespaceID               string  `json:"namespaceID,omitempty"`
	Details                   Details `json:"details,omitempty"`
	Device_Posture_Configured bool    `json:"devicePostureConfigured,omitempty"`
	Type                      string  `json:"type,omitempty"`
	DisplayName               string  `json:"displayName,omitempty"`
}
type Details struct {
	Actions  []string        `json:"actions,omitempty"`
	Matches  MatchDetailList `json:"matches,omitempty"`
	Apps     Values          `json:"apps,omitempty"`
	Networks Values          `json:"networks,omitempty"`
}

type Details_ds struct {
	Actions  []string        `json:"actions,omitempty"`
	Matches  MatchDetailList `json:"matches,omitempty"`
	Apps     *Values         `json:"apps,omitempty"`
	Networks *Values         `json:"networks,omitempty"`
}
type MatchDetailList struct {
	UserList               []MatchDetails `json:"users,omitempty"`
	Networks               []MatchDetails `json:"networks,omitempty"`
	ProsimoNetworks        []MatchDetails `json:"prosimoNetworks,omitempty"`
	NetworkACL             []MatchDetails `json:"networkACL,omitempty"`
	IDP                    []MatchDetails `json:"idp,omitempty"`
	Devices                []MatchDetails `json:"devices,omitempty"`
	Time                   []MatchDetails `json:"time,omitempty"`
	Advanced               []MatchDetails `json:"advanced,omitempty"`
	FQDN                   []MatchDetails `json:"fqdn,omitempty"`
	URL                    []MatchDetails `json:"url,omitempty"`
	Device_Posture_Profile []MatchDetails `json:"devicePostureProfiles,omitempty"`
}
type MatchDetails struct {
	Property   string `json:"property,omitempty"`
	Operations string `json:"operation,omitempty"`
	Values     Values `json:"values,omitempty"`
}

type Values struct {
	InputItems     []InputItems `json:"inputItems,omitempty"`
	SelectedItems  []InputItems `json:"selectedItems,omitempty"`
	SelectedGroups []InputItems `json:"selectedGroups,omitempty"`
}
type InputItems struct {
	ItemID          string     `json:"id,omitempty"`
	ItemName       string     `json:"name,omitempty"`
	CityCode        int        `json:"cityCode,omitempty"`
	CityName        string     `json:"cityName,omitempty"`
	CountryCodeISO2 string     `json:"countryCodeISO2,omitempty"`
	Region          string     `json:"region,omitempty"`
	StateName       string     `json:"stateName,omitempty"`
	CountryName     string     `json:"countryName,omitempty"`
	KeyValues       *KeyValues `json:"keyValues,omitempty"`
}

type KeyValues struct {
	SourceIp   []string `json:"sourceIP,omitempty"`
	TargetIp   []string `json:"targetIP,omitempty"`
	Protocol   []string `json:"protocol,omitempty"`
	SourcePort []string `json:"sourcePort,omitempty"`
	TargetPort []string `json:"targetPort,omitempty"`
}

type MatchItem struct {
	Users                  PropertyItemList `json:"Users,omitempty"`
	Location               PropertyItemList `json:"Location,omitempty"`
	IDP                    PropertyItemList `json:"IDP,omitempty"`
	Devices                PropertyItemList `json:"Devices,omitempty"`
	Time                   PropertyItemList `json:"Time,omitempty"`
	URL                    PropertyItemList `json:"URL,omitempty"`
	FQDN                   PropertyItemList `json:"fqdn,omitempty"`
	Advanced               PropertyItemList `json:"Advanced,omitempty"`
	Device_Posture_Profile PropertyItemList `json:"Device_Posture_Profile,omitempty"`
	Networks               PropertyItemList `json:"Networks,omitempty"`
	NetworkACL             PropertyItemList `json:"NetworkACL,omitempty"`
}
type PropertyItemList struct {
	Property []PropertyItem
}
type PropertyItem struct {
	User_Property   string          `json:"user_property,omitempty"`
	Server_Property string          `json:"server_property,omitempty"`
	Operations      []OperationItem `json:"operations,omitempty"`
}
type OperationItem struct {
	User_Operation_Name   string `json:"user_operation_name,omitempty"`
	Server_Operation_Name string `json:"server_operation_name,omitempty"`
}

func (prosimoClient *ProsimoClient) ReadJson() MatchItem {

	return GetPolicyServerTemplate()
}

func (prosimoClient *ProsimoClient) GetPolicy(ctx context.Context) ([]*Policy, error) {

	PolicySearchInput := Policy{}
	req, err := prosimoClient.api_client.NewRequest("POST", GetPolicyEndpoint, PolicySearchInput)
	if err != nil {
		return nil, err
	}

	policyListData := &PolicyListDataResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, policyListData)
	if err != nil {
		return nil, err
	}

	return policyListData.Data.Records, nil

}
func (prosimoClient *ProsimoClient) GetInternetEgressControl(ctx context.Context) ([]*Policy, error) {

	PolicySearchInput := Policy{}
	req, err := prosimoClient.api_client.NewRequest("POST", GetIECEndpoint , PolicySearchInput)
	if err != nil {
		return nil, err
	}

	policyListData := &PolicyListDataResponse{}
	_, err = prosimoClient.api_client.Do(ctx, req, policyListData)
	if err != nil {
		return nil, err
	}

	return policyListData.Data.Records, nil

}
func (prosimoClient *ProsimoClient) GetPolicyByName(ctx context.Context, policyName string) (*Policy, error) {

	policyList, err := prosimoClient.GetPolicy(ctx)
	if err != nil {
		return nil, err
	}

	var policy *Policy
	for _, returnedPolicy := range policyList {
		if returnedPolicy.Name == policyName {
			policy = returnedPolicy
			break
		}
	}

	if policy == nil {
		return nil, errors.New("Policy doesnt exists")
	}

	return policy, nil

}

func (prosimoClient *ProsimoClient) GetInternetEgressControlByName(ctx context.Context, policyName string) (*Policy, error) {

	policyList, err := prosimoClient.GetInternetEgressControl(ctx)
	if err != nil {
		return nil, err
	}

	var policy *Policy
	for _, returnedPolicy := range policyList {
		if returnedPolicy.Name == policyName {
			policy = returnedPolicy
			break
		}
	}

	if policy == nil {
		return nil, errors.New("Policy doesnt exists")
	}

	return policy, nil

}

func (prosimoClient *ProsimoClient) GetPolicyID(ctx context.Context, policyName string) (string, error) {

	policyList, err := prosimoClient.GetPolicy(ctx)
	if err != nil {
		return "", err
	}

	var policy *Policy
	for _, returnedPolicy := range policyList {
		if returnedPolicy.Name == policyName {
			policy = returnedPolicy
			break
		}
	}

	if policy == nil {
		return "", errors.New("Policy doesnt exists")
	}

	return policy.ID, nil

}

func (prosimoClient *ProsimoClient) GetPolicyByID(ctx context.Context, policyID string) (*Policy, error) {

	policyList, err := prosimoClient.GetPolicy(ctx)
	if err != nil {
		return nil, err
	}

	var policy *Policy
	for _, returnedPolicy := range policyList {
		if returnedPolicy.ID == policyID {
			policy = returnedPolicy
			break
		}
	}
	if policy == nil {
		return nil, errors.New("Policy doesnt exists")
	}

	return policy, nil

}

func (prosimoClient *ProsimoClient) CreatePolicy(ctx context.Context, inputPolicy *Policy) (*PolicytData, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", PolicyEndpoint, inputPolicy)
	if err != nil {
		return nil, err
	}

	policyListData := &PolicytData{}
	_, err = prosimoClient.api_client.Do(ctx, req, policyListData)
	if err != nil {
		return nil, err
	}

	return policyListData, nil

}

func (prosimoClient *ProsimoClient) UpdatePolicy(ctx context.Context, inputPolicy *Policy) (*PolicytData, error) {

	updatePolicyEndpoint := fmt.Sprintf("%s/%s", PolicyEndpoint, inputPolicy.ID)
	req, err := prosimoClient.api_client.NewRequest("PUT", updatePolicyEndpoint, inputPolicy)
	if err != nil {
		return nil, err
	}
	policyListData := &PolicytData{}

	_, err = prosimoClient.api_client.Do(ctx, req, policyListData)
	if err != nil {
		return nil, err
	}
	return policyListData, nil
}

func (prosimoClient *ProsimoClient) DeletePolicy(ctx context.Context, policyID string) error {

	DeletePolicyEndpoint := fmt.Sprintf("%s/%s", PolicyEndpoint, policyID)
	req, err := prosimoClient.api_client.NewRequest("DELETE", DeletePolicyEndpoint, nil)
	if err != nil {
		return err
	}

	//policyListData := &PolicytData{}
	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}

type AppDetailsData struct {
	GetApp []Record `json:"data,omitempty"`
}

type Record struct {
	ID            string `json:"id,omitempty"`
	TeamID        string `json:"teamID,omitempty"`
	IDPID         string `json:"idpID,omitempty"`
	AppAccessType string `json:"appAccessType,omitempty"`
	PolicyGroupID string `json:"policyGroupID,omitempty"`
	AppName       string `json:"appName,omitempty"`
	Status        string `json:"astatus,omitempty"`
}

func (prosimoClient *ProsimoClient) GetAppID(ctx context.Context, appName string) (string, error) {
	var appId string
	req, err := prosimoClient.api_client.NewRequest("GET", AppListEndpoint, nil)
	if err != nil {
		return "", err
	}
	applist := &AppDetailsData{}
	_, err = prosimoClient.api_client.Do(ctx, req, applist)

	if err != nil {
		return "", err
	}
	fmt.Println("applist", applist)
	for _, val := range applist.GetApp {
		if val.AppName == appName {
			appId = val.ID
			fmt.Println("appId", appId)
		}
	}

	return appId, nil
}

func (prosimoClient *ProsimoClient) GetNetworkID(ctx context.Context, networkName string) (string, error) {
	var networkId string
	req, err := prosimoClient.api_client.NewRequest("GET", NetworkListEndpoint, nil)
	if err != nil {
		return "", err
	}
	networklist := &NetworkOnboardResList{}
	_, err = prosimoClient.api_client.Do(ctx, req, networklist)

	if err != nil {
		return "", err
	}
	fmt.Println("networklist", networklist)
	for _, val := range networklist.NetworkOnboardResponseDataList {
		if val.Name == networkName {
			networkId = val.ID
			fmt.Println("appId", networkId)
		}
	}

	return networkId, nil
}

type LocationSearchPayload struct {
	Value           string `json:"value,omitempty"`
	CountryCodeISO2 string `json:"countryCodeISO2,omitempty"`
}

type GeoLocRes struct {
	LocDetails []InputItems `json:"data,omitempty"`
}

func (prosimoClient *ProsimoClient) FetchGeoLocation(ctx context.Context, payload *LocationSearchPayload) ([]InputItems, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", GeoLocEndpoint, payload)
	if err != nil {
		return nil, err
	}

	geoLocData := &GeoLocRes{}
	_, err = prosimoClient.api_client.Do(ctx, req, geoLocData)
	if err != nil {
		return nil, err
	}

	return geoLocData.LocDetails, nil

}
