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

//SelectExecCCIP 选择一个Admin权限的节点cli去操作更新配置块等操作
var SelectExecCCIP string

//ExecCCIPList 用来存放所有部署过合约的远端机器IP
var ExecCCIPList []string

//###################################创建新增合约脚本#####################################

//MakeStepGeneralCreateChainCodeScript 创建脚本
func MakeStepGeneralCreateChainCodeScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add chainCode script start"},
		})

		err := GeneralCreateChainCodeScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChainCodeScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add chainCode script end"},
		})
		return nil
	}
}

//GeneralCreateChainCodeScript 为一次多安装合约
func GeneralCreateChainCodeScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralCreateChainCodeScriptFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChainCodeScript",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralCreateChainCodeScriptFile  创建脚本
func GeneralCreateChainCodeScriptFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	folderName := fmt.Sprintf("addChainCode-%s-%s", ccName, cc.Version)
	//获取此任务之前最新订单信息
	indent, err := dmysql.GetStartTaskBeforIndent(general)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChainCodeScript",
			Error: []string{"get the latest indent error"},
		})
		return err
	}

	//组织数组
	var endorseOrg = make(map[string]objectdefine.OrgType, len(cc.EndorsementOrg))
	for _, orgName := range cc.EndorsementOrg {
		endorseOrg[orgName] = indent.Org[orgName]
	}

	var orgName string
	var peerIP string
	var cliName string
	var orgNameMSP string
	var orgPeerNUM string
	var peerTLSCaList string
	var peerUsersMspList string
	var peerAddressList string
	var peerHostsList string
	for _, org := range endorseOrg {
		if len(orgName) == 0 {
			orgName = org.Name
		}
		if len(orgNameMSP) == 0 {
			orgNameMSP = fmt.Sprintf("%sMSP", org.Name)
		} else {
			orgNameMSP = orgNameMSP + " " + fmt.Sprintf("%sMSP", org.Name)
		}
		if len(orgPeerNUM) == 0 {
			orgPeerNUM = strconv.Itoa(len(org.Peer))
		} else {
			orgPeerNUM = orgPeerNUM + " " + strconv.Itoa(len(org.Peer))
		}
		//找到一个主节点 用来对合约的操作
		for _, peer := range org.Peer {
			if len(peerIP) == 0 {
				if peer.User == "Admin" {
					peerIP = peer.IP
					SelectExecCCIP = peer.IP
					cliName = peer.CliName
					peerTLSCaList = filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", org.OrgDomain, "peers", peer.Domain, "tls", "ca.crt"))
					peerUsersMspList = filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", org.OrgDomain, "users", fmt.Sprintf("Admin@%s", org.OrgDomain), "msp"))
					peerAddressList = fmt.Sprintf("%s:%d", peer.Domain, peer.Port)
					peerHostsList = fmt.Sprintf("\"%s %s\"", peer.IP, peer.Domain)
				}
			}
		}
	}

	outputRoot := filepath.Join(general.BaseOutput, "addChainCode", peerIP, folderName, folderName)
	err = os.MkdirAll(outputRoot, 0777)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChainCodeScript",
			Error: []string{"Build output folder errorr"},
		})
		return errors.WithMessage(err, "Build output folder error")
	}
	//拷贝证书
	src := filepath.Join(general.SourceBaseOutput, "crypto-config")
	dst := filepath.Join(outputRoot, "crypto-config")
	_, err = os.Stat(dst)
	if err != nil {
		err = os.MkdirAll(dst, 0777)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateChainCodeScript",
				Error: []string{"Build output folder errorr"},
			})
			return errors.WithMessage(err, "Build output folder error")
		}
	}
	tools.CopyFolder(dst, src)
	//安装链码节点环境变量组装
	for _, org := range endorseOrg {
		for _, peer := range org.Peer {
			if cliName != peer.CliName {
				if len(peerTLSCaList) == 0 || len(peerUsersMspList) == 0 || len(peerAddressList) == 0 {
					return errors.WithMessage(err, "org noe include admin peer ")
				}
				peerTLSCa := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", org.OrgDomain, "peers", peer.Domain, "tls", "ca.crt"))
				peerTLSCaList = peerTLSCaList + " " + peerTLSCa
				peerUsersMsp := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", org.OrgDomain, "users", fmt.Sprintf("Admin@%s", org.OrgDomain), "msp"))
				peerUsersMspList = peerUsersMspList + " " + peerUsersMsp
				peerAddress := fmt.Sprintf("%s:%d", peer.Domain, peer.Port)
				peerAddressList = peerAddressList + " " + peerAddress
				peerHostsList = peerHostsList + " " + fmt.Sprintf("\"%s %s\"", peer.IP, peer.Domain)
			}
		}
	}
	//脚本配置
	//cliName := fmt.Sprintf("cli-%s", orgName)
	replaceMap := make(map[string]string)
	replaceMap["target"] = folderName
	replaceMap["cliName"] = cliName
	replaceMap["installCCName"] = ccName
	//脚本step
	exchangeMap := make(map[string]string)
	exchangeMap["channelName"] = general.ChannelName
	var ordererPort int
	var ordererOrgDomain string
	var ordererDomain string
	for _, orderer := range indent.Orderer {
		//ordererName = orderer.Name
		ordererPort = orderer.Port
		ordererDomain = orderer.Domain
		ordererOrgDomain = orderer.OrgDomain
		break
	}
	ordererAddress := fmt.Sprintf("%s:%d", ordererDomain, ordererPort)
	exchangeMap["ordererAddress"] = ordererAddress
	ordererTLSCa := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations", ordererOrgDomain, "orderers", ordererDomain, "tls", "ca.crt"))
	exchangeMap["ordererTlsCa"] = ordererTLSCa
	//链码配置
	exchangeMap["installCCName"] = ccName
	exchangeMap["installCCVersion"] = cc.Version
	exchangeMap["installCCPath"] = fmt.Sprintf("github.com/chaincode/%s", ccName)
	exchangeMap["installCCPolicy"] = cc.Policy
	//判断是否第一次实例化
	whether, err := dmysql.GetCCInstantiatedTime(general.ChannelName, ccName, cc.Version)
	if err != nil || whether == -1 {
		return errors.WithMessage(err, "mysql get chaincode Time fail")
	}
	exchangeMap["isFristInstantiated"] = strconv.Itoa(whether)
	//节点配置
	exchangeMap["orgMSPID"] = fmt.Sprintf("(%s)", orgNameMSP)
	exchangeMap["orgPeerNum"] = fmt.Sprintf("(%s)", orgPeerNUM)
	exchangeMap["peerTlsCa"] = fmt.Sprintf("(%s)", peerTLSCaList)
	exchangeMap["peerUsersMsp"] = fmt.Sprintf("(%s)", peerUsersMspList)
	exchangeMap["peerAddress"] = fmt.Sprintf("(%s)", peerAddressList)
	exchangeMap["peerHosts"] = fmt.Sprintf("(%s)", peerHostsList)
	//创建脚本文件
	step1Buff, err := dcache.GetChainCodeScriptStep1Template(general.Version, replaceMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChainCodeScript",
			Error: []string{"Build addChainCodeStep1.sh replace parame error"},
		})
		err = errors.WithMessage(err, "Build addChainCodeStep1.sh replace parame error")
		return err
	}
	scriptStep1SavePath := filepath.Join(general.BaseOutput, "addChainCode", peerIP, folderName)
	err = ioutil.WriteFile(filepath.Join(scriptStep1SavePath, "addChainCodeStep1.sh"), step1Buff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChainCodeScript",
			Error: []string{"replace parame writer addChainCodeStep1.sh error"},
		})
		err = errors.WithMessage(err, "replace parame writer addChainCodeStep1.sh error")
		return err
	}
	step2Buff, err := dcache.GetChainCodeScriptStep2Template(general.Version, exchangeMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChainCodeScript",
			Error: []string{"Build addPeerStep2.sh exchange parame error"},
		})
		err = errors.WithMessage(err, "Build addChainCodeStep2.sh exchange parame error")
		return err
	}
	err = ioutil.WriteFile(filepath.Join(outputRoot, "addChainCodeStep2.sh"), step2Buff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChainCodeScript",
			Error: []string{"exchange parame writer addChainCodeStep2.sh error"},
		})
		err = errors.WithMessage(err, "exchange parame writer addChainCodeStep2.sh error")
		return err
	}

	//拷贝合约
	dstCC := filepath.Join(outputRoot, "chaincode", ccName)
	_, err = os.Stat(dstCC)
	if err != nil {
		err = os.MkdirAll(dst, 0777)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateChainCodeScript",
				Error: []string{"Build output folder errorr"},
			})
			return errors.WithMessage(err, "Build output folder error")
		}
	}
	ccSavepath := filepath.ToSlash(filepath.Join(dcache.GetTemplatePathByVersion(general.Version), "chaincode", ccName))
	_, err = os.Stat(ccSavepath)
	if err != nil {
		return errors.WithMessage(err, "Build chaincode file is not exist")
	}
	tools.CopyFolder(dstCC, ccSavepath)
	return nil
}

//###################################创建升级合约脚本#####################################

//MakeStepGeneralUpgradeChainCodeScript 创建升级合约脚本
func MakeStepGeneralUpgradeChainCodeScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create upgrade chainCode script start"},
		})

		err := GeneralUpgradeChainCodeScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailUpgradeChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateUpgradeChainCodeScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create upgrade chainCode script end"},
		})
		return nil
	}
}

//GeneralUpgradeChainCodeScript 为一次多安装合约
func GeneralUpgradeChainCodeScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralUpgradeChainCodeScriptFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateUpgradeChainCodeScript",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralUpgradeChainCodeScriptFile  创建脚本
func GeneralUpgradeChainCodeScriptFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	folderName := fmt.Sprintf("upgradeChainCode-%s-%s", ccName, cc.Version)
	//获取此任务之前最新订单信息
	indent, err := dmysql.GetStartTaskBeforIndent(general)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateUpgradeChainCodeScript",
			Error: []string{"get the latest indent error"},
		})
		return err
	}
	//组织数组
	var endorseOrg = make(map[string]objectdefine.OrgType, len(cc.EndorsementOrg))
	for _, orgName := range cc.EndorsementOrg {
		endorseOrg[orgName] = indent.Org[orgName]
	}
	var orgName string
	var peerIP string
	var cliName string
	var orgNameMSP string
	var orgPeerNUM string
	var peerTLSCaList string
	var peerUsersMspList string
	var peerAddressList string
	var peerHostsList string

	for _, org := range endorseOrg {
		if len(orgName) == 0 {
			orgName = org.Name
		}
		if len(orgNameMSP) == 0 {
			orgNameMSP = fmt.Sprintf("%sMSP", org.Name)
		} else {
			orgNameMSP = orgNameMSP + " " + fmt.Sprintf("%sMSP", org.Name)
		}
		if len(orgPeerNUM) == 0 {
			orgPeerNUM = strconv.Itoa(len(org.Peer))
		} else {
			orgPeerNUM = orgPeerNUM + " " + strconv.Itoa(len(org.Peer))
		}
		for _, peer := range org.Peer {
			if len(peerIP) == 0 {
				if peer.User == "Admin" {
					peerIP = peer.IP
					SelectExecCCIP = peer.IP
					cliName = peer.CliName
					peerTLSCaList = filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", org.OrgDomain, "peers", peer.Domain, "tls", "ca.crt"))
					peerUsersMspList = filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", org.OrgDomain, "users", fmt.Sprintf("Admin@%s", org.OrgDomain), "msp"))
					peerAddressList = fmt.Sprintf("%s:%d", peer.Domain, peer.Port)
					peerHostsList = fmt.Sprintf("\"%s %s\"", peer.IP, peer.Domain)
				}
			}
		}
	}

	outputRoot := filepath.Join(general.BaseOutput, "upgradeChainCode", peerIP, folderName, folderName)
	err = os.MkdirAll(outputRoot, 0777)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateUpgradeChainCodeScript",
			Error: []string{"Build output folder errorr"},
		})
		return errors.WithMessage(err, "Build output folder error")
	}

	//安装链码节点环境变量组装
	for _, org := range endorseOrg {
		for _, peer := range org.Peer {
			if cliName != peer.CliName {
				if len(peerTLSCaList) == 0 || len(peerUsersMspList) == 0 || len(peerAddressList) == 0 {
					return errors.WithMessage(err, "org noe include admin peer ")
				}
				peerTLSCa := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", org.OrgDomain, "peers", peer.Domain, "tls", "ca.crt"))
				peerTLSCaList = peerTLSCaList + " " + peerTLSCa
				peerUsersMsp := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", org.OrgDomain, "users", fmt.Sprintf("Admin@%s", org.OrgDomain), "msp"))
				peerUsersMspList = peerUsersMspList + " " + peerUsersMsp
				peerAddress := fmt.Sprintf("%s:%d", peer.Domain, peer.Port)
				peerAddressList = peerAddressList + " " + peerAddress
				peerHostsList = peerHostsList + " " + fmt.Sprintf("\"%s %s\"", peer.IP, peer.Domain)
			}
		}
	}
	//拷贝证书
	src := filepath.Join(general.SourceBaseOutput, "crypto-config")
	dst := filepath.Join(outputRoot, "crypto-config")
	_, err = os.Stat(dst)
	if err != nil {
		err = os.MkdirAll(dst, 0777)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateUpgradeChainCodeScript",
				Error: []string{"Build output folder errorr"},
			})
			return errors.WithMessage(err, "Build output folder error")
		}
	}
	tools.CopyFolder(dst, src)
	//脚本配置
	//cliName := fmt.Sprintf("cli-%s", orgName)
	replaceMap := make(map[string]string)
	replaceMap["target"] = folderName
	replaceMap["cliName"] = cliName
	replaceMap["installCCName"] = fmt.Sprintf("%s-%s", ccName, cc.Version)
	//脚本step
	exchangeMap := make(map[string]string)
	exchangeMap["channelName"] = general.ChannelName
	var ordererPort int
	var ordererOrgDomain string
	var ordererDomain string
	for _, orderer := range indent.Orderer {
		//ordererName = orderer.Name
		ordererPort = orderer.Port
		ordererDomain = orderer.Domain
		ordererOrgDomain = orderer.OrgDomain
		break
	}
	ordererAddress := fmt.Sprintf("%s:%d", ordererDomain, ordererPort)
	exchangeMap["ordererAddress"] = ordererAddress
	ordererTLSCa := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations", ordererOrgDomain, "orderers", ordererDomain, "tls", "ca.crt"))
	exchangeMap["ordererTlsCa"] = ordererTLSCa
	//合约配置
	exchangeMap["installCCName"] = ccName
	exchangeMap["installCCVersion"] = cc.Version
	exchangeMap["installCCPath"] = fmt.Sprintf("github.com/chaincode/%s", fmt.Sprintf("%s-%s", ccName, cc.Version))
	exchangeMap["installCCPolicy"] = cc.Policy
	//环境变量
	exchangeMap["orgMSPID"] = fmt.Sprintf("(%s)", orgNameMSP)
	exchangeMap["orgPeerNum"] = fmt.Sprintf("(%s)", orgPeerNUM)
	exchangeMap["peerTlsCa"] = fmt.Sprintf("(%s)", peerTLSCaList)
	exchangeMap["peerUsersMsp"] = fmt.Sprintf("(%s)", peerUsersMspList)
	exchangeMap["peerAddress"] = fmt.Sprintf("(%s)", peerAddressList)
	exchangeMap["peerHosts"] = fmt.Sprintf("(%s)", peerHostsList)

	//创建脚本文件
	step1Buff, err := dcache.GetUpgradeChainCodeScriptStep1Template(general.Version, replaceMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateUpgradeChainCodeScript",
			Error: []string{"Build upgradeChainCodeStep1.sh replace parame error"},
		})
		err = errors.WithMessage(err, "Build upgradeChainCodeStep1.sh replace parame error")
		return err
	}
	scriptStep1SavePath := filepath.Join(general.BaseOutput, "upgradeChainCode", peerIP, folderName)
	err = ioutil.WriteFile(filepath.Join(scriptStep1SavePath, "upgradeChainCodeStep1.sh"), step1Buff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateUpgradeChainCodeScript",
			Error: []string{"replace parame writer upgradeChainCodeStep1.sh error"},
		})
		err = errors.WithMessage(err, "replace parame writer upgradeChainCodeStep1.sh error")
		return err
	}
	step2Buff, err := dcache.GetUpgradeChainCodeScriptStep2Template(general.Version, exchangeMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateUpgradeChainCodeScript",
			Error: []string{"Build upgradeChainCodeStep2.sh exchange parame error"},
		})
		err = errors.WithMessage(err, "Build upgradeChainCodeStep2.sh exchange parame error")
		return err
	}
	err = ioutil.WriteFile(filepath.Join(outputRoot, "upgradeChainCodeStep2.sh"), step2Buff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateUpgradeChainCodeScript",
			Error: []string{"exchange parame writer upgradeChainCodeStep2.sh error"},
		})
		err = errors.WithMessage(err, "exchange parame writer upgradeChainCodeStep2.sh error")
		return err
	}

	//拷贝合约
	dCCName := fmt.Sprintf("%s-%s", ccName, cc.Version)
	dst = filepath.Join(outputRoot, "chaincode", dCCName)
	_, err = os.Stat(dst)
	if err != nil {
		err = os.MkdirAll(dst, 0777)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateUpgradeChainCodeScript",
				Error: []string{"Build output folder errorr"},
			})
			return errors.WithMessage(err, "Build output folder error")
		}
	}
	ccSavepath := filepath.ToSlash(filepath.Join(dcache.GetTemplatePathByVersion(general.Version), "chaincode", ccName))
	tools.CopyFolder(dst, ccSavepath)
	return nil
}

//###################################创建停用合约脚本#####################################

//MakeStepGeneralDisableChainCodeScript 创建停用合约脚本
func MakeStepGeneralDisableChainCodeScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create disable chainCode script start"},
		})

		err := GeneralDisableChainCodeScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDisableChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateDisableChainCodeScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create disable chainCode script end"},
		})
		return nil
	}
}

//GeneralDisableChainCodeScript 为一次多安装合约
func GeneralDisableChainCodeScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralDisableChainCodeScriptFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateDisableChainCodeScript",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralDisableChainCodeScriptFile  创建脚本
func GeneralDisableChainCodeScriptFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	folderName := fmt.Sprintf("disableChainCode-%s-%s", ccName, cc.Version)
	//获取此任务之前最新订单信息
	indent, err := dmysql.GetStartTaskBeforIndent(general)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateDisableChainCodeScript",
			Error: []string{"get the latest indent error"},
		})
		return err
	}
	//组织数组
	if len(ExecCCIPList) != 0 {
		ExecCCIPList = nil
	}
	var endorseOrg = make(map[string]objectdefine.OrgType, len(cc.EndorsementOrg))
	for _, orgName := range cc.EndorsementOrg {
		endorseOrg[orgName] = indent.Org[orgName]
	}

	peerIPMapInfo := make(map[string][]objectdefine.PeerType)
	for _, org := range endorseOrg {
		for _, peer := range org.Peer {
			peerArray := make([]objectdefine.PeerType, 0)
			if v, ok := peerIPMapInfo[peer.IP]; ok {
				peerArray = v
				peerArray = append(peerArray, peer)
				peerIPMapInfo[peer.IP] = peerArray
			} else {
				ExecCCIPList = append(ExecCCIPList, peer.IP)
				peerArray = append(peerArray, peer)
				peerIPMapInfo[peer.IP] = peerArray
			}
		}
	}

	for ip, peerA := range peerIPMapInfo {
		//以每台机器为单位 创建一个脚本去执行
		outputRoot := filepath.Join(general.BaseOutput, "disableChainCode", ip, folderName, folderName)
		_, err := os.Stat(outputRoot)
		if err != nil {
			err = os.MkdirAll(outputRoot, 0777)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateDisableChainCodeScript",
					Error: []string{"Build output folder errorr"},
				})
				return errors.WithMessage(err, "Build output folder error")
			}
		}
		var ccContainerNameList string
		var peerContainerNameList string
		for _, peer := range peerA {
			ccContainerName := fmt.Sprintf("dev-%s-%s-%s", peer.Domain, ccName, cc.Version)
			if len(ccContainerNameList) == 0 {
				ccContainerNameList = ccContainerName
			} else {
				ccContainerNameList = ccContainerNameList + " " + ccContainerName
			}
			peerContainerName := peer.Domain
			if len(peerContainerNameList) == 0 {
				peerContainerNameList = peerContainerName
			} else {
				peerContainerNameList = peerContainerNameList + " " + peerContainerName
			}
		}
		replaceMap := make(map[string]string)
		replaceMap["ccContainerName"] = fmt.Sprintf("(%s)", ccContainerNameList)
		replaceMap["peerContainerName"] = fmt.Sprintf("(%s)", peerContainerNameList)
		replaceMap["ccName"] = fmt.Sprintf("%s.%s", ccName, cc.Version)
		//创建脚本文件
		step1Buff, err := dcache.GetDisableChainCodeScriptStepTemplate(general.Version, replaceMap)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateDisableChainCodeScript",
				Error: []string{"Build disableChainCodeScript.sh replace parame error"},
			})
			err = errors.WithMessage(err, "Build disableChainCodeScript.sh replace parame error")
			return err
		}
		err = ioutil.WriteFile(filepath.Join(outputRoot, "disableChainCodeScript.sh"), step1Buff, 0644)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateDisableChainCodeScript",
				Error: []string{"replace parame writer disableChainCodeScript.sh error"},
			})
			err = errors.WithMessage(err, "replace parame writer disableChainCodeScript.sh error")
			return err
		}
	}

	return nil
}

//###################################创建启用合约脚本#####################################

//MakeStepGeneralEnableChainCodeScript 创建启用合约脚本
func MakeStepGeneralEnableChainCodeScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create enable chainCode script start"},
		})

		err := GeneralEnableChainCodeScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailEnableChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateEnableChainCodeScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create enable chainCode script end"},
		})
		return nil
	}
}

//GeneralEnableChainCodeScript 为一次多安装合约
func GeneralEnableChainCodeScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralEnableChainCodeScriptFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateEnableChainCodeScript",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralEnableChainCodeScriptFile  创建脚本
func GeneralEnableChainCodeScriptFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	folderName := fmt.Sprintf("enableChainCode-%s-%s", ccName, cc.Version)
	//获取此任务之前最新订单信息
	indent, err := dmysql.GetStartTaskBeforIndent(general)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateEnableChainCodeScript",
			Error: []string{"get the latest indent error"},
		})
		return err
	}
	//组织数组
	if len(ExecCCIPList) != 0 {
		ExecCCIPList = nil
	}
	var endorseOrg = make(map[string]objectdefine.OrgType, len(cc.EndorsementOrg))
	for _, orgName := range cc.EndorsementOrg {
		endorseOrg[orgName] = indent.Org[orgName]
	}

	peerIPMapInfo := make(map[string][]objectdefine.PeerType)
	for _, org := range endorseOrg {
		for _, peer := range org.Peer {
			peerArray := make([]objectdefine.PeerType, 0)
			if v, ok := peerIPMapInfo[peer.IP]; ok {
				peerArray = v
				peerArray = append(peerArray, peer)
				peerIPMapInfo[peer.IP] = peerArray
			} else {
				ExecCCIPList = append(ExecCCIPList, peer.IP)
				peerArray = append(peerArray, peer)
				peerIPMapInfo[peer.IP] = peerArray
			}
		}
	}

	for ip, peerA := range peerIPMapInfo {
		//以每台机器为单位 创建一个脚本去执行
		outputRoot := filepath.Join(general.BaseOutput, "enableChainCode", ip, folderName, folderName)
		_, err := os.Stat(outputRoot)
		if err != nil {
			err = os.MkdirAll(outputRoot, 0777)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateEnableChainCodeScript",
					Error: []string{"Build output folder errorr"},
				})
				return errors.WithMessage(err, "Build output folder error")
			}
		}

		var peerContainerNameList string
		for _, peer := range peerA {

			peerContainerName := peer.Domain
			if len(peerContainerNameList) == 0 {
				peerContainerNameList = peerContainerName
			} else {
				peerContainerNameList = peerContainerNameList + " " + peerContainerName
			}
		}
		replaceMap := make(map[string]string)

		replaceMap["peerContainerName"] = fmt.Sprintf("(%s)", peerContainerNameList)
		replaceMap["ccName"] = fmt.Sprintf("%s.%s", ccName, cc.Version)
		//创建脚本文件
		step1Buff, err := dcache.GetEnableChainCodeScriptStepTemplate(general.Version, replaceMap)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateEnableChainCodeScript",
				Error: []string{"Build enableChainCodeScript.sh replace parame error"},
			})
			err = errors.WithMessage(err, "Build enableChainCodeScript.sh replace parame error")
			return err
		}
		err = ioutil.WriteFile(filepath.Join(outputRoot, "enableChainCodeScript.sh"), step1Buff, 0644)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateEnableChainCodeScript",
				Error: []string{"replace parame writer enableChainCodeScript.sh error"},
			})
			err = errors.WithMessage(err, "replace parame writer enableChainCodeScript.sh error")
			return err
		}
	}

	return nil
}

//##########################创建删除合约脚本#############################################

//MakeStepGeneralDeleteChainCodeScript 创建删除脚本
func MakeStepGeneralDeleteChainCodeScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"delete chainCode script start"},
		})

		err := GeneralDeleteChainCodeScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDeleteChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateDeleteChainCodeScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"delete chainCode script end"},
		})
		return nil
	}
}

//GeneralDeleteChainCodeScript 删除合约
func GeneralDeleteChainCodeScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralDeleteChainCodeScriptFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateDeleteChainCodeScript",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralDeleteChainCodeScriptFile 创建删除合约脚本
func GeneralDeleteChainCodeScriptFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	folderName := fmt.Sprintf("deleteChainCode-%s-%s", ccName, cc.Version)
	//获取此任务之前最新订单信息
	indent, err := dmysql.GetStartTaskBeforIndent(general)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateDeleteChainCodeScript",
			Error: []string{"get the latest indent error"},
		})
		return err
	}
	//组织数组
	if len(ExecCCIPList) != 0 {
		ExecCCIPList = nil
	}
	var endorseOrg = make(map[string]objectdefine.OrgType, len(cc.EndorsementOrg))
	for _, orgName := range cc.EndorsementOrg {
		endorseOrg[orgName] = indent.Org[orgName]
	}

	peerIPMapInfo := make(map[string][]objectdefine.PeerType)
	for _, org := range endorseOrg {
		for _, peer := range org.Peer {
			peerArray := make([]objectdefine.PeerType, 0)
			if v, ok := peerIPMapInfo[peer.IP]; ok {
				peerArray = v
				peerArray = append(peerArray, peer)
				peerIPMapInfo[peer.IP] = peerArray
			} else {
				ExecCCIPList = append(ExecCCIPList, peer.IP)
				peerArray = append(peerArray, peer)
				peerIPMapInfo[peer.IP] = peerArray
			}
		}
	}

	for ip, peerA := range peerIPMapInfo {
		//以每台机器为单位 创建一个脚本去执行
		outputRoot := filepath.Join(general.BaseOutput, "deleteChainCode", ip, folderName, folderName)
		_, err := os.Stat(outputRoot)
		if err != nil {
			err = os.MkdirAll(outputRoot, 0777)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateDeleteChainCodeScript",
					Error: []string{"Build output folder errorr"},
				})
				return errors.WithMessage(err, "Build output folder error")
			}
		}
		var ccContainerNameList string
		var peerContainerNameList string
		for _, peer := range peerA {
			ccContainerName := fmt.Sprintf("dev-%s-%s-%s", peer.Domain, ccName, cc.Version)
			if len(ccContainerNameList) == 0 {
				ccContainerNameList = ccContainerName
			} else {
				ccContainerNameList = ccContainerNameList + " " + ccContainerName
			}
			peerContainerName := peer.Domain
			if len(peerContainerNameList) == 0 {
				peerContainerNameList = peerContainerName
			} else {
				peerContainerNameList = peerContainerNameList + " " + peerContainerName
			}
		}
		replaceMap := make(map[string]string)
		replaceMap["ccContainerName"] = fmt.Sprintf("(%s)", ccContainerNameList)
		replaceMap["peerContainerName"] = fmt.Sprintf("(%s)", peerContainerNameList)
		replaceMap["ccName"] = fmt.Sprintf("%s.%s", ccName, cc.Version)
		//创建脚本文件
		step1Buff, err := dcache.GetDeleteChainCodeScriptStepTemplate(general.Version, replaceMap)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateDeleteChainCodeScript",
				Error: []string{"Build deleteChainCodeScript.sh replace parame error"},
			})
			err = errors.WithMessage(err, "Build deleteChainCodeScript.sh replace parame error")
			return err
		}
		err = ioutil.WriteFile(filepath.Join(outputRoot, "deleteChainCodeScript.sh"), step1Buff, 0644)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateDeleteChainCodeScript",
				Error: []string{"replace parame writer deleteChainCodeScript.sh error"},
			})
			err = errors.WithMessage(err, "replace parame writer deleteChainCodeScript.sh error")
			return err
		}
	}

	return nil
}
