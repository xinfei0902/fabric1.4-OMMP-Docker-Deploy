package daction

import (
	"deploy-server/app/dconfig"
	"deploy-server/app/dmysql"
	"deploy-server/app/objectdefine"
	"fmt"
	"runtime"
	// "time"
    // "os"
	// "io"
	"os/exec"
)

func MakeStepGeneralCreateCompleteDeleteDeployRemote(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
        //从数据库获取所有远端节点ip
		two, err := dmysql.GetIPListFromIndent()
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		//清理
		fmt.Println("clean",two)
		for ip,_ := range two {
			connectPort := dconfig.GetStringByKey("toolsPort")
			url := fmt.Sprintf("http://%s:%s/command/exec", ip, connectPort)
			output.AppendLog(&objectdefine.StepHistory{
				Name: "GeneralDeleteDeployExecScriptStep1",
				Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
			})
			command := "start"
			buff := "docker rm -f $(docker ps -aq) && docker volume prune -f && docker network rm anhui_default && cd /blockchainData && ls | grep -v deploy-tools | xargs -i rm -rf {}"
			fmt.Println("command ==",buff)
			args := []string{buff}
			err := MakeHTTPRemoteCmd(url, command, args, output)
			if err != nil {
				// output.AppendLog(&objectdefine.StepHistory{
				// 	Name:  "GeneralDeleteDeployExecScriptStep1",
				// 	Error: []string{err.Error()},
				// })
				// return err
			}
            //
			buffPackage := "cd /blockchainData/deploy-tools && rm -rf add* && rm -rf create* && rm -rf disable* && rm -rf enable* && rm -rf upgrade* && rm -rf deploy && rm -rf deploy.tar.gz && rm -rf "+ip+"*"
			fmt.Println("command ==",buffPackage)
			ppargs := []string{buffPackage}
			err = MakeHTTPRemoteCmd(url, command, ppargs, output)
			if err != nil {
				// output.AppendLog(&objectdefine.StepHistory{
				// 	Name:  "GeneralDeleteDeployExecScriptStep1",
				// 	Error: []string{err.Error()},
				// })
				// return err
			}
		}


		//从数据库获取所有远端节点ip
		orderer, err := dmysql.GetIPListFromOrderer()
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		//清理
		fmt.Println("clean",orderer)
		for ip,_ := range orderer {
			connectPort := dconfig.GetStringByKey("toolsPort")
			url := fmt.Sprintf("http://%s:%s/command/exec", ip, connectPort)
			output.AppendLog(&objectdefine.StepHistory{
				Name: "GeneralDeleteDeployExecScriptStep1",
				Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
			})
			command := "start"
			
			buff := "docker rm -f $(docker ps -aq) && docker volume prune -f && docker network rm anhui_default && cd /blockchainData && ls | grep -v deploy-tools | xargs -i rm -rf {}"
			fmt.Println("command ==",buff)
			args := []string{buff}
			err := MakeHTTPRemoteCmd(url, command, args, output)
			if err != nil {
				// output.AppendLog(&objectdefine.StepHistory{
				// 	Name:  "GeneralDeleteDeployExecScriptStep1",
				// 	Error: []string{err.Error()},
				// })
				// return err
			}
			//
			buffPackage := "cd /blockchainData/deploy-tools && rm -rf add* && rm -rf create* && rm -rf disable* && rm -rf enable* && rm -rf upgrade* && rm -rf deploy && rm -rf deploy.tar.gz"
			fmt.Println("command ==",buffPackage)
			ppargs := []string{buffPackage}
			err = MakeHTTPRemoteCmd(url, command, ppargs, output)
			if err != nil {
				// output.AppendLog(&objectdefine.StepHistory{
				// 	Name:  "GeneralDeleteDeployExecScriptStep1",
				// 	Error: []string{err.Error()},
				// })
				// return err
			}
		}
		return nil
	}
}


func MakeStepGeneralCreateCompleteDeleteDeployLocal(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
        
		var cmd *exec.Cmd
		execCommand := fmt.Sprintf("rm -rf ./output/* && rm -rf ./version/fabric1.4.4/chaincode/*")
		fmt.Println("execCommand",execCommand)
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/c", execCommand)
		} else {
			cmd = exec.Command("/bin/bash", "-c", execCommand)
		}
		
		if err := cmd.Start(); err != nil {
			return err
		}

		if err := cmd.Wait(); err != nil {
			return err
		}

		return nil
	}
}


func MakeStepGeneralCreateCompleteDeleteDB(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
        
		var cmd *exec.Cmd
			execCommand := fmt.Sprintf("rm -rf ./output/*")
			fmt.Println("execCommand",execCommand)
			if runtime.GOOS == "windows" {
				cmd = exec.Command("cmd", "/c", execCommand)
			} else {
				cmd = exec.Command("/bin/bash", "-c", execCommand)
			}
			
			if err := cmd.Start(); err != nil {
				return err
			}
	
			if err := cmd.Wait(); err != nil {
				return err
			}
		return nil
	}
}