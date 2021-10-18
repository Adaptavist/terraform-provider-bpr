package bitbucket

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testBitbucketPipelinesRun = `
resource "bpr_run" "this" {
	workspace  = "my-workspace"
	repository = "my-repository"
	pipeline   = "custom"
	tag        = "v0.0.0"
	variables  = {
		"var_key" = "var_value"
	}
}
`

func TestBitbucketPipelines_Run(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testBitbucketPipelinesRun,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("bpr_run.this", "pipeline", regexp.MustCompile("custom")),
					resource.TestMatchResourceAttr("bpr_run.this", "outputs.string_output", regexp.MustCompile("value")),
					resource.TestMatchResourceAttr("bpr_run.this", "outputs.map_output.map_key", regexp.MustCompile("value")),
				),
			},
		},
	})
}

//{
//	"key_output": "value",
//	"map_output": {
//	"map_key": "map_value"
//	}
//}
