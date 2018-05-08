package machine

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize/api"
)

func readResource(d *schema.ResourceData, meta interface{}) error {
	client := meta.(api.BaseClient)

	resourceID := d.Id()

	client.GetMachine(resourceID)
	return nil
}
