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

	"github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
	"github.com/pkg/errors"
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
	templateBlock, err := GetGenesisCompleteBlockTemplate(general, output)
	if err != nil {
		err = errors.WithMessage(err, "Load template for genesisblock error")
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateCompleteDeployBlock",
			Error: []string{err.Error()},
		})
		return err
	}
	genesis := CreateGenesisCompleteBlockConfig(general, templateBlock)
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

	output.AppendLog(history)
	return nil
}

//GetGenesisCompleteBlockTemplate  获取模板转结构
func GetGenesisCompleteBlockTemplate(general *objectdefine.Indent, output *objectdefine.TaskNode) (ret *localconfig.TopLevel, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Errorf("Get template of genesis block config error: %v", e)
			ret = nil
		}
	}()
	VersionRoot := dcache.GetTemplatePathByVersion(general.Version)
	path := filepath.Join(VersionRoot, "template", "channel")

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
	ret.Channel = template.Channel
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
	}
	one.Organizations = template.Orderer.Organizations
	one.Policies = template.Orderer.Policies
	return &one
}

func genesisCompleteBlockProfiles(general *objectdefine.Indent, template *localconfig.TopLevel, self *localconfig.TopLevel) map[string]*localconfig.Profile {
	ret := make(map[string]*localconfig.Profile)
	//创世区块的联盟配置
	genesis := &localconfig.Profile{}
	genesis.Capabilities = self.Channel.Capabilities
	genesis.Orderer = self.Orderer
	ordererOrganizations := make([]*localconfig.Organization, 1)
	ordererOrganizations[0] = self.Organizations[0]
	genesis.Orderer.Organizations = ordererOrganizations
	genesis.Orderer.Capabilities = self.Capabilities["Orderer"]
	genesis.Application = self.Application
	genesis.Application.Organizations = ordererOrganizations
	genesis.Consortiums = make(map[string]*localconfig.Consortium)
	genesisOrganizations := make([]*localconfig.Organization, len(general.Org))

	for i := 0; i < len(general.Org); i++ {
		genesisOrganizations[i] = self.Organizations[i+1]
	}
	genesis.Consortiums = map[string]*localconfig.Consortium{
		"GuizhouPrisonConsortium": &localconfig.Consortium{
			Organizations: genesisOrganizations,
		},
	}
	ret["GuizhouPrisonNetworkGenesis"] = genesis
	//每一个通道的配置
	for channelName, deploy := range general.Deploy {
		channelConfig := genesisCompleteBlockChannelConfig(deploy, self)
		ret[channelName] = channelConfig
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
	channel.Consortium = "GuizhouPrisonConsortium"
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

	filename := filepath.Join(general.BaseOutput, "crypto-config.yaml")
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
	err = os.MkdirAll(general.BaseOutput, 0777)
	if nil != err {
		err = errors.WithMessage(err, "Make Output folder error")
		return
	}

	filename := filepath.Join(general.BaseOutput, "configtx.yaml")
	err = tools.WriteYamlFile(filename, config)
	if nil != err {
		err = errors.WithMessage(err, "Write genesis config error")
		return
	}

	return nil
}

//CreateGenesisBlockAndCifeExecComanndLine 构造目录 组装生成创世区块和证书命令
func CreateGenesisBlockAndCifeExecComanndLine(general *objectdefine.Indent) (outputPath string, ret []*objectdefine.CommandList) {

	outputPath = filepath.Join(general.BaseOutput, "channel-artifacts")

	readconfig := filepath.Join(general.BaseOutput, "crypto-config.yaml")
	var env []string
	env = []string{fmt.Sprintf("FABRIC_CFG_PATH=%s", filepath.ToSlash(general.BaseOutput))}

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
	//ChannelName := general.ChannelName

	outputBlocks := filepath.Join(output, "genesis.block")

	ret := make([]objectdefine.ProcessPair, 0)

	//创建证书命令
	ret = append(ret, objectdefine.ProcessPair{
		Exec:        cryptogen,
		Args:        []string{`generate`, `--config`, readconfig},
		Dir:         general.BaseOutput,
		Environment: env,
	})

	//block file 用来创建创世区块的
	ret = append(ret, objectdefine.ProcessPair{
		Exec:        configtxgen,
		Args:        []string{`-profile`, `GuizhouPrisonNetworkGenesis`, `-channelID`, `byfn-sys-channel`, `-outputBlock`, outputBlocks},
		Dir:         general.BaseOutput,
		Environment: env,
	})

	for channelName, deploy := range general.Deploy {
		// channel file
		outputChannelTx := filepath.Join(output, channelName+".tx")
		ret = append(ret, objectdefine.ProcessPair{
			Exec:        configtxgen,
			Args:        []string{`-profile`, channelName, `-channelID`, channelName, `-outputCreateChannelTx`, outputChannelTx},
			Dir:         general.BaseOutput,
			Environment: env,
		})
		for _, org := range deploy.JoinOrg {
			if len(org.Peer) == 0 {
				continue
			}
			for _, peer := range org.Peer {
				if "Admin" != peer.User {
					continue
				}
				orgMsp := fmt.Sprintf("%sMSP", org.Name)
				ret = append(ret, objectdefine.ProcessPair{
					Exec:        configtxgen,
					Args:        []string{`-profile`, channelName, `-channelID`, channelName, `-asOrg`, orgMsp, `-outputAnchorPeersUpdate`, filepath.Join(output, fmt.Sprintf("%sMSPanchors@%s.tx", org.Name, channelName))},
					Dir:         general.BaseOutput,
					Environment: env,
				})
			}
		}
	}
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
