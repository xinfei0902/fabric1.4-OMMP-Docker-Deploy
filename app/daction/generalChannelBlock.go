package daction

import (
	"deploy-server/app/dcache"
	"deploy-server/app/dmysql"
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

//MakeStepGeneralCreateBlock 创建configtx.yaml
func MakeStepGeneralCreateBlock(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create channel block configtx.yaml start"},
		})

		err := GeneralCreateChannelBlock(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateChannelTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChannelTx",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create channel block configtx.yaml end"},
		})
		return nil
	}
}

//GeneralCreateChannelBlock 创建通道 .tx文件 以及锚节点 .tx文件
func GeneralCreateChannelBlock(general *objectdefine.Indent, output *objectdefine.TaskNode) error {

	templateBlock, err := GetGenesisBlockTemplate(general, output)
	if err != nil {
		err = errors.WithMessage(err, "Load template for genesisblock error")
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChannelTx",
			Error: []string{err.Error()},
		})
		return err
	}
	//注释部分为原先用来生成证书已经补全configtx.yaml文件的 初版先不使用
	//genesis := CreateGenesisBlockConfig(general, templateBlock)
	// crypto := &objectdefine.CryptoForOrgConfig{}
	// if general.IsNewOrgCreateChannel == true {
	// 	crypto = CreateCryptoForOrgConfig(general)
	// 	err := CreateCryptoForOrg(general, crypto)
	// 	if err != nil {
	// 		err = errors.WithMessage(err, "Create crypto files config error")
	// 		output.AppendLog(&objectdefine.StepHistory{
	// 			Name:  "generalCreateBlock",
	// 			Error: []string{err.Error()},
	// 		})
	// 		return err
	// 	}
	// }

	err = CreateGenesisBlock(general, templateBlock)
	if err != nil {
		err = errors.WithMessage(err, "Create genesisblock config error")
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChannelTx",
			Error: []string{err.Error()},
		})
		return err
	}
	//拷贝工具
	srcBinPath := dcache.GetBinPathByVersion(general.Version)
	_, err = os.Stat(general.BaseOutput)
	if err != nil {
		err = os.MkdirAll(general.BaseOutput, 0777)
		if nil != err {
			return errors.WithMessage(err, "Make Output folder error")
		}
	}
	tools.CopyFolder(general.BaseOutput, srcBinPath)
	//拷贝证书
	src := filepath.Join(general.SourceBaseOutput, "crypto-config")
	dst := filepath.Join(general.BaseOutput, "crypto-config")
	_, err = os.Stat(dst)
	if err != nil {
		err = os.MkdirAll(dst, 0777)
		if nil != err {

			return errors.WithMessage(err, "Make Output folder error")
		}
	}
	tools.CopyFolder(dst, src)
	//用工具生成
	outputPath, calls := V1CreateScriptForChannelArtifacts(general)
	err = os.MkdirAll(outputPath, 0777)
	if err != nil {
		err = errors.WithMessage(err, "Prepare path for crypto files error")
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChannelTx",
			Error: []string{err.Error()},
		})
		return err
	}
	logs, err := V1CallScriptForChannelArtifacts(general, calls)

	history := &objectdefine.StepHistory{
		Name: "generalCreateChannelTx",
		Log:  logs,
	}
	if err != nil {
		history.Error = []string{err.Error()}
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreateChannelTx",
			Error: []string{err.Error()},
		})
	}

	output.AppendLog(history)
	return nil
}

//GetGenesisBlockTemplate  获取模板转结构
func GetGenesisBlockTemplate(general *objectdefine.Indent, output *objectdefine.TaskNode) (ret *localconfig.TopLevel, err error) {
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

//CreateGenesisBlockConfig  补全构造新的 configtx.yaml文件
func CreateGenesisBlockConfig(general *objectdefine.Indent, template *localconfig.TopLevel) (ret *localconfig.TopLevel) {
	ret = &localconfig.TopLevel{}
	ret.Organizations = genesisBlockOrganizations(general, template)
	ret.Capabilities = template.Capabilities
	ret.Application = template.Application
	ret.Channel = template.Channel
	ret.Orderer = genesisBlockOrderer(general, template, ret.Organizations)
	ret.Profiles = genesisBlockProfiles(general, template, ret)
	return ret
}

func genesisBlockOrganizations(general *objectdefine.Indent, template *localconfig.TopLevel) []*localconfig.Organization {
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

func genesisBlockOrderer(general *objectdefine.Indent, template *localconfig.TopLevel, orgs []*localconfig.Organization) *localconfig.Orderer {

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

func genesisBlockProfiles(general *objectdefine.Indent, template *localconfig.TopLevel, self *localconfig.TopLevel) map[string]*localconfig.Profile {
	ret := make(map[string]*localconfig.Profile)
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
	genesosOrganizations := make([]*localconfig.Organization, len(general.Org))

	for i := 0; i < len(general.Org); i++ {
		genesosOrganizations[i] = self.Organizations[i+1]
	}
	genesis.Consortiums = map[string]*localconfig.Consortium{
		"GuizhouPrisonConsortium": &localconfig.Consortium{
			Organizations: genesosOrganizations,
		},
	}
	ret["GuizhouPrisonNetworkGenesis"] = genesis
	channel := &localconfig.Profile{}
	channel.Consortium = "GuizhouPrisonConsortium"
	// channel.Capabilities = self.Channel.Capabilities
	// //one.Policies = self.Channel.Policies
	// channel.Orderer = self.Orderer
	// organizationsOrderer := make([]*localconfig.Organization, 1)
	// organizationsOrderer[0] = self.Organizations[0]
	// channel.Orderer.Organizations = organizationsOrderer
	channel.Application = self.Application
	channel.Application.Organizations = genesosOrganizations
	channel.Application.Capabilities = self.Capabilities["Application"]
	// channel.Consortiums = make(map[string]*localconfig.Consortium)
	// organizations := make([]*localconfig.Organization, len(general.Org))
	// for i := 0; i < len(general.Org); i++ {
	// 	organizations[i] = self.Organizations[i+1]
	// }
	// channel.Consortiums = map[string]*localconfig.Consortium{
	// 	"SampleConsortium": &localconfig.Consortium{
	// 		Organizations: organizations,
	// 	},
	// }
	ret["TwoOrgsChannel"] = channel

	return ret
}

//CreateGenesisBlock 创建configtx.yaml文件
func CreateGenesisBlock(general *objectdefine.Indent, config *localconfig.TopLevel) (err error) {
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

	// V2
	// task.TargetPath /--
	//                  |- SubPathOfCryptoForOrgConfig
	//                  |- SubPathOfGenesisBlockConfig
	//                  |- channel-artifacts /--
	//                                        |- SubPathOfGenesisBlock
	//                                        |- ChannelName.tx

	return nil
}

//CreateCryptoForOrgConfig 创建生成证书的crypto-config.yaml文件
func CreateCryptoForOrgConfig(general *objectdefine.Indent) *objectdefine.CryptoForOrgConfig {
	ret := &objectdefine.CryptoForOrgConfig{}

	//ret.OrdererOrgs = make([]objectdefine.OrgSpec, 0, 1)
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
	// two := objectdefine.OrgSpec{}
	// var ordererName string
	// var ordererDomain string
	// two.Specs = make([]objectdefine.NodeSpec, len(general.Orderer))
	// for i, orderer := range general.Orderer {
	// 	if len(ordererName) == 0 {
	// 		ordererDomain = orderer.OrgDomain
	// 	}
	// 	//hostName := fmt.Sprintf("%s:%d", orderer.Domain, orderer.Port)
	// 	two.Specs[i].Hostname = strings.ToLower(orderer.Name)
	// }

	// two.Name = "Orderer"
	// two.Domain = ordererDomain
	// two.EnableNodeOUs = true
	// two.Users = objectdefine.UsersSpec{
	// 	Count: 1,
	// }
	// ret.OrdererOrgs = append(ret.OrdererOrgs, two)
	return ret
}

//CreateCryptoForOrg 生成证书
func CreateCryptoForOrg(general *objectdefine.Indent, config *objectdefine.CryptoForOrgConfig) (err error) {
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

	src := filepath.Join(general.SourceBaseOutput, "crypto-config", "ordererOrganizations")
	dst := filepath.Join(general.BaseOutput, "crypto-config", "ordererOrganizations")
	_, err = os.Stat(dst)
	if err != nil {
		err = os.MkdirAll(dst, 0777)
		if nil != err {
			err = errors.WithMessage(err, "Make Output folder error")
			return
		}
	}
	tools.CopyFolder(dst, src)

	// V2
	// task.TargetPath /--
	//                  |- SubPathOfCryptoForOrgConfig
	//                  |- SubPathOfGenesisBlockConfig
	//                  |- channel-artifacts /--
	//                                        |- SubPathOfGenesisBlock
	//                                        |- ChannelName.tx

	return nil
}

//V1CreateScriptForChannelArtifacts 构造目录 组装生成.tx文件命令
func V1CreateScriptForChannelArtifacts(general *objectdefine.Indent) (outputPath string, ret []*objectdefine.CommandList) {
	// task.TargetPath /--
	//                  |- SubPathOfCryptoForOrgConfig
	//                  |- SubPathOfGenesisBlockConfig
	//                  |- generateBlock.ps1
	//                  |- generateBlock.bat
	//                  |- generateBlock.sh
	//                  |- channel-artifacts /--
	//                                        |- SubPathOfGenesisBlock
	//                                        |- ChannelName.tx

	outputPath = filepath.Join(general.BaseOutput, "channel-artifacts")

	readconfig := filepath.Join(general.BaseOutput, "crypto-config.yaml")
	var env []string
	env = []string{fmt.Sprintf("FABRIC_CFG_PATH=%s", filepath.ToSlash(general.BaseOutput))}

	ret = make([]*objectdefine.CommandList, 0, 2)

	toolsPath := dcache.GetBinPathByVersion(general.Version)

	pair := dcache.GetBinToolsPair()

	bins, ok := pair[runtime.GOOS]
	if !ok {
		return "", nil
	}
	one := &objectdefine.CommandList{}
	one.OS = runtime.GOOS
	one.Call = v1MakeGenerateExec(env, general, outputPath, readconfig,
		filepath.Join(toolsPath, bins.CryptoGen),
		filepath.Join(toolsPath, bins.ConfigGen))

	ret = append(ret, one)

	return
}

//v1MakeGenerateExec 组装生成.tx文件命令
func v1MakeGenerateExec(env []string, general *objectdefine.Indent, output, readconfig, cryptogen, configtxgen string) []objectdefine.ProcessPair {
	ChannelName := general.ChannelName

	//outputBlocks := filepath.Join(output, "genesis.block")
	outputChannelTx := filepath.Join(output, ChannelName+".tx")

	ret := make([]objectdefine.ProcessPair, 0, 16)

	// crypto files  原先用来生成证书 初版先不用
	// if general.IsNewOrgCreateChannel == true {
	// 	ret = append(ret, objectdefine.ProcessPair{
	// 		Exec:        cryptogen,
	// 		Args:        []string{`generate`, `--config`, readconfig},
	// 		Dir:         general.BaseOutput,
	// 		Environment: env,
	// 	})
	// }

	//block file 用来创建创世区块的
	// ret = append(ret, objectdefine.ProcessPair{
	// 	Exec:        configtxgen,
	// 	Args:        []string{`-profile`, `GuizhouPrisonNetworkGenesis`, `-channelID`, `byfn-sys-channel`, `-outputBlock`, outputBlocks},
	// 	Dir:         general.BaseOutput,
	// 	Environment: env,
	// })

	// channel file
	ret = append(ret, objectdefine.ProcessPair{
		Exec:        configtxgen,
		Args:        []string{`-profile`, `TwoOrgsChannel`, `-channelID`, ChannelName, `-outputCreateChannelTx`, outputChannelTx},
		Dir:         general.BaseOutput,
		Environment: env,
	})

	for _, one := range general.Org {
		if len(one.Peer) == 0 {
			continue
		}

		for _, two := range one.Peer {
			if "Admin" != two.User {
				continue
			}
			orgMsp := fmt.Sprintf("%sMSP", one.Name)
			ret = append(ret, objectdefine.ProcessPair{
				Exec:        configtxgen,
				Args:        []string{`-profile`, `TwoOrgsChannel`, `-channelID`, ChannelName, `-asOrg`, orgMsp, `-outputAnchorPeersUpdate`, filepath.Join(output, one.Name+"MSPanchors.tx")},
				Dir:         general.BaseOutput,
				Environment: env,
			})
		}
	}

	return ret
}

//V1CallScriptForChannelArtifacts 运行命令 生成.tx文件
func V1CallScriptForChannelArtifacts(general *objectdefine.Indent, list []*objectdefine.CommandList) (logs []string, err error) {

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
