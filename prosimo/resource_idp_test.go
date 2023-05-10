package prosimo

import (
	"fmt"
	"testing"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceIDP_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIDPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIDP_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafExists("prosimo_idp.test"),
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

func TestAccResourceIDP_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIDPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIDPPre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDPExists("prosimo_idp.test_update"),
				),
			},
			{
				Config: testAccResourceIDPPost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDPExists("prosimo_idp.test_update"),
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

func testAccResourceIDP_basic(name string) string {
	return fmt.Sprintf(` 
	resource "prosimo_idp" "test" {
		idp_account = "google"
		auth_type = "oidc"
		account_url = "https://accounts.google.com"
		api_cred_provided = "yes"
		oidc {
		  client_id = "227796663766-6db0damqp8r104269fs2ih8i3e3uijhf.apps.googleusercontent.com"
		  secret_id = "var.google-secret"
		  admin_email = "test0319@myeventarena.com"
		  customer_id = "C02tdt1kg"
		  domain = "myevetarena7.com"
		  # file_path = "/User/file/path"
		}
		select_type = "partner"
		partner {
		user_domain =  ["myevetarena.com"]
		apps = ["all"]
		}
	}
	  
`)
}

func testAccResourceIDPPre() string {
	return fmt.Sprintf(`
	resource "prosimo_idp" "test_update" {
		idp_account = "google"
		auth_type = "oidc"
		account_url = "https://accounts.google.com"
		api_cred_provided = "yes"
		oidc {
		  client_id = "227796663766-6db0damqp8r104269fs2ih8i3e3uijhf.apps.googleusercontent.com"
		  secret_id = "var.google-secret"
		  admin_email = "test0319@myeventarena.com"
		  customer_id = "C02tdt1kg"
		  domain = "myevetarena7.com"
		  # file_path = "/User/file/path"
		}
		select_type = "partner"
		partner {
		user_domain =  ["myevetarena.com"]
		apps = ["none"]
		}
	}
`)
}

func testAccResourceIDPPost() string {
	return fmt.Sprintf(`
	resource "prosimo_idp" "test_update" {
		idp_account = "google"
		auth_type = "oidc"
		account_url = "https://accounts.google.com"
		api_cred_provided = "yes"
		oidc {
		  client_id = "227796663766-6db0damqp8r104269fs2ih8i3e3uijhf.apps.googleusercontent.com"
		  secret_id = "var.google-secret"
		  admin_email = "test0319@myeventarena.com"
		  customer_id = "C02tdt1kg"
		  domain = "myevetarena7.com"
		  # file_path = "/User/file/path"
		}
		select_type = "partner"
		partner {
		user_domain =  ["myevetarena.com"]
		apps = ["all"]
		}
	}
`)
}

func testAccCheckIDPExists(resource string) resource.TestCheckFunc {
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

func testAccCheckIDPDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "prosimo_idp" {
			continue
		}
		_, flag := client.GetIDPfiltered(rs.Primary.ID)
		if !flag {
			return fmt.Errorf("Resource still avaible in portal %s", rs.Primary.ID)
		}
	}

	return nil
}
