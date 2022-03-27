package dmysql

import (
	"deploy-server/app/objectdefine"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

//MakeStepGeneralUpdateCreateCompleteDeployIndent  首先把数据插入数据库
func MakeStepGeneralUpdateCreateCompleteDeployIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create channel info insert mysql ommp_indent table start"},
		})
		err := insertCreateCompleteDeployIndent(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create channel info insert mysql ommp_indent table end"},
		})
		return nil
	}
}


//insertCreateCompleteDeployIndent 通道数据插入数据库语句以及插入数据
func insertCreateCompleteDeployIndent(indent *objectdefine.Indent) error {
	taskID := indent.ID
	sourceID := indent.SourceID
	// if len(sourceID) == 0 {
	// 	sourceID = "sourceid"
	// }
	channelname := indent.ChannelName
	consensus := indent.Consensus
	version := indent.Version
	for _, org := range indent.Org {
		orgName := org.Name
		orgDomain := org.OrgDomain
		var peerName, peerNickName, peerAccessKey, peerIP, peerUser, peerDomain, cliName, caIP, caName string
		var peerID, peerPort, couchdbPort, ccPort, caPort int
		caIP = org.CA.IP
		caName = org.CA.Name
		caPort = org.CA.Port
		for _, peer := range org.Peer {
			peerIP = peer.IP
			peerPort = peer.Port
			peerName = peer.Name
			peerNickName = peer.NickName
			peerAccessKey = peer.AccessKey
			peerUser = peer.User
			peerID = peer.PeerID
			peerDomain = peer.Domain
			cliName = peer.CliName
			couchdbPort = peer.CouchdbPort
			ccPort = peer.ChaincodePort
			sqlI := "insert into ommp_indent (task_id,source_id,channelname,consensus,org_name,org_domain,peer_ip,peer_name,nick_name,accesskey,peer_user,peer_id,peer_port,couchdb_port,cc_port,peer_domain,ca_name,ca_ip,ca_port,cli_name,version,channel_status,org_status,peer_status,peer_run_status) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
			result, err := db.Exec(sqlI, taskID, sourceID, channelname, consensus, orgName, orgDomain, peerIP, peerName, peerNickName, peerAccessKey, peerUser, peerID, peerPort, couchdbPort, ccPort, peerDomain, caName, caIP, caPort, cliName, version,1,1,1,1)
			if err != nil {
				return errors.Errorf("channel info insert mysql ommp_indent fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("channel info insert mysql ommp_indent fail")
			}
		}
	}

	for _,orderer := range indent.Orderer{
		sqlI := "insert into ommp_orderer (channelname,orderer_name,orderer_domain,orderer_ip,orderer_port,orderer_orgdomain) values (?,?,?,?,?,?)"
		result, err := db.Exec(sqlI, channelname, orderer.Name, orderer.Domain, orderer.IP, orderer.Port, orderer.OrgDomain)
		if err != nil {
			return errors.Errorf("channel info insert mysql ommp_indent fail %s", err)
		}
		rows, _ := result.RowsAffected()
		if rows <= 0 {
			return errors.New("channel info insert mysql ommp_indent fail")
		}
		
	}
	return nil

}

//MakeStepGeneralInsertCreateChannelIndent 首先把数据插入数据库
func MakeStepGeneralInsertCreateChannelIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create channel info insert mysql ommp_indent table start"},
		})
		err := insertCreateChannelIndent(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create channel info insert mysql ommp_indent table end"},
		})
		return nil
	}
}

//MakeStepGeneralUpdateCreateChannelIndent 执行成功之后 更新数据库数据状态
func MakeStepGeneralUpdateCreateChannelIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"insert mysql ommp_indent table start"},
		})
		err := updateSuccessCreateChannelTaskStatus(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := UpdateFailCreateChannelTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChannelUpdateDB",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"insert mysql ommp_indent table end"},
		})
		return nil
	}
}

//MakeStepGeneralInsertAddOrgIndent 首先把数据插入数据库
func MakeStepGeneralInsertAddOrgIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create org info insert mysql ommp_indent table start"},
		})
		err := insertAddOrgIndent(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create org info insert mysql ommp_indent table end"},
		})
		return nil
	}
}

//MakeStepGeneralUpdateDelOrgIndent 删除组织 更改状态
func MakeStepGeneralUpdateDelOrgIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_indent table start"},
		})
		err := updateDeleteOrgIndent(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_indent table end"},
		})
		return nil
	}
}

//MakeStepGeneralUpdateAddOrgIndent 执行成功之后 更新数据库数据状态
func MakeStepGeneralUpdateAddOrgIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_indent table start"},
		})
		err := updateSuccessCreateOrgTaskStatus(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := UpdateFailCreateOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateOrgUpdateDB",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_indent table end"},
		})
		return nil
	}
}

//MakeStepGeneralUpdateDeleteOrgIndent 执行成功之后 更新数据库数据状态
func MakeStepGeneralUpdateDeleteOrgIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_indent table start"},
		})
		err := updateSuccessDeleteOrgTaskStatus(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := UpdateFailDeleteOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeleteOrgUpdateDB",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_indent table end"},
		})
		return nil
	}
}

//MakeStepGeneralInsertAddPeerIndent 首先把数据插入数据库
func MakeStepGeneralInsertAddPeerIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"insert mysql ommp_indent table start"},
		})
		err := insertAddPeerIndent(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"insert mysql ommp_indent table end"},
		})
		return nil
	}
}

//MakeStepGeneralUpdateDelpeerIndent 删除节点 更改状态
func MakeStepGeneralUpdateDelpeerIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_indent table start"},
		})
		err := updateDeletePeerIndent(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_indent table end"},
		})
		return nil
	}
}

//MakeStepGeneralUpdatePeerIndent 执行成功之后 更新数据库数据状态
func MakeStepGeneralUpdatePeerIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"insert mysql ommp_indent table start"},
		})
		err := updateSuccessCreatePeerTaskStatus(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := UpdateFailCreatePeerTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreatePeerUpdateDB",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"insert mysql ommp_indent table end"},
		})
		return nil
	}
}

//MakeStepGeneralUpdateDeletePeerIndent 执行成功之后 更新数据库数据状态
func MakeStepGeneralUpdateDeletePeerIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_indent table start"},
		})
		err := updateSuccessDeletePeerTaskStatus(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := UpdateFailDeletePeerTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeletePeerUpdateDB",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_indent table end"},
		})
		return nil
	}
}

//MakeStepGeneralUpdatePeerDisableStatus 更新节点数据状态为停用
func MakeStepGeneralUpdatePeerDisableStatus(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_indent table peer status start"},
		})
		err := updatePeerDisableStatus(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_indent table  peer status end"},
		})
		return nil
	}
}

//MakeStepGeneralUpdatePeerEnableStatus 更新节点数据状态为启用
func MakeStepGeneralUpdatePeerEnableStatus(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_indent table peer status start"},
		})
		err := updatePeerEnableStatus(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_indent table  peer status end"},
		})
		return nil
	}
}

//MakeStepGeneralUpdatePeerModiflyNickName 更新节点昵称
func MakeStepGeneralUpdatePeerModiflyNickName(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_indent table peer nickname start"},
		})
		err := UpdateIndetPeerNickName(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_indent table  peer nickname end"},
		})
		return nil
	}
}

//MakeStepGeneralInsertAddChainCodeIndent 首先把数据插入数据库
func MakeStepGeneralInsertAddChainCodeIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"insert mysql ommp_chaincode table start"},
		})
		err := insertAddChainCodeIndent(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"insert mysql ommp_chaincode table  end"},
		})
		return nil
	}
}

//MakeStepGeneralUpdateDeleteChainCodeIndent 删除合约先把状态置为正在删除
func MakeStepGeneralUpdateDelChainCodeIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_chaincode table status start"},
		})
		err := updateChainCodeDeleteStatus(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_chaincode table status end"},
		})
		return nil
	}
}

//MakeStepGeneralUpdateAddChainCodeIndent  执行新增合约成功之后 更新数据库数据状态
func MakeStepGeneralUpdateAddChainCodeIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_chaincode table start"},
		})
		err := updateSuccessCreateChaincodeTaskStatus(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := UpdateFailCreateChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateChainCodeUpdateDB",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_chaincode table end"},
		})
		return nil
	}
}

//MakeStepGeneralUpdateDeleteChainCodeIndent  执行删除合约成功之后 更新数据库数据状态
func MakeStepGeneralUpdateDeleteChainCodeIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_chaincode table status start"},
		})
		err := updateSuccessDeleteChaincodeTaskStatus(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := UpdateFailDeleteChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeleteChainCodeUpdateDB",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_chaincode table status end"},
		})
		return nil
	}
}

//MakeStepGeneralInsertUpgradeChainCodeIndent  首先把数据插入数据库
func MakeStepGeneralInsertUpgradeChainCodeIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"insert mysql ommp_chaincode table start"},
		})
		err := UpdateUpgradeChainCodeIndent(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"insert mysql ommp_chaincode table  end"},
		})
		return nil
	}
}

//MakeStepGeneralDisableChainCodeUpdateIndent  停用合约更新状态为正在停用
func MakeStepGeneralDisableChainCodeUpdateIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_chaincode table start"},
		})
		err := UpdateDisableChainCodeIndent(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_chaincode table  end"},
		})
		return nil
	}
}

//MakeStepGeneralEnableChainCodeUpdateIndent  启用合约更新状态为正在启用
func MakeStepGeneralEnableChainCodeUpdateIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_chaincode table start"},
		})
		err := UpdateEnableChainCodeIndent(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_chaincode table  end"},
		})
		return nil
	}
}

//MakeStepGeneralUpdateUpgradeChainCodeIndent 执行成功之后 更新数据库数据状态
func MakeStepGeneralUpdateUpgradeChainCodeIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_chaincode table start"},
		})
		err := updateSuccessUpgradeChaincodeTaskStatus(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := UpdateFailUpgradeChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalUpgradeChainCodeUpdateDB",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_chaincode table end"},
		})
		return nil
	}
}

//MakeStepGeneralUpdateDisableChainCodeIndent 更新合约为停用状态
func MakeStepGeneralUpdateDisableChainCodeIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_chaincode table start"},
		})
		err := UpdateSuccessDisableChainCodeIndent(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := UpdateFailDisableChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDisableChainCodeWriteDB",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_chaincode table  end"},
		})
		return nil
	}
}

//MakeStepGeneralUpdateEnableChainCodeIndent 更新合约状态为启用
func MakeStepGeneralUpdateEnableChainCodeIndent(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_chaincode table start"},
		})
		err := UpdateSuccessEnableChainCodeIndent(general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := UpdateFailEnableChaincodeTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalEnableChainCodeWriteDB",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"update mysql ommp_chaincode table  end"},
		})
		return nil
	}
}

//insertCreateChannelIndent 通道数据插入数据库语句以及插入数据
func insertCreateChannelIndent(indent *objectdefine.Indent) error {
	taskID := indent.ID
	sourceID := indent.SourceID
	if len(sourceID) == 0 {
		sourceID = "sourceid"
	}
	channelname := indent.ChannelName
	consensus := indent.Consensus
	version := indent.Version
	for _, org := range indent.Org {
		orgName := org.Name
		orgDomain := org.OrgDomain
		var peerName, peerNickName, peerAccessKey, peerIP, peerUser, peerDomain, cliName, caIP, caName string
		var peerID, peerPort, couchdbPort, ccPort, caPort int
		caIP = org.CA.IP
		caName = org.CA.Name
		caPort = org.CA.Port
		for _, peer := range org.Peer {
			peerIP = peer.IP
			peerPort = peer.Port
			peerName = peer.Name
			peerNickName = peer.NickName
			peerAccessKey = peer.AccessKey
			peerUser = peer.User
			peerID = peer.PeerID
			peerDomain = peer.Domain
			cliName = peer.CliName
			couchdbPort = peer.CouchdbPort
			ccPort = peer.ChaincodePort
			sqlI := "insert into ommp_indent (task_id,source_id,channelname,consensus,org_name,org_domain,peer_ip,peer_name,nick_name,accesskey,peer_user,peer_id,peer_port,couchdb_port,cc_port,peer_domain,ca_name,ca_ip,ca_port,cli_name,version,peer_run_status) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
			result, err := db.Exec(sqlI, taskID, sourceID, channelname, consensus, orgName, orgDomain, peerIP, peerName, peerNickName, peerAccessKey, peerUser, peerID, peerPort, couchdbPort, ccPort, peerDomain, caName, caIP, caPort, cliName, version, 1)
			if err != nil {
				return errors.Errorf("channel info insert mysql ommp_indent fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("channel info insert mysql ommp_indent fail")
			}
		}
	}
	return nil
}

//insertAddOrgIndent 组织数据插入数据库语句以及插入数据
func insertAddOrgIndent(indent *objectdefine.Indent) error {
	taskID := indent.ID
	sourceID := indent.SourceID
	if len(sourceID) == 0 {
		sourceID = "sourceid"
	}
	channelname := indent.ChannelName
	consensus := indent.Consensus
	version := indent.Version
	for _, org := range indent.Org {
		orgName := org.Name
		orgDomain := org.OrgDomain
		var peerName, peerNickName, peerAccessKey, peerIP, peerUser, peerDomain, cliName, caIP, caName string
		var peerID, peerPort, couchdbPort, ccPort, caPort int
		caIP = org.CA.IP
		caName = org.CA.Name
		caPort = org.CA.Port
		for _, peer := range org.Peer {
			peerIP = peer.IP
			peerPort = peer.Port
			peerName = peer.Name
			peerNickName = peer.NickName
			peerAccessKey = peer.AccessKey
			peerUser = peer.User
			peerID = peer.PeerID
			fmt.Println("new add org peerID", peerID)
			peerDomain = peer.Domain
			cliName = peer.CliName
			couchdbPort = peer.CouchdbPort
			ccPort = peer.ChaincodePort
			sqlI := "insert into ommp_indent (task_id,source_id,channelname,consensus,org_name,org_domain,peer_ip,peer_name,nick_name,accesskey,peer_user,peer_id,peer_port,couchdb_port,cc_port,peer_domain,ca_name,ca_ip,ca_port,cli_name,version,channel_status) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
			result, err := db.Exec(sqlI, taskID, sourceID, channelname, consensus, orgName, orgDomain, peerIP, peerName, peerNickName, peerAccessKey, peerUser, peerID, peerPort, couchdbPort, ccPort, peerDomain, caName, caIP, caPort, cliName, version, 1)
			if err != nil {
				return errors.Errorf("mysql insert ommp_indent fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql insert ommp_indent fail")
			}
		}
	}

	return nil
}

//updateDeleteOrgIndent 删除组织数据语句以及更新状态
func updateDeleteOrgIndent(indent *objectdefine.Indent) error {
	channelname := indent.ChannelName

	for _, org := range indent.Org {
		orgName := org.Name
		orgDomain := org.OrgDomain
		//for _, peer := range org.Peer {
		//	peerDomain := peer.Domain
		sqlI := "update ommp_indent set org_status=?,peer_status=?,peer_run_status=? where channelname=? and org_name=? and org_domain=?"
		result, err := db.Exec(sqlI, 3, 3, 0, channelname, orgName, orgDomain)
		if err != nil {
			return errors.Errorf("mysql update ommp_indent fail %s", err)
		}
		rows, _ := result.RowsAffected()
		if rows <= 0 {
			return errors.New("mysql update ommp_indent fail")
		}
		//}
	}

	return nil
}

//insertAddPeerIndent 节点数据插入数据库语句以及插入数据
func insertAddPeerIndent(indent *objectdefine.Indent) error {
	taskID := indent.ID
	sourceID := indent.SourceID
	if len(sourceID) == 0 {
		sourceID = "sourceid"
	}
	channelname := indent.ChannelName
	consensus := indent.Consensus
	version := indent.Version
	for _, org := range indent.Org {
		orgName := org.Name
		orgDomain := org.OrgDomain
		var peerName, peerNickName, peerAccessKey, peerIP, peerUser, peerDomain, cliName, caIP, caName string
		var peerID, peerPort, couchdbPort, ccPort, caPort int
		caIP = org.CA.IP
		caName = org.CA.Name
		caPort = org.CA.Port
		for _, peer := range org.Peer {
			peerIP = peer.IP
			peerPort = peer.Port
			peerName = peer.Name
			peerNickName = peer.NickName
			peerAccessKey = peer.AccessKey
			peerUser = peer.User
			peerID = peer.PeerID
			peerDomain = peer.Domain
			cliName = peer.CliName
			couchdbPort = peer.CouchdbPort
			ccPort = peer.ChaincodePort
			sqlI := "insert into ommp_indent (task_id,source_id,channelname,consensus,org_name,org_domain,peer_ip,peer_name,nick_name,accesskey,peer_user,peer_id,peer_port,couchdb_port,cc_port,peer_domain,ca_name,ca_ip,ca_port,cli_name,version,channel_status,org_status) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
			result, err := db.Exec(sqlI, taskID, sourceID, channelname, consensus, orgName, orgDomain, peerIP, peerName, peerNickName, peerAccessKey, peerUser, peerID, peerPort, couchdbPort, ccPort, peerDomain, caName, caIP, caPort, cliName, version, 1, 1)
			if err != nil {
				return errors.Errorf("mysql insert ommp_indent fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql insert ommp_indent fail")
			}
		}
	}

	return nil
}

//updateDeletePeerIndent 删除节点数据语句以及更新状态
func updateDeletePeerIndent(indent *objectdefine.Indent) error {
	channelname := indent.ChannelName
	for _, org := range indent.Org {
		orgName := org.Name
		for _, peer := range org.Peer {
			sqlI := "update ommp_indent set peer_status=?,peer_run_status=? where channelname=? and org_name=? and peer_domain=?"
			result, err := db.Exec(sqlI, 3, 0, channelname, orgName, peer.Domain)
			if err != nil {
				return errors.Errorf("mysql update ommp_indent fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update ommp_indent fail")
			}
		}
	}

	return nil
}

//updatePeerDisableStatus 更新数据库节点状态信息信息
func updatePeerDisableStatus(indent *objectdefine.Indent) error {
	for _, org := range indent.Org {
		for _, peer := range org.Peer {
			peerIP := peer.IP
			peerName := peer.Name
			peerDomain := peer.Domain

			sqlU := "update ommp_indent set peer_run_status=? where peer_ip=? and peer_name=? and peer_domain=?"
			result, err := db.Exec(sqlU, 0, peerIP, peerName, peerDomain)
			if err != nil {
				return errors.Errorf("mysql update ommp_indent fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update ommp_indent fail")
			}
		}
	}

	return nil
}

//updatePeerEnableStatus 更新数据库节点状态信息信息
func updatePeerEnableStatus(indent *objectdefine.Indent) error {
	for _, org := range indent.Org {
		for _, peer := range org.Peer {
			peerIP := peer.IP
			peerName := peer.Name
			peerDomain := peer.Domain

			sqlU := "update ommp_indent set peer_run_status=? where peer_ip=? and peer_name=? and peer_domain=?"
			result, err := db.Exec(sqlU, 1, peerIP, peerName, peerDomain)
			if err != nil {
				return errors.Errorf("mysql update ommp_indent fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update ommp_indent fail")
			}
		}
	}

	return nil
}

//insertAddChainCodeIndent 合约数据 插入数据库语句以及插入数据
func insertAddChainCodeIndent(indent *objectdefine.Indent) error {
	taskID := indent.ID
	for ccName, ccA := range indent.Chaincode {
		for _, cc := range ccA {
			ccVersion := cc.Version
			ccPolicy := cc.Policy
			ccChannel := indent.ChannelName
			var ccDesc string
			if len(cc.Describe) != 0 {
				ccDesc = cc.Describe
			}
			endorsejson, _ := json.Marshal(cc.EndorsementOrg)
			endorseString := string(endorsejson)
			//sqlI := "insert into chaincode (cc_name,cc_version,cc_policy,cc_org,channelname,detail) values (?,?,?,?,?,?)"
			sqlU := "update ommp_chaincode set task_id=?,cc_policy=?, cc_org=?,channelname=?,detail=?,is_install=?,status=? where cc_name=? and cc_version=? and is_install=? and status=?"
			result, err := db.Exec(sqlU, taskID, ccPolicy, endorseString, ccChannel, ccDesc, 0, 0, ccName, ccVersion, 0, 0)
			if err != nil {
				return errors.Errorf("mysql update ommp_chaincode fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update ommp_chaincodet fail")
			}
		}
	}

	return nil
}

//updateChainCodeDeleteStatus 删除合约 更新数据库状态信息
func updateChainCodeDeleteStatus(indent *objectdefine.Indent) error {
	for ccName, ccA := range indent.Chaincode {
		for _, cc := range ccA {
			endorsejson, _ := json.Marshal(cc.EndorsementOrg)
			endorseString := string(endorsejson)
			//sqlI := "insert into chaincode (cc_name,cc_version,cc_policy,cc_org,channelname,detail) values (?,?,?,?,?,?)"
			sqlU := "update ommp_chaincode set is_install=?,status=? where cc_name=? and cc_version=? and channelname=? and cc_org=?"
			result, err := db.Exec(sqlU, 0, 3, ccName, cc.Version, indent.ChannelName, endorseString)
			if err != nil {
				return errors.Errorf("mysql update ommp_chaincode fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update ommp_chaincode fail")
			}
		}
	}

	return nil
}

//UpdateUpgradeChainCodeIndent 升级合于 更新数据库信息
func UpdateUpgradeChainCodeIndent(indent *objectdefine.Indent) error {
	taskID := indent.ID
	for ccName, ccA := range indent.Chaincode {
		for _, cc := range ccA {
			ccVersion := cc.Version
			ccPolicy := cc.Policy
			ccChannel := indent.ChannelName
			var ccDesc string
			if len(cc.Describe) != 0 {
				ccDesc = cc.Describe
			}
			endorsejson, _ := json.Marshal(cc.EndorsementOrg)
			endorseString := string(endorsejson)
			//sqlI := "insert into chaincode (cc_name,cc_version,cc_policy,cc_org,channelname,detail) values (?,?,?,?,?,?)"
			sqlU := "update ommp_chaincode set task_id=?,cc_version=?,cc_policy=?, cc_org=?,detail=?,is_install=?,status=? where cc_name=? and channelname=?"
			result, err := db.Exec(sqlU, taskID, ccVersion, ccPolicy, endorseString, ccDesc, 0, 0, ccName, ccChannel)
			if err != nil {
				return errors.Errorf("mysql update ommp_chaincode fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update ommp_chaincode fail")
			}
		}
	}

	return nil
}

//UpdateDisableChainCodeIndent 停用合约 更新数据库状态
func UpdateDisableChainCodeIndent(indent *objectdefine.Indent) error {
	for ccName, ccA := range indent.Chaincode {
		for _, cc := range ccA {
			endorsejson, _ := json.Marshal(cc.EndorsementOrg)
			endorseString := string(endorsejson)
			//sqlI := "insert into chaincode (cc_name,cc_version,cc_policy,cc_org,channelname,detail) values (?,?,?,?,?,?)"
			sqlU := "update ommp_chaincode set status=? where cc_name=? and cc_version=? and channelname=? and  cc_org=?"
			result, err := db.Exec(sqlU, 3, ccName, cc.Version, indent.ChannelName, endorseString)
			if err != nil {
				return errors.Errorf("mysql update ommp_chaincode fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update ommp_chaincode fail")
			}
		}
	}
	return nil
}

//UpdateEnableChainCodeIndent 启用合约 更新数据库状态
func UpdateEnableChainCodeIndent(indent *objectdefine.Indent) error {
	for ccName, ccA := range indent.Chaincode {
		for _, cc := range ccA {
			endorsejson, _ := json.Marshal(cc.EndorsementOrg)
			endorseString := string(endorsejson)
			//sqlI := "insert into chaincode (cc_name,cc_version,cc_policy,cc_org,channelname,detail) values (?,?,?,?,?,?)"
			sqlU := "update ommp_chaincode set status=? where cc_name=? and cc_version=? and channelname=? and  cc_org=?"
			result, err := db.Exec(sqlU, 3, ccName, cc.Version, indent.ChannelName, endorseString)
			if err != nil {
				return errors.Errorf("mysql update ommp_chaincode fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update ommp_chaincode fail")
			}
		}
	}
	return nil
}

//ReceiveChainCodeUploadFile 接受包存链码文件压缩包
func ReceiveChainCodeUploadFile(w http.ResponseWriter, r *http.Request) error {
	file, header, err := r.FormFile("file")
	filename := header.Filename
	workPath, _ := os.Getwd()
	fileSavePath := filepath.ToSlash(filepath.Join(workPath, "version", "fabric1.4.4", "chaincode", filename))
	if runtime.GOOS == "windows" {
		fileSavePath = strings.Replace(fileSavePath, "/", "\\", -1)
	}

	str := strings.Split(filename, ".tar")
	var ccName string
	if len(str) > 0 {
		ccName = str[0]
	} else {
		return errors.New("chaincode Name error:please check chaincode Name")
	}

	//ccPath := filepath.ToSlash(filepath.Join(workPath, "version", "fabric1.4.4", "chaincode", ccName))
	//检测是否重复上传
	errV := CheckAddChaincodeVersionISExist(ccName, "1.0")
	if errV != nil {
		//web.OutputEnter(w, "", nil, errors.WithMessage(errV, "chaincode version error:please check chaincode version"))
		return errors.WithMessage(errV, "chaincode name already exist: please check chaincode name")
	}

	out, err := os.Create(fileSavePath)
	if err != nil {
		return errors.WithMessage(err, "create save file "+filename+"err info:")
	}
	_, err = io.Copy(out, file)
	if err != nil {
		return errors.WithMessage(err, "copy save file "+filename+"err info:")
	}
	defer out.Close()
	var cmd *exec.Cmd
	ccPath := filepath.ToSlash(filepath.Join(workPath, "version", "fabric1.4.4", "chaincode"))
	execCommand := fmt.Sprintf("cd %s && tar -xvf %s", ccPath, fileSavePath)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", execCommand)
	} else {
		cmd = exec.Command("/bin/bash", "-c", execCommand)
	}
	if err := cmd.Start(); err != nil {
		return errors.WithMessage(err, "exec command1 start")
	}
	if err := cmd.Wait(); err != nil {
		return errors.WithMessage(err, "exec command1 wait")
	}
	execCommand = fmt.Sprintf("rm -rf %s/%s", ccPath, filename)
	if runtime.GOOS != "windows" {
		cmd = exec.Command("/bin/bash", "-c", execCommand)
		if err := cmd.Start(); err != nil {
			return errors.WithMessage(err, "exec command2 wait")
		}
		if err := cmd.Wait(); err != nil {
			return errors.WithMessage(err, "exec command2 wait")
		}
	}

	//把新上传的合约写入数据库
	err = UploadChainCodeWriteDB(filename)
	if err != nil {
		return errors.WithMessage(err, "upload chaincode write db error")
	}
	return nil
}

//UploadChainCodeWriteDB 新合约上传信息写入数据库
func UploadChainCodeWriteDB(ccName string) error {
	ccVersion := "1.0"
	str := strings.Split(ccName, ".tar")
	ccName = str[0]
	sqlI := "insert into ommp_chaincode (task_id,cc_name,cc_version,is_install,status) values (?,?,?,?,?)"
	result, err := db.Exec(sqlI, "upload", ccName, ccVersion, 0, 0)
	if err != nil {
		return errors.Errorf("mysql insert ommp_chaincode fail %s", err)
	}
	rows, _ := result.RowsAffected()
	if rows <= 0 {
		return errors.New("mysql insert ommp_chaincode fail")
	}
	return nil
}

//ReceiveUpgradeChainCodeUploadFile 接受升级链码文件压缩包
func ReceiveUpgradeChainCodeUploadFile(w http.ResponseWriter, r *http.Request) error {
	file, header, err := r.FormFile("file")
	ccVersion := r.FormValue("version")
	if len(ccVersion) == 0 {
		return errors.New("upload file lack version")
	}
	filename := header.Filename
	workPath, _ := os.Getwd()
	fileSavePath := filepath.ToSlash(filepath.Join(workPath, "version", "fabric1.4.4", "chaincode", filename))
	if runtime.GOOS == "windows" {
		fileSavePath = strings.Replace(fileSavePath, "/", "\\", -1)
	}

	str := strings.Split(filename, ".tar")
	var ccName string
	if len(str) > 0 {
		ccName = str[0]
	} else {
		return errors.New("chaincode Name error:please check chaincode Name")
	}
	ccPath := filepath.ToSlash(filepath.Join(workPath, "version", "fabric1.4.4", "chaincode", ccName))
	//检测版本是否已经存在了
	errV := CheckUpgradeChaincodeVersionISExist(ccName, ccVersion)
	if errV != nil {
		//web.OutputEnter(w, "", nil, errors.WithMessage(errV, "chaincode version error:please check chaincode version"))
		return errors.WithMessage(errV, "chaincode version error:please check chaincode version")
	}
	var cmd *exec.Cmd
	//检测文件是否存在 存在删除 以防备没有替换成功
	_, err = os.Stat(ccPath)
	if err == nil {
		execCommand := fmt.Sprintf("rm -rf %s/%s", ccPath, filename)
		if runtime.GOOS != "windows" {
			cmd = exec.Command("/bin/bash", "-c", execCommand)
			if err := cmd.Start(); err != nil {
				//web.OutputEnter(w, "", nil, errors.WithMessage(err, "exec command2 wait"))
				return err
			}
			if err := cmd.Wait(); err != nil {
				//web.OutputEnter(w, "", nil, errors.WithMessage(err, "exec command2 wait"))
				return err
			}
		}
	}
	out, err := os.Create(fileSavePath)
	if err != nil {
		//web.OutputEnter(w, "", nil, errors.WithMessage(err, "create save file "+filename+"err info:"))
		return err
	}
	_, err = io.Copy(out, file)
	if err != nil {
		//web.OutputEnter(w, "", nil, errors.WithMessage(err, "copy save file "+filename+"err info:"))
		return err
	}
	defer out.Close()

	ccPath = filepath.ToSlash(filepath.Join(workPath, "version", "fabric1.4.4", "chaincode"))
	execCommand := fmt.Sprintf("cd %s && tar -xvf %s", ccPath, fileSavePath)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", execCommand)
	} else {
		cmd = exec.Command("/bin/bash", "-c", execCommand)
	}
	if err := cmd.Start(); err != nil {
		//web.OutputEnter(w, "", nil, errors.WithMessage(err, "exec command1 start"))
		return err
	}
	if err := cmd.Wait(); err != nil {
		//web.OutputEnter(w, "", nil, errors.WithMessage(err, "exec command1 wait"))
		return err
	}
	execCommand = fmt.Sprintf("rm -rf %s/%s", ccPath, filename)
	if runtime.GOOS != "windows" {
		cmd = exec.Command("/bin/bash", "-c", execCommand)
		if err := cmd.Start(); err != nil {
			//web.OutputEnter(w, "", nil, errors.WithMessage(err, "exec command2 wait"))
			return err
		}
		if err := cmd.Wait(); err != nil {
			//web.OutputEnter(w, "", nil, errors.WithMessage(err, "exec command2 wait"))
			return err
		}
	}

	//把新上传的合约写入数据库
	err = UploadUpgradeChainCodeWriteDB(ccName, ccVersion)
	if err != nil {
		//web.OutputEnter(w, "", nil, errors.WithMessage(err, "upload chaincode write db error"))
		return err
	}
	return nil
}

//CheckUpgradeChaincodeVersionISExist 检测升级上传版本是否存在
func CheckUpgradeChaincodeVersionISExist(ccName, ccVersion string) error {

	sqlS := "select cc_version from ommp_chaincode where cc_name=\"" + ccName + "\""
	row, err := db.Query(sqlS)
	defer row.Close()
	if err != nil {
		return errors.WithMessage(err, "mysql query ommp_chaincode version info fail")
	}
	ccVersionMap, err := GetRowsValues(row)
	if err != nil {
		return errors.WithMessage(err, "mysql query ommp_chaincode version result rows to map fail")
	}
	var ccVersionList []string
	for _, ccV := range ccVersionMap {
		ccVersion := ccV["cc_version"]
		ccVersionList = append(ccVersionList, ccVersion)
	}
	for _, ccv := range ccVersionList {
		if ccv >= ccVersion {
			return errors.New("The chaincode version is lower than the latest version")
		}
	}
	return nil
}

//CheckAddChaincodeVersionISExist 检测新增上传版本是否存在
func CheckAddChaincodeVersionISExist(ccName, ccVersion string) error {
	sqlS := "select count(*) from ommp_chaincode where cc_name=\"" + ccName + "\" and cc_version=\"" + ccVersion + "\" and is_install=0 and status=0"
	fmt.Println("sqls info ", ccName, ccVersion)
	var count int
	err := db.QueryRow(sqlS).Scan(&count)
	if err != nil {
		return errors.Errorf("mysql check chaincode fail %s", err)
	}
	if count > 0 {
		return errors.New("mysql check chaincode name already exist")
	}
	return nil
}

//UploadUpgradeChainCodeWriteDB 用来升级的合约信息写入数据库
func UploadUpgradeChainCodeWriteDB(ccName, ccVersion string) error {

	sqlI := "insert into ommp_chaincode (task_id,cc_name,cc_version,is_install,status) values (?,?,?,?,?)"
	result, err := db.Exec(sqlI, "upgred", ccName, ccVersion, 0, 0)
	if err != nil {
		return errors.Errorf("mysql insert upgrade chaincode fail %s", err)
	}
	rows, _ := result.RowsAffected()
	if rows <= 0 {
		return errors.New("mmysql insert upgrade chaincode fail")
	}
	return nil
}

//UpdateSuccessDisableChainCodeIndent 更新合约数据库信息
func UpdateSuccessDisableChainCodeIndent(indent *objectdefine.Indent) error {
	for ccName, ccA := range indent.Chaincode {
		for _, cc := range ccA {
			endorsejson, _ := json.Marshal(cc.EndorsementOrg)
			endorseString := string(endorsejson)
			//sqlI := "insert into chaincode (cc_name,cc_version,cc_policy,cc_org,channelname,detail) values (?,?,?,?,?,?)"
			sqlU := "update ommp_chaincode set status=? where cc_name=? and cc_version=? and channelname=? and cc_org=?"
			result, err := db.Exec(sqlU, 0, ccName, cc.Version, indent.ChannelName, endorseString)
			if err != nil {
				return errors.Errorf("mysql update chaincode fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update chaincodet fail")
			}
		}
	}

	return nil
}

//UpdateSuccessEnableChainCodeIndent 启用成功更新合约数据库状态
func UpdateSuccessEnableChainCodeIndent(indent *objectdefine.Indent) error {

	for ccName, ccA := range indent.Chaincode {
		for _, cc := range ccA {
			endorsejson, _ := json.Marshal(cc.EndorsementOrg)
			endorseString := string(endorsejson)
			//sqlI := "insert into chaincode (cc_name,cc_version,cc_policy,cc_org,channelname,detail) values (?,?,?,?,?,?)"
			sqlU := "update ommp_chaincode set status=? where cc_name=? and cc_version=? and channelname=? and cc_org=?"
			result, err := db.Exec(sqlU, 1, ccName, cc.Version, indent.ChannelName, endorseString)
			if err != nil {
				return errors.Errorf("mysql update chaincode fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update chaincodet fail")
			}
		}
	}

	return nil
}

//CheckChannelIsExist 检测通道是否已经存在
func CheckChannelIsExist(channelName string) (bool, error) {

	var count int
	sql := "select count(*) from  ommp_indent where channelname = \"" + channelName + "\""
	err := db.QueryRow(sql).Scan(&count)
	if err != nil {
		return false, errors.Errorf("mysql select form ommp_indent fail %s", err)
	}
	if count > 0 {
		return false, errors.Errorf("check channel=%s  is exist", channelName)
	}
	return true, nil
}

//GetCCInstantiatedTime 查询合约是否存在过安装实例化
func GetCCInstantiatedTime(channelName, ccName, ccVersion string) (int, error) {
	sqlU := "select * from ommp_chaincode where cc_name=? and cc_version=? and channelname=? and status=?"
	result, err := db.Exec(sqlU, ccName, ccVersion, channelName, 4)
	if err != nil {
		return -1, errors.Errorf("mysql update chaincode fail %s", err)
	}
	rows, _ := result.RowsAffected()
	if rows <= 0 {
		return 0, nil
	}
	return 1, nil
}

//CheckPeerPort 检测并补全节点端口
func CheckPeerPort(ip string) (int, int, int, error) {
	sql := "select distinct peer_port,couchdb_port from  ommp_indent where peer_ip = \"" + ip + "\""
	row, err := db.Query(sql)
	defer row.Close()
	if err != nil {
		return 0, 0, 0, errors.WithMessage(err, "mysql query ommp_indent fail")
	}
	portArrayMap, err := GetRowsValues(row)
	if err != nil {
		return 0, 0, 0, errors.WithMessage(err, "mysql query ommp_indent result rows to map fail")
	}
	allPort := make(map[string]bool, 0)
	for _, port := range portArrayMap {
		peerPort := port["peer_port"]
		allPort[peerPort] = true
	}
	peerPort := 7051
	couchdbPort := 5084
	ccPort := 7052
LOOP:
	for {
		peerPort += 100
		couchdbPort += 100
		ccPort += 100
		pS := strconv.Itoa(peerPort)
		if _, ok := allPort[pS]; ok {
			goto LOOP
		}
		goto END
	}
END:

	return peerPort, couchdbPort, ccPort, nil
}

//CheckCAPort 检测并补全CA端口
func CheckCAPort(ip string) (int, error) {
	sqlC := "select distinct ca_port from  ommp_indent where peer_ip = \"" + ip + "\""
	rowC, err := db.Query(sqlC)
	defer rowC.Close()
	if err != nil {
		return 0, errors.WithMessage(err, "mysql query ommp_indent fail")
	}
	caArrayMap, err := GetRowsValues(rowC)
	if err != nil {
		return 0, errors.WithMessage(err, "mysql query ommp_indent result rows to map fail")
	}
	caPortMap := make(map[string]bool, 0)
	for _, ca := range caArrayMap {
		caPort := ca["ca_port"]
		caPortMap[caPort] = true
	}
	caPort := 7054
CLOOP:
	for {
		caPort += 1000
		cS := strconv.Itoa(caPort)
		if _, ok := caPortMap[cS]; ok {
			goto CLOOP
		}
		goto CEND
	}
CEND:

	return caPort, nil
}

//UpdateIndetPeerNickName 更新节点昵称
func UpdateIndetPeerNickName(indent *objectdefine.Indent) error {
	channelName := indent.ChannelName
	for _, org := range indent.Org {
		for _, peer := range org.Peer {
			if len(peer.IP) == 0 || len(peer.Name) == 0 || len(peer.Domain) == 0 {
				return errors.New("update ommp_indent peer nickname fail: incoming parameter error")
			}
			sqlU := "update ommp_indent set nick_name=? where channelname=? and peer_ip=? and peer_name=? and peer_domain=?"
			result, err := db.Exec(sqlU, peer.NickName, channelName, peer.IP, peer.Name, peer.Domain)
			if err != nil {
				return errors.Errorf("mysql update ommp_indent peer nickname fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update ommp_indent peer nickname fail fail")
			}
		}
	}

	return nil
}

/*
*  新增通道 组织 节点 合约 成功之后更新各自的状态
*  0 表示-正在创建   1-表示运行中  2- 表示失败
 */

//updateSuccessCreateChannelTaskStatus 成功时:更新数据库创建通道数据库状态信息
func updateSuccessCreateChannelTaskStatus(indent *objectdefine.Indent) error {
	channelName := indent.ChannelName
	for _, org := range indent.Org {
		for _, peer := range org.Peer {
			peerIP := peer.IP
			peerName := peer.Name
			peerDomain := peer.Domain

			sqlU := "update ommp_indent set channel_status=?,org_status=?,peer_status=?,peer_run_status=? where peer_ip=? and peer_name=? and peer_domain=? and channelname=?"
			result, err := db.Exec(sqlU, 1, 1, 1, 1, peerIP, peerName, peerDomain, channelName)
			if err != nil {
				return errors.Errorf("mysql update ommp_indent fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update ommp_indent fail")
			}
		}
	}

	return nil
}

//updateSuccessCreateOrgTaskStatus 成功时:更新数据库创建组织数据库状态信息
func updateSuccessCreateOrgTaskStatus(indent *objectdefine.Indent) error {
	channelName := indent.ChannelName
	for _, org := range indent.Org {
		for _, peer := range org.Peer {
			peerIP := peer.IP
			peerName := peer.Name
			peerDomain := peer.Domain

			sqlU := "update ommp_indent set channel_status=?,org_status=?,peer_status=?,peer_run_status=? where peer_ip=? and peer_name=? and peer_domain=? and channelname=?"
			result, err := db.Exec(sqlU, 1, 1, 1, 1, peerIP, peerName, peerDomain, channelName)
			if err != nil {
				return errors.Errorf("mysql update ommp_indent fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update ommp_indent fail")
			}
		}
	}

	return nil
}

//updateSuccessDeleteOrgTaskStatus 成功时:更新数据库删除组织数据库状态信息
func updateSuccessDeleteOrgTaskStatus(indent *objectdefine.Indent) error {

	channelName := indent.ChannelName
	for _, org := range indent.Org {
		orgName := org.Name

		sqlU := "update ommp_indent set org_status=?,peer_status=?,peer_run_status where org_name=? and channelname=?"
		result, err := db.Exec(sqlU, 4, 4, 0, orgName, channelName)
		if err != nil {
			return errors.Errorf("mysql update ommp_indent fail %s", err)
		}
		rows, _ := result.RowsAffected()
		if rows <= 0 {
			return errors.New("mysql update ommp_indent fail")
		}

	}

	return nil
}

//updateSuccessCreatePeerTaskStatus 成功时:更新数据库创建节点数据库状态信息
func updateSuccessCreatePeerTaskStatus(indent *objectdefine.Indent) error {
	channelName := indent.ChannelName
	for _, org := range indent.Org {
		for _, peer := range org.Peer {
			peerIP := peer.IP
			peerName := peer.Name
			peerDomain := peer.Domain

			sqlU := "update ommp_indent set channel_status=?,org_status=?,peer_status=?,peer_run_status=? where peer_ip=? and peer_name=? and peer_domain=? and channelname=?"
			result, err := db.Exec(sqlU, 1, 1, 1, 1, peerIP, peerName, peerDomain, channelName)
			if err != nil {
				return errors.Errorf("mysql update ommp_indent fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update ommp_indent fail")
			}
		}
	}

	return nil
}

//updateSuccessDeletePeerTaskStatus 成功时:更新数据库删除节点数据库状态信息
func updateSuccessDeletePeerTaskStatus(indent *objectdefine.Indent) error {

	channelName := indent.ChannelName
	for _, org := range indent.Org {
		orgName := org.Name
		for _, peer := range org.Peer {

			sqlU := "update ommp_indent set peer_status=?,peer_run_status where org_name=? and channelname=? and peer_domain=?"
			result, err := db.Exec(sqlU, 4, 0, orgName, channelName, peer.Domain)
			if err != nil {
				return errors.Errorf("mysql update ommp_indent fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update ommp_indent fail")
			}
		}

	}

	return nil
}

//updateSuccessCreateChaincodeTaskStatus 成功时:更新数据库创建合约数据库状态信息
func updateSuccessCreateChaincodeTaskStatus(indent *objectdefine.Indent) error {
	for ccName, ccA := range indent.Chaincode {
		for _, cc := range ccA {
			ccVersion := cc.Version
			ccPolicy := cc.Policy
			ccChannel := indent.ChannelName
			var ccDesc string
			if len(cc.Describe) != 0 {
				ccDesc = cc.Describe
			}
			endorsejson, _ := json.Marshal(cc.EndorsementOrg)
			endorseString := string(endorsejson)
			//sqlI := "insert into chaincode (cc_name,cc_version,cc_policy,cc_org,channelname,detail) values (?,?,?,?,?,?)"
			sqlU := "update ommp_chaincode set cc_policy=?, cc_org=?,channelname=?,detail=?,is_install=?,status=? where cc_name=? and cc_version=? and is_install=? and status=?"
			result, err := db.Exec(sqlU, ccPolicy, endorseString, ccChannel, ccDesc, 1, 1, ccName, ccVersion, 0, 0)
			if err != nil {
				return errors.Errorf("mysql update ommp_chaincode fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update ommp_chaincodet fail")
			}
		}
	}

	return nil
}

//updateSuccessDeleteChaincodeTaskStatus 成功时:更新数据库删除合约数据库状态信息
func updateSuccessDeleteChaincodeTaskStatus(indent *objectdefine.Indent) error {
	for ccName, ccA := range indent.Chaincode {
		for _, cc := range ccA {
			endorsejson, _ := json.Marshal(cc.EndorsementOrg)
			endorseString := string(endorsejson)
			//sqlI := "insert into chaincode (cc_name,cc_version,cc_policy,cc_org,channelname,detail) values (?,?,?,?,?,?)"
			sqlU := "update ommp_chaincode set is_install=?,status=? where cc_name=? and cc_version=? and channelname=? and cc_org=?"
			result, err := db.Exec(sqlU, 0, 4, ccName, cc.Version, indent.ChannelName, endorseString)
			if err != nil {
				return errors.Errorf("mysql update ommp_chaincode fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update ommp_chaincodet fail")
			}
		}
	}

	return nil
}

//updateSuccessUpgradeChaincodeTaskStatus 成功时:更新数据库升级合约数据库状态信息
func updateSuccessUpgradeChaincodeTaskStatus(indent *objectdefine.Indent) error {
	for ccName, ccA := range indent.Chaincode {
		for _, cc := range ccA {
			ccVersion := cc.Version
			ccPolicy := cc.Policy
			ccChannel := indent.ChannelName
			var ccDesc string
			if len(cc.Describe) != 0 {
				ccDesc = cc.Describe
			}
			endorsejson, _ := json.Marshal(cc.EndorsementOrg)
			endorseString := string(endorsejson)
			//sqlI := "insert into chaincode (cc_name,cc_version,cc_policy,cc_org,channelname,detail) values (?,?,?,?,?,?)"
			sqlU := "update ommp_chaincode set cc_policy=?, cc_org=?,cc_version=?,detail=?,is_install=?,status=? where cc_name=? and channelname=?"
			result, err := db.Exec(sqlU, ccPolicy, endorseString, ccVersion, ccDesc, 1, 1, ccName, ccChannel)
			if err != nil {
				return errors.Errorf("mysql update ommp_chaincode fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update ommp_chaincodet fail")
			}
		}
	}

	return nil
}

/*
*  新增通道   失败之后更新各自的状态
*  0 表示-正在创建   1-表示运行中  2- 表示失败
 */

//UpdateFailCreateChannelTaskStatus  失败时:更新数据库创建通道数据库状态信息
func UpdateFailCreateChannelTaskStatus(indent *objectdefine.Indent) error {
	channelName := indent.ChannelName
	for _, org := range indent.Org {
		for _, peer := range org.Peer {
			peerIP := peer.IP
			peerName := peer.Name
			peerDomain := peer.Domain

			sqlU := "update ommp_indent set channel_status=?,org_status=?,peer_status=?,peer_run_status=? where peer_ip=? and peer_name=? and peer_domain=? and channelname=?"
			result, err := db.Exec(sqlU, 2, 1, 1, 1, peerIP, peerName, peerDomain, channelName)
			if err != nil {
				return errors.Errorf("create channel fail : mysql update ommp_indent fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("create channel fail : mysql update ommp_indent fail")
			}
		}
	}

	return nil
}

//UpdateFailCreateOrgTaskStatus   失败时:更新数据库创建组织数据库状态信息
func UpdateFailCreateOrgTaskStatus(indent *objectdefine.Indent) error {
	channelName := indent.ChannelName
	for _, org := range indent.Org {
		for _, peer := range org.Peer {
			peerIP := peer.IP
			peerName := peer.Name
			peerDomain := peer.Domain

			sqlU := "update ommp_indent set channel_status=?,org_status=?,peer_status=?,peer_run_status=? where peer_ip=? and peer_name=? and peer_domain=? and channelname=?"
			result, err := db.Exec(sqlU, 1, 2, 2, 2, peerIP, peerName, peerDomain, channelName)
			if err != nil {
				return errors.Errorf("create org fail : mysql update ommp_indent fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("create org fail : mysql update ommp_indent fail")
			}
		}
	}

	return nil
}

//UpdateFailDeleteOrgTaskStatus   失败时:更新数据库删除组织数据库状态信息
func UpdateFailDeleteOrgTaskStatus(indent *objectdefine.Indent) error {
	channelName := indent.ChannelName
	for _, org := range indent.Org {
		orgName := org.Name

		sqlU := "update ommp_indent set org_status=? where org_name=? and channelname=?"
		result, err := db.Exec(sqlU, 2, orgName, channelName)
		if err != nil {
			return errors.Errorf("delete org fail : mysql update ommp_indent fail %s", err)
		}
		rows, _ := result.RowsAffected()
		if rows <= 0 {
			return errors.New("delete org fail : mysql update ommp_indent fail")
		}

	}

	return nil
}

//UpdateFailCreatePeerTaskStatus 失败时:更新数据库创建节点数据库状态信息
func UpdateFailCreatePeerTaskStatus(indent *objectdefine.Indent) error {
	channelName := indent.ChannelName
	for _, org := range indent.Org {
		for _, peer := range org.Peer {
			peerIP := peer.IP
			peerName := peer.Name
			peerDomain := peer.Domain

			sqlU := "update ommp_indent set channel_status=?,org_status=?,peer_status=?,peer_run_status=? where peer_ip=? and peer_name=? and peer_domain=? and channelname=?"
			result, err := db.Exec(sqlU, 1, 1, 2, 2, peerIP, peerName, peerDomain, channelName)
			if err != nil {
				return errors.Errorf("mysql update ommp_indent fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update ommp_indent fail")
			}
		}
	}

	return nil
}

//UpdateFailDeletePeerTaskStatus   失败时:更新数据库删除节点数据库状态信息
func UpdateFailDeletePeerTaskStatus(indent *objectdefine.Indent) error {
	channelName := indent.ChannelName
	for _, org := range indent.Org {
		orgName := org.Name
		for _, peer := range org.Peer {
			sqlU := "update ommp_indent set peer_status=?,peer_run_status=? where org_name=? and channelname=? and peer_domain=?"
			result, err := db.Exec(sqlU, 2, 0, orgName, channelName, peer.Domain)
			if err != nil {
				return errors.Errorf("delete peer fail : mysql update ommp_indent fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("delete peer fail : mysql update ommp_indent fail")
			}
		}

	}

	return nil
}

//UpdateFailCreateChaincodeTaskStatus 失败时:更新数据库创建合约数据库状态信息
func UpdateFailCreateChaincodeTaskStatus(indent *objectdefine.Indent) error {
	for ccName, ccA := range indent.Chaincode {
		for _, cc := range ccA {
			ccVersion := cc.Version
			ccPolicy := cc.Policy
			ccChannel := indent.ChannelName
			var ccDesc string
			if len(cc.Describe) != 0 {
				ccDesc = cc.Describe
			}
			endorsejson, _ := json.Marshal(cc.EndorsementOrg)
			endorseString := string(endorsejson)
			//sqlI := "insert into chaincode (cc_name,cc_version,cc_policy,cc_org,channelname,detail) values (?,?,?,?,?,?)"
			sqlU := "update ommp_chaincode set cc_policy=?, cc_org=?,channelname=?,detail=?,is_install=?,status=? where cc_name=? and cc_version=? and is_install=? and status=?"
			result, err := db.Exec(sqlU, ccPolicy, endorseString, ccChannel, ccDesc, 0, 2, ccName, ccVersion, 0, 0)
			if err != nil {
				return errors.Errorf("mysql update ommp_chaincode fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update ommp_chaincodet fail")
			}
		}
	}

	return nil
}

//UpdateFailDeleteChaincodeTaskStatus 失败时:更新数据库删除合约数据库状态信息
func UpdateFailDeleteChaincodeTaskStatus(indent *objectdefine.Indent) error {
	for ccName, ccA := range indent.Chaincode {
		for _, cc := range ccA {
			endorsejson, _ := json.Marshal(cc.EndorsementOrg)
			endorseString := string(endorsejson)
			//sqlI := "insert into chaincode (cc_name,cc_version,cc_policy,cc_org,channelname,detail) values (?,?,?,?,?,?)"
			sqlU := "update ommp_chaincode set is_install=?,status=? where cc_name=? and cc_version=? and channelname=? and cc_org=?"
			result, err := db.Exec(sqlU, 1, 2, ccName, cc.Version, indent.ChannelName, endorseString)
			if err != nil {
				return errors.Errorf("delete chaincode fail : mysql update ommp_chaincode fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("delete chaincode fail : mysql update ommp_chaincodet fail")
			}
		}
	}

	return nil
}

//UpdateFailUpgradeChaincodeTaskStatus 失败时:更新数据库升级合约数据库状态信息
func UpdateFailUpgradeChaincodeTaskStatus(indent *objectdefine.Indent) error {
	for ccName, ccA := range indent.Chaincode {
		for _, cc := range ccA {
			ccVersion := cc.Version
			ccPolicy := cc.Policy
			ccChannel := indent.ChannelName
			var ccDesc string
			if len(cc.Describe) != 0 {
				ccDesc = cc.Describe
			}
			endorsejson, _ := json.Marshal(cc.EndorsementOrg)
			endorseString := string(endorsejson)
			//sqlI := "insert into chaincode (cc_name,cc_version,cc_policy,cc_org,channelname,detail) values (?,?,?,?,?,?)"
			sqlU := "update ommp_chaincode set cc_policy=?,cc_version=?, cc_org=?,detail=?,is_install=?,status=? where cc_name=? and channelname=?"
			result, err := db.Exec(sqlU, ccPolicy, ccVersion, endorseString, ccDesc, 0, 2, ccName, ccChannel)
			if err != nil {
				return errors.Errorf("update chaincode fail : mysql update ommp_chaincode fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("update chaincode fail : mysql update ommp_chaincodet fail")
			}
		}
	}

	return nil
}

//UpdateFailDisableChaincodeTaskStatus 失败时:更新数据库停用合约数据库状态信息
func UpdateFailDisableChaincodeTaskStatus(indent *objectdefine.Indent) error {
	for ccName, ccA := range indent.Chaincode {
		for _, cc := range ccA {
			endorsejson, _ := json.Marshal(cc.EndorsementOrg)
			endorseString := string(endorsejson)
			//sqlI := "insert into chaincode (cc_name,cc_version,cc_policy,cc_org,channelname,detail) values (?,?,?,?,?,?)"
			sqlU := "update ommp_chaincode set status=? where cc_name=? and cc_version=? and channelname=? and cc_org=?"
			result, err := db.Exec(sqlU, 2, ccName, cc.Version, indent.ChannelName, endorseString)
			if err != nil {
				return errors.Errorf("disable chaincode fail : mysql update ommp_chaincode fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("disable chaincode fail : mysql update ommp_chaincodet fail")
			}
		}
	}

	return nil
}

//UpdateFailEnableChaincodeTaskStatus 失败时:启用合约更新合约数据库状态
func UpdateFailEnableChaincodeTaskStatus(indent *objectdefine.Indent) error {

	for ccName, ccA := range indent.Chaincode {
		for _, cc := range ccA {
			endorsejson, _ := json.Marshal(cc.EndorsementOrg)
			endorseString := string(endorsejson)
			//sqlI := "insert into chaincode (cc_name,cc_version,cc_policy,cc_org,channelname,detail) values (?,?,?,?,?,?)"
			sqlU := "update ommp_chaincode set status=? where cc_name=? and cc_version=? and channelname=? and cc_org=?"
			result, err := db.Exec(sqlU, 2, ccName, cc.Version, indent.ChannelName, endorseString)
			if err != nil {
				return errors.Errorf("mysql update chaincode fail %s", err)
			}
			rows, _ := result.RowsAffected()
			if rows <= 0 {
				return errors.New("mysql update chaincodet fail")
			}
		}
	}

	return nil
}



//MakeStepGeneralRemoteServerUpdateDB 主机服务插入信息
func MakeStepGeneralRemoteServerUpdateDB(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		taskID := general.ID
		sourceID := general.SourceID
		server_name := general.Server.ServerName
		server_des :=general.Server.ServerDes
		server_extip := general.Server.ServerExtIp
		server_intip  := general.Server.ServerIntIp
		server_user := general.Server.ServerUser 
		server_password := general.Server.ServerPassword
		server_num := 0
		server_status := 1
		sqlI := "insert into ommp_server (task_id,source_id,server_name,server_des,server_extip,server_intip,server_user,server_password,server_num,server_status) values (?,?,?,?,?,?,?,?,?,?)"
		result, err := db.Exec(sqlI, taskID, sourceID, server_name, server_des, server_extip, server_intip, server_user, server_password, server_num, server_status)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return errors.Errorf("channel info insert mysql ommp_server fail %s", err)
		}
		rows, _ := result.RowsAffected()
		if rows <= 0 {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return errors.New("channel info insert mysql ommp_server fail")
		}
		return nil
	}
}

//CheckServiceExist 检测主机服务是否已经存在
func CheckServiceExist(extip,intip string) (bool, error) {

	var count int
	sql := "select count(*) from  ommp_server where server_extip = \"" + extip + "\" and server_intip=\"" + intip + "\""
	err := db.QueryRow(sql).Scan(&count)
	if err != nil {
		return false, errors.Errorf("mysql select form ommp_server fail %s", err)
	}
	if count > 0 {
		return false, errors.Errorf("check extip=%s intip=%s is exist", extip,intip)
	}
	return true, nil
}

// GetServiceInfoFromExtIP 根据外部ip获取主机服务信息
func GetServiceInfoFromExtIP(extip string) (*objectdefine.IndentServer,error){
	sqlU := "select * from ommp_server where server_extip="+extip+""
	row, err := db.Query(sqlU)
	defer row.Close()
	if err != nil {
		return nil,errors.WithMessage(err, "mysql query ommp_server fail")
	}

	indentArrayMap, err := GetRowsValues(row)
	if err != nil {
		return nil,errors.WithMessage(err, "mysql query ommp_server result rows to map fail")
	}
    if len(indentArrayMap) >1 {
		return nil,errors.WithMessage(err, "mysql query ommp_server result fail")
	}
	serverInfo := &objectdefine.IndentServer{}
	for _, allInfo := range indentArrayMap {
		serverInfo.ServerName = allInfo["server_name"]
		serverInfo.ServerExtIp = allInfo["server_extip"]
		serverInfo.ServerIntIp = allInfo["server_intip"]
		serverInfo.ServerUser = allInfo["server_user"]
		serverInfo.ServerPassword = allInfo["server_password"]
		break
	}
	return serverInfo,nil
}
//
//UpdateServerName 更新z主机服务名称
func UpdateServerName(server *objectdefine.IndentServer) error {
	sqlU := "update ommp_server set server_name=? where server_extip=? and server_intip=?"
	result, err := db.Exec(sqlU, server.ServerName, server.ServerExtIp, server.ServerIntIp)
	if err != nil {
		return errors.Errorf("mysql update ommp_server server name fail %s", err)
	}
	rows, _ := result.RowsAffected()
	if rows <= 0 {
		return errors.New("mysql update ommp_server server name fail fail")
	}
	return nil
}

//UpdateServerDes 更新主机服务描述
func UpdateServerDes(server *objectdefine.IndentServer) error {
	sqlU := "update ommp_server set server_des=? where server_extip=? and server_intip=?"
	result, err := db.Exec(sqlU, server.ServerDes, server.ServerExtIp, server.ServerIntIp)
	if err != nil {
		return errors.Errorf("mysql update ommp_server server des fail %s", err)
	}
	rows, _ := result.RowsAffected()
	if rows <= 0 {
		return errors.New("mysql update ommp_server server des fail fail")
	}
	return nil
}

//UpdateServerUser 更新主机ssh连接用户名
func UpdateServerUser(server *objectdefine.IndentServer) error {
	sqlU := "update ommp_server set server_user=?, server_password=? where server_extip=? and server_intip=?"
	result, err := db.Exec(sqlU, server.ServerUser, server.ServerPassword,server.ServerExtIp, server.ServerIntIp)
	if err != nil {
		return errors.Errorf("mysql update ommp_server server ssh user password fail %s", err)
	}
	rows, _ := result.RowsAffected()
	if rows <= 0 {
		return errors.New("mysql update ommp_server server ssh user password fail fail")
	}
	return nil
}

//CheckServerUsingStatus 检测服务是否符合删除条件
func CheckServerUsingStatus(server *objectdefine.IndentServer) (bool,error) {
	sqlU := "select server_status from ommp_server where server_extip="+server.ServerExtIp+" and server_intip="+server.ServerIntIp+" and server_num=0"
	row, err := db.Query(sqlU)
	defer row.Close()
	if err != nil {
		return false,errors.WithMessage(err, "mysql query ommp_server fail")
	}

	indentArrayMap, err := GetRowsValues(row)
	if err != nil {
		return false,errors.WithMessage(err, "mysql query ommp_server result rows to map fail")
	}

	for _, allInfo := range indentArrayMap {
		status,_ :=strconv.Atoi(allInfo["server_status"])
		if status == 1 || status ==2 {
            return true,nil
		}
	}
	return false,nil
}

//DeleteServer 删除服务
func DeleteServer(server *objectdefine.IndentServer) error {
	sqlU := "delete from ommp_server where server_extip=? and server_intip=?"
	result, err := db.Exec(sqlU, server.ServerUser, server.ServerPassword,server.ServerExtIp, server.ServerIntIp)
	if err != nil {
		return errors.Errorf("mysql update ommp_server server ssh user password fail %s", err)
	}
	rows, _ := result.RowsAffected()
	if rows <= 0 {
		return errors.New("mysql update ommp_server server ssh user password fail fail")
	}
	return nil
}


//MakeStepGeneralCreateCompleteDeleteDB 
func MakeStepGeneralCreateCompleteDeleteDB(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		err := DeleteBlockChainServer()
		if err != nil{
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
		}
		return nil
	}
}
//DeleteBlockServer 删除区块链网络服务
func DeleteBlockChainServer() error {
	sqlD := "delete from ommp_indent where id >=0"
	result, err := db.Exec(sqlD)
	if err != nil {
		return errors.Errorf("mysql delete ommp_indent server fail %s", err)
	}
	rows, _ := result.RowsAffected()
	if rows <= 0 {
		return errors.New("mysql delete ommp_indent server fail")
	}

	sqlD = "delete from ommp_orderer where id >=0"
	result, err = db.Exec(sqlD)
	if err != nil {
		return errors.Errorf("mysql delete ommp_orderer server fail %s", err)
	}
	rows, _ = result.RowsAffected()
	if rows <= 0 {
		return errors.New("mysql delete ommp_orderer server fail")
	}

	sqlD = "delete from ommp_chaincode where id >=0"
	result, err = db.Exec(sqlD)
	if err != nil {
		return errors.Errorf("mysql delete ommp_chaincode server fail %s", err)
	}
	rows, _ = result.RowsAffected()
	if rows <= 0 {
		return errors.New("mysql delete ommp_chaincode server fail")
	}


	// sqlD = "delete from ommp_server where id >=0"
	// result, err = db.Exec(sqlD)
	// if err != nil {
	// 	return errors.Errorf("mysql delete ommp_server server fail %s", err)
	// }
	// rows, _ = result.RowsAffected()
	// if rows <= 0 {
	// 	return errors.New("mysql delete ommp_server server fail")
	// }
	return nil
}