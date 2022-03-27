package daction

import (
	"deploy-server/app/dconfig"
	"deploy-server/app/objectdefine"
	"fmt"
)

//MakeStepGeneralOperateOrgExecScript 执行脚本命令
func MakeStepGeneralCreateCompleteDeployStartContainer(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec start container script start"},
		})

		err := GeneralCreateCompleteDeployStartContainer(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			// result := err
			// err := dmysql.UpdateFailCreateOrgTaskStatus(general)
			// if err != nil {
			// 	output.AppendLog(&objectdefine.StepHistory{
			// 		Name:  "generalCreateOperateOrgExecScript",
			// 		Error: []string{err.Error()},
			// 	})
			// 	return err
			// }
			// return result
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec add org Operate script end"},
		})
		return nil
	}
}

func GeneralCreateCompleteDeployStartContainer(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	if general.Consensus == "kafka"{
		for _, zookeeper := range general.Kafka.Zookeeper {
			err := GeneralCreateCompleteDeployStartZookeeperContainerStep(general, zookeeper, output)
			if err != nil {
				return err
			}
		}
		for _, kafka := range general.Kafka.Kafka {
			err := GeneralCreateCompleteDeployStartKafkaContainerStep(general, kafka, output)
			if err != nil {
				return err
			}
		}
	}

	for _, orderer := range general.Orderer {
		err := GeneralCreateCompleteDeployStartOrdererContainerStep(general, orderer, output)
		if err != nil {
			return err
		}
	}

	for _, org := range general.Org {
		for _,peer := range org.Peer{
			err := GeneralCreateCompleteDeployStartPeerContainerStep(general, peer, output)
			if err != nil {
				return err
			}
		}
	}
 
	//启动ca
	for _, org := range general.Org {
		err := GeneralCreateCompleteDeployStartCAContainerStep(general, org ,output)
		if err != nil {
			return err
		}
	}

	for _, org := range general.Org {
		for _,peer := range org.Peer{
			//if peer.User == "Admin"{
				cliName := fmt.Sprintf("cli-%s-%s", org.Name, peer.Name)
			//	if cliName == completeDeployCliName{
					err := GeneralCreateCompleteDeployStartCliContainerStep(general, cliName ,peer, output)
					if err != nil {
						return err
					}
				//}	
			//}
		}
	}
	return nil
}

func GeneralCreateCompleteDeployStartZookeeperContainerStep(general *objectdefine.Indent, zookeeper objectdefine.ZooKeeperType, output *objectdefine.TaskNode) error {

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", zookeeper.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateCompleteDeployStartZookeeperContainer",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	command := "start"
	buff := fmt.Sprintf("cd %s && COMPOSE_PROJECT_NAME=anhui docker-compose -f base.yaml up -d %s", "deploy",zookeeper.Domain)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		// output.AppendLog(&objectdefine.StepHistory{
		// 	Name:  "generalCreateOperateOrgExecScript",
		// 	Error: []string{err.Error()},
		// })
		return err
	}
	return nil
}

func GeneralCreateCompleteDeployStartKafkaContainerStep(general *objectdefine.Indent, kafka objectdefine.KafkaType, output *objectdefine.TaskNode) error {

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", kafka.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateCompleteDeployStartKafkaContainer",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	command := "start"
	buff := fmt.Sprintf("cd %s && COMPOSE_PROJECT_NAME=anhui docker-compose -f base.yaml up -d %s", "deploy",kafka.Domain)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		// output.AppendLog(&objectdefine.StepHistory{
		// 	Name:  "generalCreateOperateOrgExecScript",
		// 	Error: []string{err.Error()},
		// })
		return err
	}
	return nil
}

func GeneralCreateCompleteDeployStartOrdererContainerStep(general *objectdefine.Indent, orderer objectdefine.OrderType, output *objectdefine.TaskNode) error {

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", orderer.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateCompleteDeployStartOrdererContainer",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	command := "start"
	buff := fmt.Sprintf("cd %s && COMPOSE_PROJECT_NAME=anhui docker-compose -f base.yaml up -d %s", "deploy",orderer.Domain)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		// output.AppendLog(&objectdefine.StepHistory{
		// 	Name:  "generalCreateOperateOrgExecScript",
		// 	Error: []string{err.Error()},
		// })
		return err
	}
	return nil
}

func GeneralCreateCompleteDeployStartPeerContainerStep(general *objectdefine.Indent, peer objectdefine.PeerType, output *objectdefine.TaskNode) error {

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", peer.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateCompleteDeployStartPeerContainer",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	command := "start"
	buff := fmt.Sprintf("cd %s && COMPOSE_PROJECT_NAME=anhui docker-compose -f base.yaml up -d %s", "deploy",peer.Domain)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		// output.AppendLog(&objectdefine.StepHistory{
		// 	Name:  "generalCreateOperateOrgExecScript",
		// 	Error: []string{err.Error()},
		// })
		return err
	}
	return nil
}


func GeneralCreateCompleteDeployStartCAContainerStep(general *objectdefine.Indent, org objectdefine.OrgType,output *objectdefine.TaskNode) error {

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", org.CA.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateCompleteDeployStartPeerContainer",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	command := "start"
	buff := fmt.Sprintf("cd %s && COMPOSE_PROJECT_NAME=anhui docker-compose -f base.yaml up -d %s", "deploy",org.CA.Name)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		// output.AppendLog(&objectdefine.StepHistory{
		// 	Name:  "generalCreateOperateOrgExecScript",
		// 	Error: []string{err.Error()},
		// })
		return err
	}
	return nil
}

func GeneralCreateCompleteDeployStartCliContainerStep(general *objectdefine.Indent, cliName string,peer objectdefine.PeerType, output *objectdefine.TaskNode) error {

	connectPort := dconfig.GetStringByKey("toolsPort")
	url := fmt.Sprintf("http://%s:%s/command/exec", peer.IP, connectPort)
	output.AppendLog(&objectdefine.StepHistory{
		Name: "generalCreateCompleteDeployStartCliContainer",
		Log:  []string{fmt.Sprintf("connect remote tools url:%s", url)},
	})
	command := "start"
	buff := fmt.Sprintf("cd %s && COMPOSE_PROJECT_NAME=anhui docker-compose -f base.yaml up -d %s", "deploy",cliName)
	args := []string{buff}
	err := MakeHTTPRemoteCmd(url, command, args, output)
	if err != nil {
		// output.AppendLog(&objectdefine.StepHistory{
		// 	Name:  "generalCreateOperateOrgExecScript",
		// 	Error: []string{err.Error()},
		// })
		return err
	}
	return nil
}