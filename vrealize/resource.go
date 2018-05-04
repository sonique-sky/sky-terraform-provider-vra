package vrealize

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize/api"
)

func ResourceMachine() *schema.Resource {
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
		"catalog_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"catalog_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"businessgroup_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"wait_timeout": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  15,
		},
		"request_status": {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
		"failed_message": {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
			Optional: true,
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
		"catalog_configuration": {
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

//Function use - to create machine
//Terraform call - terraform apply
func changeTemplateValue(templateInterface map[string]interface{}, field string, value interface{}) (map[string]interface{}, bool) {
	var replaced bool
	//Iterate over the map to get field provided as an argument
	for i := range templateInterface {
		//If value type is map then set recursive call which will fiend field in one level down of map interface
		if reflect.ValueOf(templateInterface[i]).Kind() == reflect.Map {
			template, _ := templateInterface[i].(map[string]interface{})
			templateInterface[i], replaced = changeTemplateValue(template, field, value)
		} else if i == field {
			//If value type is not map then compare field name with provided field name
			//If both matches then update field value with provided value
			templateInterface[i] = value
			return templateInterface, true
		}
	}
	//Return updated map interface type
	return templateInterface, replaced
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

//Function use - to set a create resource call
//Terraform call - terraform apply
func createResource(d *schema.ResourceData, meta interface{}) error {
	//Log file handler to generate logs for debugging purpose
	//Get client handle
	client := meta.(*api.Client)

	//If catalog_name and catalog_id both not provided then throw an error
	if len(d.Get("catalog_name").(string)) <= 0 && len(d.Get("catalog_id").(string)) <= 0 {
		return fmt.Errorf("Either catalog_name or catalog_id should be present in given configuration")
	}

	if len(d.Get("catalog_name").(string)) > 0 {
		catalogID, returnErr := client.ReadCatalogIDByName(d.Get("catalog_name").(string))
		log.Printf("createResource->catalog_id %v\n", catalogID)
		if returnErr != nil {
			return fmt.Errorf("%v", returnErr)
		}
		if catalogID == nil || catalogID == "" {
			return fmt.Errorf("no catalog found with name %v", d.Get("catalog_name").(string))
		}
		d.Set("catalog_id", catalogID.(string))
	} else if len(d.Get("catalog_id").(string)) > 0 {
		CatalogName, nameError := client.ReadCatalogNameByID(d.Get("catalog_id").(string))
		if nameError != nil {
			return fmt.Errorf("%v", nameError)
		}
		if nameError != nil {
			d.Set("catalog_name", CatalogName.(string))
		}
	}
	//Get catalog blueprint
	templateCatalogItem, err := client.GetCatalogItem(d.Get("catalog_id").(string))
	log.Printf("createResource->templateCatalogItem %v\n", templateCatalogItem)

	catalogConfiguration, _ := d.Get("catalog_configuration").(map[string]interface{})
	for field1 := range catalogConfiguration {
		templateCatalogItem.Data[field1] = catalogConfiguration[field1]

	}
	log.Printf("createResource->templateCatalogItem.Data %v\n", templateCatalogItem.Data)

	if len(d.Get("businessgroup_id").(string)) > 0 {
		templateCatalogItem.BusinessGroupID = d.Get("businessgroup_id").(string)
	}

	//Get all resource keys from blueprint in array
	var keyList []string
	for field := range templateCatalogItem.Data {
		if reflect.ValueOf(templateCatalogItem.Data[field]).Kind() == reflect.Map {
			keyList = append(keyList, field)
		}
	}
	log.Printf("createResource->key_list %v\n", keyList)

	//Arrange keys in descending order of text length
	for field1 := range keyList {
		for field2 := range keyList {
			if len(keyList[field1]) > len(keyList[field2]) {
				temp := keyList[field1]
				keyList[field1], keyList[field2] = keyList[field2], temp
			}
		}
	}

	//array to keep track of resource values that have been used
	usedConfigKeys := []string{}
	var replaced bool

	//Update template field values with user configuration
	resourceConfiguration, _ := d.Get("resource_configuration").(map[string]interface{})
	for configKey := range resourceConfiguration {
		for dataKey := range keyList {
			//compare resource list (resource_name) with user configuration fields (resource_name+field_name)
			if strings.Contains(configKey, keyList[dataKey]) {
				//If user_configuration contains resource_list element
				// then split user configuration key into resource_name and field_name
				splitedArray := strings.SplitN(configKey, keyList[dataKey]+".", 2)
				if len(splitedArray) != 2 {
					return fmt.Errorf("resource_configuration key is not in correct format. Expected %s to start with %s\n", configKey, keyList[dataKey]+".")
				}
				//Function call which changes the template field values with  user values
				templateCatalogItem.Data[keyList[dataKey]], replaced = changeTemplateValue(
					templateCatalogItem.Data[keyList[dataKey]].(map[string]interface{}),
					splitedArray[1],
					resourceConfiguration[configKey])
				if replaced {
					usedConfigKeys = append(usedConfigKeys, configKey)
				} else {
					log.Printf("%s was not replaced", configKey)
				}
			}
		}
	}

	//Add remaining keys to template vs updating values
	// first clean out used values
	for usedKey := range usedConfigKeys {
		delete(resourceConfiguration, usedConfigKeys[usedKey])
	}
	log.Println("Entering Add Loop")
	for configKey2 := range resourceConfiguration {
		for dataKey := range keyList {
			log.Printf("Add Loop: configKey2=[%s] keyList[%d] =[%v]", configKey2, dataKey, keyList[dataKey])
			if strings.Contains(configKey2, keyList[dataKey]) {
				splitArray := strings.Split(configKey2, keyList[dataKey]+".")
				log.Printf("Add Loop Contains %+v", splitArray[1])
				resourceItem := templateCatalogItem.Data[keyList[dataKey]].(map[string]interface{})
				resourceItem = addTemplateValue(
					resourceItem["data"].(map[string]interface{}),
					splitArray[1],
					resourceConfiguration[configKey2])
			}
		}
	}
	//update template with deployment level config
	// limit to description and reasons as other things could get us into trouble
	deploymentConfiguration, _ := d.Get("deployment_configuration").(map[string]interface{})
	for depField := range deploymentConfiguration {
		fieldstr := fmt.Sprintf("%s", depField)
		switch fieldstr {
		case "description":
			templateCatalogItem.Description = deploymentConfiguration[depField].(string)
		case "reasons":
			templateCatalogItem.Reasons = deploymentConfiguration[depField].(string)
		default:
			log.Printf("unknown option [%s] with value [%s] ignoring\n", depField, deploymentConfiguration[depField])
		}
	}
	//Log print of template after values updated
	log.Printf("Updated template - %v\n", templateCatalogItem.Data)

	//Through an exception if there is any error while getting catalog template
	if err != nil {
		return fmt.Errorf("Invalid CatalogItem ID %v", err)
	}

	//Set a  create machine function call
	requestMachine, err := client.RequestMachine(templateCatalogItem)

	//Check if error got while create machine call
	//If Error is occured, through an exception with an error message
	if err != nil {
		return fmt.Errorf("Resource Machine Request Failed: %v", err)
	}

	//Set request ID
	d.SetId(requestMachine.ID)
	//Set request status
	d.Set("request_status", "SUBMITTED")

	waitTimeout := d.Get("wait_timeout").(int) * 60

	for i := 0; i < waitTimeout/30; i++ {
		time.Sleep(3e+10)
		readResource(d, meta)

		if d.Get("request_status") == "SUCCESSFUL" {
			return nil
		}
		if d.Get("request_status") == "FAILED" {
			//If request is failed during the time then
			//unset resource details from state.
			d.SetId("")
			return fmt.Errorf("instance got failed while creating." +
				" kindly check detail for more information")
		}
	}
	if d.Get("request_status") == "IN_PROGRESS" {
		//If request is in_progress state during the time then
		//keep resource details in state files and throw an error
		//so that the child resource won't go for create call.
		//If execution gets timed-out and status is in progress
		//then dependent machine won't be get created in this iteration.
		//A user needs to ensure that the status should be a success state
		//using terraform refresh command and hit terraform apply again.
		return fmt.Errorf("resource is still being created")
	}

	return nil
}

//Function use - to update centOS 6.3 machine present in state file
//Terraform call - terraform refresh
func updateResource(d *schema.ResourceData, meta interface{}) error {
	log.Println(d)
	return nil
}

func readResource(d *schema.ResourceData, meta interface{}) error {
	requestMachineID := d.Id()
	client := meta.(*api.Client)
	resourceTemplate, errTemplate := client.GetRequestStatus(requestMachineID)

	if errTemplate != nil {
		return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}

	d.Set("request_status", resourceTemplate.Phase)
	//If request is failed then set failed message in state file
	if resourceTemplate.Phase == "FAILED" {
		d.Set("failed_message", resourceTemplate.RequestCompletion.CompletionDetails)
	}
	return nil
}

//Function use - To delete resources which are created by terraform and present in state file
//Terraform call - terraform destroy
func deleteResource(d *schema.ResourceData, meta interface{}) error {
	requestMachineID := d.Id()
	client := meta.(*api.Client)

	if len(d.Id()) == 0 {
		return fmt.Errorf("resource not found")
	}
	//If resource create status is in_progress then skip delete call and through an exception
	if d.Get("request_status").(string) != "SUCCESSFUL" {
		if d.Get("request_status").(string) == "FAILED" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("machine cannot be deleted while in-progress state. Please try later")

	}
	//Fetch machine template
	templateResources, errTemplate := client.GetResourceViews(requestMachineID)

	client.DestroyMachine(templateResources)

	if errTemplate != nil {
		return fmt.Errorf("Resource view failed to load:  %v", errTemplate)
	}

	d.SetId("")
	return nil
}
