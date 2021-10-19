package bitbucket

import (
	"context"
	"github.com/adaptavist/bitbucket_pipelines_client/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ProviderName = "bpr"

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"base_url": {
					Type:     schema.TypeString,
					Optional: true,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{
						"BITBUCKET_BASE_URL",
					}, "https://api.bitbucket.org"),
					Description: "Bitbucket pipelines app password",
				},
				"username": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_USERNAME", nil),
					Description: "Bitbucket pipelines username",
				},
				"password": {
					Type:     schema.TypeString,
					Optional: true,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{
						"BITBUCKET_PASSWORD",
						"BITBUCKET_APP_PASSWORD",
					}, nil),
					Description: "Bitbucket pipelines app password",
				},
				"workspace": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_WORKSPACE", nil),
					Description: "Bitbucket workspace slug",
				},
				"repository": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_REPOSITORY", nil),
					Description: "Bitbucket repository slug",
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"bpr_run": resourceBitbucketPipelineRun(),
			},
		}
		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		// Setup a User-Agent for your API client (replace the provider name for yours):
		// userAgent := p.UserAgent("terraform-provider-scaffolding", version)
		// TODO: myClient.UserAgent = userAgent

		return &client.Client{
			Config: &client.Config{
				Username:   d.Get("username").(string),
				Password:   d.Get("password").(string),
				BaseURL:    d.Get("base_url").(string),
				Workspace:  valueAsPointer("workspace", d),
				Repository: valueAsPointer("repository", d),
			},
		}, nil
	}
}

func valueAsPointer(k string, d *schema.ResourceData) *string {
	if v := d.Get(k).(string); v != "" {
		return &v
	}
	return nil
}
