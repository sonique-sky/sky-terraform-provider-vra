package machine

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize/api"
)

func readResource(d *schema.ResourceData, meta interface{}) error {
	client := meta.(api.BaseClient)

	resourceID := d.Id()

	resource, err := client.GetMachine(resourceID)
	if err != nil {
		return err
	}

	d.Set("hostname", resource.Name)
	d.Set("description", resource.Description)

	if val, found := resource.IntValue("cpu"); found {
		d.Set("cpu", val)
	}

	if val, found := resource.StringValue("ip_address"); found {
		d.Set("ip_address", val)
	}

	return nil
}
