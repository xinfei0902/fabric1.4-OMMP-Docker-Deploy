package dtask

import (
	"deploy-server/app/daction"
	"deploy-server/app/dmysql"
	"deploy-server/app/objectdefine"
	"deploy-server/app/dssh"
)

//NewGeneralCompleteDeployStep 一键部署 
func NewGeneralCompleteDeployStep() []objectdefine.SubTaskPair {
	return []objectdefine.SubTaskPair{
		//创建创世区块
		MakeTaskPair("generalCreateCompleteDeployBlock", daction.MakeStepGeneralCreateCompleteDeployBlock),
		//("generalCreateCompleteDeployBlock", daction.MakeStepGeneralCreateCompleteDeployBlockAnther),
		
		//创建 base.yaml文件
		MakeTaskPair("generalCreateCompleteDeployBaseYaml", daction.MakeStepGeneralCreateCompleteDeployBaseYaml),
		// //创建deploy.sh脚本
		// MakeTaskPair("generalCreateCompleteDeployScript", daction.MakeStepGeneralCreateChannelScript),
		//创建script.sh脚本
		MakeTaskPair("generalCreateCompleteDeployScript", daction.MakeStepGeneralCreateCompleteDeployScript),
		//打包
		MakeTaskPair("generalCreateCompleteDeployZip", daction.MakeStepGeneralCreateCompleteDeployZip),
		//上传
		MakeTaskPair("generalCreateCompleteDeployZipUpload", daction.MakeStepGeneralCreateCompleteDeployZipUpload),
		//解压
		MakeTaskPair("generalCreateCompleteDeployUnZip", daction.MakeStepGeneralCreateCompleteDeployUnzip),
		//启动容器
		MakeTaskPair("generalCreateCompleteDeployStartContainerStep", daction.MakeStepGeneralCreateCompleteDeployStartContainer),
		//执行创建通道 加入通道 安装合约脚本
		MakeTaskPair("generalCreateCompleteDeployExecScriptStep", daction.MakeStepGeneralCreateCompleteDeployExecScriptStep),
		//信息写入数据库
		MakeTaskPair("generalCreateCompleteDeployUpdateDB", dmysql.MakeStepGeneralUpdateCreateCompleteDeployIndent),
	}
}



//NewGeneralCompleteDeleteDeployStep 一键删除部署
func NewGeneralCompleteDeleteDeployStep() []objectdefine.SubTaskPair {
	return []objectdefine.SubTaskPair{
		//清除所有区块链机器容器 挂载  挂载目录 docker网络
		MakeTaskPair("generalCreateCompleteDeleteDeployRemote", daction.MakeStepGeneralCreateCompleteDeleteDeployRemote),
		//清除本地 任务包 
        MakeTaskPair("generalCreateCompleteDeleteLocal",daction.MakeStepGeneralCreateCompleteDeleteDeployLocal),
		//清除数据库
		MakeTaskPair("generalCreateCompleteDeployDeleteDB", dmysql.MakeStepGeneralCreateCompleteDeleteDB),
	}
}

//NewGeneralCreateChannelStep 创建通道所有执行步骤
func NewGeneralCreateChannelStep() []objectdefine.SubTaskPair {

	return []objectdefine.SubTaskPair{
		//1.初步写入数据库
		MakeTaskPair("generalCreateChannelWriteDB", dmysql.MakeStepGeneralInsertCreateChannelIndent),
		//2.创建configtx.yaml文件
		MakeTaskPair("generalCreateChannelTx", daction.MakeStepGeneralCreateBlock),
		//3.创建脚本
		MakeTaskPair("generalCreateChannelScript", daction.MakeStepGeneralCreateChannelScript),
		//4.打包
		MakeTaskPair("generalCreateChannelZip", daction.MakeStepGeneralChannelZip),
		//5.上传
		MakeTaskPair("generalCreateChannelZipUpload", daction.MakeStepGeneralChannelZipUpload),
		//6.解压
		MakeTaskPair("generalCreateChannelUnZip", daction.MakeStepGeneralChannelUnzip),
		//执行分两步
		//7.启动组织
		MakeTaskPair("generalCreateChannelExecScriptStep", daction.MakeStepGeneralChannelExecScriptStep1),
		// //2.启动创建命令
		// MakeTaskPair("generalCreateOrgExecScriptStep2", daction.MakeStepGeneralChannelExecScriptStep2),
		//8.更新通道状态数据库
		MakeTaskPair("generalCreateChannelUpdateDB", dmysql.MakeStepGeneralUpdateCreateChannelIndent),
	}
}

//NewGeneralAddOrgStep 新增组织具体操作
func NewGeneralAddOrgStep() []objectdefine.SubTaskPair {

	return []objectdefine.SubTaskPair{
		//初步写入数据库
		MakeTaskPair("generalCreateOrgWriteDB", dmysql.MakeStepGeneralInsertAddOrgIndent),
		//1.创建证书
		MakeTaskPair("generalCreateOrgCife", daction.MakeStepGeneralAddOrgCife),
		//2. 创建新增组织脚本
		MakeTaskPair("generalCreateOrgScriptStep", daction.MakeStepGeneralOrgScriptStep),
		//3.创建新的configtx.yaml文件
		MakeTaskPair("generalCreateOrgConfigtx", daction.MakeStepGeneralOrgConfigtx),
		//4. 创建base.yaml文件
		MakeTaskPair("generalCreateOrgBaseFile", daction.MakeStepGeneralOrgBaseFile),
		//5. 打包所有文件
		MakeTaskPair("generalCreateOrgZip", daction.MakeStepGeneralOrgZip),
		//6. 调用工具命令上传
		MakeTaskPair("generalCreateOrgZipUpload", daction.MakeStepGeneralOrgZipUpload),
		//7. 调用工具解压
		MakeTaskPair("generalCreateOrgUnZip", daction.MakeStepGeneralOrgUnzip),
		//8. 调用工具执行脚本先更新块配置
		MakeTaskPair("generalCreateOperateOrgExecScript", daction.MakeStepGeneralOperateOrgExecScript),
		//9. 调用工具执行脚本
		MakeTaskPair("generalCreateOrgExecScript", daction.MakeStepGeneralOrgExecScript),
		//10. 调用工具启动检测节点状态脚本
		MakeTaskPair("generalCreateOrgExecPeerStatusScript", daction.MakeStepGeneralOrgExecPeerStatusScript),
		//字段写入数据ku
		MakeTaskPair("generalCreateOrgUpdateDB", dmysql.MakeStepGeneralUpdateAddOrgIndent),
	}
}

//NewGeneralDeleteOrgStep 删除组织具体操作
func NewGeneralDeleteOrgStep() []objectdefine.SubTaskPair {

	return []objectdefine.SubTaskPair{
		//1.初步更改数据库
		MakeTaskPair("generalDeleteOrgUpdateDB", dmysql.MakeStepGeneralUpdateDelOrgIndent),
		//2. 创建删除组织脚本
		MakeTaskPair("generalDeleteOrgScriptStep", daction.MakeStepGeneralDeleteOrgScriptStep),
		//3. 打包所有文件
		MakeTaskPair("generalDeleteOrgZip", daction.MakeStepGeneralDeleteOrgZip),
		//4. 调用工具命令上传
		MakeTaskPair("generalDeleteOrgZipUpload", daction.MakeStepGeneralDeleteOrgZipUpload),
		//5. 调用工具解压
		MakeTaskPair("generalDeleteOrgUnZip", daction.MakeStepGeneralDeleteOrgUnzip),
		//6. 调用工具执行脚本先更新块配置
		MakeTaskPair("generalDeleteOperateOrgExecScript", daction.MakeStepGeneralDeleteOperateOrgExecScript),
		//7. 调用工具执行脚本
		MakeTaskPair("generalDeleteOrgExecScript", daction.MakeStepGeneralDeleteOrgExecScript),
		//8. 调用工具停掉检测节点状态脚本
		MakeTaskPair("generalDeleteOrgExecPeerStatusScript", daction.MakeStepGeneralDeleteOrgExecPeerStatusScript),
		//9. 删除本地证书
		MakeTaskPair("generalDeleteOrgLocalCife", daction.MakeStepDeleteOrgLocalCife),
		//字段写入数据库
		MakeTaskPair("generalDeleteOrgUpdateDB", dmysql.MakeStepGeneralUpdateDeleteOrgIndent),
	}
}

//NewGeneralAddPeerStep 新增节点具体操作
func NewGeneralAddPeerStep() []objectdefine.SubTaskPair {

	return []objectdefine.SubTaskPair{
		//初步写入数据库
		MakeTaskPair("generalCreatePeerWriteDB", dmysql.MakeStepGeneralInsertAddPeerIndent),
		//1.创建证书Peer
		MakeTaskPair("generalCreatePeerCife", daction.MakeStepGeneralAddPeerCife),
		//2. 创建新增节点脚本Peer
		MakeTaskPair("generalCreatePeerScriptStep", daction.MakeStepGeneralPeerScriptStep),
		//3. 创建base.yaml文件
		MakeTaskPair("generalCreatePeerBaseFile", daction.MakeStepGeneralPeerBaseFile),
		//5. 打包所有文件
		MakeTaskPair("generalCreatePeerZip", daction.MakeStepGeneralPeerZip),
		//6. 调用工具命令上传
		MakeTaskPair("generalCreatePeerZipUpload", daction.MakeStepGeneralPeerZipUpload),
		//7. 调用工具解压
		MakeTaskPair("generalCreatePeerUnZip", daction.MakeStepGeneralPeerUnzip),
		//8. 调用工具执行脚本
		MakeTaskPair("generalCreatePeerExecScript", daction.MakeStepGeneralPeerExecScript),
		//9. 调用工具启动检测节点状态脚本
		MakeTaskPair("generalCreatePeerExecPeerStatusScript", daction.MakeStepGeneralPeerExecPeerStatusScript),
		//字段写入数据ku
		MakeTaskPair("generalCreatePeerUpdateDB", dmysql.MakeStepGeneralUpdatePeerIndent),
	}
}

//NewGeneralDeletePeerStep 删除节点具体操作
func NewGeneralDeletePeerStep() []objectdefine.SubTaskPair {

	return []objectdefine.SubTaskPair{
		//1.初步更改数据库
		MakeTaskPair("generalDeletePeerUpdateDB", dmysql.MakeStepGeneralUpdateDelpeerIndent),
		//2. 创建删除节点脚本
		MakeTaskPair("generalDeletePeerScriptStep", daction.MakeStepGeneralDeletePeerScriptStep),
		//3. 打包所有文件
		MakeTaskPair("generalDeletePeerZip", daction.MakeStepGeneralDeletePeerZip),
		//4. 调用工具命令上传
		MakeTaskPair("generalDeletePeerZipUpload", daction.MakeStepGeneralDeletePeerZipUpload),
		//5. 调用工具解压
		MakeTaskPair("generalDeletePeerUnZip", daction.MakeStepGeneralDeletePeerUnzip),
		//6. 调用工具执行脚本
		MakeTaskPair("generalDeletePeerExecScript", daction.MakeStepGeneralDeletePeerExecScript),
		//7. 调用工具停掉检测节点状态脚本
		MakeTaskPair("generalDeletePeerExecPeerStatusScript", daction.MakeStepGeneralDeletePeerExecPeerStatusScript),
		//8. 删除本地证书
		MakeTaskPair("generalDeletePeerLocalCife", daction.MakeStepDeletePeerLocalCife),
		//字段写入数据库
		MakeTaskPair("generalDeletePeerUpdateDB", dmysql.MakeStepGeneralUpdateDeletePeerIndent),
	}
}

//NewGeneralDisablePeerStep 节点停用
func NewGeneralDisablePeerStep() []objectdefine.SubTaskPair {

	return []objectdefine.SubTaskPair{
		//1. 调用工具执行命令
		MakeTaskPair("generalDisablePeerExecCommand", daction.MakeStepGeneralDisablePeerExecCommand),
		//更新数据库状态
		MakeTaskPair("generalDisablePeerUpdateStatus", dmysql.MakeStepGeneralUpdatePeerDisableStatus),
	}
}

//NewGeneralEnablePeerStep 节点启用
func NewGeneralEnablePeerStep() []objectdefine.SubTaskPair {

	return []objectdefine.SubTaskPair{
		//1. 调用工具执行命令
		MakeTaskPair("generalEnablePeerExecCommand", daction.MakeStepGeneralEnablePeerExecCommand),
		//更新数据库状态
		MakeTaskPair("generalEnablePeerUpdateStatus", dmysql.MakeStepGeneralUpdatePeerEnableStatus),
	}
}

//NewGeneralModiflyPeerStep 节点启用
func NewGeneralModiflyPeerStep() []objectdefine.SubTaskPair {

	return []objectdefine.SubTaskPair{
		//1. 调用工具执行命令
		MakeTaskPair("generalModiflyPeerUpdateDB", dmysql.MakeStepGeneralUpdatePeerModiflyNickName),
	}
}

//NewGeneralChainCodeAddStep 新增合约
func NewGeneralChainCodeAddStep() []objectdefine.SubTaskPair {

	return []objectdefine.SubTaskPair{
		//初步写入数据库
		MakeTaskPair("generalCreateChainCodeWriteDB", dmysql.MakeStepGeneralInsertAddChainCodeIndent),
		//1.创建脚本文件 拷贝合约文件
		MakeTaskPair("generalCreateChainCodeScript", daction.MakeStepGeneralCreateChainCodeScript),
		//2. 打包
		MakeTaskPair("generalCreateChainCodeZip", daction.MakeStepGeneralChainCodeZip),
		//3. 上传
		MakeTaskPair("generalCreateChainCodeZipUpload", daction.MakeStepGeneralChainCodeZipUpload),
		//4. 调用工具解压
		MakeTaskPair("generalCreateChainCodeUnZip", daction.MakeStepGeneralChainCodeUnzip),
		//5. 调用工具执行脚本
		MakeTaskPair("generalCreateChainCodeExecScript", daction.MakeStepGeneralChainCodeExecScript),
		//字段写入数据ku
		MakeTaskPair("generalCreateChainCodeUpdateDB", dmysql.MakeStepGeneralUpdateAddChainCodeIndent),
	}
}

//NewGeneralChainCodeDeleteStep 删除合约
func NewGeneralChainCodeDeleteStep() []objectdefine.SubTaskPair {

	return []objectdefine.SubTaskPair{
		// //初步写入数据库
		// MakeTaskPair("generalDeleteChainCodeWriteDB", dmysql.MakeStepGeneralUpdateDelChainCodeIndent),
		//检测失败链码是否实例化 实例化不可删除
		MakeTaskPair("generalDeleteChaincodeCheckChaincodeIsInstantiated",dssh.MakeStepGeneralCheckChaincodeIsInstantiated),
		//1.调用命令删除合约
		MakeTaskPair("generalCreateDeleteChainCodeCommand", daction.MakeStepGeneralDeleteChainCodeCommand),
		//2. 打包
		//MakeTaskPair("generalCreateDeleteChainCodeZip", daction.MakeStepGeneralDeleteChainCodeZip),
		//3. 上传
		//MakeTaskPair("generalCreateDeleteChainCodeZipUpload", daction.MakeStepGeneralDeleteChainCodeZipUpload),
		//4. 调用工具解压
		//MakeTaskPair("generalDeleteChainCodeUnZip", daction.MakeStepGeneralDeleteChainCodeUnzip),
		//5. 调用工具执行脚本
		//MakeTaskPair("generalDeleteChainCodeExecScript", daction.MakeStepGeneralDeleteChainCodeExecScript),
		//字段写入数据ku
		MakeTaskPair("generalDeleteChainCodeUpdateDB", dmysql.MakeStepGeneralUpdateDeleteChainCodeIndent),
	}
}

//NewGeneralChainCodeUpgradeStep 升级合约
func NewGeneralChainCodeUpgradeStep() []objectdefine.SubTaskPair {

	return []objectdefine.SubTaskPair{
		//初步写入数据库
		MakeTaskPair("generalUpgradeChainCodeWriteDB", dmysql.MakeStepGeneralInsertUpgradeChainCodeIndent),
		//1.创建脚本文件 拷贝合约文件
		MakeTaskPair("generalCreateUpgradeChainCodeScript", daction.MakeStepGeneralUpgradeChainCodeScript),
		//2. 打包
		MakeTaskPair("generalCreateUpgradeChainCodeZip", daction.MakeStepGeneralUpgradeChainCodeZip),
		//3. 上传
		MakeTaskPair("generalUpgradeChainCodeZipUpload", daction.MakeStepGeneralUpgradeChainCodeZipUpload),
		//4. 调用工具解压
		MakeTaskPair("generalUpgradeChainCodeUnZip", daction.MakeStepGeneralUpgradeChainCodeUnzip),
		//5. 调用工具执行脚本
		MakeTaskPair("generalUpgradeChainCodeExecScript", daction.MakeStepGeneralUpgradeChainCodeExecScript),
		//更新数据库
		MakeTaskPair("generalUpgradeChainCodeUpdateDB", dmysql.MakeStepGeneralUpdateUpgradeChainCodeIndent),
	}
}

//NewGeneralChainCodeDisableStep 合约停用
func NewGeneralChainCodeDisableStep() []objectdefine.SubTaskPair {

	return []objectdefine.SubTaskPair{
		//初步写入数据库
		MakeTaskPair("generalDisableChainCodeUpdateDB", dmysql.MakeStepGeneralDisableChainCodeUpdateIndent),
		//1.创建脚本文件
		MakeTaskPair("generalCreateDisableChainCodeScript", daction.MakeStepGeneralDisableChainCodeScript),
		//2. 打包
		MakeTaskPair("generalCreateDisableChainCodeZip", daction.MakeStepGeneralDisableChainCodeZip),
		//3. 上传
		MakeTaskPair("generalDisableChainCodeZipUpload", daction.MakeStepGeneralDisableChainCodeZipUpload),
		//4. 调用工具解压
		MakeTaskPair("generalDisableChainCodeUnZip", daction.MakeStepGeneralDisableChainCodeUnzip),
		//5. 调用工具执行脚本
		MakeTaskPair("generalDisableChainCodeExecScript", daction.MakeStepGeneralDisableChainCodeExecScript),
		//字段写入数据ku
		MakeTaskPair("generalDisableChainCodeWriteDB", dmysql.MakeStepGeneralUpdateDisableChainCodeIndent),
	}
}

//NewGeneralChainCodeEnableStep 合约启用
func NewGeneralChainCodeEnableStep() []objectdefine.SubTaskPair {

	return []objectdefine.SubTaskPair{
		//初步写入数据库
		MakeTaskPair("generalEnableChainCodeUpdateDB", dmysql.MakeStepGeneralEnableChainCodeUpdateIndent),
		//1.创建脚本文件
		MakeTaskPair("generalCreateEnableChainCodeScript", daction.MakeStepGeneralEnableChainCodeScript),
		//2. 打包
		MakeTaskPair("generalCreateEnableChainCodeZip", daction.MakeStepGeneralEnableChainCodeZip),
		//3. 上传
		MakeTaskPair("generalEnableChainCodeZipUpload", daction.MakeStepGeneralEnableChainCodeZipUpload),
		//4. 调用工具解压
		MakeTaskPair("generalEnableChainCodeUnZip", daction.MakeStepGeneralEnableChainCodeUnzip),
		//5. 调用工具执行脚本
		MakeTaskPair("generalEnableChainCodeExecScript", daction.MakeStepGeneralEnableChainCodeExecScript),
		//字段写入数据ku
		MakeTaskPair("generalEnableChainCodeWriteDB", dmysql.MakeStepGeneralUpdateEnableChainCodeIndent),
	}
}
 
//NewGeneralChainCodeBulidStep 合约编译
func NewGeneralChainCodeBulidStep() []objectdefine.SubTaskPair {

	return []objectdefine.SubTaskPair{
		//编译脚本拷贝到合约目录下
        MakeTaskPair("generalChainCodeBulidScriptCopy", daction.MakeStepGeneralChainCodeBulidScripCopy),
		//脚本执行
		MakeTaskPair("generalChainCodeBulidExecScript", daction.MakeStepGeneralChainCodeBulidExecScript),
	}
}

//NewGeneralAddServerInfo
func NewGeneralAddServerInfoStep() []objectdefine.SubTaskPair {
	return []objectdefine.SubTaskPair{
		//检测连接 部署工具
		MakeTaskPair("generalAddServerInfoCheck", dssh.MakeStepGeneralCheckRemoteServerSShConnect),
		//检测连接 部署工具
		MakeTaskPair("generalAddServerInfoToolStart", dssh.MakeStepGeneralCheckRemotStartServerTools),
		//构建环境
		MakeTaskPair("generalAddServerInfoEnv", dssh.MakeStepGeneralRemoteServerBulidEnv),
		//服务机器信息写入数据库
		MakeTaskPair("generalAddServerInfoInsertDB", dmysql.MakeStepGeneralRemoteServerUpdateDB),
	}
}