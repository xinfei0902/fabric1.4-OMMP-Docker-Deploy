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
	"strings"

	"github.com/pkg/errors"
)

var (
	//ChannelPeerIP 选择admin节点所在IP
	ChannelPeerIP string
	//ChannelPeerIPList 存放ip和目录对应 目前初版放弃
	ChannelPeerIPList = make(map[string][]string)
)

//MakeStepGeneralCreateChannelScript 创建脚本
func MakeStepGeneralCreateChannelScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create channel script start"},
		})

		err := GeneralCreateChannelScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateChannelTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChannelScript",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create channel  script end"},
		})
		return nil
	}
}

//GeneralCreateChannelScript 生成脚本
func GeneralCreateChannelScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreateChannelScriptFile(general, org, peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChannelScript",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralCreateChannelScriptFile 生成脚本
func GeneralCreateChannelScriptFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folderName := fmt.Sprintf("createChannel-%s", general.ChannelName)
	// foladaArray := make([]string, 0)
	// if foldList, ok := ChannelPeerIPList[peerOrder.IP]; ok {
	// 	foladaArray := foldList
	// 	foladaArray = append(foladaArray, folderName)
	// 	ChannelPeerIPList[peerOrder.IP] = foladaArray
	// } else {
	// 	foladaArray = append(foladaArray, folderName)
	// 	ChannelPeerIPList[peerOrder.IP] = foladaArray
	// }
	outputRoot := filepath.Join(general.BaseOutput, "createChannel", peerOrder.IP, folderName, folderName)
	err := os.MkdirAll(outputRoot, 0777)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChannelScript",
			Error: []string{"Build output folder errorr"},
		})
		return errors.WithMessage(err, "Build output folder error")
	}
	//拷贝channel-artifacts
	src := filepath.Join(general.BaseOutput, "channel-artifacts")
	dst := filepath.Join(outputRoot, "channel-artifacts")
	_, err = os.Stat(dst)
	if err != nil {
		err = os.MkdirAll(dst, 0777)
		if nil != err {

			return errors.WithMessage(err, "Make Output folder error")
		}
	}
	tools.CopyFolder(dst, src)
	//
	replaceMap := make(map[string]string)
	replaceMap["target"] = folderName
	replaceMap["newOrgCliName"] = peerOrder.CliName
	replaceMap["sourceOrgCliName"] = peerOrder.CliName
	orgMspName := fmt.Sprintf("%sMSP", orgOrder.Name)
	replaceMap["orgMSPName"] = orgMspName
	caFolder := fmt.Sprintf("ca-%s", orgOrder.Name)
	replaceMap["newOrgCAName"] = caFolder
	//目前创建通道每个组织就一个节点 故此启动CA 如果后期支持多节点 需要修改
	replaceMap["startCA"] = "1"
	replaceMap["peerName"] = peerOrder.Domain
	//第二步脚本
	exchangeMap := make(map[string]string)
	exchangeMap["target"] = folderName
	exchangeMap["newOrgCliName"] = peerOrder.CliName
	exchangeMap["channelName"] = general.ChannelName
	var ordererPort int
	var ordererOrgDomain string
	var ordererDomain string
	for _, orderer := range general.Orderer {
		ordererPort = orderer.Port
		ordererDomain = orderer.Domain
		ordererOrgDomain = orderer.OrgDomain
		break
	}
	ordererAddress := fmt.Sprintf("%s:%d", ordererDomain, ordererPort)
	exchangeMap["ordererAddress"] = ordererAddress
	ordererTLSCa := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations", ordererOrgDomain, "orderers", ordererDomain, "tls", "ca.crt"))
	exchangeMap["ordererTlsCa"] = ordererTLSCa
	//节点
	var orgName string
	//var peerIP string
	var orgNameMSP string
	var orgPeerNUM string
	var peerTLSCaList string
	var peerUsersMspList string
	var peerAddressList string
	var peerAdminTLSCaList string
	var peerAdminUsersMspList string
	var peerAdminAddressList string

	for _, org := range general.Org {
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
			peerTLSCA := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", org.OrgDomain, "peers", peer.Domain, "tls", "ca.crt"))
			if len(peerTLSCaList) == 0 {
				peerTLSCaList = peerTLSCA
			} else {
				peerTLSCaList = peerTLSCaList + " " + peerTLSCA
			}
			peerUserMSP := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", org.OrgDomain, "users", fmt.Sprintf("Admin@%s", org.OrgDomain), "msp"))
			if len(peerUsersMspList) == 0 {
				peerUsersMspList = peerUserMSP
			} else {
				peerUsersMspList = peerUsersMspList + " " + peerUserMSP
			}
			peerAddress := fmt.Sprintf("%s:%d", peer.Domain, peer.Port)
			if len(peerAddressList) == 0 {
				peerAddressList = peerAddress
			} else {
				peerAddressList = peerAddressList + " " + peerAddress
			}
			if len(ChannelPeerIP) == 0 {
				if peer.User == "Admin" {
					ChannelPeerIP = peer.IP
				}
			}
			if peer.User == "Admin" {
				peerAdminTLSCa := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", org.OrgDomain, "peers", peer.Domain, "tls", "ca.crt"))
				if len(peerAdminTLSCaList) == 0 {
					peerAdminTLSCaList = peerAdminTLSCa
				} else {
					peerAdminTLSCaList = peerAdminTLSCaList + " " + peerAdminTLSCa
				}
				peerAdminUserMSP := filepath.ToSlash(filepath.Join("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations", org.OrgDomain, "users", fmt.Sprintf("Admin@%s", org.OrgDomain), "msp"))
				if len(peerAdminUsersMspList) == 0 {
					peerAdminUsersMspList = peerAdminUserMSP
				} else {
					peerAdminUsersMspList = peerAdminUsersMspList + " " + peerAdminUserMSP
				}
				peerAdminAddress := fmt.Sprintf("%s:%d", peer.Domain, peer.Port)
				if len(peerAdminAddressList) == 0 {
					peerAdminAddressList = peerAdminAddress
				} else {
					peerAdminAddressList = peerAdminAddressList + " " + peerAdminAddress
				}
			}
		}
	}
	exchangeMap["orgMSPID"] = orgNameMSP
	exchangeMap["peerTlsCa"] = peerTLSCaList
	exchangeMap["peerUsersMsp"] = peerUsersMspList
	exchangeMap["peerAddress"] = peerAddressList

	//生成脚本

	step1Buff, err := dcache.GetChannelScriptStep1Template(general.Version, replaceMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChannelScript",
			Error: []string{"Build createChannelStep1.sh replace parame error"},
		})
		err = errors.WithMessage(err, "Build createChannelStep1.sh replace parame error")
		return err
	}
	scriptStep1SavePath := filepath.Join(general.BaseOutput, "createChannel", peerOrder.IP, folderName)
	err = ioutil.WriteFile(filepath.Join(scriptStep1SavePath, "createChannelStep1.sh"), step1Buff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChannelScript",
			Error: []string{"replace parame writer createChannelStep1.sh error"},
		})
		err = errors.WithMessage(err, "replace parame writer createChannelStep1.sh error")
		return err
	}
	step2Buff, err := dcache.GetChannelScriptStep2Template(general.Version, exchangeMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChannelScript",
			Error: []string{"Build createChannelStep2.sh exchange parame error"},
		})
		err = errors.WithMessage(err, "Build createChannelStep2.sh exchange parame error")
		return err
	}
	err = ioutil.WriteFile(filepath.Join(outputRoot, "createChannelStep2.sh"), step2Buff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChannelScript",
			Error: []string{"exchange parame writer createChannelStep2.sh error"},
		})
		err = errors.WithMessage(err, "exchange parame writer createChannelStep2.sh error")
		return err
	}

	return nil
}

//GeneralCreateChannelBaseYamlFile 构建base.yaml文件
func GeneralCreateChannelBaseYamlFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, outputRoot string, output *objectdefine.TaskNode) error {
	// folderName := fmt.Sprintf("addOrg-%s", orgOrder.Name)
	// outputRoot := filepath.Join(general.BaseOutput, "addOrg", peerOrder.IP, folderName, folderName)
	orgOrgDomain := orgOrder.OrgDomain

	peerPort := peerOrder.Port
	peerDomain := peerOrder.Domain
	couchdbPort := peerOrder.CouchdbPort
	peerChaincodePort := peerOrder.ChaincodePort
	// var caPort int
	// for _, org := range general.Org {
	// 	// for _, ca := range org.CA {
	// 	caPort = org.CA.Port
	// 	// 	break
	// 	// }
	// 	break
	// }
	////构建map 填写需要的组变量
	replaceMap := make(map[string]string)
	caFolder := fmt.Sprintf("ca-%s", orgOrder.Name)
	replaceMap["caName"] = caFolder
	cifeCaPath := filepath.ToSlash(filepath.Join(general.BaseOutput, "crypto-config/peerOrganizations", orgOrgDomain, "ca"))
	var cifeCaCertPath string
	var cifeCaSKPath string
	filepath.Walk(cifeCaPath, func(path string, info os.FileInfo, err error) error {
		folder := info.Name()
		if strings.Contains(folder, "_sk") {
			cifeCaSKPath = folder
		} else if strings.Contains(folder, ".pem") {
			cifeCaCertPath = folder
		}
		return nil
	})
	cifeTLSPath := filepath.ToSlash(filepath.Join(general.BaseOutput, "crypto-config/peerOrganizations", orgOrgDomain, "tlsca"))
	var cifeTLSCertPath string
	var cifeTLSSKPath string
	filepath.Walk(cifeTLSPath, func(path string, info os.FileInfo, err error) error {
		folder := info.Name()
		if strings.Contains(folder, "_sk") {
			cifeTLSSKPath = folder
		} else if strings.Contains(folder, ".pem") {
			cifeTLSCertPath = folder
		}
		return nil
	})
	replaceMap["caCertPem"] = cifeCaCertPath
	replaceMap["caSK"] = cifeCaSKPath
	replaceMap["tlsCertPem"] = cifeTLSCertPath
	replaceMap["tlsSk"] = cifeTLSSKPath
	replaceMap["caPort"] = strconv.Itoa(orgOrder.CA.Port)
	caCertPath := filepath.ToSlash(filepath.Join("crypto-config/peerOrganizations", orgOrgDomain, "ca"))
	replaceMap["caCertPath"] = caCertPath
	tlsCaPath := filepath.ToSlash(filepath.Join("crypto-config/peerOrganizations", orgOrgDomain, "tlsca"))
	replaceMap["tlsCaPath"] = tlsCaPath
	//couchdb数据库配置
	couchdbName := fmt.Sprintf("couchdb-%s", peerDomain)
	replaceMap["couchdbName"] = couchdbName
	replaceMap["couchdbPort"] = strconv.Itoa(couchdbPort)
	//peer配置
	replaceMap["peerName"] = peerDomain
	replaceMap["peerDomain"] = peerDomain
	replaceMap["peerPort"] = strconv.Itoa(peerPort)
	replaceMap["peerChaincodePort"] = strconv.Itoa(peerChaincodePort)
	orgName := fmt.Sprintf("%sMSP", orgOrder.Name)
	replaceMap["orgNameMSP"] = orgName
	peerMspPath := filepath.ToSlash(filepath.Join("crypto-config/peerOrganizations", orgOrder.OrgDomain, "peers", peerDomain, "msp"))
	replaceMap["peerMspPath"] = peerMspPath
	peerTLSPath := filepath.ToSlash(filepath.Join("crypto-config/peerOrganizations", orgOrder.OrgDomain, "peers", peerDomain, "tls"))
	replaceMap["peerTlsPath"] = peerTLSPath
	//cli配置
	replaceMap["cliName"] = peerOrder.CliName
	peerTLSServerCrtPath := filepath.ToSlash(filepath.Join(orgOrder.OrgDomain, "peers", peerDomain, "tls", "server.crt"))
	replaceMap["peerTlsServerCrtPath"] = peerTLSServerCrtPath
	peerTLSServerKeyPath := filepath.ToSlash(filepath.Join(orgOrder.OrgDomain, "peers", peerDomain, "tls", "server.key"))
	replaceMap["peerTlsServerKeyPath"] = peerTLSServerKeyPath
	peerCliTLSPath := filepath.ToSlash(filepath.Join(orgOrder.OrgDomain, "peers", peerDomain, "tls", "ca.crt"))
	replaceMap["peerCliTlsPath"] = peerCliTLSPath
	peerCliMspPath := filepath.ToSlash(filepath.Join(orgOrder.OrgDomain, "users", fmt.Sprintf("Admin@%s", orgOrder.OrgDomain), "msp"))
	replaceMap["peerCliMspPath"] = peerCliMspPath

	//创建base.yaml文件
	//创建三个脚本文件
	baseBuff, err := dcache.GetOrgBaseFileTemplate(general.Version, replaceMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalOrgBaseFile",
			Error: []string{"Build base.yaml replace parame error"},
		})
		err = errors.WithMessage(err, "Build base.yaml replace parame error")
		return err
	}
	err = ioutil.WriteFile(filepath.Join(outputRoot, "base.yaml"), baseBuff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalOrgBaseFile",
			Error: []string{"replace parame writer base.yaml error"},
		})
		err = errors.WithMessage(err, "replace parame writer base.yaml error")
		return err
	}
	return nil
}
