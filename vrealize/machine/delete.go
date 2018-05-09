package machine

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize/api"
	"fmt"
)

func deleteResource(d *schema.ResourceData, meta interface{}) error {
	requestMachineID := d.Id()
	client := meta.(api.Client)

	if len(d.Id()) == 0 {
		return fmt.Errorf("resource not found")
	}

	client.DestroyMachine(requestMachineID)

	d.SetId("")
	return nil
}
