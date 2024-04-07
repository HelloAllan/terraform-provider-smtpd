package smtpd

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
		"smtpd": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("SMTPD_API_KEY"); err == "" {
		t.Fatal("SMTPD_API_KEY must be set for acceptance tests")
	}
	if err := os.Getenv("SMTPD_API_SECRET"); err == "" {
		t.Fatal("SMTPD_API_SECRET must be set for acceptance tests")
	}
}
