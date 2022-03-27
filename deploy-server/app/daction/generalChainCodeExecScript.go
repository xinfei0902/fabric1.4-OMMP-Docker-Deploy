package daction

import (
	"deploy-server/app/dconfig"
	"deploy-server/app/dmysql"
	"deploy-server/app/objectdefine"
	"fmt"
)

//####################新增合约#################################

//MakeStepGeneralChainCodeExecScript 执行脚本命令
func MakeStepGeneralChainCodeExecScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec add chaincode script start"},
		})

		err := GeneralCreateChainCodeExecScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChainCodeExecScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec add chaincode script end"},
		})
		return nil
	}
}

//GeneralCreateChainCodeExecScript 以后多链码扩展
func GeneralCreateChainCodeExecScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralCreateChainCodeExecScriptFile(general, &cc, ccName, output)
			if err != nil {
				// output.AppendLog(&objectdefine.StepHistory{
				// 	Name:  "generalCreateChainCodeExecScript",
				// 	Error: []string{err.Error()},
				// })
				return err
			}
		}
	}
	return nil
}

//GeneralCreateChainCodeExecScriptFile 执行脚本
func GeneralCreateChainCodeExecScriptFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("addChainCode-%s-%s", ccName, cc.Version)
	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", SelectExecCCIP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateChainCodeExecScript",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	command := "start"
	buff := fmt.Sprintf("cd %s && chmod +x addChainCodeStep1.sh && ./addChainCodeStep1.sh", folder)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		// output.AppendLog(&objectdefine.StepHistory{
		// 	Name:  "generalCreateChainCodeExecScript",
		// 	Error: []string{err.Error()},
		// })

		return err
	}

	return nil
}

//####################删除合约#################################

//MakeStepGeneralDeleteChainCodeExecScript 合约删除执行脚本命令
func MakeStepGeneralDeleteChainCodeExecScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec delete chaincode script start"},
		})

		err := GeneralDeleteChainCodeExecScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec delete chaincode script end"},
		})
		return nil
	}
}

//GeneralDeleteChainCodeExecScript 以后多链码扩展
func GeneralDeleteChainCodeExecScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralDeleteChainCodeExecScriptFile(general, &cc, ccName, output)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//GeneralDeleteChainCodeExecScriptFile 执行脚本
func GeneralDeleteChainCodeExecScriptFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {

	for _, ipaddress := range ExecCCIPList {
		folder := fmt.Sprintf("deleteChainCode-%s-%s", ccName, cc.Version)
		connectPort := dconfig.GetStringByKey("toolsPort")
		url := fmt.Sprintf("http://%s:%s/command/exec", ipaddress, connectPort)
		output.AppendLog(&objectdefine.StepHistory{
			Name: "generalDeleteChainCodeExecScript",
			Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
		})
		command := "start"
		buff := fmt.Sprintf("cd ./%s/%s && chmod +x deleteChainCodeScript.sh && ./deleteChainCodeScript.sh", folder, folder)
		args := []string{buff}
		err := MakeHTTPRemoteCmd(url, command, args, output)
		if err != nil {
			return err
		}
	}

	return nil
}

//####################升级合约#################################

//MakeStepGeneralUpgradeChainCodeExecScript 执行脚本命令
func MakeStepGeneralUpgradeChainCodeExecScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec upgrade chaincode script start"},
		})

		err := GeneralUpgradeChainCodeExecScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailUpgradeChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalUpgradeChainCodeExecScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec upgrade chaincode script end"},
		})
		return nil
	}
}

//GeneralUpgradeChainCodeExecScript 以后多链码扩展
func GeneralUpgradeChainCodeExecScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralUpgradeChainCodeExecScriptFile(general, &cc, ccName, output)
			if err != nil {
				// output.AppendLog(&objectdefine.StepHistory{
				// 	Name:  "generalCreateChainCodeScript",
				// 	Error: []string{err.Error()},
				// })
				return err
			}
		}
	}
	return nil
}

//GeneralUpgradeChainCodeExecScriptFile 执行脚本
func GeneralUpgradeChainCodeExecScriptFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("upgradeChainCode-%s-%s", ccName, cc.Version)
	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", SelectExecCCIP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalUpgradeChainCodeExecScript",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	command := "start"
	buff := fmt.Sprintf("cd %s && chmod +x upgradeChainCodeStep1.sh && ./upgradeChainCodeStep1.sh", folder)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		// output.AppendLog(&objectdefine.StepHistory{
		// 	Name:  "generalUpgradeChainCodeExecScript",
		// 	Error: []string{err.Error()},
		// })
		return err
	}
	return nil
}

//####################停用合约#################################

//MakeStepGeneralDisableChainCodeExecScript 合约停用执行脚本命令
func MakeStepGeneralDisableChainCodeExecScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec disable chaincode script start"},
		})

		err := GeneralDisableChainCodeExecScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDisableChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDisableChainCodeExecScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec disable chaincode script end"},
		})
		return nil
	}
}

//GeneralDisableChainCodeExecScript 以后多链码扩展
func GeneralDisableChainCodeExecScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralDisableChainCodeExecScriptFile(general, &cc, ccName, output)
			if err != nil {
				// output.AppendLog(&objectdefine.StepHistory{
				// 	Name:  "generalCreateDisableChainCodeScript",
				// 	Error: []string{err.Error()},
				// })
				return err
			}
		}
	}
	return nil
}

//GeneralDisableChainCodeExecScriptFile 执行脚本
func GeneralDisableChainCodeExecScriptFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {

	for _, ipaddress := range ExecCCIPList {
		folder := fmt.Sprintf("disableChainCode-%s-%s", ccName, cc.Version)
		connectPort := dconfig.GetStringByKey("toolsPort")
		url := fmt.Sprintf("http://%s:%s/command/exec", ipaddress, connectPort)
		output.AppendLog(&objectdefine.StepHistory{
			Name: "generalDisableChainCodeExecScript",
			Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
		})
		command := "start"
		buff := fmt.Sprintf("cd ./%s/%s && chmod +x disableChainCodeScript.sh && ./disableChainCodeScript.sh", folder, folder)
		args := []string{buff}
		err := MakeHTTPRemoteCmd(url, command, args, output)
		if err != nil {
			// output.AppendLog(&objectdefine.StepHistory{
			// 	Name:  "generalDisableChainCodeExecScript",
			// 	Error: []string{err.Error()},
			// })
			return err
		}
	}

	return nil
}

//####################启用合约#################################

//MakeStepGeneralEnableChainCodeExecScript 合约启用执行脚本命令
func MakeStepGeneralEnableChainCodeExecScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec enable chaincode script start"},
		})

		err := GeneralEnableChainCodeExecScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailEnableChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalEnableChainCodeExecScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec enable chaincode script end"},
		})
		return nil
	}
}

//GeneralEnableChainCodeExecScript 以后多链码扩展
func GeneralEnableChainCodeExecScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralEnableChainCodeExecScriptFile(general, &cc, ccName, output)
			if err != nil {
				// output.AppendLog(&objectdefine.StepHistory{
				// 	Name:  "generalCreateEnableChainCodeScript",
				// 	Error: []string{err.Error()},
				// })
				return err
			}
		}
	}
	return nil
}

//GeneralEnableChainCodeExecScriptFile 执行脚本
func GeneralEnableChainCodeExecScriptFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {

	for _, ipaddress := range ExecCCIPList {
		folder := fmt.Sprintf("enableChainCode-%s-%s", ccName, cc.Version)
		connectPort := dconfig.GetStringByKey("toolsPort")
		url := fmt.Sprintf("http://%s:%s/command/exec", ipaddress, connectPort)
		output.AppendLog(&objectdefine.StepHistory{
			Name: "generalEnableChainCodeExecScript",
			Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
		})
		command := "start"
		buff := fmt.Sprintf("cd ./%s/%s && chmod +x enableChainCodeScript.sh && ./enableChainCodeScript.sh", folder, folder)
		args := []string{buff}
		err := MakeHTTPRemoteCmd(url, command, args, output)
		if err != nil {
			// output.AppendLog(&objectdefine.StepHistory{
			// 	Name:  "generalEnableChainCodeExecScript",
			// 	Error: []string{err.Error()},
			// })
			return err
		}
	}

	return nil
}
