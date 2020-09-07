package daction

import (
	"deploy-server/app/dconfig"
	"deploy-server/app/objectdefine"
	"fmt"
)

//MakeStepGeneralDisablePeerExecCommand 停用peer
func MakeStepGeneralDisablePeerExecCommand(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec disable peer command start"},
		})

		err := GeneralDisablePeerExec(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec disable peer command end"},
		})
		return nil
	}
}

//GeneralDisablePeerExec 多组织
func GeneralDisablePeerExec(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralDisablePeerExecCommand(general, org, peer, output)
			if err != nil {
				// output.AppendLog(&objectdefine.StepHistory{
				// 	Name:  "generalDisablePeerExecCommand",
				// 	Error: []string{err.Error()},
				// })
				return err
			}
		}
	}
	return nil
}

//GeneralDisablePeerExecCommand 执行命令
func GeneralDisablePeerExecCommand(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalDisablePeerExecCommand",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	command := "start"
	buff := fmt.Sprintf("docker stop %s", peerOrder.Domain)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		// output.AppendLog(&objectdefine.StepHistory{
		// 	Name:  "generalDisablePeerExecCommand",
		// 	Error: []string{err.Error()},
		// })
		return err
	}
	return nil
}

//MakeStepGeneralEnablePeerExecCommand 启用peer
func MakeStepGeneralEnablePeerExecCommand(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec enable peer command start"},
		})

		err := GeneralEnablePeerExec(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec enable peer command end"},
		})
		return nil
	}
}

//GeneralEnablePeerExec 多组织
func GeneralEnablePeerExec(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralEnablePeerExecCommand(general, org, peer, output)
			if err != nil {
				// output.AppendLog(&objectdefine.StepHistory{
				// 	Name:  "generalenablePeerExecCommand",
				// 	Error: []string{err.Error()},
				// })
				return err
			}
		}
	}
	return nil
}

//GeneralEnablePeerExecCommand 执行命令
func GeneralEnablePeerExecCommand(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalEnablePeerExecCommand",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	command := "start"
	buff := fmt.Sprintf("docker start %s", peerOrder.Domain)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		// output.AppendLog(&objectdefine.StepHistory{
		// 	Name:  "generalenablePeerExecCommand",
		// 	Error: []string{err.Error()},
		// })
		return err
	}
	return nil
}
