package test

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

const TEST_RESOURCE_FOLDER = "testdata/"

func Content(t testing.TB, resourceName string) []byte {
	f, err := ioutil.ReadFile(TEST_RESOURCE_FOLDER + resourceName)
	if nil != err {
		t.Errorf("Error while loading "+resourceName+": %+v", err)
	}
	return f
}

func JsonContent(t testing.TB, resourceName string, v interface{}) {
	content := Content(t, resourceName)
	err := json.Unmarshal(content, v)
	if nil != err {
		t.Errorf("Error when parsing json from "+resourceName+": %+v", err)
	}
}
