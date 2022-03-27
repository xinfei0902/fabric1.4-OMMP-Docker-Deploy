package daction

import (
	"deploy-server/app/dconfig"
	
	"deploy-server/app/objectdefine"
	"fmt"
)

//MakeStepGeneralCreateCompleteDeployExecScriptStep 执行脚本命令
func MakeStepGeneralCreateCompleteDeployExecScriptStep(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec create channel script start"},
		})

		err := GeneralCreateCreateCompleteDeployExecScriptStep1(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			// result := err
			// err := dmysql.UpdateFailCreateChannelTaskStatus(general)
			// if err != nil {
			// 	output.AppendLog(&objectdefine.StepHistory{
			// 		Name:  "generalCreateChannelExecScriptStep",
			// 		Error: []string{err.Error()},
			// 	})
			// 	return err
			// }
			// return result
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec create channel script end"},
		})
		return nil
	}
}

//GeneralCreateCreateCompleteDeployExecScriptStep1 执行脚本
func GeneralCreateCreateCompleteDeployExecScriptStep1(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", CliIP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "GeneralCreateCreateCompleteDeployExecScriptStep1",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	command := "start"
	buff := fmt.Sprintf("cd %s && chmod +x deploy.sh && ./deploy.sh", "deploy")
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		// output.AppendLog(&objectdefine.StepHistory{
		// 	Name:  "generalCreateOperateOrgExecScript",
		// 	Error: []string{err.Error()},
		// })
		return err
	}
	return nil
}