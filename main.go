package main

import (
	"github.com/3pjgames/terraform-provider-ucloud/ucloud"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ucloud.Provider,
	})
}
