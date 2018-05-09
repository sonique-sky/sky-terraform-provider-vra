package machine

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize/api"
)

func readResource(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Start Read")
	client := meta.(api.Client)


	resourceID := d.Id()
	log.Printf("Resource ID: %s", resourceID)

	resource, err := client.GetMachine(resourceID)
	if err != nil {
		log.Printf("Read Got Error: %v", err)
		return err
	}
	log.Printf("Read resource: %v", resource)

	//d.Set("hostname", resource.Name)
	//d.Set("description", resource.Description)

	//if val, found := resource.IntValue("cpu"); found {
	//	d.Set("cpu", val)
	//}
	//
	//if val, found := resource.StringValue("ip_address"); found {
	//	d.Set("ip_address", val)
	//}

	return nil
}
