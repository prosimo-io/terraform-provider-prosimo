package prosimo

import (
	"context"
	"fmt"
	"log"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validate_primaryIDP(ctx context.Context, idpName string, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	prosimoClient := meta.(*client.ProsimoClient)
	idpList, err := prosimoClient.GetIDP(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	idpExists := false
	for _, idpDB := range idpList.IDPs {
		if idpDB.IDPName == idpName {
			idpExists = true
		}
	}
	if !idpExists {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Primary IDP doesnt exist with valid status",
			Detail:   "Primary IDP doesnt exist with valid status",
		})

		return diags
	}
	return diags
}

func retryUntilTaskComplete_appOnboard(ctx context.Context, d *schema.ResourceData, meta interface{}, taskID string, appOnboardSettingsOpts *client.AppOnboardSettingsOpts) resource.RetryFunc {
	var diags diag.Diagnostics
	prosimoClient := meta.(*client.ProsimoClient)
	return func() *resource.RetryError {
		getTaskStatus, err := prosimoClient.GetTaskStatus(ctx, taskID)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if getTaskStatus.TaskDetails.Status == "IN-PROGRESS" {
			return resource.RetryableError(fmt.Errorf("task %s is not completed yet", taskID))
		} else if getTaskStatus.TaskDetails.Status == "FAILURE" {
			for _, subtask := range getTaskStatus.ItemList {
				if subtask.Status == "FAILURE" {
					log.Printf("[ERROR]: task %s has failed at step %s, rolling back", taskID, subtask.Name)
				}
			}
			resourceAppOnboardingRead(ctx, d, meta)
			log.Println("[DEBUG]: offboarding app")
			appOnboardSettingsOpts.ID = d.Id()
			diags = offboardApp(ctx, d, meta, appOnboardSettingsOpts)
			if diags != nil {
				return resource.NonRetryableError(err)
			}

		}
		return nil
	}
}

func retryUntilTaskComplete_appOffboard(ctx context.Context, d *schema.ResourceData, meta interface{}, taskID string, appOnboardSettingsOpts *client.AppOnboardSettingsOpts) resource.RetryFunc {
	var diags diag.Diagnostics
	prosimoClient := meta.(*client.ProsimoClient)
	return func() *resource.RetryError {
		getTaskStatus, err := prosimoClient.GetTaskStatus(ctx, taskID)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if getTaskStatus.TaskDetails.Status == "IN-PROGRESS" {
			return resource.RetryableError(fmt.Errorf("task %s is not completed yet", taskID))
		} else if getTaskStatus.TaskDetails.Status == "FAILURE" {
			for _, subtask := range getTaskStatus.ItemList {
				if subtask.Status == "FAILURE" {
					log.Printf("[ERROR]: task %s has failed at step %s, rolling back", taskID, subtask.Name)
				}
			}
			if d.Get("force_offboard").(bool) {
				resourceAppOnboardingRead(ctx, d, meta)
				log.Println("[DEBUG]: Force offboarding app")
				appOnboardSettingsOpts.ID = d.Id()
				diags = forceOffboardApp(ctx, d, meta, appOnboardSettingsOpts)
				if diags != nil {
					return resource.NonRetryableError(err)
				}
			}

		}
		return nil
	}
}

func createAppOnboardSettings(ctx context.Context, d *schema.ResourceData, meta interface{}, appOnboardSettingsOpts *client.AppOnboardSettingsOpts) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	appOnboardSettings := appOnboardSettingsOpts.GetAppOnboardSettings()
	createdAppOnboardSettings, err := prosimoClient.CreateAppOnboardSettings(ctx, appOnboardSettings)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Created AppOnboard Settings for appName - %s, appAccessType - (%s), id - (%s)",
		appOnboardSettings.App_Name, appOnboardSettings.App_Access_Type, createdAppOnboardSettings.ResourceData.ID)
	d.SetId(createdAppOnboardSettings.ResourceData.ID)

	return diags

}

func createAppOnboardCloudConfigs(ctx context.Context, d *schema.ResourceData, meta interface{}, appOnboardSettingsOpts *client.AppOnboardSettingsOpts) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	appOnboardSettingsID := appOnboardSettingsOpts.ID

	cloudConfigList := []*client.AppURL{}

	for _, appURLOpts := range appOnboardSettingsOpts.AppURLsOpts {

		cloudConfigOpts := appURLOpts.CloudConfigOpts

		appURL := appURLOpts.GetAppURL()
		if cloudConfigOpts.AppHOstedType == client.HostedPrivate {
			cloudCreds, err := prosimoClient.GetCloudCredsPrivate(ctx)
			if err != nil {
				return diag.FromErr(err)
			}
			for _, cloudCred := range cloudCreds.CloudCreds {
				if cloudCred.Nickname == cloudConfigOpts.CloudCredsName {
					appURL.PrivateDcID = cloudCred.ID
				}
			}
			appURL.DCAappIP = cloudConfigOpts.DCAappIP
		} else {
			// log.Println("cloudConfigOpts.CloudCredsName", cloudConfigOpts.CloudCredsName)
			cloudCreds, err := prosimoClient.GetCloudCredsByName(ctx, cloudConfigOpts.CloudCredsName)
			if err != nil {
				return diag.FromErr(err)
			}

			appURL.CloudKeyID = cloudCreds.ID
		}
		appURL.ConnectionOption = cloudConfigOpts.ConnectionOption

		globalRegionList := []*client.AppOnboardCloudConfigRegions{}
		// globalCloudConfigRegion := &client.AppOnboardCloudConfigRegions{}
		appOnboardCloudConfigRegionsList := []*client.AppOnboardCloudConfigRegions{}

		for _, appOnboardCloudConfigRegionOpts := range cloudConfigOpts.Regions {
			globalCloudConfigRegion := &client.AppOnboardCloudConfigRegions{}
			if appOnboardSettingsOpts.AppType == client.CloudPlatformType || appOnboardSettingsOpts.AppOnboardType == client.TypeCloudSvc {
				appOnboardCloudConfigRegions := &client.AppOnboardCloudConfigRegions{}
				appOnboardCloudConfigRegions.ID = appOnboardCloudConfigRegionOpts.ID
				appOnboardCloudConfigRegions.Name = appOnboardCloudConfigRegionOpts.Name
				appOnboardCloudConfigRegions.RegionType = appOnboardCloudConfigRegionOpts.RegionType
				appOnboardCloudConfigRegions.Buckets = appOnboardCloudConfigRegionOpts.Buckets
				globalCloudConfigRegion.Name = appOnboardCloudConfigRegions.Name
				globalCloudConfigRegion.RegionType = appOnboardCloudConfigRegions.RegionType
				globalCloudConfigRegion.Buckets = appOnboardCloudConfigRegions.Buckets
				appOnboardCloudConfigRegionsList = append(appOnboardCloudConfigRegionsList, appOnboardCloudConfigRegions)
				appURL.Regions = appOnboardCloudConfigRegionsList
			} else {
				if appOnboardCloudConfigRegionOpts.BackendIPAddressDiscover {
					appOnboardCloudConfigRegions := &client.AppOnboardCloudConfigRegions{}
					appOnboardCloudConfigRegions.Name = appOnboardCloudConfigRegionOpts.Name
					appOnboardCloudConfigRegions.ConnOption = appOnboardCloudConfigRegionOpts.ConnOption
					if appOnboardCloudConfigRegions.ConnOption == client.OptntransitGateway || appOnboardCloudConfigRegions.ConnOption == client.OptnazureTransitVnet || appOnboardCloudConfigRegions.ConnOption == client.OptnvwanHub {
						tgwendptInputList := []*client.AppOnboardCloudRegionEndpoints{}
						tgwendptInput := &client.AppOnboardCloudRegionEndpoints{}
						tgwendptInput.AppNetworkID = appOnboardCloudConfigRegionOpts.AppnetworkID
						tgwendptInput.AttachPointID = appOnboardCloudConfigRegionOpts.AttachPointID
						tgwendptInputList = append(tgwendptInputList, tgwendptInput)
						appOnboardCloudConfigRegions.Endpoints = tgwendptInputList
					}
					appOnboardCloudConfigRegions.ID = appOnboardCloudConfigRegionOpts.ID
					appOnboardCloudConfigRegions.RegionType = appOnboardCloudConfigRegionOpts.RegionType
					if cloudConfigOpts.AppHOstedType == client.HostedPrivate {
						CloudRegionDetails, err := prosimoClient.GetCloudRegion(ctx, appURL.CloudKeyID)
						if err != nil {
							return diag.FromErr(err)
						}
						for _, CloudRegion := range CloudRegionDetails.CloudRegionList {
							appOnboardCloudConfigRegions.LocationID = CloudRegion.LocationID
						}
					}
					globalCloudConfigRegion.LocationID = appOnboardCloudConfigRegions.LocationID
					globalCloudConfigRegion.Name = appOnboardCloudConfigRegionOpts.Name
					globalCloudConfigRegion.ConnOption = appOnboardCloudConfigRegionOpts.ConnOption
					globalCloudConfigRegion.RegionType = appOnboardCloudConfigRegionOpts.RegionType
					// appOnboardCloudConfigRegionsList := []*client.AppOnboardCloudConfigRegions{}
					appOnboardCloudConfigRegions.InputType = "discovered"
					if appOnboardCloudConfigRegionOpts.ConnOption == client.OptntransitGateway {
						globalCloudConfigRegion.ModifyTgwAppRouteTable = appOnboardCloudConfigRegionOpts.ModifyTgwAppRouteTable
					}
					appOnboardCloudConfigRegionsList = append(appOnboardCloudConfigRegionsList, appOnboardCloudConfigRegions)
					appURL.Regions = appOnboardCloudConfigRegionsList

					appOnboardCloudConfigResponseRegionsList, err := prosimoClient.DiscoverAppOnboardEndpoint(ctx, appOnboardSettingsID, appURL)
					if err != nil {
						return diag.FromErr(err)
					}

					if appOnboardCloudConfigRegionOpts.AttachPointID != "" {
						endpointList := []*client.AppOnboardCloudRegionEndpoints{}
						NetworkIDFlag := false
						AttachIDFlag := false
						for _, appOnboardCloudConfigResponseRegions := range appOnboardCloudConfigResponseRegionsList {
							for _, endpoints := range appOnboardCloudConfigResponseRegions.Endpoints {
								if endpoints.AppNetworkID == appOnboardCloudConfigRegionOpts.AppnetworkID {
									endpointList = append(endpointList, endpoints)
									NetworkIDFlag = true
								}
							}
							if !NetworkIDFlag {
								diags = append(diags, diag.Diagnostic{
									Severity: diag.Error,
									Summary:  "Invalid  AppNetworkID.",
									Detail:   "Input AppNetworkID does not exist.",
								})
								return diags
							}
							// }

							globalCloudConfigRegion.Endpoints = endpointList
							// log.Println("globalCloudConfigRegion.Endpoints", globalCloudConfigRegion.Endpoints)
							// }
							for _, attachpoint := range appOnboardCloudConfigResponseRegions.Attchpoints {
								if attachpoint.AttachPointID == appOnboardCloudConfigRegionOpts.AttachPointID {
									globalCloudConfigRegion.Endpoints[0].AttachPointID = attachpoint.AttachPointID
									AttachIDFlag = true
									// appOnboardCloudConfigResponseRegionsList[0].Endpoints[0].AttachPointID = attachpoint.AttachPointID
								}
							}
							if !AttachIDFlag {
								diags = append(diags, diag.Diagnostic{
									Severity: diag.Error,
									Summary:  "Invalid  AttachPointID.",
									Detail:   "Input AttachPointID does not exist.",
								})
								return diags
							}
							// }
						}
					} else {

						globalCloudConfigRegion.Endpoints = appOnboardCloudConfigResponseRegionsList[0].Endpoints
					}

				} else if appOnboardSettingsOpts.Dns_Discovery || appURL.SubdomainIncluded {
					appOnboardDnsCustomoptns := &client.AppOnboardDnsCustom{}
					// if appURL.SubdomainIncluded {
					appOnboardDnsCustomoptns.IsHealthCheckEnabled = appURLOpts.DnsCustom.IsHealthCheckEnabled
					if appURLOpts.DnsCustom.DnsAppName != "" {
						apppnboardSearchops := &client.AppOnboardSearch{}
						statusFilterops := "DEPLOYED"
						statusFilter := []string{}
						statusFilter = append(statusFilter, statusFilterops)
						apppnboardSearchops.StatusFilter = statusFilter
						apppnboardSearchops.NetworkService = "dns"
						returnedappList, err := prosimoClient.SearchAppOnboardApps(ctx, apppnboardSearchops)
						if err != nil {
							return diag.FromErr(err)
						}
						if len(returnedappList.Data.Records) > 0 {
							for _, onboardApp := range returnedappList.Data.Records {
								if onboardApp.App_Name == appURLOpts.DnsCustom.DnsAppName {
									appOnboardDnsCustomoptns.DnsAppID = onboardApp.ID
								} else {
									log.Println("[ERROR]: Input DNS app doesnot exist")
								}
							}
						} else {
							log.Println("[ERROR]: Input DNS app doesnot exist")
						}
					} else {
						appOnboardDnsCustomoptns.DnsServers = appURLOpts.DnsCustom.DnsServers
					}
					appURL.DnsCustom = appOnboardDnsCustomoptns
				} else {

					for _, backendIPAddress := range appOnboardCloudConfigRegionOpts.BackendIPAddressEntry {

						appOnboardCloudConfigRegions := &client.AppOnboardCloudConfigRegions{}
						appOnboardCloudConfigRegions.Name = appOnboardCloudConfigRegionOpts.Name
						appOnboardCloudConfigRegions.ConnOption = appOnboardCloudConfigRegionOpts.ConnOption
						if appOnboardCloudConfigRegions.ConnOption == client.OptntransitGateway || appOnboardCloudConfigRegions.ConnOption == client.OptnazureTransitVnet || appOnboardCloudConfigRegions.ConnOption == client.OptnvwanHub {
							tgwendptInputList := []*client.AppOnboardCloudRegionEndpoints{}
							tgwendptInput := &client.AppOnboardCloudRegionEndpoints{}
							tgwendptInput.AppNetworkID = appOnboardCloudConfigRegionOpts.AppnetworkID
							tgwendptInput.AttachPointID = appOnboardCloudConfigRegionOpts.AttachPointID
							tgwendptInputList = append(tgwendptInputList, tgwendptInput)
							appOnboardCloudConfigRegions.Endpoints = tgwendptInputList
						}
						appOnboardCloudConfigRegions.RegionType = appOnboardCloudConfigRegionOpts.RegionType
						appOnboardCloudConfigRegions.ID = appOnboardCloudConfigRegionOpts.ID
						if cloudConfigOpts.AppHOstedType == client.HostedPrivate {
							CloudRegionDetails, err := prosimoClient.GetCloudRegion(ctx, appURL.CloudKeyID)
							if err != nil {
								return diag.FromErr(err)
							}
							for _, CloudRegion := range CloudRegionDetails.CloudRegionList {
								appOnboardCloudConfigRegions.LocationID = CloudRegion.LocationID
							}
						}

						globalCloudConfigRegion.LocationID = appOnboardCloudConfigRegions.LocationID
						globalCloudConfigRegion.Name = appOnboardCloudConfigRegionOpts.Name
						globalCloudConfigRegion.RegionType = appOnboardCloudConfigRegionOpts.RegionType
						globalCloudConfigRegion.ConnOption = appOnboardCloudConfigRegionOpts.ConnOption
						globalCloudConfigRegion.ID = appOnboardCloudConfigRegionOpts.ID
						globalCloudConfigRegion.InputType = "entry"
						if appOnboardCloudConfigRegionOpts.ConnOption == client.OptntransitGateway {
							globalCloudConfigRegion.ModifyTgwAppRouteTable = appOnboardCloudConfigRegionOpts.ModifyTgwAppRouteTable
						}
						appURL.IP = backendIPAddress
						appOnboardCloudConfigRegions.InputType = "entry"
						// appOnboardCloudConfigRegionsList := []*client.AppOnboardCloudConfigRegions{}
						appOnboardCloudConfigRegionsList = append(appOnboardCloudConfigRegionsList, appOnboardCloudConfigRegions)
						appURL.Regions = appOnboardCloudConfigRegionsList

						appOnboardCloudConfigResponseRegionsList, err := prosimoClient.DiscoverAppOnboardEndpoint(ctx, appOnboardSettingsID, appURL)
						if err != nil {
							return diag.FromErr(err)
						}
						if appOnboardCloudConfigRegionOpts.AttachPointID != "" {
							endpointList := []*client.AppOnboardCloudRegionEndpoints{}
							NetworkIDFlag := false
							AttachIDFlag := false
							for _, appOnboardCloudConfigResponseRegions := range appOnboardCloudConfigResponseRegionsList {
								for _, endpoints := range appOnboardCloudConfigResponseRegions.Endpoints {
									if endpoints.AppNetworkID == appOnboardCloudConfigRegionOpts.AppnetworkID {
										endpointList = append(endpointList, endpoints)
										NetworkIDFlag = true
									}
								}
								if !NetworkIDFlag {
									diags = append(diags, diag.Diagnostic{
										Severity: diag.Error,
										Summary:  "Invalid  AppNetworkID.",
										Detail:   "Input AppNetworkID does not exist.",
									})
									return diags
								}
								// }

								globalCloudConfigRegion.Endpoints = endpointList
								// log.Println("globalCloudConfigRegion.Endpoints", globalCloudConfigRegion.Endpoints)
								// }
								for _, attachpoint := range appOnboardCloudConfigResponseRegions.Attchpoints {
									if attachpoint.AttachPointID == appOnboardCloudConfigRegionOpts.AttachPointID {
										globalCloudConfigRegion.Endpoints[0].AttachPointID = attachpoint.AttachPointID
										AttachIDFlag = true
									}
								}
								if !AttachIDFlag {
									diags = append(diags, diag.Diagnostic{
										Severity: diag.Error,
										Summary:  "Invalid  AttachPointID.",
										Detail:   "Input AttachPointID does not exist.",
									})
									return diags
								}
								// }
							}
						} else {

							globalCloudConfigRegion.Endpoints = appOnboardCloudConfigResponseRegionsList[0].Endpoints
						}
						// globalCloudConfigRegion.Endpoints = append(globalCloudConfigRegion.Endpoints, appOnboardCloudConfigResponseRegionsList[0].Endpoints...)
						globalCloudConfigRegion.InputEndpoints = []string{}
						globalCloudConfigRegion.SelectedEndpoints = []string{}

					}
				}
			}
			globalRegionList = append(globalRegionList, globalCloudConfigRegion)
		}
		appURL.IP = ""
		if appURL.DnsCustom == nil {
			appURL.Regions = globalRegionList
		} else {
			appURL.Regions = nil
		}
		cloudConfigList = append(cloudConfigList, appURL)

	}

	appOnboardSettingsCloudConfigData := &client.AppOnboardSettings{}
	appOnboardSettingsCloudConfigData.AppURLs = cloudConfigList
	appOnboardSettingsCloudConfigData.Dns_Discovery = appOnboardSettingsOpts.Dns_Discovery
	_, err := prosimoClient.CreateAppOnboardCloudConfig(ctx, appOnboardSettingsID, appOnboardSettingsCloudConfigData)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func createAppOnboardDNSService(ctx context.Context, d *schema.ResourceData, meta interface{}, appOnboardSettingsOpts *client.AppOnboardSettingsOpts) diag.Diagnostics {
	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	appOnboardSettingsID := appOnboardSettingsOpts.ID

	dnsServiceAppList := []*client.AppURL{}

	for _, appURLOpts := range appOnboardSettingsOpts.AppURLsOpts {

		dnsServiceOpts := appURLOpts.DNSServiceOpts

		appURL := appURLOpts.GetAppURL()

		dnsService := &client.DNSService{}

		if dnsServiceOpts.Type == client.ManualDNSServiceType || dnsServiceOpts.Type == client.ProsimoDNSServiceType {
			dnsService.Type = dnsServiceOpts.Type

		} else {
			cloudCreds, err := prosimoClient.GetCloudCredsByName(ctx, dnsServiceOpts.CloudCredsName)
			if err != nil {
				return diag.FromErr(err)
			}

			dnsService.ID = cloudCreds.ID
			dnsService.Type = dnsServiceOpts.Type
		}

		dnsAppURL := &client.AppURL{}
		dnsAppURL.ID = appURL.ID
		dnsAppURL.DNSService = dnsService
		dnsServiceAppList = append(dnsServiceAppList, dnsAppURL)

	}

	appOnboardSettingsDNSServiceData := &client.AppOnboardSettings{}
	appOnboardSettingsDNSServiceData.AppURLs = dnsServiceAppList
	_, err := prosimoClient.CreateDNSService(ctx, appOnboardSettingsID, appOnboardSettingsDNSServiceData)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func createAppOnboardSSLCert(ctx context.Context, d *schema.ResourceData, meta interface{}, appOnboardSettingsOpts *client.AppOnboardSettingsOpts) diag.Diagnostics {
	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	appOnboardSettingsID := appOnboardSettingsOpts.ID

	sslCertAppList := []*client.AppURL{}
	appOnboardSettingsSSLCertData := &client.AppOnboardSettings{}
	if appOnboardSettingsOpts.ClientCert != "" {
		res, err := prosimoClient.GetCertDetails(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(res) > 0 {
			flag := false
			for _, cert := range res {
				if cert.URL == appOnboardSettingsOpts.ClientCert && cert.Type == "Client" {
					flag = true
					appOnboardSettingsSSLCertData.CertID = cert.ID
				}
			}
			if !flag {
				log.Println("[ERROR]: Client Certficate does not exist.")
			}
		} else {
			log.Println("[ERROR]: Client Certficate does not exist.")
		}
	}
	for _, appURLOpts := range appOnboardSettingsOpts.AppURLsOpts {

		sslCertOpts := appURLOpts.SSLCertOpts

		appURL := appURLOpts.GetAppURL()
		sslCertAppURL := &client.AppURL{}
		sslCertAppURL.ID = appURL.ID
		if appURLOpts.DomainType == client.ProsimoAppDomain {
			res, err := prosimoClient.GetCertDetails(ctx)
			if err != nil {
				return diag.FromErr(err)
			}
			if len(res) > 0 {
				for _, cert := range res {
					if cert.ISTeamCert {
						sslCertAppURL.CertID = cert.ID
					}
				}
			}
		} else {
			if sslCertOpts.GenerateCert {
				generateCert := &client.GenerateCert{}
				generateCert.CA = "Let’sEncrypt"
				generateCert.URL = appURL.InternalDomain
				generateCert.TermsOfServiceAgreed = true

				generateCertResponseData, err := prosimoClient.GenerateCert(ctx, generateCert)
				if err != nil {
					return diag.FromErr(err)
				}
				sslCertAppURL.CertID = generateCertResponseData.ResourceData.ID

			} else if sslCertOpts.UploadCert != nil {
				uploadCertResponseData, err := prosimoClient.UploadCert(ctx, sslCertOpts.UploadCert.CertPath, sslCertOpts.UploadCert.KeyPath)
				if err == nil {
					sslCertAppURL.CertID = uploadCertResponseData.ResourceData.ID
					// return diag.FromErr(err)
				} else {

					log.Println("[ERROR]: Certificate upload failed")
					return diag.FromErr(err)
				}
			} else if sslCertOpts.ExistingCert != "" {
				res, err := prosimoClient.GetCertDetails(ctx)
				if err != nil {
					return diag.FromErr(err)
				}
				if len(res) > 0 {
					flag := false
					for _, cert := range res {
						if cert.URL == sslCertOpts.ExistingCert {
							flag = true
							sslCertAppURL.CertID = cert.ID
						}
					}
					if !flag {
						log.Println("[ERROR]: Certficate does not exist.")
					}
				} else {
					log.Println("[ERROR]: Certficate does not exist.")
				}
			}

		}
		sslCertAppList = append(sslCertAppList, sslCertAppURL)

	}
	appOnboardSettingsSSLCertData.AppURLs = sslCertAppList
	_, err := prosimoClient.CreateAppOnboardCert(ctx, appOnboardSettingsID, appOnboardSettingsSSLCertData)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func createAppOnboardOptOption(ctx context.Context, d *schema.ResourceData, meta interface{}, appOnboardSettingsOpts *client.AppOnboardSettingsOpts) diag.Diagnostics {
	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	appOnboardSettingsID := appOnboardSettingsOpts.ID

	optOptionAppList := []*client.AppURL{}

	for _, appURLOpts := range appOnboardSettingsOpts.AppURLsOpts {

		cacheRuleName := appURLOpts.CacheRuleName

		appURL := appURLOpts.GetAppURL()
		optOptionAppURL := &client.AppURL{}
		optOptionAppURL.ID = appURL.ID

		if cacheRuleName == "" {
			cacheRuleName = "Default Cache"
		}

		cacheRuleDbObj, err := prosimoClient.GetCacheRuleByName(ctx, cacheRuleName)
		if err != nil {
			return diag.FromErr(err)
		}

		optOptionAppURL.CacheRuleID = cacheRuleDbObj.ID

		optOptionAppList = append(optOptionAppList, optOptionAppURL)

	}

	appOnboardSettingsOptOptionData := &client.AppOnboardSettings{}
	appOnboardSettingsOptOptionData.AppURLs = optOptionAppList
	appOnboardSettingsOptOptionData.EnableMultiCloud = appOnboardSettingsOpts.EnableMultiCloud
	appOnboardSettingsOptOptionData.OptOption = appOnboardSettingsOpts.OptOption
	_, err := prosimoClient.CreateAppOnboardOptOption(ctx, appOnboardSettingsID, appOnboardSettingsOptOptionData)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func createAppOnboardSecurity(ctx context.Context, d *schema.ResourceData, meta interface{}, appOnboardSettingsOpts *client.AppOnboardSettingsOpts) diag.Diagnostics {
	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	appOnboardSettingsID := appOnboardSettingsOpts.ID
	appOnboardPolicy := &client.AppOnboardPolicy{}

	for _, policyName := range appOnboardSettingsOpts.PolicyName {

		policyDbObj, err := prosimoClient.GetPolicyByName(ctx, policyName)
		if err != nil {
			return diag.FromErr(err)
		}

		appOnboardPolicy.PolicyIDs = append(appOnboardPolicy.PolicyIDs, policyDbObj.ID)
	}

	_, err := prosimoClient.CreateAppOnboardPolicy(ctx, appOnboardSettingsID, appOnboardPolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	wafAppList := []*client.AppURL{}

	for _, appURLOpts := range appOnboardSettingsOpts.AppURLsOpts {

		wafName := appURLOpts.WafPolicyName

		appURL := appURLOpts.GetAppURL()
		wafAppURL := &client.AppURL{}
		wafAppURL.ID = appURL.ID

		if wafName != "" {
			wafDbObj, err := prosimoClient.GetWafByName(ctx, wafName)
			if err != nil {
				return diag.FromErr(err)
			}
			wafAppURL.WafHTTP = wafDbObj.ID
		}

		wafAppList = append(wafAppList, wafAppURL)

	}

	appOnboardSettingsWafData := &client.AppOnboardSettings{}
	appOnboardSettingsWafData.AppURLs = wafAppList
	_, err = prosimoClient.CreateAppOnboardWaf(ctx, appOnboardSettingsID, appOnboardSettingsWafData)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func onboardApp(ctx context.Context, d *schema.ResourceData, meta interface{}, appOnboardSettingsOpts *client.AppOnboardSettingsOpts) diag.Diagnostics {
	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	appOnboardSettingsID := appOnboardSettingsOpts.ID

	if appOnboardSettingsOpts.OnboardApp {
		appOnboardDeploymentData := &client.AppOnboardSettings{}
		appOnboardResData, err := prosimoClient.OnboardAppDeployment(ctx, appOnboardSettingsID, appOnboardDeploymentData)
		if err != nil {
			return diag.FromErr(err)
		}
		if d.Get("wait_for_rollout").(bool) {
			log.Printf("[DEBUG] Waiting for task id %s to complete", appOnboardResData.ResourceData.ID)
			err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
				retryUntilTaskComplete_appOnboard(ctx, d, meta, appOnboardResData.ResourceData.ID, appOnboardSettingsOpts))
			// retryUntilTaskComplete(ctx, d, meta, appOnboardResData.ResourceData.ID))
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[INFO] task %s is successful", appOnboardResData.ResourceData.ID)
		}
	}

	return diags
}

func offboardApp(ctx context.Context, d *schema.ResourceData, meta interface{}, appOnboardSettingsOpts *client.AppOnboardSettingsOpts) diag.Diagnostics {
	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	appOnboardSettingsID := appOnboardSettingsOpts.ID

	appOffboardResData, err := prosimoClient.OffboardAppDeployment(ctx, appOnboardSettingsID)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.Get("wait_for_rollout").(bool) {
		log.Printf("[DEBUG] Waiting for task id %s to complete", appOffboardResData.ResourceData.ID)
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
			retryUntilTaskComplete_appOffboard(ctx, d, meta, appOffboardResData.ResourceData.ID, appOnboardSettingsOpts))
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[DEBUG] task %s is successful", appOffboardResData.ResourceData.ID)
	}

	return diags
}

func forceOffboardApp(ctx context.Context, d *schema.ResourceData, meta interface{}, appOnboardSettingsOpts *client.AppOnboardSettingsOpts) diag.Diagnostics {
	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	appOnboardSettingsID := appOnboardSettingsOpts.ID

	appOffboardResData, err := prosimoClient.ForceOffboardAppDeployment(ctx, appOnboardSettingsID)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.Get("wait_for_rollout").(bool) {
		log.Printf("[DEBUG] Waiting for task id %s to complete", appOffboardResData.ResourceData.ID)
		err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
			retryUntilTaskComplete(ctx, d, meta, appOffboardResData.ResourceData.ID))
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[DEBUG] task %s is successful", appOffboardResData.ResourceData.ID)
	}

	return diags
}

func updateAppOnboardSettings(ctx context.Context, d *schema.ResourceData, meta interface{},
	appOnboardSettingsOpts *client.AppOnboardSettingsOpts, appOnboardSettingsID string) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	appOnboardSettings := appOnboardSettingsOpts.GetAppOnboardSettings()
	updatedAppOnboardSettings, err := prosimoClient.UpdateAppOnboardSettings(ctx, appOnboardSettings)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Updated AppOnboard Settings for appName - %s, appAccessType - (%s), id - (%s)",
		appOnboardSettings.App_Name, appOnboardSettings.App_Access_Type, updatedAppOnboardSettings.ResourceData.ID)
	d.SetId(updatedAppOnboardSettings.ResourceData.ID)

	return diags

}

func resourceAppOnboardingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)
	// log.Printf("resourceAppOnboardingRead %s", d.Id())
	appOnboardSettingsDbObj, err := prosimoClient.GetAppOnboardSettings(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(appOnboardSettingsDbObj.ID)
	d.Set("app_name", appOnboardSettingsDbObj.App_Name)
	d.Set("optimization_option", appOnboardSettingsDbObj.OptOption)
	d.Set("enable_multi_cloud_access", appOnboardSettingsDbObj.EnableMultiCloud)
	d.Set("onboard_app", appOnboardSettingsDbObj.Deployed)
	d.Set("decommission_app", d.Get("decommission_app").(bool))

	// for _, policyIDs := range appOnboardSettingsDbObj.PolicyIDs {
	// 	policyDbObj, err := prosimoClient.GetPolicyByID(ctx, policyIDs)
	// 	if err != nil {
	// 		return diag.FromErr(err)
	// 	}
	// 	d.Set("policy_name", policyDbObj.Name)
	// }

	return diags

}
