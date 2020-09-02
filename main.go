package main

import (
	"github.com/hashicorp/terraform/plugin"
	gitops "github.com/tyler-technologies/terraform-provider-gitops/gitops"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gitops.Provider,
	})
}
