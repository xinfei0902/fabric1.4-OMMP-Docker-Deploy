package daction

import (
	"deploy-server/app/dcache"
	"deploy-server/app/objectdefine"
	"deploy-server/app/tools"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"os/exec"
	//"github.com/hyperledger/fabric/common/channelconfig"
   //"io/ioutil"
	"github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
	"github.com/hyperledger/fabric/protos/orderer/etcdraft"
//"github.com/hyperledger/fabric/common/tools/configtxgen/encoder"
	"github.com/pkg/errors"
	//cb "github.com/hyperledger/fabric/protos/common"
	//pb "github.com/hyperledger/fabric/protos/peer"
	//"github.com/hyperledger/fabric/protos/utils"
)

//MakeStepGeneralCreateCompleteDeployBlock 创建一键部署创世区块
func MakeStepGeneralCreateCompleteDeployBlock(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create block start"},
		})

		err := GeneralCreateCompleteBlock(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
			// result := err
			// err := dmysql.UpdateFailCreateChannelTaskStatus(general)
			// if err != nil {
			// 	output.AppendLog(&objectdefine.StepHistory{
			// 		Name:  "generalCreateChannelTx",
			// 		Error: []string{err.Error()},
			// 	})
			// 	return err
			// }
			// return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create block end"},
		})
		return nil
	}
}

//GeneralCreateCompleteBlock  创建block过程
func GeneralCreateCompleteBlock(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	// templateBlock,err := GetGenesisCompleteBlockTemplate(general, output)
	// if err != nil {
	// 	err = errors.WithMessage(err, "Load template for genesisblock error")
	// 	output.AppendLog(&objectdefine.StepHistory{
	// 		Name:  "generalCreateCompleteDeployBlock",
	// 		Error: []string{err.Error()},
	// 	})
	// 	return err
	// }
	genesis, err := MakeCreatBlockInfo(general,output)
	if err != nil{
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateCompleteDeployBlock",
			Error: []string{err.Error()},
		})
		return err
	}
	// pgen := encoder.New(profileConfig)
	//genesis := CreateGenesisCompleteBlockConfig(general, templateBlock)
	crypto := &objectdefine.CryptoForOrgConfig{}
	crypto = CreateCompleteCryptoForOrgConfig(general)
	err = CreateCompleteCryptoForOrg(general, crypto)
	if err != nil {
		err = errors.WithMessage(err, "Create crypto files config error")
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateCompleteDeployBlock",
			Error: []string{err.Error()},
		})
		return err
	}
	err = CreateCompleteDeployBlock(general, genesis)
	if err != nil {
		err = errors.WithMessage(err, "Create genesisblock config error")
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateCompleteDeployBlock",
			Error: []string{err.Error()},
		})
		return err
	}

	outputPath, calls := CreateGenesisBlockAndCifeExecComanndLine(general)
	err = os.MkdirAll(outputPath, 0777)
	if err != nil {
		err = errors.WithMessage(err, "Prepare path for crypto files error")
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateCompleteDeployBlock",
			Error: []string{err.Error()},
		})
		return err
	}
	logs, err := CallExecAllCommand(general, calls)
	history := &objectdefine.StepHistory{
		Name: "generalCreateCompleteDeployBlock",
		Log:  logs,
	}
	if err != nil {
		history.Error = []string{err.Error()}
	}
	//拷贝bin
	srcBin := dcache.GetBinPathByVersion(general.Version)
	srcBinPath := fmt.Sprintf("%s/configtxgen",srcBin)
	dstBinPath := fmt.Sprintf("%s/deploy/configtxgen",general.BaseOutput)
	tools.CopyFileOne(dstBinPath, srcBinPath)
    //cmd执行
	for _,org := range general.Org{
		outputtxPath:= fmt.Sprintf("%s/deploy/channel-artifacts/%sMSPanchors.tx",general.BaseOutput,org.Name)
		asOrgMSP := fmt.Sprintf("%sMSP",org.Name)
		var cmd *exec.Cmd
		execCommand := fmt.Sprintf("cd %s && chmod +x configtxgen && export FABRIC_CFG_PATH=%s && ./configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate %s -channelID %s -asOrg %s",filepath.ToSlash(filepath.Join(general.BaseOutput, "deploy")),filepath.ToSlash(filepath.Join(general.BaseOutput, "deploy")),outputtxPath,general.ChannelName,asOrgMSP)
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
	}
	
    //证书拷贝
	souridPath := filepath.ToSlash(dcache.GetOutputSubPath("sourceid", ""))
	_,err = os.Stat(souridPath)
	if err != nil {
		err = os.MkdirAll(souridPath, 0777)
		if err != nil {
			err = errors.WithMessage(err, "Prepare path for crypto files error")
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateCompleteDeployBlock",
				Error: []string{err.Error()},
			})
			return err
		}
	}
	src := filepath.Join(general.BaseOutput, "deploy","crypto-config")
	dst := filepath.Join(souridPath, "crypto-config")
	tools.CopyFolder(dst, src)
	
	//status.json
	 vPath := dcache.GetVersionRootPathByVersion(general.Version)
	 ssrc :=filepath.ToSlash(filepath.Join(vPath,"template","deploy","status.json"))
	// fmt.Println("src---",ssrc)
	// err = tools.CopyFileOne(souridPath,ssrc)
	// if err != nil {
	// 	return err
	// }
	var cmd *exec.Cmd
	execCommand := fmt.Sprintf("cp -rf %s %s/",ssrc,souridPath)
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
	// src = filepath.Join(general.BaseOutput, "deploy","crypto-config")
	// dst := filepath.Join(souridPath, "crypto-config")
	// tools.CopyFolder(dst, src)
	//把一键部署通道名称写入文件
	channFilePath :=filepath.ToSlash(dcache.GetOutputSubPath("sourceid", "channelFile.txt"))
	f, err := os.OpenFile(channFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
	   return err
	} else {
	   n, _ := f.Seek(0, os.SEEK_END)
	   _, err = f.WriteAt([]byte(general.ChannelName), n)
	   fmt.Println("write succeed!")
	   defer f.Close()
	}
	output.AppendLog(history)

	return nil
}

//GetGenesisCompleteBlockTemplate  获取模板转结构
func GetGenesisCompleteBlockTemplate(general *objectdefine.Indent, output *objectdefine.TaskNode) (ret *localconfig.TopLevel,err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("Get template of genesis block config error: %v", e)
			ret = nil
		}
	}()
	VersionRoot := dcache.GetTemplatePathByVersion(general.Version)
	path := filepath.Join(VersionRoot, "template", "channel")
	//var profileConfig *localconfig.Profile
    //profileConfig = localconfig.Load("SampleMultiNodeEtcdRaft")
	ret = localconfig.LoadTopLevel(path)
	if len(ret.Organizations) == 0 {  
		err = errors.New("Read template organizations is empty")
		ret = nil
	}
	return
}

//CreateGenesisCompleteBlockConfig  补全构造新的 configtx.yaml文件
func CreateGenesisCompleteBlockConfig(general *objectdefine.Indent, template *localconfig.TopLevel) (ret *localconfig.TopLevel) {
	ret = &localconfig.TopLevel{}
	ret.Organizations = genesisCompleteBlockOrganizations(general, template)
	ret.Capabilities = template.Capabilities
	ret.Application = template.Application
	ret.Application.Organizations = ret.Organizations
	ret.Application.Policies = template.Application.Policies
	ret.Application.Capabilities =ret.Capabilities["application"]
	ret.Channel = template.Channel
	// ret.Channel.Policies = template.Channel.Policies
	// ret.Channel.Capabilities = ret.Capabilities["channel"]
	ret.Orderer = genesisCompleteBlockOrderer(general, template, ret.Organizations)
	ret.Profiles = genesisCompleteBlockProfiles(general, template, ret)

	return ret
}

func genesisCompleteBlockOrganizations(general *objectdefine.Indent, template *localconfig.TopLevel) []*localconfig.Organization {
	ret := make([]*localconfig.Organization, 0, len(general.Org)+1)

	for _, orderer := range general.Orderer {
		OrganizationsTemp := &localconfig.Organization{}
		OrganizationsTemp.ID = "OrdererMSP"
		OrganizationsTemp.Name = "OrdererOrg"
		OrganizationsTemp.AdminPrincipal = template.Organizations[0].AdminPrincipal
		OrganizationsTemp.MSPType = template.Organizations[0].MSPType
		OrganizationsTemp.Policies = make(map[string]*localconfig.Policy)
		for k, v := range template.Organizations[0].Policies {
			one := &localconfig.Policy{}
			one.Type = v.Type
			one.Rule = v.Rule
			OrganizationsTemp.Policies[k] = one
		}
		mspDir := fmt.Sprintf("crypto-config/ordererOrganizations/%s/msp", orderer.OrgDomain)
		OrganizationsTemp.MSPDir = mspDir

		for k, v := range OrganizationsTemp.Policies {
			v.Rule = strings.Replace(v.Rule, "$[org]", "OrdererMSP", -1)
			OrganizationsTemp.Policies[k] = v
		}
		ret = append(ret, OrganizationsTemp)
		break
	}
	
	for _, orgObj := range general.Org {
		if len(orgObj.Peer) == 0 {
			continue
		}

		OrganizationsTemp := &localconfig.Organization{}
		orgNameMSP := fmt.Sprintf("%sMSP", orgObj.Name)
		OrganizationsTemp.ID = orgNameMSP
		OrganizationsTemp.Name = orgNameMSP
		OrganizationsTemp.AdminPrincipal = template.Organizations[1].AdminPrincipal
		OrganizationsTemp.MSPType = template.Organizations[1].MSPType
		OrganizationsTemp.Policies = make(map[string]*localconfig.Policy)
		for k, v := range template.Organizations[1].Policies {
			one := &localconfig.Policy{}
			one.Type = v.Type
			one.Rule = v.Rule
			OrganizationsTemp.Policies[k] = one
		}
		mspDir := fmt.Sprintf("crypto-config/peerOrganizations/%s/msp", orgObj.OrgDomain)
		OrganizationsTemp.MSPDir = mspDir

		for k, v := range OrganizationsTemp.Policies {
			v.Rule = strings.Replace(v.Rule, "$[org]", orgNameMSP, -1)
			OrganizationsTemp.Policies[k] = v
		}

		anchors := make([]*localconfig.AnchorPeer, 0, len(orgObj.Peer))
		for _, one := range orgObj.Peer {
			if "Admin" != one.User {
				continue
			}
			anchors = append(anchors, &localconfig.AnchorPeer{
				Host: one.Domain,
				Port: one.Port,
			})
		}

		if len(anchors) > 0 {
			OrganizationsTemp.AnchorPeers = anchors
		}

		ret = append(ret, OrganizationsTemp)
	}

	return ret
}

func genesisCompleteBlockOrderer(general *objectdefine.Indent, template *localconfig.TopLevel, orgs []*localconfig.Organization) *localconfig.Orderer {

	one := *template.Orderer
	one.OrdererType = general.Consensus
	one.Addresses = make([]string, 0, 32)
	for _, orderer := range general.Orderer {
		url := fmt.Sprintf("%s:%d", orderer.Domain, orderer.Port)
		one.Addresses = append(one.Addresses, url)
	}
	one.BatchTimeout = template.Orderer.BatchTimeout
	one.BatchSize = template.Orderer.BatchSize

	if one.OrdererType == objectdefine.ConsensusKafka && general.Kafka != nil {

		pair := make([]string, 0, len(general.Kafka.Kafka))
		for _, one := range general.Kafka.Kafka {
			pair = append(pair, fmt.Sprintf("%s:%d", one.Domain, one.Port))
		}
		one.Kafka.Brokers = pair
	}else{
		one.Kafka = template.Orderer.Kafka
	}
	Consenters := make([]*etcdraft.Consenter,0)
	Addresses := make([]string,0)
	for _, orderer := range general.Orderer{
		consenter := &etcdraft.Consenter{
			Host: orderer.Domain,
			Port: uint32(orderer.Port),
			ClientTlsCert: []byte(fmt.Sprintf("crypto-config/ordererOrganizations/example.com/orderers/%s/tls/server.crt",orderer.Domain)),
		    ServerTlsCert: []byte(fmt.Sprintf("crypto-config/ordererOrganizations/example.com/orderers/%s/tls/server.crt",orderer.Domain)),
		}
		Addresses = append(Addresses,fmt.Sprintf("%s:%d",orderer.Domain,orderer.Port))
		Consenters= append(Consenters,consenter)
	}
	one.EtcdRaft.Consenters = Consenters
	one.Organizations = template.Orderer.Organizations
	one.Policies = template.Orderer.Policies
	return &one
}

func genesisCompleteBlockProfiles(general *objectdefine.Indent, template *localconfig.TopLevel, self *localconfig.TopLevel) map[string]*localconfig.Profile {
	ret := make(map[string]*localconfig.Profile)
	//创世区块的联盟配置
	genesis := &localconfig.Profile{}
	genesis.Capabilities = make(map[string]bool,5)
	genesis.Capabilities = self.Channel.Capabilities
	genesis.Orderer = self.Orderer
	ordererOrganizations := make([]*localconfig.Organization, 1)
	ordererOrganizations[0] = self.Organizations[0]
	//raft
	genesis.Orderer.OrdererType = "etcdraft"
	Consenters := make([]*etcdraft.Consenter,0)
	Addresses := make([]string,0)
	for _, orderer := range general.Orderer{
		consenter := &etcdraft.Consenter{
			Host: orderer.Domain,
			Port: uint32(orderer.Port),
			ClientTlsCert: []byte(fmt.Sprintf("crypto-config/ordererOrganizations/example.com/orderers/%s/tls/server.crt",orderer.Domain)),
		    ServerTlsCert: []byte(fmt.Sprintf("crypto-config/ordererOrganizations/example.com/orderers/%s/tls/server.crt",orderer.Domain)),
		}
		Addresses = append(Addresses,fmt.Sprintf("%s:%d",orderer.Domain,orderer.Port))
		Consenters= append(Consenters,consenter)
	}
	genesis.Orderer.EtcdRaft.Consenters = Consenters
	genesis.Orderer.Addresses = Addresses
	genesis.Orderer.Organizations = ordererOrganizations
	genesis.Orderer.Capabilities = make(map[string]bool,5)
	genesis.Orderer.Capabilities = self.Capabilities["orderer"]
	genesis.Application = self.Application
	genesis.Application.Organizations = ordererOrganizations
	genesis.Consortiums = make(map[string]*localconfig.Consortium)
	genesisOrganizations := make([]*localconfig.Organization, len(general.Org))

	for i := 0; i < len(general.Org); i++ {
		genesisOrganizations[i] = self.Organizations[i+1]
	}
	// genesis.Consortiums = map[string]*localconfig.Consortium{
	// 	"PrisonConsortium": &localconfig.Consortium{
	// 		Organizations: genesisOrganizations,
	// 	},
	// }
	genesis.Consortiums = map[string]*localconfig.Consortium{
		"SampleConsortium": &localconfig.Consortium{
			Organizations: genesisOrganizations,
		},
	}
	//ret["PrisonNetworkGenesis"] = genesis
	ret["SampleMultiNodeEtcdRaft"] = genesis
	channel := &localconfig.Profile{}
	channel.Consortium = "SampleConsortium"
	channel.Application = self.Application
	channel.Application.Organizations = genesisOrganizations
	channel.Application.Capabilities = make(map[string]bool,5)
	channel.Application.Capabilities = self.Capabilities["application"]

	ret["TwoOrgsChannel"] = channel
	//每一个通道的配置
	if len(general.Deploy) == 0{

	}else{
		for channelName, deploy := range general.Deploy {
			channelConfig := genesisCompleteBlockChannelConfig(deploy, self)
			ret[channelName] = channelConfig
		}
	}
	

	return ret
}

//genesisCompleteBlockChannelConfig  config.yaml 多通道配置
func genesisCompleteBlockChannelConfig(deploy objectdefine.DeployType, self *localconfig.TopLevel) *localconfig.Profile {
	cgenesosOrganizations := make([]*localconfig.Organization, 0)
	for _, org := range deploy.JoinOrg {
		orgName := fmt.Sprintf("%sMSP", org.Name)
		for _, configOrg := range self.Organizations {
			if orgName == configOrg.ID {
				cgenesosOrganizations = append(cgenesosOrganizations, configOrg)
			}
		}
	}
	channel := &localconfig.Profile{}
	channel.Consortium = "PrisonConsortium"
	channel.Application = &localconfig.Application{}
	channel.Application.Organizations = cgenesosOrganizations
	channel.Application.Capabilities = self.Capabilities["Application"]
	return channel
}

//CreateCompleteCryptoForOrgConfig 创建生成证书的crypto-config.yaml文件
func CreateCompleteCryptoForOrgConfig(general *objectdefine.Indent) *objectdefine.CryptoForOrgConfig {
	ret := &objectdefine.CryptoForOrgConfig{}

	ret.OrdererOrgs = make([]objectdefine.OrgSpec, 0, 1)
	ret.PeerOrgs = make([]objectdefine.OrgSpec, 0, len(general.Org))
	for _, one := range general.Org {
		if len(one.Peer) == 0 {
			continue
		}

		two := objectdefine.OrgSpec{}
		two.Name = one.Name
		two.Domain = one.OrgDomain
		two.EnableNodeOUs = true
		two.Template = objectdefine.NodeTemplate{
			Count: len(one.Peer),
		}
		two.Users = objectdefine.UsersSpec{

			Count: len(one.Peer) - 1,
		}
		ret.PeerOrgs = append(ret.PeerOrgs, two)
	}
	two := objectdefine.OrgSpec{}
	var ordererName string
	var ordererDomain string
	two.Specs = make([]objectdefine.NodeSpec, len(general.Orderer))
	for i, orderer := range general.Orderer {
		if len(ordererName) == 0 {
			ordererDomain = orderer.OrgDomain
		}
		//hostName := fmt.Sprintf("%s:%d", orderer.Domain, orderer.Port)
		two.Specs[i].Hostname = strings.ToLower(orderer.Name)
	}

	two.Name = "Orderer"
	two.Domain = ordererDomain
	two.EnableNodeOUs = true
	ordererCount := len(general.Orderer)
	two.Users = objectdefine.UsersSpec{
		Count: ordererCount,
	}
	ret.OrdererOrgs = append(ret.OrdererOrgs, two)
	return ret
}

//CreateCompleteCryptoForOrg 创建生成证书的配置文件
func CreateCompleteCryptoForOrg(general *objectdefine.Indent, config *objectdefine.CryptoForOrgConfig) (err error) {
	err = os.MkdirAll(general.BaseOutput, 0777)
	if nil != err {
		err = errors.WithMessage(err, "Make Output folder error")
		return
	}
	certPath := filepath.Join(general.BaseOutput,"deploy")
	_,err = os.Stat(certPath)
	if err != nil{
		err = os.MkdirAll(certPath, 0777)
		if nil != err {
			err = errors.WithMessage(err, "Make Output folder error")
			return
		}
	}
   
	filename := filepath.Join(certPath, "crypto-config.yaml")
	err = tools.WriteYamlFile(filename, config)
	if nil != err {
		err = errors.WithMessage(err, "Write crypto for org config error")
		return
	}

	// src := filepath.Join(general.SourceBaseOutput, "crypto-config", "ordererOrganizations")
	// dst := filepath.Join(general.BaseOutput, "crypto-config", "ordererOrganizations")
	// _, err = os.Stat(dst)
	// if err != nil {
	// 	err = os.MkdirAll(dst, 0777)
	// 	if nil != err {
	// 		err = errors.WithMessage(err, "Make Output folder error")
	// 		return
	// 	}
	// }
	// tools.CopyFolder(dst, src)

	// V2
	// task.TargetPath /--
	//                  |- SubPathOfCryptoForOrgConfig
	//                  |- SubPathOfGenesisBlockConfig
	//                  |- channel-artifacts /--
	//                                        |- SubPathOfGenesisBlock
	//                                        |- ChannelName.tx

	return nil
}

//CreateCompleteDeployBlock 创建configtx.yaml文件
func CreateCompleteDeployBlock(general *objectdefine.Indent, config *localconfig.TopLevel) (err error) {
	// write config
	// write script
	// call script
	_,err = os.Stat(general.BaseOutput)
	if err != nil{
		err = os.MkdirAll(general.BaseOutput, 0777)
		if nil != err {
			err = errors.WithMessage(err, "Make Output folder error")
			return
		}
	}
	
	configtPath := filepath.Join(general.BaseOutput,"deploy")
	filename := filepath.Join(configtPath, "configtx.yaml")
	err = tools.WriteYamlFile(filename, config)
	if nil != err {
		err = errors.WithMessage(err, "Write genesis config error")
		return
	}

	return nil
}

//CreateGenesisBlockAndCifeExecComanndLine 构造目录 组装生成创世区块和证书命令
func CreateGenesisBlockAndCifeExecComanndLine(general *objectdefine.Indent) (outputPath string, ret []*objectdefine.CommandList) {

	outputPath = filepath.Join(general.BaseOutput, "deploy","channel-artifacts")

	readconfig := filepath.Join(general.BaseOutput, "deploy","crypto-config.yaml")
	var env []string
	env = []string{fmt.Sprintf("FABRIC_CFG_PATH=%s", filepath.ToSlash(filepath.Join(general.BaseOutput, "deploy")))}

    os.Setenv("FABRIC_CFG_PATH",filepath.ToSlash(filepath.Join(general.BaseOutput, "deploy")))
	var cmd *exec.Cmd
	execCommand := fmt.Sprintf("export FABRIC_CFG_PATH=%s",filepath.ToSlash(filepath.Join(general.BaseOutput, "deploy")))
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", execCommand)
	} else {
		cmd = exec.Command("/bin/bash", "-c", execCommand)
	}
	if err := cmd.Start(); err != nil {
		return "",nil
	}
	if err := cmd.Wait(); err != nil {
		return "",nil
	}
	ret = make([]*objectdefine.CommandList, 0)

	toolsPath := dcache.GetBinPathByVersion(general.Version)

	pair := dcache.GetBinToolsPair()

	bins, ok := pair[runtime.GOOS]
	if !ok {
		return "", nil
	}
	one := &objectdefine.CommandList{}
	one.OS = runtime.GOOS
	one.Call = MakeGenerateExecCommand(env, general, outputPath, readconfig,
		filepath.Join(toolsPath, bins.CryptoGen),
		filepath.Join(toolsPath, bins.ConfigGen))

	ret = append(ret, one)

	return
}

//MakeGenerateExecCommand 组装生成.tx文件命令
func MakeGenerateExecCommand(env []string, general *objectdefine.Indent, output, readconfig, cryptogen, configtxgen string) []objectdefine.ProcessPair {
	ChannelName := general.ChannelName

	outputBlocks := filepath.Join(output, "genesis.block")
    baseOut := filepath.Join(general.BaseOutput, "deploy")
	ret := make([]objectdefine.ProcessPair, 0)

	//创建证书命令
	ret = append(ret, objectdefine.ProcessPair{
		Exec:        cryptogen,
		Args:        []string{`generate`, `--config`, readconfig},
		Dir:         baseOut,
		Environment: env,
	})

	//block file 用来创建创世区块的
	ret = append(ret, objectdefine.ProcessPair{
		Exec:        configtxgen,
		Args:        []string{`-profile`, `SampleMultiNodeEtcdRaft`, `-channelID`, `byfn-sys-channel`, `-outputBlock`, outputBlocks},
		Dir:         baseOut,
		Environment: env,
	})
	// ret = append(ret, objectdefine.ProcessPair{
	// 	Exec:        configtxgen,
	// 	Args:        []string{`-profile`, `samplemultinodeetcdraft`, `-channelID`, `byfn-sys-channel`, `-outputBlock`, outputBlocks},
	// 	Dir:         baseOut,
	// 	Environment: env,
	// })
	//通道 ***.tx
	outputChannelTx := filepath.Join(output, ChannelName+".tx")
	ret = append(ret, objectdefine.ProcessPair{
		Exec:        configtxgen,
		Args:        []string{`-profile`, `TwoOrgsChannel`, `-channelID`, ChannelName, `-outputCreateChannelTx`, outputChannelTx},
		Dir:         baseOut,
		Environment: env,
	}) 
	// outputChannelTx := filepath.Join(output, ChannelName+".tx")
	// ret = append(ret, objectdefine.ProcessPair{
	// 	Exec:        configtxgen,
	// 	Args:        []string{`-profile`, `twoorgschannel`, `-channelID`, ChannelName, `-outputCreateChannelTx`, outputChannelTx},
	// 	Dir:         baseOut,
	// 	Environment: env,
	// })
	// if len(general.Org) == 0{
	// 	fmt.Println("exec ancher tx create nil")
	// }
	// for _,org :=range general.Org{
	// 	orgMsp := fmt.Sprintf("%sMSP", org.Name)
	// 	ancherName := filepath.Join(output, fmt.Sprintf("%sMSPanchors.tx", org.Name))
	// 	ret = append(ret, objectdefine.ProcessPair{
	// 		Exec:        configtxgen,
	// 		Args:        []string{`-profile`, `TwoOrgsChannel `, `-channelID`, ChannelName, `-asOrg`, orgMsp, `-outputAnchorPeersUpdate`, ancherName},
	// 		Dir:         baseOut,
	// 		Environment: env,
	// 	})
	// }

	// for channelName, deploy := range general.Deploy {
	// 	// channel file
	// 	outputChannelTx := filepath.Join(output, channelName+".tx")
	// 	ret = append(ret, objectdefine.ProcessPair{
	// 		Exec:        configtxgen,
	// 		Args:        []string{`-profile`, channelName, `-channelID`, channelName, `-outputCreateChannelTx`, outputChannelTx},
	// 		Dir:         general.BaseOutput,
	// 		Environment: env,
	// 	})
	// 	for _, org := range deploy.JoinOrg {
	// 		if len(org.Peer) == 0 {
	// 			continue
	// 		}
	// 		for _, peer := range org.Peer {
	// 			if "Admin" != peer.User {
	// 				continue
	// 			}
	// 			orgMsp := fmt.Sprintf("%sMSP", org.Name)
	// 			ret = append(ret, objectdefine.ProcessPair{
	// 				Exec:        configtxgen,
	// 				Args:        []string{`-profile`, channelName, `-channelID`, channelName, `-asOrg`, orgMsp, `-outputAnchorPeersUpdate`, filepath.Join(output, fmt.Sprintf("%sMSPanchors@%s.tx", org.Name, channelName))},
	// 				Dir:         general.BaseOutput,
	// 				Environment: env,
	// 			})
	// 		}
	// 	}
	// }
	return ret
}

func MakeGenerateCryptogenExecCommand(env []string, general *objectdefine.Indent, readconfig, cryptogen string) []objectdefine.ProcessPair {
	ret := make([]objectdefine.ProcessPair, 0)
	//创建证书命令
	ret = append(ret, objectdefine.ProcessPair{
		Exec:        cryptogen,
		Args:        []string{`extend`, `--config`, readconfig},
		Dir:         general.SourceBaseOutput,
		Environment: env,
	})
	return ret
}

//CallExecAllCommand 运行命令 生成.tx文件
func CallExecAllCommand(general *objectdefine.Indent, list []*objectdefine.CommandList) (logs []string, err error) {

	match := -1
	notWindows := -1
	for j, one := range list {
		if one.OS == runtime.GOOS {
			match = j
			break
		}
		if notWindows < 0 && one.OS != `windows` {
			notWindows = j
		}
	}

	if match < 0 {
		if notWindows < 0 {
			return nil, errors.New("no match script for channel artifacts")
		}
		match = notWindows
	}

	call := list[match]
	logs, err = call.Run()
	return
}





//MakeStepGeneralCreateCompleteDeployBlock 创建一键部署创世区块
func MakeStepGeneralCreateCompleteDeployBlockAnther(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create block start"},
		})

		err := GeneralCreateCompleteBlockAnther(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create block end"},
		})
		return nil
	}
}

//GeneralCreateCompleteBlock  创建block过程
func GeneralCreateCompleteBlockAnther(general *objectdefine.Indent, output *objectdefine.TaskNode) error {

	crypto := &objectdefine.CryptoForOrgConfig{}
	crypto = CreateCompleteCryptoForOrgConfig(general)
	err := CreateCompleteCryptoForOrg(general, crypto)
	if err != nil {
		err = errors.WithMessage(err, "Create crypto files config error")
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateCompleteDeployBlock",
			Error: []string{err.Error()},
		})
		return err
	}

	// err = CreateCompleteDeployBlock(general, genesis)
	// if err != nil {
	// 	err = errors.WithMessage(err, "Create genesisblock config error")
	// 	output.AppendLog(&objectdefine.StepHistory{
	// 		Name:  "generalCreateCompleteDeployBlock",
	// 		Error: []string{err.Error()},
	// 	})
	// 	return err
	// }

	_, calls := CreateGenesisBlockAndCifeExecComanndLineAnther(general)
	
	logs, err := CallExecAllCommand(general, calls)
	history := &objectdefine.StepHistory{
		Name: "generalCreateCompleteDeployBlock",
		Log:  logs,
	}
	if err != nil {
		history.Error = []string{err.Error()}
	}

    _,err = MakeCreatBlockInfo(general,output)
   if err != nil{
	output.AppendLog(&objectdefine.StepHistory{
		Name:  "generalCreateCompleteDeployBlock",
		Error: []string{err.Error()},
	})
   }
	//拷贝bin
	srcBin := dcache.GetBinPathByVersion(general.Version)
	srcBinPath := fmt.Sprintf("%s/configtxgen",srcBin)
	dstBinPath := fmt.Sprintf("%s/deploy/configtxgen",general.BaseOutput)
	tools.CopyFileOne(dstBinPath, srcBinPath)
    //cmd执行
	for _,org := range general.Org{
		outputtxPath:= fmt.Sprintf("%s/deploy/channel-artifacts/%sMSPanchors.tx",general.BaseOutput,org.Name)
		asOrgMSP := fmt.Sprintf("%sMSP",org.Name)
		var cmd *exec.Cmd
		execCommand := fmt.Sprintf("cd %s && chmod +x configtxgen && export FABRIC_CFG_PATH=%s && ./configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate %s -channelID %s -asOrg %s",filepath.ToSlash(filepath.Join(general.BaseOutput, "deploy")),filepath.ToSlash(filepath.Join(general.BaseOutput, "deploy")),outputtxPath,general.ChannelName,asOrgMSP)
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
	}
	
    //证书拷贝
	souridPath := filepath.ToSlash(dcache.GetOutputSubPath("sourceid", ""))
	_,err = os.Stat(souridPath)
	if err != nil {
		err = os.MkdirAll(souridPath, 0777)
		if err != nil {
			err = errors.WithMessage(err, "Prepare path for crypto files error")
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreateCompleteDeployBlock",
				Error: []string{err.Error()},
			})
			return err
		}
	}
	src := filepath.Join(general.BaseOutput, "deploy","crypto-config")
	dst := filepath.Join(souridPath, "crypto-config")
	tools.CopyFolder(dst, src)
	
	//status.json
	 vPath := dcache.GetVersionRootPathByVersion(general.Version)
	 ssrc :=filepath.ToSlash(filepath.Join(vPath,"template","deploy","status.json"))
	// fmt.Println("src---",ssrc)
	// err = tools.CopyFileOne(souridPath,ssrc)
	// if err != nil {
	// 	return err
	// }
	var cmd *exec.Cmd
	execCommand := fmt.Sprintf("cp -rf %s %s/",ssrc,souridPath)
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
	// src = filepath.Join(general.BaseOutput, "deploy","crypto-config")
	// dst := filepath.Join(souridPath, "crypto-config")
	// tools.CopyFolder(dst, src)
	 output.AppendLog(history)

	return nil
}

func MakeCreatBlockInfo(general *objectdefine.Indent,output *objectdefine.TaskNode)(*localconfig.TopLevel,error){
	
	VersionRoot := dcache.GetTemplatePathByVersion(general.Version)
	path := filepath.Join(VersionRoot, "template", "channel")
	
	template := localconfig.LoadTopLevel(path)
	if len(template.Organizations) == 0 {
		
	}
	Capabilities := make(map[string]map[string]bool)

	channelCap := make(map[string]bool)
	channelCap["V1_4_3"] = true
	channelCap["V1_3"] = false
	channelCap["V1_1"] = false
	Capabilities["Channel"] =channelCap
    
	ordererCap := make(map[string]bool)
	ordererCap["V1_4_2"] = true
	ordererCap["V1_1"] = false
	Capabilities["Orderer"] =ordererCap

	applicationCap := make(map[string]bool)
	applicationCap["V1_4_2"] = true
	applicationCap["V1_3"] = false
	applicationCap["V1_2"] = false
	applicationCap["V1_1"] = false
	Capabilities["Application"] =applicationCap

	OrganizationsArray := make([]*localconfig.Organization,0)
	for _, orderer := range general.Orderer {
		OrganizationsTemp := &localconfig.Organization{}
		OrganizationsTemp.ID = "OrdererMSP"
		OrganizationsTemp.Name = "OrdererOrg"
		OrganizationsTemp.AdminPrincipal = template.Organizations[0].AdminPrincipal
		OrganizationsTemp.MSPType = template.Organizations[0].MSPType
		OrganizationsTemp.Policies = make(map[string]*localconfig.Policy)
		for k, v := range template.Organizations[0].Policies {
			one := &localconfig.Policy{}
			one.Type = v.Type
			one.Rule = v.Rule
			OrganizationsTemp.Policies[k] = one
		}
		mspDir := fmt.Sprintf("crypto-config/ordererOrganizations/%s/msp", orderer.OrgDomain)
		OrganizationsTemp.MSPDir = mspDir

		for k, v := range OrganizationsTemp.Policies {
			v.Rule = strings.Replace(v.Rule, "$[org]", "OrdererMSP", -1)
			OrganizationsTemp.Policies[k] = v
		}
		OrganizationsArray = append(OrganizationsArray, OrganizationsTemp)
		break
	}
	
		//ret := &localconfig.TopLevel{}
		
	for _, orgObj := range general.Org {
		if len(orgObj.Peer) == 0 {
			continue
		}

		OrganizationsTemp := &localconfig.Organization{}
		orgNameMSP := fmt.Sprintf("%sMSP", orgObj.Name)
		OrganizationsTemp.ID = orgNameMSP
		OrganizationsTemp.Name = orgNameMSP
		sss := template.Organizations[1].AdminPrincipal
		fmt.Println("ssss",sss)
		OrganizationsTemp.AdminPrincipal = template.Organizations[1].AdminPrincipal
		OrganizationsTemp.MSPType = template.Organizations[1].MSPType
		OrganizationsTemp.Policies = make(map[string]*localconfig.Policy)
		for k, v := range template.Organizations[1].Policies {
			one := &localconfig.Policy{}
			one.Type = v.Type
			one.Rule = v.Rule
			OrganizationsTemp.Policies[k] = one
		}
		mspDir := fmt.Sprintf("crypto-config/peerOrganizations/%s/msp", orgObj.OrgDomain)
		OrganizationsTemp.MSPDir = mspDir

		for k, v := range OrganizationsTemp.Policies {
			v.Rule = strings.Replace(v.Rule, "$[org]", orgNameMSP, -1)
			OrganizationsTemp.Policies[k] = v
		}

		anchors := make([]*localconfig.AnchorPeer, 0, len(orgObj.Peer))
		for _, one := range orgObj.Peer {
			if "Admin" != one.User {
				continue
			}
			anchors = append(anchors, &localconfig.AnchorPeer{
				Host: one.Domain,
				Port: one.Port,
			})
		}

		if len(anchors) > 0 {
			OrganizationsTemp.AnchorPeers = anchors
		}

		OrganizationsArray = append(OrganizationsArray, OrganizationsTemp)
	}
		
	template.Organizations = OrganizationsArray
	
     //template.Channel 模块
	 template.Channel.Capabilities =Capabilities["Channel"]
	 for name,value:= range template.Channel.Policies{
		 if name =="admins"{
			 template.Channel.Policies["Admins"] = value
			delete(template.Channel.Policies,"admins")
		 }
		 
		if name == "readers"{
		 template.Channel.Policies["Readers"]=value
			delete(template.Channel.Policies,"readers")
		}
 
		if name == "writers"{
		 template.Channel.Policies["Writers"]=value
			delete(template.Channel.Policies,"writers")
		}
	}
	
	//template.Application 模块
	template.Application.Capabilities =Capabilities["Application"]
	for name,value:= range template.Application.Policies{
		if name =="admins"{
			template.Application.Policies["Admins"] = value
		   delete(template.Application.Policies,"admins")
		}
		
	   if name == "readers"{
		template.Application.Policies["Readers"]=value
		   delete(template.Application.Policies,"readers")
	   }
 
	   if name == "writers"{
		template.Application.Policies["Writers"]=value
		   delete(template.Application.Policies,"writers")
	   }
   }
	//template.Orderer = template.Orderer
	template.Orderer.OrdererType = "etcdraft"
	Addresses := make([]string, 0, 32)
	for _, orderer := range general.Orderer {
		url := fmt.Sprintf("%s:7050", orderer.Domain)
		Addresses = append(Addresses, url)
	}
	template.Orderer.Addresses = Addresses
	Consenters := make([]*etcdraft.Consenter,0)
	Addressesss := make([]string,0)
	for _, orderer := range general.Orderer{
		consenter := &etcdraft.Consenter{
			Host: orderer.Domain,
			Port: uint32(orderer.Port),
			ClientTlsCert: []byte(fmt.Sprintf("crypto-config/ordererOrganizations/example.com/orderers/%s/tls/server.crt",orderer.Domain)),
			ServerTlsCert: []byte(fmt.Sprintf("crypto-config/ordererOrganizations/example.com/orderers/%s/tls/server.crt",orderer.Domain)),
		}
		Addressesss = append(Addressesss,fmt.Sprintf("%s:%d",orderer.Domain,orderer.Port))
		Consenters= append(Consenters,consenter)
	}
	template.Orderer.EtcdRaft.Consenters = Consenters
	
    //template.Orderer.Capabilities =Capabilities["Application"]
	for name,value:= range template.Orderer.Policies{
		if name =="admins"{
			template.Orderer.Policies["Admins"] = value
			delete(template.Orderer.Policies,"admins")
		}
		
		if name == "readers"{
		template.Orderer.Policies["Readers"]=value
			delete(template.Orderer.Policies,"readers")
		}

		if name == "writers"{
		template.Orderer.Policies["Writers"]=value
			delete(template.Orderer.Policies,"writers")
		}
		if name == "blockvalidation"{
		    template.Orderer.Policies["BlockValidation"]=value
			delete(template.Orderer.Policies,"blockvalidation")
		}	
		
	}

	////template.Capabilities 模块
	template.Capabilities = Capabilities
	//profile
	template.Profiles["samplemultinodeetcdraft"].Application.Organizations[0] =  template.Organizations[0]
	template.Profiles["samplemultinodeetcdraft"].Orderer.Addresses = Addresses
	template.Profiles["samplemultinodeetcdraft"].Orderer.Addresses = Addresses
	template.Profiles["samplemultinodeetcdraft"].Orderer.EtcdRaft.Consenters = Consenters
	template.Profiles["samplemultinodeetcdraft"].Orderer.Organizations[0] =  template.Organizations[0]
	
	genesisOrganizations := make([]*localconfig.Organization, len(general.Org))
	
	for i := 0; i < len(general.Org); i++ {
		genesisOrganizations[i] = template.Organizations[i+1]
	}
	template.Profiles["samplemultinodeetcdraft"].Consortiums["sampleconsortium"].Organizations = genesisOrganizations
	
	template.Profiles["twoorgschannel"].Application.Organizations =genesisOrganizations
	template.Profiles["twoorgschannel"].Application.Organizations =genesisOrganizations
	for name,value:= range template.Profiles{
		if name =="samplemultinodeetcdraft"{
            //application
			template.Profiles["samplemultinodeetcdraft"].Application.Capabilities = Capabilities["Application"]
            //policies
			for name,value:= range template.Profiles["samplemultinodeetcdraft"].Application.Policies{
				if name =="admins"{
					template.Profiles["samplemultinodeetcdraft"].Application.Policies["Admins"] = value
					delete(template.Profiles["samplemultinodeetcdraft"].Application.Policies,"admins")
				}
				
				if name == "readers"{
					template.Profiles["samplemultinodeetcdraft"].Application.Policies["Readers"]=value
					delete(template.Profiles["samplemultinodeetcdraft"].Application.Policies,"readers")
				}
		
				if name == "writers"{
				template.Profiles["samplemultinodeetcdraft"].Application.Policies["Writers"]=value
					delete(template.Profiles["samplemultinodeetcdraft"].Application.Policies,"writers")
				}
			}

			//orderer
			template.Profiles["samplemultinodeetcdraft"].Orderer.Capabilities = Capabilities["Orderer"]
            //policies
			for name,value:= range template.Profiles["samplemultinodeetcdraft"].Orderer.Policies{
				if name =="admins"{
					template.Profiles["samplemultinodeetcdraft"].Orderer.Policies["Admins"] = value
					delete(template.Profiles["samplemultinodeetcdraft"].Orderer.Policies,"admins")
				}
				
				if name == "readers"{
					template.Profiles["samplemultinodeetcdraft"].Orderer.Policies["Readers"]=value
					delete(template.Profiles["samplemultinodeetcdraft"].Orderer.Policies,"readers")
				}
		
				if name == "writers"{
				template.Profiles["samplemultinodeetcdraft"].Orderer.Policies["Writers"]=value
					delete(template.Profiles["samplemultinodeetcdraft"].Orderer.Policies,"writers")
				}

				if name == "blockvalidation"{
					template.Profiles["samplemultinodeetcdraft"].Orderer.Policies["BlockValidation"]=value
					delete(template.Profiles["samplemultinodeetcdraft"].Orderer.Policies,"blockvalidation")
				}
			}
          
			//cap
			template.Profiles["samplemultinodeetcdraft"].Capabilities = Capabilities["Channel"]
			for name,value:= range template.Profiles["samplemultinodeetcdraft"].Policies{
				if name =="admins"{
					template.Profiles["samplemultinodeetcdraft"].Policies["Admins"] = value
					delete(template.Profiles["samplemultinodeetcdraft"].Policies,"admins")
				}
				
				if name == "readers"{
					template.Profiles["samplemultinodeetcdraft"].Policies["Readers"]=value
					delete(template.Profiles["samplemultinodeetcdraft"].Policies,"readers")
				}
		
				if name == "writers"{
				template.Profiles["samplemultinodeetcdraft"].Policies["Writers"]=value
					delete(template.Profiles["samplemultinodeetcdraft"].Policies,"writers")
				}
			}
			for cName,cvalue := range value.Consortiums{
                if cName == "sampleconsortium"{
					value.Consortiums["SampleConsortium"]=cvalue
					delete(value.Consortiums,"sampleconsortium")
				}
			}
			template.Profiles["SampleMultiNodeEtcdRaft"] = value
			delete(template.Profiles,"samplemultinodeetcdraft")

		}
		
		if name == "twoorgschannel"{
			template.Profiles["twoorgschannel"].Application.Capabilities = Capabilities["Application"]
			for name,value:= range template.Profiles["twoorgschannel"].Application.Policies{
				if name =="admins"{
					template.Profiles["twoorgschannel"].Application.Policies["Admins"] = value
					delete(template.Profiles["twoorgschannel"].Application.Policies,"admins")
				}
				
				if name == "readers"{
				template.Profiles["twoorgschannel"].Application.Policies["Readers"]=value
					delete(template.Profiles["twoorgschannel"].Application.Policies,"readers")
				}
		
				if name == "writers"{
				template.Profiles["twoorgschannel"].Application.Policies["Writers"]=value
					delete(template.Profiles["twoorgschannel"].Application.Policies,"writers")
				}				
			}
            //cap
			template.Profiles["twoorgschannel"].Capabilities = Capabilities["Channel"]
			//policies
			for name,value:= range template.Profiles["twoorgschannel"].Policies{
				if name =="admins"{
					template.Profiles["twoorgschannel"].Policies["Admins"] = value
					delete(template.Profiles["twoorgschannel"].Policies,"admins")
				}
				
				if name == "readers"{
				template.Profiles["twoorgschannel"].Policies["Readers"]=value
					delete(template.Profiles["twoorgschannel"].Policies,"readers")
				}
		
				if name == "writers"{
				template.Profiles["twoorgschannel"].Policies["Writers"]=value
					delete(template.Profiles["twoorgschannel"].Policies,"writers")
				}				
			}
			template.Profiles["TwoOrgsChannel"]=template.Profiles["twoorgschannel"]
			delete(template.Profiles,"twoorgschannel")
		}
	}

	return template,nil

		// //创世取款
		// config := template.Profiles["samplemultinodeetcdraft"]
		// pgen := encoder.New(config)
		// if config.Orderer == nil {
		// 	return errors.Errorf("refusing to generate block which is missing orderer section")
		// }
		// if config.Consortiums == nil {
		// 	return errors.New("Genesis block does not contain a consortiums")
		// }
		// genesisBlock := pgen.GenesisBlockForChannel("byfn-sys-channel")
		// outputPath := filepath.Join(general.BaseOutput, "deploy","channel-artifacts")
		// outputBlocks := filepath.Join(outputPath, "genesis.block")
		// err := ioutil.WriteFile(outputBlocks, utils.MarshalOrPanic(genesisBlock), 0644)
		// if err != nil {
		// 	return fmt.Errorf("Error writing genesis block: %s", err)
		// }
		// //通道tx 
		// chconfig := template.Profiles["twoorgschannel"]
		// var configtx *cb.Envelope
		
		// configtx, err = encoder.MakeChannelCreationTransaction(general.ChannelName, nil, chconfig)
		
		// if err != nil {
		// 	return err
		// }
		// outputChannelCreateTx := filepath.Join(outputPath, general.ChannelName+".tx")
		// err = ioutil.WriteFile(outputChannelCreateTx, utils.MarshalOrPanic(configtx), 0644)
		// if err != nil {
		// 	return fmt.Errorf("Error writing channel create tx: %s", err)
		// }
		// //锚节点tx
		// for _,org := range general.Org{
		// 	anconfig := template.Profiles["twoorgschannel"]
		// 	if anconfig.Application == nil {
		// 		return fmt.Errorf("Cannot update anchor peers without an application section")
		// 	}
		// 	asOrg :=  fmt.Sprintf("%sMSP",org.Name)
		// 	var org *localconfig.Organization
		// 	for _, iorg := range anconfig.Application.Organizations {
		// 		if iorg.Name == asOrg {
		// 			org = iorg
		// 		}
		// 	}
		
		// 	if org == nil {
		// 		return fmt.Errorf("No organization name matching: %s", asOrg)
		// 	}
		// 	anchorPeers := make([]*pb.AnchorPeer, len(org.AnchorPeers))
		// 	for i, anchorPeer := range org.AnchorPeers {
		// 		anchorPeers[i] = &pb.AnchorPeer{
		// 			Host: anchorPeer.Host,
		// 			Port: int32(anchorPeer.Port),
		// 		}
		// 	}
		// 	configUpdate := &cb.ConfigUpdate{
		// 		ChannelId: general.ChannelName,
		// 		WriteSet:  cb.NewConfigGroup(),
		// 		ReadSet:   cb.NewConfigGroup(),
		// 	}

        //     // Add all the existing config to the readset
		// 	configUpdate.ReadSet.Groups[channelconfig.ApplicationGroupKey] = cb.NewConfigGroup()
		// 	configUpdate.ReadSet.Groups[channelconfig.ApplicationGroupKey].Version = 1
		// 	configUpdate.ReadSet.Groups[channelconfig.ApplicationGroupKey].ModPolicy = channelconfig.AdminsPolicyKey
		// 	configUpdate.ReadSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name] = cb.NewConfigGroup()
		// 	configUpdate.ReadSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Values[channelconfig.MSPKey] = &cb.ConfigValue{}
		// 	configUpdate.ReadSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Policies[channelconfig.ReadersPolicyKey] = &cb.ConfigPolicy{}
		// 	configUpdate.ReadSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Policies[channelconfig.WritersPolicyKey] = &cb.ConfigPolicy{}
		// 	configUpdate.ReadSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Policies[channelconfig.AdminsPolicyKey] = &cb.ConfigPolicy{}
		
		// 	// Add all the existing at the same versions to the writeset
		// 	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey] = cb.NewConfigGroup()
		// 	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Version = 1
		// 	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].ModPolicy = channelconfig.AdminsPolicyKey
		// 	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name] = cb.NewConfigGroup()
		// 	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Version = 1
		// 	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].ModPolicy = channelconfig.AdminsPolicyKey
		// 	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Values[channelconfig.MSPKey] = &cb.ConfigValue{}
		// 	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Policies[channelconfig.ReadersPolicyKey] = &cb.ConfigPolicy{}
		// 	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Policies[channelconfig.WritersPolicyKey] = &cb.ConfigPolicy{}
		// 	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Policies[channelconfig.AdminsPolicyKey] = &cb.ConfigPolicy{}
		// 	configUpdate.WriteSet.Groups[channelconfig.ApplicationGroupKey].Groups[org.Name].Values[channelconfig.AnchorPeersKey] = &cb.ConfigValue{
		// 		Value:     utils.MarshalOrPanic(channelconfig.AnchorPeersValue(anchorPeers).Value()),
		// 		ModPolicy: channelconfig.AdminsPolicyKey,
		// 	}

		// 	configUpdateEnvelope := &cb.ConfigUpdateEnvelope{
		// 		ConfigUpdate: utils.MarshalOrPanic(configUpdate),
		// 	}
		// 	update := &cb.Envelope{
		// 		Payload: utils.MarshalOrPanic(&cb.Payload{
		// 			Header: &cb.Header{
		// 				ChannelHeader: utils.MarshalOrPanic(&cb.ChannelHeader{
		// 					ChannelId: general.ChannelName,
		// 					Type:      int32(cb.HeaderType_CONFIG_UPDATE),
		// 				}),
		// 			},
		// 			Data: utils.MarshalOrPanic(configUpdateEnvelope),
		// 		}),
		// 	}
		// 	outputAnchorPeersUpdate := fmt.Sprintf("%s/deploy/channel-artifacts/%sMSPanchors.tx",general.BaseOutput,org.Name)
		// 	err := ioutil.WriteFile(outputAnchorPeersUpdate, utils.MarshalOrPanic(update), 0644)
		// 	if err != nil {
		// 		return fmt.Errorf("Error writing channel anchor peer update: %s", err)
		// 	}
		// 	return nil
		// }
		
	
}

func CreateGenesisBlockAndCifeExecComanndLineAnther(general *objectdefine.Indent) (outputPath string, ret []*objectdefine.CommandList) {


	readconfig := filepath.Join(general.BaseOutput, "deploy","crypto-config.yaml")
	var env []string
	env = []string{fmt.Sprintf("FABRIC_CFG_PATH=%s", filepath.ToSlash(filepath.Join(general.BaseOutput, "deploy")))}

    os.Setenv("FABRIC_CFG_PATH",filepath.ToSlash(filepath.Join(general.BaseOutput, "deploy")))
	var cmd *exec.Cmd
	execCommand := fmt.Sprintf("export FABRIC_CFG_PATH=%s",filepath.ToSlash(filepath.Join(general.BaseOutput, "deploy")))
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", execCommand)
	} else {
		cmd = exec.Command("/bin/bash", "-c", execCommand)
	}
	if err := cmd.Start(); err != nil {
		return "",nil
	}
	if err := cmd.Wait(); err != nil {
		return "",nil
	}
	ret = make([]*objectdefine.CommandList, 0)

	toolsPath := dcache.GetBinPathByVersion(general.Version)

	pair := dcache.GetBinToolsPair()

	bins, ok := pair[runtime.GOOS]
	if !ok {
		return "", nil
	}
	one := &objectdefine.CommandList{}
	one.OS = runtime.GOOS
	one.Call = MakeGenerateExecCommandAnther(env, general, outputPath, readconfig,
		filepath.Join(toolsPath, bins.CryptoGen),
		filepath.Join(toolsPath, bins.ConfigGen))

	ret = append(ret, one)

	return
}

func MakeGenerateExecCommandAnther(env []string, general *objectdefine.Indent, output, readconfig, cryptogen, configtxgen string) []objectdefine.ProcessPair {
	//ChannelName := general.ChannelName

	//outputBlocks := filepath.Join(output, "genesis.block")
    baseOut := filepath.Join(general.BaseOutput, "deploy")
	ret := make([]objectdefine.ProcessPair, 0)

	//创建证书命令
	ret = append(ret, objectdefine.ProcessPair{
		Exec:        cryptogen,
		Args:        []string{`generate`, `--config`, readconfig},
		Dir:         baseOut,
		Environment: env,
	})

	return ret
}