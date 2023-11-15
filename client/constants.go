package client

// Prosimo API constants
const (
	AWSCloudType     = "AWS"
	AzureCloudType   = "AZURE"
	GCPCloudType     = "GCP"
	PrivateCloudType = "PRIVATE"

	AWSKeyType     = "AWSKEY"
	AzureKeyType   = "AZUREKEY"
	GCPKeyType     = "GCPKEY"
	PrivateKeyType = "PRIVATE"

	TypeWEB       = "web"
	TypeFQDN      = "fqdn"
	TypeIP        = "ip"
	TypeCloudSvc  = "cloud-svc"
	TypeJumpBox   = "jumpbox"
	TypeCitrixVDI = "vdi"

	PatternEqual       = "=="
	PatternEqualNoCase = "=*"
	PatternNotEqual    = "!="
	PatternContains    = "=@"
	PatternOR          = ","
	PatternAND         = "&"

	ModifyTgwAppRoutetableType   = "MODIFY"
	MaintainTgwAppRoutetableType = "MAINTAIN"

	PrivateConnection_type = "PRIVATE"
	PublicConnection_type  = "PUBLIC"

	AWSIAMRoleAuth   = "AWSIAM"
	AWSAccessKeyAuth = "AWSKEY"

	IDPOktaAccount         = "okta"
	IDPAzureADAccount      = "azure_ad"
	IDPOneLoginAccount     = "one_login"
	IDPOtherAccount        = "other"
	IDPPingOneAccount      = "ping-one"
	IDPPingFederateAccount = "ping-federate"
	IDPGoogleWSAccount     = "google"

	IDPOIDCAuth = "oidc"
	IDPSAMLAuth = "saml"

	APICredYes = "yes"
	APICredNo  = "no"

	AddTgwAction = "ADD"
	ModTgwAction = "MOD"
	DelTgwAction = "DEL"

	EdgeTgwConnection = "EDGE"
	VPCTgwConnection  = "VPC"
	VNETTgwConnection = "VNET"

	IDPEURegion      = "eu"
	IDPAsiaRegion    = "asia"
	IDPDefaultRegion = "default"
	IDPUSRegion      = "us"

	IDPPrimaryType = "primary"
	IDPPartnerType = "partner"

	AgentAppAccessType     = "agent"
	AgentlessAppAccessType = "agentless"

	DefaultAppType    = "default"
	CitrixAppType     = "citrix"
	TCPAppType        = "tcp"
	CloudPlatformType = "cloud-svc"
	JumpBoxType       = "jumpbox"

	CustomAppDomain  = "custom"
	ProsimoAppDomain = "prosimo"

	SamlAuth  = "saml"
	OidcAuth  = "oidc"
	OtherAuth = "other"

	AutoServiceIP   = "auto"
	ManualServiceIP = "manual"

	HTTPAppProtocol  = "http"
	HTTPSAppProtocol = "https"
	SSHAppProtocol   = "ssh"
	RDPAppProtocol   = "rdp"
	VNCAppProtocol   = "vnc"
	TCPAppProtocol   = "tcp"
	UDPAppProtocol   = "udp"

	HTTPU2AALWebAppProtocol  = "http"
	HTTPSU2AALWebAppProtocol = "https"
	SSHU2AALWebAppProtocol   = "ssh"
	RDPU2AALWebAppProtocol   = "rdp"
	VNCU2AALWebAppProtocol   = "vnc"

	TCPfqdnAppProtocol = "tcp"
	UDPfqdnAppProtocol = "udp"
	DNSfqdnAppProtocol = "dns"

	DRNameAlert    = "alert"
	DRNameMfa      = "mfa"
	DRNameLockUser = "lockUser"

	PolicyTransit = "transit"
	PolicyAccess  = "access"

	PolicyActionAllow   = "allow"
	PolicyActionDeny    = "deny"
	PolicyActionBypass  = "bypass"
	PolicyActionSkipWAF = "skipWAF"
	PolicyActionAlert   = "alert"
	PolicyActionMFA     = "mfa"

	PolicyMatchUser        = "users"
	PolicyMatchNetwork     = "networks"
	PolicyMatchDevice      = "devices"
	PolicyMatchTime        = "time"
	PolicyMatchUrl         = "url"
	PolicyMatchApplication = "fqdn"
	PolicyMatchAppAdvanced = "advanced"
	PolicyMatchAppIDP      = "idp"
	PolicyMatchAppDP       = "device-posture"
	PolicyMatchPosition    = "location"
	PolicyMatchNetworkACL  = "networkacl"

	PolicyMatchOperationIs             = "Is"
	PolicyMatchOperationIsNot          = "Is NOT"
	PolicyMatchOperationContains       = "Contains"
	PolicyMatchOperationDoesNotContain = "Does NOT contain"
	PolicyMatchOperationStartsWith     = "Starts with"
	PolicyMatchOperationEndsWith       = "Ends with"
	PolicyMatchOperationIn             = "In"
	PolicyMatchOperationNotIn          = "NOT in"
	PolicyMatchOperationIsAtLeast      = "Is at least"
	PolicyMatchOperationBetween        = "Between"

	CacheTimeUnitDays    = "Days"
	CacheTimeUnitHours   = "Hours"
	CacheTimeUnitMinutes = "Minutes"
	CacheTimeUnitSeconds = "Seconds"

	CacheAPITimeUnitDays    = "days"
	CacheAPITimeUnitHours   = "hours"
	CacheAPITimeUnitMinutes = "minutes"
	CacheAPITimeUnitSeconds = "seconds"

	CacheTypeDynamic          = "Dynamic"
	CacheTypeStaticLongLived  = "Static - Long lived"
	CacheTypeStaticShortLived = "Static - Short lived"

	CacheAPIInputDynamic          = "dynamic"
	CacheAPIInputStaticLongLived  = "static-long-lived"
	CacheAPIInputStaticShortLived = "static-short-lived"

	DefaultPolicyNetwork = "DEFAULT-MCN-POLICY"

	ExistingVPCSource = "Existing"
	ProsimoVPCSource  = "Prosimo"

	AzureBandwidth     = "small"
	AzureBandwidthName = "<1 Gbps"
	AzureInstanceType  = "Standard_A2_v2"

	LessThan1GBPS   = "<1 Gbps"
	OneToFiveGBPS   = "1-5 Gbps"
	FiveToTenGBPS   = "5-10 Gbps"
	MoreThanTenGBPS = ">10 Gbps"

	AWST3Medium     = "t3.medium"
	AWsT3aMedium    = "t3a.medium"
	AWSC5Large      = "c5.large"
	AWSC5aLarge     = "c5a.large"
	AWSC5xLarge     = "c5.xlarge"
	AWSC5axLarge    = "c5a.xlarge"
	AWSC5nxLarge    = "c5n.xlarge"
	AWSC5a8xLarge   = "c5a.8xlarge"
	AWSC59xLarge    = "c5.9xlarge"
	AWSC5n9xLarge   = "c5n.9xlarge"
	AWSC5a16xLarge  = "c5a.16xlarge"
	AWSC518xLarge   = "c5.18xlarge"
	AWSC5n18xLarge  = "c5n.18xlarge"
	GCPE2Standard2  = "e2-standard-2"
	GCPE2Standard4  = "e2-standard-4"
	GCPE2Standard8  = "e2-standard-8"
	GCPE2Standard16 = "e2-standard-16"
	GCPECStandard16 = "c2-standard-16"

	ConnectorSizeSmall      = "small"
	ConnectorSizeMedium     = "medium"
	ConnectorSizeLarge      = "large"
	ConnectorSizeExtraLarge = "extra-large"

	Optnpeering          = "peering"
	OptnpeeringInput     = "private" //parse it to private
	OptnawsPrivateLink   = "awsPrivateLink"
	OptnawsVpnGateway    = "awsVpnGateway"
	OptnazurePrivateLink = "azurePrivateLink"
	OptntransitGateway   = "transitGateway"
	OptnazureTransitVnet = "azureTransitVnet"
	OptnvwanHub          = "vwanHub"
	Optnpublic           = "public"
	Optnprivate          = "private"

	HostedPrivate = "PRIVATE"
	HostedPublic  = "PUBLIC"

	PublicCloudConnectionOption  = "public"
	PrivateCloudConnectionOption = "private"

	WorkloadVpcConnectorPlacementOptions = "Workload VPC"
	InfraVPCConnectorPlacementOptions    = "Infra VPC"
	AppConnectorPlacementOptions         = "app"
	NoneConnectorPlacementOptions        = "none"
	InfraConnectorPlacementOptions       = "infra"

	AutoServiceInsertionOptions   = "auto"
	NoneServiceInsertionOptions   = "none"
	ManualServiceInsertionOptions = "manual"

	VpcPeeringConnectivityType = "vpc-peering"
	TGWConnectivityType = "transit-gateway"
	PublicConnectivityType = "public"
	VnetPeeringConnectivityType = "vnet-peering"
	VwanHubConnectivityType = "vwan-hub"


	ManualDNSServiceType     = "manual"
	AWSRoute53DNSServiceType = "aws_route53"
	ProsimoDNSServiceType    = "prosimo"

	PerformanceEnhancedOptOption = "PerformanceEnhanced"
	CostSavingOptOption          = "CostSaving"
	FastLaneOptOption            = "FastLane"

	AllowAllPolicyOption        = "ALLOW-ALL-USERS"
	AllowAllPolicyOptionNetwork = "ALLOW-ALL-NETWORKS"
	DenyAllPolicyOption         = "DENY-ALL-USERS"
	DenyAllPolicyOptionNetwork  = "DENY-ALL-NETWORKS"
	CustomizePolicyOption       = "CUSTOMIZE"

	BehindFabricAppOnboardType    = "behind_fabric"
	AccessingFabricAppOnboardType = "accessing_fabric"
	BothAppOnboardType            = "both"
	BothAppOnboardServerType      = "behind_accessing_fabric"

	UserToAppInteractionType = "userToApp"
	AppToAppInteractionType  = "appToApp"

	FQDNAddessType = "fqdn"
	IPAddressType  = "ip_address"

	ActiveRegionType = "active"
	BackUpRegionType = "backup"

	WafModeEnforce = "enforce"
	WafModeDetect  = "detect"

	TypeCrowdStrike = "CrowdStrike"

	ApiCrowdStrike          = "https://api.crowdstrike.com"
	ApiUSCrowdStrike        = "https://api.us-2.crowdstrike.com"
	ApiEUCrowdStrike        = "https://api.eu-1.crowdstrike.com"
	ApiLaggerGCWCrowdStrike = "https://api.laggar.gcw.crowdstrike.com"

	TypeEnabled  = "enabled"
	TypeDisabled = "disabled"
	TypeNa       = "na"

	RiskLevelHigh   = "high"
	RiskLevelMedium = "medium"
	RiskLevelLow    = "low"

	OsWindows = "windows"
	OsMac     = "mac"

	TypeIS        = "is"
	TypeISNOT     = "is-not"
	TypeISATLEAST = "is-atleast"

	WafRuleSetBasic = "Basic rule-set"
	WafRuleSetOWASP = "OWASP Modsecurity Core Ruleset v3.2"

	TypeUSER     = "USER"
	TypeAPP      = "APP"
	TypeDEVICE   = "DEVICE"
	TypeTIME     = "TIME"
	TypeIP_RANGE = "IP_RANGE"
	TypeGEO      = "GEO"

	APIPrefix = "api/"

	APITokenEndpoint                     = APIPrefix + "token"
	RoleEndpoint                         = APIPrefix + "role"
	IPPoolEndpoint                       = APIPrefix + "ippool"
	CloudCredsEndpoint                   = APIPrefix + "cloud/creds"
	CloudCredsEndpointPrivate            = APIPrefix + "cloud/private"
	TaskEndpoint                         = APIPrefix + "task"
	CloudS3Endpoint                      = APIPrefix + "cloud/service/amazon/s3/discovery"
	SharedServiceEndpoint                = APIPrefix + "shared-svc"
	SharedServiceDeploymentEndpoint      = APIPrefix + "shared-svc/deployment"
	GetSharedServiceEndpoint             = APIPrefix + "shared-svc/search"
	ServiceInsertionEndpoint             = APIPrefix + "svc-insert"
	GetServiceInsertionEndpoint          = APIPrefix + "svc-insert/search"
	PrivateLinkSourceEndpoint            = APIPrefix + "private-link-source"
	GetPrivateLinkSourceEndpoint         = APIPrefix + "private-link-source/search"
	DiscoverPVSNetworkEndpoint           = APIPrefix + "private-link-source/discovery/networks"
	DiscoverPVSSubnetEndpoint            = APIPrefix + "private-link-source/discovery/subnets"
	PrivateLinkMappingEndpoint           = APIPrefix + "private-link-source/policy"
	GetPrivateLinkMappingEndpoint        = APIPrefix + "private-link-source/policy/search"
	EdgeEndpoint                         = APIPrefix + "prosimo/app"
	PatchEdgeEndpoint                    = APIPrefix + "prosimo/app"
	PatchEdgeSubnetEndpoint              = APIPrefix + "subnet/cluster"
	AppDeploymentEndpoint                = APIPrefix + "prosimo/app/deployment"
	CloudRegionEndpoint                  = APIPrefix + "cloud/%s/region"
	CloudVpcDiscoveryEndpoint            = APIPrefix + "cloud/vpc/discovery"
	CloudEPDiscoveryEndpoint             = APIPrefix + "cloud/endpoint/discovery"
	CloudTransitHubDiscoveryEndpoint     = APIPrefix + "cloud/transit-hub/discovery"
	CloudSubnetDiscoveryEndpoint         = APIPrefix + "cloud/subnet/discovery"
	IDPEndpoint                          = APIPrefix + "idp"
	LogExporterEndpoint                  = APIPrefix + "logexporter"
	GroupingEndpoint                     = APIPrefix + "groupings"
	DiscoveredNetworksapi                = APIPrefix + "network/discovery"
	OnboardAppEndpoint                   = APIPrefix + "app/onboard"
	OnboardNetworkEndpoint               = APIPrefix + "network/onboard"
	OnboardAppSearchEndpoint             = APIPrefix + "app/onboard/search"
	OnboardNetworkSearchEndpoint         = APIPrefix + "network/onboard/search"
	OnboardAppSettingsEndpoint           = APIPrefix + "app/onboard/settings"
	OnboardAppSettingsUpdateEndpoint     = APIPrefix + "app/onboard/%s/settings/%s"
	OnboardAppSettingsGetEndpoint        = APIPrefix + "app/onboard/%s/settings"
	OnboardAppSummaryGetEndpoint         = APIPrefix + "app/onboard/%s/summary"
	OnboardAppCloudConfigEndpoint        = APIPrefix + "app/onboard/%s/cloud"
	OnboardNetworkCloudEndpoint          = APIPrefix + "network/onboard/%s/cloud"
	OnboardAppEndpointDiscoveredEndpoint = APIPrefix + "app/onboard/%s/endpoint/discovered"
	OnboardAppCloudEndpoint              = APIPrefix + "app/onboard/%s/cloud"
	OnboardAppDNSServiceEndpoint         = APIPrefix + "app/onboard/%s/dns"
	OnboardAppCertEndpoint               = APIPrefix + "app/onboard/%s/cert"
	OnboardAppOptOptionEndpoint          = APIPrefix + "app/onboard/%s/appopt"
	OnboardAppPolicyEndpoint             = APIPrefix + "app/onboard/%s/policy"
	OnboardNetworkPolicyEndpoint         = APIPrefix + "network/onboard/%s/policy"
	OnboardAppWafEndpoint                = APIPrefix + "app/onboard/%s/waf"
	OnboardAppDeploymentEndpoint         = APIPrefix + "app/deployment/%s"
	OnboardNetworkDeploymentEndpoint     = APIPrefix + "network/deployment/%s"
	NetworkOnboardEndpoint               = APIPrefix + "network/onboard/%s"
	OnboardAppClientAppEndpoint          = APIPrefix + "capp"
	AppDomainEndpoint                    = APIPrefix + "app/domain"
	WafAppDomainEndpoint                 = APIPrefix + "waf/%s/domain"
	WafRuleSetEndpoint                   = APIPrefix + "waf/ruleset"
	GenerateCertEndpoint                 = APIPrefix + "cert/domain"
	GetCertEndpoint                      = APIPrefix + "cert"
	// InternalCertEndpoint = APIPrefix + "cert/internal"
	// GenerateCertEndpoint     = APIPrefix + "cert/domain"
	UploadCACertEndpoint            = APIPrefix + "cert/ca"
	UploadClientCertEndpoint        = APIPrefix + "cert/client"
	CacheRuleEndpoint               = APIPrefix + "cacherule"
	WafEndpoint                     = APIPrefix + "waf"
	GetPolicyEndpoint               = APIPrefix + "policy/search"
	PolicyEndpoint                  = APIPrefix + "policy"
	IPRepEndpoint                   = APIPrefix + "ip-reputation"
	DynamicRiskEndpoint             = APIPrefix + "dynamicrisk"
	GeoVelocityEndpoint             = APIPrefix + "geo/allowlist"
	GetCityCode                     = APIPrefix + "geo/country"
	UserAllowlistEndpoint           = APIPrefix + "usersettings/excludeduser"
	AppListEndpoint                 = APIPrefix + "app/onboard"
	NetworkListEndpoint             = APIPrefix + "network/onboard"
	EDRConfigEndpoint               = APIPrefix + "edr"
	EDRProfileEndpoint              = APIPrefix + "edr-profile"
	DPProfileEndpoint               = APIPrefix + "device-posture"
	GeoLocEndpoint                  = APIPrefix + "geo/search"
	ConnectorPlacementEndpoint      = APIPrefix + "app/onboard/advanced/connector/placement"
	NameSpaceEndpoint               = APIPrefix + "namespace"
	GetNameSpaceEndpoint            = APIPrefix + "namespace/search"
	AssignNetworkEndpoint           = APIPrefix + "namespace/assign"
	ExportNetworkEndpoint           = APIPrefix + "namespace/export"
	WithdrawNetworkEndpoint         = APIPrefix + "namespace/withdraw"
	CloudGatewayEndpoint            = APIPrefix + "cloud-gateway"
	ConnectivityOptionsEndpont      = APIPrefix + "prosimo/app/connectivity/options"
	GetTransitSetupEndpoint         = APIPrefix + "network/transit/setup/search"
	GetNetworkEndpoint              = APIPrefix + "network/transit/setup/cloud-network/search"
	CreateTransitSetupEndpoint      = APIPrefix + "network/transit/setup"
	TransitSetupSummaryEndpoint     = APIPrefix + "network/transit/setup/summary"
	DeleteTransitSetupEndpoint      = APIPrefix + "network/transit/setup/bulk-delete"
	CreateTransitDeploymentEndpoint = APIPrefix + "network/transit/deployment"
	CreateRouteEntryEndpoint        = APIPrefix + "route-entry"
	GetRouteEntryEndpoint           = APIPrefix + "route-entry/search"
)
