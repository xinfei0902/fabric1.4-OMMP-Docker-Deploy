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

//MakeStepGeneralOrgConfigtx 函数有执行任务顺序调用，实现构建configtx.yaml 文件以及证书等
func MakeStepGeneralOrgConfigtx(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add org configtx.yaml start"},
		})

		err := GeneralCreateOrgFolder(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateOrgConfigtx",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add org configtx.yaml end"},
		})
		return nil
	}
}

//GeneralCreateOrgFolder 创建用来生成新增组织整体目录
func GeneralCreateOrgFolder(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	//考虑如果支持多个组织新增情况 目前先暂定一次只新增一个组织 先写循环以备以后扩展
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreateConfigtxFile(general, org, peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateOrgConfigtx",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralCreateConfigtxFile 构建配置文件具体操作
func GeneralCreateConfigtxFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	//创建依赖包存放目录
	folder := fmt.Sprintf("addOrg-%s", orgOrder.Name)
	outputRoot := filepath.Join(general.BaseOutput, "addOrg", OperateAddOrgIP, folder, folder)
	startRoot := filepath.Join(general.BaseOutput, "addOrg", peerOrder.IP, folder, folder)
	err := os.MkdirAll(outputRoot, 0777)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgConfigtx",
			Error: []string{"Build output folder errorr"},
		})
		return errors.WithMessage(err, "Build output folder error")
	}
	//拷贝bin文件 包含configtxgen 等
	//src bin目录地址
	binPath := dcache.GetBinPathByVersion(general.Version)
	//dst 拷贝目的地址
	dst := filepath.Join(outputRoot, "bin")
	err = tools.CopyFolder(dst, binPath)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgConfigtx",
			Error: []string{"copy bin file error"},
		})
		return errors.WithMessage(err, "copy bin file error")
	}
	//拷贝证书
	src := filepath.Join(general.BaseOutput, "crypto-config")
	dst = filepath.Join(outputRoot, "crypto-config")
	tools.CopyFolder(dst, src)
	src = filepath.Join(general.BaseOutput, "crypto-config")
	dst = filepath.Join(startRoot, "crypto-config")
	tools.CopyFolder(dst, src)
	//map存放配置文件需要变量
	//目前新增组织只是限定一个
	replaceMap := make(map[string]string)
	orgName := fmt.Sprintf("%sMSP", orgOrder.Name)
	replaceMap["orgName"] = orgName
	orgMspPath := filepath.ToSlash(filepath.Join("crypto-config/peerOrganizations", orgOrder.OrgDomain, "msp"))
	replaceMap["mspPath"] = orgMspPath
	var peerDomain string
	var peerPort int
	for _, peer := range orgOrder.Peer {
		peerDomain = peer.Domain
		peerPort = peer.Port
		break
	}

	replaceMap["host"] = peerDomain
	replaceMap["port"] = strconv.Itoa(peerPort)
	//传进参数 返回整个构建完毕的配置文件
	buff, err := dcache.GetOrgConfigTemplate(general.Version, replaceMap)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgConfigtx",
			Error: []string{"Build config replace parame error"},
		})
		err = errors.WithMessage(err, "Build config error")
		return err
	}
	err = ioutil.WriteFile(filepath.Join(outputRoot, "configtx.yaml"), buff, 0644)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgConfigtx",
			Error: []string{"parame writer configtx.yaml error"},
		})
		err = errors.WithMessage(err, "Write config error")
		return err
	}
	return nil
}
