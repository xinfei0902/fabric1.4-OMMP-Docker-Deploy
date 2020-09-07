package objectdefine

// Indent for API
type Indent struct {
	ID                    string `json:"id,omitempty"` //任务ID
	SourceID              string `json:"sourceid,omitempty"`
	Desc                  string `json:"desc,omitempty"`
	IsNewOrgCreateChannel bool   `json:"isneworgcreatechannel,omitempty"` //创建通道 是基于已经存在组织还是重新生成组织

	BaseOutput       string `json:"-"` //以任务ID名称为目录的路径
	SourceBaseOutput string `json:"-"`

	ChannelName string                     `json:"channelname,omitempty"` //通道名称
	Deploy      map[string]DeployType      `json:"deploy,omitempty"`      //通道名称 为一键部署准备字段
	Chaincode   map[string][]ChainCodeType `json:"chaincode,omitempty"`   //为了获取一个合约多次版本更替，用来检测合约设置为数组

	Model string `json:"model,omitempty"`

	Consensus string `json:"consensus,omitempty"` //共识模式  solo kafka raft ....

	Kafka *KafkaStruct `json:"kafka,omitempty"`

	Org map[string]OrgType `json:"org,omitempty"` //组织信息

	Orderer []OrderType `json:"orderer,omitempty"` //orderer信息

	//OrdererBuild map[string]OrdererBuildType `json:"ordererbuild"`

	Version string `json:"version,omitempty"`

	FireWall map[string]map[int]bool `json:"firewall,omitempty"` //存放所有IP下使用的端口
	//SetupPath map[string]SetupPathType `json:"setup,omitempty"`
	//Secret *IndentSecret `json:"secret,omitempty"`
}
