
# Terraform Provider Prosimo.

* Website: https://prosimo.io

The Terraform provider for Prosimo is a Terraform plugin to enable full
lifecycle management of Prosimo resources. The provider is maintained
internally by Prosimo engineering team.


## Development Requirements

-	[Terraform Latest Release](https://developer.hashicorp.com/terraform/downloads)
-	[Go](https://golang.org/doc/install) 1.16+

Building The Provider
---------------------
- Clone the repo and get into terraform-provider-prosimo folder.
- In the Make file update your OS and os architecture details.
### MacOS / Linux
```sh
GO := GOPATH=$(GOPATH) GOPRIVATE=git.prosimo.io/prosimoio  GOOS=darwin/linux GOARCH=amd64 go
$ make install
```

### Windows
```sh
go fmt
go install
```

Using The Provider (Terraform v1.0+)
------------------------------------
To use a local built version follow the steps below.
Steps:

- Download and install Terraform.

- Add the following snippet to every module that you want to use the provider.
```hcl
terraform {
  required_providers {
    prosimo = {
      version = "x.0.0"
      source  = "prosimo.com/prosimo/prosimo"
    }
  }
}

```
Examples
--------

- Look at [docs] folder for documentation and examples.


