package prosimo

import (
	"fmt"
	"testing"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceAppOnboardingCloudSVC_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOnboardingCloudSVCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAppOnboardingCloudSVC_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOnboardingCloudSVCExists("prosimo_app_onboarding_cloudsvc.test"),
				),
				// TODO there is a bug around maps that causes a permadiff for empty maps
				// ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceAppOnboardingCloudSVC_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppOnboardingCloudSVCDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAppOnboardingCloudSVCPre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOnboardingCloudSVCExists("prosimo_app_onboarding_cloudsvc.test_update"),
				),
			},
			{
				Config: testAccResourceAppOnboardingCloudSVCPost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppOnboardingCloudSVCExists("prosimo_app_onboarding_cloudsvc.test_update"),
				),
			},
		},
	})
}

func testAccResourceAppOnboardingCloudSVC_basic(name string) string {
	return fmt.Sprintf(` 
resource "prosimo_app_onboarding_cloudsvc" "test" {

	app_name = "common-app-s3"
	cloud_svc = "amazon-s3"
	// idp_name = "azure_ad"
	app_urls {
		internal_domain = "qingscnetworkers.info"
		app_fqdn = "qingscnetworkers.info"
		subdomain_included = false
		cloud_config {
			# connection_option = "private"
			cloud_creds_name = "prosimo-aws-app-iam"
			edge_regions {
				region_name = "us-east-2"
				region_type = "active"
				buckets = ["qing-panbootstrap"]
			}
		}
		dns_service {
			type = "manual"
			}
		ssl_cert {
			generate_cert = true
		}
	}
	optimization_option = "PerformanceEnhanced"
	policy_name = ["ALLOW-ALL-USERS"]
	onboard_app = false
	decommission_app = false
}
`)
}

func testAccResourceAppOnboardingCloudSVCPre() string {
	return fmt.Sprintf(`
resource "prosimo_app_onboarding_cloudsvc" "test_update" {

	app_name = "common-app-s3-new"
	cloud_svc = "amazon-s3"
	// idp_name = "azure_ad"
	app_urls {
		internal_domain = "qingscnetworkers.info"
		app_fqdn = "qingscnetworkers.info"
		subdomain_included = false
		cloud_config {
			# connection_option = "private"
			cloud_creds_name = "prosimo-aws-app-iam"
			edge_regions {
				region_name = "us-east-2"
				region_type = "active"
				buckets = ["qing-panbootstrap"]
			}
		}
		dns_service {
			type = "manual"
			}
		ssl_cert {
			generate_cert = true
		}
	}
	optimization_option = "PerformanceEnhanced"
	policy_name = ["ALLOW-ALL-USERS"]
	onboard_app = false
	decommission_app = false
}
`)
}

func testAccResourceAppOnboardingCloudSVCPost() string {
	return fmt.Sprintf(`
resource "prosimo_app_onboarding_cloudsvc" "test_update" {

	app_name = "common-app-s3-new"
	cloud_svc = "amazon-s3"
	// idp_name = "azure_ad"
	app_urls {
		internal_domain = "qingscnetworkers.info"
		app_fqdn = "qingscnetworkers.info"
		subdomain_included = false
		cloud_config {
			# connection_option = "private"
			cloud_creds_name = "prosimo-aws-app-iam"
			edge_regions {
				region_name = "us-east-2"
				region_type = "active"
				buckets = ["qing-panbootstrap"]
			}
		}
		dns_service {
			type = "manual"
			}
		ssl_cert {
			generate_cert = true
		}
	}
	optimization_option = "PerformanceEnhanced"
	policy_name = ["DENY-ALL-USERS"]
	onboard_app = false
	decommission_app = false
}
`)
}

func testAccCheckAppOnboardingCloudSVCExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		name := rs.Primary.ID
		err, flag := client.GetAppOnboardingfiltered(name)
		if flag {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckAppOnboardingCloudSVCDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "prosimo_app_onboarding_cloudsvc" {
			continue
		}
		_, flag := client.GetAppOnboardingfiltered(rs.Primary.ID)
		if !flag {
			return fmt.Errorf("Resource still avaible in portal %s", rs.Primary.ID)
		}
	}

	return nil
}
