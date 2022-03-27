package daction

import (
	"deploy-server/app/objectdefine"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"io/ioutil"
	"testing"
	"encoding/json"
	"github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
	"deploy-server/app/tools"
	"runtime"
	"os/exec"
	"github.com/hyperledger/fabric/protos/orderer/etcdraft"
)

func Test_arrat(t *testing.T) {
	sPath := "D:/GoPath/src/deploy-server/output/sourceid/crypto-config/peerOrganizations/org1.example.com/tlsca/"
	var st string
	filepath.Walk(sPath, func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)        //打印path信息
		fmt.Println(info.Name()) //打印文件或目录名
		folder := info.Name()
		if strings.Contains(folder, "_sk") {
			st = folder
			fmt.Println(st)
		} else {
			fmt.Println("456")
		}

		return nil
	})
	fmt.Println("st", st)
}

func Test_RnName(t *testing.T) {
	file := "D:/GoPath/src/deploy-server/deploy.sh"
	err := os.Rename(file, "123.sh")
	if err != nil {

	}
}

func Test_MapTT(t *testing.T) {
	mapp := make(map[string]objectdefine.DeployType, 0)
	orgmp := make(map[string]objectdefine.OrgType, 0)
	orgmp["Org1"] = objectdefine.OrgType{}
	mapp["org"] = objectdefine.DeployType{
		JoinOrg: orgmp,
	}

	for dd, deploy := range mapp {
		for cc, _ := range deploy.JoinOrg {
			fmt.Println(dd)
			fmt.Println(cc)
		}
	}
}

func Test_checkIPFormat(t *testing.T) {
	ip := "114.256.114.114"
	address := net.ParseIP(ip)
	if address == nil {
		fmt.Println("22222", address)
	} else {
		fmt.Println("1111", address)
	}
}

func Test_FolderFile(t *testing.T) {
	sPath := "D:/testFolder"

	filepath.Walk(sPath, func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)        //打印path信息

		return nil
	})
	fmt.Println("end")
}

func Test_block(t *testing.T){
	indentPath := "D:/gopath/src/deploy-server/indent.json"
	b, err := ioutil.ReadFile(indentPath) // just pass the file name
    if err != nil {
        fmt.Print(err)
    }
	os.Setenv("FABRIC_CFG_PATH",filepath.ToSlash("D:/gopath/src/deploy-server"))
	general := &objectdefine.Indent{}
    err = json.Unmarshal(b, general)
    if err != nil {
        fmt.Println("解析数据失败", err)
        return
    }

	//GeneralCreateCompleteBlock(general,output)
	// template, err := GetGenesisCompleteBlockTemplate(general, output)
	// VersionRoot := dcache.GetTemplatePathByVersion(general.Version)
	path := "D:/gopath/src/deploy-server"

	template := localconfig.LoadTopLevel(path)
	if len(template.Organizations) == 0 {
		
	}
	//ret := &localconfig.TopLevel{}
    OrganizationsArray := make([]*localconfig.Organization,0)
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
	template.Organizations[1] = OrganizationsArray[0]
	for i:=1;i<len(general.Org);i++ {
		template.Organizations = append(template.Organizations,OrganizationsArray[1])
	}

	//template.Orderer = template.Orderer
	template.Orderer.OrdererType = general.Consensus
	Addresses := make([]string, 0, 32)
	for _, orderer := range general.Orderer {
		url := fmt.Sprintf("%s:%d", orderer.Domain, orderer.Port)
		Addresses = append(Addresses, url)
	}
	template.Orderer.Addresses = Addresses

	ret := make(map[string]*localconfig.Profile)
	//创世区块的联盟配置
	genesis := &localconfig.Profile{}
	genesis.Capabilities = template.Channel.Capabilities
	genesis.Orderer = template.Orderer
	ordererOrganizations := make([]*localconfig.Organization, 1)
	ordererOrganizations[0] = template.Organizations[0]
	//raft
	genesis.Orderer.OrdererType = "etcdraft"
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
	genesis.Orderer.EtcdRaft.Consenters = Consenters
	genesis.Orderer.Addresses = Addressesss
	genesis.Orderer.Organizations = ordererOrganizations
	genesis.Orderer.Capabilities = template.Capabilities["Orderer"]
	genesis.Application = template.Application
	genesis.Application.Organizations = ordererOrganizations
	genesis.Consortiums = make(map[string]*localconfig.Consortium)
	genesisOrganizations := make([]*localconfig.Organization, len(general.Org))

	for i := 0; i < len(general.Org); i++ {
		genesisOrganizations[i] = template.Organizations[i+1]
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
	channel.Application = template.Application
	channel.Application.Organizations = genesisOrganizations
	channel.Application.Capabilities = template.Capabilities["Application"]
	ret["TwoOrgsChannel"] = channel

	template.Profiles = ret
	configtPath := "D:/gopath/src/deploy-server"
	filename := filepath.Join(configtPath, "testconfigtx.yaml")
	err = tools.WriteYamlFile(filename, template)
	if nil != err {
		
	}
}

func Test_base(t *testing.T){

	indentPath := "D:/gopath/src/deploy-server/indent.json"
	b, err := ioutil.ReadFile(indentPath) // just pass the file name
    if err != nil {
        fmt.Print(err)
    }
	os.Setenv("FABRIC_CFG_PATH",filepath.ToSlash("D:/gopath/src/deploy-server"))
	general := &objectdefine.Indent{}
    err = json.Unmarshal(b, general)
    if err != nil {
        fmt.Println("解析数据失败", err)
        return
    }

	//GeneralCreateCompleteBlock(general,output)
	// template, err := GetGenesisCompleteBlockTemplate(general, output)
	// VersionRoot := dcache.GetTemplatePathByVersion(general.Version)
	path := "D:/gopath/src/deploy-server"

	template := localconfig.LoadTopLevel(path)
	if len(template.Organizations) == 0 {
		
	}
	genesis := CreateGenesisCompleteBlockConfig(general, template)
	configtPath := "D:/gopath/src/deploy-server"
	filename := filepath.Join(configtPath, "testconfigtx.yaml")
	err = tools.WriteYamlFile(filename, genesis)
	if nil != err {
		
	}
}

func Test_blockFiel(t *testing.T){
	indentPath := "D:/gopath/src/deploy-server/indent.json"
	b, err := ioutil.ReadFile(indentPath) // just pass the file name
    if err != nil {
        fmt.Print(err)
    }
	os.Setenv("FABRIC_CFG_PATH",filepath.ToSlash("D:/gopath/src/deploy-server"))
	general := &objectdefine.Indent{}
    err = json.Unmarshal(b, general)
    if err != nil {
        fmt.Println("解析数据失败", err)
        return
    }

	//GeneralCreateCompleteBlock(general,output)
	// template, err := GetGenesisCompleteBlockTemplate(general, output)
	// VersionRoot := dcache.GetTemplatePathByVersion(general.Version)


	path := "D:/gopath/src/deploy-server"

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

	//Policies := make(map[string]map[string]bool)
	
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

	// template.Organizations[1] = OrganizationsArray[0]
	// for i:=1;i<len(general.Org);i++ {
	// 	template.Organizations = append(template.Organizations,OrganizationsArray[1])
	// }
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
		url := fmt.Sprintf("%s:%d", orderer.Domain, orderer.Port)
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


	// ret := make(map[string]*localconfig.Profile)


	// //创世区块的联盟配置
	// genesis := &localconfig.Profile{}
	// //genesis.Capabilities = template.Channel.Capabilities
	// genesis.Orderer = template.Orderer
	// ordererOrganizations := make([]*localconfig.Organization, 1)
	// ordererOrganizations[0] = template.Organizations[0]
	// //raft
	// // genesis.Orderer.OrdererType = "etcdraft"
	// // Consenters := make([]*etcdraft.Consenter,0)
	// // Addressesss := make([]string,0)
	// // for _, orderer := range general.Orderer{
	// // 	consenter := &etcdraft.Consenter{
	// // 		Host: orderer.Domain,
	// // 		Port: uint32(orderer.Port),
	// // 		ClientTlsCert: []byte(fmt.Sprintf("crypto-config/ordererOrganizations/example.com/orderers/%s/tls/server.crt",orderer.Domain)),
	// // 	    ServerTlsCert: []byte(fmt.Sprintf("crypto-config/ordererOrganizations/example.com/orderers/%s/tls/server.crt",orderer.Domain)),
	// // 	}
	// // 	Addressesss = append(Addressesss,fmt.Sprintf("%s:%d",orderer.Domain,orderer.Port))
	// // 	Consenters= append(Consenters,consenter)
	// // }
	// genesis.Orderer.EtcdRaft.Consenters = Consenters
	// genesis.Orderer.Addresses = Addressesss
	// genesis.Orderer.Organizations = ordererOrganizations
	// //genesis.Orderer.Capabilities = template.Capabilities["Orderer"]
	// genesis.Application = template.Application
	// genesis.Application.Organizations = ordererOrganizations
	// genesis.Consortiums = make(map[string]*localconfig.Consortium)
	// genesisOrganizations := make([]*localconfig.Organization, len(general.Org))

	// for i := 0; i < len(general.Org); i++ {
	// 	genesisOrganizations[i] = template.Organizations[i+1]
	// }
	// // genesis.Consortiums = map[string]*localconfig.Consortium{
	// // 	"PrisonConsortium": &localconfig.Consortium{
	// // 		Organizations: genesisOrganizations,
	// // 	},
	// // }
	// genesis.Consortiums = map[string]*localconfig.Consortium{
	// 	"SampleConsortium": &localconfig.Consortium{
	// 		Organizations: genesisOrganizations,
	// 	},
	// }
	// //ret["PrisonNetworkGenesis"] = genesis
	// ret["SampleMultiNodeEtcdRaft"] = genesis
	
	// channel := &localconfig.Profile{}
	// channel.Consortium = "SampleConsortium"
	// channel.Application = template.Application
	// channel.Application.Organizations = genesisOrganizations
	// channel.Application.Capabilities = template.Capabilities["Application"]
	// ret["TwoOrgsChannel"] = channel

	// template.Profiles = ret
	configtPath := "D:/gopath/src/deploy-server"
	filename := filepath.Join(configtPath, "testconfigtx.yaml")
	err = tools.WriteYamlFile(filename, template)
	if nil != err {
		
	}
}



func Test_blockSedl(t *testing.T){
	indentPath := "D:/gopath/src/deploy-server/indent.json"
	b, err := ioutil.ReadFile(indentPath) // just pass the file name
    if err != nil {
        fmt.Print(err)
    }
	os.Setenv("FABRIC_CFG_PATH",filepath.ToSlash("D:/gopath/src/deploy-server"))
	general := &objectdefine.Indent{}
    err = json.Unmarshal(b, general)
    if err != nil {
        fmt.Println("解析数据失败", err)
        return
    }
    ConfigtxPath := "D:/gopath/src/deploy-server/"
	var ordererOrgDomain string
	ordererHost := make([]string,0)
	Consenters := make(map[string][]interface{},len(general.Orderer))
	for _, orderer := range general.Orderer {
        ordererOrgDomain = orderer.OrgDomain
		
		ordererHost = append(ordererHost,fmt.Sprintf("%s:%d",orderer.Domain,orderer.Port))
		consenttersValue := make([]interface{},3)

		consenttersValue = append(consenttersValue,orderer.Port)
		ordererCtlsCert := fmt.Sprintf("crypto-config/ordererOrganizations/%s/orderers/%s/tls/server.crt",orderer.OrgDomain,orderer.Domain)
	    ordererStlsCert := fmt.Sprintf("crypto-config/ordererOrganizations/%s/orderers/%s/tls/server.crt",orderer.OrgDomain,orderer.Domain)
		consenttersValue = append(consenttersValue,ordererCtlsCert)
		consenttersValue = append(consenttersValue,ordererStlsCert)
 
		Consenters[orderer.Domain]=consenttersValue	
       // ordererDomain := orderer.Domain
		//ordererPort := orderer.Port	
	}

	var cmd *exec.Cmd 
	execCommand := fmt.Sprintf("cd %s && sed -i 's/ORDERERORGDOMAIN/"+ordererOrgDomain+"/g' configtx.yaml",ConfigtxPath)
	fmt.Println("execCommand",execCommand)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", execCommand)
	} else {
		cmd = exec.Command("/bin/bash", "-c", execCommand)
	}
	if err := cmd.Start(); err != nil {
        fmt.Println("11111",err)
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("11111",err)
	}
    var SourceV string
	var change string
	for i,v := range ordererHost{

        if i == 0{
			var cmd *exec.Cmd 
            SourceV = fmt.Sprintf("- %s",v)
			execCommand := fmt.Sprintf("cd %s && sed -i 's/ORDERERHOST/"+SourceV+"/g' configtx.yaml",ConfigtxPath)
			fmt.Println("execCommand",execCommand)
			if runtime.GOOS == "windows" {
				cmd = exec.Command("cmd", "/c", execCommand)
			} else {
				cmd = exec.Command("/bin/bash", "-c", execCommand)
			}
			if err := cmd.Start(); err != nil {
				fmt.Println("11111",err)
			}
			if err := cmd.Wait(); err != nil {
				fmt.Println("11111",err)
			}
		}else{
			var cmd *exec.Cmd 
			if len(change) == 0{
				execCommand := fmt.Sprintf("cd %s && sed -i '/"+SourceV+"/a\\        - "+v+"' configtx.yaml",ConfigtxPath)
				fmt.Println("execCommand",execCommand)
				if runtime.GOOS == "windows" {
					cmd = exec.Command("cmd", "/c", execCommand)
				} else {
					cmd = exec.Command("/bin/bash", "-c", execCommand)
				}
				if err := cmd.Start(); err != nil {
					fmt.Println("11111",err)
				}
				if err := cmd.Wait(); err != nil {
					fmt.Println("11111",err)
				}
				change = fmt.Sprintf("- %s",v)
			}else{
				//var cmd *exec.Cmd 
				execCommand := fmt.Sprintf("cd %s && sed -i '/"+change+"/a\\        - "+v+"' configtx.yaml",ConfigtxPath)
				fmt.Println("execCommand",execCommand)
				if runtime.GOOS == "windows" {
					cmd = exec.Command("cmd", "/c", execCommand)
				} else {
					cmd = exec.Command("/bin/bash", "-c", execCommand)
				}
				if err := cmd.Start(); err != nil {
					fmt.Println("11111",err)
				}
				if err := cmd.Wait(); err != nil {
					fmt.Println("11111",err)
				}
				change = fmt.Sprintf("- %s",v)
			}
            
			
		}
	}
}

func Test_writeFile(t *testing.T){
	channFilePath :="../../channelFile.txt"
	f, err := os.OpenFile(channFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
	  fmt.Println("11111111")
	} else {
	   n, _ := f.Seek(0, os.SEEK_END)
	   _, err = f.WriteAt([]byte("#!/bin/bash\ncurPath=$PWD\nCCFILE=$1\n## 构建go mod\ngo mod init $CCFILE\n\n## 检测\ngo mod tidy \n\n## vendor\ngo bulid \n"), n)
	   fmt.Println("write succeed!")
	   defer f.Close()
	}
}
func Test_readFile(t *testing.T){
	channFilePath :="../../channelFile.txt"
	channName, err := ioutil.ReadFile(channFilePath) // just pass the file name
		if err != nil {
			fmt.Print(err)
		}

     chann:= string(channName)
	 fmt.Println(chann)
}