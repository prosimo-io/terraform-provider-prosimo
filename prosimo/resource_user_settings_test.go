package prosimo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceUserSettings_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserSettings_basic(name),
				Check:  resource.ComposeTestCheckFunc(
				// testAccCheckUserSettingsExists("prosimo_user_settings.test"),
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

func TestAccResourceUserSettings_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserSettingsPre(),
				Check:  resource.ComposeTestCheckFunc(
				// testAccCheckUserSettingsExists("prosimo_user_settings.test_update"),
				),
			},
			{
				Config: testAccResourceUserSettingsPost(),
				Check:  resource.ComposeTestCheckFunc(
				// testAccCheckUserSettingsExists("prosimo_user_settings.test_update"),
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

func testAccResourceUserSettings_basic(name string) string {
	return fmt.Sprintf(` 
	resource "prosimo_user_settings" "test" {
		allow_list {
		email   = "def@def.gz"  
		reason = "def" 
	  }    
		allow_list {
		email   = "abd@def.gz"  
		reason = "def abc" 
	  }         
	}
`)
}

func testAccResourceUserSettingsPre() string {
	return fmt.Sprintf(`
	resource "prosimo_user_settings" "test_update" {
		allow_list {
		email   = "def@def.gz"  
		reason = "def" 
	  }           
	}
`)
}

func testAccResourceUserSettingsPost() string {
	return fmt.Sprintf(`
	resource "prosimo_user_settings" "test_update" {
		allow_list {
		email   = "def@def.gz"  
		reason = "def" 
	  }    
		allow_list {
		email   = "abd@def.gz"  
		reason = "def abc" 
	  }         
	}
`)
}
