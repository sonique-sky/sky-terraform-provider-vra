package machine

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	"sort"
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
	inty := make(map[string]interface{})
	inty["foo"] = "baz"

	changeValue(inty, "foo", "bbbar")

	fmt.Println(inty["foo"])

}
func changeValue(templateInterface map[string]interface{}, field string, value interface{}) (map[string]interface{}, bool) {
	templateInterface[field] = value
	return templateInterface, true
}
