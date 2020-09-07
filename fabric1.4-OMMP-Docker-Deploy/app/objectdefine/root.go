package objectdefine

import "github.com/pkg/errors"

////////////////////////////// net ////////////////////////////////////////
const (
	ConsensusKafka    = "kafka"
	ConsensusSolo     = "solo"
	ConsensusEtcdRaft = "etcdraft"

	DefaultBaseDomain  = "example.com"
	DefaultChannelName = "mychannel"
	DefaultBaseRoot    = "/etc/sinochain"

	DefaultKafkaLogDirs = "tmp/kafka-logs"

	EnvKeyJavaHome      = "JAVA_HOME"
	EnvKeyFabricCfgPath = "FABRIC_CFG_PATH"
)
const (
	//DEFAULTUSERNUMBER 呃
	DEFAULTUSERNUMBER = 100000
)

// Errors
var (
	ErrorMissOrg         = errors.New("Miss Org Name in configs")
	ErrorMissKafka       = errors.New("Miss kafka configs in kafka consensus")
	ErrorMissOrderer     = errors.New("Miss Orderer configs in EtcdRaft consensus")
	ErrorMissPeerIP      = makeErrorMissIP("Peer")
	ErrorMissOrdererIP   = makeErrorMissIP("Orderer")
	ErrorMissBaaSIP      = makeErrorMissIP("BaaS")
	ErrorMissCAIP        = makeErrorMissIP("CA")
	ErrorMissKafkaIP     = makeErrorMissIP("Kafka")
	ErrorMissZookeeperIP = makeErrorMissIP("Zookeeper")
)

func makeErrorMissIP(name string) error {
	return errors.Errorf("Miss %s IP configs", name)
}
