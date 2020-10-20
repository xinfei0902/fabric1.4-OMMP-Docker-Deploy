package dcache

import (
	"deploy-server/app/objectdefine"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//GlobalOneLocal 获取本地loacl.json
func GlobalOneLocal() *localCacheV1 {
	return globalLocalObject
}

//GetBinToolsPair 获取生成证书 创世区块工具
func GetBinToolsPair() map[string]objectdefine.BinToolsType {
	core := GlobalOneLocal()
	return core.local.Bin
}

//GetStoreRoot 获取local.json 中storeRoot字段
func GetStoreRoot() string {
	core := GlobalOneLocal()
	return core.local.StoreRoot
}

//GetDefaultVersion 获取local.json默认版本
func GetDefaultVersion() string {
	core := GlobalOneLocal()
	if len(core.local.Versions) == 0 {
		return ""
	}
	return core.local.Versions[0].Version
}

//GetOutputSubPath 获取输出路径
func GetOutputSubPath(id, sub string) string {
	base := GetStoreRoot()
	if len(sub) == 0 {
		return filepath.Join(base, id)
	}
	return filepath.Join(base, id, sub)
}

//GetVersionRootPathByVersion 获取路径
func GetVersionRootPathByVersion(version string) string {
	core := GlobalOneLocal()
	l := len(core.local.Versions)

	for i := 0; i < l; i++ {
		one := &core.local.Versions[i]
		if one.Version == version {
			return one.VersionRoot
		}
	}
	return ""
}

//GetBinPathByVersion 获取文件
func GetBinPathByVersion(version string) string {
	core := GlobalOneLocal()
	l := len(core.local.Versions)

	for i := 0; i < l; i++ {
		one := &core.local.Versions[i]
		if one.Version == version {
			return one.BinRoot
		}
	}
	return ""
}

//GetTemplatePathByVersion 获取local.json 中versionRoot路径
func GetTemplatePathByVersion(version string) string {
	core := GlobalOneLocal()
	l := len(core.local.Versions)

	for i := 0; i < l; i++ {
		one := &core.local.Versions[i]
		if one.Version == version {
			return one.VersionRoot
		}
	}
	return ""
}

// //ReceiveChainCodeUploadFile 接受包存链码文件压缩包
// func ReceiveChainCodeUploadFile(w http.ResponseWriter, r *http.Request) error {
// 	file, header, err := r.FormFile("file")
// 	filename := header.Filename
// 	workPath, _ := os.Getwd()

// 	fileSavePath := filepath.ToSlash(filepath.Join(workPath, "version", "fabric1.4.4", "chaincode", filename))
// 	if runtime.GOOS == "windows" {
// 		fileSavePath = strings.Replace(fileSavePath, "/", "\\", -1)
// 	}

// 	out, err := os.Create(fileSavePath)
// 	if err != nil {
// 		web.OutputEnter(w, "", nil, errors.WithMessage(err, "create save file "+filename+"err info:"))
// 		return err
// 	}
// 	_, err = io.Copy(out, file)
// 	if err != nil {
// 		web.OutputEnter(w, "", nil, errors.WithMessage(err, "copy save file "+filename+"err info:"))
// 		return err
// 	}
// 	defer out.Close()
// 	var cmd *exec.Cmd
// 	ccPath := filepath.ToSlash(filepath.Join(workPath, "version", "fabric1.4.4", "chaincode"))
// 	execCommand := fmt.Sprintf("cd %s && tar -xvf %s", ccPath, fileSavePath)
// 	if runtime.GOOS == "windows" {
// 		cmd = exec.Command("cmd", "/c", execCommand)
// 	} else {
// 		cmd = exec.Command("/bin/bash", "-c", execCommand)
// 	}
// 	if err := cmd.Start(); err != nil {
// 		web.OutputEnter(w, "", nil, errors.WithMessage(err, "exec command1 start"))
// 		return err
// 	}
// 	if err := cmd.Wait(); err != nil {
// 		web.OutputEnter(w, "", nil, errors.WithMessage(err, "exec command1 wait"))
// 		return err
// 	}
// 	execCommand = fmt.Sprintf("rm -rf %s/%s", ccPath, filename)
// 	if runtime.GOOS != "windows" {
// 		cmd = exec.Command("/bin/bash", "-c", execCommand)
// 		if err := cmd.Start(); err != nil {
// 			web.OutputEnter(w, "", nil, errors.WithMessage(err, "exec command2 wait"))
// 			return err
// 		}
// 		if err := cmd.Wait(); err != nil {
// 			web.OutputEnter(w, "", nil, errors.WithMessage(err, "exec command2 wait"))
// 			return err
// 		}
// 	}

// 	//把新上传的合约写入数据库
// 	err = dmysql.UploadChainCodeWriteDB(filename)
// 	if err != nil {
// 		web.OutputEnter(w, "", nil, errors.WithMessage(err, "upload chaincode write db error"))
// 		return err
// 	}
// 	return nil
// }

//GetChainCodeList 获取链码列表
func GetChainCodeList(w http.ResponseWriter, r *http.Request) ([]string, error) {
	workPath, _ := os.Getwd()
	fileSavePath := filepath.ToSlash(filepath.Join(workPath, "version", "fabric1.4.4", "chaincode"))
	if runtime.GOOS == "windows" {
		fileSavePath = strings.Replace(fileSavePath, "/", "\\", -1)
	}
	var ccList []string
	dir, err := ioutil.ReadDir(fileSavePath)
	if err != nil {
		return nil, err
	}
	for _, fi := range dir {
		if fi.IsDir() {
			ccList = append(ccList, fi.Name())
		}
	}
	return ccList, nil
}
