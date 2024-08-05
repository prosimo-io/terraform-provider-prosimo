package client

import (
	"encoding/json"
	"log"
)

const (
	policyJson = `{
		"Users": {
			"property": [{
				"user_property": "User",
				"server_property": "email",    
				"operations": [{
					"user_operation_name": "Is",
					"server_operation_name": "Is"
				}, {
					"user_operation_name": "Contains",
					"server_operation_name": "Contains"
				}, {
					"user_operation_name": "Is NOT",
					"server_operation_name": "Is-not"
				}, {
					"user_operation_name": "Does NOT contain",
					"server_operation_name": "Does-not-contain"
				}, {
					"user_operation_name": "Starts with",
					"server_operation_name": "Starts-with"
				},{
					"user_operation_name": "Ends with",
					"server_operation_name": "Ends-with"
				}]
			}]
		},
		"Location": {
			"property": [{
				"user_property": "IP Prefix/Address",
				"server_property": "ipPrefix",
				"operations": [{
					"user_operation_name": "In",
					"server_operation_name": "In"
				}, {
					"user_operation_name": "NOT in",
					"server_operation_name": "Not-in"
				}]
			},{
				"user_property": "Country",
				"server_property": "geoip_country_code",
				"operations": [{
					"user_operation_name": "Is",
					"server_operation_name": "Is"
				}, {
					"user_operation_name": "Is NOT",
					"server_operation_name": "Is-not"
				}]
			},{
				"user_property": "State",
				"server_property": "geoip_country_code",
				"operations": [{
					"user_operation_name": "Is",
					"server_operation_name": "Is"
				}, {
					"user_operation_name": "Is NOT",
					"server_operation_name": "Is-not"
				}]
			},{
				"user_property": "City",
				"server_property": "geoip_country_code",
				"operations": [{
					"user_operation_name": "Is",
					"server_operation_name": "Is"
				}, {
					"user_operation_name": "Is NOT",
					"server_operation_name": "Is-not"
				}]
			}]
		},
		"Networks": {
			"property": [{
				"user_property": "Network",
				"server_property": "network",
				"operations": [{
					"user_operation_name": "Is",
					"server_operation_name": "Is"
				}]
			}]
		},
		"NetworkACL": {
			"property": [{
				"server_property": "networkACL",
				"operations": [{
					"server_operation_name": "In"
				}]
			}]
		},
		"EgressFqdn": {
			"property": [{
				"user_property": "EgressFqdn",
				"server_property": "fqdn",
				"operations": [{
					"user_operation_name": "Is",
					"server_operation_name": "is"
				}]
			}]
		},
		"IDP": {
			"property": [{
				"user_property": "",
				"server_property": "",
				"operations": [{
					"user_operation_name": "Is",
					"server_operation_name": "Is"
				}, {
					"user_operation_name": "Is NOT",
					"server_operation_name": "Is-not"
				}, {
					"user_operation_name": "Does NOT contain",
					"server_operation_name": "Does-not-contain"
				}, {
					"user_operation_name": "Starts with",
					"server_operation_name": "Starts-with"
				}, {
					"user_operation_name": "Ends with",
					"server_operation_name": "Ends-with"
				}, {
					"user_operation_name": "Contains",
					"server_operation_name": "Contains"
				}]
			}]

		},
		"Devices": {
			"property": [{
				"user_property": "Device OS",
				"server_property": "ua_os_name",
				"operations": [{
					"user_operation_name": "Is",
					"server_operation_name": "Is"
				}, {
					"user_operation_name": "Is NOT",
					"server_operation_name": "Is-not"
				}]
			}, {
				"user_property": "Device OS Version",
				"server_property": "ua_os_version",
				"operations": [{
					"user_operation_name": "Is",
					"server_operation_name": "Is"
				}, {
					"user_operation_name": "Is NOT",
					"server_operation_name": "Is-not"
				}, {
					"user_operation_name": "Is at least",
					"server_operation_name": "Is-atleast"
				}]
			}, {
				"user_property": "Device Category",
				"server_property": "ua_device_category",
				"operations": [{
					"user_operation_name": "Is",
					"server_operation_name": "Is"
				}, {
					"user_operation_name": "Is NOT",
					"server_operation_name": "Is-not"
				}]
			}, {
				"user_property": "Browser",
				"server_property": "ua_browser_name",
				"operations": [{
					"user_operation_name": "Is",
					"server_operation_name": "Is"
				}, {
					"user_operation_name": "Is NOT",
					"server_operation_name": "Is-not"
				}]
			}, {
				"user_property": "Browser Version",
				"server_property": "ua_browser_version",
				"operations": [{
					"user_operation_name": "Is",
					"server_operation_name": "Is"
				}, {
					"user_operation_name": "Is NOT",
					"server_operation_name": "Is-not"
				}, {
					"user_operation_name": "Is at least",
					"server_operation_name": "Is-atleast"
				}]
			}, {
				"user_property": "Trusted Device Certificate",
				"server_property": "client_cert_issuer_dn",
				"operations": [{
					"user_operation_name": "Is",
					"server_operation_name": "Is"
				}, {
					"user_operation_name": "Is NOT",
					"server_operation_name": "Is-not"
				}]
			}]

		},
		"Time": {
			"property": [{
				"user_property": "Time",
				"server_property": "time",
				"operations": [{
					"user_operation_name": "Between",
					"server_operation_name": "Between"
				}]
			}]

		},
		"URL": {
			"property": [{
				"user_property": "URL",
				"server_property": "url",
				"operations": [{
					"user_operation_name": "Does NOT contain",
					"server_operation_name": "Does-not-contain"
				}, {
					"user_operation_name": "Contains",
					"server_operation_name": "Contains"
				}]
			}]

		},
		"FQDN": {
			"property": [{
				"user_property": "FQDN",
				"server_property": "fqdn",
				"operations": [{
					"user_operation_name": "Is",
					"server_operation_name": "Is"
				}]
			}]

		},
		"Device_Posture_Profile": {
			"property": [{
				"user_property": "Risk Level",
				"server_property": "device_risk_level",
				"operations": [{
					"user_operation_name": "Is",
					"server_operation_name": "Is"
				}, {
					"user_operation_name": "Is NOT",
					"server_operation_name": "Is-not"
				}]
			}]
	
		},
		"Advanced": {
			"property": [{
				"user_property": "HTTP Method",
				"server_property": "method",
				"operations": [{
					"user_operation_name": "Is",
					"server_operation_name": "Is"
				}, {
					"user_operation_name": "Is NOT",
					"server_operation_name": "Is-not"
				}]
			}]

		}
	}
`

	internetEgressJson = `{
		"FQDN": {
			"property": [{
				"user_property": "FQDN",
				"server_property": "fqdn",
				"operations": [{
					"user_operation_name": "Is",
					"server_operation_name": "is"
				}, {
					"user_operation_name": "Is NOT",
					"server_operation_name": "is-not"
				}]
			},
			{
				"user_property": "Domain",
				"server_property": "domain",
				"operations": [{
					"user_operation_name": "Is",
					"server_operation_name": "is"
				}]
			}]
		}
	}
`
)

func GetPolicyServerTemplate() MatchItem {

	var input MatchItem
	err := json.Unmarshal([]byte(policyJson), &input)
	if err != nil {
		log.Println("Error unmarshaling policyJson:", err)
	}
	return input
}

func GetInternetEgressPolicyServerTemplate() IEMatchItem {

	var input IEMatchItem
	err := json.Unmarshal([]byte(internetEgressJson), &input)
	if err != nil {
		log.Println("Error unmarshaling internetEgressJson:", err)
	}
	return input
}
