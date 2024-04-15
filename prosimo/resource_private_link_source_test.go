package prosimo

import (
	"fmt"
	"testing"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourcePrivateLinkSource_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPrivateLinkSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePrivateLinkSource_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateLinkSourceExists("prosimo_private_link_source.test"),
				),
				// TODO there is a bug around maps that causes a permadiff for empty maps
				// ExpectNonEmptyPlan: true,
			},
			// {
			// 	ResourceName: "prosimo_app_onboarding_citrixvdi.test",
			// 	ImportState:  true,
			// 	// ImportStateVerify: true,
			// },
		},
	})
}

func TestAccResourcePrivateLinkSource_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPrivateLinkSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePrivateLinkSourcePre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateLinkSourceExists("prosimo_private_link_source.test_update"),
				),
			},
			{
				Config: testAccResourcePrivateLinkSourcePost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateLinkSourceExists("prosimo_private_link_source.test_update"),
				),
			},
			// {
			// 	ResourceName: "prosimo_app_onboarding_citrixvdi.test_update",
			// 	ImportState:  true,
			// 	// ImportStateVerify: true,
			// },
		},
	})
}

func testAccResourcePrivateLinkSource_basic(name string) string {
	return fmt.Sprintf(` 
	resource "prosimo_private_link_source" "test" {
		name = "terraform-test"
		cloud_creds_name = "prosimo-aws-app-iam"
		cloud_region = "us-east-2"
		cloud_sources {
			cloud_network {
				name = "cloud-census-us-east-2-vpc"
			}
			subnets {
				cidr = "192.13.0.0/24"
			}
		}
	}
`)
}

func testAccResourcePrivateLinkSourcePre() string {
	return fmt.Sprintf(`
	resource "prosimo_private_link_source" "test_update" {
		name = "terraform-test"
		cloud_creds_name = "prosimo-aws-app-iam"
		cloud_region = "us-east-2"
		cloud_sources {
			cloud_network {
				name = "cloud-census-us-east-2-vpc"
			}
			subnets {
				cidr = "192.13.0.1/24"
			}
		}
	}
`)
}

func testAccResourcePrivateLinkSourcePost() string {
	return fmt.Sprintf(`
	resource "prosimo_private_link_source" "test_update" {
		name = "terraform-test"
		cloud_creds_name = "prosimo-aws-app-iam"
		cloud_region = "us-east-2"
		cloud_sources {
			cloud_network {
				name = "cloud-census-us-east-2-vpc"
			}
			subnets {
				cidr = "192.13.0.0/24"
			}
		}
	}
`)
}

func testAccCheckPrivateLinkSourceExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		err, flag := client.GetPrivateLinkSourcefiltered(rs.Primary.ID)
		if flag {
			return fmt.Errorf("error fetching item with resource %s. %p", resource, err)
		}
		return nil
	}
}

func testAccCheckPrivateLinkSourceDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "prosimo_private_link_source" {
			continue
		}
		_, flag := client.GetPrivateLinkSourcefiltered(rs.Primary.ID)
		if !flag {
			return fmt.Errorf("Resource still avaible in portal %s", rs.Primary.ID)
		}
	}

	return nil
}
