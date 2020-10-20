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
	"strings"

	"github.com/pkg/errors"
)

//MakeStepGeneralOrgBaseFile 构建base.yaml文件
func MakeStepGeneralOrgBaseFile(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add org base.yaml start"},
		})

		err := GeneralCreateOrgBase(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateOrgBaseFile",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add org base.yaml end"},
		})
		return nil
	}
}

//GeneralCreateOrgBase 以后多组织扩展
func GeneralCreateOrgBase(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreateOrgBaseFile(general, org, peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateOrgBaseFile",
					Error: []string{err.Error()},
				})
				return err
			}
		}

	}
	return nil
}

//GeneralCreateOrgBaseFile 构建base.yaml文件
func GeneralCreateOrgBaseFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folderName := fmt.Sprintf("addOrg-%s", orgOrder.Name)
	outputRoot := filepath.Join(general.BaseOutput, "addOrg", peerOrder.IP, folderName, folderName)
	orgOrgDomain := orgOrder.OrgDomain

	peerPort := peerOrder.Port
	peerDomain := peerOrder.Domain
	couchdbPort := peerOrder.CouchdbPort
	peerChaincodePort := peerOrder.ChaincodePort
	var caPort int
	for _, org := range general.Org {
		// for _, ca := range org.CA {
		caPort = org.CA.Port
		// 	break
		// }
		break
	}
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
	replaceMap["caPort"] = strconv.Itoa(caPort)
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

	//peerStatusMonitor脚本
	exChange := make(map[string]string)
	exChange["checkPeerPort"] = strconv.Itoa(peerPort)
	exChange["peerNick"] = peerOrder.NickName
	exChange["accessKey"] = peerOrder.AccessKey
	// peerCount, err := dmysql.GetAllPeerNum()
	// if err != nil {
	// 	output.AppendLog(&objectdefine.StepHistory{
	// 		Name:  "generalOrgBaseFile",
	// 		Error: []string{"get peer count error from mysql tables indent"},
	// 	})
	// 	err = errors.WithMessage(err, "get peer count error from mysql tables indent")
	// 	return err
	// }
	exChange["peerID"] = strconv.Itoa(peerOrder.PeerID)
	exChange["peerIP"] = peerOrder.IP
	exChange["peerPort"] = strconv.Itoa(peerPort)
	//创建base.yaml文件
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

	peerStatusBuff, err := dcache.GetPeerStatusFileTemplate(general.Version, exChange)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgBaseFile",
			Error: []string{"Build peerStatusMonitor.sh replace parame error"},
		})
		err = errors.WithMessage(err, "Build peerStatusMonitor.sh replace parame error")
		return err
	}
	statusFileName := fmt.Sprintf("peerStatusMonitor-%s.sh", peerDomain)
	err = ioutil.WriteFile(filepath.Join(outputRoot, statusFileName), peerStatusBuff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgBaseFile",
			Error: []string{"replace parame writer peerStatusMonitor.sh error"},
		})
		err = errors.WithMessage(err, "replace parame writer peerStatusMonitor.sh error")
		return err
	}

	return nil
}
