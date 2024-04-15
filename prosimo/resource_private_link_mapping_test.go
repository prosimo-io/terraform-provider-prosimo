package prosimo

import (
	"fmt"
	"testing"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourcePrivateLinkMapping_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPrivateLinkMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePrivateLinkMapping_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateLinkMappingExists("prosimo_private_link_mapping.test"),
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

func TestAccResourcePrivateLinkMapping_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPrivateLinkMappingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePrivateLinkMappingPre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateLinkMappingExists("prosimo_private_link_mapping.test_update"),
				),
			},
			{
				Config: testAccResourcePrivateLinkMappingPost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPrivateLinkMappingExists("prosimo_private_link_mapping.test_update"),
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

func testAccResourcePrivateLinkMapping_basic(name string) string {
	return fmt.Sprintf(` 
	resource "prosimo_private_link_mapping" "test" {
		source = "terraform-test"
		target = "common-app"
		hosted_zones {
			vpc_name = "cloud-census-us-east-2-vpc"
			domain_name = "app-aws-us-west-2-1681064493897.myeventarena.com"
			private_hosted_zone = "myeventarena.com."
		}
			hosted_zones {
			vpc_name = "cloud-census-us-east-2-vpc"
			domain_name = "speedtest-server-us-west-2-1681064493897.myeventarena.com"
			private_hosted_zone = "myeventarena.com."
		}
	}
`)
}

func testAccResourcePrivateLinkMappingPre() string {
	return fmt.Sprintf(`
	resource "prosimo_private_link_mapping" "test_update" {
		source = "terraform-test"
		target = "common-app"
		hosted_zones {
			vpc_name = "cloud-census-us-east-2-vpc"
			domain_name = "app-aws-us-west-2-1681064493897.myeventarena.com"
			private_hosted_zone = "myeventarena.com."
		}
			hosted_zones {
			vpc_name = "cloud-census-us-east-2-vpc"
			domain_name = "speedtest-server-us-west-2-1681064493897.myeventarena.com"
			private_hosted_zone = ".com."
		}
	}
`)
}

func testAccResourcePrivateLinkMappingPost() string {
	return fmt.Sprintf(`
	resource "prosimo_private_link_mapping" "test_update" {
		source = "terraform-test"
		target = "common-app"
		hosted_zones {
			vpc_name = "cloud-census-us-east-2-vpc"
			domain_name = "app-aws-us-west-2-1681064493897.myeventarena.com"
			private_hosted_zone = "myeventarena.com."
		}
			hosted_zones {
			vpc_name = "cloud-census-us-east-2-vpc"
			domain_name = "speedtest-server-us-west-2-1681064493897.myeventarena.com"
			private_hosted_zone = "myeventarena.com."
		}
	}
`)
}

func testAccCheckPrivateLinkMappingExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		err, flag := client.GetPrivateLinkMappingfiltered(rs.Primary.ID)
		if flag {
			return fmt.Errorf("error fetching item with resource %s. %p", resource, err)
		}
		return nil
	}
}

func testAccCheckPrivateLinkMappingDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "prosimo_private_link_mapping" {
			continue
		}
		_, flag := client.GetPrivateLinkMappingfiltered(rs.Primary.ID)
		if !flag {
			return fmt.Errorf("Resource still avaible in portal %s", rs.Primary.ID)
		}
	}

	return nil
}
