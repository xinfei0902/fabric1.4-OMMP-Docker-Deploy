package dconfig

import (
	"encoding/json"
	"io/ioutil"
)

//ReadConfigFile 读文件
func ReadConfigFile(name string) (ret map[string]interface{}, err error) {
	buff, err := ioutil.ReadFile(name)
	if err != nil {
		return
	}
	ret = make(map[string]interface{})
	err = json.Unmarshal(buff, &ret)
	return
}
