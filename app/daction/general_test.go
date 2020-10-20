package daction

import (
	"deploy-server/app/objectdefine"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_arrat(t *testing.T) {
	sPath := "D:/GoPath/src/deploy-server/output/sourceid/crypto-config/peerOrganizations/baiyun.example.com/tlsca/"
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
	orgmp["baiyun"] = objectdefine.OrgType{}
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
