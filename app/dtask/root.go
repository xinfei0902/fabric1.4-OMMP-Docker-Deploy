package dtask

import "deploy-server/app/objectdefine"

//NewGeneralCompleteDeploy  一键部署
func NewGeneralCompleteDeploy() *objectdefine.TaskType {
	return objectdefine.NewTaskType(
		NewGeneralCompleteDeployStep(),
		MakeGeneralStartStep(),
		MakeGeneralEndStep())
}

//NewGeneralCreateChannel 创建通道
func NewGeneralCreateChannel() *objectdefine.TaskType {
	return objectdefine.NewTaskType(
		NewGeneralCreateChannelStep(),
		MakeGeneralStartStep(),
		MakeGeneralEndStep())
}

//NewGeneralAddOrgTask 创建新组织
func NewGeneralAddOrgTask() *objectdefine.TaskType {

	return objectdefine.NewTaskType(
		NewGeneralAddOrgStep(),
		MakeGeneralStartStep(),
		MakeGeneralEndStep())
}

//NewGeneralDeleteOrgTask 删除组织
func NewGeneralDeleteOrgTask() *objectdefine.TaskType {

	return objectdefine.NewTaskType(
		NewGeneralDeleteOrgStep(),
		MakeGeneralStartStep(),
		MakeGeneralEndStep())
}

//NewGeneralAddPeerTask 创建新的节点
func NewGeneralAddPeerTask() *objectdefine.TaskType {

	return objectdefine.NewTaskType(
		NewGeneralAddPeerStep(),
		MakeGeneralStartStep(),
		MakeGeneralEndStep())
}

//NewGeneralDeletePeerTask 删除节点
func NewGeneralDeletePeerTask() *objectdefine.TaskType {

	return objectdefine.NewTaskType(
		NewGeneralDeletePeerStep(),
		MakeGeneralStartStep(),
		MakeGeneralEndStep())
}

//NewGeneralDisablePeerTask 停用节点
func NewGeneralDisablePeerTask() *objectdefine.TaskType {

	return objectdefine.NewTaskType(
		NewGeneralDisablePeerStep(),
		MakeGeneralStartStep(),
		MakeGeneralEndStep())
}

//NewGeneralEnablePeerTask 启用节点
func NewGeneralEnablePeerTask() *objectdefine.TaskType {

	return objectdefine.NewTaskType(
		NewGeneralEnablePeerStep(),
		MakeGeneralStartStep(),
		MakeGeneralEndStep())
}

//NewGeneralModiflyPeerTask 修改节点昵称
func NewGeneralModiflyPeerTask() *objectdefine.TaskType {

	return objectdefine.NewTaskType(
		NewGeneralModiflyPeerStep(),
		MakeGeneralStartStep(),
		MakeGeneralEndStep())
}

//NewGeneralChaincCodeAdd 创建新的合约
func NewGeneralChaincCodeAdd() *objectdefine.TaskType {

	return objectdefine.NewTaskType(
		NewGeneralChainCodeAddStep(),
		MakeGeneralStartStep(),
		MakeGeneralEndStep())
}

//NewGeneralChaincCodeDelete 删除合约
func NewGeneralChaincCodeDelete() *objectdefine.TaskType {

	return objectdefine.NewTaskType(
		NewGeneralChainCodeDeleteStep(),
		MakeGeneralStartStep(),
		MakeGeneralEndStep())
}

//NewGeneralChainCodeUpgrade 合约升级
func NewGeneralChainCodeUpgrade() *objectdefine.TaskType {

	return objectdefine.NewTaskType(
		NewGeneralChainCodeUpgradeStep(),
		MakeGeneralStartStep(),
		MakeGeneralEndStep())
}

//NewGeneralChainCodeDisable 合约停用
func NewGeneralChainCodeDisable() *objectdefine.TaskType {

	return objectdefine.NewTaskType(
		NewGeneralChainCodeDisableStep(),
		MakeGeneralStartStep(),
		MakeGeneralEndStep())
}

//NewGeneralChainCodeEnable 合约启用
func NewGeneralChainCodeEnable() *objectdefine.TaskType {

	return objectdefine.NewTaskType(
		NewGeneralChainCodeEnableStep(),
		MakeGeneralStartStep(),
		MakeGeneralEndStep())
}
