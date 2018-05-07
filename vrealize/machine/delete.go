package machine

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize/api"
	"fmt"
)

func deleteResource(d *schema.ResourceData, meta interface{}) error {
	requestMachineID := d.Id()
	client := meta.(api.DeleteClient)

	if len(d.Id()) == 0 {
		return fmt.Errorf("resource not found")
	}

	if d.Get("request_status").(string) != "SUCCESSFUL" {
		if d.Get("request_status").(string) == "FAILED" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("machine cannot be deleted while in-progress state. Please try later")

	}

	templateResources, errTemplate := client.GetResourceViews(requestMachineID)

	if errTemplate != nil {
		return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}

	client.DestroyMachine(templateResources)

	d.SetId("")
	return nil
}