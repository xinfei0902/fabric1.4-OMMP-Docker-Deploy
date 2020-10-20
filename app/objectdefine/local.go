package objectdefine

// LocalType for config
type LocalType struct {
	Consensus *ConsensusList `json:"consensus,omitempty"`

	Versions []VersionType `json:"version,omitempty"`

	Bin map[string]BinToolsType `json:"bin,omitempty"`

	StoreRoot string `json:"storepath,omitempty"`
}

//BinToolsType  创建区块和证书工具
type BinToolsType struct {
	CryptoGen string `json:"crypto,omitempty"`
	ConfigGen string `json:"configtx,omitempty"`
}

// Const values for file name
const (
	ConstHistoryTaskStatusFileName  = "status.json"
	ConstHistoryTaskIndentFileName  = "indent.json"
	ConstHistoryTaskGeneralFileName = "NewIndent.json"
)

// LocalStruct for local configs
const (
	BuildTypeDeploy    = "deploy"
	BuildTypeChannel   = "channel"
	BuildTypeChainCode = "chaincode"
	BuildTypeDaemon    = "daemon"
	BuildTypeSDK       = "sdk"
	BuildTypeCA        = "ca"
	BuildTypePeer      = "peer"
	BuildTypeOrderer   = "orderer"
	BuildToolsLog      = "toolslog"
	BuildTypeKafka     = "kafka"
	BuildTypeZookeeper = "zookeeper"
	BuildTypeOrg       = "org"
	TemplateSystemd    = "systemd"
	TemplateSDK        = "sdk"
)
