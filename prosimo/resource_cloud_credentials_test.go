package prosimo

import (
	"fmt"
	"testing"

	"git.prosimo.io/prosimoio/tools/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceCloudCred_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudCredDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCloudCred_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudCredExists("prosimo_cloud_creds.test"),
				),
				// TODO there is a bug around maps that causes a permadiff for empty maps
				// ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:      "prosimo_cloud_creds.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceCloudCred_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudCredDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCloudCred_pre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudCredExists("prosimo_cloud_creds.test_update"),
				),
			},
			{
				Config: testAccResourceCloudCred_post(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudCredExists("prosimo_cloud_creds.test_update"),
				),
			},
			{
				ResourceName:      "prosimo_cloud_creds.test_update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceCloudCred_basic(name string) string {
	return fmt.Sprintf(` 

resource "prosimo_cloud_creds" "test" {
	cloud_type = "GCP"
	nickname   = "demo-1"
	gcp {
		file_path = "/Users/sibaprasadtripathy/Downloads/prosimo-test-bf8bbf15b37c_copy.json"
	}
	} 
`)
}

func testAccResourceCloudCred_pre() string {
	return fmt.Sprintf(` 

resource "prosimo_cloud_creds" "test_update" {
	cloud_type = "GCP"
	nickname   = "demo-2"
	gcp {
		file_path = "/Users/sibaprasadtripathy/Downloads/prosimo-test-bf8bbf15b37c_copy.json"
	}
	} 
`)
}

func testAccResourceCloudCred_post() string {
	return fmt.Sprintf(` 

resource "prosimo_cloud_creds" "test_update" {
	cloud_type = "GCP"
	nickname   = "demo-1"
	gcp {
		file_path = "/Users/sibaprasadtripathy/Downloads/prosimo-test-bf8bbf15b37c_copy.json"
	}
	} 
`)
}

func testAccCheckCloudCredExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		err, flag := client.GetCloudCredfiltered(rs.Primary.ID)
		if flag {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckCloudCredDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "prosimo_cloud_creds" {
			continue
		}
		_, flag := client.GetCloudCredfiltered(rs.Primary.ID)
		if !flag {
			return fmt.Errorf("Resource still avaible in portal %s", rs.Primary.ID)
		}
	}

	return nil
}
