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
	var channeNameArray string
	var OrgNameArray string
	var peerConfigArray string
	var ccNameArray string
	var ccEndorOrgArray string
	var ccInstallConfigArray string
	//处理变量
	for channelName, deploy := range general.Deploy {
		if len(channelName) == 0 {
			channeNameArray = "\"" + channelName + "\""
		} else {
			channeNameArray = channeNameArray + " " + "\"" + channelName + "\""
		}
		var orgNameString string
		for orgName, org := range deploy.JoinOrg {
			if len(orgNameString) == 0 {
				orgNameString = orgName
			} else {
				orgNameString = orgNameString + " " + orgName
			}
			var peerConfig string
			for _, peer := range org.Peer {
				if len(peerConfig) == 0 {
					peerConfig = orgName
					peerConfig = peerConfig + " " + fmt.Sprintf("%sMSP", orgName)
				} else {
					peerConfig = peerConfig + " " + fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/%s/peers/%s/tls/ca.crt", org.OrgDomain, peer.Domain)
					peerConfig = peerConfig + " " + fmt.Sprintf("/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/%s/users/Admin@%s/msp", org.OrgDomain, org.OrgDomain)
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
			OrgNameArray = "\"" + orgNameString + "\""
		} else {
			OrgNameArray = OrgNameArray + " " + "\"" + orgNameString + "\""
		}

		//合约
		var ccNameString string
		var ccInstllString string
		for ccName, cc := range deploy.JoinCC {
			if len(ccNameString) == 0 {
				ccNameString = ccName
			} else {
				ccNameString = ccNameString + " " + ccName
			}
			var endorOrgString string
			for _, endorOrg := range cc.EndorsementOrg {
				if len(endorOrgString) == 0 {
					endorOrgString = endorOrg
				} else {
					endorOrgString = endorOrgString + " " + endorOrg
				}
			}
			if len(ccEndorOrgArray) == 0 {
				ccEndorOrgArray = "\"" + endorOrgString + "\""
			} else {
				ccEndorOrgArray = ccEndorOrgArray + " " + "\"" + endorOrgString + "\""
			}
			if len(ccInstllString) == 0 {
				ccInstllString = ccName + " " + cc.Version + " " + cc.Policy
			} else {
				ccInstllString = ccInstllString + " " + ccName + " " + cc.Version + " " + cc.Policy
			}
		}
		if len(ccNameArray) == 0 {
			ccNameArray = "\"" + ccNameString + "\""
		} else {
			ccNameArray = ccNameArray + " " + "\"" + ccNameString + "\""
		}
		if len(ccInstallConfigArray) == 0 {
			ccInstallConfigArray = "\"" + ccInstllString + "\""
		} else {
			ccInstallConfigArray = ccInstallConfigArray + " " + "\"" + ccInstllString + "\""
		}

	}
	//###########################
	replaceMap["channelNameArray"] = fmt.Sprintf("(%s)", channeNameArray)
	replaceMap["orgNameArray"] = fmt.Sprintf("(%s)", OrgNameArray)
	replaceMap["peerConfigArray"] = fmt.Sprintf("(%s)", peerConfigArray)
	replaceMap["ccNameArray"] = fmt.Sprintf("(%s)", ccNameArray)
	replaceMap["ccOrgArray"] = fmt.Sprintf("(%s)", ccEndorOrgArray)
	replaceMap["ccInstallConfigArray"] = fmt.Sprintf("(%s)", ccInstallConfigArray)
	step1Buff, err := dcache.GetCompleteDeployScriptTemplate(general.Version, replaceMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateCompleteDeployScript",
			Error: []string{"Build script.sh replace parame error"},
		})
		err = errors.WithMessage(err, "Build script.sh replace parame error")
		return err
	}
	scriptStep1SavePath := filepath.Join(general.BaseOutput, "script")
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
	return nil
}
