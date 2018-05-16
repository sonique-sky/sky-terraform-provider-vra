package vrealize

import (
	"testing"
	"github.com/stretchr/testify/assert"
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