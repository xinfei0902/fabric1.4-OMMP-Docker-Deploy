package dssh

import (
	"deploy-server/app/objectdefine"
	"deploy-server/app/dcache"
	"deploy-server/app/dconfig"
	"deploy-server/app/dmysql"
	"fmt"
	"bytes"
	"time"
	"encoding/json"
	"os"
	"io"
	"path/filepath"
	"golang.org/x/crypto/ssh"
	"github.com/pkg/sftp"
	"io/ioutil"
	"net/http"
	 "os/exec"
	 "runtime"
	 "bufio"
	 "strings"
	//"github.com/tmc/scp"
	"github.com/pkg/errors"
)
func MakeStepGeneralCheckRemoteServerSShConnect(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"start check ssh connect"},
		})
		ret := general.Server
		user := ret.ServerUser
		pass := ret.ServerPassword
		ip := ret.ServerExtIp
		client,err := getSSHClientSession(ip,user,pass)
		if err != nil{
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		session, err := client.NewSession()
		defer session.Close()
		//创建存放路径
		//destDir := "/root/blockChain/deploy-tools"
		destDir := dconfig.GetStringByKey("toolsDeployPath")
	    err = session.Run("mkdir -p "+destDir+"/conf")
		if err != nil {
			return err
		}

		version := general.Version
		rootPath := dcache.GetVersionRootPathByVersion(version)

		//创建deployEnv.sh
		exchangeMap := make(map[string]string)
		//exchangeMap["target"] = "deploy"
		exchangeMap["DEPLOY_TOOLS_PATH"] = destDir
		step2Buff, err := dcache.GetDeployEnvScriptTemplate(general.Version, exchangeMap)
		if err != nil {
			err = errors.WithMessage(err, "Build deployEnv.sh replace parame error")
			return err
		}
	
		err = ioutil.WriteFile(filepath.Join(rootPath,"deployTools", "deployEnv.sh"), step2Buff, 0644)
		if err != nil {
			err = errors.WithMessage(err, "replace parame writer script.sh error")
			return err
		}

        //打成压缩包
		// fmt.Println("tar path in",rootPath)
		// err = daction.GeneralDeployEnvZipFile(rootPath)
		// if err != nil {
		// 	err = errors.WithMessage(err, "Build deployTools.tar error")
		// 	return err
		// }
		//scp 传输工具
		//srcDir := fmt.Sprintf("%s/deployTools/deployTools.tar",rootPath)
	    deployToolsFolder := []string{"app.conf","deployEnv.sh","deploy-tools","go1.15.13.linux-amd64.tar.gz"}
		for _,file := range deployToolsFolder{
			srcDir := fmt.Sprintf("%s/deployTools/"+file+"",rootPath)
			if file == "app.conf"{
				srcDir = fmt.Sprintf("%s/deployTools/conf/"+file+"",rootPath)
			}
			
			sftpClient,err:=sftp.NewClient(client)
			//defer sftpClient.Close()
			if err != nil{
				output.AppendLog(&objectdefine.StepHistory{
					Name:  name,
					Error: []string{err.Error()},
				})
				return err
			}
			srcFile, err := os.Open(srcDir)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  name,
					Error: []string{err.Error()},
				})
				return err
			}
			//defer srcFile.Close()
			//dstFile, err := sftpClient.Create(filepath.Join(destDir, filepath.Base(srcDir)))
			fffolder := filepath.Join(destDir, file)
			if filepath.Base(srcDir) == "app.conf"{
				fffolder = filepath.Join(destDir,"conf", file)
			}
			fmt.Println("fffolder",fffolder)
			dstFile, err := sftpClient.Create(fffolder)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  name,
					Error: []string{err.Error()},
				})
			}
			//defer dstFile.Close()
			buf := make([]byte, 1024)
			for {
				n, _ := srcFile.Read(buf)
				if n == 0 {
					break
				}
				dstFile.Write(buf[0:n])
			}	
			//defer srcFile.Close()
			//defer dstFile.Close()
			srcFile.Close()
			dstFile.Close()
			sftpClient.Close()
		}
		

		//执行命令


		// execCommand := fmt.Sprintf("cd %s && chmod +x deploy-tools && nohup ./deploy-tools >> deploy-tools.log 2>&1 &",destDir)
        // fmt.Println("execCommand",execCommand)
		// session1, err := client.NewSession()
		// //defer session1.Close()
		// //buff, err := session1.Output("cd "+destDir+"; tar -xf  deployTools.tar; cd deployTools;chmod +x deploy-tools; nohup ./deploy-tools >> deploy-tools.log 2>&1 &")
		// buff, err := session1.Output(execCommand)

		// fmt.Println("ssh exec ++++")
		// if err != nil {
		// 	fmt.Println("ssh exec -----")
		// 	output.AppendLog(&objectdefine.StepHistory{
		// 		Name:  name,
		// 		Error: []string{err.Error()},
		// 	})
		//     return err
		// }
		// fmt.Println("exec ress")
		// fmt.Println("buff",string(buff))
		// session1.Close()
		time.Sleep(5 * time.Second)
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"end check ssh connect"},
		})
		return nil
	}
}


func MakeStepGeneralCheckRemotStartServerTools(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"enable tools start"},
		})
		ret := general.Server
		user := ret.ServerUser
		pass := ret.ServerPassword
		ip := ret.ServerExtIp
		client,err := getSSHClientSession(ip,user,pass)
		if err != nil{
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		destDir := dconfig.GetStringByKey("toolsDeployPath")
		//执行命令
		execCommand := fmt.Sprintf("cd %s; chmod +x deploy-tools; nohup ./deploy-tools >> deploy-tools.log 2>&1 &",destDir)
        fmt.Println("execCommand",execCommand)
		session1, err := client.NewSession()
		defer session1.Close()
		//buff, err := session1.Output("cd "+destDir+"; tar -xf  deployTools.tar; cd deployTools;chmod +x deploy-tools; nohup ./deploy-tools >> deploy-tools.log 2>&1 &")
		buff, err := session1.Output(execCommand)

		fmt.Println("ssh exec ++++")
		if err != nil {
			fmt.Println("ssh exec -----")
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
		    return err
		}
		fmt.Println("exec ress")
		fmt.Println("buff",string(buff))
		time.Sleep(5 * time.Second)
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"enable tools end"},
		})
		return nil
	}
}

func MakeStepGeneralRemoteServerBulidEnv(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"start bulid server env"},
		})
        //ip := general.Server.ServerIntIp
		ip := general.Server.ServerExtIp
		connectPort := dconfig.GetStringByKey("toolsPort")
		url := fmt.Sprintf("http://%s:%s/command/exec", ip, connectPort)
		output.AppendLog(&objectdefine.StepHistory{
			Name: "generalRemoteServerBulidEnv",
			Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
		})
		command := "start"
		//destDir :=filepath.Join(dconfig.GetStringByKey("toolsDeployPath"),"deployTools")
		destDir :=dconfig.GetStringByKey("toolsDeployPath")

		buff := fmt.Sprintf("cd "+destDir+" && chmod +x deployEnv.sh && ./deployEnv.sh && source /etc/profile")
		args := []string{buff}
		err := MakeHTTPRemoteCmd(url, command, args, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalRemoteServerBulidEnv",
				Error: []string{err.Error()},
			})
			return err
		}

		// ret := general.Server
		// user := ret.ServerUser
		// pass := ret.ServerPassword
		// ip := ret.ServerExtIp
		// client,err := getSSHClientSession(ip,user,pass)
		// if err != nil{
		// 	output.AppendLog(&objectdefine.StepHistory{
		// 		Name:  name,
		// 		Error: []string{err.Error()},
		// 	})
		// 	return err
		// }
		// session, err := client.NewSession()
		// defer session.Close()
		// _, err = session.Output("cd /root/blockChain/deploy-tools && chmod +x deployEnv.sh && ./deployEnv.sh")
		// fmt.Println("ssh exec ++++")
		// if err != nil {
		// 	fmt.Println("ssh exec -----")
		// 	output.AppendLog(&objectdefine.StepHistory{
		// 		Name:  name,
		// 		Error: []string{err.Error()},
		// 	})
		//     return err
		// }
		
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"end bulid server env"},
		})
		return nil
	}
}

func CheckRemoteServerSShConnect(server *objectdefine.IndentServer)error{
	    user := server.ServerUser
		pass := server.ServerPassword
		ip := server.ServerExtIp
		address := fmt.Sprintf("%s:22",ip)
		_, err := ssh.Dial("tcp", address, &ssh.ClientConfig{
			User:            user,
			Auth:            []ssh.AuthMethod{ssh.Password(pass)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
		if err != nil {
			return err
		}
		return nil
}

func getSSHClientSession(ip, user, pass string) (*ssh.Client,error){
	address := fmt.Sprintf("%s:22",ip)
	client, err := ssh.Dial("tcp", address, &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return nil,err
	}
	return client ,nil
}


func MakeStepGeneralCheckChaincodeIsInstantiated(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"check chaincode instantiated start"},
		})
		sourceIndent, err := dmysql.GetStartTaskBeforIndent(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		//
		var checkCCOrgName string
		var checkGccName string
		var checkGccVersion string
		for gccName,gcc := range general.Chaincode{
            for _,cc :=range gcc{
				checkGccVersion = cc.Version
				for _,cci := range cc.EndorsementOrg{
					checkCCOrgName = cci
					checkGccName = gccName
					break
				}
			}
		}
		orgInfo := make(map[string]objectdefine.OrgType, 0)
		for _, org := range sourceIndent.Org{
			if org.Name == checkCCOrgName{
				orgInfo[org.Name] = org
				break
			}
		}
		var peerAdminIP string
		var peerCliName string
		for _, org := range orgInfo{
            for _,peer := range org.Peer{
				peerAdminIP = peer.IP
				peerCliName = peer.CliName
				break
			}
		}
		//mysql查询主机服务
		serverInfo,err := dmysql.GetServiceInfoFromExtIP(peerAdminIP)
		if err != nil{
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{fmt.Sprintf("query server info fail:%s", err.Error())},
			})
			return err
		}
		//远端ssh
		//ret := general.Server
		user := serverInfo.ServerUser
		pass := serverInfo.ServerPassword
		ip := serverInfo.ServerExtIp
		client,err := getSSHClientSession(ip,user,pass)
		if err != nil{
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		session, err := client.NewSession()
		defer session.Close()
		//执行是否实例化命令
		execCommand := fmt.Sprintf("rm -rf peerInstantiatedStatus.log && docker exec -it %s /bin/bash -c 'peer chaincode list --instantiated -C %s' > peerInstantiatedStatus.log 2>&1 &",peerCliName,general.ChannelName)
        fmt.Println("execCommand",execCommand)
		buff, err := session.Output(execCommand)
		fmt.Println("ssh exec ++++")
		if err != nil {
			fmt.Println("ssh exec -----")
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
		    return err
		}
		fmt.Println("buff",string(buff))
		time.Sleep(5 * time.Second)
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"check chaincode instantiated exec command"},
		})
		//获取输出文件
		var cmd *exec.Cmd 
		execrmCommand := "rm -rf ./peerInstantiatedStatus.log"
		fmt.Println("execrmCommand",execrmCommand)
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/c", execrmCommand)
		} else {
			cmd = exec.Command("/bin/bash", "-c", execrmCommand)
		}
		if err := cmd.Start(); err != nil {
			return err
		}
		if err := cmd.Wait(); err != nil {
			return err
		}
		sftpClient,err:=sftp.NewClient(client)
		defer sftpClient.Close()
		if err != nil{
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		//remotePath :=  dconfig.GetStringByKey("toolsDeployPath")
		remotepath := fmt.Sprintf("%s/peerInstantiatedStatus.log",dconfig.GetStringByKey("toolsDeployPath"))
		srcFile, err := sftpClient.Open(remotepath)
		if err != nil {
			fmt.Println("文件读取失败", err)
			return err
		}
	    defer srcFile.Close()
	    dstFile, e := os.Create("peerInstantiatedStatus.log")
		if e != nil {
			fmt.Println("文件创建失败", e)
			return err
		}
	    defer dstFile.Close()
		if _, err1 := srcFile.WriteTo(dstFile); err1 != nil {
			fmt.Println("文件写入失败", err1)
			return err
		}
		//解析文件内容判断是否实例化
		f, err := os.Open("./peerInstantiatedStatus.log")
		if err != nil {
			return err
		}
		defer f.Close()
		rd := bufio.NewReader(f)
		for {
			line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
			if err != nil || io.EOF == err {
				break
			}
			fmt.Println(line)
			checkString := fmt.Sprintf("Name: %s, Version: %s",checkGccName,checkGccVersion)
			isContains := strings.Contains(line, checkString)
			if isContains {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  name,
					Error: []string{"chaincode has been instantiated,Can not be deleted"},
				})
				err = errors.WithMessage(nil,"chaincode has been instantiated,Can not be deleted")
				return err
			}
		}  
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"check chaincode instantiated end"},
		})
		return nil
	}
}

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