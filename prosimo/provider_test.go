package prosimo

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"prosimo": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("PROSIMO_BASE_URL"); v == "" {
		t.Fatal("PROSIMO_BASE_URL must be set for acceptance tests.")
	}
	if v := os.Getenv("PROSIMO_TOKEN"); v == "" {
		t.Fatal("PROSIMO_TOKEN must be set for acceptance tests.")
	}
	if v := os.Getenv("PROSIMO_INSECURE"); v == "" {
		t.Fatal("PROSIMO_INSECURE must be set for acceptance tests.")
	}
}
