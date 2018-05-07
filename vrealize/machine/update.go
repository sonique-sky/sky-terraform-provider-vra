package machine

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func updateResource(d *schema.ResourceData, meta interface{}) error {
	log.Println(d)
	return nil
}

