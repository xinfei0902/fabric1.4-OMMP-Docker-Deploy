package daction

import (
	"deploy-server/app/dcache"
	"deploy-server/app/dmysql"
	"deploy-server/app/objectdefine"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
)

//MakeStepGeneralPeerBaseFile 构建base.yaml文件
func MakeStepGeneralPeerBaseFile(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add peer base.yaml start"},
		})

		err := GeneralCreatePeerBase(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreatePeerTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreatePeerBaseFile",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add peer base.yaml end"},
		})
		return nil
	}
}

//GeneralCreatePeerBase 以后多节点扩展
func GeneralCreatePeerBase(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreatePeerBaseFile(general, org, peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreatePeerBaseFile",
					Error: []string{err.Error()},
				})
				return err
			}
		}

	}
	return nil
}

//GeneralCreatePeerBaseFile 构建base.yaml文件
func GeneralCreatePeerBaseFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folderName := fmt.Sprintf("addPeer-%s-%s", orgOrder.Name, peerOrder.Name)
	outputRoot := filepath.Join(general.BaseOutput, "addPeer", peerOrder.IP, folderName, folderName)
	peerPort := peerOrder.Port
	peerDomain := peerOrder.Domain
	couchdbPort := peerOrder.CouchdbPort
	peerChaincodePort := peerOrder.ChaincodePort
	replaceMap := make(map[string]string)
	//couchdb数据库配置
	couchdbName := fmt.Sprintf("couchdb-%s", peerDomain)
	replaceMap["couchdbName"] = couchdbName
	replaceMap["couchdbPort"] = strconv.Itoa(couchdbPort)
	//peer配置
	otherPeerList, err := dmysql.GetOtherPeerList(general.ChannelName, orgOrder.Name)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerBaseFile",
			Error: []string{"Build base.yaml replace parame error"},
		})
		err = errors.WithMessage(err, "Build base.yaml replace parame error")
		return err
	}
	replaceMap["peerName"] = peerDomain
	replaceMap["peerDomain"] = peerDomain
	replaceMap["peerPort"] = strconv.Itoa(peerPort)
	replaceMap["peerGossipStrap"] = otherPeerList
	replaceMap["peerChaincodePort"] = strconv.Itoa(peerChaincodePort)
	orgName := fmt.Sprintf("%sMSP", orgOrder.Name)
	replaceMap["orgNameMSP"] = orgName
	peerMspPath := filepath.ToSlash(filepath.Join("crypto-config/peerOrganizations", orgOrder.OrgDomain, "peers", peerDomain, "msp"))
	replaceMap["peerMspPath"] = peerMspPath
	peerTLSPath := filepath.ToSlash(filepath.Join("crypto-config/peerOrganizations", orgOrder.OrgDomain, "peers", peerDomain, "tls"))
	replaceMap["peerTlsPath"] = peerTLSPath
	//cli配置
	replaceMap["cliName"] = fmt.Sprintf("cli-%s-%s", orgOrder.Name, peerOrder.Name)
	peerTLSServerCrtPath := filepath.ToSlash(filepath.Join(orgOrder.OrgDomain, "peers", peerDomain, "tls", "server.crt"))
	replaceMap["peerTlsServerCrtPath"] = peerTLSServerCrtPath
	peerTLSServerKeyPath := filepath.ToSlash(filepath.Join(orgOrder.OrgDomain, "peers", peerDomain, "tls", "server.key"))
	replaceMap["peerTlsServerKeyPath"] = peerTLSServerKeyPath
	peerCliTLSPath := filepath.ToSlash(filepath.Join(orgOrder.OrgDomain, "peers", peerDomain, "tls", "ca.crt"))
	replaceMap["peerCliTlsPath"] = peerCliTLSPath
	peerCliMspPath := filepath.ToSlash(filepath.Join(orgOrder.OrgDomain, "users", fmt.Sprintf("Admin@%s", orgOrder.OrgDomain), "msp"))
	replaceMap["peerCliMspPath"] = peerCliMspPath

	//peerStatusMonitor脚本
	exChange := make(map[string]string)
	exChange["checkPeerPort"] = strconv.Itoa(peerPort)
	exChange["peerNick"] = peerOrder.NickName
	exChange["accessKey"] = peerOrder.AccessKey
	exChange["peerID"] = strconv.Itoa(peerOrder.PeerID)
	exChange["peerIP"] = peerOrder.IP
	exChange["peerPort"] = strconv.Itoa(peerPort)
	//创建base.yaml文件
	baseBuff, err := dcache.GetPeerBaseFileTemplate(general.Version, replaceMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerBaseFile",
			Error: []string{"Build base.yaml replace parame error"},
		})
		err = errors.WithMessage(err, "Build base.yaml replace parame error")
		return err
	}
	err = ioutil.WriteFile(filepath.Join(outputRoot, "base.yaml"), baseBuff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerBaseFile",
			Error: []string{"replace parame writer base.yaml error"},
		})
		err = errors.WithMessage(err, "replace parame writer base.yaml error")
		return err
	}

	peerStatusBuff, err := dcache.GetPeerStatusFileTemplate(general.Version, exChange)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerBaseFile",
			Error: []string{"Build peerStatusMonitor.sh replace parame error"},
		})
		err = errors.WithMessage(err, "Build peerStatusMonitor.sh replace parame error")
		return err
	}
	statusFileName := fmt.Sprintf("peerStatusMonitor-%s.sh", peerDomain)
	err = ioutil.WriteFile(filepath.Join(outputRoot, statusFileName), peerStatusBuff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerBaseFile",
			Error: []string{"replace parame writer peerStatusMonitor.sh error"},
		})
		err = errors.WithMessage(err, "replace parame writer peerStatusMonitor.sh error")
		return err
	}
	return nil
}
