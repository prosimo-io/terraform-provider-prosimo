package client

import (
	"context"
	"fmt"
	"log"
)

type AppOnboardSettingsOpts struct {
	AppID                   string          `json:"appID,omitempty"`
	App_Name                string          `json:"appName,omitempty"`
	AppOnboardType          string          `json:"appOnboardType,omitempty"`
	App_Access_Type         string          `json:"appAccessType,omitempty"`
	AppURLsOpts             []*AppURLOpts   `json:"appURLs,omitempty"`
	AppType                 string          `json:"appType,omitempty"`
	CitrixIP                []string        `json:"citrixIP,omitempty"`
	CloudService            string          `json:"cloudService,omitempty"`
	LogoURL                 string          `json:"logoURL,omitempty"`
	AppSamlRewrite          *AppSamlRewrite `json:"samlRewrite"`
	OnboardType             string          `json:"onboardType,omitempty"`
	InterActionType         string          `json:"interactionType,omitempty"`
	AddressType             string          `json:"addressType,omitempty"`
	ID                      string          `json:"id,omitempty"`
	Team_ID                 string          `json:"teamID,omitempty"`
	IDP_ID                  string          `json:"idpID,omitempty"`
	PolicyGroupID           string          `json:"policyGroupID,omitempty"`
	Optimize_App_Experience bool            `json:"optimizeAppExperience,omitempty"`
	OptOption               string          `json:"optOption,omitempty"`
	ClientCert              string          `json:"certID,omitempty"`
	EnableMultiCloud        bool            `json:"enableMultiCloud,omitempty"`
	PolicyName              []string        `json:"policyName,omitempty"`
	Custompolicy            *CustomPolicy   `json:"custompolicyName,omitempty"`
	Status                  string          `json:"status,omitempty"`
	Progress                int             `json:"progress,omitempty"`
	Deployed                bool            `json:"deployed,omitempty"`
	Source                  string          `json:"source,omitempty"`
	OnboardApp              bool            `json:"onboardApp,omitempty"`
	DecommissionApp         bool            `json:"decommissionApp,omitempty"`
	Dns_Discovery           bool            `json:"dnsDiscovery,omitempty"`
}

func (appOnboardSettingsOpts *AppOnboardSettingsOpts) GetAppOnboardSettings() *AppOnboardSettings {
	appOnboardSettings := &AppOnboardSettings{}
	appOnboardSettings.AppID = appOnboardSettingsOpts.AppID
	appOnboardSettings.App_Name = appOnboardSettingsOpts.App_Name
	appOnboardSettings.App_Access_Type = appOnboardSettingsOpts.App_Access_Type
	appOnboardSettings.AppType = appOnboardSettingsOpts.AppType
	appOnboardSettings.CitrixIP = appOnboardSettingsOpts.CitrixIP
	appOnboardSettings.CloudService = appOnboardSettingsOpts.CloudService
	appOnboardSettings.LogoURL = appOnboardSettingsOpts.LogoURL
	appOnboardSettings.AppSamlRewrite = appOnboardSettingsOpts.AppSamlRewrite
	appOnboardSettings.ID = appOnboardSettingsOpts.ID
	appOnboardSettings.Team_ID = appOnboardSettingsOpts.Team_ID
	appOnboardSettings.IDP_ID = appOnboardSettingsOpts.IDP_ID
	appOnboardSettings.PolicyGroupID = appOnboardSettingsOpts.PolicyGroupID
	appOnboardSettings.Optimize_App_Experience = appOnboardSettingsOpts.Optimize_App_Experience
	appOnboardSettings.OptOption = appOnboardSettingsOpts.OptOption
	appOnboardSettings.EnableMultiCloud = appOnboardSettingsOpts.EnableMultiCloud
	appOnboardSettings.Status = appOnboardSettingsOpts.Status
	appOnboardSettings.Progress = appOnboardSettingsOpts.Progress
	appOnboardSettings.Deployed = appOnboardSettingsOpts.Deployed
	appOnboardSettings.Source = appOnboardSettingsOpts.Source
	appOnboardSettings.OnboardType = appOnboardSettingsOpts.OnboardType
	appOnboardSettings.AddressType = appOnboardSettingsOpts.AddressType
	appOnboardSettings.InterActionType = appOnboardSettingsOpts.InterActionType
	appOnboardSettings.AppOnboardType = appOnboardSettingsOpts.AppOnboardType
	appOnboardSettings.Dns_Discovery = appOnboardSettingsOpts.Dns_Discovery

	appURLList := []*AppURL{}
	for _, appURLOpts := range appOnboardSettingsOpts.AppURLsOpts {
		appURLList = append(appURLList, appURLOpts.GetAppURL())
	}
	appOnboardSettings.AppURLs = appURLList

	return appOnboardSettings
}

type AppURLOpts struct {
	ID                string                     `json:"id,omitempty"`
	InternalDomain    string                     `json:"internalDomain,omitempty"`
	DomainType        string                     `json:"domainType,omitempty"`
	AppFqdn           string                     `json:"appFqdn,omitempty"`
	ServiceIpType     string                     `json:"serviceIPType,omitempty"`
	ServiceIp         string                     `json:"serviceIP,omitempty"`
	SubdomainIncluded bool                       `json:"subdomainIncluded,omitempty"`
	Protocols         []*AppProtocol             `json:"protocols,omitempty"`
	ExtProtocols      []*AppProtocol             `json:"extProtocols,omitempty"`
	HealthCheckInfo   *AppHealthCheckInfo        `json:"healthCheckInfo,omitempty"`
	CloudConfigOpts   *AppOnboardCloudConfigOpts `json:"cloudConfigOpts,omitempty"`
	DNSServiceOpts    *DNSServiceOpts            `json:"dnsServiceOpts,omitempty"`
	SSLCertOpts       *SSLCertOpts               `json:"sslCertOpts,omitempty"`
	WafPolicyName     string                     `json:"wafPolicyName,omitempty"`

	TeamID           string                          `json:"teamID,omitempty"`
	PappFqdn         string                          `json:"pappFqdn,omitempty"`
	CloudKeyID       string                          `json:"cloudKeyID,omitempty"`
	PrivateDcID      string                          `json:"dcID,omitempty"`
	DCAappIP         string                          `json:"dcAppIp,omitempty"`
	CertID           string                          `json:"certID,omitempty"`
	CacheRuleID      string                          `json:"cacheRuleID,omitempty"`
	CacheRuleName    string                          `json:"cacheRuleName,omitempty"`
	DNSService       *DNSService                     `json:"dnsService,omitempty"`
	ConnectionOption string                          `json:"connectionOption,omitempty"`
	Deployed         bool                            `json:"deployed,omitempty"`
	IP               string                          `json:"ip,omitempty"`
	Regions          []*AppOnboardCloudConfigRegions `json:"regions,omitempty"`
	DnsCustom        *AppOnboardDnsCustom            `json:"dnsCustom,omitempty"`
}

func (appURLOpts *AppURLOpts) GetAppURL() *AppURL {
	appURL := &AppURL{}
	appURL.ID = appURLOpts.ID
	appURL.InternalDomain = appURLOpts.InternalDomain
	appURL.DomainType = appURLOpts.DomainType
	appURL.AppFqdn = appURLOpts.AppFqdn
	appURL.ServiceIpType = appURLOpts.ServiceIpType
	appURL.ServiceIp = appURLOpts.ServiceIp
	appURL.SubdomainIncluded = appURLOpts.SubdomainIncluded
	appURL.Protocols = appURLOpts.Protocols
	appURL.ExtProtocols = appURLOpts.ExtProtocols
	appURL.HealthCheckInfo = appURLOpts.HealthCheckInfo
	appURL.TeamID = appURLOpts.TeamID
	appURL.PappFqdn = appURLOpts.PappFqdn
	appURL.CloudKeyID = appURLOpts.CloudKeyID
	appURL.PrivateDcID = appURLOpts.PrivateDcID
	appURL.DCAappIP = appURLOpts.DCAappIP
	appURL.CertID = appURLOpts.CertID
	appURL.CacheRuleID = appURLOpts.CacheRuleID
	appURL.DNSService = appURLOpts.DNSService
	appURL.ConnectionOption = appURLOpts.ConnectionOption
	appURL.Deployed = appURLOpts.Deployed
	appURL.IP = appURLOpts.IP
	appURL.Regions = appURLOpts.Regions
	appURL.DnsCustom = appURLOpts.DnsCustom

	return appURL
}

type AppOnboardSettingsResponseData struct {
	Data *AppOnboardSettings `json:"data,omitempty"`
}
type AppOnboardSearchResponse struct {
	Records    []*AppOnboardSettings `json:"records,omitempty"`
	TotalCount int                   `json:"totalCount,omitempty"`
}
type AppOnboardSearchResponseData struct {
	Data *AppOnboardSearchResponse `json:"data,omitempty"`
}

type AppOnboardSettings struct {
	AppID                   string          `json:"appID,omitempty"`
	App_Name                string          `json:"appName,omitempty"`
	AppOnboardType          string          `json:"appOnboardType,omitempty"`
	App_Access_Type         string          `json:"appAccessType,omitempty"`
	AppURLs                 []*AppURL       `json:"appURLs,omitempty"`
	AppType                 string          `json:"appType,omitempty"`
	CitrixIP                []string        `json:"citrixIP,omitempty"`
	CloudService            string          `json:"cloudService,omitempty"`
	LogoURL                 string          `json:"logoURL,omitempty"`
	AppSamlRewrite          *AppSamlRewrite `json:"samlRewrite"`
	OnboardType             string          `json:"onboardType,omitempty"`
	InterActionType         string          `json:"interactionType,omitempty"`
	AddressType             string          `json:"addressType,omitempty"`
	CertID                  string          `json:"certID,omitempty"`
	ID                      string          `json:"id,omitempty"`
	Team_ID                 string          `json:"teamID,omitempty"`
	IDP_ID                  string          `json:"idpID,omitempty"`
	PolicyGroupID           string          `json:"policyGroupID,omitempty"`
	PolicyIDs               []string        `json:"policyIDs,omitempty"`
	Optimize_App_Experience bool            `json:"optimizeAppExperience,omitempty"`
	OptOption               string          `json:"optOption,omitempty"`
	EnableMultiCloud        bool            `json:"enableMultiCloud,omitempty"`
	Status                  string          `json:"status,omitempty"`
	Progress                int             `json:"progress,omitempty"`
	Deployed                bool            `json:"deployed,omitempty"`
	Source                  string          `json:"source,omitempty"`
	Dns_Discovery           bool            `json:"dnsDiscovery"`
	PolicyUpdated           bool            `json:"policyUpdated"`
}
type AppURL struct {
	ID                string              `json:"id,omitempty"`
	InternalDomain    string              `json:"internalDomain,omitempty"`
	DomainType        string              `json:"domainType,omitempty"`
	AppFqdn           string              `json:"appFqdn,omitempty"`
	ServiceIpType     string              `json:"serviceIPType,omitempty"`
	ServiceIp         string              `json:"serviceIP,omitempty"`
	SubdomainIncluded bool                `json:"subdomainIncluded,omitempty"`
	Protocols         []*AppProtocol      `json:"protocols,omitempty"`
	ExtProtocols      []*AppProtocol      `json:"extProtocols,omitempty"`
	HealthCheckInfo   *AppHealthCheckInfo `json:"healthCheckInfo,omitempty"`
	WafHTTP           string              `json:"wafHTTP,omitempty"`

	TeamID           string                          `json:"teamID,omitempty"`
	PappFqdn         string                          `json:"pappFqdn,omitempty"`
	CloudKeyID       string                          `json:"cloudKeyID,omitempty"`
	PrivateDcID      string                          `json:"dcID,omitempty"`
	DCAappIP         string                          `json:"dcAppIp,omitempty"`
	CertID           string                          `json:"certID,omitempty"`
	CacheRuleID      string                          `json:"cacheRuleID,omitempty"`
	DNSService       *DNSService                     `json:"dnsService,omitempty"`
	ConnectionOption string                          `json:"connectionOption,omitempty"`
	Deployed         bool                            `json:"deployed,omitempty"`
	IP               string                          `json:"ip,omitempty"`
	Regions          []*AppOnboardCloudConfigRegions `json:"regions,omitempty"`
	DnsCustom        *AppOnboardDnsCustom            `json:"dnsCustom,omitempty"`
}

func (appURL *AppURL) String() string {
	return fmt.Sprintf("ID:%s, InternalDomain:%s, DomainType:%s, AppFqdn:%s, SubdomainIncluded:%t, PappFqdn:%s, CloudKeyID:%s, CertID:%s, CacheRuleID:%s, Deployed:%t, IP:%s",
		appURL.ID, appURL.InternalDomain, appURL.DomainType, appURL.AppFqdn, appURL.SubdomainIncluded, appURL.PappFqdn, appURL.CloudKeyID, appURL.CertID, appURL.CacheRuleID, appURL.Deployed, appURL.IP)
}

type AppProtocol struct {
	Protocol            string   `json:"protocol,omitempty"`
	Port                int      `json:"port,omitempty"`
	PortList            []string `json:"portList,omitempty"`
	WebSocketEnabled    bool     `json:"webSocketEnabled,omitempty"`
	IsValidProtocolPort bool     `json:"isValidProtocolPort,omitempty"`
	Paths               []string `json:"paths,omitempty"`
}

type AppHealthCheckInfo struct {
	Enabled  bool   `json:"enabled,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
}

type AppSamlRewrite struct {
	Metadata         string `json:"metadata,omitempty"`
	MetadataURL      string `json:"metadataURL,omitempty"`
	SelectedAuthType string `json:"selectedAuthType,omitempty"`
}

type DNSService struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
}

type DNSServiceOpts struct {
	Type           string `json:"type,omitempty"`
	CloudCredsName string `json:"name,omitempty"`
	AppURL         string `json:"appURL,omitempty"`
}

func (appSettings AppOnboardSettings) String() string {
	return fmt.Sprintf("{AppName:%s, AppAccessType:%s", appSettings.App_Name, appSettings.App_Access_Type)
}

// Step 1 : settings

func (prosimoClient *ProsimoClient) CreateAppOnboardSettings(ctx context.Context, appOnboardSettings *AppOnboardSettings) (*ResourcePostResponseData, error) {

	return prosimoClient.api_client.PostRequest(ctx, OnboardAppSettingsEndpoint, appOnboardSettings)

}

func (prosimoClient *ProsimoClient) UpdateAppOnboardSettings(ctx context.Context, appOnboardSettings *AppOnboardSettings) (*ResourcePostResponseData, error) {

	updateOnboardAppSettingsEndpoint := fmt.Sprintf(OnboardAppSettingsUpdateEndpoint, appOnboardSettings.ID, appOnboardSettings.ID)

	return prosimoClient.api_client.PutRequest(ctx, updateOnboardAppSettingsEndpoint, appOnboardSettings)

}

func (prosimoClient *ProsimoClient) GetAppOnboardSettings(ctx context.Context, appOnboardID string) (*AppOnboardSettings, error) {

	getAppOnboardSettingsEndpt := fmt.Sprintf(OnboardAppSettingsGetEndpoint, appOnboardID)

	req, err := prosimoClient.api_client.NewRequest("GET", getAppOnboardSettingsEndpt, nil)
	if err != nil {
		return nil, err
	}

	appOnboardSettingsResponseData := &AppOnboardSettingsResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, appOnboardSettingsResponseData)
	if err != nil {
		return nil, err
	}

	return appOnboardSettingsResponseData.Data, nil

}

func (prosimoClient *ProsimoClient) GetAppOnboardSummary(ctx context.Context, appOnboardID string) (*AppOnboardSettings, error) {

	getAppOnboardSummaryEndpt := fmt.Sprintf(OnboardAppSummaryGetEndpoint, appOnboardID)

	req, err := prosimoClient.api_client.NewRequest("GET", getAppOnboardSummaryEndpt, nil)
	if err != nil {
		return nil, err
	}

	appOnboardSettingsResponseData := &AppOnboardSettingsResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, appOnboardSettingsResponseData)
	if err != nil {
		return nil, err
	}

	return appOnboardSettingsResponseData.Data, nil

}

func (prosimoClient *ProsimoClient) DeleteAppOnboardSettings(ctx context.Context, appOnboardSettingsID string) error {

	deleteAppOnboardSettingsEndpt := fmt.Sprintf("%s/%s", OnboardAppSettingsEndpoint, appOnboardSettingsID)

	return prosimoClient.api_client.DeleteRequest(ctx, deleteAppOnboardSettingsEndpt)

}

// Step 1 : cloud configs

type AppOnboardCloudConfigOpts struct {
	AppURL                     string                             `json:"appURL,omitempty"`
	CloudCredsName             string                             `json:"cloudCredsName,omitempty"`
	ConnectionOption           string                             `json:"connectionOption,omitempty"`
	Regions                    []*AppOnboardCloudConfigRegionOpts `json:"regions,omitempty"`
	IsShowConnectionOptions    bool                               `json:"isShowConnectionOptions,omitempty"`
	HasPrivateConnectionOption bool                               `json:"hasPrivateConnectionOptions,omitempty"`
	AppHOstedType              string                             `json:"cloudType,omitempty"`
	DCAappIP                   string                             `json:"dcAppIp,omitempty"`
}

type AppOnboardCloudConfigRegionOpts struct {
	ID                       string               `json:"id,omitempty"`
	ModifyTgwAppRouteTable   string               `json:"modifyTgwAppRouteTable,omitempty"`
	AppnetworkID             string               `json:"appNetworkID,omitempty"`
	AttachPointID            string               `json:"attachPointID,omitempty"`
	Name                     string               `json:"name,omitempty"`
	ConnOption               string               `json:"connOption,omitempty"`
	RegionType               string               `json:"regionType,omitempty"`
	BackendIPAddressDiscover bool                 `json:"backendIPAddressDiscover,omitempty"`
	BackendIPAddressDns      bool                 `json:"backendIPAddressDns,omitempty"`
	DnsCustom                *AppOnboardDnsCustom `json:"dnsCustom,omitempty"`
	BackendIPAddressEntry    []string             `json:"backendIPAddressEntry,omitempty"`
	Buckets                  []string             `json:"buckets,omitempty"`
}

type AppOnboardCloudConfigRegions struct {
	ID                     string                             `json:"id,omitempty"`
	Name                   string                             `json:"name,omitempty"`
	Endpoints              []*AppOnboardCloudRegionEndpoints  `json:"endpoints,omitempty"`
	Attchpoints            []*AppOnboardCloudRegionAttchPoint `json:"attachPoints,omitempty"`
	ConnOption             string                             `json:"connOption,omitempty"`
	RegionType             string                             `json:"regionType,omitempty"`
	InputType              string                             `json:"inputType,omitempty"`
	ModifyTgwAppRouteTable string                             `json:"modifyTgwAppRouteTable,omitempty"`
	InputEndpoints         []string                           `json:"inputEndpoints,omitempty"`
	SelectedEndpoints      []string                           `json:"selectedEndpoints,omitempty"`
	LocationID             string                             `json:"locationID,omitempty"`
	Buckets                []string                           `json:"buckets,omitempty"`
	DnsCustom              *AppOnboardDnsCustom               `json:"dnsCustom,omitempty"`
}

func (cloudConfigRegions AppOnboardCloudConfigRegions) String() string {
	return fmt.Sprintf("{ID:%s, Name:%s, Endpoints:%v", cloudConfigRegions.ID, cloudConfigRegions.Name, cloudConfigRegions.Endpoints)
}

type AppOnboardCloudConfigRegionsResponseData struct {
	Data []*AppOnboardCloudConfigRegions `json:"data,omitempty"`
}

type AppOnboardCloudRegionEndpoints struct {
	AppIP                  string `json:"appIP,omitempty"`
	AppNetworkID           string `json:"appNetworkID,omitempty"`
	AttachPointDisplayName string `json:"attachPointDisplayName,omitempty"`
	AttachPointID          string `json:"attachPointID,omitempty"`
	AttachPointIP          string `json:"attachPointIP,omitempty"`
	DisplayName            string `json:"displayName,omitempty"`
}

type AppOnboardCloudRegionAttchPoint struct {
	AttachPointDisplayName string `json:"attachPointDisplayName,omitempty"`
	AttachPointID          string `json:"attachPointID,omitempty"`
}

type AppOnboardDnsCustom struct {
	DnsAppName           string   `json:"appname,omitempty"`
	DnsAppID             string   `json:"appId,omitempty"`
	DnsServers           []string `json:"servers,omitempty"`
	IsHealthCheckEnabled bool     `json:"isHealthCheckEnabled,omitempty"`
}

func (endpoints AppOnboardCloudRegionEndpoints) String() string {
	return fmt.Sprintf("{AppIP:%s, AppNetworkID:%s, DisplayName:%v", endpoints.AppIP, endpoints.AppNetworkID, endpoints.DisplayName)
}

func (prosimoClient *ProsimoClient) CreateAppOnboardCloudConfig(ctx context.Context, appOnboardSettingsID string, appOnboardSettings *AppOnboardSettings) (*ResourcePostResponseData, error) {

	postOnboardAppCloudConfigEndpoint := fmt.Sprintf(OnboardAppCloudEndpoint, appOnboardSettingsID)

	return prosimoClient.api_client.PutRequest(ctx, postOnboardAppCloudConfigEndpoint, appOnboardSettings)

}

func (prosimoClient *ProsimoClient) DiscoverAppOnboardEndpoint(ctx context.Context, appURL *AppURL, appOnboardType string) ([]*AppOnboardCloudConfigRegions, error) {
	appOnboardSettingsCloudConfigData := &AppOnboardSettings{}
	cloudConfigList := []*AppURL{}
	cloudConfigList = append(cloudConfigList, appURL)
	appOnboardSettingsCloudConfigData.AppURLs = cloudConfigList
	appOnboardSettingsCloudConfigData.AppOnboardType = appOnboardType
    
	req, err := prosimoClient.api_client.NewRequest("POST", OnboardAppEndpointDiscoveredEndpoint, appOnboardSettingsCloudConfigData)
	if err != nil {
		return nil, err
	}

	appOnboardCloudConfigRegionsResponseData := &AppOnboardCloudConfigRegionsResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, appOnboardCloudConfigRegionsResponseData)
	if err != nil {
		return nil, err
	}

	return appOnboardCloudConfigRegionsResponseData.Data, nil

}

// Step 3 : DNS Service

func (prosimoClient *ProsimoClient) CreateDNSService(ctx context.Context, appOnboardSettingsID string, appOnboardSettings *AppOnboardSettings) (*ResourcePostResponseData, error) {

	postOnboardAppDNSServiceEndpoint := fmt.Sprintf(OnboardAppDNSServiceEndpoint, appOnboardSettingsID)

	return prosimoClient.api_client.PutRequest(ctx, postOnboardAppDNSServiceEndpoint, appOnboardSettings)

}

// Step 4 : Generate Cert

type GenerateCert struct {
	CA                   string `json:"ca,omitempty"`
	TermsOfServiceAgreed bool   `json:"termsOfServiceAgreed,omitempty"`
	URL                  string `json:"url,omitempty"`
}

type SSLCertOpts struct {
	GenerateCert bool        `json:"generateCert,omitempty"`
	UploadCert   *UploadCert `json:"uploadCert,omitempty"`
	ExistingCert string
}

type CustomPolicy struct {
	Name string
}

func (prosimoClient *ProsimoClient) GenerateCert(ctx context.Context, generateCert *GenerateCert) (*ResourcePostResponseData, error) {

	return prosimoClient.api_client.PostRequest(ctx, GenerateCertEndpoint, generateCert)

}

func (prosimoClient *ProsimoClient) CreateAppOnboardCert(ctx context.Context, appOnboardSettingsID string, appOnboardSettings *AppOnboardSettings) (*ResourcePostResponseData, error) {

	postOnboardAppCertEndpoint := fmt.Sprintf(OnboardAppCertEndpoint, appOnboardSettingsID)

	return prosimoClient.api_client.PutRequest(ctx, postOnboardAppCertEndpoint, appOnboardSettings)

}

// Step 5 : Cache Rule

func (prosimoClient *ProsimoClient) CreateAppOnboardOptOption(ctx context.Context, appOnboardSettingsID string, appOnboardSettings *AppOnboardSettings) (*ResourcePostResponseData, error) {

	postOnboardAppOptOptionEndpoint := fmt.Sprintf(OnboardAppOptOptionEndpoint, appOnboardSettingsID)

	return prosimoClient.api_client.PutRequest(ctx, postOnboardAppOptOptionEndpoint, appOnboardSettings)

}

// Step 6 : Policy and WAF

type AppOnboardPolicy struct {
	PolicyIDs []string `json:"policyIDs,omitempty"`
}

func (prosimoClient *ProsimoClient) CreateAppOnboardPolicy(ctx context.Context, appOnboardSettingsID string, appOnboardPolicy *AppOnboardPolicy) (*ResourcePostResponseData, error) {

	postOnboardAppPolicyEndpoint := fmt.Sprintf(OnboardAppPolicyEndpoint, appOnboardSettingsID)

	return prosimoClient.api_client.PutRequest(ctx, postOnboardAppPolicyEndpoint, appOnboardPolicy)

}

func (prosimoClient *ProsimoClient) CreateAppOnboardWaf(ctx context.Context, appOnboardSettingsID string, appOnboardSettings *AppOnboardSettings) (*ResourcePostResponseData, error) {

	postOnboardAppWafEndpoint := fmt.Sprintf(OnboardAppWafEndpoint, appOnboardSettingsID)

	return prosimoClient.api_client.PutRequest(ctx, postOnboardAppWafEndpoint, appOnboardSettings)

}

// Step 7 : Onboard App

func (prosimoClient *ProsimoClient) OnboardAppDeployment(ctx context.Context, appOnboardSettingsID string, appOnboardSettings *AppOnboardSettings) (*ResourcePostResponseData, error) {

	postOnboardAppDeploymentEndpoint := fmt.Sprintf(OnboardAppDeploymentEndpoint, appOnboardSettingsID)

	return prosimoClient.api_client.PutRequest(ctx, postOnboardAppDeploymentEndpoint, appOnboardSettings)

}
func (prosimoClient *ProsimoClient) OnboardAppDeploymentV2(ctx context.Context, appOnboardSettings *AppOnboardSettings, paramValue string) (*ResourcePostResponseData, error) {
	// log.Println("Entering reboard block", networkOnboardsops)
	log.Println("Inside app onboard v2")

	OnboardAppkDeploymentEndpoint := fmt.Sprintf("%s?%s=%s", AppOnboardEndpointNew, ParamName, paramValue)

	// onboardrequest := &NetworkDeploymentres{}
	req, err := prosimoClient.api_client.NewRequest("POST", OnboardAppkDeploymentEndpoint, appOnboardSettings)
	if err != nil {
		return nil, err
	}

	onboardresponse := &ResourcePostResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, onboardresponse)
	if err != nil {
		return nil, err
	}

	return onboardresponse, nil

}

func (prosimoClient *ProsimoClient) OffboardAppDeployment(ctx context.Context, appOnboardSettingsID string) (*ResourcePostResponseData, error) {

	deleteOnboardAppDeploymentEndpoint := fmt.Sprintf(OnboardAppDeploymentEndpoint, appOnboardSettingsID)

	return prosimoClient.api_client.AppOffboarding_DeleteRequest(ctx, deleteOnboardAppDeploymentEndpoint)

}

func (prosimoClient *ProsimoClient) ForceOffboardAppDeployment(ctx context.Context, appOnboardSettingsID string) (*ResourcePostResponseData, error) {

	deleteOnboardAppDeploymentEndpoint := fmt.Sprintf(OnboardAppDeploymentEndpoint, appOnboardSettingsID)

	return prosimoClient.api_client.AppOffboarding_Force_DeleteRequest(ctx, deleteOnboardAppDeploymentEndpoint)

}

type AppDomainData struct {
	Data []*AppDomain `json:"data,omitempty"`
}

type AppDomain struct {
	Domain     string `json:"domain,omitempty"`
	ID         string `json:"id,omitempty"`
	WafDetails string `json:"wafDetails,omitempty"`
}

func (prosimoClient *ProsimoClient) GetAppDomains(ctx context.Context) ([]*AppDomain, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", AppDomainEndpoint, nil)
	if err != nil {
		return nil, err
	}

	appDomainData := &AppDomainData{}
	_, err = prosimoClient.api_client.Do(ctx, req, appDomainData)
	if err != nil {
		return nil, err
	}

	return appDomainData.Data, nil

}

type AppOnboardSearch struct {
	StatusFilter   []string `json:"statusFilter,omitempty"`
	NetworkService string   `json:"networkService,omitempty"`
	Category       string   `json:"category,omitempty"`
}

func (prosimoClient *ProsimoClient) SearchAppOnboardApps(ctx context.Context, AppOnboardSearchInput *AppOnboardSearch) (*AppOnboardSearchResponseData, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", OnboardAppSearchEndpoint, AppOnboardSearchInput)
	if err != nil {
		return nil, err
	}

	appOnboardSettingsResponseData := &AppOnboardSearchResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, appOnboardSettingsResponseData)
	if err != nil {
		return nil, err
	}

	return appOnboardSettingsResponseData, nil

}
func (prosimoClient *ProsimoClient) DeleteApp(ctx context.Context, appOnboardSettingsID string) error {

	deleteOnboardAppDeploymentEndpoint := fmt.Sprintf("%s/%s", OnboardAppEndpoint, appOnboardSettingsID)

	return prosimoClient.api_client.DeleteRequest(ctx, deleteOnboardAppDeploymentEndpoint)

}
