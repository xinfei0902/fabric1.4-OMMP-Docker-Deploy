package daction

import (
	"deploy-server/app/dconfig"
	"deploy-server/app/dmysql"
	"deploy-server/app/objectdefine"
	"fmt"
)

//MakeStepGeneralChannelExecScriptStep1 执行脚本命令
func MakeStepGeneralChannelExecScriptStep1(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec create channel script start"},
		})

		err := GeneralCreateChannelExecScriptStep1(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateChannelTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChannelExecScriptStep",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec create channel script end"},
		})
		return nil
	}
}

//GeneralCreateChannelExecScriptStep1 执行脚本
func GeneralCreateChannelExecScriptStep1(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreateChannelExecScriptStep1File(general, &peer, output)
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

//GeneralCreateChannelExecScriptStep1File 执行脚本
//这个所有操作的执行操作都需要改动 考虑多台机器，最好是一个步骤都执行完毕之后 再执行下一个步骤 而不是一台机器全部执行完毕 再去执行别的机器
func GeneralCreateChannelExecScriptStep1File(general *objectdefine.Indent, peerOrder *objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("createChannel-%s", general.ChannelName)
	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateChannelExecScriptStep",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	command := "start"
	buff := fmt.Sprintf("cd %s && chmod +x createChannelStep1.sh && ./createChannelStep1.sh", folder)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		return err
	}

	return nil
}

//MakeStepGeneralChannelExecScriptStep2 执行脚本命令
func MakeStepGeneralChannelExecScriptStep2(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec create channel script step2 start"},
		})

		err := GeneralCreateChannelExecScriptStep2(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec create channel script step2 end"},
		})
		return nil
	}
}

//GeneralCreateChannelExecScriptStep2 以后多组织扩展
func GeneralCreateChannelExecScriptStep2(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ipaddress, folderA := range ChannelPeerIPList {
		err := GeneralCreateChannelExecScriptStep2File(general, ipaddress, folderA, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateChannelExecScriptStep2",
				Error: []string{err.Error()},
			})
			return err
		}
	}
	return nil
}

//GeneralCreateChannelExecScriptStep2File 执行脚本
//这个素有操作的执行操作都需要改动 考虑多台机器，最好是一个步骤都执行完毕之后 再执行下一个步骤 而不是一台机器全部执行完毕 再去执行别的机器
func GeneralCreateChannelExecScriptStep2File(general *objectdefine.Indent, ipaddress string, folderA []string, output *objectdefine.TaskNode) error {
	//folder := fmt.Sprintf("createChannel-%s-%s-%s", general.ChannelName, orgOrder.Name, peerOrder.Name)
	isEventExec := false
	for _, folder := range folderA {
		if general.IsNewOrgCreateChannel == true {
			// url := fmt.Sprintf("%s/command/exec", one.Connect)
			if ipaddress == ChannelPeerIP && !isEventExec {
				isEventExec = true
				connectPort := dconfig.GetStringByKey("toolsPort")
				url := fmt.Sprintf("http://%s:%s/command/exec", ipaddress, connectPort)
				output.AppendLog(&objectdefine.StepHistory{
					Name: "generalCreateChannelExecScriptStep2",
					Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
				})
				command := "start"
				buff := fmt.Sprintf("cd %s && chmod +x createNewOrgChannelStep2.sh && ./createNewOrgChannelStep2.sh", folder)
				args := []string{buff}
				err := MakeHTTPRemoteCmd(url, command, args, output)
				if err != nil {
					output.AppendLog(&objectdefine.StepHistory{
						Name:  "generalCreateChannelExecScriptStep2",
						Error: []string{err.Error()},
					})
					return err
				}
			}

		} else {

			connectPort := dconfig.GetStringByKey("toolsPort")
			url := fmt.Sprintf("http://%s:%s/command/exec", ipaddress, connectPort)
			output.AppendLog(&objectdefine.StepHistory{
				Name: "generalCreateChannelExecScriptStep2",
				Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
			})
			command := "start"
			buff := fmt.Sprintf("cd %s && chmod +x createChannelStep2.sh && ./createChannelStep2.sh", folder)
			args := []string{buff}
			err := MakeHTTPRemoteCmd(url, command, args, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChannelExecScriptStep2",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}
