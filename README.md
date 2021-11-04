# Terraform provider for running Bitbucket pipelines

Run Bitbucket pipelines straight from Terraform!

## WARNING 

__This is an experimental provider__

Adaptavist are using this internally to test how we can provide services to our engineering teams via BitBucket
pipelines. Use at your own risk.

## Usage

```terraform
provider "bpr" {
  // Optional - set the default workspace for all bpr resources
  workspace = "the-owner-account"
  // Optional - set the default repository for all bpr resources
  repository = "youre-repo-slug"
  // Set username with BITBUCKET_USERNAME
  // Set password with BITBUCKET_PASSWORD or BITBUCKET_APP_PASSWORD
}

resource "bpr_run" "this" {
  // Optional - defaults to the providers workspace if empty
  workspace = "the-owner-account"
  // Optional - defaults to the provider repository if empty
  repository = "youre-repo-slug"
  // required - unless branch is set
  tag = "v1.0.0"
  // Required - unless tag is set
  branch = "master"
  // Optional - defaults to "default"
  pipeline = "custom_pipeline"
  // Optional - variables you want to pass to the pipeline as environment variables, values are also encoded to JSON
  variables = {
    VAR_NAME = "VAR_VALUE"
  }
}

// Outputs are optional and are only available if the pipeline outputs them using the conventions specified later.
output "this" {
  value = bpr_run.this.outputs
}
```

We strongly recommend you have a backend or commit the state file when using this provider, otherwise Terraform will keep
running your pipelines and using all your units up!

## Supporting outputs

Currently, output support is a quick and dirty hack using substrings and JSON decoding which could and should be
improved, but for now here's the current convention.

Output from your pipelines should render like so.

```text
--- OUTPUT JSON START ---
YOUR JSON OUTPUT HERE
--- OUTPUT JSON STOP ---
```

This isn't easily achieved in Bitbucket pipelines as it will print your command because printing its output causing
unwanted behavior. So our best advice is to create a wrapper script to print your outputs and call it from your 
pipeline - similar to below.

__terraform-output.sh__

```bash
#!/usr/bin/env sh
echo --- OUTPUT JSON START ---
terraform output -json
echo --- OUTPUT JSON STOP ---
```

__bitbucket-pipelines.yml__

```bash
pipelines:
  custom:
    do_a_thing:
    - step:
        script:
        - sh terraform-output.sh
```

### Complex output types

These are supported by the provider already, however they do get flattened before being made available.

#### Example

Output looking like this

```json
{
  "nested_object": {
    "nested_value": "example"
  }
}
```

Would look like this in the provider

```go
map[string]interface{}{
	"nested_object.nested_value": "example"
}
```

Fortunately this has no impact on how you access the outputs within Terraform like so

```terraform
output "this" {
  value = bpr_run.this.outputs.nested_object.nested_value
}
```