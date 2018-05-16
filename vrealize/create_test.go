package vrealize

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	"sort"
	"io/ioutil"
	"encoding/json"
	"github.com/sonique-sky/sky-terraform-provider-vra/vrealize/api"
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

func Test_checkKey_OK(t *testing.T) {
	configKey := "Machine.foo.bar.baz"
	dataKey := "Machine"

	res, found := checkKey(dataKey, configKey)

	assert.True(t, found)
	assert.Equal(t, "foo.bar.baz", res)
}

func Test_checkKey_Fail(t *testing.T) {
	configKey := "Machine.foo.bar.baz"
	dataKey := "NotMachine"

	_, found := checkKey(dataKey, configKey)

	assert.False(t, found)
}

func Test_keylistSorting(t *testing.T) {
	var keyList = []string{"first", "longestest", "lil"}
	//Arrange keys in descending order of text length

	sort.Slice(keyList, func(i1, i2 int) bool {
		return len(keyList[i1]) > len(keyList[i2])
	})

	fmt.Println(keyList)
}

func Test_foof(t *testing.T) {
	template := new(api.RequestTemplate)
	bytes, _ := ioutil.ReadFile("../api/test_data/catalog_item_request_template.json")
	json.Unmarshal(bytes, template)

	templateData := template.Data["CentOS_6.3"].(map[string]interface{})

	changeTemplateValue(templateData, "foo", "bbbar")

	fmt.Println(template)
}

