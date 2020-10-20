package daction

import (
	"deploy-server/app/dconfig"
	"deploy-server/app/dmysql"
	"deploy-server/app/objectdefine"
	"fmt"
)

//MakeStepGeneralPeerExecScript 执行脚本命令
func MakeStepGeneralPeerExecScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec add peer script start"},
		})

		err := GeneralCreatePeerExecScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreatePeerTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreatePeerExecScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec add peer script end"},
		})
		return nil
	}
}

//GeneralCreatePeerExecScript 以后多组织扩展
func GeneralCreatePeerExecScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreatePeerExecScriptFile(general, org, peer, output)
			if err != nil {
				// output.AppendLog(&objectdefine.StepHistory{
				// 	Name:  "generalCreatePeerExecScript",
				// 	Error: []string{err.Error()},
				// })
				return err
			}
		}

	}
	return nil
}

//GeneralCreatePeerExecScriptFile 执行脚本
func GeneralCreatePeerExecScriptFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("addPeer-%s-%s", orgOrder.Name, peerOrder.Name)

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreatePeerExecScript",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	command := "start"
	buff := fmt.Sprintf("cd %s && chmod +x addPeerStep1.sh && ./addPeerStep1.sh", folder)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		// output.AppendLog(&objectdefine.StepHistory{
		// 	Name:  "generalCreatePeerExecScript",
		// 	Error: []string{err.Error()},
		// })
		return err
	}
	return nil
}

//####################节点状态脚本#############################

//MakeStepGeneralPeerExecPeerStatusScript 执行脚本命令
func MakeStepGeneralPeerExecPeerStatusScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec peerStatusMonitor script start"},
		})

		err := GeneralCreatePeerExecPeerStatusScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreatePeerTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreatePeerExecPeerStatusScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"execpeerStatusMonitor script end"},
		})
		return nil
	}
}

//GeneralCreatePeerExecPeerStatusScript 以后多组织扩展
func GeneralCreatePeerExecPeerStatusScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreatePeerExecPeerStatusScriptFile(general, org, peer, output)
			if err != nil {
				// output.AppendLog(&objectdefine.StepHistory{
				// 	Name:  "generalCreatePeerExecPeerStatusScript",
				// 	Error: []string{err.Error()},
				// })
				return err
			}
		}

	}
	return nil
}

//GeneralCreatePeerExecPeerStatusScriptFile 执行脚本
func GeneralCreatePeerExecPeerStatusScriptFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("addPeer-%s-%s", orgOrder.Name, peerOrder.Name)

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreatePeerExecPeerStatusScript",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	statusFileName := fmt.Sprintf("peerStatusMonitor-%s.sh", peerOrder.Domain)
	command := "status"
	buff := fmt.Sprintf("cd %s/%s && chmod +x %s && source /etc/profile && nohup ./%s >/dev/null 2>&1 &", folder, folder, statusFileName, statusFileName)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		// output.AppendLog(&objectdefine.StepHistory{
		// 	Name:  "generalCreatePeerExecPeerStatusScript",
		// 	Error: []string{err.Error()},
		// })
		return err
	}
	return nil
}

//#######################删除节点执行脚本#################################

//MakeStepGeneralDeletePeerExecScript 执行脚本命令
func MakeStepGeneralDeletePeerExecScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec delete peer script start"},
		})

		err := GeneralDeletePeerExecScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDeletePeerTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeletePeerExecScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec delete peer script end"},
		})
		return nil
	}
}

//GeneralDeletePeerExecScript 以后多组织扩展
func GeneralDeletePeerExecScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralDeletePeerExecScriptFile(general, org, peer, output)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

//GeneralDeletePeerExecScriptFile 执行脚本
func GeneralDeletePeerExecScriptFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("deletePeer-%s-%s", orgOrder.Name, peerOrder.Name)

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalDeletePeerExecScript",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	command := "start"
	buff := fmt.Sprintf("cd %s && chmod +x deletePeerStep1.sh && ./deletePeerStep1.sh", folder)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		return err
	}
	return nil
}

//#########################删除节点状态执行脚本###########################

//MakeStepGeneralDeletePeerExecPeerStatusScript 执行脚本命令
func MakeStepGeneralDeletePeerExecPeerStatusScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec peerStatusMonitor script start"},
		})

		err := GeneralDeletePeerExecPeerStatusScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDeletePeerTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeletePeerExecPeerStatusScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"execpeerStatusMonitor script end"},
		})
		return nil
	}
}

//GeneralDeletePeerExecPeerStatusScript 以后多组织扩展
func GeneralDeletePeerExecPeerStatusScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralDeletePeerExecPeerStatusScriptFile(general, org, peer, output)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

//GeneralDeletePeerExecPeerStatusScriptFile 执行脚本
func GeneralDeletePeerExecPeerStatusScriptFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("deletePeer-%s-%s", orgOrder.Name, peerOrder.Name)

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalDeletePeerExecPeerStatusScript",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	statusFileName := fmt.Sprintf("peerStatusMonitor-%s.sh", peerOrder.Domain)
	command := "status"
	buff := fmt.Sprintf("cd %s/%s && chmod +x %s && source /etc/profile && ./%s", folder, folder, statusFileName, statusFileName)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		return err
	}
	return nil
}
