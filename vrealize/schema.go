package vrealize

import "github.com/hashicorp/terraform/helper/schema"

func resourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"catalog_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"hostname": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"catalog_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"wait_timeout": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  15,
		},
		"deployment_configuration": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
		},
		"resource_configuration": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
		},
	}
}

