package vrealize

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
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
		return nil, fmt.Errorf("catalog lookup failed: %v", catalogErr)
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

	//Update template field values with user configuration
	resourceConfiguration := terraformRequest.ResourceConfiguration
	for configKey, configValue := range terraformRequest.ResourceConfiguration {
		for _, dataKey := range keyList {
			templateData := requestTemplate.Data[dataKey].(map[string]interface{})
			if field, found := checkKey(dataKey, configKey); found {
				if changeTemplateValue(templateData, field, configValue) {
					delete(terraformRequest.ResourceConfiguration, configKey)
				} else {
					log.Printf("%s was not replaced", configKey)
				}
			}
		}
	}

	log.Println("Entering Add Loop")
	for configKey, configVal := range resourceConfiguration {
		for _, dataKey := range keyList {
			log.Printf("Add Loop: configKey=[%s] dataKey=[%s]", configKey, dataKey)
			if field, found := checkKey(dataKey, configKey); found {
				log.Printf("Add Loop Contains %s for %s", dataKey, configKey)
				resourceItem := requestTemplate.Data[dataKey].(map[string]interface{})
				resourceItem = addTemplateValue(
					resourceItem["data"].(map[string]interface{}),
					field,
					configVal)
			}
		}
	}

	//update template with deployment level config
	// limit to description and reasons as other things could get us into trouble
	for depField, depValue := range terraformRequest.DeploymentConfiguration {
		fieldstr := fmt.Sprintf("%s", depField)
		switch fieldstr {
		case "description":
			requestTemplate.Description = depValue.(string)
		case "reasons":
			requestTemplate.Reasons = depValue.(string)
		default:
			log.Printf("unknown option [%s] with value [%s] ignoring\n", depField, depValue)
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
			template := val.(map[string]interface{})
			replaced = changeTemplateValue(template, field, value)
			templateInterface[key] = template
		} else if key == field {
			templateInterface[key] = value
			return true
		}
	}
	return replaced
}

func addTemplateValue(templateInterface map[string]interface{}, field string, value interface{}) map[string]interface{} {
	//simplest case is adding a simple value. Leaving as a func in case there's a need to do more complicated additions later
	//	templateInterface[data]
	for key,val := range templateInterface {
		// Recurse into any sub-structures named 'data'
		if reflect.ValueOf(val).Kind() == reflect.Map && key == "data" {
			template := val.(map[string]interface{})
			val = addTemplateValue(template, field, value)
		} else {
			templateInterface[field] = value
		}
	}
	return templateInterface
}
