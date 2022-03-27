package daction

import (
	"deploy-server/app/dcache"
	"deploy-server/app/dmysql"
	"deploy-server/app/objectdefine"
	"deploy-server/app/tools"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
)

//MakeStepGeneralPeerScriptStep 构建第一个脚本 构建需求分步执行
func MakeStepGeneralPeerScriptStep(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add peer script start"},
		})

		err := GeneralCreatePeerScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreatePeerTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreatePeerScriptStep",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add peer script end"},
		})
		return nil
	}
}

//GeneralCreatePeerScript 后期实现多节点构建，目前虽然是遍历 但仅新增一个节点
func GeneralCreatePeerScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreatePeerScriptFile(general, org, peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreatePeerScriptStep",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralCreatePeerScriptFile 创建脚本文件
func GeneralCreatePeerScriptFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folderName := fmt.Sprintf("addPeer-%s-%s", orgOrder.Name, peerOrder.Name)
	outputRoot := filepath.Join(general.BaseOutput, "addPeer", peerOrder.IP, folderName, folderName)
	err := os.MkdirAll(outputRoot, 0777)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerScriptStep",
			Error: []string{"Build output folder errorr"},
		})
		return errors.WithMessage(err, "Build output folder error")
	}
	//拷贝证书
	src := filepath.Join(general.BaseOutput, "crypto-config")
	dst := filepath.Join(outputRoot, "crypto-config")
	tools.CopyFolder(dst, src)
	//获取此任务之前最新订单信息
	indent, err := dmysql.GetStartTaskBeforIndent(general)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerScriptStep",
			Error: []string{"get the latest indent error"},
		})
		return err
	}
	//var oldOrgCliName string
	var installCCNameList string
	var installCCVersion string
	var installCCPath string
	var isExecInstallCC string
	peerHostsList := ""
	for _,orderer := range indent.Orderer{
		if len(peerHostsList)==0{
			ordererIntIP,err := dmysql.GetIntIPFromExtIP(orderer.IP)
			if err != nil{
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChainCodeScript",
					Error: []string{err.Error()},
				})
				return err
			}
			peerHostsList = fmt.Sprintf("\"%s %s\"", ordererIntIP, orderer.Domain)
		}else{
			ordererIntIP,err := dmysql.GetIntIPFromExtIP(orderer.IP)
			if err != nil{
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChainCodeScript",
					Error: []string{err.Error()},
				})
				return err
			}
			peerHostsList = peerHostsList + " " + fmt.Sprintf("\"%s %s\"", ordererIntIP, orderer.Domain)
		}
	}
    //已有网络 hosts
	for _, org := range indent.Org {
		for _, peer := range org.Peer {
			if len(peerHostsList)==0{
				peerIntIP,err := dmysql.GetIntIPFromExtIP(peer.IP)
				if err != nil{
					output.AppendLog(&objectdefine.StepHistory{
						Name:  "generalCreateChainCodeScript",
						Error: []string{err.Error()},
					})
					return err
				}
				peerHostsList = fmt.Sprintf("\"%s %s\"", peerIntIP, peer.Domain)
			}else{
				peerIntIP,err := dmysql.GetIntIPFromExtIP(peer.IP)
				if err != nil{
					output.AppendLog(&objectdefine.StepHistory{
						Name:  "generalCreateChainCodeScript",
						Error: []string{err.Error()},
					})
					return err
				}
				peerHostsList = peerHostsList + " " + fmt.Sprintf("\"%s %s\"", peerIntIP, peer.Domain)
			}
		}
	}

	//添加要添加hosts

	if len(peerHostsList)==0{
		peerIntIP,err := dmysql.GetIntIPFromExtIP(peerOrder.IP)
		if err != nil{
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateChainCodeScript",
				Error: []string{err.Error()},
			})
			return err
		}
		peerHostsList = fmt.Sprintf("\"%s %s\"", peerIntIP, peerOrder.Domain)
	}else{
		peerIntIP,err := dmysql.GetIntIPFromExtIP(peerOrder.IP)
		if err != nil{
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateChainCodeScript",
				Error: []string{err.Error()},
			})
			return err
		}
		peerHostsList = peerHostsList + " " + fmt.Sprintf("\"%s %s\"", peerIntIP, peerOrder.Domain)
	}
	
	for _, org := range indent.Org {
		if org.Name == orgOrder.Name {
			if len(org.ChainCode) == 0 {
				isExecInstallCC = "0"
			} else {
				//#################################################
				//此处如果需要安装链码 那么需要拷贝链码到chaincode文件
				//################################################
				chaincodePath := filepath.Join(outputRoot, "chaincode")
				err := os.MkdirAll(chaincodePath, 0777)
				if err != nil {
					output.AppendLog(&objectdefine.StepHistory{
						Name:  "generalCreatePeerScriptStep",
						Error: []string{"Build output folder errorr"},
					})
					return errors.WithMessage(err, "Build chaincode folder error")
				}
				for ccName, cc := range org.ChainCode {
					if len(installCCNameList) == 0 {
						installCCNameList = ccName
						srcChaincodePath := filepath.Join(dcache.GetTemplatePathByVersion(general.Version), "template", "chaincode")
						src := filepath.Join(srcChaincodePath, ccName)
						tools.CopyFolder(chaincodePath, src)
					} else {
						installCCNameList = installCCNameList + " " + ccName
						srcChaincodePath := filepath.Join(dcache.GetTemplatePathByVersion(general.Version), "template", "chaincode")
						src := filepath.Join(srcChaincodePath, ccName)
						tools.CopyFolder(chaincodePath, src)
					}
					if len(installCCVersion) == 0 {
						installCCVersion = cc.Version
					} else {
						installCCVersion = installCCVersion + " " + cc.Version
					}
					ccPath := fmt.Sprintf("github.com/chaincode/%s", ccName)
					if len(installCCPath) == 0 {
						installCCPath = ccPath
					} else {
						installCCPath = installCCPath + " " + ccPath
					}
					isExecInstallCC = "1"
				}
			}
		}

	}

	orgOrgDomain := orgOrder.OrgDomain
	peerPort := peerOrder.Port
	peerDomain := peerOrder.Domain
	//构建map 填写需要的组变量
	replaceMap := make(map[string]string)
	replaceMap["target"] = folderName
	replaceMap["newOrgCliName"] = peerOrder.CliName
	replaceMap["peerName"] = peerOrder.Domain
	replaceMap["isExecInstallCC"] = isExecInstallCC
	replaceMap["peerHosts"] = fmt.Sprintf("(%s)", peerHostsList)
	//脚本2 需要变量 主要为加入通道
	exchangeMap := make(map[string]string)
	exchangeMap["target"] = folderName
	exchangeMap["orgMSPID"] = fmt.Sprintf("%sMSP", orgOrder.Name)
	exchangeMap["channelName"] = general.ChannelName
	//var ordererIP string
	var ordererPort int
	var ordererOrgDomain string
	var ordererDomain string
	for _, orderer := range indent.Orderer {
		//ordererIP = orderer.IP
		ordererPort = orderer.Port
		ordererDomain = orderer.Domain
		ordererOrgDomain = orderer.OrgDomain
		break
	}
	ordererAddress := fmt.Sprintf("%s:%d", ordererDomain, ordererPort)
	exchangeMap["ordererAddress"] = ordererAddress
	ordererTLSCa := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations", ordererOrgDomain, "orderers", ordererDomain, "tls", "ca.crt"))
	exchangeMap["ordererTlsCa"] = ordererTLSCa
	exchangeMap["orgMSPID"] = fmt.Sprintf("%sMSP", orgOrder.Name)
	peerTLSCa := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", orgOrgDomain, "peers", peerDomain, "tls", "ca.crt"))
	exchangeMap["peerTlsCa"] = peerTLSCa
	peerUsersMsp := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", orgOrgDomain, "users", fmt.Sprintf("Admin@%s", orgOrgDomain), "msp"))
	exchangeMap["peerUsersMsp"] = peerUsersMsp
	peerAddress := fmt.Sprintf("%s:%d", peerDomain, peerPort)
	exchangeMap["peerAddress"] = peerAddress
	// ordererIntIP,err := dmysql.GetIntIPFromExtIP(ordererIP)
	// if err != nil{
	// 	output.AppendLog(&objectdefine.StepHistory{
	// 		Name:  "generalCreatePeerScriptStep",
	// 		Error: []string{err.Error()},
	// 	})
	// 	return err
	// }
	// peerHostList := fmt.Sprintf("\"%s %s\"", ordererIntIP, ordererDomain)
	// peerIntIP,err := dmysql.GetIntIPFromExtIP(peerOrder.IP)
	// if err != nil{
	// 	output.AppendLog(&objectdefine.StepHistory{
	// 		Name:  "generalCreatePeerScriptStep",
	// 		Error: []string{err.Error()},
	// 	})
	// 	return err
	// }
	// peerHostList = peerHostList + " " + fmt.Sprintf("\"%s %s\"", peerIntIP, peerOrder.Domain)
	exchangeMap["peerHosts"] = fmt.Sprintf("(%s)", peerHostsList)
	//创建三个脚本文件
	step1Buff, err := dcache.GetPeerScriptStep1Template(general.Version, replaceMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerScriptStep",
			Error: []string{"Build addPeerStep1.sh replace parame error"},
		})
		err = errors.WithMessage(err, "Build addPeerStep1.sh replace parame error")
		return err
	}
	scriptStep1SavePath := filepath.Join(general.BaseOutput, "addPeer", peerOrder.IP, folderName)
	err = ioutil.WriteFile(filepath.Join(scriptStep1SavePath, "addPeerStep1.sh"), step1Buff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerScriptStep",
			Error: []string{"replace parame writer addPeerStep1.sh error"},
		})
		err = errors.WithMessage(err, "replace parame writer addPeerStep1.sh error")
		return err
	}
	step2Buff, err := dcache.GetPeerScriptStep2Template(general.Version, exchangeMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerScriptStep",
			Error: []string{"Build addPeerStep2.sh exchange parame error"},
		})
		err = errors.WithMessage(err, "Build addPeerStep2.sh exchange parame error")
		return err
	}
	err = ioutil.WriteFile(filepath.Join(outputRoot, "addPeerStep2.sh"), step2Buff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerScriptStep",
			Error: []string{"exchange parame writer addPeerStep2.sh error"},
		})
		err = errors.WithMessage(err, "exchange parame writer addOrgStep2.sh error")
		return err
	}
	if isExecInstallCC == "1" {
		//构建脚本3 主要安装链码
		change := make(map[string]string)
		change["target"] = folderName
		change["channelName"] = general.ChannelName
		change["orgMSPID"] = fmt.Sprintf("%sMSP", orgOrder.Name)
		change["peerTlsCa"] = peerTLSCa
		change["peerUsersMsp"] = peerUsersMsp
		change["peerAddress"] = peerAddress
		change["installCCName"] = fmt.Sprintf("(%s)", installCCNameList)
		change["installCCVersion"] = fmt.Sprintf("(%s)", installCCVersion)
		change["installCCPath"] = fmt.Sprintf("(%s)", installCCPath)
		step3Buff, err := dcache.GetPeerScriptStep3Template(general.Version, change)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreatePeerScriptStep",
				Error: []string{"Build addPeerStep3.sh change  parame error"},
			})
			err = errors.WithMessage(err, "Build addPeerStep3.sh change  parame error")
			return err
		}
		err = ioutil.WriteFile(filepath.Join(outputRoot, "addPeerStep3.sh"), step3Buff, 0644)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreatePeerScriptStep",
				Error: []string{"change parame writer addPeerStep3.sh error"},
			})
			err = errors.WithMessage(err, "change parame writer addPeerStep3.sh error")
			return err
		}
	}
	return nil
}

//MakeStepGeneralDeletePeerScriptStep 构建脚本 构建需求分步执行
func MakeStepGeneralDeletePeerScriptStep(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"delete peer script start"},
		})
		//防备没有替换成功 保证执行动作之前 字符串为空
		if len(OperateAddOrgIP) != 0 {
			OperateAddOrgIP = ""
		}
		err := GeneralDeletePeerScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDeletePeerTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeletePeerScriptStep",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"cdelete peer script end"},
		})
		return nil
	}
}

//GeneralDeletePeerScript 后期实现多组织构建，目前虽然是遍历 但仅新增一个组织
func GeneralDeletePeerScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralDeletePeerScriptFile(general, org, peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeletePeerScriptStep",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralDeletePeerScriptFile 创建脚本文件
func GeneralDeletePeerScriptFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folderName := fmt.Sprintf("deletePeer-%s-%s", orgOrder.Name, peerOrder.Name)

	//构建脚本 主要停掉节点
	//var peerContainerArray string
	peerIPPairContainer := make(map[string]string, 0)
	peerIPPairCCImages := make(map[string]string)

	for _, peer := range orgOrder.Peer {
		if str, ok := peerIPPairContainer[peer.IP]; ok {
			str = str + " " + peer.Domain + " " + fmt.Sprintf("couchdb-%s", peer.Domain) + " " + peer.CliName
			peerIPPairContainer[peer.IP] = str
			images := peerIPPairCCImages[peer.IP]
			images = images + " " + fmt.Sprintf("dev-%s", peer.Domain)
			peerIPPairCCImages[peer.IP] = images
		} else {
			peerIPPairContainer[peer.IP] = peer.Domain + " " + fmt.Sprintf("couchdb-%s", peer.Domain) + " " + peer.CliName
			peerIPPairCCImages[peer.IP] = fmt.Sprintf("dev-%s", peer.Domain)
		}
		//peerStatusMonitor脚本
		exChange := make(map[string]string)
		exChange["checkPeerPort"] = strconv.Itoa(peer.Port)
		exChange["peerNick"] = peer.NickName
		exChange["accessKey"] = peer.AccessKey
		exChange["peerID"] = strconv.Itoa(peer.PeerID)
		exChange["peerIP"] = peer.IP
		exChange["peerPort"] = strconv.Itoa(peer.Port)
		peerStatusBuff, err := dcache.GetDeletePeerStatusFileTemplate(general.Version, exChange)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalDeletePeerScriptStep",
				Error: []string{"Build deletePeerStatusMonitor.sh replace parame error"},
			})
			err = errors.WithMessage(err, "Build deletePeerStatusMonitor.sh replace parame error")
			return err
		}
		outputRoot := filepath.Join(general.BaseOutput, "deletePeer", peer.IP, folderName, folderName)
		statusFileName := fmt.Sprintf("deletePeerStatusMonitor-%s.sh", peer.Domain)
		err = ioutil.WriteFile(filepath.Join(outputRoot, statusFileName), peerStatusBuff, 0644)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalDeletePeerScriptStep",
				Error: []string{"replace parame writer deletePeerStatusMonitor.sh error"},
			})
			err = errors.WithMessage(err, "replace parame writer deletePeerStatusMonitor.sh error")
			return err
		}
	}
	for ip, value := range peerIPPairContainer {
		change := make(map[string]string)
		change["peerContainerArray"] = fmt.Sprintf("(%s)", value)
		ccImages := peerIPPairCCImages[ip]
		change["deleteCCImages"] = fmt.Sprintf("(%s)", ccImages)
		scriptStep1SavePath := filepath.Join(general.BaseOutput, "deletePeer", ip, folderName)
		_, err := os.Stat(scriptStep1SavePath)
		if err != nil {
			err := os.MkdirAll(scriptStep1SavePath, 0777)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeletePeerScriptStep",
					Error: []string{"Build output folder errorr"},
				})
				return errors.WithMessage(err, "Build output folder error")
			}
		}
		step1Buff, err := dcache.GetDeletePeerScriptStep1Template(general.Version, change)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalDeletePeerScriptStep",
				Error: []string{"Build deletePeerStep1.sh change  parame error"},
			})
			err = errors.WithMessage(err, "Build deletePeerStep1.sh change  parame error")
			return err
		}
		err = ioutil.WriteFile(filepath.Join(scriptStep1SavePath, "deletePeerStep1.sh"), step1Buff, 0644)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalDeletePeerScriptStep",
				Error: []string{"change parame writer deletePeerStep1.sh error"},
			})
			err = errors.WithMessage(err, "change parame writer deletePeerStep1.sh error")
			return err
		}
	}

	return nil
}
