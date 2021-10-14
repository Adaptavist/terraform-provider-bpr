provider "bpr" {
  workspace = "my-org"
}

resource "bpr_run" "account" {
  repository = "my-account-repo"
  tag        = "v1.0.0"
  pipeline   = "deploy"
  variables = {
    ENV_VAR_KEY = "value"
  }
}