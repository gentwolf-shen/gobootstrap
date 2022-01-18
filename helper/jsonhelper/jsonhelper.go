package jsonhelper

import (
	"io/ioutil"

	json "github.com/json-iterator/go"
)

func ToObj(bytes []byte, obj interface{}) error {
	return json.Unmarshal(bytes, obj)
}

func ToByte(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func ToStr(obj interface{}) string {
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}

func FileToObj(filename string, obj interface{}) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return ToObj(bytes, obj)
}
