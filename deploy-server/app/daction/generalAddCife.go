package daction

import (
	"deploy-server/app/dcache"
	"deploy-server/app/dmysql"
	"deploy-server/app/objectdefine"
	"deploy-server/app/tools"
	"fmt"
	"io"
	"os"
	"path/filepath"
	//"strings"
	"os/exec"
	"runtime"
	"github.com/pkg/errors"

	"github.com/hyperledger/fabric/common/tools/cryptogen/ca"
	"github.com/hyperledger/fabric/common/tools/cryptogen/csp"
	"github.com/hyperledger/fabric/common/tools/cryptogen/msp"
)

//NodeSpec 用于构建证书
type NodeSpec struct {
	Hostname           string   `yaml:"Hostname"`
	CommonName         string   `yaml:"CommonName"`
	Country            string   `yaml:"Country"`
	Province           string   `yaml:"Province"`
	Locality           string   `yaml:"Locality"`
	OrganizationalUnit string   `yaml:"OrganizationalUnit"`
	StreetAddress      string   `yaml:"StreetAddress"`
	PostalCode         string   `yaml:"PostalCode"`
	SANS               []string `yaml:"SANS"`
}

//DefaultCASpec 用于构建证书
func DefaultCASpec(cn string) *NodeSpec {
	return &NodeSpec{
		CommonName:         cn,
		Country:            "",
		Province:           "",
		Locality:           "",
		OrganizationalUnit: "",
		StreetAddress:      "",
		PostalCode:         "",
	}
}

//MakeStepGeneralCreateChanelCife //创建通道生成新的组织证书
func MakeStepGeneralCreateChanelCife(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {

		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create channel new org cife start"},
		})
		err := GeneralAddCertificate("org", general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			return err
		}
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create channel new org cife end"},
		})
		//在新的任务ID下保留一份证书
		dst := filepath.Join(general.BaseOutput, "crypto-config")
		src := filepath.Join(general.SourceBaseOutput, "crypto-config")
		tools.CopyFolder(dst, src)
		return nil
	}
}

//MakeStepGeneralAddPeerCife 创建peer证书
func MakeStepGeneralAddPeerCife(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {

		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create peer cife start"},
		})
		//err := GeneralAddCertificate("peer", general)
		err := GeneralCryptogenExternCertificate(general,output)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreatePeerTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreatePeerCife",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}

		//拷贝证书到sourceid
		dst := filepath.Join(general.BaseOutput, "crypto-config")
		src	:= filepath.Join(general.SourceBaseOutput, "crypto-config")
		tools.CopyFolder(dst, src)

		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create peer cife end"},
		})

		return nil
	}
}

//MakeStepGeneralAddOrgCife 创建组织证书
func MakeStepGeneralAddOrgCife(name string) objectdefine.RunTaskHandle {
	return func(general *objectdefine.Indent, output *objectdefine.TaskNode) error {
		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create org cife start"},
		})
		err := GeneralAddCertificate("org", general)
		if err != nil {
			output.AppendLog(&objectdefine.StepHistory{
				Name:  name,
				Error: []string{err.Error()},
			})
			result := err
			err := dmysql.UpdateFailCreateOrgTaskStatus(general)
			if err != nil {
				output.AppendLog(&objectdefine.StepHistory{
					Name:  "generalCreateOrgCife",
					Error: []string{err.Error()},
				})
				return err
			}
			return result
		}
		dst := filepath.Join(general.BaseOutput, "crypto-config")
		src := filepath.Join(general.SourceBaseOutput, "crypto-config")
		tools.CopyFolder(dst, src)

		output.AppendLog(&objectdefine.StepHistory{
			Name: name,
			Log:  []string{"create org cife end"},
		})

		return nil
	}
}

//GeneralAddCertificate 根据handle 执行不同创建证书过程
func GeneralAddCertificate(handle string, input *objectdefine.Indent) error {
	switch handle {
	case "peer":
		for _, org := range input.Org {
			for _, peer := range org.Peer {
				if len(peer.Domain) == 0 {
					return errors.New("lack domain")
				}
				if len(peer.Name) == 0 {
					return errors.New("lack peer name")
				}
				if len(peer.User) == 0 {
					return errors.New("lack UserName")
				}
				cife := fmt.Sprintf("crypto-config/peerOrganizations/%s", peer.OrgDomain)
				sourceOutput := filepath.ToSlash(dcache.GetOutputSubPath(input.SourceID, ""))
				CifeBuildPath := filepath.ToSlash(filepath.Join(sourceOutput, cife))
				err := BuildNewPeerInMSP(CifeBuildPath, peer.OrgDomain, peer.Domain, peer.User)
				if err != nil {
					return err
				}
			}
		}

	case "org":
		for _, org := range input.Org {
			for _, peer := range org.Peer {
				if len(input.BaseOutput) == 0 {
					return errors.New("lack certificate path")
				}
				if len(peer.Org) == 0 {
					return errors.New("lack OrgName")
				}
				if len(peer.Domain) == 0 {
					return errors.New("lack domain")
				}
				cife := fmt.Sprintf("crypto-config/peerOrganizations/%s", peer.OrgDomain)
				sourceOutput := filepath.ToSlash(dcache.GetOutputSubPath(input.SourceID, ""))
				CifeBuildPath := filepath.ToSlash(filepath.Join(sourceOutput, cife))
				if peer.User == "Admin" {
					err := BuildOneOrganization(CifeBuildPath, peer.Org, DefaultCASpec(peer.OrgDomain))
					if err != nil {
						return err
					}
					err = BuildNewPeerInMSP(CifeBuildPath, peer.OrgDomain, peer.Domain, "")
					if err != nil {
						return err
					}
				} else {
					err := BuildNewPeerInMSP(CifeBuildPath, peer.OrgDomain, peer.Domain, peer.User)
					if err != nil {
						return err
					}
				}
			}
		}
	case "nodes":
		sub := os.Args[2]
		switch sub {
		case "peer":
		case "orderer":
		case "sdk":
		}
	case "help":
	case "version":

	default:
		return errors.New("please select add cife is peer or org")
	}

	return nil
}

//BuildOneOrganization 创建组织下面目录以及证书
func BuildOneOrganization(orgDir, orgName string, orgSpec *NodeSpec) error {
	caDir := filepath.Join(orgDir, "ca")
	tlsCADir := filepath.Join(orgDir, "tlsca")
	mspDir := filepath.Join(orgDir, "msp")
	usrDir := filepath.Join(orgDir, "users")
	adminCertsDir := filepath.Join(mspDir, "admincerts")

	// generate signing CA
	signCA, err := ca.NewCA(caDir, orgName, "ca."+orgSpec.CommonName, orgSpec.Country, orgSpec.Province, orgSpec.Locality, orgSpec.OrganizationalUnit, orgSpec.StreetAddress, orgSpec.PostalCode)
	if err != nil {
		err = errors.WithMessage(err, "Error generating signCA for org "+orgName)
		return err
	}

	// generate TLS CA
	tlsCA, err := ca.NewCA(tlsCADir, orgName, "tlsca."+orgSpec.CommonName, orgSpec.Country, orgSpec.Province, orgSpec.Locality, orgSpec.OrganizationalUnit, orgSpec.StreetAddress, orgSpec.PostalCode)
	if err != nil {
		err = errors.WithMessage(err, "Error generating tlsCA for org "+orgName)
		return err
	}

	err = msp.GenerateVerifyingMSP(mspDir, signCA, tlsCA, false)
	if err != nil {
		err = errors.WithMessage(err, "Error generating MSP for org "+orgName)
		return err
	}

	// admin
	err = msp.GenerateLocalMSP(filepath.Join(usrDir, "Admin@"+orgSpec.CommonName), "Admin@"+orgSpec.CommonName, []string{}, signCA, tlsCA, msp.CLIENT, false)
	if err != nil {
		err = errors.WithMessage(err, "Error generating Admin User for org "+orgSpec.CommonName)
		return err
	}

	// copy the admin cert to the org's MSP admincerts
	err = copyAdminCert(usrDir, adminCertsDir, "Admin@"+orgSpec.CommonName)
	if err != nil {
		err = errors.WithMessage(err, "Error copying Admin cert for org "+orgSpec.CommonName)
		return err
	}

	return nil
}

//BuildNewPeerInMSP 创建peers目录下节点证书
func BuildNewPeerInMSP(OrgPath, baseDomain, host, user string) (err error) {
	// OrgPath /---
	//           |- peers /- host.baseDomain: should not be exist
	//           |- users /- user@baseDomain: should not be exist
	//           |- ca : should be exist
	//           |- tlsca : should be exist

	// 1. check ca files
	signCA, tlsCA, err := LoadCAObjectFromFiles(OrgPath, baseDomain)
	if err != nil {
		return err
	}

	// 2. check peer files
	// 3. make peer files
	if len(host) > 0 {
		err = subBuildNewPeerInMSP(signCA, tlsCA, OrgPath, baseDomain, host)
		if err != nil {
			err = errors.WithMessage(err, "Build peer licenses files error")
			return err
		}
	}

	// 4. check user files
	// 5. make user files
	if len(user) > 0 {
		err = subBuildNewPeerUserInMSP(signCA, tlsCA, OrgPath, baseDomain, user)
		if err != nil {
			err = errors.WithMessage(err, "Build user licenses files error")
			return err
		}
	}

	return nil
}

func copyAdminCert(usersDir, adminCertsDir, adminUserName string) error {
	if _, err := os.Stat(filepath.Join(adminCertsDir,
		adminUserName+"-cert.pem")); err == nil {
		return nil
	}
	// delete the contents of admincerts
	err := os.RemoveAll(adminCertsDir)
	if err != nil {
		return err
	}
	// recreate the admincerts directory
	err = os.MkdirAll(adminCertsDir, 0777)
	if err != nil {
		return err
	}
	err = copyFile(filepath.Join(usersDir, adminUserName, "msp", "signcerts",
		adminUserName+"-cert.pem"), filepath.Join(adminCertsDir,
		adminUserName+"-cert.pem"))
	if err != nil {
		return err
	}
	return nil

}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}

//LoadCAObjectFromFiles 创建节点证书时需要知道组织根ca证书路径
func LoadCAObjectFromFiles(OrgPath string, baseDomain string) (signCA *ca.CA, tlsCA *ca.CA, err error) {
	signCA, err = GetCA(filepath.Join(OrgPath, "ca"), DefaultCASpec("ca."+baseDomain))
	if err != nil {
		err = errors.WithMessage(err, "load ca error")
		return
	}

	tlsCA, err = GetCA(filepath.Join(OrgPath, "tlsca"), DefaultCASpec("tlsca."+baseDomain))
	if err != nil {
		err = errors.WithMessage(err, "load tlsca error")
		return
	}
	return
}

//MakeCA ...
func MakeCA(caDir string, spec *NodeSpec, orgName string) (*ca.CA, error) {
	return ca.NewCA(caDir, orgName, spec.CommonName, spec.Country, spec.Province, spec.Locality, spec.OrganizationalUnit, spec.StreetAddress, spec.PostalCode)
}

//GetCA 获取CA
func GetCA(caDir string, spec *NodeSpec) (*ca.CA, error) {
	_, signer, err := csp.LoadPrivateKey(caDir)
	if err != nil {
		err = errors.WithMessage(err, "Load private key error")
		return nil, err
	}
	cert, err := ca.LoadCertificateECDSA(caDir)
	if err != nil {
		err = errors.WithMessage(err, "Load cert key error")
		return nil, err
	}

	return &ca.CA{
		Name:               spec.CommonName,
		Signer:             signer,
		SignCert:           cert,
		Country:            spec.Country,
		Province:           spec.Province,
		Locality:           spec.Locality,
		OrganizationalUnit: spec.OrganizationalUnit,
		StreetAddress:      spec.StreetAddress,
		PostalCode:         spec.PostalCode,
	}, nil
}

func subBuildNewPeerInMSP(signca, tlsca *ca.CA, OrgPath, baseDomain, host string) (err error) {
	// OrgPath /---
	//           |- peers /- host.baseDomain: should not be exist

	url := host
	path := filepath.Join(OrgPath, "peers", url)
	_, err = os.Stat(filepath.Join(path, "msp"))
	if err == nil {
		return errors.Errorf("Peer [%s] is exist", url)
	}

	err = msp.GenerateLocalMSP(path, url, []string{url}, signca, tlsca, msp.PEER, false)
	if err != nil {
		err = errors.WithMessage(err, "Build MSP error")
		return
	}

	err = copyAdminCert(filepath.Join(OrgPath, "users"), filepath.Join(path, "msp", "admincerts"), "Admin@"+baseDomain)
	if err != nil {
		err = errors.WithMessage(err, "Copy Admin user cert error")
		return
	}
	return nil
}

func subBuildNewPeerUserInMSP(signca, tlsca *ca.CA, OrgPath, baseDomain, user string) (err error) {
	uri := user + "@" + baseDomain
	path := filepath.Join(OrgPath, "users", uri)
	_, err = os.Stat(filepath.Join(path, "msp"))
	if err == nil {
		return errors.Errorf("User [%s] is exist", uri)
	}

	err = msp.GenerateLocalMSP(path, uri, []string{}, signca, tlsca, msp.CLIENT, false)
	if err != nil {
		err = errors.WithMessage(err, "build MSP error")
		return
	}

	return nil
}


//使用configtxgen extend 扩展证书
func GeneralCryptogenExternCertificate(general *objectdefine.Indent,output *objectdefine.TaskNode) error {
	//获取已有网络信息
	indent, err := dmysql.GetStartTaskBeforIndent(general)
	if err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerScriptStep",
			Error: []string{"get the latest indent error"},
		})
		return err
	}
	
	//获取本来已有订单要在那个组织节点个数
	var orgName string
	for _,org := range general.Org{
		orgName = org.Name
	}
	var SPeerNum int
	for _,org := range indent.Org{
		if org.Name == orgName {
			SPeerNum = len(org.Peer)
		}
	}
	//组装crypto-config.yaml结构配置
	ret := &objectdefine.CryptoForOrgConfig{}
	//ret.OrdererOrgs = make([]objectdefine.OrgSpec, 0, 1)
	ret.PeerOrgs = make([]objectdefine.OrgSpec, 0, len(general.Org))
	for _, one := range general.Org {
		if len(one.Peer) == 0 {
			continue
		}
	    addPeerNum := len(one.Peer)
		totalPeerNum := addPeerNum +SPeerNum
		two := objectdefine.OrgSpec{}
		two.Name = one.Name
		two.Domain = one.OrgDomain
		two.EnableNodeOUs = true
		two.Template = objectdefine.NodeTemplate{
			Count: totalPeerNum,
		}
		two.Users = objectdefine.UsersSpec{

			Count: totalPeerNum - 1,
		}
		ret.PeerOrgs = append(ret.PeerOrgs, two)
	}
	// two := objectdefine.OrgSpec{}
	// var ordererName string
	// var ordererDomain string
	// two.Specs = make([]objectdefine.NodeSpec, len(general.Orderer))
	// for i, orderer := range indent.Orderer {
	// 	if len(ordererName) == 0 {
	// 		ordererDomain = orderer.OrgDomain
	// 	}
	// 	//hostName := fmt.Sprintf("%s:%d", orderer.Domain, orderer.Port)
	// 	two.Specs[i].Hostname = strings.ToLower(orderer.Name)
	// }

	// two.Name = "Orderer"
	// two.Domain = ordererDomain
	// two.EnableNodeOUs = true
	// ordererCount := len(general.Orderer)
	// two.Users = objectdefine.UsersSpec{
	// 	Count: ordererCount,
	// }
	// ret.OrdererOrgs = append(ret.OrdererOrgs, two)
	//创建配置文件
	_,err = os.Stat(general.BaseOutput)
	if err != nil{
		err = os.MkdirAll(general.BaseOutput, 0777)
		if nil != err {
			err = errors.WithMessage(err, "Make Output folder error")
			output.AppendLog(&objectdefine.StepHistory{
				Name:  "generalCreatePeerScriptStep",
				Error: []string{"Make Output folder error"},
			})
			return err
		}
	}
	

	filename := filepath.Join(general.SourceBaseOutput, "crypto-config.yaml")
	//先删除再写入
	var env []string
	env = []string{fmt.Sprintf("FABRIC_CFG_PATH=%s", general.SourceBaseOutput)}
	os.Setenv("FABRIC_CFG_PATH",general.SourceBaseOutput)
	var cmd *exec.Cmd
	execCommand := fmt.Sprintf("cd %s && rm -rf  crypto-config.yaml && export FABRIC_CFG_PATH=%s",general.SourceBaseOutput,general.SourceBaseOutput)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", execCommand)
	} else {
		cmd = exec.Command("/bin/bash", "-c", execCommand)
	}
	if err := cmd.Start(); err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerScriptStep",
			Error: []string{"create peer cert: exec cmd commad error"},
		})
		return err
	}
	if err := cmd.Wait(); err != nil {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerScriptStep",
			Error: []string{"create peer cert: exec cmd commad error"},
		})
		return err
	}
	err = tools.WriteYamlFile(filename, ret)
	if nil != err {
		err = errors.WithMessage(err, "Write crypto for org config error")
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerScriptStep",
			Error: []string{"Write crypto for org config error"},
		})
		return err
	}
	//创建命令
	//outputPath := filepath.Join(general.BaseOutput,"channel-artifacts")

	readconfig := filepath.Join(general.SourceBaseOutput,"crypto-config.yaml")
	
	
	cmdRet := make([]*objectdefine.CommandList, 0)

	toolsPath := dcache.GetBinPathByVersion(general.Version)

	pair := dcache.GetBinToolsPair()

	bins, ok := pair[runtime.GOOS]
	if !ok {
		output.AppendLog(&objectdefine.StepHistory{
			Name:  "generalCreatePeerScriptStep",
			Error: []string{"cryptogen.bin : create cert tools not exist"},
		})
		return err
	}
	one := &objectdefine.CommandList{}
	one.OS = runtime.GOOS
	// one.Call = MakeGenerateExecCommand(env, general, outputPath, readconfig,
	// 	filepath.Join(toolsPath, bins.CryptoGen),
	// 	filepath.Join(toolsPath, bins.ConfigGen))
    one.Call = MakeGenerateCryptogenExecCommand(env, general,  readconfig,
		filepath.Join(toolsPath, bins.CryptoGen))
	cmdRet = append(cmdRet, one)
	logs, err := CallExecAllCommand(general, cmdRet)
	history := &objectdefine.StepHistory{
		Name: "generalCreateCompleteDeployBlock",
		Log:  logs,
	}
	if err != nil {
		history.Error = []string{err.Error()}
	}
	return nil
}