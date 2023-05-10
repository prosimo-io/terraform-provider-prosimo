package prosimo

import (
	"fmt"
	"testing"

	"git.prosimo.io/prosimoio/prosimo/terraform-provider-prosimo.git/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPolicy_basic(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists("prosimo_policy.test"),
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

func TestAccResourcePolicy_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyPre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists("prosimo_policy.test_update"),
				),
			},
			{
				Config: testAccResourcePolicyPost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists("prosimo_policy.test_update"),
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

func testAccResourcePolicy_basic(name string) string {
	return fmt.Sprintf(` 
	resource "prosimo_policy" "test" {
		name = "psonar-test-policy-transit"
		app_access_type = "transit"
		details {
		   actions = "allow"
		   # lock_users = false
		   # alert = true
		   # mfa = true
		   # lock_users = false
			 matches {
			   match_entries {
				   property = "URL"
				   operation = "Does NOT contain"
				   type     = "url"
				   values {
					   inputitems {
						   name = "1235678"                      
					   }
				   }
			   }
			   match_entries {
				   property = "Network"
				   operation = "Is"
				   type     = "networks"
				   values {
					   selecteditems {
						   name = "aws-uswest1-spoke1-infra-tf"                  
					   }
				   }
			   }
			   match_entries {
				   type     = "networkacl"
				   values {
					   inputitems {
						  ip_details {
						   source_ip = ["any"]
						   target_ip = ["any"]
						   protocol = ["tcp"]
						   source_port = ["any"]
						   target_port = ["any"]
						  }               
					   }
				   }
			   }
			   match_entries {
				   property = "Time"
				   operation = "Between"
				   type     = "time"
				   values {
					   inputitems {
						   name = "10"                    
					   }
				   }
				}
			 
			 match_entries {
				   property = "FQDN"
				   operation = "Is"
				   type     = "fqdn"
				   values {
					   inputitems {
						   name = "1234"                    
					   }
				   }
			   }
				match_entries {
				   property = "HTTP Method"
				   operation = "Is"
				   type     = "advanced"
				   values {
					   selecteditems {
						   name = "GET"                     
					   }
					   selecteditems {
						   name = "POST"                     
					   }
					   selecteditems {
						   name = "HEAD"                     
					   }
					   selecteditems {
						   name = "DELETE"                     
					   }
					   selecteditems {
						   name = "CONNECT"                     
					   }
				   }
				}
			 }
		   networks {
			   selecteditems {
				   name = "gcp-usw1-vpc1-tf"
			   }
		   }
		}
	}
`)
}

func testAccResourcePolicyPre() string {
	return fmt.Sprintf(`
	resource "prosimo_policy" "test_update" {
		name = "psonar-test-policy-transit"
		app_access_type = "transit"
		details {
		   actions = "allow"
		   # lock_users = false
		   # alert = true
		   # mfa = true
		   # lock_users = false
			 matches {
			   match_entries {
				   property = "URL"
				   operation = "Does NOT contain"
				   type     = "url"
				   values {
					   inputitems {
						   name = "1235678"                      
					   }
				   }
			   }
			   match_entries {
				   property = "Network"
				   operation = "Is"
				   type     = "networks"
				   values {
					   selecteditems {
						   name = "aws-uswest1-spoke1-infra-tf"                  
					   }
				   }
			   }
			   match_entries {
				   type     = "networkacl"
				   values {
					   inputitems {
						  ip_details {
						   source_ip = ["any"]
						   target_ip = ["any"]
						   protocol = ["tcp"]
						   source_port = ["any"]
						   target_port = ["any"]
						  }               
					   }
				   }
			   }
			   match_entries {
				   property = "Time"
				   operation = "Between"
				   type     = "time"
				   values {
					   inputitems {
						   name = "10"                    
					   }
				   }
				}
			 
			 match_entries {
				   property = "FQDN"
				   operation = "Is"
				   type     = "fqdn"
				   values {
					   inputitems {
						   name = "1234"                    
					   }
				   }
			   }
				match_entries {
				   property = "HTTP Method"
				   operation = "Is"
				   type     = "advanced"
				   values {
					   selecteditems {
						   name = "GET"                     
					   }
					   selecteditems {
						   name = "POST"                     
					   }
					   selecteditems {
						   name = "HEAD"                     
					   }
					   selecteditems {
						   name = "DELETE"                     
					   }
					   selecteditems {
						   name = "CONNECT"                     
					   }
				   }
				}
			 }
		}
	}
`)
}

func testAccResourcePolicyPost() string {
	return fmt.Sprintf(`
	resource "prosimo_policy" "test_update" {
		name = "psonar-test-policy-transit"
		app_access_type = "transit"
		details {
		   actions = "allow"
		   # lock_users = false
		   # alert = true
		   # mfa = true
		   # lock_users = false
			 matches {
			   match_entries {
				   property = "URL"
				   operation = "Does NOT contain"
				   type     = "url"
				   values {
					   inputitems {
						   name = "1235678"                      
					   }
				   }
			   }
			   match_entries {
				   property = "Network"
				   operation = "Is"
				   type     = "networks"
				   values {
					   selecteditems {
						   name = "aws-uswest1-spoke1-infra-tf"                  
					   }
				   }
			   }
			   match_entries {
				   type     = "networkacl"
				   values {
					   inputitems {
						  ip_details {
						   source_ip = ["any"]
						   target_ip = ["any"]
						   protocol = ["tcp"]
						   source_port = ["any"]
						   target_port = ["any"]
						  }               
					   }
				   }
			   }
			   match_entries {
				   property = "Time"
				   operation = "Between"
				   type     = "time"
				   values {
					   inputitems {
						   name = "10"                    
					   }
				   }
				}
			 
			 match_entries {
				   property = "FQDN"
				   operation = "Is"
				   type     = "fqdn"
				   values {
					   inputitems {
						   name = "1234"                    
					   }
				   }
			   }
				match_entries {
				   property = "HTTP Method"
				   operation = "Is"
				   type     = "advanced"
				   values {
					   selecteditems {
						   name = "GET"                     
					   }
					   selecteditems {
						   name = "POST"                     
					   }
					   selecteditems {
						   name = "HEAD"                     
					   }
					   selecteditems {
						   name = "DELETE"                     
					   }
					   selecteditems {
						   name = "CONNECT"                     
					   }
				   }
				}
			 }
		   networks {
			   selecteditems {
				   name = "gcp-usw1-vpc1-tf"
			   }
		   }
		}
	}
`)
}

func testAccCheckPolicyExists(resource string) resource.TestCheckFunc {
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

func testAccCheckPolicyDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "prosimo_policy" {
			continue
		}
		_, flag := client.GetPolicyfiltered(rs.Primary.ID)
		if !flag {
			return fmt.Errorf("Resource still avaible in portal %s", rs.Primary.ID)
		}
	}

	return nil
}
