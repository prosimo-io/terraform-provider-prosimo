---
page_title: "Prosimo Terraform Provider"
subcategory: ""
description: |-
    This repo is used to manage and connfigure prosimo cloud infrastructure resources using terraform.
---

# Provider Initialization.


```hcl
terraform {
  required_providers {
    prosimo = {
      version = "2.9.2"
      source  = ""hashicorp/prosimo""
    }
  }
}

```


# Provider Authentication
To Authenticate with Prosimo control plane you need to generate a token from Prosimo Dashboard.
As part of authentication the credentials can be passed in couple of ways.
Static credentials can be passed as below, However this is not recommended for production usage.
```hcl
provider "prosimo" {
  base_url ="https://your_team_name.admin.prosimo.io"
  token    = "token""
  insecure = false  
}
```

You can also set  credentials by exporting them as environmental variables.

```hcl
export PROSIMO_BASE_URL="https://your_team_name.admin.prosimo.io"
export PROSIMO_TOKEN="token"
export PROSIMO_INSECURE=False
```
And then configure the provider:
```hcl
provider "prosimo" {}
```
## Schema

### Required

- `token` (String) The API token used to connect to Prosimo Dashboard. This token can be generated on the Prosimo Dashboard. If you have any questions on this, please reach out to us or go through the Prosimo documentation.
- `base_url` (String) The URL that is used to access the Prosimo dashboard. For ex. https://prosimo_dashboard_name.admin.prosimo.io. Please replace prosimo_dashboard_name with your actual dashboard name.

### Optional

- `insecure` (Boolean) Defaults to False. Enable `insecure` mode only for testing purposes

