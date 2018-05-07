package machine

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize/api"
)

func readResource(d *schema.ResourceData, meta interface{}) error {
	client := meta.(api.Client)

	requestMachineID := d.Id()
	resourceTemplate, errTemplate := client.GetRequestStatus(requestMachineID)

	if errTemplate != nil {
		return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}

	d.Set("request_status", resourceTemplate.Phase)

	if resourceTemplate.Phase == "FAILED" {
		d.Set("failed_message", resourceTemplate.RequestCompletion.CompletionDetails)
	}
	return nil
}
