package machine

import (
	"reflect"
	"github.com/hashicorp/terraform/helper/schema"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Create: createResource,
		Read:   readResource,
		Update: updateResource,
		Delete: deleteResource,
		Schema: setResourceSchema(),

	}
}

func setResourceSchema() map[string]*schema.Schema {
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

func changeTemplateValue(templateInterface map[string]interface{}, field string, value interface{}) bool {
	var replaced bool
	//Iterate over the map to get field provided as an argument
	for key, val := range templateInterface {
		//If value type is map then set recursive call which will find field in one level down of map interface
		if reflect.ValueOf(val).Kind() == reflect.Map {
			template, _ := val.(map[string]interface{})
			replaced = changeTemplateValue(template, field, value)
			templateInterface[key] = template
		} else if key == field {
			templateInterface[key] = value
			return true
		}
	}
	return replaced
}

//modeled after changeTemplateValue, for values being added to template vs updating existing ones
func addTemplateValue(templateInterface map[string]interface{}, field string, value interface{}) map[string]interface{} {
	//simplest case is adding a simple value. Leaving as a func in case there's a need to do more complicated additions later
	//	templateInterface[data]
	for i := range templateInterface {
		if reflect.ValueOf(templateInterface[i]).Kind() == reflect.Map && i == "data" {
			template, _ := templateInterface[i].(map[string]interface{})
			templateInterface[i] = addTemplateValue(template, field, value)
		} else { //if i == "data" {
			templateInterface[field] = value
		}
	}
	//Return updated map interface type
	return templateInterface
}
