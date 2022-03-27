package objectdefine

// EtcdRaftStruct for Indent
type EtcdRaftStruct struct {
	Consenters []*Consenter `protobuf:"bytes,1,rep,name=consenters" json:"consenters,omitempty"`
}

// Consenter for Indent
type Consenter struct {
	Host          string `protobuf:"bytes,1,opt,name=host" json:"host,omitempty"`
	Port          uint32 `protobuf:"varint,2,opt,name=port" json:"port,omitempty"`
	ClientTLSCert []byte `protobuf:"bytes,3,opt,name=client_tls_cert,json=clientTlsCert,proto3" json:"client_tls_cert,omitempty"`
	ServerTLSCert []byte `protobuf:"bytes,4,opt,name=server_tls_cert,json=serverTlsCert,proto3" json:"server_tls_cert,omitempty"`
}

// FabricStruct for Indent
type FabricStruct struct {
	BatchTimeout string    `json:"BatchTimeout,omitempty"`
	BatchSize    BatchSize `json:"BatchSize,omitempty"`
}

// BatchSize for FabricStruct
type BatchSize struct {
	MaxMessageCount   uint32 `json:"MaxMessageCount,omitempty"`
	AbsoluteMaxBytes  string `json:"AbsoluteMaxBytes,omitempty"`
	PreferredMaxBytes string `json:"PreferredMaxBytes,omitempty"`
}

// KafkaStruct for Indent
type KafkaStruct struct {
	Kafka     []KafkaType     `json:"kafka,omitempty"`
	Zookeeper []ZooKeeperType `json:"zookeeper,omitempty"`
}

// KafkaType for KafkaStruct
type KafkaType struct {
	ConfigID string `json:"configid,omitempty"`

	IP   string `json:"ip"`
	Name string `json:"name,omitempty"`
	Port int    `json:"port,omitempty"`
	//LodDirs string `json:"logdirs,omitempty"`
	BrokeID int `json:"brokerid,omitempty"`

	//SetupPath     string            `json:"setuppath,omitempty"`
	//LogPath       string            `json:"logpath,omitempty"`
	//DataPath      string            `json:"datapath,omitempty"`
	Folder string `json:"folder,omitempty"`
	Domain string `json:"domain,omitempty"`
	//OldScriptName string `json:"-"`
	//NewScriptName string `json:"-"`
	//Build         *KafkaBuildPair   `json:"build"`
	//Service *KafkaServiceType `json:"-"`
}

// KafkaBuildPair for KafkaType
type KafkaBuildPair struct {
	// Output for build step
	OutputRoot string `json:"pathbuild,omitempty"`

	Env map[string]string `json:"env,omitempty"`
}
type KafkaServiceType struct {
	SetupPath          string `json:"-"`
	ServiceName        string `json:"-"`
	SetlogPath         string `json:"-"`
	StartCommand       string `json:"-"`
	KafkaConfigMark    string `json:"-"`
	StartCommandParams string `json:"-"`
	StopCommand        string `json:"-"`
	StopCommandParams  string `json:"-"`
	JavaPath           string `json:"-"`
	Env                string `json:"-"`
}

// ZooKeeperType for KafkaStruct
type ZooKeeperType struct {
	ConfigID string `json:"configid,omitempty"`

	IP     string `json:"ip"`
	Name   string `json:"name,omitempty"`
	Port   int    `json:"port,omitempty"`
	Follow int    `json:"follow,omitempty"`
	Vote   int    `json:"vote,omitempty"`
	ID     int    `json:"id,omitempty"`

	//SetupPath string `json:"setuppath,omitempty"`
	//LogPath   string `json:"logpath,omitempty"`
	//DataPath  string `json:"datapath,omitempty"`
	Folder string `json:"folder,omitempty"`
	Domain string `json:"domain,omitempty"`

	//OldScriptName string                `json:"-"`
	//NewScriptName string                `json:"-"`
	//Build         *ZookeeperBuildType   `json:"build"`
	//Service       *ZookeeperServiceType `json:"-"`
}

// ZookeeperBuildType for ZooKeeperType
type ZookeeperBuildType struct {
	// Output for build step
	OutputRoot string `json:"pathbuild,omitempty"`

	Env map[string]string `json:"env,omitempty"`
}

type ZookeeperServiceType struct {
	SetupPath     string `json:"-"`
	ServiceName   string `json:"-"`
	SetlogPath    string `json:"-"`
	Command       string `json:"-"`
	CommandParams string `json:"-"`
	JavaPath      string `json:"-"`
	Env           string `json:"-"`
}

// OrgType for Indent
type OrgType struct {
	Name string `json:"name,omitempty"`
	//BaseDomain string `json:"basedomain,omitempty"`
	OrgDomain string `json:"orgdomain,omitempty"`

	ChainCode map[string]ChainCodeType `json:"chaincode,omitempty"`
	Peer      []PeerType               `json:"peer,omitempty"`
	CA        *CAType                  `json:"ca,omitempty"`
	CliName   string                   `json:"cliname,omitempty"`
}

// ChainCodeType for Indent
type ChainCodeType struct {
	Version        string   `json:"version,omitempty"`
	Threshold      int      `json:"threshold,omitempty"`
	Policy         string   `json:"policy,omitempty"`
	EndorsementOrg []string `json:"endorsorg,omitempty"`
	Describe       string   `json:"desc,omitempty"`
	IsInstall      int      `json:"isinstall,omitempty"`
	Status         int      `json:"status,omitempty"`
}

// PeerType for OrgType
type PeerType struct {
	ConfigID string `json:"configid,omitempty"`

	IP         string `json:"ip"`
	Name       string `json:"name,omitempty"`
	NickName   string `json:"nickname,omitempty"`
	AccessKey  string `json:"accesskey,omitempty"`
	Org        string `json:"org,omitempty"`
	BaseDomain string `json:"basedomain,omitempty"`
	OrgDomain  string `json:"orgdomain,omitempty"`

	User string `json:"user,omitempty"`

	PeerID int `json:"peerid,omitempty"`
	//NetworkID string `json:"networkid,omitempty"`

	Port          int  `json:"port,omitempty"`
	CouchdbPort   int  `json:"couchdbport,omitempty"`
	ChaincodePort int  `json:"chaincodeport,omitempty"`
	CC            bool `json:"cc,omitempty"`
	Anchor        bool `json:"anchor,omitempty"`

	ChainCode map[string]ChainCodeType `json:"chaincode,omitempty"`

	SetupPath string `json:"setuppath,omitempty"`
	LogPath   string `json:"logpath,omitempty"`
	DataPath  string `json:"datapath,omitempty"`
	Folder    string `json:"folder,omitempty"`
	Domain    string `json:"domain,omitempty"`
	CliName   string `json:"cliname,omitempty"`
	CliInstallStatus   string `json:"clistatus,omitempty"`
	Status    int    `json:"-"`
	RunStatus int    `json:"-"`
}

// Keys in map that BaaSType.Config
const (
	KeyInBaaSTypeConfigLink   = "link"
	KeyInBaaSTypeConfigReport = "report"
)

// CAType for OrgType
type CAType struct {
	ConfigID string `json:"configid,omitempty"`

	IP string `json:"ip"`

	Name       string `json:"name,omitempty"`
	Org        string `json:"org,omitempty"`
	BaseDomain string `json:"basedomain,omitempty"`
	OrgDomain  string `json:"orgdomain,omitempty"`

	Port int `json:"port,omitempty"`

	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`

	SetupPath string `json:"setuppath,omitempty"`
	Folder    string `json:"folder,omitempty"`
	Domain    string `json:"domain,omitempty"`
}

// OrderType for Indent
type OrderType struct {
	ConfigID string `json:"configid,omitempty"`

	IP string `json:"ip"`

	Name string `json:"name,omitempty"`

	Org        string `json:"org,omitempty"`
	BaseDomain string `json:"basedomain,omitempty"`
	OrgDomain  string `json:"orgdomain,omitempty"`

	Domain string `json:"domain,omitempty"`

	Port int `json:"port,omitempty"`

	ClientTLSCert []byte `json:"-"`
	ServerTLSCert []byte `json:"-"`

	SetupPath     string `json:"setuppath,omitempty"`
	LogPath       string `json:"logpath,omitempty"`
	DataPath      string `json:"datapath,omitempty"`
	Folder        string `json:"folder,omitempty"`
	OldScriptName string `json:"-"`
	NewScriptName string `json:"-"`
}

type OrdererBuildType struct {
	Name      string
	OrgDomain string
	Hosts     []string
	URLs      []string

	SubMSP string
}
type OrdererServiceType struct {
	Folder         string `json:"-"`
	StartCommand   string `json:"-"`
	StartSub       string `json:"-"`
	StartSubParams string `json:"-"`
	StopCommand    string `json:"-"`
	SetupPath      string `json:"-"`
	ServiceName    string `json:"-"`
	SetlogPath     string `json:"-"`
	Env            string `json:"-"`
	ToolsBin       string `json:"-"`
}

// OrdererBuildPair for OrderType
type OrdererBuildPair struct {
	// Output for build step
	OutputRoot string `json:"pathbuild,omitempty"`

	// Input for build step
	SubTLS          string `json:"pathbuildtls,omitempty"`
	SubMSP          string `json:"pathbuildmsp,omitempty"`
	SubGenesisBlock string `json:"pathgenesisblock,omitempty"`
}

// DaemonType for Indent
type DaemonType struct {
	Report        string `json:"report,omitempty"`
	Snap          string `json:"snap,omitempty"`
	KeepAlive     string `json:"keepalive,omitempty"`
	OldScriptName string `json:"-"`
	NewScriptName string `json:"-"`
}

// SetupPathType for Indent
type SetupPathType struct {
	IP   string `json:"ip,omitempty"`
	Path string `json:"Path,omitempty"`
}

//DeployType 一键部署信息
type DeployType struct {
	JoinOrg map[string]OrgType       `json:"org,omitempty"`
	JoinCC  map[string]ChainCodeType `json:"cc,omitempty"`
}
