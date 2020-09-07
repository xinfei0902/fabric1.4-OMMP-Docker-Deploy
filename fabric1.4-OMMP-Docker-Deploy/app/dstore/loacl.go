package dstore

import (
	"deploy-server/app/objectdefine"
	"deploy-server/app/tools"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func makeErrorEmptyBuff(name string) error {
	return errors.Errorf("Empty %s buff in local config", name)
}

//NewLocalTypeFromJSON 把读取local.json里[]byte字节转LocalType{}结构 检测必需的文件和工具是否存在
func NewLocalTypeFromJSON(buff []byte) (ret *objectdefine.LocalType, err error) {
	ret = &objectdefine.LocalType{}
	err = json.Unmarshal(buff, ret)
	if err != nil {
		return
	}
	err = CheckLocalBase(ret)
	return
}

//CheckLocalBase 检测必要文件已经工具是否存在
func CheckLocalBase(localconfig *objectdefine.LocalType) (err error) {
	err = CheckConsensus(localconfig)
	if err != nil {
		return
	}

	err = CheckVersionFiles(localconfig)
	if err != nil {
		return
	}

	err = CheckOutput(localconfig)
	if err != nil {
		return
	}

	return nil
}

//CheckConsensus 检测共识
func CheckConsensus(localconfig *objectdefine.LocalType) (err error) {
	if localconfig.Consensus == nil || len(localconfig.Consensus.Consensus) == 0 {
		return makeErrorEmptyBuff("Consensus")
	}

	return nil
}

//CheckVersionFiles 检测版本下合约和工具是否存在
func CheckVersionFiles(localconfig *objectdefine.LocalType) (err error) {
	if len(localconfig.Versions) == 0 {
		return makeErrorEmptyBuff("Versions")
	}

	for _, one := range localconfig.Versions {
		if len(one.Version) == 0 {
			return makeErrorEmptyBuff("Version.Name")
		}

		if len(one.VersionRoot) == 0 {
			return makeErrorEmptyBuff("Version.Root")
		}

		info, err := os.Stat(one.VersionRoot)
		if err != nil {
			return errors.WithMessage(err, "Check Version Root error for "+one.Version)
		}

		if info.IsDir() != true {
			return errors.New("Version Root is not folder for " + one.Version)
		}

		for name, pairs := range one.ChainCode {
			if len(pairs.Version) == 0 {
				return makeErrorEmptyBuff("Version.Chaincode.Version")
			}

			two, ok := one.Build[name]
			if !ok {
				return errors.New("Chaincode Path is empty for" + one.Version + ".ChainCode." + name)
			}
			if two.Type != objectdefine.BuildTypeChainCode {
				return errors.New("Chaincode Type is wrong for" + one.Version + ".ChainCode." + name)
			}
		}

		for name, value := range one.Build {
			path := GetBuildPath(one.VersionRoot, value.BinName)
			info, err := os.Stat(path)
			if err != nil {
				return errors.WithMessage(err, "Check Version.Build BinName error for "+one.Version+".Build."+name)
			}
			if info.IsDir() == false && info.Size() == 0 {
				return errors.New("Chaincode bin file is empty for" + one.Version + ".Build." + name)
			}
		}
	}

	return nil
}

//CheckOutput  检测路径
func CheckOutput(localconfig *objectdefine.LocalType) (err error) {
	if len(localconfig.StoreRoot) == 0 {
		return makeErrorEmptyBuff("StorePath")
	}
	path := tools.STDPath(localconfig.StoreRoot)
	info, err := os.Stat(path)
	if err != nil {
		err = errors.WithMessage(err, "Check Store folder error")
		return
	}
	if info.IsDir() == false {
		err = errors.New("Store folder is not dir " + localconfig.StoreRoot)
		return
	}
	localconfig.StoreRoot = path
	return nil
}

//GetBuildPath 获取路径
func GetBuildPath(root string, name string) string {
	return filepath.Join(root, "build", name)
}

//LocalObejctIntoJSONBuff 结构转[]byte
func LocalObejctIntoJSONBuff(one *objectdefine.LocalType) ([]byte, error) {
	return json.Marshal(one)
}
