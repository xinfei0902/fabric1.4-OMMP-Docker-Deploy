package daction

import (
	"deploy-server/app/objectdefine"
	"deploy-server/app/tools"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

//MakeStepGeneralCreateCompleteDeployBaseYaml 创建一键部署base.yaml文件
func MakeStepGeneralCreateCompleteDeployBaseYaml(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create complete deploy base.yaml start"},
		})

		err := GeneralCreateCompleteDeployBase(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create complete deploy base.yaml"},
		})
		return nil
	}
}

//GeneralCreateCompleteDeployBase  创建base.yaml
func GeneralCreateCompleteDeployBase(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	ret := &objectdefine.BaseYaml{}
	// buff, err := dcache.GetDeployBaseYamlTemplate(general.Version)
	// if err != nil {
	// 	return err
	// }
	// err = yaml.Unmarshal(buff, ret)
	// if err != nil {
	// 	return errors.WithMessage(err, "Parse template error")
	// }
	service := make(map[string]*objectdefine.DockerConfig)
	ret.Version = fmt.Sprintf("'%s'", "2")
	zkExtraHosts := make([]string, 0)
	var zookeeperConnect string
	for _, zookeeper := range general.Kafka.Zookeeper {
		zkExtraHosts = append(zkExtraHosts, fmt.Sprintf("%s:%s", zookeeper.Domain, zookeeper.IP))
		if len(zookeeperConnect) == 0 {
			zookeeperConnect = fmt.Sprintf("%s:%d", zookeeper.Domain, zookeeper.Port)
		} else {
			zookeeperConnect = zookeeperConnect + "," + fmt.Sprintf("%s:%d", zookeeper.Domain, zookeeper.Port)
		}

	}
	for _, kafka := range general.Kafka.Kafka {
		zkExtraHosts = append(zkExtraHosts, fmt.Sprintf("%s:%s", kafka.Domain, kafka.IP))
	}
	for i, zookeeper := range general.Kafka.Zookeeper {
		zookeeperInfo := &objectdefine.DockerConfig{}
		zookeeperInfo.ContainerName = zookeeper.Domain
		zookeeperInfo.Images = "hyperledger/fabric-zookeeper:latest"
		enviroment := make([]string, 0)
		enviroment = append(enviroment, fmt.Sprintf("ZOO_MY_ID=%d", i+1))
		enviroment = append(enviroment, "ZOO_SERVERS=server.1=zookeeper0:2888:3888 server.2=zookeeper1:2888:3888 server.3=zookeeper2:2888:3888 quorumListenOnAllIPs=true")
		enviroment = append(enviroment, "ZOOKEEPER_CLIENT_PORT=2181")
		enviroment = append(enviroment, "ZOOKEEPER_TICK_TIME=5000")
		zookeeperInfo.Environment = enviroment
		volumes := make([]string, 0)
		volumes = append(volumes, fmt.Sprintf("/blockchainData/%s/data:/data", zookeeper.Domain))
		volumes = append(volumes, fmt.Sprintf("/blockchainData/%s/datalog:/datalog", zookeeper.Domain))
		zookeeperInfo.Volumes = volumes
		ports := make([]string, 0)
		ports = append(ports, fmt.Sprintf("%d:2181", zookeeper.Port))
		ports = append(ports, fmt.Sprintf("%d:2888", zookeeper.Follow))
		ports = append(ports, fmt.Sprintf("%d:3888", zookeeper.Vote))
		zookeeperInfo.Ports = ports
		zookeeperInfo.ExtraHosts = zkExtraHosts
		service[zookeeper.Domain] = zookeeperInfo
	}
	for i, kafka := range general.Kafka.Kafka {
		kafkaInfo := &objectdefine.DockerConfig{}
		kafkaInfo.ContainerName = kafka.Domain
		kafkaInfo.Images = "hyperledger/fabric-kafka:latest"
		enviroment := make([]string, 0)
		enviroment = append(enviroment, fmt.Sprintf("KAFKA_BROKER_ID=%d", i+1))
		enviroment = append(enviroment, fmt.Sprintf("KAFKA_ZOOKEEPER_CONNECT=%s", zookeeperConnect))
		enviroment = append(enviroment, fmt.Sprintf("KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://%s:%d", kafka.Domain, kafka.Port))
		enviroment = append(enviroment, "ZKAFKA_MESSAGE_MAX_BYTES=1048576")
		enviroment = append(enviroment, "KAFKA_REPLICA_FETCH_MAX_BYTES=1048576")
		enviroment = append(enviroment, "KAFKA_UNCLEAN_LEADER_ELECTION_ENABLE=false")
		enviroment = append(enviroment, "KAFKA_LOG_RETENTION_MS=-1")
		enviroment = append(enviroment, "KAFKA_MIN_INSYNC_REPLICAS=2")
		enviroment = append(enviroment, "KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=3")
		enviroment = append(enviroment, "KAFKA_DEFAULT_REPLICATION_FACTOR=3")
		enviroment = append(enviroment, "KAFKA_LOG.DIRS=/opt/kafka/kafka-logs")
		kafkaInfo.Environment = enviroment
		volumes := make([]string, 0)
		volumes = append(volumes, fmt.Sprintf("/blockchainData/%s/kafka-logs:/opt/kafka/kafka-logs", kafka.Domain))
		kafkaInfo.Volumes = volumes
		ports := make([]string, 0)
		ports = append(ports, fmt.Sprintf("%d:9092", kafka.Port))
		kafkaInfo.Ports = ports
		kafkaInfo.ExtraHosts = zkExtraHosts
		service[kafka.Domain] = kafkaInfo
	}
	peerConnectHosts := make([]string, 0)
	for _, orderer := range general.Orderer {
		peerConnectHosts = append(peerConnectHosts, fmt.Sprintf("%s:%s", orderer.Domain, orderer.IP))
		ordererInfo := &objectdefine.DockerConfig{}
		ordererInfo.ContainerName = orderer.Domain
		ordererInfo.Images = "hyperledger/fabric-orderer:latest"
		enviroment := make([]string, 0)
		enviroment = append(enviroment, "GODEBUG=netdns=go")
		enviroment = append(enviroment, "ORDERER_GENERAL_LISTENADDRESS=0.0.0.0")
		enviroment = append(enviroment, "ORDERER_GENERAL_GENESISMETHOD=file")
		enviroment = append(enviroment, "ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block")
		enviroment = append(enviroment, "ORDERER_GENERAL_LOCALMSPID=OrdererMSP")
		enviroment = append(enviroment, "ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp")
		enviroment = append(enviroment, "ORDERER_GENERAL_TLS_ENABLED=true")
		enviroment = append(enviroment, "ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key")
		enviroment = append(enviroment, "ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt")
		enviroment = append(enviroment, "ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]")
		enviroment = append(enviroment, "ORDERER_KAFKA_TOPIC_REPLICATIONFACTOR=3")
		enviroment = append(enviroment, "ORDERER_KAFKA_BROKERS=[kafka0:9092,kafka1:9092,kafka2:9092,kafka3:9092]")
		enviroment = append(enviroment, "ORDERER_GENERAL_CLUSTER_CLIENTCERTIFICATE=/var/hyperledger/orderer/tls/server.crt")
		enviroment = append(enviroment, "ORDERER_GENERAL_CLUSTER_CLIENTPRIVATEKEY=/var/hyperledger/orderer/tls/server.key")
		enviroment = append(enviroment, "ORDERER_GENERAL_CLUSTER_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]")
		ordererInfo.Environment = enviroment
		ordererInfo.WorkingDir = "/opt/gopath/src/github.com/hyperledger/fabric"
		ordererInfo.Command = "orderer"
		volumes := make([]string, 0)
		volumes = append(volumes, "./channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block")
		volumes = append(volumes, fmt.Sprintf("./crypto-config/ordererOrganizations/%s/orderers/%s/msp:/var/hyperledger/orderer/msp", orderer.OrgDomain, orderer.Domain))
		volumes = append(volumes, fmt.Sprintf("./crypto-config/ordererOrganizations/%s/orderers/%s/tls:/var/hyperledger/orderer/tls", orderer.OrgDomain, orderer.Domain))
		volumes = append(volumes, fmt.Sprintf("/blockchainData/%s:/var/hyperledger/production/orderer", orderer.Domain))
		ordererInfo.Volumes = volumes
		ports := make([]string, 0)
		ports = append(ports, fmt.Sprintf("%d:7050", orderer.Port))
		ordererInfo.Ports = ports
		ordererInfo.ExtraHosts = zkExtraHosts
		service[orderer.Domain] = ordererInfo
	}
	orgPeerPairMap := make(map[string][]string, 0)
	for _, org := range general.Org {
		peerA := make([]string, 0)
		for _, peer := range org.Peer {
			peerConnectHosts = append(peerConnectHosts, fmt.Sprintf("%s:%s", peer.Domain, peer.IP))
			peerA = append(peerA, fmt.Sprintf("%s:%d", peer.Domain, peer.Port))
		}
		orgPeerPairMap[org.Name] = peerA
	}
	var peercli string
	for _, org := range general.Org {
		//ca
		caInfo := &objectdefine.DockerConfig{}
		caInfo.ContainerName = org.CA.Name
		caInfo.Images = "hyperledger/fabric-ca:latest"
		enviroment := make([]string, 0)
		enviroment = append(enviroment, "FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server")
		enviroment = append(enviroment, fmt.Sprintf("FABRIC_CA_SERVER_CA_NAME=%s", org.CA.Name))
		cifeCaPath := filepath.ToSlash(filepath.Join(general.BaseOutput, "crypto-config/peerOrganizations", org.OrgDomain, "ca"))
		var cifeCaCertPath string
		var cifeCaSKPath string
		filepath.Walk(cifeCaPath, func(path string, info os.FileInfo, err error) error {
			folder := info.Name()
			if strings.Contains(folder, "_sk") {
				cifeCaSKPath = folder
			} else if strings.Contains(folder, ".pem") {
				cifeCaCertPath = folder
			}
			return nil
		})
		enviroment = append(enviroment, fmt.Sprintf("FABRIC_CA_SERVER_CA_CERTFILE=/etc/hyperledger/fabric-ca-server-config/ca/%s", cifeCaCertPath))
		enviroment = append(enviroment, fmt.Sprintf("FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config/ca/%s", cifeCaSKPath))

		cifeTLSPath := filepath.ToSlash(filepath.Join(general.BaseOutput, "crypto-config/peerOrganizations", org.OrgDomain, "tlsca"))
		var cifeTLSCertPath string
		var cifeTLSSKPath string
		filepath.Walk(cifeTLSPath, func(path string, info os.FileInfo, err error) error {
			folder := info.Name()
			if strings.Contains(folder, "_sk") {
				cifeTLSSKPath = folder
			} else if strings.Contains(folder, ".pem") {
				cifeTLSCertPath = folder
			}
			return nil
		})
		enviroment = append(enviroment, fmt.Sprintf("FABRIC_CA_SERVER_TLS_CERTFILE=/etc/hyperledger/fabric-ca-server-config/tlsca/%s", cifeTLSCertPath))
		enviroment = append(enviroment, fmt.Sprintf("FABRIC_CA_SERVER_TLS_KEYFILE=/etc/hyperledger/fabric-ca-server-config/tlsca/%s", cifeTLSSKPath))
		enviroment = append(enviroment, "FABRIC_CA_SERVER_TLS_ENABLED=true")
		caInfo.Environment = enviroment
		caInfo.Command = "sh -c 'fabric-ca-server start -b admin:heitu-admin -d'"
		volumes := make([]string, 0)
		volumes = append(volumes, fmt.Sprintf("./crypto-config/peerOrganizations/%s/ca:/etc/hyperledger/fabric-ca-server-config/ca", org.OrgDomain))
		volumes = append(volumes, fmt.Sprintf("./crypto-config/peerOrganizations/%s/tlsca:/etc/hyperledger/fabric-ca-server-config/tlscas", org.OrgDomain))
		volumes = append(volumes, fmt.Sprintf("/blockchainData/%s:/etc/hyperledger/fabric-ca-server", org.CA.Name))
		caInfo.Volumes = volumes
		ports := make([]string, 0)
		ports = append(ports, fmt.Sprintf("%d:7054", org.CA.Port))
		caInfo.Ports = ports
		service[org.CA.Name] = caInfo
		//peer couchdb cli
		for _, peer := range org.Peer {
			//couchdb
			couchdbInfo := &objectdefine.DockerConfig{}
			couchdbInfo.ContainerName = fmt.Sprintf("couchdb-%s", peer.Domain)
			couchdbInfo.Images = "hyperledger/fabric-couchdb:latest"
			cenviroment := make([]string, 0)
			cenviroment = append(cenviroment, "COUCHDB_USER=")
			cenviroment = append(cenviroment, "COUCHDB_PASSWORD=")
			couchdbInfo.Environment = cenviroment
			cports := make([]string, 0)
			cports = append(cports, fmt.Sprintf("%d:5984", peer.CouchdbPort))
			couchdbInfo.Ports = cports
			cvolumes := make([]string, 0)
			cvolumes = append(cvolumes, fmt.Sprintf("/blockchainData/%s:/opt/couchdb/data", fmt.Sprintf("couchdb-%s", peer.Domain)))
			couchdbInfo.Volumes = cvolumes
			service[fmt.Sprintf("couchdb-%s", peer.Domain)] = couchdbInfo
			//peer
			peerInfo := &objectdefine.DockerConfig{}
			peerInfo.ContainerName = peer.Domain
			peerInfo.Images = "hyperledger/fabric-peer:latest"
			penviroment := make([]string, 0)
			penviroment = append(penviroment, fmt.Sprintf("CORE_PEER_ID=%s", peer.Domain))
			penviroment = append(penviroment, fmt.Sprintf("CORE_PEER_ADDRESS=%s:%d", peer.Domain, peer.Port))
			penviroment = append(penviroment, fmt.Sprintf("CORE_PEER_LISTENADDRESS=0.0.0.0:%d", peer.Port))
			penviroment = append(penviroment, fmt.Sprintf("CORE_PEER_CHAINCODEADDRESS=%s:%d", peer.Domain, peer.ChaincodePort))
			penviroment = append(penviroment, fmt.Sprintf("CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:%d", peer.ChaincodePort))
			var peerhostsAS string
			for _, ps := range orgPeerPairMap[org.Name] {
				if len(peerhostsAS) == 0 {
					peerhostsAS = ps
				} else {
					peerhostsAS = peerhostsAS + " " + ps
				}
			}
			penviroment = append(penviroment, fmt.Sprintf("CORE_PEER_GOSSIP_BOOTSTRAP=%s", peerhostsAS))
			penviroment = append(penviroment, fmt.Sprintf("CORE_PEER_GOSSIP_EXTERNALENDPOINT=%s:%d", peer.Domain, peer.Port))
			penviroment = append(penviroment, "CORE_PEER_GOSSIP_USELEADERELECTION=true")
			penviroment = append(penviroment, "CORE_PEER_GOSSIP_ORGLEADER=false")
			penviroment = append(penviroment, fmt.Sprintf("CORE_PEER_LOCALMSPID=%s", fmt.Sprintf("%sMSP", org.Name)))
			penviroment = append(penviroment, "CORE_PEER_TLS_ENABLED=true")
			penviroment = append(penviroment, "CORE_PEER_PROFILE_ENABLED=true")
			penviroment = append(penviroment, "CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt")
			penviroment = append(penviroment, "CORE_LEDGER_STATE_STATEDATABASE=CouchDB")
			penviroment = append(penviroment, fmt.Sprintf("CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=%s:5984", fmt.Sprintf("couchdb-%s", peer.Domain)))
			penviroment = append(penviroment, "CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=")
			penviroment = append(penviroment, "CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=")
			penviroment = append(penviroment, "GODEBUG=netdns=go")
			penviroment = append(penviroment, "CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock")
			penviroment = append(penviroment, "CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=baiyun_default")
			penviroment = append(penviroment, "FABRIC_LOGGING_SPEC=INFO")
			penviroment = append(penviroment, "CORE_CHAINCODE_LOGGING_LEVEL=DEBUG")
			penviroment = append(penviroment, "CORE_CHAINCODE_LOGGING_SHIM=DEBUG")
			penviroment = append(penviroment, "CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key")
			penviroment = append(penviroment, "CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt")
			peerInfo.Environment = penviroment
			peerInfo.WorkingDir = "/opt/gopath/src/github.com/hyperledger/fabric/peer"
			peerInfo.Command = "peer node start"
			pvolumes := make([]string, 0)
			pvolumes = append(pvolumes, "/var/run/:/host/var/run/")
			pvolumes = append(pvolumes, fmt.Sprintf("./crypto-config/peerOrganizations/%s/peers/%s/msp:/etc/hyperledger/fabric/msp", org.OrgDomain, peer.Domain))
			pvolumes = append(pvolumes, fmt.Sprintf("./crypto-config/peerOrganizations/%s/peers/%s/tls:/etc/hyperledger/fabric/tls", org.OrgDomain, peer.Domain))
			pvolumes = append(pvolumes, fmt.Sprintf("/blockchainData/%s:/var/hyperledger/production", peer.Domain))
			peerInfo.Volumes = pvolumes
			pports := make([]string, 0)
			pports = append(pports, fmt.Sprintf("%d:%d", peer.Port, peer.Port))
			peerInfo.Ports = pports
			pdependsOn := make([]string, 0)
			pdependsOn = append(pdependsOn, fmt.Sprintf("couchdb-%s", peer.Domain))
			peerInfo.DependsOn = pdependsOn
			peerInfo.ExtraHosts = peerConnectHosts
			service[peer.Domain] = peerInfo
			//cli
			if len(peercli) == 0 {
				peercli = peer.Name
				cliInfo := &objectdefine.DockerConfig{}
				cliInfo.ContainerName = fmt.Sprintf("cli-%s-%s", org.Name, peer.Name)
				cliInfo.Images = "hyperledger/fabric-tools:latest"
				cliInfo.TTY = true
				cliInfo.StdinOpen = true
				clienvironment := make([]string, 0)
				clienvironment = append(clienvironment, "SYS_CHANNEL=system_channel")
				clienvironment = append(clienvironment, "GOPATH=/opt/gopath")
				clienvironment = append(clienvironment, "CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock")
				clienvironment = append(clienvironment, "FABRIC_LOGGING_SPEC=DEBUG")
				clienvironment = append(clienvironment, "CORE_PEER_ID=cli")
				clienvironment = append(clienvironment, fmt.Sprintf("CORE_PEER_ADDRESS=%s:%d", peer.Domain, peer.Port))
				clienvironment = append(clienvironment, fmt.Sprintf("CORE_PEER_LOCALMSPID=%sMSP", org.Name))
				clienvironment = append(clienvironment, "CORE_PEER_TLS_ENABLED=true")
				clienvironment = append(clienvironment, fmt.Sprintf("CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/%s/peers/%s/tls/server.crt", org.OrgDomain, peer.Domain))
				clienvironment = append(clienvironment, fmt.Sprintf("CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/%s/peers/%s/tls/server.key", org.OrgDomain, peer.Domain))
				clienvironment = append(clienvironment, fmt.Sprintf("CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/%s/peers/%s/tls/ca.crt", org.OrgDomain, peer.Domain))
				clienvironment = append(clienvironment, fmt.Sprintf("CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/%s/users/Admin@%s/msp", org.OrgDomain, org.OrgDomain))
				cliInfo.Environment = clienvironment
				cliInfo.WorkingDir = "/opt/gopath/src/github.com/hyperledger/fabric/peer"
				cliInfo.Command = "/bin/bash"
				clivolumes := make([]string, 0)
				clivolumes = append(clivolumes, "/var/run/:/host/var/run/")
				clivolumes = append(clivolumes, "./chaincode/:/opt/gopath/src/github.com/chaincode")
				clivolumes = append(clivolumes, "/crypto-config:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/")
				clivolumes = append(clivolumes, "./scripts:/opt/gopath/src/github.com/hyperledger/fabric/peer/scripts/")
				clivolumes = append(clivolumes, "./channel-artifacts:/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts")
				cliInfo.Volumes = clivolumes
				cliInfo.ExtraHosts = peerConnectHosts
				service[fmt.Sprintf("cli-%s-%s", org.Name, peer.Name)] = cliInfo
			}
		}
	}
	ret.Service = service
	err := tools.WriteYamlFile(filepath.Join(general.BaseOutput, "base.yaml"), ret)
	if err != nil {
		err = errors.WithMessage(err, "Write sinoconfig.yaml file error")
		return err
	}
	return nil
}
