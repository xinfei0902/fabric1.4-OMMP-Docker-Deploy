package daction

import (
	"deploy-server/app/dcache"
	"deploy-server/app/dmysql"
	"deploy-server/app/objectdefine"
	"deploy-server/app/tools"
	"fmt"
	
	"io/ioutil"
	"net"
	
	"strings"
	"path/filepath"
	"github.com/pkg/errors"
)

//CheckCompletionIndent  说明:  检测订单 补全
func CheckCompletionIndent(indent *objectdefine.Indent, model string) (*objectdefine.Indent, error) {
	ret := &objectdefine.Indent{}
	ret = indent
	// creatChannelName := indent.ChannelName
	// if model == "createChannel" {
	// 	indent.ChannelName = "score-channel"
	// }
    // if len(indent.ChannelName)==0{
    //   return nil, errors.New("Empty buff for channelName")
	// }
	creatChannelName := indent.ChannelName
	ret.ID = strings.ToLower(indent.ID)
	if len(ret.ID) == 0 {
		return nil, errors.New("Empty buff for ID")
	}
	if model != "createServer"{
		if len(indent.ChannelName) == 0 {
			return nil, errors.New("channel name cannot be empty")
		}
		ret.ChannelName = indent.ChannelName
		fmt.Println("ret.ChannelName",ret.ChannelName)
		ret.Model = model
		if len(indent.Consensus) == 0 {
			ret.Consensus = "etcdraft"
		} else {
			ret.Consensus = indent.Consensus
		}
	}
	
	if len(indent.Version) == 0 {
		ret.Version = dcache.GetDefaultVersion()
	} else {
		ret.Version = indent.Version
	}
	

	switch model {
	case "completeDeploy":
		wholeIndent, err := checkCompleteDeployIndent(indent)
		if err != nil {
			return nil, errors.WithMessage(err, " check rapid deploy blockchain fail")
		}

		checkPortIndent, err := checkCompleteDeployIndentPort(wholeIndent)
		if err != nil {
			return nil, errors.WithMessage(err, " check rapid deploy blockchain fail")
		}

		checkDeployIndent, err := checkCompleteDeployIndentDeploy(checkPortIndent)
		if err != nil {
			return nil, errors.WithMessage(err, " check rapid deploy blockchain fail")
		}
		fmt.Println("checkDeployIndent")
		ret = checkDeployIndent
		break
	case "createChannel":
		//获取原通道
		channFilePath :=filepath.ToSlash(dcache.GetOutputSubPath("sourceid", "channelFile.txt"))
		channName, err := ioutil.ReadFile(channFilePath) // just pass the file name
		if err != nil {
			fmt.Print(err)
		}
        readFileChannelName := strings.Replace(string(channName), " ", "", -1)  
		readFileChannelName = strings.Replace(string(channName), "\n", "", -1) 
		indent.ChannelName = readFileChannelName 
	
		sourceIndent, err := dmysql.GetStartTaskBeforIndent(indent)
		if err != nil {
			return nil, errors.New("get mysql indent data fail")
		}	
		indent.ChannelName = creatChannelName
		orgType, err := checkCreateChannelIndent(sourceIndent, indent)
		if err != nil {
			return nil, errors.WithMessage(err, " check create channel indent fail")
		}
		ret.Org = orgType
		ret.Orderer = sourceIndent.Orderer
		break
	case "addOrg":
		sourceIndent, err := dmysql.GetStartTaskBeforIndent(indent)
		if err != nil {
			return nil, errors.New("get mysql indent data fail")
		}
		if len(indent.Org) == 0 {
			return nil, errors.New("Empty buff indent.Org configs")
		}
		orgType, err := checkAddOrgIndent(sourceIndent, indent)
		if err != nil {
			return nil, errors.WithMessage(err, " check add org indent fail")
		}
		ret.Org = orgType
		break
	case "deleteOrg":
		sourceIndent, err := dmysql.GetDeleteTaskBeforIndent(indent)
		if err != nil {
			return nil, errors.New("get mysql indent data fail")
		}
		if len(indent.Org) == 0 {
			return nil, errors.New("Empty buff indent.Org configs")
		}
		orgType, err := checkDeleteOrgIndent(sourceIndent, indent)
		if err != nil {
			return nil, errors.WithMessage(err, " check delete org indent fail")
		}
		ret.Org = orgType
		break
	case "addPeer":
		sourceIndent, err := dmysql.GetStartTaskBeforIndent(indent)
		if err != nil {
			return nil, errors.New("get mysql indent data fail")
		}
		if len(indent.Org) == 0 {
			return nil, errors.New("Empty buff indent.Org configs")
		}
		fmt.Println("add peer chech peer info start")
		orgType, err := checkAddPeerIndent(sourceIndent, indent)
		if err != nil {
			return nil, errors.WithMessage(err, " check add peer indent fail")
		}
		fmt.Println("add peer chech peer info end", orgType)
		ret.Org = orgType
		break
	case "deletePeer":
		sourceIndent, err := dmysql.GetStartTaskBeforIndent(indent)
		if err != nil {
			return nil, errors.New("get mysql indent data fail")
		}
		if len(indent.Org) == 0 {
			return nil, errors.New("Empty buff indent.Org configs")
		}
		orgType, err := checkDeletePeerIndent(sourceIndent, indent)
		if err != nil {
			return nil, errors.WithMessage(err, " check delete peer indent fail")
		}
		ret.Org = orgType
		break
	case "disablePeer":
		sourceIndent, err := dmysql.GetStartTaskBeforIndent(indent)
		if err != nil {
			return nil, errors.New("get mysql indent data fail")
		}
		if len(indent.Org) == 0 {
			return nil, errors.New("Empty buff indent.Org configs")
		}
		orgType, err := checkDisablePeerIndent(sourceIndent, indent)
		if err != nil {
			return nil, errors.WithMessage(err, " check disable peer indent fail")
		}
		ret.Org = orgType
		break
	case "enablePeer":
		sourceIndent, err := dmysql.GetStartTaskBeforIndent(indent)
		if err != nil {
			return nil, errors.New("get mysql indent data fail")
		}
		if len(indent.Org) == 0 {
			return nil, errors.New("Empty buff indent.Org configs")
		}
		orgType, err := checkEnablePeerIndent(sourceIndent, indent)
		if err != nil {
			return nil, errors.WithMessage(err, " check enable peer indent fail")
		}
		ret.Org = orgType
		break
	case "modiflyPeer":
		sourceIndent, err := dmysql.GetStartTaskBeforIndent(indent)
		if err != nil {
			return nil, errors.New("get mysql indent data fail")
		}
		if len(indent.Org) == 0 {
			return nil, errors.New("Empty buff indent.Org configs")
		}
		orgType, err := checkModiflyPeerIndent(sourceIndent, indent)
		if err != nil {
			return nil, errors.WithMessage(err, " check enable peer indent fail")
		}
		ret.Org = orgType
		break
	case "chaincodeAdd":
		sourceIndent, err := dmysql.GetStartTaskBeforIndent(indent)
		if err != nil {
			return nil, errors.New("get mysql indent data fail")
		}
		ccType, err := checkAddChainCodeIndent(sourceIndent, indent)
		if err != nil {
			return nil, errors.WithMessage(err, " check add chaincode indent fail")
		}
		ret.Chaincode = ccType
		break
	case "chaincodeDelete":
		sourceIndent, err := dmysql.GetStartTaskBeforIndent(indent)
		if err != nil {
			return nil, errors.New("get mysql indent data fail")
		}
		ccType, err := checkDeleteChainCodeIndent(sourceIndent, indent)
		if err != nil {
			return nil, errors.WithMessage(err, " check delete chaincode indent fail")
		}
		ret.Chaincode = ccType
		break
	case "chaincodeUpgrade":
		sourceIndent, err := dmysql.GetStartTaskBeforIndent(indent)
		if err != nil {
			return nil, errors.New("get mysql indent data fail")
		}
		ccType, err := checkUpgradeChainCodeIndent(sourceIndent, indent)
		if err != nil {
			return nil, errors.WithMessage(err, " check upgrade chaincode indent fail")
		}
		ret.Chaincode = ccType
		break
	case "chaincodeDisable":
		sourceIndent, err := dmysql.GetStartTaskBeforIndent(indent)
		if err != nil {
			return nil, errors.New("get mysql indent data fail")
		}
		ccType, err := checkDisableChainCodeIndent(sourceIndent, indent)
		if err != nil {
			return nil, errors.WithMessage(err, " check disable chaincode indent fail")
		}
		ret.Chaincode = ccType
		break
	case "chaincodeEnable":
		sourceIndent, err := dmysql.GetStartTaskBeforIndent(indent)
		if err != nil {
			return nil, errors.New("get mysql indent data fail")
		}
		ccType, err := checkEnableChainCodeIndent(sourceIndent, indent)
		if err != nil {
			return nil, errors.WithMessage(err, " check enable chaincode indent fail")
		}
		ret.Chaincode = ccType
		break
	case "createServer":
		serverInfo,err := checkAddServerInfoIndent(indent)
		if err != nil{
			return nil, errors.WithMessage(err, " check add server info parm fail")
		}
		ret.Server = serverInfo
		break	
	}


	return ret, nil
}

//checkCompleteDeployIndent 检测一键部署订单信息 补全遗漏项信息
func checkCompleteDeployIndent(indent *objectdefine.Indent) (*objectdefine.Indent, error) {
	ret := &objectdefine.Indent{}
	ret.ID = strings.ToLower(indent.ID)
	ret.SourceID = strings.ToLower(indent.ID)
	if len(ret.ID) == 0 {
		return nil, errors.New("Empty buff for ID")
	}
	if len(indent.ChannelName) == 0 {
		return nil, errors.New("channel name cannot be empty")
	}
	ret.ChannelName = indent.ChannelName
	ret.Deploy = indent.Deploy
	if len(indent.Version) == 0 {
		ret.Version = dcache.GetDefaultVersion()
	} else {
		ret.Version = indent.Version
	}
	if len(indent.Consensus) == 0 {
		if indent.Kafka != nil {
			ret.Consensus = "kafka"
		} else {
			if len(indent.Orderer) > 1 {
				return nil, errors.New("Consensus [kafka] need more Zookeeper & Kafka configs")
			}
			ret.Consensus = objectdefine.ConsensusSolo
		}
	} else {
		ret.Consensus = indent.Consensus
	}
	ret.BaseOutput = indent.BaseOutput
	switch ret.Consensus {
	case objectdefine.ConsensusKafka:
		if indent.Kafka == nil || len(indent.Kafka.Kafka) == 0 || len(indent.Kafka.Zookeeper) == 0 {
			return nil, errors.New("Consensus [" + objectdefine.ConsensusKafka + "] need more Zookeeper & Kafka configs")
		}
		ret.Kafka = &objectdefine.KafkaStruct{}
		ret.Kafka.Kafka = make([]objectdefine.KafkaType, len(indent.Kafka.Kafka))

		for i, one := range indent.Kafka.Kafka {
			if len(one.IP) == 0 {
				return nil, errors.Errorf("IP Must exist in indent.Kafka.Kafka[%d]", i)
			}
			address := net.ParseIP(one.IP)
			if address == nil {
				return nil, errors.Errorf("IP wrong fromat in indent.Kafka.Kafka[%d]", i)
			}
			if len(one.Name) == 0 {
				one.Name = fmt.Sprintf("kafka%d", i)
			}
			if one.Port < 1 || one.Port > 65534 {
				one.Port = 9092
			}

			if one.BrokeID < 1 {
				one.BrokeID = i
			}

			// build
			if len(one.Domain) == 0 {
				one.Domain = one.Name
			}
			one.Domain = strings.ToLower(one.Domain)

			ret.Kafka.Kafka[i] = one
		}

		ret.Kafka.Zookeeper = make([]objectdefine.ZooKeeperType, len(indent.Kafka.Zookeeper))
		for i, one := range indent.Kafka.Zookeeper {
			if len(one.IP) == 0 {
				return nil, errors.Errorf("IP Must exist in indent.Kafka.Zookeeper[%d]", i)
			}
			address := net.ParseIP(one.IP)
			if address == nil {
				return nil, errors.Errorf("IP wrong fromat in indent.Kafka.Zookeeper[%d]", i)
			}
			if len(one.Name) == 0 {
				one.Name = fmt.Sprintf("zookeeper%d", i)
			}
			if one.Port < 1 || one.Port > 65534 {
				one.Port = 2181
			}
			if one.Follow < 1 || one.Follow > 65534 {
				one.Follow = 2888
			}
			if one.Vote < 1 || one.Vote > 65534 {
				one.Vote = 3888
			}

			if one.ID < 1 {
				one.ID = i
			}

			// build
			if len(one.Domain) == 0 {
				one.Domain = one.Name
			}
			one.Domain = strings.ToLower(one.Domain)
			ret.Kafka.Zookeeper[i] = one
		}

	default:
		ret.Kafka = nil
	}
	if len(indent.Org) == 0 {
		return nil, errors.New("Empty buff indent.Org configs")
	}
	ret.Org = make(map[string]objectdefine.OrgType, len(indent.Org))
	countID := 1
	for _, org := range indent.Org {
		if len(org.Name) == 0 {
			return nil, errors.New("Empty buff indent.Org.Name configs")
		}
		if len(org.OrgDomain) == 0 {
			//return nil, errors.New("Empty buff indent.Org.orgDomain configs")
			org.OrgDomain = fmt.Sprintf("%s.example.com", strings.ToLower(org.Name))
		}

		//##############peer##########################
		var anchorIP string

		peerList := make([]objectdefine.PeerType, len(org.Peer))
		for i, one := range org.Peer {
			if len(one.IP) == 0 {
				return nil, errors.Errorf("IP Must exist in indent.Org.Peer[%d]", i)
			}
			address := net.ParseIP(one.IP)
			if address == nil {
				return nil, errors.Errorf("IP wrong fromat in indent.Org.Peer[%d]", i)
			}
			one.OrgDomain = org.OrgDomain
			if len(one.Name) == 0 {
				one.Name = fmt.Sprintf("peer%d", i)
			}
			if len(one.Org) == 0 {
				one.Org = org.Name
			}

			if len(one.User) == 0 {
				if i == 0 {
					one.User = "Admin"
					anchorIP = one.IP
				} else {
					one.User = fmt.Sprintf("User%d", i)
				}
			}
			one.PeerID = countID - 1
			countID++

			if one.Port < 1 || one.Port > 65532 {
				one.Port = 7051
			}
			if one.CouchdbPort < 1 || one.CouchdbPort > 65532 {
				one.CouchdbPort = 5084
			}
			if one.ChaincodePort < 1 || one.ChaincodePort > 65532 {
				one.ChaincodePort = 7052
			}

			if len(one.Domain) == 0 {
				one.Domain = one.Name + "." + one.OrgDomain
			}
			one.Domain = strings.ToLower(one.Domain)
			if len(one.CliName) == 0 {
				one.CliName = fmt.Sprintf("cli-%s-%s", org.Name, one.Name)
			}
			peerList[i] = one
		}

		org.Peer = peerList

		//##################CA#################
		caType := &objectdefine.CAType{}
		if len(caType.Name) == 0 {
			caType.Name = fmt.Sprintf("ca-%s", org.Name)
		}
		if len(caType.IP) == 0 {
			caType.IP = anchorIP
		}
		address := net.ParseIP(caType.IP)
		if address == nil {
			return nil, errors.Errorf("IP wrong fromat in indent.Org[%s].CA", org.Name)
		}
		if caType.Port < 1 || caType.Port > 65535 {
			caType.Port = 7054
		}
		org.CA = caType
		ret.Org[org.Name] = org
	}

	//################orderer#######################
	orderList := make([]objectdefine.OrderType, len(indent.Orderer))
	for i, one := range indent.Orderer {
		if len(one.IP) == 0 {
			return nil, errors.Errorf("IP Must exist in indent.orderer[%d]", i)
		}
		address := net.ParseIP(one.IP)
		if address == nil {
			return nil, errors.Errorf("IP wrong fromat in indent.orderer[%d]", i)
		}
		if len(one.Name) == 0 {
			one.Name = fmt.Sprintf("orderer%d", i)
		}

		if one.Port < 1 || one.Port > 65534 {
			one.Port = 7050
		}

		// build
		if len(one.Domain) == 0 {
			one.Domain = fmt.Sprintf("%s.%s", one.Name, "example.com")
		}
		one.Domain = strings.ToLower(one.Domain)
		if len(one.OrgDomain) == 0 {
			one.OrgDomain = "example.com"
		}
		orderList[i] = one
	}

	ret.Orderer = orderList
	return ret, nil
}

//checkCompleteDeployIndentPort  检测一键已部署 端口是否重复 重复之后修改其值
func checkCompleteDeployIndentPort(indent *objectdefine.Indent) (*objectdefine.Indent, error) {
	//portPair := indent.FireWall

	fireWall := make(map[string]map[int]bool)
	if indent.Consensus == "kafka"{
		for i, one := range indent.Kafka.Zookeeper {
			if v, ok := fireWall[one.IP]; ok {
				if _, ok := v[one.Port]; ok {
				LOOP:
					for {
						one.Port += 1000
						if _, ok := v[one.Port]; ok {
							goto LOOP
						}
						goto END
					}
				END:
					v[one.Port] = true
				}
				v[one.Port] = true
				if _, ok := v[one.Follow]; ok {
				FLOOP:
					for {
						one.Follow += 1000
						if _, ok := v[one.Follow]; ok {
							goto FLOOP
						}
						goto FEND
					}
				FEND:
					v[one.Follow] = true
				}
	
				if _, ok := v[one.Vote]; ok {
				VLOOP:
					for {
						one.Vote += 1000
						if _, ok := v[one.Vote]; ok {
							goto VLOOP
						}
						goto VEND
					}
				VEND:
					v[one.Vote] = true
				}
				fireWall[one.IP] = v
			} else {
				portPair := make(map[int]bool, 0)
				portPair[one.Port] = true
				portPair[one.Follow] = true
				portPair[one.Vote] = true
				fireWall[one.IP] = portPair
			}
			indent.Kafka.Zookeeper[i] = one
		}
	
		for i, one := range indent.Kafka.Kafka {
			if v, ok := fireWall[one.IP]; ok {
				if _, ok := v[one.Port]; ok {
				KLOOP:
					for {
						one.Port += 100
						if _, ok := v[one.Port]; ok {
							goto KLOOP
						}
						goto KEND
					}
				KEND:
				}
				v[one.Port] = true
				fireWall[one.IP] = v
			} else {
				portPair := make(map[int]bool, 0)
				portPair[one.Port] = true
				fireWall[one.IP] = portPair
			}
			indent.Kafka.Kafka[i] = one
		}
	}
	

	for i, one := range indent.Orderer {
		if v, ok := fireWall[one.IP]; ok {
			if _, ok := v[one.Port]; ok {
			OLOOP:
				for {
					one.Port += 100
					if _, ok := v[one.Port]; ok {
						goto OLOOP
					}
					goto OEND
				}
			OEND:
			}
			v[one.Port] = true
			fireWall[one.IP] = v
		} else {
			portPair := make(map[int]bool, 0)
			portPair[one.Port] = true
			fireWall[one.IP] = portPair
		}
		indent.Orderer[i] = one
	}

	for orgName, org := range indent.Org {
		for i, one := range org.Peer {
			if v, ok := fireWall[one.IP]; ok {
				if _, ok := v[one.Port]; ok {
				PLOOP:
					for {
						one.Port += 100
						one.CouchdbPort += 100
						one.ChaincodePort += 100
						if _, ok := v[one.Port]; ok {
							goto PLOOP
						}
						goto PEND
					}
				PEND:
				}

				v[one.Port] = true
				v[one.CouchdbPort] = true
				v[one.ChaincodePort] = true
				fireWall[one.IP] = v
			} else {
				portPair := make(map[int]bool, 0)
				portPair[one.Port] = true
				portPair[one.CouchdbPort] = true
				portPair[one.ChaincodePort] = true
				fireWall[one.IP] = portPair
			}
			org.Peer[i] = one
		}
		if v, ok := fireWall[org.CA.IP]; ok {
			if _, ok := v[org.CA.Port]; ok {
			CLOOP:
				for {
					org.CA.Port += 1000
					if _, ok := v[org.CA.Port]; ok {
						goto CLOOP
					}
					goto CEND
				}
			CEND:
			}
			v[org.CA.Port] = true
			fireWall[org.CA.IP] = v
			org.CA = org.CA
		} else {
			portPair := make(map[int]bool, 0)
			portPair[org.CA.Port] = true
			fireWall[org.CA.IP] = portPair
		}
		indent.Org[orgName] = org
	}
	indent.FireWall = fireWall
	return indent, nil
}

//checkCompleteDeployIndentDeploy 检测补全通道组织 合约
func checkCompleteDeployIndentDeploy(indent *objectdefine.Indent) (*objectdefine.Indent, error) {
	if len(indent.Deploy) ==0 {
	  fmt.Println("CompleteDeploy blockChain not need install")
	  return indent,nil
	}else{
	deployMod := make(map[string]objectdefine.DeployType, 0)
	deployInfo := objectdefine.DeployType{}

	ccInfo := make(map[string]objectdefine.ChainCodeType, 0)
	for channelName, deploy := range indent.Deploy {
		orgInfo := make(map[string]objectdefine.OrgType, 0)
		for orgName, _ := range deploy.JoinOrg {
			if len(orgName) == 0 {
				return nil, errors.New("channel:+\"" + channelName + "\" join org is not empty")
			}
			if v, ok := indent.Org[orgName]; ok {
				orgInfo[orgName] = v
			} else {
				return nil, errors.New("channel:+\"" + channelName + "\" join org is not exist")
			}
		}
		for ccName, ccS := range deploy.JoinCC {
			if len(ccName) == 0 {
				break
			}
			if len(ccS.Version) == 0 {
				ccS.Version = "1.0"
			}
			if len(ccS.EndorsementOrg) == 0 {
				return nil, errors.New("chaincode endorsement org  is empty")
			}
			for _, orgName := range ccS.EndorsementOrg {
				if _, ok := indent.Org[orgName]; !ok {
					return nil, errors.New("chaincode endorsement org  is not exist")
				}
			}
			endorOrgNum := len(ccS.EndorsementOrg)
			endList := make([]string, 0)
			if endorOrgNum == 1 {
				endOrgName := fmt.Sprintf("%sMSP.member", ccS.EndorsementOrg[0])
				ccS.Policy = "OR('" + endOrgName + "')"
				ccInfo[ccName] = ccS
			} else {
				for i := 0; i < endorOrgNum; i++ {
					endOrgName := fmt.Sprintf("%sMSP.member", ccS.EndorsementOrg[i])
					endList = append(endList, "'"+endOrgName+"'")
				}
				ccS.Policy = "OR(" + strings.Join(endList, ",") + ")"
				ccInfo[ccName] = ccS
			}
		}
		deployInfo.JoinOrg = orgInfo
		deployInfo.JoinCC = ccInfo
		deployMod[channelName] = deployInfo
	}
	indent.Deploy = deployMod  
	fmt.Println("checkCompleteDeployIndentDeploy")
    }
	return indent, nil
}

//checkCreateChannelIndent 检测创建通道订单
func checkCreateChannelIndent(sourceIndent *objectdefine.Indent, general *objectdefine.Indent) (map[string]objectdefine.OrgType, error) {
	//分为两种情况 一种用已有组织创建通道 一种是启用新的组织创建通道
	orgType := make(map[string]objectdefine.OrgType, len(general.Org))
	peerType := make([]objectdefine.PeerType, 0)

	//检测通道是否存在
	status, err := dmysql.CheckChannelIsExist(general.ChannelName)
	if !status && err != nil {
		return nil, errors.WithMessage(err, "check channel is exist")
	}

	if len(sourceIndent.Org) == 0 {
		return nil, errors.WithMessage(err, "check org is empty")
	}
	for _, org := range sourceIndent.Org {
		for _,peer := range org.Peer{
			if peer.User == "Admin" && len(peerType) == 0  {
				fmt.Println("admin org",org.Name)
				peerType = append(peerType, peer)
				org.Peer = peerType
				orgType[org.Name] = org
				break
			}
		}
	}

	// allIP, err := dmysql.GetAllInfoformIndent()
	// if err != nil {
	// 	return nil, errors.New("get mysql indent data fail")
	// }
	//以为目前只支持新组织创建所以此处默认为true
	// general.IsNewOrgCreateChannel = true
	// if general.IsNewOrgCreateChannel == true {
	// 	//获取所有ip 对应端口
	// 	//allIP, err := dmysql.GetAllInfoformDeploy()
	// 	for _, org := range general.Org {
	// 		if len(org.Name) == 0 {
	// 			return nil, errors.New("Empty buff indent.Org.Name configs")
	// 		}
	// 		if len(org.OrgDomain) == 0 {
	// 			return nil, errors.New("Empty buff indent.Org.orgDomain configs")
	// 		}
	// 		if len(org.Peer) == 0 {
	// 			return nil, errors.New("Empty buff indent.Org.peer info configs")
	// 		}
	// 		peerList := make([]objectdefine.PeerType, len(org.Peer))
	// 		for i, one := range org.Peer {
	// 			if len(one.IP) == 0 {
	// 				return nil, errors.Errorf("IP Must exist in indent.Org.Peer[%d]", i)
	// 			}
	// 			one.OrgDomain = org.OrgDomain
	// 			if len(one.Name) == 0 {
	// 				return nil, errors.Errorf("name Must exist in indent.Org.Peer[%d]", i)
	// 			}
	// 			if len(one.User) == 0 {
	// 				//数据库查询当前组织有多少个节点
	// 				count, err := dmysql.GetCurrtOrgPeerNum(org.Name, general.ChannelName)
	// 				if count == -1 || err != nil {
	// 					return nil, errors.WithMessage(err, "mysql query fail")
	// 				}
	// 				if count == 0 {
	// 					one.User = "Admin"
	// 				} else {
	// 					one.User = fmt.Sprintf("User%d", count-1)
	// 				}
	// 			}
	// 			if one.Port < 1 || one.Port > 65535 {
	// 				return nil, errors.Errorf("indent Peer:%s prot %d Out of range", one.Name, one.Port)
	// 			}
	// 			if port, ok := allIP[one.IP]; ok {
	// 				if _, ok := port[one.Port]; ok {
	// 					return nil, errors.Errorf("indent Peer:%s prot %d is exist", one.Name, one.Port)
	// 				}
	// 			}
	// 			if len(one.Domain) == 0 {
	// 				one.Domain = one.Name + "." + one.OrgDomain
	// 			}
	// 			if len(one.CliName) == 0 {
	// 				one.CliName = fmt.Sprintf("cli-%s-%s", org.Name, one.Name)
	// 			}
	// 			peerList[i] = one
	// 		}
	// 		org.Peer = peerList
	// 		// caType := &objectdefine.CAType{}
	// 		caType := org.CA
	// 		if len(caType.Name) == 0 {
	// 			caType.Name = fmt.Sprintf("ca-%s", org.Name)
	// 		}
	// 		if len(caType.IP) == 0 {
	// 			return nil, errors.Errorf("IP Must exist in indent.Org.CA")
	// 		}

	// 		if caType.Port < 1 || caType.Port > 65535 {
	// 			return nil, errors.Errorf("indent Peer:%s prot %d Out of range", caType.Name, caType.Port)
	// 		}
	// 		if port, ok := allIP[caType.IP]; ok {
	// 			if _, ok := port[caType.Port]; ok {
	// 				return nil, errors.Errorf("indent Peer:%s prot %d is exist", caType.Name, caType.Port)
	// 			}
	// 		}
	// 		org.CA = caType
	// 		orgType[org.Name] = org
	// 	}
	// } else {
	// 	//目前初版先支持新组织 后期再添加此功能
	// }
	return orgType, nil
}

func checkAddOrgIndent(srcIndent, general *objectdefine.Indent) (map[string]objectdefine.OrgType, error) {
	orgType := make(map[string]objectdefine.OrgType, len(general.Org))
	for _, org := range general.Org {
		if len(org.Name) == 0 {
			return nil, errors.New("Empty buff indent.Org.Name configs")
		}
		//用来保证组织名称首字母大写
		org.Name = tools.InitialsToUpper(org.Name)
		if _, ok := srcIndent.Org[org.Name]; ok {
			return nil, errors.New("new add org is exist")
		}
		if len(org.OrgDomain) == 0 {
			//return nil, errors.New("Empty buff indent.Org.orgDomain configs")
			org.OrgDomain = fmt.Sprintf("%s.example", strings.ToLower(org.Name))
		}
		if len(org.Peer) == 0 {
			return nil, errors.New("Empty buff indent.Org.peer info configs")
		}
		var caIP string //用来如果没有指定ca地址 那就使用主节点的ip 把ca和主节点放到一台机器上
		peerList := make([]objectdefine.PeerType, len(org.Peer))
		for i, one := range org.Peer {
			if len(one.Name) == 0 {
				return nil, errors.Errorf("name Must exist in indent.Org.Peer[%d]", i)
			}
			if len(one.IP) == 0 {
				return nil, errors.Errorf("IP Must exist in indent.Org.Peer.Name=%s", one.Name)
			}
			address := net.ParseIP(one.IP)
			if address == nil {
				return nil, errors.Errorf("IP wrong fromat in indent.Org.Peer.Name=%s", one.Name)
			}
			one.OrgDomain = org.OrgDomain

			if len(one.NickName) == 0 {
				return nil, errors.Errorf("nickname Must exist in indent.Org.Peer.Name=%s", one.Name)
			}
			if len(one.AccessKey) == 0 {
				return nil, errors.Errorf("accesskey Must existin indent.Org.Peer.Name=%s", one.Name)
			}
			if len(one.Org) == 0 {
				one.Org = org.Name
			}
			//需要获取总节点个数用来推送节点状态的peerID
			peerCount, err := dmysql.GetAllPeerNum()
			if err != nil || peerCount == -1 {
				return nil, errors.WithMessage(err, "mysql query peer total count fail")
			}
			one.PeerID = peerCount
			fmt.Println("new add org peerID", one.PeerID)
			if len(one.User) == 0 {
				//数据库查询当前组织有多少个节点
				count, err := dmysql.GetCurrtOrgPeerNum(org.Name, general.ChannelName)
				if count == -1 || err != nil {
					return nil, errors.WithMessage(err, "mysql query fail")
				}
				if count == 0 {
					one.User = "Admin"
				} else {
					one.User = fmt.Sprintf("User%d", count-1)
				}
			}
			if one.User == "Admin" && len(caIP) == 0 {
				caIP = one.IP
			}
			if one.Port < 1 || one.Port > 65535 {
				peerPort, couchdbPort, ccPort, err := dmysql.CheckPeerPort(one.IP)
				if err != nil {
					return nil, errors.Errorf("indent Peer:%s prot %d Out of range", one.Name, one.Port)
				}
				one.Port = peerPort
				one.CouchdbPort = couchdbPort
				one.ChaincodePort = ccPort
			}
			// if port, ok := srcIndent.FireWall[one.IP]; ok {
			// 	if _, ok := port[one.Port]; ok {
			// 		return nil, errors.Errorf("indent Peer:%s prot %d is exist", one.Name, one.Port)
			// 	}
			// }
			// if port, ok := srcIndent.FireWall[one.IP]; ok {
			// 	if _, ok := port[one.CouchdbPort]; ok {
			// 		return nil, errors.Errorf("indent Peer:%s prot %d is exist", one.Name, one.Port)
			// 	}
			// }

			if len(one.Domain) == 0 {
				one.Domain = one.Name + "." + one.OrgDomain
			}
			if len(one.CliName) == 0 {
				one.CliName = fmt.Sprintf("cli-%s-%s", org.Name, one.Name)
			}
			peerList[i] = one
		}
		org.Peer = peerList
		caType := &objectdefine.CAType{}
		//caType := org.CA
		if len(caType.Name) == 0 {
			caType.Name = fmt.Sprintf("ca-%s", org.Name)
		}
		if len(caType.IP) == 0 {
			if len(caIP) == 0 {
				return nil, errors.Errorf("IP Must exist in indent.Org.CA")
			}
			caType.IP = caIP
			address := net.ParseIP(caType.IP)
			if address == nil {
				return nil, errors.Errorf("IP wrong fromat in indent.Org[%s].CA", org.Name)
			}
		}

		if caType.Port < 1 || caType.Port > 65535 {
			caPort, err := dmysql.CheckCAPort(caType.IP)
			if err != nil {
				return nil, errors.Errorf("indent ca:%s prot %d Out of range", caType.Name, caType.Port)
			}
			caType.Port = caPort
		}
		// if port, ok := srcIndent.FireWall[caType.IP]; ok {
		// 	if _, ok := port[caType.Port]; ok {
		// 		return nil, errors.Errorf("indent ca:%s prot %d is exist", caType.Name, caType.Port)
		// 	}
		// }
		org.CA = caType
		orgType[org.Name] = org
	}
	return orgType, nil
}

//checkDeleteOrgIndent 检测补全所有组织 用于签名
func checkDeleteOrgIndent(srcIndent, general *objectdefine.Indent) (map[string]objectdefine.OrgType, error) {
	orgType := make(map[string]objectdefine.OrgType, len(general.Org))
	for _, org := range general.Org {
		if len(org.Name) == 0 {
			return nil, errors.New("Empty buff indent.Org.Name configs")
		}
		//用来保证组织名称首字母大写
		org.Name = tools.InitialsToUpper(org.Name)
		if _, ok := srcIndent.Org[org.Name]; !ok {
			return nil, errors.New("delete org dont exist")

		}
		orgType[org.Name] = srcIndent.Org[org.Name]
	}
	return orgType, nil
}

func checkAddPeerIndent(srcIndent, general *objectdefine.Indent) (map[string]objectdefine.OrgType, error) {
	orgType := make(map[string]objectdefine.OrgType, len(general.Org))
	for _, org := range general.Org {
		if len(org.Name) == 0 {
			return nil, errors.New("Empty buff indent.Org.Name configs")
		}
		if _, ok := srcIndent.Org[org.Name]; !ok {
			return nil, errors.New("indent.Org.Name is not exist")
		}
		if len(org.OrgDomain) == 0 {
			org.OrgDomain = srcIndent.Org[org.Name].OrgDomain
		}
		if len(org.Peer) == 0 {
			return nil, errors.New("Empty buff indent.Org.peer info configs")
		}
		peerList := make([]objectdefine.PeerType, len(org.Peer))
		for i, one := range org.Peer {
			if len(one.IP) == 0 {
				return nil, errors.Errorf("IP Must exist in indent.Org.Peer[%d]", i)
			}
			address := net.ParseIP(one.IP)
			if address == nil {
				return nil, errors.Errorf("IP wrong fromat in indent.Org.Peer.Name=%s", one.Name)
			}
			one.OrgDomain = org.OrgDomain
			if len(one.Name) == 0 {
				return nil, errors.Errorf("name Must exist in indent.Org.Peer[%d]", i)
			}
			//检测同组织节点名称是否相同
			orgAll := srcIndent.Org[org.Name]
			for _, peer := range orgAll.Peer {
				if peer.Name == one.Name {
					return nil, errors.Errorf("alike org[%s] peer name[%s] is same", org.Name, one.Name)
				}
			}
			if len(one.NickName) == 0 {
				return nil, errors.Errorf("nickname Must exist in indent.Org.Peer[%d]", i)
			}
			if len(one.AccessKey) == 0 {
				return nil, errors.Errorf("accesskey Must exist in indent.Org.Peer[%d]", i)
			}
			//需要获取总节点个数用来推送节点状态的peerID
			peerCount, err := dmysql.GetAllPeerNum()
			if err != nil || peerCount == -1 {
				return nil, errors.WithMessage(err, "mysql query peer total count fail")
			}
			one.PeerID = peerCount
			if len(one.User) == 0 {
				//数据库查询当前组织有多少个节点
				count, err := dmysql.GetCurrtOrgPeerNum(org.Name, general.ChannelName)
				if count == -1 || err != nil {
					return nil, errors.WithMessage(err, "mysql query fail")
				}
				if count == 0 {
					one.User = "Admin"
				} else {
					one.User = fmt.Sprintf("User%d", count)
				}
			}
			if one.Port < 1 || one.Port > 65535 {
				peerPort, couchdbPort, ccPort, err := dmysql.CheckPeerPort(one.IP)
				if err != nil {
					return nil, errors.Errorf("indent Peer:%s prot %d Out of range", one.Name, one.Port)
				}
				one.Port = peerPort
				one.CouchdbPort = couchdbPort
				one.ChaincodePort = ccPort
			}

			// if port, ok := srcIndent.FireWall[one.IP]; ok {
			// 	if _, ok := port[one.Port]; ok {
			// 		return nil, errors.Errorf("indent Peer:%s prot %d is exist", one.Name, one.Port)
			// 	}
			// }

			if len(one.Domain) == 0 {
				one.Domain = one.Name + "." + one.OrgDomain
			}
			if len(one.CliName) == 0 {
				one.CliName = fmt.Sprintf("cli-%s-%s", org.Name, one.Name)
			}
			peerList[i] = one
		}
		org.Peer = peerList
		for _, orgS := range srcIndent.Org {
			if orgS.Name == org.Name {
				org.CA = orgS.CA
			}
		}
		orgType[org.Name] = org
	}
	return orgType, nil
}

//checkDeletePeerIndent 检测补全节点
func checkDeletePeerIndent(srcIndent, general *objectdefine.Indent) (map[string]objectdefine.OrgType, error) {
	orgType := make(map[string]objectdefine.OrgType, len(general.Org))
	peerType := make([]objectdefine.PeerType, 0)
	for _, org := range general.Org {
		if len(org.Name) == 0 {
			return nil, errors.New("Empty buff indent.Org.Name configs")
		}
		//用来保证组织名称首字母大写
		org.Name = tools.InitialsToUpper(org.Name)
		if _, ok := srcIndent.Org[org.Name]; !ok {
			return nil, errors.New("delete org dont exist")

		}
		orgInfo := srcIndent.Org[org.Name]
		for _, generalPeer := range org.Peer {
			for _, peer := range orgInfo.Peer {
				if generalPeer.Name == peer.Name {
					peerType = append(peerType, peer)
				}
			}
		}
		org.Peer = peerType
		orgType[org.Name] = org
	}
	return orgType, nil
}

func checkDisablePeerIndent(srcIndent, general *objectdefine.Indent) (map[string]objectdefine.OrgType, error) {
	orgType := make(map[string]objectdefine.OrgType, len(general.Org))
	for _, org := range general.Org {
		if len(org.Name) == 0 {
			return nil, errors.New("Empty buff indent.Org.Name configs")
		}
		if _, ok := srcIndent.Org[org.Name]; !ok {
			return nil, errors.New("indent.Org.Name is not exist")
		}
		if len(org.OrgDomain) == 0 {
			org.OrgDomain = srcIndent.Org[org.Name].OrgDomain
		}
		if len(org.Peer) == 0 {
			return nil, errors.New("Empty buff indent.Org.peer info configs")
		}
		peerList := make([]objectdefine.PeerType, len(org.Peer))
		srcOrg := srcIndent.Org[org.Name]
		for i, one := range org.Peer {
			for _, two := range srcOrg.Peer {
				if len(one.Name) == 0 {
					return nil, errors.Errorf("name Must exist in indent.Org.Peer[%d]", i)
				}
				if len(one.IP) == 0 {
					return nil, errors.Errorf("IP Must exist in indent.Org.Peer[%d]", i)
				}
				address := net.ParseIP(one.IP)
				if address == nil {
					return nil, errors.Errorf("IP wrong fromat in indent.Org.Peer.Name=%s", one.Name)
				}
				if one.Name == two.Name {
					if two.RunStatus == 0 {
						return nil, errors.Errorf("indent.Org.Peer[%d] already disable", i)
					}
					one.OrgDomain = org.OrgDomain
					if one.Port == 0 {
						return nil, errors.Errorf("indent Peer:%s prot %d is empty", one.Name, one.Port)
					}
					if len(one.Domain) == 0 {
						one.Domain = one.Name + "." + one.OrgDomain
					}
					peerList[i] = one
				}
			}
		}
		org.Peer = peerList
		orgType[org.Name] = org
	}
	return orgType, nil
}

func checkEnablePeerIndent(srcIndent, general *objectdefine.Indent) (map[string]objectdefine.OrgType, error) {
	orgType := make(map[string]objectdefine.OrgType, len(general.Org))
	for _, org := range general.Org {
		if len(org.Name) == 0 {
			return nil, errors.New("Empty buff indent.Org.Name configs")
		}
		if _, ok := srcIndent.Org[org.Name]; !ok {
			return nil, errors.New("indent.Org.Name is not exist")
		}
		if len(org.OrgDomain) == 0 {
			org.OrgDomain = srcIndent.Org[org.Name].OrgDomain
		}
		if len(org.Peer) == 0 {
			return nil, errors.New("Empty buff indent.Org.peer info configs")
		}
		peerList := make([]objectdefine.PeerType, len(org.Peer))
		srcOrg := srcIndent.Org[org.Name]
		for i, one := range org.Peer {
			for _, two := range srcOrg.Peer {
				if len(one.Name) == 0 {
					return nil, errors.Errorf("name Must exist in indent.Org.Peer[%d]", i)
				}
				if len(one.IP) == 0 {
					return nil, errors.Errorf("IP Must exist in indent.Org.Peer[%d]", i)
				}
				address := net.ParseIP(one.IP)
				if address == nil {
					return nil, errors.Errorf("IP wrong fromat in indent.Org.Peer.Name=%s", one.Name)
				}
				if one.Name == two.Name {
					if two.RunStatus == 1 {
						return nil, errors.Errorf("indent.Org.Peer[%d] already enable", i)
					}
					one.OrgDomain = org.OrgDomain
					if one.Port == 0 {
						return nil, errors.Errorf("indent Peer:%s prot %d is empty", one.Name, one.Port)
					}
					if len(one.Domain) == 0 {
						one.Domain = one.Name + "." + one.OrgDomain
					}
					peerList[i] = one
				}
			}
		}
		org.Peer = peerList
		orgType[org.Name] = org
	}
	return orgType, nil
}

func checkModiflyPeerIndent(srcIndent, general *objectdefine.Indent) (map[string]objectdefine.OrgType, error) {
	orgType := make(map[string]objectdefine.OrgType, len(general.Org))
	for _, org := range general.Org {
		if len(org.Name) == 0 {
			return nil, errors.New("Empty buff indent.Org.Name configs")
		}
		if _, ok := srcIndent.Org[org.Name]; !ok {
			return nil, errors.New("indent.Org.Name is not exist")
		}
		if len(org.OrgDomain) == 0 {
			org.OrgDomain = srcIndent.Org[org.Name].OrgDomain
		}
		if len(org.Peer) == 0 {
			return nil, errors.New("Empty buff indent.Org.peer info configs")
		}
		peerList := make([]objectdefine.PeerType, len(org.Peer))
		//srcOrg := srcIndent.Org[org.Name]
		for i, one := range org.Peer {
			if len(one.Name) == 0 {
				return nil, errors.Errorf("name Must exist in indent.Org.Peer[%d]", i)
			}
			if len(one.IP) == 0 {
				return nil, errors.Errorf("IP Must exist in indent.Org.Peer[%d]", i)
			}
			address := net.ParseIP(one.IP)
			if address == nil {
				return nil, errors.Errorf("IP wrong fromat in indent.Org.Peer.Name=%s", one.Name)
			}
			if len(one.Domain) == 0 {
				one.Domain = one.Name + "." + one.OrgDomain
			}
			if len(one.NickName) == 0 {
				return nil, errors.Errorf("IP Must exist indent.Org.Peer.NickName")
			}
			peerList[i] = one
		}
		org.Peer = peerList
		orgType[org.Name] = org
	}
	return orgType, nil
}

func checkAddChainCodeIndent(srcIndent, general *objectdefine.Indent) (map[string][]objectdefine.ChainCodeType, error) {
	ccType := make(map[string][]objectdefine.ChainCodeType)

	srcChainCode := srcIndent.Chaincode
	gChainCode := general.Chaincode
	if len(gChainCode) == 0 {
		return nil, errors.New("Empty buff indent chaincode configs")
	}

	for gccName, gcc := range general.Chaincode {

		if len(gccName) == 0 {
			return nil, errors.New("chaincode name empty")
		}
		ccArray := make([]objectdefine.ChainCodeType, 0)
		for _, cc := range gcc {
			if len(cc.Version) == 0 {
				cc.Version = "1.0"
				//return nil, errors.New("chaincode version empty")
			}
			srcCC := srcChainCode[gccName]
			for _, scc := range srcCC {
				if scc.Version == cc.Version {
					if scc.IsInstall != 0 {
						return nil, errors.New("chaincode version is install")
					}
				}
			}
			if len(cc.EndorsementOrg) == 0 {
				return nil, errors.New("chaincode endorsement org  is empty")
			}
			for _, orgName := range cc.EndorsementOrg {
				if _, ok := srcIndent.Org[orgName]; !ok {
					return nil, errors.New("chaincode endorsement org  is not exist")
				}
			}
			endorOrgNum := len(cc.EndorsementOrg)
			endList := make([]string, 0)
			if endorOrgNum == 1 {
				endOrgName := fmt.Sprintf("%sMSP.member", cc.EndorsementOrg[0])
				cc.Policy = "OR('" + endOrgName + "')"
				ccArray = append(ccArray, cc)
			} else {
				for i := 0; i < endorOrgNum; i++ {
					endOrgName := fmt.Sprintf("%sMSP.member", cc.EndorsementOrg[i])
					endList = append(endList, "'"+endOrgName+"'")
				}
				cc.Policy = "OR(" + strings.Join(endList, ",") + ")"
				ccArray = append(ccArray, cc)
			}

		}
		ccType[gccName] = ccArray
	}

	return ccType, nil
}

func checkDeleteChainCodeIndent(srcIndent, general *objectdefine.Indent) (map[string][]objectdefine.ChainCodeType, error) {
	ccType := make(map[string][]objectdefine.ChainCodeType)

	srcChainCode := srcIndent.Chaincode
	gChainCode := general.Chaincode
	if len(gChainCode) == 0 {
		return nil, errors.New("Empty buff indent chaincode configs")
	}

	for gccName, gcc := range general.Chaincode {

		if len(gccName) == 0 {
			return nil, errors.New("chaincode name empty")
		}
		ccArray := make([]objectdefine.ChainCodeType, 0)
		for _, cc := range gcc {
			if len(cc.Version) == 0 {
				return nil, errors.New("chaincode version empty")
			}
			srcCC := srcChainCode[gccName]
			for _, scc := range srcCC {
				if scc.Version == cc.Version {
					if scc.IsInstall == 0 || scc.Status == 3 {
						return nil, errors.New("chaincode version dont install or chaincode being delete")
					}
				}
			}
			if len(cc.EndorsementOrg) == 0 {
				return nil, errors.New("chaincode endorsement org  is empty")
			}
			for _, orgName := range cc.EndorsementOrg {
				if _, ok := srcIndent.Org[orgName]; !ok {
					return nil, errors.New("chaincode endorsement org  is not exist")
				}
			}

			ccArray = append(ccArray, cc)

		}
		ccType[gccName] = ccArray
	}

	return ccType, nil
}

func checkUpgradeChainCodeIndent(srcIndent, general *objectdefine.Indent) (map[string][]objectdefine.ChainCodeType, error) {
	ccType := make(map[string][]objectdefine.ChainCodeType)

	srcChainCode := srcIndent.Chaincode
	gChainCode := general.Chaincode
	if len(gChainCode) == 0 {
		return nil, errors.New("Empty buff indent chaincode configs")
	}

	for gccName, gcc := range general.Chaincode {

		ccArray := make([]objectdefine.ChainCodeType, 0)
		for _, cc := range gcc {
			if len(cc.Version) == 0 {
				return nil, errors.New("chaincode version empty")
			}
			if _, ok := srcChainCode[gccName]; !ok {
				//去检查是否存在新上传的合约
				err := dmysql.CheckUpgradeChaincodeVersionISExist(gccName, cc.Version)
				if err == nil {
					return nil, errors.New("chaincode is not exist")
				}
			}
			srcCC := srcChainCode[gccName]
			if len(srcCC) != 0 {
				for _, scc := range srcCC {
					if scc.Version == cc.Version {
						if scc.IsInstall != 0 {
							return nil, errors.New("chaincode version is install")
						}
					}
				}
			}
			if len(cc.EndorsementOrg) == 0 {
				return nil, errors.New("chaincode endorsement org  is empty")
			}
			for _, orgName := range cc.EndorsementOrg {
				if _, ok := srcIndent.Org[orgName]; !ok {

					return nil, errors.New("chaincode endorsement org  is not exist")
				}
			}
			endorOrgNum := len(cc.EndorsementOrg)
			endList := make([]string, 0)
			if endorOrgNum == 1 {
				endOrgName := fmt.Sprintf("%sMSP.member", cc.EndorsementOrg[0])
				cc.Policy = "OR('" + endOrgName + "')"
				ccArray = append(ccArray, cc)

			} else {
				for i := 0; i < endorOrgNum; i++ {
					endOrgName := fmt.Sprintf("%sMSP.member", cc.EndorsementOrg[i])
					endList = append(endList, "'"+endOrgName+"'")
				}
				cc.Policy = "OR(" + strings.Join(endList, ",") + ")"
				ccArray = append(ccArray, cc)
			}
		}
		ccType[gccName] = ccArray
	}

	return ccType, nil
}

func checkDisableChainCodeIndent(srcIndent, general *objectdefine.Indent) (map[string][]objectdefine.ChainCodeType, error) {
	ccType := make(map[string][]objectdefine.ChainCodeType)

	srcChainCode := srcIndent.Chaincode
	gChainCode := general.Chaincode
	if len(gChainCode) == 0 {
		return nil, errors.New("Empty buff indent chaincode configs")
	}

	for gccName, gcc := range general.Chaincode {
		if _, ok := srcChainCode[gccName]; !ok {
			return nil, errors.New("chaincode is not exist")
		}
		ccArray := make([]objectdefine.ChainCodeType, 0)
		for _, cc := range gcc {
			if len(cc.Version) == 0 {
				return nil, errors.New("chaincode version empty")
			}
			srcCC := srcChainCode[gccName]
			for _, scc := range srcCC {
				if scc.Version == cc.Version {
					if scc.IsInstall == 0 || scc.Status == 0 {
						return nil, errors.New("chaincode version dont install or already disable")
					}
					if len(cc.EndorsementOrg) == 0 {
						cc.EndorsementOrg = scc.EndorsementOrg
					}
					if len(cc.Policy) == 0 {
						cc.Policy = scc.Policy
					}
				}
			}

			ccArray = append(ccArray, cc)
		}
		ccType[gccName] = ccArray
	}

	return ccType, nil
}

func checkEnableChainCodeIndent(srcIndent, general *objectdefine.Indent) (map[string][]objectdefine.ChainCodeType, error) {
	ccType := make(map[string][]objectdefine.ChainCodeType)
	srcChainCode := srcIndent.Chaincode
	gChainCode := general.Chaincode
	if len(gChainCode) == 0 {
		return nil, errors.New("Empty buff indent chaincode configs")
	}
	for gccName, gcc := range general.Chaincode {
		if _, ok := srcChainCode[gccName]; !ok {
			return nil, errors.New("chaincode is not exist")
		}
		ccArray := make([]objectdefine.ChainCodeType, 0)
		for _, cc := range gcc {
			if len(cc.Version) == 0 {
				return nil, errors.New("chaincode version empty")
			}
			srcCC := srcChainCode[gccName]
			for _, scc := range srcCC {
				if scc.Version == cc.Version {
					if scc.IsInstall == 0 || scc.Status == 1 {
						return nil, errors.New("chaincode version not install or already enable")
					}
					if len(cc.EndorsementOrg) == 0 {
						cc.EndorsementOrg = scc.EndorsementOrg
					}
					if len(cc.Policy) == 0 {
						cc.Policy = scc.Policy
					}
				}
			}
			ccArray = append(ccArray, cc)
		}
		ccType[gccName] = ccArray
	}

	return ccType, nil
}

func checkAddServerInfoIndent(general *objectdefine.Indent)(*objectdefine.IndentServer, error){
  
	if len(general.Server.ServerName) == 0 || len(general.Server.ServerExtIp)==0 || len(general.Server.ServerIntIp)==0 ||len(general.Server.ServerUser)==0||len(general.Server.ServerPassword)==0{	
			return nil, errors.New("Missing required parameters")		
	}
	status, err := dmysql.CheckServiceExist(general.Server.ServerExtIp,general.Server.ServerIntIp)
	if !status && err != nil {
		return nil, errors.WithMessage(err, "check server is exist,dont repeat add")
	}
   return general.Server,nil
}