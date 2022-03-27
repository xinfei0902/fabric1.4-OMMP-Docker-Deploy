package daction

import (
	"deploy-server/app/dmysql"
	"deploy-server/app/objectdefine"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

//MakeStepDeleteOrgLocalCife 删除本地生成的证书
func MakeStepDeleteOrgLocalCife(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec delete org cife start"},
		})
		err := GeneralDeleteOrgCife(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDeleteOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeleteOrgLocalCife",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec delete org cife start"},
		})
		return nil
	}
}

//GeneralDeleteOrgCife 以后多组织扩展
func GeneralDeleteOrgCife(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		err := GeneralDeleteOrgCifeFile(general, org, output)
		if err != nil {
			return err
		}
	}
	return nil
}

//GeneralDeleteOrgCifeFile 删除证书
func GeneralDeleteOrgCifeFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, output *objectdefine.TaskNode) error {
	cifePath := filepath.ToSlash(filepath.Join(general.SourceBaseOutput, "crypto-config", "peerOrganizations", orgOrder.OrgDomain))
	info, err := os.Stat(cifePath)
	if err != nil {
		return nil
	}
	if false == info.IsDir() {
		return errors.New("Is not folder")
	}

	return os.RemoveAll(cifePath)
}

//################删除节点证书############################

//MakeStepDeletePeerLocalCife 删除本地生成的证书
func MakeStepDeletePeerLocalCife(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec delete peer cife start"},
		})
		err := GeneralDeletePeerCife(general, output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailDeletePeerTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalDeletePeerLocalCife",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"exec delete peer cife start"},
		})
		return nil
	}
}

//GeneralDeletePeerCife 以后多组织扩展
func GeneralDeletePeerCife(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
	for _, org := range general.Org {
		for _, peer := range org.Peer {
			err := GeneralDeletePeerCifeFile(general, org, peer, output)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//GeneralDeletePeerCifeFile 删除证书(包含peers和users)
func GeneralDeletePeerCifeFile(general *objectdefine.Indent, orgOrder objectdefine.OrgType, peerOrder objectdefine.PeerType, output *objectdefine.TaskNode) error {
	peerCifePath := filepath.ToSlash(filepath.Join(general.SourceBaseOutput, "crypto-config", "peerOrganizations", orgOrder.OrgDomain, "peers", peerOrder.Domain))
	info, err := os.Stat(peerCifePath)
	if err != nil {
		return nil
	}
	if false == info.IsDir() {
		return errors.New("Is not folder")
	}
	os.RemoveAll(peerCifePath)
	userCifePath := filepath.ToSlash(filepath.Join(general.SourceBaseOutput, "crypto-config", "peerOrganizations", orgOrder.OrgDomain, "users", fmt.Sprintf("%s@%s", peerOrder.User, orgOrder.OrgDomain)))
	info, err = os.Stat(userCifePath)
	if err != nil {
		return nil
	}
	if false == info.IsDir() {
		return errors.New("Is not folder")
	}

	return os.RemoveAll(userCifePath)
}
