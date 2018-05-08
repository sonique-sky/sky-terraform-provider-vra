package machine

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize/api"
)

func readResource(d *schema.ResourceData, meta interface{}) error {
	client := meta.(api.Client)

	resourceID := d.Id()
	resourceTemplate, errTemplate := client.GetRequestStatus(resourceID)

	if errTemplate != nil {
		return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}

	d.Set("request_status", resourceTemplate.Phase)

	if resourceTemplate.Phase == "FAILED" {
		d.Set("failed_message", resourceTemplate.RequestCompletion.CompletionDetails)
	}
	return nil
}

func readRequestStatus(requestId string, meta interface{}) (string, error) {
	client := meta.(api.Client)

	resourceTemplate, errTemplate := client.GetRequestStatus(requestId)

	if errTemplate != nil {
		return "", fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}
	return resourceTemplate.Phase, nil
}
