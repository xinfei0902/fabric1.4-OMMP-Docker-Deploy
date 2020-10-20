package dservice

import (
	"deploy-server/app/web"
	"net/http"
)

//RegisterWebAPI 注册接口地址与接口函数之间关系
func RegisterWebAPI(debug bool, db bool) (err error) {

	base := "/api/deploy/manage/"
	// with / without db
	list := []httpHandleFuncPair{
		{base + "indent/status", makeIndentStatus},
		{base + "indent/config", makeIndentConfig},
		{base + "indent/complete/deploy", makeIndentCompleteDeploy},
		//获取ip列表以及档期机器安装几个节点
		{base + "indent/get/ip/list", makeGetIPList},
		//当前机器已经使用的端口列表
		{base + "indent/get/ip/port/list", makeGetIPPortList},
		//获取通道列表
		{base + "indent/get/channel/list", makeGetChannelList},
		//获取通道下组织列表
		{base + "indent/get/org/query/list", makeGetChannelOrgInfoList},
		//获取通道组织下所有节点列表
		{base + "indent/get/peer/query/list", makeGetChannelOrgPeerInfoList},
		//获取合约列表
		{base + "indent/get/chaincode/query/list", makeGetchainCodeInfoList},
		//通道
		{base + "indent/general/channel/create", makeGeneralCreateChannel},
		//组织
		{base + "indent/general/org/add", makeGeneralAddOrg},
		{base + "indent/general/org/delete", makeGeneralDeleteOrg},
		//节点
		{base + "indent/general/peer/add", makeGeneralAddPeer},
		{base + "indent/general/peer/delete", makeGeneralDeletePeer},
		{base + "indent/general/peer/disable", makeGeneralDisablePeer},
		{base + "indent/general/peer/enable", makeGeneralEnablePeer},
		{base + "indent/general/peer/nickname/modify", makeGeneralPeerNickNameModify},
		//合约
		//需要新增一个接受上传合约链码接口
		{base + "indent/general/create/chaincode/upload", makeGeneralCreateChainCodeUpload},
		{base + "indent/general/upgrade/chaincode/upload", makeGeneralUpgradeChainCodeUpload},
		//获取本地合约列表
		//{base + "indent/general/chaincode/list", makeGeneralChainCodeList},
		{base + "indent/general/chaincode/add", makeGeneralChainCodeAdd},
		{base + "indent/general/chaincode/delete", makeGeneralChainCodeDelete},
		{base + "indent/general/chaincode/upgrade", makeGeneralChainCodeUpgrade},
		{base + "indent/general/chaincode/disable", makeGeneralChainCodeDisable},
		{base + "indent/general/chaincode/enable", makeGeneralChainCodeEnable},
	}

	// API
	for _, one := range list {
		if one.f == nil {
			continue
		}
		f := one.f(debug, db)
		if f == nil {
			continue
		}
		err = web.PushHandleFunc(one.p, f)
		if err != nil {
			return
		}
	}

	listHandler := []httpHandlePair{
		// {base + "store/packages/download/", makePackageDownload},
	}

	for _, one := range listHandler {
		if one.f == nil {
			continue
		}
		f := one.f(debug, db)
		if f == nil {
			continue
		}

		err = web.PushHandle(one.p, http.StripPrefix(one.p, f))
		if err != nil {
			return
		}
	}

	return nil
}
