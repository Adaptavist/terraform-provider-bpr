package bitbucket

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"
)

var providerFactories = map[string]func() (*schema.Provider, error){
	ProviderName: func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}


func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("BITBUCKET_USERNAME"); err == "" {
		t.Fatal("BITBUCKET_USERNAME must be set for acceptance tests")
	}
	if err := os.Getenv("BITBUCKET_PASSWORD"); err == "" {
		t.Fatal("BITBUCKET_PASSWORD must be set for acceptance tests")
	}
	if err := os.Getenv("BITBUCKET_BASE_URL"); err == "" {
		t.Fatal("BITBUCKET_BASE_URL must be set for acceptance tests")
	}
}