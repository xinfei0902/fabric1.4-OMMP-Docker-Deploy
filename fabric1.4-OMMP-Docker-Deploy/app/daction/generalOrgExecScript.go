package daction

import (
	"deploy-server/app/dconfig"
	"deploy-server/app/dmysql"
	"deploy-server/app/objectdefine"
	"fmt"
)

//MakeStepGeneralOperateOrgExecScript 执行脚本命令
func MakeStepGeneralOperateOrgExecScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec add org Operate script start"},
		})

		err := GeneralCreateOperateOrgExecScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateOperateOrgExecScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec add org Operate script end"},
		})
		return nil
	}
}

//GeneralCreateOperateOrgExecScript 以后多组织扩展
func GeneralCreateOperateOrgExecScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		err := GeneralCreateOperateOrgExecScriptFile(general, org, output)
		if err != nil {
			// output.AppendLog(&objectdefine.StepHistory{
			// 	Name:  "generalCreateOperateOrgExecScript",
			// 	Error: []string{err.Error()},
			// })
			return err
		}
	}
	return nil
}

//GeneralCreateOperateOrgExecScriptFile 执行脚本
func GeneralCreateOperateOrgExecScriptFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("addOrg-%s", orgOrder.Name)
	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", OperateAddOrgIP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateOperateOrgExecScript",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	command := "start"
	buff := fmt.Sprintf("cd %s && chmod +x addOrgStep1.sh && ./addOrgStep1.sh", folder)
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

//MakeStepGeneralDeleteOperateOrgExecScript 执行删除组织更新配置块脚本命令
func MakeStepGeneralDeleteOperateOrgExecScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec delete org script start"},
		})

		err := GeneralDeleteOperateOrgExecScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDeleteOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeleteOperateOrgExecScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec delete org script end"},
		})
		return nil
	}
}

//GeneralDeleteOperateOrgExecScript 以后多组织扩展
func GeneralDeleteOperateOrgExecScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		err := GeneralDeleteOperateOrgExecScriptFile(general, org, output)
		if err != nil {
			// output.AppendLog(&objectdefine.StepHistory{
			// 	Name:  "generalCreateOperateOrgExecScript",
			// 	Error: []string{err.Error()},
			// })
			return err
		}
	}
	return nil
}

//GeneralDeleteOperateOrgExecScriptFile 执行脚本
func GeneralDeleteOperateOrgExecScriptFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("deleteOrg-%s", orgOrder.Name)
	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", OperateAddOrgIP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalDeleteOperateOrgExecScript",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	command := "start"
	buff := fmt.Sprintf("cd %s && chmod +x deleteOrgStep1.sh && ./deleteOrgStep1.sh", folder)
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

//MakeStepGeneralOrgExecScript 执行脚本命令
func MakeStepGeneralOrgExecScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec add org script start"},
		})

		err := GeneralCreateOrgExecScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateOrgExecScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec add org script end"},
		})
		return nil
	}
}

//GeneralCreateOrgExecScript 以后多组织扩展
func GeneralCreateOrgExecScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreateOrgExecScriptFile(general, org, peer, output)
			if err != nil {
				// output.AppendLog(&objectdefine.StepHistory{
				// 	Name:  "generalCreateOrgExecScript",
				// 	Error: []string{err.Error()},
				// })
				return err
			}
		}
	}
	return nil
}

//GeneralCreateOrgExecScriptFile 执行脚本
func GeneralCreateOrgExecScriptFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("addOrg-%s", orgOrder.Name)
	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateOrgExecScript",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	command := "start"
	buff := fmt.Sprintf("cd %s && chmod +x addOrgStep3.sh && ./addOrgStep3.sh", folder)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		// output.AppendLog(&objectdefine.StepHistory{
		// 	Name:  "generalCreateOrgExecScript",
		// 	Error: []string{err.Error()},
		// })
		return err
	}
	return nil
}

//MakeStepGeneralDeleteOrgExecScript 执行删除组织下删除节点脚本命令
func MakeStepGeneralDeleteOrgExecScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec delete org script start"},
		})

		err := GeneralDeleteOrgExecScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDeleteOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeleteOrgExecScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec delete org script end"},
		})
		return nil
	}
}

//GeneralDeleteOrgExecScript 以后多组织扩展
func GeneralDeleteOrgExecScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	peerAllIP := make(map[string]string, 0)
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			if _, ok := peerAllIP[peer.IP]; !ok {
				peerAllIP[peer.IP] = peer.Name
				err := GeneralDeleteOrgExecScriptFile(general, org, peer, output)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

//GeneralDeleteOrgExecScriptFile 执行脚本
func GeneralDeleteOrgExecScriptFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("deleteOrg-%s", orgOrder.Name)
	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalDeleteOrgExecScript",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	command := "start"
	buff := fmt.Sprintf("cd %s && chmod +x deleteOrgStep3.sh && ./deleteOrgStep3.sh", folder)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		return err
	}
	return nil
}

//MakeStepGeneralOrgExecPeerStatusScript 执行脚本命令
func MakeStepGeneralOrgExecPeerStatusScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec peerStatusMonitor script start"},
		})

		err := GeneralCreateOrgExecPeerStatusScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateOrgExecPeerStatusScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec peerStatusMonitor script end"},
		})
		return nil
	}
}

//GeneralCreateOrgExecPeerStatusScript 以后多组织扩展
func GeneralCreateOrgExecPeerStatusScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreateOrgExecPeerStatusScriptFile(general, org, peer, output)
			if err != nil {
				// output.AppendLog(&objectdefine.StepHistory{
				// 	Name:  "generalCreateOrgExecPeerStatusScript",
				// 	Error: []string{fmt.Sprintf("build exec script cmd error:%s", err.Error())},
				// })
				return err
			}
		}
	}
	return nil
}

//GeneralCreateOrgExecPeerStatusScriptFile 执行脚本
func GeneralCreateOrgExecPeerStatusScriptFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("addOrg-%s", orgOrder.Name)
	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateOrgExecPeerStatusScript",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	statusFileName := fmt.Sprintf("peerStatusMonitor-%s.sh", peerOrder.Domain)
	command := "status"
	buff := fmt.Sprintf("cd %s/%s && chmod +x %s && source /etc/profile && nohup ./%s >/dev/null 2>&1 &", folder, folder, statusFileName, statusFileName)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		// output.AppendLog(&objectdefine.StepHistory{
		// 	Name:  "generalCreateOrgExecPeerStatusScript",
		// 	Error: []string{fmt.Sprintf("exec http remote cmd error:%s", err.Error())},
		// })
		return err
	}
	return nil
}

//MakeStepGeneralDeleteOrgExecPeerStatusScript 执行脚本命令
func MakeStepGeneralDeleteOrgExecPeerStatusScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"delete peerStatusMonitor script start"},
		})

		err := GeneralDeleteOrgExecPeerStatusScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDeleteOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeleteOrgExecPeerStatusScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"delete peerStatusMonitor script end"},
		})
		return nil
	}
}

//GeneralDeleteOrgExecPeerStatusScript 以后多组织扩展
func GeneralDeleteOrgExecPeerStatusScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralDeleteOrgExecPeerStatusScriptFile(general, org, peer, output)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//GeneralDeleteOrgExecPeerStatusScriptFile 执行脚本
func GeneralDeleteOrgExecPeerStatusScriptFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("deleteOrg-%s", orgOrder.Name)
	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalDeleteOrgExecPeerStatusScript",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	statusFileName := fmt.Sprintf("deletePeerStatusMonitor-%s.sh", peerOrder.Domain)
	command := "status"
	buff := fmt.Sprintf("cd %s/%s && chmod +x %s && source /etc/profile &&  ./%s ", folder, folder, statusFileName, statusFileName)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		return err
	}
	return nil
}
