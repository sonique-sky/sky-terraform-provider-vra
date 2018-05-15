package machine

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize/api"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

func Test_addTemplateValue(t *testing.T) {


	template := make(map[string]interface{})
	template["foo"] = "bar"

	data := make(map[string]interface{})
	data["boo"] = "hoho"
	template["data"] = data

	addTemplateValue(template, "foo", "var")
	addTemplateValue(template, "boo", "hehe")


	assert.Equal(t, "var", template["foo"])

	assert.Equal(t, "hehe", template["data"].(map[string]interface{})["boo"])


}


func Test_keylistSorting(t *testing.T) {
	var requestTemplate = new(api.RequestTemplate)

	bytes, _ := ioutil.ReadFile("../api/test_data/request_template.json")

	json.Unmarshal(bytes, requestTemplate)

	var keyList []string
	for field := range requestTemplate.Data {
		if reflect.ValueOf(requestTemplate.Data[field]).Kind() == reflect.Map {
			keyList = append(keyList, field)
		}
	}

	//Arrange keys in descending order of text length
	for field1 := range keyList {
		for field2 := range keyList {
			if len(keyList[field1]) > len(keyList[field2]) {
				temp := keyList[field1]
				keyList[field1], keyList[field2] = keyList[field2], temp
			}
		}
	}

	fmt.Print(keyList)
}