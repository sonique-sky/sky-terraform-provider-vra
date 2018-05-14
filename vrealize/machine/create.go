package machine

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize/api"
)

func createResource(d *schema.ResourceData, meta interface{}) error {
	client := meta.(api.Client)

	var requestTemplate = new(api.RequestTemplate)
	var catalogErr = *new(error)

	if catalogId, idGiven := d.GetOk("catalog_id"); idGiven {
		requestTemplate, catalogErr = client.ReadCatalogByID(catalogId.(string))
	} else if catalogName, nameGiven := d.GetOk("catalog_name"); nameGiven {
		requestTemplate, catalogErr = client.ReadCatalogByName(catalogName.(string))
	} else {
		return fmt.Errorf("cannot retrieve catalog without 'catalog_id' or 'catalog_name'")
	}

	if catalogErr != nil {
		return fmt.Errorf("catalog Lookup failed %v", catalogErr)
	}

	log.Printf("createResource->requestTemplate %v\n", requestTemplate)

	catalogConfiguration, _ := d.Get("catalog_configuration").(map[string]interface{})
	for field1 := range catalogConfiguration {
		requestTemplate.Data[field1] = catalogConfiguration[field1]

	}
	log.Printf("createResource->requestTemplate.Data %v\n", requestTemplate.Data)

	//Get all resource keys from blueprint in array
	var keyList []string
	for field := range requestTemplate.Data {
		if reflect.ValueOf(requestTemplate.Data[field]).Kind() == reflect.Map {
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
				requestTemplate.Data[keyList[dataKey]], replaced = changeTemplateValue(
					requestTemplate.Data[keyList[dataKey]].(map[string]interface{}),
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
				resourceItem := requestTemplate.Data[keyList[dataKey]].(map[string]interface{})
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
			requestTemplate.Description = deploymentConfiguration[depField].(string)
		case "reasons":
			requestTemplate.Reasons = deploymentConfiguration[depField].(string)
		default:
			log.Printf("unknown option [%s] with value [%s] ignoring\n", depField, deploymentConfiguration[depField])
		}
	}
	log.Printf("Updated template - %v\n", requestTemplate.Data)

	requestMachine, err := client.RequestMachine(requestTemplate)

	if err != nil {
		return fmt.Errorf("resource machine request failed: %v", err)
	}

	waitTimeout := d.Get("wait_timeout").(int) * 60

	request := new(api.Request)
	for i := 0; i < waitTimeout/30; i++ {
		time.Sleep(3e+10)

		request, _ = client.GetRequest(requestMachine.ID)
		if request.Phase == "FAILED" {
			return fmt.Errorf("instance got failed while creating - check detail for more information")
		}
		if request.Phase == "SUCCESSFUL" {
			resourceViews, e := client.GetResource(requestMachine.ID, "Infrastructure.Virtual")
			if e != nil {
				return e
			}

			if len(resourceViews.Resources) == 0 {
				return fmt.Errorf("could not find expected resource")
			}

			resource := resourceViews.Resources[0]
			d.SetId(resource.ID)
			readResource(d, meta)
			return nil
		}
	}

	if request == nil || request.Phase == "IN_PROGRESS" {
		//If request is in_progress state during the time then
		//keep resource details in state files and throw an error
		//so that the child resource won't go for create call.
		//If execution gets timed-out and request is in progress
		//then dependent machine won't be get created in this iteration.
		//A user needs to ensure that the request should be a success state
		//using terraform refresh command and hit terraform apply again.
		return fmt.Errorf("resource is still being created")
	}

	return nil
}
