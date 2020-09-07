package daction

import (
	"bytes"
	"deploy-server/app/objectdefine"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

//GetLatestIndentJSONFIleInfo 获取最新的联盟信息 包含组织 排序节点  节点
func GetLatestIndentJSONFIleInfo(sourceBaseOutput string) (indent *objectdefine.Indent, err error) {
	input := &objectdefine.Indent{}

	newIndentPath := filepath.Join(sourceBaseOutput, "NewIndent.json")
	_, err = os.Stat(newIndentPath)
	if err == nil {
		data, err := ioutil.ReadFile(newIndentPath)
		err = json.Unmarshal(data, input)
		if err != nil {
			return nil, errors.WithMessage(err, "get NewIndent.json json2struct error")
		}
	} else {
		oldfilePath := filepath.Join(sourceBaseOutput, "indent.json")
		buff, err := ioutil.ReadFile(oldfilePath)
		ioutil.WriteFile(newIndentPath, buff, 0777)
		err = json.Unmarshal(buff, input)
		if err != nil {
			return nil, errors.WithMessage(err, "get indent.json json2struct error")
		}
	}
	return input, nil
}

//GetURLByIP 根据传进来的节点ip 获取连接工具URL地址
func GetURLByIP(ip string, secret []objectdefine.IPSecret) *objectdefine.IPSecret {
	for i := range secret {
		if secret[i].IP == ip {
			return &secret[i]
		}
	}
	return nil
}

//newfileUploadRequest 上传文件
func newfileUploadRequest(uri string, params map[string]string, paramName, path, dstPath string, output *objectdefine.TaskNode) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "upload",
			Error: []string{("os open url fail" + err.Error())},
		})
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "upload",
			Error: []string{("create file" + err.Error())},
		})
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "upload",
			Error: []string{("writer close fail" + err.Error())},
		})
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "upload",
			Error: []string{("http Post fail" + err.Error())},
		})
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

//MakeHTTPRemoteCmd 与工具连接 发送命令 判断返回的结果
func MakeHTTPRemoteCmd(url, command string, args []string, output *objectdefine.TaskNode) (err error) {
	one := &objectdefine.RequestBody{
		Command: command,
		Args:    args,
	}

	buff, err := json.Marshal(one)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "remote cmd",
			Error: []string{fmt.Sprintf("Analysis request body fail:%s", err.Error())},
		})
		return err
	}
	client := &http.Client{}
	reader := bytes.NewBuffer(buff)
	request, _ := http.NewRequest("POST", url, reader)
	request.Header.Set("Content-type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "remote cmd",
			Error: []string{fmt.Sprintf("exec request clent.Do fail:%s", err.Error())},
		})
		return err
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	if response.StatusCode != 200 {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "remote cmd",
			Error: []string{fmt.Sprintf("exec request return statusCode:%d fail:%s", response.StatusCode, err.Error())},
		})
		return errors.New(fmt.Sprintf("exec request return statusCode:%d fail:%s", response.StatusCode, err.Error()))
	}
	buff, err = ioutil.ReadAll(response.Body)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "remote cmd",
			Error: []string{fmt.Sprintf("exec request return response body:%s fail:%s", string(buff), err.Error())},
		})
		return errors.New("read response body false")
	}
	two := &objectdefine.ReponseBody{}
	err = json.Unmarshal(buff, two)
	if two.Code.Code == 200 {
		return nil
	}
	output.AppendLog(&objectdefine.StepHistory{
		Name:  "remote cmd",
		Error: []string{fmt.Sprintf("exec script return Error:%s", two.Code.Message)},
	})
	return errors.New(fmt.Sprintf("exec script return Error:%s", two.Code.Message))
}

//GetTaskIDRunStatus 暂且先不用
//GetTaskIDRunStatus 获取任务ID的状态
func GetTaskIDRunStatus(baseOutput string) (bool, error) {

	name := filepath.Join(baseOutput, objectdefine.ConstHistoryTaskStatusFileName)
	buff, err := ioutil.ReadFile(name)
	if err != nil {
		return false, errors.WithMessage(err, "Read logs file error")
	}

	one := &objectdefine.IndentStatus{}
	err = json.Unmarshal(buff, one)
	if err != nil {
		return false, errors.WithMessage(err, "Parse indent file error")
	}
	for _, plains := range one.Plains {
		if plains.Done == true && plains.Success == false {
			return false, nil
		}
	}
	return true, nil
}
