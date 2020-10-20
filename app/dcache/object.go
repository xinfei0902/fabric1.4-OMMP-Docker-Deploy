package dcache

import (
	"deploy-server/app/objectdefine"
	"encoding/json"
	"sync"

	"github.com/pkg/errors"
)

type localCacheV1 struct {
	local *objectdefine.LocalType

	ConsensusBuff []byte
	VersionBuff   []byte
}

var initOne sync.Once
var globalLocalObject *localCacheV1

//SaveStaticLocal 保存内容到globalLocalObject结构 也缓存所有脚本模板文件
func SaveStaticLocal(local *objectdefine.LocalType) (err error) {

	one := new(localCacheV1)
	one.local = local
	err = one.Init()
	if err != nil {
		return
	}
	err = errors.New("Double install local config")
	initOne.Do(func() {
		globalLocalObject = one
		err = InitGlobalTemplateBuffCache()
		if err != nil {
			return
		}

	})
	return
}

//Init 缓存local.json信息
func (core *localCacheV1) Init() error {

	/// Consensus
	shell := &objectdefine.BaseReponse{
		Success: true,
		Data:    core.local.Consensus,
	}
	buff, err := json.MarshalIndent(shell, "", "")
	if err != nil {
		err = errors.WithMessage(err, "prepare consensus error")
		return err
	}

	core.ConsensusBuff = buff

	temp := make([]objectdefine.VersionType, len(core.local.Versions))
	for i, one := range core.local.Versions {
		temp[i] = one
		temp[i].Build = nil
	}
	shell.Data = temp
	buff, err = json.MarshalIndent(shell, "", "")
	if err != nil {
		err = errors.WithMessage(err, "prepare consensus error")
		return err
	}

	core.VersionBuff = buff

	return nil
}
