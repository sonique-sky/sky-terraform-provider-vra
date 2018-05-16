package vrealize

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"fmt"
	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize/api"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Create: createResource,
		Read:   readResource,
		Update: updateResource,
		Delete: deleteResource,
		Schema: resourceSchema(),
	}
}

type TerraformRequest struct {
	RequestId               string
	CatalogId               string
	CatalogName             string
	WaitTimeOut             int
	CatalogConfiguration    map[string]interface{}
	DeploymentConfiguration map[string]interface{}
	ResourceConfiguration   map[string]interface{}
}

func createResource(d *schema.ResourceData, meta interface{}) error {
	var terraformRequest = new(TerraformRequest)

	if catalogId, idGiven := d.GetOk("catalog_id"); idGiven {
		terraformRequest.CatalogId = catalogId.(string)
	} else if catalogName, nameGiven := d.GetOk("catalog_name"); nameGiven {
		terraformRequest.CatalogName = catalogName.(string)
	} else {
		return fmt.Errorf("cannot retrieve catalog without 'catalog_id' or 'catalog_name'")
	}

	if catalogConfiguration, given := d.GetOk("catalog_configuration"); given {
		terraformRequest.CatalogConfiguration = catalogConfiguration.(map[string]interface{})
	}

	if resourceConfiguration, given := d.GetOk("resource_configuration"); given {
		terraformRequest.ResourceConfiguration = resourceConfiguration.(map[string]interface{})
	}

	if deploymentConfiguration, given := d.GetOk("deployment_configuration"); given {
		terraformRequest.DeploymentConfiguration = deploymentConfiguration.(map[string]interface{})
	}

	if waitTimeout, given := d.GetOk("wait_timeout"); given {
		terraformRequest.WaitTimeOut = waitTimeout.(int)
	}

	resource, err := Create(meta.(api.Client), terraformRequest)

	if err != nil {
		return err
	}

	if resource != nil {
		d.SetId(resource.ID)
		readResource(d, meta)
	}
	return nil
}

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

	d.Set("hostname", resource.Name)
	d.Set("request_id", resource.RequestID)
	//d.Set("description", resource.Description)

	if val, found := resource.StringValue("ip_address"); found {
		d.Set("ip_address", val)
	}

	request, reqErr := client.GetRequest(resource.RequestID)

	if reqErr != nil {
		return reqErr
	}

	d.Set("foo", request.Phase)
	return nil
}

func updateResource(d *schema.ResourceData, meta interface{}) error {
	log.Println(d)
	return nil
}

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
