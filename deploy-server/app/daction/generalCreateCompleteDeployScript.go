package daction

import (
	"deploy-server/app/dcache"
	"deploy-server/app/objectdefine"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

//MakeStepGeneralCreateCompleteDeployScript 创建一键部署script文件
func MakeStepGeneralCreateCompleteDeployScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create complete deploy base.yaml start"},
		})

		err := GeneralCreateCompleteDeployScript(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create complete deploy base.yaml"},
		})
		return nil
	}
}

func GeneralCreateCompleteDeployScript(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	replaceMap := make(map[string]string)
	//配置脚本orderer环境变量
	var ordererPort int
	var ordererOrgDomain string
	var ordererDomain string
	for _, orderer := range general.Orderer {
		//ordererIP = orderer.IP
		ordererPort = orderer.Port
		ordererDomain = orderer.Domain
		ordererOrgDomain = orderer.OrgDomain
		break
	}
	ordererAddress := fmt.Sprintf("%s:%d", ordererDomain, ordererPort)
	replaceMap["ordererAddress"] = ordererAddress
	ordererTLSCa := filepath.ToSlash(filepath.Join(ordererOrgDomain, "orderers", ordererDomain, "tls", "ca.crt"))
	replaceMap["ordererTlsCa"] = ordererTLSCa
	//###############################
	var SHChanneName string
	var SHCreateOrgMSPID string
	var SHCreateTlsRoot string
	var SHCreateMSP string
	var SHCreatePeerAddress string
	var OrgNameArray string
	var peerConfigArray string
	var ccNameArray string
	var ccEndorOrgArray string
	var ccInstallConfigArray string
	var isInstllCC string
	//
	if len(general.Deploy) == 0{
		isInstllCC = "off"
		channelName := general.ChannelName
		if len(channelName) == 0 {
			return errors.New("channel name empty")
		} else {
			SHChanneName = "\"" + channelName + "\""
		}
		var orgNameString string
		for orgName,org := range general.Org{
			if len(orgNameString) == 0 {
				orgNameString = orgName
			} else {
				orgNameString = orgNameString + " " + orgName
			}
			var peerConfig string
		     peertlsConfig := make(map[string]string,len(org.Peer))
			 peermspConfig := make(map[string]string,len(org.Peer))
			for _,peer := range org.Peer{
				if peer.User == "Admin"{
				    tlsConfig := fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/%s/peers/%s/tls/ca.crt", org.OrgDomain, peer.Domain)
					mspConfig := fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/%s/users/Admin@%s/msp", org.OrgDomain, org.OrgDomain)
				    value := fmt.Sprintf("%s",org.Name+peer.Name)
					peertlsConfig[value] = tlsConfig
					peermspConfig[value] = mspConfig
				}else{
					tlsConfig := fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/%s/peers/%s/tls/ca.crt", org.OrgDomain, peer.Domain)
					mspConfig := fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/%s/users/Admin@%s/msp", org.OrgDomain,org.OrgDomain)
					value := fmt.Sprintf("%s",org.Name+peer.Name)
					peertlsConfig[value] = tlsConfig
					peermspConfig[value] = mspConfig
				}
			}
			for _, peer := range org.Peer {
				if fmt.Sprintf("cli-%s-%s", org.Name, peer.Name) == completeDeployCliName{
					SHCreateOrgMSPID = fmt.Sprintf("%sMSP", orgName)
					SHCreateTlsRoot = fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/%s/peers/%s/tls/ca.crt", org.OrgDomain, peer.Domain)
					SHCreateMSP = fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/%s/users/Admin@%s/msp", org.OrgDomain, org.OrgDomain)
					SHCreatePeerAddress = fmt.Sprintf("%s:%d", peer.Domain, peer.Port)
				}
				if len(peerConfig) == 0 {
					peerConfig = orgName
					peerConfig = peerConfig + " " + fmt.Sprintf("%sMSP", orgName)
					value := fmt.Sprintf("%s",org.Name+peer.Name)
					peerConfig =peerConfig + " " +peertlsConfig[value]
					peerConfig =peerConfig + " " +peermspConfig[value]
					peerConfig = peerConfig + " " + fmt.Sprintf("%s:%d", peer.Domain, peer.Port)
				} else {
					value := fmt.Sprintf("%s",org.Name+peer.Name)
					peerConfig =peerConfig + " " +peertlsConfig[value]
					peerConfig =peerConfig + " " +peermspConfig[value]
					peerConfig = peerConfig + " " + fmt.Sprintf("%s:%d", peer.Domain, peer.Port)
				}
			}
			if len(peerConfigArray) == 0 {
				peerConfigArray = "\"" + peerConfig + "\""
			} else {
				peerConfigArray = peerConfigArray + " " + "\"" + peerConfig + "\""
			}

		}
		if len(OrgNameArray) == 0 {
			OrgNameArray = " "+ orgNameString +" "
		} else {
			OrgNameArray = OrgNameArray + " " + "" + orgNameString + ""
		}
	}else{
		isInstllCC = "on"
    } 
	//###########################
	replaceMap["channelName"] = fmt.Sprintf("%s", SHChanneName)
	replaceMap["createChannelMSPID"] = SHCreateOrgMSPID
	replaceMap["createChannelTlsRoot"] = SHCreateTlsRoot
	replaceMap["createChannelMsp"] = SHCreateMSP
	replaceMap["createChannelPeerAddress"] = SHCreatePeerAddress

	replaceMap["orgNameArray"] = fmt.Sprintf("(%s)", OrgNameArray)
	replaceMap["peerConfigArray"] = fmt.Sprintf("(%s)", peerConfigArray)
	replaceMap["ccNameArray"] = fmt.Sprintf("(%s)", ccNameArray)
	replaceMap["ccOrgArray"] = fmt.Sprintf("(%s)", ccEndorOrgArray)
	replaceMap["ccInstallConfigArray"] = fmt.Sprintf("(%s)", ccInstallConfigArray)
	replaceMap["isInstallCC"] = isInstllCC
	step1Buff, err := dcache.GetCompleteDeployScriptTemplate(general.Version, replaceMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateCompleteDeployScript",
			Error: []string{"Build script.sh replace parame error"},
		})
		err = errors.WithMessage(err, "Build script.sh replace parame error")
		return err
	}
	scriptStep1SavePath := filepath.Join(general.BaseOutput, "deploy","scripts")
	_, err = os.Stat(scriptStep1SavePath)
	if err != nil {
		err := os.MkdirAll(scriptStep1SavePath, 0777)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateCompleteDeployScript",
				Error: []string{"Build output folder errorr"},
			})
			return errors.WithMessage(err, "Build output folder error")
		}
	}

	err = ioutil.WriteFile(filepath.Join(scriptStep1SavePath, "script.sh"), step1Buff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateCompleteDeployScript",
			Error: []string{"replace parame writer script.sh error"},
		})
		err = errors.WithMessage(err, "replace parame writer script.sh error")
		return err
	}

	exchangeMap := make(map[string]string)
	//exchangeMap["target"] = "deploy"
	exchangeMap["cliName"] = CliName
	step2Buff, err := dcache.GetCompleteDeployScriptDeployTemplate(general.Version, exchangeMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateCompleteDeployScript",
			Error: []string{"Build deploy.sh replace parame error"},
		})
		err = errors.WithMessage(err, "Build deploy.sh replace parame error")
		return err
	}
	_, err = os.Stat(general.BaseOutput)
	if err != nil {
		err := os.MkdirAll(scriptStep1SavePath, 0777)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateCompleteDeployScript",
				Error: []string{"Build output folder errorr"},
			})
			return errors.WithMessage(err, "Build output folder error")
		}
	}

	err = ioutil.WriteFile(filepath.Join(general.BaseOutput,"deploy", "deploy.sh"), step2Buff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateCompleteDeployScript",
			Error: []string{"replace parame writer script.sh error"},
		})
		err = errors.WithMessage(err, "replace parame writer script.sh error")
		return err
	}

	return nil
}
