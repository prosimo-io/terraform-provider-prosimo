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

// var appOnboardSettings *client.AppOnboardSettings

// var appURL *client.AppURL
// var cloudConfigList []*client.AppURL

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
		if getTaskStatus.Status == "IN-PROGRESS" {
			return resource.RetryableError(fmt.Errorf("task %s is not completed yet", taskID))
		} else if getTaskStatus.Status == "FAILURE" {
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
		if getTaskStatus.Status == "IN-PROGRESS" {
			return resource.RetryableError(fmt.Errorf("task %s is not completed yet", taskID))
		} else if getTaskStatus.Status == "FAILURE" {
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

func createAppOnboardConfigs(ctx context.Context, d *schema.ResourceData, meta interface{}, appOnboardSettingsOpts *client.AppOnboardSettingsOpts) (diag.Diagnostics, *client.AppOnboardSettings) {

	var diags diag.Diagnostics
	var cloudConfigList []*client.AppURL
	appOnboardSettings := appOnboardSettingsOpts.GetAppOnboardSettings()
	log.Println("appOnboardSettingsOpts", appOnboardSettingsOpts)

	prosimoClient := meta.(*client.ProsimoClient)

	for _, appURLOpts := range appOnboardSettingsOpts.AppURLsOpts {
		appURL := appURLOpts.GetAppURL()
		//------------------------------------------------------------------------------------------------------
		//		Cloud Configuration
		//------------------------------------------------------------------------------------------------------

		cloudConfigOpts := appURLOpts.CloudConfigOpts

		if cloudConfigOpts.AppHOstedType == client.HostedPrivate {
			cloudCreds, err := prosimoClient.GetCloudCredsPrivate(ctx)
			if err != nil {
				return diag.FromErr(err), nil
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
				return diag.FromErr(err), nil
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
							return diag.FromErr(err), nil
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

					appOnboardCloudConfigResponseRegionsList, err := prosimoClient.DiscoverAppOnboardEndpoint(ctx, appURL, appOnboardSettingsOpts.AppOnboardType)
					if err != nil {
						return diag.FromErr(err), nil
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
								return diags, nil
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
								return diags, nil
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
							return diag.FromErr(err), nil
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
								return diag.FromErr(err), nil
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

						appOnboardCloudConfigResponseRegionsList, err := prosimoClient.DiscoverAppOnboardEndpoint(ctx, appURL, appOnboardSettingsOpts.AppOnboardType)
						if err != nil {
							return diag.FromErr(err), nil
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
									return diags, nil
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
									return diags, nil
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

		//------------------------------------------------------------------------------------------------------
		//		DNS Service Configuration
		//------------------------------------------------------------------------------------------------------
		if appOnboardSettingsOpts.AppOnboardType == client.TypeWEB || appOnboardSettingsOpts.AppOnboardType == client.TypeJumpBox || appOnboardSettingsOpts.AppOnboardType == client.TypeCloudSvc || appOnboardSettingsOpts.AppOnboardType == client.TypeCitrixVDI {
			log.Println("inside dns block")
			dnsServiceOpts := appURLOpts.DNSServiceOpts
			dnsService := &client.DNSService{}
			if dnsServiceOpts.Type == client.ManualDNSServiceType || dnsServiceOpts.Type == client.ProsimoDNSServiceType {
				dnsService.Type = dnsServiceOpts.Type

			} else {
				cloudCreds, err := prosimoClient.GetCloudCredsByName(ctx, dnsServiceOpts.CloudCredsName)
				if err != nil {
					return diag.FromErr(err), nil
				}

				dnsService.ID = cloudCreds.ID
				dnsService.Type = dnsServiceOpts.Type
			}
			appURL.DNSService = dnsService
		}
		//------------------------------------------------------------------------------------------------------
		//		Certificate Configuration
		//------------------------------------------------------------------------------------------------------
		if appOnboardSettingsOpts.AppOnboardType == client.TypeWEB || appOnboardSettingsOpts.AppOnboardType == client.TypeJumpBox || appOnboardSettingsOpts.AppOnboardType == client.TypeCloudSvc || appOnboardSettingsOpts.AppOnboardType == client.TypeCitrixVDI {
			appOnboardSettingsSSLCertData := &client.AppOnboardSettings{}
			if appOnboardSettingsOpts.ClientCert != "" {
				res, err := prosimoClient.GetCertDetails(ctx)
				if err != nil {
					return diag.FromErr(err), nil
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
			sslCertOpts := appURLOpts.SSLCertOpts
			if appURLOpts.DomainType == client.ProsimoAppDomain {
				res, err := prosimoClient.GetCertDetails(ctx)
				if err != nil {
					return diag.FromErr(err), nil
				}
				if len(res) > 0 {
					for _, cert := range res {
						if cert.ISTeamCert {
							appURL.CertID = cert.ID
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
						return diag.FromErr(err), nil
					}
					appURL.CertID = generateCertResponseData.ResourceData.ID

				} else if sslCertOpts.UploadCert != nil {
					uploadCertResponseData, err := prosimoClient.UploadCert(ctx, sslCertOpts.UploadCert.CertPath, sslCertOpts.UploadCert.KeyPath)
					if err == nil {
						appURL.CertID = uploadCertResponseData.ResourceData.ID
						// return diag.FromErr(err)
					} else {

						log.Println("[ERROR]: Certificate upload failed")
						return diag.FromErr(err), nil
					}
				} else if sslCertOpts.ExistingCert != "" {
					res, err := prosimoClient.GetCertDetails(ctx)
					if err != nil {
						return diag.FromErr(err), nil
					}
					if len(res) > 0 {
						flag := false
						for _, cert := range res {
							if cert.URL == sslCertOpts.ExistingCert {
								flag = true
								appURL.CertID = cert.ID
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
		}
		//------------------------------------------------------------------------------------------------------
		//		Cache Rule
		//------------------------------------------------------------------------------------------------------

		cacheRuleName := appURLOpts.CacheRuleName

		if cacheRuleName == "" {
			cacheRuleName = "Default Cache"
		}

		cacheRuleDbObj, err := prosimoClient.GetCacheRuleByName(ctx, cacheRuleName)
		if err != nil {
			return diag.FromErr(err), nil
		}
		appURL.CacheRuleID = cacheRuleDbObj.ID

		//------------------------------------------------------------------------------------------------------
		//		WAF
		//------------------------------------------------------------------------------------------------------

		wafName := appURLOpts.WafPolicyName
		if wafName != "" {
			wafDbObj, err := prosimoClient.GetWafByName(ctx, wafName)
			if err != nil {
				return diag.FromErr(err), nil
			}
			appURL.WafHTTP = wafDbObj.ID
		}

		//------------------------------------------------------------------------------------------------------
		cloudConfigList = append(cloudConfigList, appURL) //Append global appURL struct to global appURL List.(look at the variables on top)
		//------------------------------------------------------------------------------------------------------

	}
	//------------------------------------------------------------------------------------------------------
	//		Policies
	//------------------------------------------------------------------------------------------------------
	appOnboardPolicy := &client.AppOnboardPolicy{}

	for _, policyName := range appOnboardSettingsOpts.PolicyName {

		policyDbObj, err := prosimoClient.GetPolicyByName(ctx, policyName)
		if err != nil {
			return diag.FromErr(err), nil
		}

		appOnboardSettings.PolicyIDs = append(appOnboardPolicy.PolicyIDs, policyDbObj.ID)
	}

	appOnboardSettings.AppURLs = cloudConfigList
	appOnboardSettings.Dns_Discovery = appOnboardSettingsOpts.Dns_Discovery
	appOnboardSettings.EnableMultiCloud = appOnboardSettingsOpts.EnableMultiCloud
	appOnboardSettings.OptOption = appOnboardSettingsOpts.OptOption
	return diags, appOnboardSettings
}

func onboardApp(ctx context.Context, d *schema.ResourceData, meta interface{}, appOnboardSettingsOpts *client.AppOnboardSettingsOpts, appOnboardSettings *client.AppOnboardSettings) diag.Diagnostics {
	var diags diag.Diagnostics

	prosimoClient := meta.(*client.ProsimoClient)

	if appOnboardSettingsOpts.OnboardApp {

		appOnboardResData, err := prosimoClient.OnboardAppDeploymentV2(ctx, appOnboardSettings, client.ParamValueDeploy)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(appOnboardResData.ResourceData.ID)

		if d.Get("wait_for_rollout").(bool) {
			log.Printf("[DEBUG] Waiting for task id %s to complete", appOnboardResData.ResourceData.TaskID)
			err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
				retryUntilTaskComplete_appOnboard(ctx, d, meta, appOnboardResData.ResourceData.TaskID, appOnboardSettingsOpts))
			// retryUntilTaskComplete(ctx, d, meta, appOnboardResData.ResourceData.ID))
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("[INFO] task %s is successful", appOnboardResData.ResourceData.TaskID)
		}
	} else {
		appOnboardResData, err := prosimoClient.OnboardAppDeploymentV2(ctx, appOnboardSettings, client.ParamValueSave)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(appOnboardResData.ResourceData.ID)

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
