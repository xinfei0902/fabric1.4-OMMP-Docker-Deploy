package dconfig

import (
	"deploy-server/app/derrors"
	"os"

	"github.com/spf13/cobra"
)

//AddFlags 添加
func AddFlags(name string, one *cobra.Command) (err error) {
	flags := one.Flags()

	globalContainer, ok := globalContainerWhole[name]
	if !ok {
		globalContainer = make(map[string]valuePair)
	}

	globalFlags, ok := globalFlagsWhole[name]
	if !ok {
		globalFlags = make(map[string]flagsPair)
	}

	for key, v := range globalFlags {
		if one, ok := globalContainer[key]; true == ok && one.Command != nil {
			continue
		}
		if v.Value == nil {
			p := flags.StringP(v.Name, v.Short, "", v.Usage)
			value, _ := globalContainer[key]
			value.Command = p
			globalContainer[key] = value
			continue
		}
		switch v.Value.(type) {
		case int, *int:
			p := flags.IntP(v.Name, v.Short, 0, v.Usage)
			value, _ := globalContainer[key]
			value.Command = p
			globalContainer[key] = value
		case bool, *bool:
			p := flags.BoolP(v.Name, v.Short, false, v.Usage)
			value, _ := globalContainer[key]
			value.Command = p
			globalContainer[key] = value
		case string, *string:
			p := flags.StringP(v.Name, v.Short, "", v.Usage)
			value, _ := globalContainer[key]
			value.Command = p
			globalContainer[key] = value
		default:
			panic("add code here")
		}
	}
	globalContainerWhole[name] = globalContainer

	return
}

//ReadFileAndHoldCommand 读取config.josn 把k-v值写入全局变量 globalContainer和globalValues
func ReadFileAndHoldCommand(name string, cmd *cobra.Command) (err error) {
	tmp := make(map[string]interface{})

	filename := ""
	//以service 为name 获取下面值
	globalContainer, ok := globalContainerWhole[name]
	if !ok {
		globalContainer = make(map[string]valuePair)
	}

	//在初始化 flages时 获取globalFileKey=config  globalContainer["config"]=config.json路径
	if len(globalFileKey) > 0 {
		v, ok := globalContainer[globalFileKey]
		if !ok {
			return derrors.ErrorKeyNotContainValuef(globalFileKey)
		}
		filename = v.GetStringValue()

		if len(filename) == 0 {
			return derrors.ErrorKeyNotContainValuef(globalFileKey)
		}
		tmp, err = ReadConfigFile(filename)
		if err != nil {
			if false == os.IsNotExist(err) {
				return err
			}
		}
	}
	for key, value := range tmp {
		one, _ := globalContainer[key]
		one.File = value
		globalContainer[key] = one
	}

	globalContainerWhole[name] = globalContainer

	for key, value := range globalContainer {
		globalValues[key] = value.GetInterface()
	}

	return nil
}
