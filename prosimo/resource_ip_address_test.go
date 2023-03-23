package prosimo

import (
	"fmt"
	"testing"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceIPAddress_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIPAddress_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPAddressExists("prosimo_ip_addresses.test"),
				),
				// TODO there is a bug around maps that causes a permadiff for empty maps
				// ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:      "prosimo_ip_addresses.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceIPAddress_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIPAddressPre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPAddressExists("prosimo_ip_addresses.test_update"),
				),
			},
			{
				Config: testAccResourceIPAddressPost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPAddressExists("prosimo_ip_addresses.test_update"),
				),
			},
			{
				ResourceName:      "prosimo_ip_addresses.test_update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceIPAddress_basic(name string) string {
	return fmt.Sprintf(` 
resource "prosimo_ip_addresses" "test" {
	cidr = "172.16.0.0/16"
	cloud_type = "AWS"
}  
`)
}

func testAccResourceIPAddressPre() string {
	return fmt.Sprintf(`
resource "prosimo_ip_addresses" "test_update" {
	cidr = "172.16.0.0/22"
	cloud_type = "AWS"
}  
`)
}

func testAccResourceIPAddressPost() string {
	return fmt.Sprintf(`
resource "prosimo_ip_addresses" "test_update" {
	cidr = "172.16.0.0/16"
	cloud_type = "AWS"
}  
`)
}

func testAccCheckIPAddressExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		name := rs.Primary.ID
		err, flag := client.GetIPPoolfiltered(name)
		if flag {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckIPAddressDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "prosimo_ip_addresses" {
			continue
		}
		_, flag := client.GetIPPoolfiltered(rs.Primary.ID)
		if !flag {
			return fmt.Errorf("Resource still avaible in portal %s", rs.Primary.ID)
		}
	}

	return nil
}
