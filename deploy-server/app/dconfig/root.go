package dconfig

import (
	"deploy-server/app/derrors"
	"deploy-server/app/tools"
	"strings"
)

func init() {
	initValues()
}

//SetFileNameKey 为了获取本地config.json 先给全局globalFileKey赋值 可以根据key获取config.json里面内容
func SetFileNameKey(name, flag, key string, defaultValue string) (err error) {
	globalFileKey = key
	return Register(name, flag, key, defaultValue, "Config file name")
}

//Register 主要把 key value值写进globalContainer Map里
func Register(name, flag, key string, defaultValue interface{}, usage string) (err error) {
	key = tools.STDString(key)
	globalFlags, ok := globalFlagsWhole[name]
	if !ok {
		globalFlags = make(map[string]flagsPair)
	}
	if len(flag) > 0 {
		one, ok := globalFlags[key]
		if ok {
			if one.Name == key && one.Short == flag && one.Value == defaultValue && one.Usage == usage {
				return nil
			}
			return derrors.ErrorSameKeyExistf(key)
		}
		globalFlags[key] = flagsPair{
			Name:  key,
			Short: flag,
			Value: defaultValue,
			Usage: usage,
		}
	}

	globalContainer, ok := globalContainerWhole[name]
	if !ok {
		globalContainer = make(map[string]valuePair)
	}
	one, _ := globalContainer[key]
	one.Default = defaultValue
	globalContainer[key] = one

	globalContainerWhole[name] = globalContainer

	return nil
}

//Get 从globalValues Map中获取值
func Get(key string) (value interface{}, ok bool) {
	//key = tools.STDString(key)
	//strings.TrimSpace 去空白
	key = strings.TrimSpace(key)
	value, ok = globalValues[key]
	return
}

//GetStringByKey 根据key获取value
func GetStringByKey(key string) string {
	v, ok := Get(key)
	if !ok {
		return ""
	}
	return GetString(v, "")
}
