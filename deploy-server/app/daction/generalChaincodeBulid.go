package daction

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"github.com/pkg/errors"
	"io"
	"os/exec"
	"fmt"
	"deploy-server/app/objectdefine"
	"deploy-server/app/tools"
	"io/ioutil"
)
var CCFolderName string
func MakeGeneralCreateChainCodeBulid(w http.ResponseWriter, r *http.Request) error{
	file, header, err := r.FormFile("file")
	filename := header.Filename
	workPath, _ := os.Getwd()
	fileSavePath := filepath.ToSlash(filepath.Join(workPath, "version", "fabric1.4.4", "chaincode", filename))
	if runtime.GOOS == "windows" {
		fileSavePath = strings.Replace(fileSavePath, "/", "\\", -1)
	}
	str := strings.Split(filename, ".tar")
	if len(str) > 0 {
		CCFolderName = str[0]
	} else {
		return errors.New("chaincode Name error:please check chaincode Name")
	}
	out, err := os.Create(fileSavePath)
	if err != nil {
		return errors.WithMessage(err, "create save file "+filename+"err info:")
	}
	_, err = io.Copy(out, file)
	if err != nil {
		return errors.WithMessage(err, "copy save file "+filename+"err info:")
	}
	defer out.Close()
	var cmd *exec.Cmd
	ccPath := filepath.ToSlash(filepath.Join(workPath, "version", "fabric1.4.4", "chaincode"))
	execCommand := fmt.Sprintf("cd %s && tar -xf %s", ccPath, fileSavePath)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", execCommand)
	} else {
		cmd = exec.Command("/bin/bash", "-c", execCommand)
	}
	if err := cmd.Start(); err != nil {
		return errors.WithMessage(err, "exec command1 start")
	}
	if err := cmd.Wait(); err != nil {
		return errors.WithMessage(err, "exec command1 wait")
	}
	return nil
}

func MakeStepGeneralChainCodeBulidScripCopy(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		workPath, _ := os.Getwd()
		dst :=filepath.ToSlash(filepath.Join(workPath, "version", "fabric1.4.4", "chaincode",CCFolderName))
		src :=filepath.ToSlash(filepath.Join(workPath, "version", "fabric1.4.4", "template","chaincode","bulid.sh"))
        err := tools.CopyFileOne(dst,src)
		if err != nil{
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		return nil
	}
}

//MakeStepGeneralChainCodeBulidExecScript 执行脚本
func MakeStepGeneralChainCodeBulidExecScript(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		workPath, _ := os.Getwd()
		filePath := filepath.ToSlash(filepath.Join(workPath, "version", "fabric1.4.4", "chaincode", CCFolderName))
		var cmd *exec.Cmd
	    execCommand := fmt.Sprintf("cd %s && chmod +x bulid.sh && rm -rf bulid.log && ./bulid >> bulid.log", filePath)
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/c", execCommand)
		} else {
			cmd = exec.Command("/bin/bash", "-c", execCommand)
		}
		if err := cmd.Start(); err != nil {
			return errors.WithMessage(err, "exec command1 start")
		}
		if err := cmd.Wait(); err != nil {
			return errors.WithMessage(err, "exec command1 wait")
		}
		return nil
	}
}


func MakeGeneralChainCodeBulidScripCopyStep(ccFileName string) error {
		workPath, _ := os.Getwd()
		dst :=filepath.ToSlash(filepath.Join(workPath, "version", "fabric1.4.4", "chaincode",ccFileName,"bulid.sh"))
		_,err := os.Stat(dst) 
		if err == nil{
			var cmd *exec.Cmd
			execCommand :="rm -rf "+dst+""
			if runtime.GOOS == "windows" {
				cmd = exec.Command("cmd", "/c", execCommand)
			} else {
				cmd = exec.Command("/bin/bash", "-c", execCommand)
			}
			if err := cmd.Start(); err != nil {
				//return errors.WithMessage(err, "exec command1 start")
			}
			if err := cmd.Wait(); err != nil {
				//return errors.WithMessage(err, "exec command1 wait")
			}
		}
		
		src :=filepath.ToSlash(filepath.Join(workPath, "version", "fabric1.4.4", "template","chaincode","bulid.sh"))
		fmt.Println("dst",dst)
		fmt.Println("src",src)
        err = tools.CopyFileOne(dst,src)
		if err != nil{
			return err
		}
		return nil
}

//MakeStepGeneralChainCodeBulidExecScript 执行脚本
func MakeGeneralChainCodeBulidExecScriptStep(ccFileName string) error {
	workPath, _ := os.Getwd()
	filePath := filepath.ToSlash(filepath.Join(workPath, "version", "fabric1.4.4", "chaincode", ccFileName))
	var cmd *exec.Cmd
	execCommand := fmt.Sprintf("cd %s && chmod +x bulid.sh && rm -rf bulid.log && ./bulid.sh "+ccFileName+" >> bulid.log 2>&1", filePath)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", execCommand)
	} else {
		cmd = exec.Command("/bin/bash", "-c", execCommand)
	}
	if err := cmd.Start(); err != nil {
		//return errors.WithMessage(err, "exec command1 start")
	}
	if err := cmd.Wait(); err != nil {
		//return errors.WithMessage(err, "exec command1 wait")
	}
	return nil

}

func MakeGeneralReadChaincodeBulidResult(ccFileName string) ([]byte,error){
	workPath, _ := os.Getwd()
	filePath := filepath.ToSlash(filepath.Join(workPath, "version", "fabric1.4.4", "chaincode", ccFileName,"bulid.log"))
	b, err := ioutil.ReadFile(filePath) // just pass the file name
    if err != nil {
        fmt.Print(err)
    }
    return b,nil
}