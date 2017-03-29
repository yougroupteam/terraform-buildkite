# Terraform provider for [buildkite](https://www.buildkite.com)

This allows you to manage buildkite pipelines with Terraform.

## Installation

Run
```bash
go get github.com/saymedia/terraform-buildkite/terraform-provider-buildkite
go install github.com/saymedia/terraform-buildkite/terraform-provider-buildkite
```
Which gives you a `terraform-provider-buildkite` in `$GOPATH/bin`.

Fom there you have a couple different options for installing that so Terraform can find it:

* Put it in the same directory as the terraform program itself.
* [Create a `.terraformrc` file](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) in your home directory (or, in CI, the home directory of whatever user runs terraform) that tells Terraform where to find the program.
As an example of the latter:

```terraform
providers {
  buildkite = "/usr/local/bin/terraform-providers-buildkite"
}
```

## Usage

```terraform
provider "buildkite" {
  # Get an API token from https://buildkite.com/user/api-access-tokens
  # Needs: read_pipelines, write_pipelines
  # Instead of embedding the API token in the .tf file,
  # it can also be passed via env variable BUILDKITE_API_TOKEN
  api_token    = "YOUR_API_TOKEN"
  # This is the part behind https://buildkite.com/, e.g. https://buildkite.com/some-org
  # Instead of embedding the org slug in the .tf file,
  # it can also be passed via env variable BUILDKITE_ORGANIZATION
  organization = "YOUR_ORG_SLUG"
}

resource "buildkite_pipeline" "terraform_test" {
  name       = "terraform-test"
  repository = "git@github.com:you/repo.git"

  step = [
    {
      type    = "script"
      name    = ":llama: Tests"
      command = "echo Hello world!"
    },
  ]
}
```

## Importing existing pipelines

You can import existing pipeline definitions by their slug:

```bash
terraform import buildkite_pipeline.my_name my-pipeline-slug
```

## Local development of this provider

To do local development you will most likely be working in a Github fork of the repository. After creating your fork
you can add it as a remote on your local repository in GOPATH:

* `cd $GOPATH/src/github.com/saymedia/terraform-buildkite`
* `git remote add mine git@github.com:yourname/terraform-buildkite`
* `git checkout -b yourbranch`
* `git push -u mine yourbranch`

After this you should be able to `git push` to your fork, and eventually open a PR if you like.

You can build like this:

* `go install github.com/saymedia/terraform-buildkite/terraform-provider-buildkite`

This should produce a file at `$GOPATH/bin/terraform-provider-buildkite`. To use this with Terraform you'll need to either move
that binary to the same directory where `terraform` lives or
[create a `.terraformrc` file](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) to help Terraform find
this file.

You can see debug output via `TF_LOG=DEBUG terraform plan`
