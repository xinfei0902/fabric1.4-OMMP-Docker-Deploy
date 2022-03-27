package dcache

import (
	"bytes"
	"deploy-server/app/dstore"
	"deploy-server/app/objectdefine"
	"deploy-server/app/tools"
	"io/ioutil"

	"github.com/pkg/errors"
)

var globalTemplateBuffCache map[string][]byte

//InitGlobalTemplateBuffCache 缓存全部模板脚本
func InitGlobalTemplateBuffCache() error {
	globalTemplateBuffCache = make(map[string][]byte)

	for _, core := range globalLocalObject.local.Versions {
		for _, pair := range [][]string{
			//deploy
			{objectdefine.BuildTypeDeploy, "base.yaml", "." + objectdefine.BuildTypeDeploy + ".base"},
			{objectdefine.BuildTypeDeploy, "script.sh", "." + objectdefine.BuildTypeDeploy + ".script"},
			{objectdefine.BuildTypeDeploy, "deploy.sh", "." + objectdefine.BuildTypeDeploy + ".deploy"},
			{objectdefine.BuildTypeDeploy,"deployEnv.sh","."+ objectdefine.BuildTypeDeploy + ".deployEnv"},
			//组织
			{objectdefine.BuildTypeOrg, "configtx.yaml", "." + objectdefine.BuildTypeOrg + ".configtx"},
			{objectdefine.BuildTypeOrg, "crypto-config.yaml", "." + objectdefine.BuildTypeOrg + ".crypto-config"},
			{objectdefine.BuildTypeOrg, "base.yaml", "." + objectdefine.BuildTypeOrg + ".base"},
			{objectdefine.BuildTypeOrg, "addOrgStep1.sh", "." + objectdefine.BuildTypeOrg + ".addOrgStep1"},
			{objectdefine.BuildTypeOrg, "addOrgStep2.sh", "." + objectdefine.BuildTypeOrg + ".addOrgStep2"},
			{objectdefine.BuildTypeOrg, "addOrgStep3.sh", "." + objectdefine.BuildTypeOrg + ".addOrgStep3"},
			{objectdefine.BuildTypeOrg, "addOrgStep4.sh", "." + objectdefine.BuildTypeOrg + ".addOrgStep4"},
			{objectdefine.BuildTypeOrg, "deleteOrgStep1.sh", "." + objectdefine.BuildTypeOrg + ".deleteOrgStep1"},
			{objectdefine.BuildTypeOrg, "deleteOrgStep2.sh", "." + objectdefine.BuildTypeOrg + ".deleteOrgStep2"},
			{objectdefine.BuildTypeOrg, "deleteOrgStep3.sh", "." + objectdefine.BuildTypeOrg + ".deleteOrgStep3"},
			{objectdefine.BuildTypeOrg, "deletePeerStatusMonitor.sh", "." + objectdefine.BuildTypeOrg + ".deletepeerStatusMonitor"},

			//节点
			{objectdefine.BuildTypePeer, "addPeerStep1.sh", "." + objectdefine.BuildTypePeer + ".addPeerStep1"},
			{objectdefine.BuildTypePeer, "addPeerStep2.sh", "." + objectdefine.BuildTypePeer + ".addPeerStep2"},
			{objectdefine.BuildTypePeer, "addPeerStep3.sh", "." + objectdefine.BuildTypePeer + ".addPeerStep3"},
			{objectdefine.BuildTypePeer, "base.yaml", "." + objectdefine.BuildTypePeer + ".base"},
			{objectdefine.BuildTypePeer, "peerStatusMonitor.sh", "." + objectdefine.BuildTypePeer + ".peerStatusMonitor"},
			{objectdefine.BuildTypePeer, "deletePeerStep1.sh", "." + objectdefine.BuildTypePeer + ".deletePeerStep1"},
			{objectdefine.BuildTypePeer, "deletePeerStatusMonitor.sh", "." + objectdefine.BuildTypePeer + ".deletepeerStatusMonitor"},

			//合约
			{objectdefine.BuildTypeChainCode, "addChainCodeStep1.sh", "." + objectdefine.BuildTypeChainCode + ".addChainCodeStep1"},
			{objectdefine.BuildTypeChainCode, "addChainCodeStep2.sh", "." + objectdefine.BuildTypeChainCode + ".addChainCodeStep2"},
			{objectdefine.BuildTypeChainCode, "deleteChainCodeScript.sh", "." + objectdefine.BuildTypeChainCode + ".deleteChainCodeScript"},
			{objectdefine.BuildTypeChainCode, "upgradeChainCodeStep1.sh", "." + objectdefine.BuildTypeChainCode + ".upgradeChainCodeStep1"},
			{objectdefine.BuildTypeChainCode, "upgradeChainCodeStep2.sh", "." + objectdefine.BuildTypeChainCode + ".upgradeChainCodeStep2"},
			{objectdefine.BuildTypeChainCode, "disableChainCodeScript.sh", "." + objectdefine.BuildTypeChainCode + ".disableChainCodeScript"},
			{objectdefine.BuildTypeChainCode, "enableChainCodeScript.sh", "." + objectdefine.BuildTypeChainCode + ".enableChainCodeScript"},
			//通道
			{objectdefine.BuildTypeChannel, "createChannelStep1.sh", "." + objectdefine.BuildTypeChannel + ".createChannelStep1"},
			{objectdefine.BuildTypeChannel, "createChannelStep2.sh", "." + objectdefine.BuildTypeChannel + ".createChannelStep2"},
			{objectdefine.BuildTypeChannel, "createNewOrgChannelStep1.sh", "." + objectdefine.BuildTypeChannel + ".createNewOrgChannelStep1"},
			{objectdefine.BuildTypeChannel, "createNewOrgChannelStep2.sh", "." + objectdefine.BuildTypeChannel + ".createNewOrgChannelStep2"},
			{objectdefine.BuildTypeChannel, "createNewOrgChannelStep3.sh", "." + objectdefine.BuildTypeChannel + ".createNewOrgChannelStep3"},
			// {objectdefine.TemplateSystemd, "simple.service", "." + objectdefine.TemplateSystemd + ".simple"},
		} {
			path := dstore.TemplateSubPath(core.VersionRoot, pair[0], pair[1])
			var buff []byte
			buff, err := ioutil.ReadFile(path)
			if err != nil {
				err = errors.WithMessage(err, "Load Template ["+path+"] error")
				return err
			}
			buff = bytes.Replace(buff, []byte("\r"), []byte{}, -1)
			globalTemplateBuffCache[core.Version+pair[2]] = buff
		}
	}
	return nil
}

//GetDeployBaseYamlTemplate 获取模板内容
func GetDeployBaseYamlTemplate(version string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeDeploy+".base"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}
	return buff, nil
}

//GetCompleteDeployScriptTemplate 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetCompleteDeployScriptTemplate(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeDeploy+".script"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetCompleteDeployScriptTemplate 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetCompleteDeployScriptDeployTemplate(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeDeploy+".deploy"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetCompleteDeployScriptTemplate 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetDeployEnvScriptTemplate(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeDeploy+".deployEnv"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetOrgConfigTemplate 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetOrgConfigTemplate(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeOrg+".configtx"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetOrgScriptStep1Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetOrgScriptStep1Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeOrg+".addOrgStep1"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetOrgScriptStep2Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetOrgScriptStep2Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeOrg+".addOrgStep2"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetOrgScriptStep3Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetOrgScriptStep3Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeOrg+".addOrgStep3"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetOrgScriptStep4Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetOrgScriptStep4Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeOrg+".addOrgStep4"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetDeleteOrgScriptStep1Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetDeleteOrgScriptStep1Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeOrg+".deleteOrgStep1"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetDeleteOrgScriptStep2Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetDeleteOrgScriptStep2Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeOrg+".deleteOrgStep2"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetDeleteOrgScriptStep3Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetDeleteOrgScriptStep3Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeOrg+".deleteOrgStep3"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetOrgBaseFileTemplate 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetOrgBaseFileTemplate(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeOrg+".base"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetPeerScriptStep1Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetPeerScriptStep1Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypePeer+".addPeerStep1"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetPeerScriptStep2Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetPeerScriptStep2Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypePeer+".addPeerStep2"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetPeerScriptStep3Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetPeerScriptStep3Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypePeer+".addPeerStep3"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetDeletePeerScriptStep1Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetDeletePeerScriptStep1Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypePeer+".deletePeerStep1"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetPeerBaseFileTemplate 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetPeerBaseFileTemplate(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypePeer+".base"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetChainCodeScriptStep1Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetChainCodeScriptStep1Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeChainCode+".addChainCodeStep1"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetChainCodeScriptStep2Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetChainCodeScriptStep2Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeChainCode+".addChainCodeStep2"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetDeleteChainCodeScriptStepTemplate 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetDeleteChainCodeScriptStepTemplate(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeChainCode+".deleteChainCodeScript"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetUpgradeChainCodeScriptStep1Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetUpgradeChainCodeScriptStep1Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeChainCode+".upgradeChainCodeStep1"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetUpgradeChainCodeScriptStep2Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetUpgradeChainCodeScriptStep2Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeChainCode+".upgradeChainCodeStep2"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetDisableChainCodeScriptStepTemplate 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetDisableChainCodeScriptStepTemplate(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeChainCode+".disableChainCodeScript"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetEnableChainCodeScriptStepTemplate 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetEnableChainCodeScriptStepTemplate(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeChainCode+".enableChainCodeScript"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetChannelScriptStep1Template 通道脚本模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetChannelScriptStep1Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeChannel+".createChannelStep1"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetChannelScriptStep2Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetChannelScriptStep2Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeChannel+".createChannelStep2"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetNewOrgChannelScriptStep1Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetNewOrgChannelScriptStep1Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeChannel+".createNewOrgChannelStep1"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetNewOrgChannelScriptStep2Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetNewOrgChannelScriptStep2Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeChannel+".createNewOrgChannelStep2"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetNewOrgChannelScriptStep3Template 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetNewOrgChannelScriptStep3Template(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypeChannel+".createNewOrgChannelStep3"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetPeerStatusFileTemplate 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetPeerStatusFileTemplate(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypePeer+".peerStatusMonitor"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}

//GetDeletePeerStatusFileTemplate 模板替换:把已经准备好参数 替换模板 生成新的[]byte字节返回
func GetDeletePeerStatusFileTemplate(version string, one map[string]string) ([]byte, error) {
	buff, ok := globalTemplateBuffCache[version+"."+objectdefine.BuildTypePeer+".deletepeerStatusMonitor"]
	if !ok {
		return nil, errors.New("No template for Version [" + version + "]")
	}

	buff = tools.ReplaceTemplateBuff(buff, one, nil, '$', '[', ']')
	return buff, nil
}
