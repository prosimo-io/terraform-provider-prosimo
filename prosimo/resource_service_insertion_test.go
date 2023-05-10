package prosimo

import (
	"fmt"
	"testing"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceServiceInsertion_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSharedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceServiceInsertion_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceInsertionExists("prosimo_service_insertion.test"),
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

func TestAccResourceServiceInsertion_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceInsertionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceServiceInsertionPre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceInsertionExists("prosimo_service_insertion.test_update"),
				),
			},
			{
				Config: testAccResourceServiceInsertionPost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceInsertionExists("prosimo_service_insertion.test_update"),
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

func testAccResourceServiceInsertion_basic(name string) string {
	return fmt.Sprintf(` 
	resource "prosimo_service_insertion" "test" {
		name = "Terraform-test"
		service_name = "terraform-test"
		source {
			networks{
				name = "terraform-si"
			}
		}
		target {
			networks{
				name = "0.0.0.0/0"
			}
		}
		ip_rules {
			source_addresses = ["any"]
			source_ports = ["any"]
			destination_addresses = ["any"]
			destination_ports = ["any"]
			protocols = ["ANY"]
		}
	}
`)
}

func testAccResourceServiceInsertionPre() string {
	return fmt.Sprintf(`
	resource "prosimo_service_insertion" "test_update" {
		name = "Terraform-test"
		service_name = "terraform-test"
		source {
			networks{
				name = "terraform-si"
			}
		}
		target {
			networks{
				name = "0.0.0.0/0"
			}
		}
		ip_rules {
			source_addresses = ["any"]
			source_ports = ["any"]
			destination_addresses = ["any"]
			destination_ports = ["any"]
			protocols = ["ANY", "TCP"]
		}
	}
`)
}

func testAccResourceServiceInsertionPost() string {
	return fmt.Sprintf(`
	resource "prosimo_service_insertion" "test_update" {
		name = "Terraform-test"
		service_name = "terraform-test"
		source {
			networks{
				name = "terraform-si"
			}
		}
		target {
			networks{
				name = "0.0.0.0/0"
			}
		}
		ip_rules {
			source_addresses = ["any"]
			source_ports = ["any"]
			destination_addresses = ["any"]
			destination_ports = ["any"]
			protocols = ["ANY"]
		}
	}
`)
}

func testAccCheckServiceInsertionExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		err, flag := client.GetSharedServicefiltered(rs.Primary.ID)
		if flag {
			return fmt.Errorf("error fetching item with resource %s. %p", resource, err)
		}
		return nil
	}
}

func testAccCheckServiceInsertionDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "prosimo_service_insertion" {
			continue
		}
		_, flag := client.GetServiceInsertionfiltered(rs.Primary.ID)
		if !flag {
			return fmt.Errorf("Resource still avaible in portal %s", rs.Primary.ID)
		}
	}

	return nil
}
