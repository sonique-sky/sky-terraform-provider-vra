package vrealize

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize/api"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema:        providerSchema(),
		ConfigureFunc: providerConfig,
		ResourcesMap:  providerResources(),
	}
}

func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"username": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Tenant administrator username.",
		},
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Tenant administrator password.",
		},
		"tenant": {
			Type:     schema.TypeString,
			Required: true,
			Description: "Specifies the tenant URL token determined by the system administrator" +
				"when creating the tenant, for example, support.",
		},
		"host": {
			Type:     schema.TypeString,
			Required: true,
			Description: "host name.domain name of the vRealize Automation server, " +
				"for example, mycompany.mktg.mydomain.com.",
		},
		"insecure": {
			Type:        schema.TypeBool,
			Default:     false,
			Optional:    true,
			Description: "Specify whether to validate TLS certificates.",
		},
	}
}

func providerConfig(r *schema.ResourceData) (interface{}, error) {
	client := api.NewClient(r.Get("username").(string),
		r.Get("password").(string),
		r.Get("tenant").(string),
		r.Get("host").(string),
		r.Get("insecure").(bool),
	)

	err := client.Authenticate()

	if err != nil {
		return nil, fmt.Errorf("unable to get auth token: %v", err)
	}

	return &client, nil
}

func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"vra7_resource": ResourceMachine(),
	}
}
