package daction

import (
	"bytes"
	"deploy-server/app/dconfig"
	"deploy-server/app/dmysql"
	"deploy-server/app/objectdefine"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/pkg/errors"
)

//######################组织上传###################################

//MakeStepGeneralOrgZipUpload 组织打包文件上传
func MakeStepGeneralOrgZipUpload(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add org upload zip start"},
		})

		err := GeneralCreateOrgZipUpload(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateOrgZipUpload",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add org upload zip end"},
		})
		return nil
	}
}

//GeneralCreateOrgZipUpload 以后多组织扩展
func GeneralCreateOrgZipUpload(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreateOrgZipUploadFile(general, org, peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateOrgZipUpload",
					Error: []string{err.Error()},
				})
				return err
			}
		}
		err := GeneralCreateOperateOrgZipUploadFile(general, org, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateOrgZipUpload",
				Error: []string{err.Error()},
			})
			return err
		}

	}
	return nil
}

//GeneralCreateOrgZipUploadFile 组织上传
func GeneralCreateOrgZipUploadFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("addOrg-%s", orgOrder.Name)
	outputRoot := filepath.Join(general.BaseOutput, "addOrg")
	path := filepath.Join(outputRoot, folder+peerOrder.IP+".tar.gz")

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/file/upload", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateOrgZipUpload",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	extraParams := map[string]string{
		"title":       "My Document",
		"author":      "Matt Aimonetti",
		"description": "A document with all the Go programming language secrets",
	}
	request, err := newfileUploadRequest(url, extraParams, "file", path, "", output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgZipUpload",
			Error: []string{fmt.Sprintf("fail upload request fail:%s", err.Error())},
		})
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgZipUpload",
			Error: []string{fmt.Sprintf("exec http post clent.DO fail:%s", err.Error())},
		})
		return err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgZipUpload",
			Error: []string{fmt.Sprintf("read request return body fail:%s", err.Error())},
		})
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 0 {
		err = errors.WithMessage(err, "return Status code err")
		return err
	}

	return nil
}

//GeneralCreateOperateOrgZipUploadFile 组织操作上传
func GeneralCreateOperateOrgZipUploadFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("addOrg-%s", orgOrder.Name)
	outputRoot := filepath.Join(general.BaseOutput, "addOrg")
	path := filepath.Join(outputRoot, folder+OperateAddOrgIP+".tar.gz")

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/file/upload", OperateAddOrgIP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateOperateOrgZipUpload",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	extraParams := map[string]string{
		"title":       "My Document",
		"author":      "Matt Aimonetti",
		"description": "A document with all the Go programming language secrets",
	}
	request, err := newfileUploadRequest(url, extraParams, "file", path, "", output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgZipUpload",
			Error: []string{fmt.Sprintf("fail upload request fail:%s", err.Error())},
		})
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgZipUpload",
			Error: []string{fmt.Sprintf("exec http post clent.DO fail:%s", err.Error())},
		})
		return err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgZipUpload",
			Error: []string{fmt.Sprintf("read request return body fail:%s", err.Error())},
		})
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 0 {
		err = errors.WithMessage(err, "return Status code err")
		return err
	}

	return nil
}

//######################删除组织上传###################################

//MakeStepGeneralDeleteOrgZipUpload 组织打包文件上传
func MakeStepGeneralDeleteOrgZipUpload(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete org upload zip start"},
		})

		err := GeneralDeleteOrgZipUpload(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDeleteOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeleteOrgZipUpload",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete org upload zip end"},
		})
		return nil
	}
}

//GeneralDeleteOrgZipUpload 以后多组织扩展
func GeneralDeleteOrgZipUpload(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	peerAllIP := make(map[string]string, 0)
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			if _, ok := peerAllIP[peer.IP]; !ok {
				peerAllIP[peer.IP] = peer.Name
				err := GeneralDeleteOrgZipUploadFile(general, org, peer, output)
				if err != nil {
					output.AppendLog(&objectdefine.StepHistory{
						Name:  "generalDeleteOrgZipUpload",
						Error: []string{err.Error()},
					})
					return err
				}
			}

		}
		err := GeneralDeleteOperateOrgZipUploadFile(general, org, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalDeleteOrgZipUpload",
				Error: []string{err.Error()},
			})
			return err
		}

	}
	return nil
}

//GeneralDeleteOrgZipUploadFile 组织上传
func GeneralDeleteOrgZipUploadFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("deleteOrg-%s", orgOrder.Name)
	outputRoot := filepath.Join(general.BaseOutput, "deleteOrg")
	path := filepath.Join(outputRoot, folder+peerOrder.IP+".tar.gz")

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/file/upload", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalDeleteOrgZipUpload",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	extraParams := map[string]string{
		"title":       "My Document",
		"author":      "Matt Aimonetti",
		"description": "A document with all the Go programming language secrets",
	}
	request, err := newfileUploadRequest(url, extraParams, "file", path, "", output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalDeleteOrgZipUpload",
			Error: []string{fmt.Sprintf("fail upload request fail:%s", err.Error())},
		})
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalDeleteOrgZipUpload",
			Error: []string{fmt.Sprintf("exec http post clent.DO fail:%s", err.Error())},
		})
		return err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalDeleteOrgZipUpload",
			Error: []string{fmt.Sprintf("read request return body fail:%s", err.Error())},
		})
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 0 {
		err = errors.WithMessage(err, "return Status code err")
		return err
	}

	return nil
}

//GeneralDeleteOperateOrgZipUploadFile 组织操作上传
func GeneralDeleteOperateOrgZipUploadFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("deleteOrg-%s", orgOrder.Name)
	outputRoot := filepath.Join(general.BaseOutput, "deleteOrg")
	path := filepath.Join(outputRoot, folder+OperateAddOrgIP+".tar.gz")

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/file/upload", OperateAddOrgIP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalDeleteOrgZipUpload",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	extraParams := map[string]string{
		"title":       "My Document",
		"author":      "Matt Aimonetti",
		"description": "A document with all the Go programming language secrets",
	}
	request, err := newfileUploadRequest(url, extraParams, "file", path, "", output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalDeleteOrgZipUpload",
			Error: []string{fmt.Sprintf("fail upload request fail:%s", err.Error())},
		})
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalDeleteOrgZipUpload",
			Error: []string{fmt.Sprintf("exec http post clent.DO fail:%s", err.Error())},
		})
		return err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalDeleteOrgZipUpload",
			Error: []string{fmt.Sprintf("read request return body fail:%s", err.Error())},
		})
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 0 {
		err = errors.WithMessage(err, "return Status code err")
		return err
	}

	return nil
}

//######################节点上传###################################

//MakeStepGeneralPeerZipUpload 节点打包文件上传
func MakeStepGeneralPeerZipUpload(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add peer upload zip start"},
		})

		err := GeneralCreatePeerZipUpload(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreatePeerTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreatePeerZipUpload",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add peer upload zip end"},
		})
		return nil
	}
}

//GeneralCreatePeerZipUpload 以后多组织扩展
func GeneralCreatePeerZipUpload(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreatePeerZipUploadFile(general, org, peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreatePeerZipUpload",
					Error: []string{err.Error()},
				})
				return err
			}
		}

	}
	return nil
}

//GeneralCreatePeerZipUploadFile 节点上传
func GeneralCreatePeerZipUploadFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("addPeer-%s-%s", orgOrder.Name, peerOrder.Name)
	outputRoot := filepath.Join(general.BaseOutput, "addPeer")
	path := filepath.Join(outputRoot, folder+".tar.gz")

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/file/upload", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreatePeerZipUpload",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	extraParams := map[string]string{
		"title":       "My Document",
		"author":      "Matt Aimonetti",
		"description": "A document with all the Go programming language secrets",
	}
	request, err := newfileUploadRequest(url, extraParams, "file", path, "", output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerZipUpload",
			Error: []string{fmt.Sprintf("fail upload request fail:%s", err.Error())},
		})
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerZipUpload",
			Error: []string{fmt.Sprintf("exec http post clent.DO fail:%s", err.Error())},
		})
		return err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerZipUpload",
			Error: []string{fmt.Sprintf("read request return body fail:%s", err.Error())},
		})
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 0 {
		err = errors.WithMessage(err, "return Status code err")
		return err
	}

	return nil
}

//######################删除节点上传###################################

//MakeStepGeneralDeletePeerZipUpload 节点打包文件上传
func MakeStepGeneralDeletePeerZipUpload(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete peer upload zip start"},
		})

		err := GeneralDeletePeerZipUpload(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDeletePeerTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeletePeerZipUpload",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete peer upload zip end"},
		})
		return nil
	}
}

//GeneralDeletePeerZipUpload 以后多组织扩展
func GeneralDeletePeerZipUpload(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralDeletePeerZipUploadFile(general, org, peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeletePeerZipUpload",
					Error: []string{err.Error()},
				})
				return err
			}
		}

	}
	return nil
}

//GeneralDeletePeerZipUploadFile 节点上传
func GeneralDeletePeerZipUploadFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("deletePeer-%s-%s", orgOrder.Name, peerOrder.Name)
	outputRoot := filepath.Join(general.BaseOutput, "deletePeer")
	path := filepath.Join(outputRoot, folder+".tar.gz")

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/file/upload", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalDeletePeerZipUpload",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	extraParams := map[string]string{
		"title":       "My Document",
		"author":      "Matt Aimonetti",
		"description": "A document with all the Go programming language secrets",
	}
	request, err := newfileUploadRequest(url, extraParams, "file", path, "", output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalDeletePeerZipUpload",
			Error: []string{fmt.Sprintf("fail upload request fail:%s", err.Error())},
		})
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalDeletePeerZipUpload",
			Error: []string{fmt.Sprintf("exec http post clent.DO fail:%s", err.Error())},
		})
		return err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalDeletePeerZipUpload",
			Error: []string{fmt.Sprintf("read request return body fail:%s", err.Error())},
		})
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 0 {
		err = errors.WithMessage(err, "return Status code err")
		return err
	}

	return nil
}

//######################合约上传###################################

//MakeStepGeneralChainCodeZipUpload 合约打包文件上传
func MakeStepGeneralChainCodeZipUpload(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add chainCode upload zip start"},
		})

		err := GeneralCreateChainCodeZipUpload(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChainCodeZipUpload",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add chainCode upload zip end"},
		})
		return nil
	}
}

//GeneralCreateChainCodeZipUpload 以后多链码扩展
func GeneralCreateChainCodeZipUpload(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralCreateChainCodeZipUploadFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChainCodeZipUpload",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralCreateChainCodeZipUploadFile 合约上传
func GeneralCreateChainCodeZipUploadFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("addChainCode-%s-%s", ccName, cc.Version)
	outputRoot := filepath.Join(general.BaseOutput, "addChainCode")
	path := filepath.Join(outputRoot, folder+".tar.gz")

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/file/upload", SelectExecCCIP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateChainCodeZipUpload",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	extraParams := map[string]string{
		"title":       "My Document",
		"author":      "Matt Aimonetti",
		"description": "A document with all the Go programming language secrets",
	}
	request, err := newfileUploadRequest(url, extraParams, "file", path, "", output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChainCodeZipUpload",
			Error: []string{fmt.Sprintf("fail upload request fail:%s", err.Error())},
		})
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChainCodeZipUpload",
			Error: []string{fmt.Sprintf("exec http post clent.DO fail:%s", err.Error())},
		})
		return err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChainCodeZipUpload",
			Error: []string{fmt.Sprintf("read request return body fail:%s", err.Error())},
		})
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 0 {
		err = errors.WithMessage(err, "return Status code err")
		return err
	}

	return nil
}

//######################合约删除上传###################################

//MakeStepGeneralDeleteChainCodeZipUpload 合约打包文件上传
func MakeStepGeneralDeleteChainCodeZipUpload(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete chainCode upload zip start"},
		})

		err := GeneralDeleteChainCodeZipUpload(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete chainCode upload zip end"},
		})
		return nil
	}
}

//GeneralDeleteChainCodeZipUpload 以后多链码扩展
func GeneralDeleteChainCodeZipUpload(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralDeleteChainCodeZipUploadFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateDeleteChainCodeZipUpload",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralDeleteChainCodeZipUploadFile 合约停用上传
func GeneralDeleteChainCodeZipUploadFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {

	for _, ipaddress := range ExecCCIPList {
		err := GeneralDeleteChainCodeZipUploadStep(general, ipaddress, output)
		if err != nil {
			return err
		}
	}
	return nil
}

//GeneralDeleteChainCodeZipUploadStep 分步上传
func GeneralDeleteChainCodeZipUploadStep(general *objectdefine.Indent, ipAddress string, output *objectdefine.TaskNode) error {
	outputRoot := filepath.Join(general.BaseOutput, "deleteChainCode")
	path := filepath.Join(outputRoot, ipAddress+".tar.gz")

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/file/upload", ipAddress, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateDeleteChainCodeZipUpload",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	extraParams := map[string]string{
		"title":       "My Document",
		"author":      "Matt Aimonetti",
		"description": "A document with all the Go programming language secrets",
	}
	request, err := newfileUploadRequest(url, extraParams, "file", path, "", output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateDeleteChainCodeZipUpload",
			Error: []string{fmt.Sprintf("fail upload request fail:%s", err.Error())},
		})
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateDeleteChainCodeZipUpload",
			Error: []string{fmt.Sprintf("exec http post clent.DO fail:%s", err.Error())},
		})
		return err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateDeleteChainCodeZipUpload",
			Error: []string{fmt.Sprintf("read request return body fail:%s", err.Error())},
		})
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 0 {
		err = errors.WithMessage(err, "return Status code err")
		return err
	}
	return nil
}

//######################合约升级上传###################################

//MakeStepGeneralUpgradeChainCodeZipUpload 合约打包文件上传
func MakeStepGeneralUpgradeChainCodeZipUpload(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create upgrade chainCode upload zip start"},
		})

		err := GeneralUpgradeChainCodeZipUpload(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailUpgradeChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalUpgradeChainCodeZipUpload",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create upgrade chainCode upload zip end"},
		})
		return nil
	}
}

//GeneralUpgradeChainCodeZipUpload 以后多链码扩展
func GeneralUpgradeChainCodeZipUpload(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralUpgradeChainCodeZipUploadFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalUpgradeChainCodeZipUpload",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralUpgradeChainCodeZipUploadFile 合约上传
func GeneralUpgradeChainCodeZipUploadFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("upgradeChainCode-%s-%s", ccName, cc.Version)
	outputRoot := filepath.Join(general.BaseOutput, "upgradeChainCode")
	path := filepath.Join(outputRoot, folder+".tar.gz")

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/file/upload", SelectExecCCIP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalUpgradeChainCodeZipUpload",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	extraParams := map[string]string{
		"title":       "My Document",
		"author":      "Matt Aimonetti",
		"description": "A document with all the Go programming language secrets",
	}
	request, err := newfileUploadRequest(url, extraParams, "file", path, "", output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalUpgradeChainCodeZipUpload",
			Error: []string{fmt.Sprintf("fail upload request fail:%s", err.Error())},
		})
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalUpgradeChainCodeZipUpload",
			Error: []string{fmt.Sprintf("exec http post clent.DO fail:%s", err.Error())},
		})
		return err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalUpgradeChainCodeZipUpload",
			Error: []string{fmt.Sprintf("read request return body fail:%s", err.Error())},
		})
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 0 {
		err = errors.WithMessage(err, "return Status code err")
		return err
	}

	return nil
}

//######################合约停用上传###################################

//MakeStepGeneralDisableChainCodeZipUpload 合约打包文件上传
func MakeStepGeneralDisableChainCodeZipUpload(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create disable chainCode upload zip start"},
		})

		err := GeneralDisableChainCodeZipUpload(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDisableChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDisableChainCodeZipUpload",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create disable chainCode upload zip end"},
		})
		return nil
	}
}

//GeneralDisableChainCodeZipUpload 以后多链码扩展
func GeneralDisableChainCodeZipUpload(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralDisableChainCodeZipUploadFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDisableChainCodeZipUpload",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralDisableChainCodeZipUploadFile 合约停用上传
func GeneralDisableChainCodeZipUploadFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {

	for _, ipaddress := range ExecCCIPList {
		err := GeneralDisableChainCodeZipUploadStep(general, ipaddress, output)
		if err != nil {
			return err
		}
	}
	return nil
}

//GeneralDisableChainCodeZipUploadStep 分步上传
func GeneralDisableChainCodeZipUploadStep(general *objectdefine.Indent, ipAddress string, output *objectdefine.TaskNode) error {
	outputRoot := filepath.Join(general.BaseOutput, "disableChainCode")
	path := filepath.Join(outputRoot, ipAddress+".tar.gz")

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/file/upload", ipAddress, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalDisableChainCodeZipUpload",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	extraParams := map[string]string{
		"title":       "My Document",
		"author":      "Matt Aimonetti",
		"description": "A document with all the Go programming language secrets",
	}
	request, err := newfileUploadRequest(url, extraParams, "file", path, "", output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalDisableChainCodeZipUpload",
			Error: []string{fmt.Sprintf("fail upload request fail:%s", err.Error())},
		})
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalDisableChainCodeZipUpload",
			Error: []string{fmt.Sprintf("exec http post clent.DO fail:%s", err.Error())},
		})
		return err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalDisableChainCodeZipUpload",
			Error: []string{fmt.Sprintf("read request return body fail:%s", err.Error())},
		})
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 0 {
		err = errors.WithMessage(err, "return Status code err")
		return err
	}
	return nil
}

//######################合约启用上传###################################

//MakeStepGeneralEnableChainCodeZipUpload 合约打包文件上传
func MakeStepGeneralEnableChainCodeZipUpload(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create enable chainCode upload zip start"},
		})

		err := GeneralEnableChainCodeZipUpload(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailEnableChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalEnableChainCodeZipUpload",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create enable chainCode upload zip end"},
		})
		return nil
	}
}

//GeneralEnableChainCodeZipUpload 以后多链码扩展
func GeneralEnableChainCodeZipUpload(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralEnableChainCodeZipUploadFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalEnableChainCodeZipUpload",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralEnableChainCodeZipUploadFile 合约停用上传
func GeneralEnableChainCodeZipUploadFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {

	for _, ipaddress := range ExecCCIPList {
		err := GeneralEnableChainCodeZipUploadStep(general, ipaddress, output)
		if err != nil {
			return err
		}
	}
	return nil
}

//GeneralEnableChainCodeZipUploadStep 分步上传
func GeneralEnableChainCodeZipUploadStep(general *objectdefine.Indent, ipAddress string, output *objectdefine.TaskNode) error {
	outputRoot := filepath.Join(general.BaseOutput, "enableChainCode")
	path := filepath.Join(outputRoot, ipAddress+".tar.gz")

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/file/upload", ipAddress, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalEnableChainCodeZipUpload",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	extraParams := map[string]string{
		"title":       "My Document",
		"author":      "Matt Aimonetti",
		"description": "A document with all the Go programming language secrets",
	}
	request, err := newfileUploadRequest(url, extraParams, "file", path, "", output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalEnableChainCodeZipUpload",
			Error: []string{fmt.Sprintf("fail upload request fail:%s", err.Error())},
		})
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalEnableChainCodeZipUpload",
			Error: []string{fmt.Sprintf("exec http post clent.DO fail:%s", err.Error())},
		})
		return err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalEnableChainCodeZipUpload",
			Error: []string{fmt.Sprintf("read request return body fail:%s", err.Error())},
		})
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 0 {
		err = errors.WithMessage(err, "return Status code err")
		return err
	}
	return nil
}

//######################通道上传###################################

//MakeStepGeneralChannelZipUpload 通道打包文件上传
func MakeStepGeneralChannelZipUpload(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create channel upload zip start"},
		})

		err := GeneralCreateChannelZipUpload(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateChannelTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChannelZipUpload",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create channel upload zip end"},
		})
		return nil
	}
}

//GeneralCreateChannelZipUpload 本地文件上传远端服务 以后多组织扩展
func GeneralCreateChannelZipUpload(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreateChannelZipUploadFile(general, &peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChannelZipUpload",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralCreateChannelZipUploadFile 通道上传-本地文件上传远端服务
func GeneralCreateChannelZipUploadFile(general *objectdefine.Indent, peerOrder *objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("createChannel-%s", general.ChannelName)
	outputRoot := filepath.Join(general.BaseOutput, "createChannel")
	path := filepath.Join(outputRoot, folder+".tar.gz")
	//获取连接远端工具端口
	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/file/upload", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateChannelZipUpload",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	extraParams := map[string]string{
		"title":       "My Document",
		"author":      "Matt Aimonetti",
		"description": "A document with all the Go programming language secrets",
	}
	request, err := newfileUploadRequest(url, extraParams, "file", path, "", output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChannelZipUpload",
			Error: []string{fmt.Sprintf("fail upload request fail:%s", err.Error())},
		})
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChannelZipUpload",
			Error: []string{fmt.Sprintf("exec http post clent.DO fail:%s", err.Error())},
		})
		return err
	}
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChannelZipUpload",
			Error: []string{fmt.Sprintf("read request return body fail:%s", err.Error())},
		})
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 0 {
		err = errors.WithMessage(err, "return Status code err")
		return err
	}

	return nil
}

//######################组织解压###################################

//MakeStepGeneralOrgUnzip 解压上传文件
func MakeStepGeneralOrgUnzip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add org Unzip start"},
		})

		err := GeneralCreateOrgUnZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateOrgUnZip",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add org Unzip end"},
		})
		return nil
	}
}

//GeneralCreateOrgUnZip 以后多组织扩展
func GeneralCreateOrgUnZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreateOrgUnZipFile(general, org, peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateOrgUnZip",
					Error: []string{err.Error()},
				})
				return err
			}
		}
		//解压用来添加新组织操作Operate
		err := GeneralCreateOperateOrgUnZipFile(general, org, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateOrgUnZip",
				Error: []string{err.Error()},
			})
			return err
		}
	}
	return nil
}

//GeneralCreateOrgUnZipFile 执行解压命令
func GeneralCreateOrgUnZipFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("addOrg-%s", orgOrder.Name)

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateOrgUnZip",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	command := "unpack"
	buff := fmt.Sprintf("tar -xvf %s%s.tar.gz", folder, peerOrder.IP)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgUnZip",
			Error: []string{err.Error()},
		})
		return err
	}
	return nil
}

//GeneralCreateOperateOrgUnZipFile 执行解压命令
func GeneralCreateOperateOrgUnZipFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("addOrg-%s", orgOrder.Name)

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", OperateAddOrgIP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateOrgUnZip",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	command := "unpack"
	buff := fmt.Sprintf("tar -xvf %s%s.tar.gz", folder, OperateAddOrgIP)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateOrgUnZip",
			Error: []string{err.Error()},
		})
		return err
	}
	return nil
}

//######################删除组织解压###################################

//MakeStepGeneralDeleteOrgUnzip 解压上传文件
func MakeStepGeneralDeleteOrgUnzip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete org Unzip start"},
		})

		err := GeneralDeleteOrgUnZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDeleteOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeleteOrgUnZip",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete org Unzip end"},
		})
		return nil
	}
}

//GeneralDeleteOrgUnZip 以后多组织扩展
func GeneralDeleteOrgUnZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	peerAllIP := make(map[string]string, 0)
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			if _, ok := peerAllIP[peer.IP]; !ok {
				peerAllIP[peer.IP] = peer.Name
				err := GeneralDeleteOrgUnZipFile(general, org, peer, output)
				if err != nil {
					output.AppendLog(&objectdefine.StepHistory{
						Name:  "GeneralDeleteOrgUnZip",
						Error: []string{err.Error()},
					})
					return err
				}
			}
		}
		//解压用来添加新组织操作Operate
		err := GeneralDeleteOperateOrgUnZipFile(general, org, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "GeneralDeleteOrgUnZip",
				Error: []string{err.Error()},
			})
			return err
		}
	}
	return nil
}

//GeneralDeleteOrgUnZipFile 执行解压命令
func GeneralDeleteOrgUnZipFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("deleteOrg-%s", orgOrder.Name)

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "GeneralDeleteOrgUnZip",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	command := "unpack"
	buff := fmt.Sprintf("tar -xvf %s%s.tar.gz", folder, peerOrder.IP)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "GeneralDeleteOrgUnZip",
			Error: []string{err.Error()},
		})
		return err
	}
	return nil
}

//GeneralDeleteOperateOrgUnZipFile 执行解压命令
func GeneralDeleteOperateOrgUnZipFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("deleteOrg-%s", orgOrder.Name)

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", OperateAddOrgIP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "GeneralDeleteOrgUnZip",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	command := "unpack"
	buff := fmt.Sprintf("tar -xvf %s%s.tar.gz", folder, OperateAddOrgIP)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "GeneralDeleteOrgUnZip",
			Error: []string{err.Error()},
		})
		return err
	}
	return nil
}

//######################节点解压###################################

//MakeStepGeneralPeerUnzip 解压上传文件
func MakeStepGeneralPeerUnzip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add peer Unzip start"},
		})

		err := GeneralCreatePeerUnZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreatePeerTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreatePeerUnZip",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add peer Unzip end"},
		})
		return nil
	}
}

//GeneralCreatePeerUnZip 以后多组织扩展
func GeneralCreatePeerUnZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreatePeerUnZipFile(general, org, peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreatePeerUnZip",
					Error: []string{err.Error()},
				})
				return err
			}
		}

	}
	return nil
}

//GeneralCreatePeerUnZipFile 执行解压命令
func GeneralCreatePeerUnZipFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("addPeer-%s-%s", orgOrder.Name, peerOrder.Name)

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreatePeerUnZip",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	command := "unpack"
	buff := fmt.Sprintf("tar -xvf %s.tar.gz", folder)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerUnZip",
			Error: []string{err.Error()},
		})
		return err
	}
	return nil
}

//######################删除节点解压###################################

//MakeStepGeneralDeletePeerUnzip 解压上传文件
func MakeStepGeneralDeletePeerUnzip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete peer Unzip start"},
		})

		err := GeneralDeletePeerUnZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDeletePeerTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeletePeerUnZip",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete peer Unzip end"},
		})
		return nil
	}
}

//GeneralDeletePeerUnZip 以后多组织扩展
func GeneralDeletePeerUnZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralDeletePeerUnZipFile(general, org, peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeletePeerUnZip",
					Error: []string{err.Error()},
				})
				return err
			}
		}

	}
	return nil
}

//GeneralDeletePeerUnZipFile 执行解压命令
func GeneralDeletePeerUnZipFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("deletePeer-%s-%s", orgOrder.Name, peerOrder.Name)

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalDeletePeerUnZip",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	command := "unpack"
	buff := fmt.Sprintf("tar -xvf %s.tar.gz", folder)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalDeletePeerUnZip",
			Error: []string{err.Error()},
		})
		return err
	}
	return nil
}

//######################合约解压###################################

//MakeStepGeneralChainCodeUnzip  合约压缩包解压
func MakeStepGeneralChainCodeUnzip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add chaincode Unzip start"},
		})

		err := GeneralCreateChainCodeUnZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChainCodeUnZip",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add chaincode  Unzip end"},
		})
		return nil
	}
}

//GeneralCreateChainCodeUnZip 以后扩展
func GeneralCreateChainCodeUnZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralCreateChainCodeUnZipFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChainCodeUnZip",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralCreateChainCodeUnZipFile 执行解压命令
func GeneralCreateChainCodeUnZipFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("addChainCode-%s-%s", ccName, cc.Version)

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", SelectExecCCIP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateChainCodeUnZip",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	command := "unpack"
	buff := fmt.Sprintf("tar -xvf %s.tar.gz", folder)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChainCodeUnZip",
			Error: []string{err.Error()},
		})
		return err
	}
	return nil
}

//######################合约删除解压###################################

//MakeStepGeneralDeleteChainCodeUnzip  合约压缩包解压
func MakeStepGeneralDeleteChainCodeUnzip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete chaincode  Unzip start"},
		})

		err := GeneralDeleteChainCodeUnZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete chaincode  Unzip end"},
		})
		return nil
	}
}

//GeneralDeleteChainCodeUnZip 以后扩展
func GeneralDeleteChainCodeUnZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralDeleteChainCodeUnZipFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeleteChainCodeUnZip",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralDeleteChainCodeUnZipFile 执行解压命令
func GeneralDeleteChainCodeUnZipFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	for _, ipaddress := range ExecCCIPList {

		connectPort := dconfig.GetStringByKey("toolsPort")
		url := fmt.Sprintf("http://%s:%s/command/exec", ipaddress, connectPort)
		output.AppendLog(&objectdefine.StepHistory{
			Name: "generalDeleteChainCodeUnZip",
			Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
		})

		command := "unpack"
		buff := fmt.Sprintf("tar -xvf %s.tar.gz", ipaddress)
		args := []string{buff}
		err := MakeHTTPRemoteCmd(url, command, args, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalDeleteChainCodeUnZip",
				Error: []string{err.Error()},
			})
			return err
		}
	}

	return nil
}

//######################合约升级解压###################################

//MakeStepGeneralUpgradeChainCodeUnzip  合约压缩包解压
func MakeStepGeneralUpgradeChainCodeUnzip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create upgrade chaincode  Unzip start"},
		})

		err := GeneralUpgradeChainCodeUnZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailUpgradeChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalUpgradeChainCodeUnZip",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create upgrade chaincode  Unzip end"},
		})
		return nil
	}
}

//GeneralUpgradeChainCodeUnZip 以后扩展
func GeneralUpgradeChainCodeUnZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralUpgradeChainCodeUnZipFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalUpgradeChainCodeUnZip",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralUpgradeChainCodeUnZipFile 执行解压命令
func GeneralUpgradeChainCodeUnZipFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("upgradeChainCode-%s-%s", ccName, cc.Version)

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", SelectExecCCIP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalUpgradeChainCodeUnZip",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	command := "unpack"
	buff := fmt.Sprintf("tar -xvf %s.tar.gz", folder)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalUpgradeChainCodeUnZip",
			Error: []string{err.Error()},
		})
		return err
	}
	return nil
}

//######################合约停用解压###################################

//MakeStepGeneralDisableChainCodeUnzip  合约压缩包解压
func MakeStepGeneralDisableChainCodeUnzip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create disable chaincode  Unzip start"},
		})

		err := GeneralDisableChainCodeUnZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDisableChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDisableChainCodeUnZip",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create disable chaincode  Unzip end"},
		})
		return nil
	}
}

//GeneralDisableChainCodeUnZip 以后扩展
func GeneralDisableChainCodeUnZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralDisableChainCodeUnZipFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDisableChainCodeUnZip",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralDisableChainCodeUnZipFile 执行解压命令
func GeneralDisableChainCodeUnZipFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	for _, ipaddress := range ExecCCIPList {

		connectPort := dconfig.GetStringByKey("toolsPort")
		url := fmt.Sprintf("http://%s:%s/command/exec", ipaddress, connectPort)
		output.AppendLog(&objectdefine.StepHistory{
			Name: "generalDisableChainCodeUnZip",
			Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
		})

		command := "unpack"
		buff := fmt.Sprintf("tar -xvf %s.tar.gz", ipaddress)
		args := []string{buff}
		err := MakeHTTPRemoteCmd(url, command, args, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalDisableChainCodeUnZip",
				Error: []string{err.Error()},
			})
			return err
		}
	}

	return nil
}

//######################合约启用解压###################################

//MakeStepGeneralEnableChainCodeUnzip  合约压缩包解压
func MakeStepGeneralEnableChainCodeUnzip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create enable chaincode  Unzip start"},
		})

		err := GeneralEnableChainCodeUnZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailEnableChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalEnableChainCodeUnZip",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create enable chaincode  Unzip end"},
		})
		return nil
	}
}

//GeneralEnableChainCodeUnZip 以后扩展
func GeneralEnableChainCodeUnZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralEnableChainCodeUnZipFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalEnableChainCodeUnZip",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralEnableChainCodeUnZipFile 执行解压命令
func GeneralEnableChainCodeUnZipFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	for _, ipaddress := range ExecCCIPList {

		connectPort := dconfig.GetStringByKey("toolsPort")
		url := fmt.Sprintf("http://%s:%s/command/exec", ipaddress, connectPort)
		output.AppendLog(&objectdefine.StepHistory{
			Name: "generalEnableChainCodeUnZip",
			Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
		})

		command := "unpack"
		buff := fmt.Sprintf("tar -xvf %s.tar.gz", ipaddress)
		args := []string{buff}
		err := MakeHTTPRemoteCmd(url, command, args, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalEnableChainCodeUnZip",
				Error: []string{err.Error()},
			})
			return err
		}
	}

	return nil
}

//######################通道解压###################################

//MakeStepGeneralChannelUnzip 解压上传文件
func MakeStepGeneralChannelUnzip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create channel Unzip start"},
		})

		err := GeneralCreateChannelUnZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateChannelTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChannelUnZip",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create channel Unzip end"},
		})
		return nil
	}
}

//GeneralCreateChannelUnZip  远端压缩包解压  以后多组织扩展
func GeneralCreateChannelUnZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreateChannelUnZipFile(general, &peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChannelUnZip",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralCreateChannelUnZipFile 执行解压命令
func GeneralCreateChannelUnZipFile(general *objectdefine.Indent, peerOrder *objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("createChannel-%s", general.ChannelName)
	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", peerOrder.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateChannelUnZip",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})

	command := "unpack"
	buff := fmt.Sprintf("tar -xvf %s.tar.gz", folder)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChannelUnZip",
			Error: []string{err.Error()},
		})
		return err
	}
	return nil
}
