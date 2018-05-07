package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize"
)

func main() {
	opts := plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return vrealize.Provider()
		},
	}

	plugin.Serve(&opts)
}