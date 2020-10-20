package daction

import (
	"archive/tar"
	"compress/gzip"
	"deploy-server/app/dmysql"
	"deploy-server/app/objectdefine"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

//###########新增组织##################

//MakeStepGeneralOrgZip 打包所有文件上传
func MakeStepGeneralOrgZip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add org zip start"},
		})

		err := GeneralCreateOrgZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateOrgZip",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add org zip end"},
		})
		return nil
	}
}

//GeneralCreateOrgZip 以后多组织扩展
func GeneralCreateOrgZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreateOrgZipFile(general, org, peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateOrgZip",
					Error: []string{err.Error()},
				})
				return err
			}
		}
		//打包需要更新块配置Operate
		GeneralCreateOperateOrgZipFile(general, org, output)
	}

	return nil
}

//GeneralCreateOrgZipFile 打包目录下所有文件
func GeneralCreateOrgZipFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("addOrg-%s", orgOrder.Name)
	base := filepath.Join(general.BaseOutput, "addOrg", peerOrder.IP)
	outputRoot := filepath.Join(general.BaseOutput, "addOrg")
	f, err := os.Create(filepath.Join(outputRoot, folder+peerOrder.IP+".tar.gz"))
	if err != nil {
		err = errors.WithMessage(err, "Create tar.gz file error")
		return err
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	err = filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if path == base {
			return nil
		}
		zipPut := strings.TrimPrefix(path, base)
		// not zip base folder
		zipPut = strings.TrimLeft(zipPut, "\\/")

		var link string
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			if link, err = os.Readlink(path); err != nil {
				return err
			}
		}

		hdr, err := tar.FileInfoHeader(info, link)
		if err != nil {
			err = errors.WithMessage(err, "Build File ["+zipPut+"] zip header error")
			return err
		}
		hdr.Name = filepath.ToSlash(zipPut)
		if runtime.GOOS == "windows" && info.IsDir() {
			hdr.Mode = hdr.Mode & 0777
		}

		err = tw.WriteHeader(hdr) //写入头文件信息
		if err != nil {
			err = errors.WithMessage(err, "Write File ["+zipPut+"] header into zip error")
			return err
		}

		if info.IsDir() {
			return nil
		}
		if !info.Mode().IsRegular() { //nothing more to do for non-regular
			return nil
		}

		fout, err := os.Open(path)
		if err != nil {
			err = errors.WithMessage(err, "Open File ["+zipPut+"] error")
			return err
		}
		_, err = io.CopyBuffer(tw, fout, nil)
		if err != nil {
			err = errors.WithMessage(err, "Write File ["+zipPut+"] into zip error")
			return err
		}
		fout.Close()

		err = tw.Flush()
		if err != nil {
			err = errors.WithMessage(err, "Flush tar file by per file error")
			return err
		}
		return nil
	})

	if err != nil {
		err = errors.WithMessage(err, "Zipping files error")
		tw.Close()
		gw.Close()
		return err
	}

	err = tw.Close()
	if err != nil {
		err = errors.WithMessage(err, "Close tar file error")
		gw.Close()
		return err
	}

	err = gw.Close()
	if err != nil {
		err = errors.WithMessage(err, "Close tar.gz file error")
		return err
	}
	return nil
}

//GeneralCreateOperateOrgZipFile 打包需要更新块配置操作
func GeneralCreateOperateOrgZipFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("addOrg-%s", orgOrder.Name)
	base := filepath.Join(general.BaseOutput, "addOrg", OperateAddOrgIP)
	outputRoot := filepath.Join(general.BaseOutput, "addOrg")
	_, err := os.Stat(filepath.Join(outputRoot, folder+OperateAddOrgIP+".tar.gz"))
	if err == nil {
		return nil
	}
	f, err := os.Create(filepath.Join(outputRoot, folder+OperateAddOrgIP+".tar.gz"))
	if err != nil {
		err = errors.WithMessage(err, "Create tar.gz file error")
		return err
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	err = filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if path == base {
			return nil
		}
		zipPut := strings.TrimPrefix(path, base)
		// not zip base folder
		zipPut = strings.TrimLeft(zipPut, "\\/")

		var link string
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			if link, err = os.Readlink(path); err != nil {
				return err
			}
		}

		hdr, err := tar.FileInfoHeader(info, link)
		if err != nil {
			err = errors.WithMessage(err, "Build File ["+zipPut+"] zip header error")
			return err
		}
		hdr.Name = filepath.ToSlash(zipPut)
		if runtime.GOOS == "windows" && info.IsDir() {
			hdr.Mode = hdr.Mode & 0777
		}

		err = tw.WriteHeader(hdr) //写入头文件信息
		if err != nil {
			err = errors.WithMessage(err, "Write File ["+zipPut+"] header into zip error")
			return err
		}

		if info.IsDir() {
			return nil
		}
		if !info.Mode().IsRegular() { //nothing more to do for non-regular
			return nil
		}
		fout, err := os.Open(path)
		if err != nil {
			err = errors.WithMessage(err, "Open File ["+zipPut+"] error")
			return err
		}
		_, err = io.CopyBuffer(tw, fout, nil)
		if err != nil {
			err = errors.WithMessage(err, "Write File ["+zipPut+"] into zip error")
			return err
		}
		fout.Close()

		err = tw.Flush()
		if err != nil {
			err = errors.WithMessage(err, "Flush tar file by per file error")
			return err
		}
		return nil
	})

	if err != nil {
		err = errors.WithMessage(err, "Zipping files error")
		tw.Close()
		gw.Close()
		return err
	}

	err = tw.Close()
	if err != nil {
		err = errors.WithMessage(err, "Close tar file error")
		gw.Close()
		return err
	}

	err = gw.Close()
	if err != nil {
		err = errors.WithMessage(err, "Close tar.gz file error")
		return err
	}
	return nil
}

//###########新增节点##################

//MakeStepGeneralPeerZip 打包所有文件上传
func MakeStepGeneralPeerZip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add peer zip start"},
		})

		err := GeneralCreatePeerZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreatePeerTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreatePeerZip",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add peer zip end"},
		})
		return nil
	}
}

//GeneralCreatePeerZip 以后多组织扩展
func GeneralCreatePeerZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreatePeerZipFile(general, org, peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreatePeerZip",
					Error: []string{err.Error()},
				})
				return err
			}
		}

	}
	return nil
}

//GeneralCreatePeerZipFile 打包目录下所有文件
func GeneralCreatePeerZipFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("addPeer-%s-%s", orgOrder.Name, peerOrder.Name)
	base := filepath.Join(general.BaseOutput, "addPeer", peerOrder.IP)
	outputRoot := filepath.Join(general.BaseOutput, "addPeer")
	f, err := os.Create(filepath.Join(outputRoot, folder+".tar.gz"))
	if err != nil {
		err = errors.WithMessage(err, "Create tar.gz file error")
		return err
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	err = filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		// not zip base folder
		if path == base {
			return nil
		}

		zipPut := strings.TrimPrefix(path, base)
		// not zip base folder
		zipPut = strings.TrimLeft(zipPut, "\\/")

		var link string
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			if link, err = os.Readlink(path); err != nil {
				return err
			}
		}
		hdr, err := tar.FileInfoHeader(info, link)
		if err != nil {
			err = errors.WithMessage(err, "Build File ["+zipPut+"] zip header error")
			return err
		}
		hdr.Name = filepath.ToSlash(zipPut)
		if runtime.GOOS == "windows" && info.IsDir() {
			hdr.Mode = hdr.Mode & 0777
		}

		err = tw.WriteHeader(hdr) //写入头文件信息
		if err != nil {
			err = errors.WithMessage(err, "Write File ["+zipPut+"] header into zip error")
			return err
		}

		if info.IsDir() {
			return nil
		}
		if !info.Mode().IsRegular() { //nothing more to do for non-regular
			return nil
		}
		fout, err := os.Open(path)
		if err != nil {
			err = errors.WithMessage(err, "Open File ["+zipPut+"] error")
			return err
		}
		_, err = io.CopyBuffer(tw, fout, nil)
		if err != nil {
			err = errors.WithMessage(err, "Write File ["+zipPut+"] into zip error")
			return err
		}
		fout.Close()

		err = tw.Flush()
		if err != nil {
			err = errors.WithMessage(err, "Flush tar file by per file error")
			return err
		}
		return nil
	})

	if err != nil {
		err = errors.WithMessage(err, "Zipping files error")
		tw.Close()
		gw.Close()
		return err
	}

	err = tw.Close()
	if err != nil {
		err = errors.WithMessage(err, "Close tar file error")
		gw.Close()
		return err
	}

	err = gw.Close()
	if err != nil {
		err = errors.WithMessage(err, "Close tar.gz file error")
		return err
	}
	return nil
}

//###########新增合约##################

//MakeStepGeneralChainCodeZip 创建合约压缩包
func MakeStepGeneralChainCodeZip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add chaincode zip start"},
		})

		err := GeneralCreateChainCodeZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChainCodeZip",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create add chaincode zip end"},
		})
		return nil
	}
}

//GeneralCreateChainCodeZip 以后多组织扩展
func GeneralCreateChainCodeZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralCreateChainCodeZipFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChainCodeZip",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralCreateChainCodeZipFile 打包目录下所有文件
func GeneralCreateChainCodeZipFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("addChainCode-%s-%s", ccName, cc.Version)
	base := filepath.Join(general.BaseOutput, "addChainCode", SelectExecCCIP)
	outputRoot := filepath.Join(general.BaseOutput, "addChainCode")
	f, err := os.Create(filepath.Join(outputRoot, folder+".tar.gz"))
	if err != nil {
		err = errors.WithMessage(err, "Create tar.gz file error")
		return err
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	err = filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		// not zip base folder
		if path == base {
			return nil
		}

		zipPut := strings.TrimPrefix(path, base)
		// not zip base folder
		zipPut = strings.TrimLeft(zipPut, "\\/")

		var link string
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			if link, err = os.Readlink(path); err != nil {
				return err
			}
		}
		hdr, err := tar.FileInfoHeader(info, link)
		if err != nil {
			err = errors.WithMessage(err, "Build File ["+zipPut+"] zip header error")
			return err
		}
		hdr.Name = filepath.ToSlash(zipPut)
		if runtime.GOOS == "windows" && info.IsDir() {
			hdr.Mode = hdr.Mode & 0777
		}

		err = tw.WriteHeader(hdr) //写入头文件信息
		if err != nil {
			err = errors.WithMessage(err, "Write File ["+zipPut+"] header into zip error")
			return err
		}

		if info.IsDir() {
			return nil
		}
		if !info.Mode().IsRegular() { //nothing more to do for non-regular
			return nil
		}
		fout, err := os.Open(path)
		if err != nil {
			err = errors.WithMessage(err, "Open File ["+zipPut+"] error")
			return err
		}
		_, err = io.CopyBuffer(tw, fout, nil)
		if err != nil {
			err = errors.WithMessage(err, "Write File ["+zipPut+"] into zip error")
			return err
		}
		fout.Close()

		err = tw.Flush()
		if err != nil {
			err = errors.WithMessage(err, "Flush tar file by per file error")
			return err
		}
		return nil
	})

	if err != nil {
		err = errors.WithMessage(err, "Zipping files error")
		tw.Close()
		gw.Close()
		return err
	}

	err = tw.Close()
	if err != nil {
		err = errors.WithMessage(err, "Close tar file error")
		gw.Close()
		return err
	}

	err = gw.Close()
	if err != nil {
		err = errors.WithMessage(err, "Close tar.gz file error")
		return err
	}
	return nil

}

//###########升级合约##################

//MakeStepGeneralUpgradeChainCodeZip 创建组织压缩包
func MakeStepGeneralUpgradeChainCodeZip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create upgrade chaincode zip start"},
		})

		err := GeneralUpgradeChainCodeZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailUpgradeChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateUpgradeChainCodeZip",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create upgrade chaincode zip end"},
		})
		return nil
	}
}

//GeneralUpgradeChainCodeZip 以后多组织扩展
func GeneralUpgradeChainCodeZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralUpgradeChainCodeZipFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateUpgradeChainCodeZip",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralUpgradeChainCodeZipFile 打包目录下所有文件
func GeneralUpgradeChainCodeZipFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("upgradeChainCode-%s-%s", ccName, cc.Version)
	base := filepath.Join(general.BaseOutput, "upgradeChainCode", SelectExecCCIP)
	outputRoot := filepath.Join(general.BaseOutput, "upgradeChainCode")
	f, err := os.Create(filepath.Join(outputRoot, folder+".tar.gz"))
	if err != nil {
		err = errors.WithMessage(err, "Create tar.gz file error")
		return err
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	err = filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		// not zip base folder
		if path == base {
			return nil
		}

		zipPut := strings.TrimPrefix(path, base)
		// not zip base folder
		zipPut = strings.TrimLeft(zipPut, "\\/")

		var link string
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			if link, err = os.Readlink(path); err != nil {
				return err
			}
		}
		hdr, err := tar.FileInfoHeader(info, link)
		if err != nil {
			err = errors.WithMessage(err, "Build File ["+zipPut+"] zip header error")
			return err
		}
		hdr.Name = filepath.ToSlash(zipPut)
		if runtime.GOOS == "windows" && info.IsDir() {
			hdr.Mode = hdr.Mode & 0777
		}

		err = tw.WriteHeader(hdr) //写入头文件信息
		if err != nil {
			err = errors.WithMessage(err, "Write File ["+zipPut+"] header into zip error")
			return err
		}

		if info.IsDir() {
			return nil
		}
		if !info.Mode().IsRegular() { //nothing more to do for non-regular
			return nil
		}
		fout, err := os.Open(path)
		if err != nil {
			err = errors.WithMessage(err, "Open File ["+zipPut+"] error")
			return err
		}
		_, err = io.CopyBuffer(tw, fout, nil)
		if err != nil {
			err = errors.WithMessage(err, "Write File ["+zipPut+"] into zip error")
			return err
		}
		fout.Close()

		err = tw.Flush()
		if err != nil {
			err = errors.WithMessage(err, "Flush tar file by per file error")
			return err
		}
		return nil
	})

	if err != nil {
		err = errors.WithMessage(err, "Zipping files error")
		tw.Close()
		gw.Close()
		return err
	}

	err = tw.Close()
	if err != nil {
		err = errors.WithMessage(err, "Close tar file error")
		gw.Close()
		return err
	}

	err = gw.Close()
	if err != nil {
		err = errors.WithMessage(err, "Close tar.gz file error")
		return err
	}
	return nil

}

//###########合约停用##################

//MakeStepGeneralDisableChainCodeZip 创建合约压缩包
func MakeStepGeneralDisableChainCodeZip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create disable chaincode zip start"},
		})

		err := GeneralDisableChainCodeZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDisableChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateDisableChainCodeZip",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create disable chaincode zip end"},
		})
		return nil
	}
}

//GeneralDisableChainCodeZip 以后多组织扩展
func GeneralDisableChainCodeZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralDisableChainCodeZipFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateDisableChainCodeZip",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralDisableChainCodeZipFile 打包目录下所有文件
func GeneralDisableChainCodeZipFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	//folder := fmt.Sprintf("disableChainCode-%s-%s", ccName, cc.Version)
	for _, ipaddress := range ExecCCIPList {
		base := filepath.Join(general.BaseOutput, "disableChainCode", ipaddress)
		outputRoot := filepath.Join(general.BaseOutput, "disableChainCode")
		f, err := os.Create(filepath.Join(outputRoot, ipaddress+".tar.gz"))
		if err != nil {
			err = errors.WithMessage(err, "Create tar.gz file error")
			return err
		}
		defer f.Close()
		gw := gzip.NewWriter(f)
		tw := tar.NewWriter(gw)
		err = filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
			// not zip base folder
			if path == base {
				return nil
			}

			zipPut := strings.TrimPrefix(path, base)
			// not zip base folder
			zipPut = strings.TrimLeft(zipPut, "\\/")

			var link string
			if info.Mode()&os.ModeSymlink == os.ModeSymlink {
				if link, err = os.Readlink(path); err != nil {
					return err
				}
			}
			hdr, err := tar.FileInfoHeader(info, link)
			if err != nil {
				err = errors.WithMessage(err, "Build File ["+zipPut+"] zip header error")
				return err
			}
			hdr.Name = filepath.ToSlash(zipPut)
			if runtime.GOOS == "windows" && info.IsDir() {
				hdr.Mode = hdr.Mode & 0777
			}

			err = tw.WriteHeader(hdr) //写入头文件信息
			if err != nil {
				err = errors.WithMessage(err, "Write File ["+zipPut+"] header into zip error")
				return err
			}

			if info.IsDir() {
				return nil
			}
			if !info.Mode().IsRegular() { //nothing more to do for non-regular
				return nil
			}

			fout, err := os.Open(path)
			if err != nil {
				err = errors.WithMessage(err, "Open File ["+zipPut+"] error")
				return err
			}
			_, err = io.CopyBuffer(tw, fout, nil)
			if err != nil {
				err = errors.WithMessage(err, "Write File ["+zipPut+"] into zip error")
				return err
			}
			fout.Close()

			err = tw.Flush()
			if err != nil {
				err = errors.WithMessage(err, "Flush tar file by per file error")
				return err
			}
			return nil
		})

		if err != nil {
			err = errors.WithMessage(err, "Zipping files error")
			tw.Close()
			gw.Close()
			return err
		}

		err = tw.Close()
		if err != nil {
			err = errors.WithMessage(err, "Close tar file error")
			gw.Close()
			return err
		}

		err = gw.Close()
		if err != nil {
			err = errors.WithMessage(err, "Close tar.gz file error")
			return err
		}
	}

	return nil

}

//###########合约启用##################

//MakeStepGeneralEnableChainCodeZip 创建合约压缩包
func MakeStepGeneralEnableChainCodeZip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create enable chaincode zip start"},
		})

		err := GeneralEnenableenableableChainCodeZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailEnableChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateEnableChainCodeZip",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create enable chaincode zip end"},
		})
		return nil
	}
}

//GeneralEnenableenableableChainCodeZip 以后多组织扩展
func GeneralEnenableenableableChainCodeZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralEnableChainCodeZipFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateEnableChainCodeZip",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralEnableChainCodeZipFile 打包目录下所有文件
func GeneralEnableChainCodeZipFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	//folder := fmt.Sprintf("disableChainCode-%s-%s", ccName, cc.Version)
	for _, ipaddress := range ExecCCIPList {
		base := filepath.Join(general.BaseOutput, "enableChainCode", ipaddress)
		outputRoot := filepath.Join(general.BaseOutput, "enableChainCode")
		f, err := os.Create(filepath.Join(outputRoot, ipaddress+".tar.gz"))
		if err != nil {
			err = errors.WithMessage(err, "Create tar.gz file error")
			return err
		}
		defer f.Close()
		gw := gzip.NewWriter(f)
		tw := tar.NewWriter(gw)
		err = filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
			// not zip base folder
			if path == base {
				return nil
			}

			zipPut := strings.TrimPrefix(path, base)
			// not zip base folder
			zipPut = strings.TrimLeft(zipPut, "\\/")

			var link string
			if info.Mode()&os.ModeSymlink == os.ModeSymlink {
				if link, err = os.Readlink(path); err != nil {
					return err
				}
			}
			hdr, err := tar.FileInfoHeader(info, link)
			if err != nil {
				err = errors.WithMessage(err, "Build File ["+zipPut+"] zip header error")
				return err
			}
			hdr.Name = filepath.ToSlash(zipPut)
			if runtime.GOOS == "windows" && info.IsDir() {
				hdr.Mode = hdr.Mode & 0777
			}

			err = tw.WriteHeader(hdr) //写入头文件信息
			if err != nil {
				err = errors.WithMessage(err, "Write File ["+zipPut+"] header into zip error")
				return err
			}

			if info.IsDir() {
				return nil
			}
			if !info.Mode().IsRegular() { //nothing more to do for non-regular
				return nil
			}
			fout, err := os.Open(path)
			if err != nil {
				err = errors.WithMessage(err, "Open File ["+zipPut+"] error")
				return err
			}
			_, err = io.CopyBuffer(tw, fout, nil)
			if err != nil {
				err = errors.WithMessage(err, "Write File ["+zipPut+"] into zip error")
				return err
			}
			fout.Close()

			err = tw.Flush()
			if err != nil {
				err = errors.WithMessage(err, "Flush tar file by per file error")
				return err
			}
			return nil
		})

		if err != nil {
			err = errors.WithMessage(err, "Zipping files error")
			tw.Close()
			gw.Close()
			return err
		}

		err = tw.Close()
		if err != nil {
			err = errors.WithMessage(err, "Close tar file error")
			gw.Close()
			return err
		}

		err = gw.Close()
		if err != nil {
			err = errors.WithMessage(err, "Close tar.gz file error")
			return err
		}
	}

	return nil

}

//###########创建通道##################

//MakeStepGeneralChannelZip 创建压缩包
func MakeStepGeneralChannelZip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create  channel zip start"},
		})

		err := GeneralChannelZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateChannelTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChannelZip",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create channel zip end"},
		})
		return nil
	}
}

//GeneralChannelZip  本地文件打包 以后多组织扩展
func GeneralChannelZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralCreateChannelZipFile(general, &peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChannelZip",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralCreateChannelZipFile 打包目录下所有文件
func GeneralCreateChannelZipFile(general *objectdefine.Indent, peerOrder *objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("createChannel-%s", general.ChannelName)
	outputRoot := filepath.Join(general.BaseOutput, "createChannel")
	base := filepath.Join(general.BaseOutput, "createChannel", peerOrder.IP)
	f, err := os.Create(filepath.Join(outputRoot, folder+".tar.gz"))
	if err != nil {
		err = errors.WithMessage(err, "Create tar.gz file error")
		return err
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	err = filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		// not zip base folder
		if path == base {
			return nil
		}

		zipPut := strings.TrimPrefix(path, base)
		// not zip base folder
		zipPut = strings.TrimLeft(zipPut, "\\/")

		var link string
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			if link, err = os.Readlink(path); err != nil {
				return err
			}
		}
		hdr, err := tar.FileInfoHeader(info, link)
		if err != nil {
			err = errors.WithMessage(err, "Build File ["+zipPut+"] zip header error")
			return err
		}
		hdr.Name = filepath.ToSlash(zipPut)
		if runtime.GOOS == "windows" && info.IsDir() {
			hdr.Mode = hdr.Mode & 0777
		}

		err = tw.WriteHeader(hdr) //写入头文件信息
		if err != nil {
			err = errors.WithMessage(err, "Write File ["+zipPut+"] header into zip error")
			return err
		}

		if info.IsDir() {
			return nil
		}
		if !info.Mode().IsRegular() { //nothing more to do for non-regular
			return nil
		}
		fout, err := os.Open(path)
		if err != nil {
			err = errors.WithMessage(err, "Open File ["+zipPut+"] error")
			return err
		}
		_, err = io.CopyBuffer(tw, fout, nil)
		if err != nil {
			err = errors.WithMessage(err, "Write File ["+zipPut+"] into zip error")
			return err
		}
		fout.Close()

		err = tw.Flush()
		if err != nil {
			err = errors.WithMessage(err, "Flush tar file by per file error")
			return err
		}
		return nil
	})

	if err != nil {
		err = errors.WithMessage(err, "Zipping files error")
		tw.Close()
		gw.Close()
		return err
	}

	err = tw.Close()
	if err != nil {
		err = errors.WithMessage(err, "Close tar file error")
		gw.Close()
		return err
	}

	err = gw.Close()
	if err != nil {
		err = errors.WithMessage(err, "Close tar.gz file error")
		return err
	}
	return nil

}

//###########删除组织##################

//MakeStepGeneralDeleteOrgZip 打包所有文件上传
func MakeStepGeneralDeleteOrgZip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete org zip start"},
		})

		err := GeneralDeleteOrgZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDeleteOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeleteOrgZip",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete org zip end"},
		})
		return nil
	}
}

//GeneralDeleteOrgZip 以后多组织扩展
func GeneralDeleteOrgZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	peerAllIP := make(map[string]string, 0)
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			if _, ok := peerAllIP[peer.IP]; !ok {
				peerAllIP[peer.IP] = peer.Name
				err := GeneralDeleteOrgZipFile(general, org, peer, output)
				if err != nil {
					output.AppendLog(&objectdefine.StepHistory{
						Name:  "generalDeleteOrgZip",
						Error: []string{err.Error()},
					})
					return err
				}
			}
		}
		//打包需要更新块配置Operate
		GeneralDeleteOperateOrgZipFile(general, org, output)
	}

	return nil
}

//GeneralDeleteOrgZipFile 打包目录下所有文件
func GeneralDeleteOrgZipFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("deleteOrg-%s", orgOrder.Name)
	base := filepath.Join(general.BaseOutput, "deleteOrg", peerOrder.IP)
	outputRoot := filepath.Join(general.BaseOutput, "deleteOrg")
	f, err := os.Create(filepath.Join(outputRoot, folder+peerOrder.IP+".tar.gz"))
	if err != nil {
		err = errors.WithMessage(err, "Create tar.gz file error")
		return err
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	err = filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if path == base {
			return nil
		}
		zipPut := strings.TrimPrefix(path, base)
		// not zip base folder
		zipPut = strings.TrimLeft(zipPut, "\\/")

		var link string
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			if link, err = os.Readlink(path); err != nil {
				return err
			}
		}

		hdr, err := tar.FileInfoHeader(info, link)
		if err != nil {
			err = errors.WithMessage(err, "Build File ["+zipPut+"] zip header error")
			return err
		}
		hdr.Name = filepath.ToSlash(zipPut)
		if runtime.GOOS == "windows" && info.IsDir() {
			hdr.Mode = hdr.Mode & 0777
		}

		err = tw.WriteHeader(hdr) //写入头文件信息
		if err != nil {
			err = errors.WithMessage(err, "Write File ["+zipPut+"] header into zip error")
			return err
		}

		if info.IsDir() {
			return nil
		}
		if !info.Mode().IsRegular() { //nothing more to do for non-regular
			return nil
		}

		fout, err := os.Open(path)
		if err != nil {
			err = errors.WithMessage(err, "Open File ["+zipPut+"] error")
			return err
		}
		_, err = io.CopyBuffer(tw, fout, nil)
		if err != nil {
			err = errors.WithMessage(err, "Write File ["+zipPut+"] into zip error")
			return err
		}
		fout.Close()

		err = tw.Flush()
		if err != nil {
			err = errors.WithMessage(err, "Flush tar file by per file error")
			return err
		}
		return nil
	})

	if err != nil {
		err = errors.WithMessage(err, "Zipping files error")
		tw.Close()
		gw.Close()
		return err
	}

	err = tw.Close()
	if err != nil {
		err = errors.WithMessage(err, "Close tar file error")
		gw.Close()
		return err
	}

	err = gw.Close()
	if err != nil {
		err = errors.WithMessage(err, "Close tar.gz file error")
		return err
	}
	return nil
}

//GeneralDeleteOperateOrgZipFile 打包需要更新块配置操作
func GeneralDeleteOperateOrgZipFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("deleteOrg-%s", orgOrder.Name)
	base := filepath.Join(general.BaseOutput, "deleteOrg", OperateAddOrgIP)
	outputRoot := filepath.Join(general.BaseOutput, "deleteOrg")
	_, err := os.Stat(filepath.Join(outputRoot, folder+OperateAddOrgIP+".tar.gz"))
	if err == nil {
		return nil
	}
	f, err := os.Create(filepath.Join(outputRoot, folder+OperateAddOrgIP+".tar.gz"))
	if err != nil {
		err = errors.WithMessage(err, "Create tar.gz file error")
		return err
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	err = filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if path == base {
			return nil
		}
		zipPut := strings.TrimPrefix(path, base)
		// not zip base folder
		zipPut = strings.TrimLeft(zipPut, "\\/")

		var link string
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			if link, err = os.Readlink(path); err != nil {
				return err
			}
		}

		hdr, err := tar.FileInfoHeader(info, link)
		if err != nil {
			err = errors.WithMessage(err, "Build File ["+zipPut+"] zip header error")
			return err
		}
		hdr.Name = filepath.ToSlash(zipPut)
		if runtime.GOOS == "windows" && info.IsDir() {
			hdr.Mode = hdr.Mode & 0777
		}

		err = tw.WriteHeader(hdr) //写入头文件信息
		if err != nil {
			err = errors.WithMessage(err, "Write File ["+zipPut+"] header into zip error")
			return err
		}

		if info.IsDir() {
			return nil
		}
		if !info.Mode().IsRegular() { //nothing more to do for non-regular
			return nil
		}
		fout, err := os.Open(path)
		if err != nil {
			err = errors.WithMessage(err, "Open File ["+zipPut+"] error")
			return err
		}
		_, err = io.CopyBuffer(tw, fout, nil)
		if err != nil {
			err = errors.WithMessage(err, "Write File ["+zipPut+"] into zip error")
			return err
		}
		fout.Close()

		err = tw.Flush()
		if err != nil {
			err = errors.WithMessage(err, "Flush tar file by per file error")
			return err
		}
		return nil
	})

	if err != nil {
		err = errors.WithMessage(err, "Zipping files error")
		tw.Close()
		gw.Close()
		return err
	}

	err = tw.Close()
	if err != nil {
		err = errors.WithMessage(err, "Close tar file error")
		gw.Close()
		return err
	}

	err = gw.Close()
	if err != nil {
		err = errors.WithMessage(err, "Close tar.gz file error")
		return err
	}
	return nil
}

//###########删除节点##################

//MakeStepGeneralDeletePeerZip 打包所有文件上传
func MakeStepGeneralDeletePeerZip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete peer zip start"},
		})

		err := GeneralDeletePeerZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDeletePeerTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeletePeerZip",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete peer zip end"},
		})
		return nil
	}
}

//GeneralDeletePeerZip 以后多组织扩展
func GeneralDeletePeerZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralDeletePeerZipFile(general, org, peer, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeletePeerZip",
					Error: []string{err.Error()},
				})
				return err
			}
		}

	}
	return nil
}

//GeneralDeletePeerZipFile 打包目录下所有文件
func GeneralDeletePeerZipFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	folder := fmt.Sprintf("deletePeer-%s-%s", orgOrder.Name, peerOrder.Name)
	base := filepath.Join(general.BaseOutput, "deletePeer", peerOrder.IP)
	outputRoot := filepath.Join(general.BaseOutput, "deletePeer")
	f, err := os.Create(filepath.Join(outputRoot, folder+".tar.gz"))
	if err != nil {
		err = errors.WithMessage(err, "Create tar.gz file error")
		return err
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	err = filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		// not zip base folder
		if path == base {
			return nil
		}

		zipPut := strings.TrimPrefix(path, base)
		// not zip base folder
		zipPut = strings.TrimLeft(zipPut, "\\/")

		var link string
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			if link, err = os.Readlink(path); err != nil {
				return err
			}
		}
		hdr, err := tar.FileInfoHeader(info, link)
		if err != nil {
			err = errors.WithMessage(err, "Build File ["+zipPut+"] zip header error")
			return err
		}
		hdr.Name = filepath.ToSlash(zipPut)
		if runtime.GOOS == "windows" && info.IsDir() {
			hdr.Mode = hdr.Mode & 0777
		}

		err = tw.WriteHeader(hdr) //写入头文件信息
		if err != nil {
			err = errors.WithMessage(err, "Write File ["+zipPut+"] header into zip error")
			return err
		}

		if info.IsDir() {
			return nil
		}
		if !info.Mode().IsRegular() { //nothing more to do for non-regular
			return nil
		}
		fout, err := os.Open(path)
		if err != nil {
			err = errors.WithMessage(err, "Open File ["+zipPut+"] error")
			return err
		}
		_, err = io.CopyBuffer(tw, fout, nil)
		if err != nil {
			err = errors.WithMessage(err, "Write File ["+zipPut+"] into zip error")
			return err
		}
		fout.Close()

		err = tw.Flush()
		if err != nil {
			err = errors.WithMessage(err, "Flush tar file by per file error")
			return err
		}
		return nil
	})

	if err != nil {
		err = errors.WithMessage(err, "Zipping files error")
		tw.Close()
		gw.Close()
		return err
	}

	err = tw.Close()
	if err != nil {
		err = errors.WithMessage(err, "Close tar file error")
		gw.Close()
		return err
	}

	err = gw.Close()
	if err != nil {
		err = errors.WithMessage(err, "Close tar.gz file error")
		return err
	}
	return nil
}

//###########删除合约#################

//MakeStepGeneralDeleteChainCodeZip 创建合约压缩包
func MakeStepGeneralDeleteChainCodeZip(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete chaincode zip start"},
		})

		err := GeneralDeleteChainCodeZip(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create delete chaincode zip end"},
		})
		return nil
	}
}

//GeneralDeleteChainCodeZip 以后多组织扩展
func GeneralDeleteChainCodeZip(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for ccName, ccA := range general.Chaincode {
		for _, cc := range ccA {
			err := GeneralDeleteChainCodeZipFile(general, &cc, ccName, output)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateDeleteChainCodeZip",
					Error: []string{err.Error()},
				})
				return err
			}
		}
	}
	return nil
}

//GeneralDeleteChainCodeZipFile 打包目录下所有文件
func GeneralDeleteChainCodeZipFile(general *objectdefine.Indent, cc *objectdefine.ChainCodeType, ccName string, output *objectdefine.TaskNode) error {
	//folder := fmt.Sprintf("disableChainCode-%s-%s", ccName, cc.Version)
	for _, ipaddress := range ExecCCIPList {
		base := filepath.Join(general.BaseOutput, "deleteChainCode", ipaddress)
		outputRoot := filepath.Join(general.BaseOutput, "deleteChainCode")
		f, err := os.Create(filepath.Join(outputRoot, ipaddress+".tar.gz"))
		if err != nil {
			err = errors.WithMessage(err, "Create tar.gz file error")
			return err
		}
		defer f.Close()
		gw := gzip.NewWriter(f)
		tw := tar.NewWriter(gw)
		err = filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
			// not zip base folder
			if path == base {
				return nil
			}

			zipPut := strings.TrimPrefix(path, base)
			// not zip base folder
			zipPut = strings.TrimLeft(zipPut, "\\/")

			var link string
			if info.Mode()&os.ModeSymlink == os.ModeSymlink {
				if link, err = os.Readlink(path); err != nil {
					return err
				}
			}
			hdr, err := tar.FileInfoHeader(info, link)
			if err != nil {
				err = errors.WithMessage(err, "Build File ["+zipPut+"] zip header error")
				return err
			}
			hdr.Name = filepath.ToSlash(zipPut)
			if runtime.GOOS == "windows" && info.IsDir() {
				hdr.Mode = hdr.Mode & 0777
			}

			err = tw.WriteHeader(hdr) //写入头文件信息
			if err != nil {
				err = errors.WithMessage(err, "Write File ["+zipPut+"] header into zip error")
				return err
			}

			if info.IsDir() {
				return nil
			}
			if !info.Mode().IsRegular() { //nothing more to do for non-regular
				return nil
			}

			fout, err := os.Open(path)
			if err != nil {
				err = errors.WithMessage(err, "Open File ["+zipPut+"] error")
				return err
			}
			_, err = io.CopyBuffer(tw, fout, nil)
			if err != nil {
				err = errors.WithMessage(err, "Write File ["+zipPut+"] into zip error")
				return err
			}
			fout.Close()

			err = tw.Flush()
			if err != nil {
				err = errors.WithMessage(err, "Flush tar file by per file error")
				return err
			}
			return nil
		})

		if err != nil {
			err = errors.WithMessage(err, "Zipping files error")
			tw.Close()
			gw.Close()
			return err
		}

		err = tw.Close()
		if err != nil {
			err = errors.WithMessage(err, "Close tar file error")
			gw.Close()
			return err
		}

		err = gw.Close()
		if err != nil {
			err = errors.WithMessage(err, "Close tar.gz file error")
			return err
		}
	}

	return nil

}
