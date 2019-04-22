package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/yougroupteam/terraform-buildkite/buildkite"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: buildkite.Provider,
	})
}
