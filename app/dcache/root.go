package dcache

import (
	"deploy-server/app/objectdefine"
	"path/filepath"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

//InitLocalConfig 初始化本地local.json 缓存local信息和脚本模板
func InitLocalConfig(one *objectdefine.LocalType) (err error) {
	err = SaveStaticLocal(one)

	if err != nil {
		return errors.WithMessage(err, "Init local config error")
	}

	return nil
}

//GetGeneralOrgFromBuff 初步构建任务id路径
func GetGeneralOrgFromBuff(buff []byte) (*objectdefine.Indent, error) {
	ret := &objectdefine.Indent{}
	var uJSON = jsoniter.ConfigCompatibleWithStandardLibrary
	err := uJSON.Unmarshal(buff, ret)
	if err != nil {
		return nil, errors.WithMessage(err, "Read body buff to Unmarshal error")
	}

	if len(ret.ID) == 0 {
		return nil, errors.New("Empty buff for ID")
	}
	if len(ret.SourceID) == 0 {
		// return nil, errors.New("Empty buff for SourceID")
		ret.SourceID = "sourceid"
	}
	if len(ret.ChannelName) == 0 {
		return nil, errors.New("Empty buff for channelname")
	}
	// if false == tools.IsAlpha(ret.ID) || len(ret.ID) > 128 {
	// 	return nil, errors.New("ID format error")
	// }
	// if false == tools.IsAlpha(ret.SourceID) || len(ret.SourceID) > 128 {
	// 	return nil, errors.New("ID format error")
	// }
	ret.ID = strings.ToLower(ret.ID)
	ret.SourceID = strings.ToLower(ret.SourceID)
	ret.BaseOutput = filepath.ToSlash(GetOutputSubPath(ret.ID, ""))
	ret.SourceBaseOutput = filepath.ToSlash(GetOutputSubPath(ret.SourceID, ""))
	return ret, nil
}
