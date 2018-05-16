package vrealize

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"sort"
	"time"

	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize/api"
)

func Create(client api.Client, terraformRequest *TerraformRequest) (*api.Resource, error) {
	var requestTemplate = new(api.RequestTemplate)
	var catalogErr = *new(error)

	if len(terraformRequest.CatalogId) > 0 {
		requestTemplate, catalogErr = client.ReadCatalogByID(terraformRequest.CatalogId)
	} else {
		requestTemplate, catalogErr = client.ReadCatalogByName(terraformRequest.CatalogName)
	}

	if catalogErr != nil {
		return nil, fmt.Errorf("catalog Lookup failed %v", catalogErr)
	}
	log.Printf("createResource->requestTemplate %v\n", requestTemplate)

	catalogConfiguration := terraformRequest.CatalogConfiguration
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
	sort.Slice(keyList, func(i1, i2 int) bool {
		return len(keyList[i1]) > len(keyList[i2])
	})

	//array to keep track of resource values that have been used
	usedConfigKeys := []string{}

	//Update template field values with user configuration
	resourceConfiguration := terraformRequest.ResourceConfiguration
	for configKey, configValue := range resourceConfiguration {
		for _, dataKey := range keyList {
			templateData := requestTemplate.Data[dataKey].(map[string]interface{})
			if field, found := checkKey(dataKey, configKey); found {
				if changeTemplateValue(templateData, field, configValue) {
					usedConfigKeys = append(usedConfigKeys, configKey)
				} else {
					log.Printf("%s was not replaced", configKey)
				}
			} else {
				return nil, fmt.Errorf("resource_configuration key is not in correct format. Expected %s to start with %s\n", configKey, dataKey+".")
			}
		}
	}

	//Add remaining keys to template vs updating values
	// first clean out used values
	for usedKey := range usedConfigKeys {
		delete(resourceConfiguration, usedConfigKeys[usedKey])
	}
	log.Println("Entering Add Loop")
	for configKey, configVal := range resourceConfiguration {
		for _, dataKey := range keyList {
			log.Printf("Add Loop: configKey2=[%s] keyList[%d] =[%v]", configKey, dataKey, dataKey)
			if strings.Contains(configKey, dataKey) {
				splitArray := strings.Split(configKey, dataKey+".")
				log.Printf("Add Loop Contains %+v", dataKey)
				resourceItem := requestTemplate.Data[dataKey].(map[string]interface{})
				resourceItem = addTemplateValue(
					resourceItem["data"].(map[string]interface{}),
					splitArray[1],
					configVal)
			}
		}
	}

	//update template with deployment level config
	// limit to description and reasons as other things could get us into trouble
	deploymentConfiguration := terraformRequest.DeploymentConfiguration
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
		return nil, fmt.Errorf("resource machine request failed: %v", err)
	}

	waitTimeout := terraformRequest.WaitTimeOut * 60

	request := new(api.Request)
	for i := 0; i < waitTimeout/30; i++ {
		time.Sleep(3e+10)

		request, _ = client.GetRequest(requestMachine.ID)
		if request.Phase == "FAILED" {
			return nil, fmt.Errorf("instance got failed while creating - check detail for more information")
		}
		if request.Phase == "SUCCESSFUL" {
			resourceViews, e := client.GetRequestResource(requestMachine.ID, "Infrastructure.Virtual")
			if e != nil {
				return nil, e
			}

			if len(resourceViews.Resources) == 0 {
				return nil, fmt.Errorf("could not find expected resource")
			}

			resource := resourceViews.Resources[0]
			return &resource, nil
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
		return nil, fmt.Errorf("resource is still being created")
	}
	return nil, nil
}

func checkKey(dataKey, configKey string) (string, bool) {
	pattern, _ := regexp.Compile("^" + dataKey + "\\.(.*)")
	res := pattern.FindStringSubmatch(configKey)
	if len(res) == 0 {
		return "", false
	}
	return res[1], true
}

func changeTemplateValue(templateInterface map[string]interface{}, field string, value interface{}) bool {
	var replaced bool
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

//Function use - to create machine
//Terraform call - terraform apply
func oldChangeTemplateValue(templateInterface map[string]interface{}, field string, value interface{}) (map[string]interface{}, bool) {
	var replaced bool
	//Iterate over the map to get field provided as an argument
	for i := range templateInterface {
		//If value type is map then set recursive call which will fiend field in one level down of map interface
		if reflect.ValueOf(templateInterface[i]).Kind() == reflect.Map {
			template, _ := templateInterface[i].(map[string]interface{})
			templateInterface[i], replaced = oldChangeTemplateValue(template, field, value)
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
