package daction

import (
	"deploy-server/app/dcache"
	"deploy-server/app/dmysql"
	"deploy-server/app/objectdefine"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
)

//OperateAddOrgIP 用来在Baiyun cli端所在执行操作ip地址 因为现在白云部署只有一个cli 为以后考虑可能是其他cli所准备
var OperateAddOrgIP string

//MakeStepGeneralOrgScriptStep 构建第一个脚本 构建需求分步执行
func MakeStepGeneralOrgScriptStep(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add org addOrgStep1.sh start"},
		})
		//防备没有替换成功 保证执行动作之前 字符串为空
		if len(OperateAddOrgIP) != 0 {
			OperateAddOrgIP = ""
		}
		err := GeneralCreateOrgScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateOrgScriptStep",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add org addOrgStep1.sh end"},
		})
		return nil
	}
}

//GeneralCreateOrgScript 后期实现多组织构建，目前虽然是遍历 但仅新增一个组织
func GeneralCreateOrgScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreateOrgScriptFile(general, org, peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateOrgScriptStep",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralCreateOrgScriptFile 创建脚本文件
func GeneralCreateOrgScriptFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folderName := fmt.Sprintf("addOrg-%s", orgOrder.Name)
	indent, err := dmysql.GetStartTaskBeforIndent(general)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgScriptStep",
			Error: []string{"get the latest indent error"},
		})
		return err
	}
	var oldOrgCliName string
	var signOrgList string
	var signOrgTLSList string
	var signOrgMSPList string
	var signOrgAddressList string
	isExistOrg := make(map[string]bool, len(general.Org))
	for _, org := range indent.Org {
		for _, peer := range org.Peer {

			if peer.User == "Admin" {
				if len(OperateAddOrgIP) == 0 && org.Name == "Baiyun" {
					oldOrgCliName = peer.CliName
					OperateAddOrgIP = peer.IP
				}
				orgMSPName := fmt.Sprintf("%sMSP", org.Name)
				if _, ok := isExistOrg[orgMSPName]; !ok {
					if len(signOrgList) == 0 {
						signOrgList = orgMSPName
						isExistOrg[orgMSPName] = true
					} else {
						signOrgList = signOrgList + " " + orgMSPName
						isExistOrg[orgMSPName] = true
					}
				}
				//固定shell数组格式
				orgTLSPath := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", org.OrgDomain, "peers", fmt.Sprintf("%s.%s", peer.Name, org.OrgDomain), "tls", "ca.crt"))
				if len(signOrgTLSList) == 0 {
					signOrgTLSList = orgTLSPath
				} else {
					signOrgTLSList = signOrgTLSList + " " + orgTLSPath
				}
				orgMspPath := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", org.OrgDomain, "users", fmt.Sprintf("Admin@%s", org.OrgDomain), "msp"))
				if len(signOrgMSPList) == 0 {
					signOrgMSPList = orgMspPath
				} else {
					signOrgMSPList = signOrgMSPList + " " + orgMspPath
				}
				orgAddress := fmt.Sprintf("%s:%d", peer.Domain, peer.Port)
				if len(signOrgAddressList) == 0 {
					signOrgAddressList = orgAddress
				} else {
					signOrgAddressList = signOrgAddressList + " " + orgAddress
				}
			}

		}

	}

	orgOrgDomain := orgOrder.OrgDomain
	peerPort := peerOrder.Port
	peerDomain := peerOrder.Domain
	//构建map 填写需要的组变量
	replaceMap := make(map[string]string)
	folder := fmt.Sprintf("addOrg-%s", orgOrder.Name)
	caFolder := fmt.Sprintf("ca-%s", orgOrder.Name)
	replaceMap["target"] = folder
	replaceMap["cliName"] = oldOrgCliName
	replaceMap["newOrgCliName"] = peerOrder.CliName
	replaceMap["caContinueName"] = caFolder
	replaceMap["peerName"] = peerOrder.Domain
	//脚本2 需要变量 主要为更新块配置
	exchangeMap := make(map[string]string)
	exchangeMap["target"] = folder
	exchangeMap["orgMSPID"] = fmt.Sprintf("%sMSP", orgOrder.Name)
	exchangeMap["channelName"] = general.ChannelName
	var ordererIP string
	var ordererPort int
	var ordererOrgDomain string
	var ordererDomain string
	for _, orderer := range indent.Orderer {
		ordererIP = orderer.IP
		ordererPort = orderer.Port
		ordererDomain = orderer.Domain
		ordererOrgDomain = orderer.OrgDomain
		break
	}
	ordererAddress := fmt.Sprintf("%s:%d", ordererDomain, ordererPort)
	exchangeMap["ordererAddress"] = ordererAddress

	ordererTLSCa := filepath.ToSlash(filepath.Join(ordererOrgDomain, "orderers", ordererDomain, "tls", "ca.crt"))
	exchangeMap["ordererTlsCa"] = ordererTLSCa
	ordererUsersMsp := filepath.ToSlash(filepath.Join(ordererOrgDomain, "users", fmt.Sprintf("Admin@%s", ordererOrgDomain), "msp"))
	exchangeMap["ordererUsersMsp"] = ordererUsersMsp
	exchangeMap["signOrgMSPID"] = fmt.Sprintf("(%s)", signOrgList)
	exchangeMap["signOrgTlsCa"] = fmt.Sprintf("(%s)", signOrgTLSList)
	exchangeMap["signOrgUsersMsp"] = fmt.Sprintf("(%s)", signOrgMSPList)
	exchangeMap["signOrgPeerHosts"] = fmt.Sprintf("(%s)", signOrgAddressList)
	//构建脚本3 主要启动组织节点
	//使用脚本1变量即可
	//构建脚本4 主要加入通道 安装链码
	change := make(map[string]string)
	change["channelName"] = general.ChannelName
	//change["ccName"] = general.ChannelName
	change["ordererAddress"] = ordererAddress
	ordererHosts := fmt.Sprintf("\"%s %s\"", ordererIP, ordererDomain)
	change["ordererHosts"] = ordererHosts
	change["ordererTlsCa"] = filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations", ordererTLSCa))
	change["orgMSPID"] = fmt.Sprintf("%sMSP", orgOrder.Name)
	peerTLSCa := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", orgOrgDomain, "peers", peerDomain, "tls", "ca.crt"))
	change["peerTlsCa"] = peerTLSCa
	peerUsersMsp := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", orgOrgDomain, "users", fmt.Sprintf("Admin@%s", orgOrgDomain), "msp"))
	change["peerUsersMsp"] = peerUsersMsp
	peerHosts := fmt.Sprintf("%s:%d", peerDomain, peerPort)
	change["peerHosts"] = peerHosts

	//创建三个脚本文件
	step1Buff, err := dcache.GetOrgScriptStep1Template(general.Version, replaceMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgScriptStep",
			Error: []string{"Build addOrgStep1.sh replace parame error"},
		})
		err = errors.WithMessage(err, "Build addOrgStep1.sh replace parame error")
		return err
	}
	scriptStep1SavePath := filepath.Join(general.BaseOutput, "addOrg", OperateAddOrgIP, folderName)
	_, err = os.Stat(scriptStep1SavePath)
	if err != nil {
		err := os.MkdirAll(scriptStep1SavePath, 0777)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateOrgScriptStep",
				Error: []string{"Build output folder errorr"},
			})
			return errors.WithMessage(err, "Build output folder error")
		}
	}

	err = ioutil.WriteFile(filepath.Join(scriptStep1SavePath, "addOrgStep1.sh"), step1Buff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgScriptStep",
			Error: []string{"replace parame writer addOrgStep1.sh error"},
		})
		err = errors.WithMessage(err, "replace parame writer addOrgStep1.sh error")
		return err
	}
	scriptStep2SavePath := filepath.Join(general.BaseOutput, "addOrg", OperateAddOrgIP, folderName, folderName)
	_, err = os.Stat(scriptStep2SavePath)
	if err != nil {
		err := os.MkdirAll(scriptStep2SavePath, 0777)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateOrgScriptStep",
				Error: []string{"Build output folder errorr"},
			})
			return errors.WithMessage(err, "Build output folder error")
		}
	}
	step2Buff, err := dcache.GetOrgScriptStep2Template(general.Version, exchangeMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgScriptStep",
			Error: []string{"Build addOrgStep2.sh exchange parame error"},
		})
		err = errors.WithMessage(err, "Build addOrgStep2.sh exchange parame error")
		return err
	}
	err = ioutil.WriteFile(filepath.Join(scriptStep2SavePath, "addOrgStep2.sh"), step2Buff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgScriptStep",
			Error: []string{"exchange parame writer addOrgStep2.sh error"},
		})
		err = errors.WithMessage(err, "exchange parame writer addOrgStep2.sh error")
		return err
	}

	//// 这里做区分 前两个脚本是用来在已有的cli更新新增组织块配置， 为了考虑不同机器 后两个脚本用来启动和加入通道 都是再一台机器上操作
	scriptStep3SavePath := filepath.Join(general.BaseOutput, "addOrg", peerOrder.IP, folderName)
	_, err = os.Stat(scriptStep3SavePath)
	if err != nil {
		err := os.MkdirAll(scriptStep3SavePath, 0777)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateOrgScriptStep",
				Error: []string{"Build output folder errorr"},
			})
			return errors.WithMessage(err, "Build output folder error")
		}
	}
	step3Buff, err := dcache.GetOrgScriptStep3Template(general.Version, replaceMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgScriptStep",
			Error: []string{"Build addOrgStep3.sh change  parame error"},
		})
		err = errors.WithMessage(err, "Build addOrgStep3.sh change  parame error")
		return err
	}
	err = ioutil.WriteFile(filepath.Join(scriptStep3SavePath, "addOrgStep3.sh"), step3Buff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgScriptStep",
			Error: []string{"change parame writer addOrgStep3.sh error"},
		})
		err = errors.WithMessage(err, "change parame writer addOrgStep3.sh error")
		return err
	}
	scriptStep4SavePath := filepath.Join(general.BaseOutput, "addOrg", peerOrder.IP, folderName, folderName)
	_, err = os.Stat(scriptStep4SavePath)
	if err != nil {
		err := os.MkdirAll(scriptStep4SavePath, 0777)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateOrgScriptStep",
				Error: []string{"Build output folder errorr"},
			})
			return errors.WithMessage(err, "Build output folder error")
		}
	}
	step4Buff, err := dcache.GetOrgScriptStep4Template(general.Version, change)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgScriptStep",
			Error: []string{"Build addOrgStep4.sh change  parame error"},
		})
		err = errors.WithMessage(err, "Build addOrgStep4.sh change  parame error")
		return err
	}
	err = ioutil.WriteFile(filepath.Join(scriptStep4SavePath, "addOrgStep4.sh"), step4Buff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgScriptStep",
			Error: []string{"change parame writer addOrgStep4.sh error"},
		})
		err = errors.WithMessage(err, "change parame writer addOrgStep4.sh error")
		return err
	}
	return nil
}

//MakeStepGeneralDeleteOrgScriptStep 构建脚本 构建需求分步执行
func MakeStepGeneralDeleteOrgScriptStep(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete org script start"},
		})
		//防备没有替换成功 保证执行动作之前 字符串为空
		if len(OperateAddOrgIP) != 0 {
			OperateAddOrgIP = ""
		}
		err := GeneralDeleteOrgScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDeleteOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeleteOrgScriptStep",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete org script end"},
		})
		return nil
	}
}

//GeneralDeleteOrgScript 后期实现多组织构建，目前虽然是遍历 但仅新增一个组织
func GeneralDeleteOrgScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		//for _, peer := range org.Peer {
		err := GeneralDeleteOrgScriptFile(general, org, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalDeleteOrgScriptStep",
				Error: []string{err.Error()},
			})
			return err
		}
		//}
	}
	return nil
}

//GeneralDeleteOrgScriptFile 创建脚本文件
func GeneralDeleteOrgScriptFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, output *objectdefine.TaskNode) error {
	folderName := fmt.Sprintf("deleteOrg-%s", orgOrder.Name)
	indent, err := dmysql.GetDeleteTaskBeforIndent(general)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalDeleteOrgScriptStep",
			Error: []string{"get the latest indent error"},
		})
		return err
	}
	var oldOrgCliName string
	var signOrgList string
	var signOrgTLSList string
	var signOrgMSPList string
	var signOrgAddressList string
	isExistOrg := make(map[string]bool, len(general.Org))
	for _, org := range indent.Org {
		for _, peer := range org.Peer {
			if peer.User == "Admin" {
				if len(OperateAddOrgIP) == 0 {
					oldOrgCliName = peer.CliName
					OperateAddOrgIP = peer.IP
				}
				orgMSPName := fmt.Sprintf("%sMSP", org.Name)
				if _, ok := isExistOrg[orgMSPName]; !ok {
					if len(signOrgList) == 0 {
						signOrgList = orgMSPName
						isExistOrg[orgMSPName] = true
					} else {
						signOrgList = signOrgList + " " + orgMSPName
						isExistOrg[orgMSPName] = true
					}
				}
				//固定shell数组格式
				orgTLSPath := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", org.OrgDomain, "peers", fmt.Sprintf("%s.%s", peer.Name, org.OrgDomain), "tls", "ca.crt"))
				if len(signOrgTLSList) == 0 {
					signOrgTLSList = orgTLSPath
				} else {
					signOrgTLSList = signOrgTLSList + " " + orgTLSPath
				}
				orgMspPath := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", org.OrgDomain, "users", fmt.Sprintf("Admin@%s", org.OrgDomain), "msp"))
				if len(signOrgMSPList) == 0 {
					signOrgMSPList = orgMspPath
				} else {
					signOrgMSPList = signOrgMSPList + " " + orgMspPath
				}
				orgAddress := fmt.Sprintf("%s:%d", peer.Domain, peer.Port)
				if len(signOrgAddressList) == 0 {
					signOrgAddressList = orgAddress
				} else {
					signOrgAddressList = signOrgAddressList + " " + orgAddress
				}
			}

		}

	}

	//构建map 填写需要的组变量
	replaceMap := make(map[string]string)
	folder := fmt.Sprintf("deleteOrg-%s", orgOrder.Name)
	//caFolder := fmt.Sprintf("ca-%s", orgOrder.Name)
	replaceMap["target"] = folder
	replaceMap["cliName"] = oldOrgCliName

	//脚本2 需要变量 主要为更新块配置
	exchangeMap := make(map[string]string)
	exchangeMap["target"] = folder
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

	ordererTLSCa := filepath.ToSlash(filepath.Join(ordererOrgDomain, "orderers", ordererDomain, "tls", "ca.crt"))
	exchangeMap["ordererTlsCa"] = ordererTLSCa
	ordererUsersMsp := filepath.ToSlash(filepath.Join(ordererOrgDomain, "users", fmt.Sprintf("Admin@%s", ordererOrgDomain), "msp"))
	exchangeMap["ordererUsersMsp"] = ordererUsersMsp
	exchangeMap["signOrgMSPID"] = fmt.Sprintf("(%s)", signOrgList)
	exchangeMap["signOrgTlsCa"] = fmt.Sprintf("(%s)", signOrgTLSList)
	exchangeMap["signOrgUsersMsp"] = fmt.Sprintf("(%s)", signOrgMSPList)
	exchangeMap["signOrgPeerHosts"] = fmt.Sprintf("(%s)", signOrgAddressList)
	//先生产两个脚本

	step1Buff, err := dcache.GetDeleteOrgScriptStep1Template(general.Version, replaceMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalDeleteOrgScriptStep",
			Error: []string{"Build deleteOrgStep1.sh replace parame error"},
		})
		err = errors.WithMessage(err, "Build deleteOrgStep1.sh replace parame error")
		return err
	}
	scriptStep1SavePath := filepath.Join(general.BaseOutput, "deleteOrg", OperateAddOrgIP, folderName)
	_, err = os.Stat(scriptStep1SavePath)
	if err != nil {
		err := os.MkdirAll(scriptStep1SavePath, 0777)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalDeleteOrgScriptStep",
				Error: []string{"Build output folder errorr"},
			})
			return errors.WithMessage(err, "Build output folder error")
		}
	}

	err = ioutil.WriteFile(filepath.Join(scriptStep1SavePath, "deleteOrgStep1.sh"), step1Buff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalDeleteOrgScriptStep",
			Error: []string{"replace parame writer deleteOrgStep1.sh error"},
		})
		err = errors.WithMessage(err, "replace parame writer deleteOrgStep1.sh error")
		return err
	}
	scriptStep2SavePath := filepath.Join(general.BaseOutput, "deleteOrg", OperateAddOrgIP, folderName, folderName)
	_, err = os.Stat(scriptStep2SavePath)
	if err != nil {
		err := os.MkdirAll(scriptStep2SavePath, 0777)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalDeleteOrgScriptStep",
				Error: []string{"Build output folder errorr"},
			})
			return errors.WithMessage(err, "Build output folder error")
		}
	}
	step2Buff, err := dcache.GetDeleteOrgScriptStep2Template(general.Version, exchangeMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalDeleteOrgScriptStep",
			Error: []string{"Build deleteOrgStep2.sh exchange parame error"},
		})
		err = errors.WithMessage(err, "Build deleteOrgStep2.sh exchange parame error")
		return err
	}
	err = ioutil.WriteFile(filepath.Join(scriptStep2SavePath, "deleteOrgStep2.sh"), step2Buff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalDeleteOrgScriptStep",
			Error: []string{"exchange parame writer deleteOrgStep2.sh error"},
		})
		err = errors.WithMessage(err, "exchange parame writer deleteOrgStep2.sh error")
		return err
	}
	//构建脚本3 主要停掉组织节点
	//var peerContainerArray string
	peerIPPairContainer := make(map[string]string, 0)
	peerIPPairCCImages := make(map[string]string)
	peerContainerCA := orgOrder.CA.Name
	peerIPPairContainer[orgOrder.CA.IP] = peerContainerCA

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
				Name:  "generalDeleteOrgScriptStep",
				Error: []string{"Build deletePeerStatusMonitor.sh replace parame error"},
			})
			err = errors.WithMessage(err, "Build deletePeerStatusMonitor.sh replace parame error")
			return err
		}
		outputRoot := filepath.Join(general.BaseOutput, "deleteOrg", peer.IP, folderName, folderName)
		statusFileName := fmt.Sprintf("deletePeerStatusMonitor-%s.sh", peer.Domain)
		err = ioutil.WriteFile(filepath.Join(outputRoot, statusFileName), peerStatusBuff, 0644)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalDeleteOrgScriptStep",
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
		scriptStep3SavePath := filepath.Join(general.BaseOutput, "deleteOrg", ip, folderName)
		_, err = os.Stat(scriptStep3SavePath)
		if err != nil {
			err := os.MkdirAll(scriptStep3SavePath, 0777)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeleteOrgScriptStep",
					Error: []string{"Build output folder errorr"},
				})
				return errors.WithMessage(err, "Build output folder error")
			}
		}
		step3Buff, err := dcache.GetDeleteOrgScriptStep3Template(general.Version, change)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalDeleteOrgScriptStep",
				Error: []string{"Build deleteOrgStep3.sh change  parame error"},
			})
			err = errors.WithMessage(err, "Build deleteOrgStep3.sh change  parame error")
			return err
		}
		err = ioutil.WriteFile(filepath.Join(scriptStep3SavePath, "deleteOrgStep3.sh"), step3Buff, 0644)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalDeleteOrgScriptStep",
				Error: []string{"change parame writer deleteOrgStep3.sh error"},
			})
			err = errors.WithMessage(err, "change parame writer deleteOrgStep3.sh error")
			return err
		}
	}

	return nil
}
