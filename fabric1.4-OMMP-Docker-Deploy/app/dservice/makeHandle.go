package dservice

import (
	"deploy-server/app/daction"
	"deploy-server/app/dcache"
	"deploy-server/app/dmysql"
	"deploy-server/app/dtask"
	"deploy-server/app/objectdefine"
	"deploy-server/app/tools"
	"deploy-server/app/web"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func responseJSONBuff(w http.ResponseWriter, buff []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.Write(buff)
}

func makeIndentCompleteDeploy(debug, db bool) http.HandlerFunc {
	return makeGeneralTaskBase(dtask.TaskGeneralComplete, "completeDeploy", true, false, false, false)
}

func makeGeneralCreateChannel(debug, db bool) http.HandlerFunc {
	return makeGeneralTaskBase(dtask.TaskGeneralCreateChannel, "createChannel", true, false, false, false)
}

func makeGeneralAddOrg(debug, db bool) http.HandlerFunc {
	return makeGeneralTaskBase(dtask.TaskGeneralAddOrg, "addOrg", true, false, false, false)
}

func makeGeneralDeleteOrg(debug, db bool) http.HandlerFunc {
	return makeGeneralTaskBase(dtask.TaskGeneralDeleteOrg, "deleteOrg", true, false, false, false)
}

func makeGeneralAddPeer(debug, db bool) http.HandlerFunc {
	return makeGeneralTaskBase(dtask.TaskGeneralAddPeer, "addPeer", true, false, false, false)
}

func makeGeneralDeletePeer(debug, db bool) http.HandlerFunc {
	return makeGeneralTaskBase(dtask.TaskGeneralDeletePeer, "deletePeer", true, false, false, false)
}

func makeGeneralDisablePeer(debug, db bool) http.HandlerFunc {
	return makeGeneralTaskBase(dtask.TaskGeneralDisablePeer, "disablePeer", true, false, false, false)
}

func makeGeneralEnablePeer(debug, db bool) http.HandlerFunc {
	return makeGeneralTaskBase(dtask.TaskGeneralEnablePeer, "enablePeer", true, false, false, false)
}

//makeGeneralPeerNickNameModify 节点昵称修改
func makeGeneralPeerNickNameModify(debug, db bool) http.HandlerFunc {
	return makeGeneralTaskBase(dtask.TaskGeneralModiflyPeer, "modiflyPeer", true, false, false, false)
	// return func(w http.ResponseWriter, r *http.Request) {
	// 	_, body := web.GetParamsBody(r)
	// 	if len(body) == 0 {
	// 		web.OutputEnter(w, "", nil, errors.WithMessage(nil, "Read body empty"))
	// 		return
	// 	}
	// 	ret := &objectdefine.Indent{}
	// 	var uJSON = jsoniter.ConfigCompatibleWithStandardLibrary
	// 	err := uJSON.Unmarshal(body, ret)
	// 	if err != nil {
	// 		web.OutputEnter(w, "", nil, errors.WithMessage(nil, "Read body json unmarshal error"))
	// 		return
	// 	}
	// 	if len(ret.Org) == 0 {
	// 		web.OutputEnter(w, "", nil, errors.WithMessage(nil, "indent org is empty"))
	// 		return
	// 	}
	// 	for _, org := range ret.Org {
	// 		for _, peer := range org.Peer {
	// 			err := dmysql.UpdateIndetPeerNickName(&peer)
	// 			if err != nil {
	// 				web.OutputEnter(w, "", nil, errors.WithMessage(nil, "indent update peer nickname error"))
	// 				return
	// 			}
	// 		}
	// 	}
	// 	web.OutputEnter(w, "", nil, nil)
	// }
}

func makeGeneralCreateChainCodeUpload(debug, db bool) http.HandlerFunc {
	return makeGeneralTaskBase(dtask.TaskGeneralChainCodeUpload, "newChaincodeUpload", false, false, false, false)
}
func makeGeneralUpgradeChainCodeUpload(debug, db bool) http.HandlerFunc {
	return makeGeneralTaskBase(dtask.TaskGeneralChainCodeUpload, "upgradeChaincodeUpload", false, false, false, false)
}

func makeGeneralChainCodeList(debug, db bool) http.HandlerFunc {
	return makeGeneralTaskBase(dtask.TaskGeneralChainCodeList, "chaincodelist", true, false, false, false)
}

func makeGeneralChainCodeAdd(debug, db bool) http.HandlerFunc {
	return makeGeneralTaskBase(dtask.TaskGeneralChainCodeAdd, "chaincodeAdd", true, false, false, false)
}

func makeGeneralChainCodeDelete(debug, db bool) http.HandlerFunc {
	return makeGeneralTaskBase(dtask.TaskGeneralChainCodeDelete, "chaincodeDelete", true, false, false, false)
}

func makeGeneralChainCodeUpgrade(debug, db bool) http.HandlerFunc {
	return makeGeneralTaskBase(dtask.TaskGeneralChainCodeUpgrade, "chaincodeUpgrade", true, false, false, false)
}

func makeGeneralChainCodeDisable(debug, db bool) http.HandlerFunc {
	return makeGeneralTaskBase(dtask.TaskGeneralChainCodeDisable, "chaincodeDisable", true, false, false, false)
}

func makeGeneralChainCodeEnable(debug, db bool) http.HandlerFunc {
	return makeGeneralTaskBase(dtask.TaskGeneralChainCodeEnable, "chaincodeEnable", true, false, false, false)
}

//makeIndentStatus 获取任务状态
func makeIndentStatus(debug, db bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, body := web.GetParamsBody(r)
		ret := &objectdefine.Indent{}
		err := json.Unmarshal(body, ret)
		if err != nil {
			web.OutputEnter(w, "", nil, errors.WithMessage(err, "Read body error"))
			return
		}

		if len(ret.ID) == 0 {
			web.OutputEnter(w, "", nil, errors.WithMessage(err, "Empty buff for ID"))
			return
		}
		if false == tools.IsAlpha(ret.ID) || len(ret.ID) > 128 {
			web.OutputEnter(w, "", nil, errors.WithMessage(err, "ID format error"))
			return
		}
		ret.ID = strings.ToLower(ret.ID)

		ret.BaseOutput = dcache.GetOutputSubPath(ret.ID, "")

		status := dtask.GetTaskStatus(ret.ID)
		if status == nil {
			web.OutputEnter(w, "", nil, errors.New("Task ["+ret.ID+"] is not exist"))
			return
		}
		web.OutputEnter(w, "", status, nil)
	}
}

//makeIndentConfig 返回订单补全的信息
func makeIndentConfig(debug, db bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, body := web.GetParamsBody(r)
		one, err := dcache.GetGeneralOrgFromBuff(body)
		if err != nil {
			web.OutputEnter(w, "", nil, errors.WithMessage(err, "Read body error"))
			return
		}
		two, err := dmysql.GetStartTaskBeforIndent(one)
		if err != nil {
			web.OutputEnter(w, "", nil, err)
			return
		}
		web.OutputEnter(w, "", two, nil)
	}
}

//makeGetIPList 获取ip列表
func makeGetIPList(debug, db bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		two, err := dmysql.GetIPListFromIndent()
		if err != nil {
			web.OutputEnter(w, "", nil, err)
			return
		}
		web.OutputEnter(w, "", two, nil)
	}
}

//makeGetIPPortList 获取当期ip使用过的所有端口
func makeGetIPPortList(debug, db bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		queryForm, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil && len(queryForm["ip"]) == 0 {
			web.OutputEnter(w, "", nil, err)
			return
		}
		two, err := dmysql.GetIPPortListFromIndent(queryForm["ip"][0])
		if err != nil {
			web.OutputEnter(w, "", nil, err)
			return
		}
		web.OutputEnter(w, "", two, nil)
	}
}

//makeGetChannelList 获取通道列表
func makeGetChannelList(debug, db bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		queryForm, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			web.OutputEnter(w, "", nil, err)
			return
		}
		var pageSize, pageNum int
		if len(queryForm["pagenum"]) == 0 {
			pageNum = 0
		} else {
			pageNum, _ = strconv.Atoi(queryForm["pagenum"][0])
		}
		if len(queryForm["pagesize"]) == 0 {
			pageSize = 10
		} else {
			pageSize, _ = strconv.Atoi(queryForm["pagesize"][0])
		}
		if pageSize == 0 {
			web.OutputEnter(w, "", nil, errors.New("pageSize cannot be zero "))
			return
		}
		two, totalPage, err := dmysql.GetChannelListFromIndent(pageNum, pageSize)
		if err != nil || totalPage == -1 {
			web.OutputEnter(w, "", nil, err)
			return
		}

		web.OutputListEnter(w, "", totalPage, two, nil)

	}
}

//makeGetChannelOrgInfoList 获取通道下组织列表
func makeGetChannelOrgInfoList(debug, db bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		queryForm, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			web.OutputEnter(w, "", nil, err)
			return
		}
		var channelName string
		if len(queryForm["channel"]) == 0 {
			channelName = ""
		} else {
			channelName = queryForm["channel"][0]
		}

		var pageSize, pageNum int
		if len(queryForm["pagenum"]) == 0 {
			pageNum = 0
		} else {
			pageNum, _ = strconv.Atoi(queryForm["pagenum"][0])
		}
		if len(queryForm["pagesize"]) == 0 {
			pageSize = 10
		} else {
			pageSize, _ = strconv.Atoi(queryForm["pagesize"][0])
		}
		if pageSize == 0 {
			web.OutputEnter(w, "", nil, errors.New("pageSize cannot be zero "))
			return
		}
		two, totalPage, err := dmysql.GetChannelOrgInfoListFromIndent(channelName, pageNum, pageSize)
		if err != nil || totalPage == -1 {
			web.OutputEnter(w, "", nil, err)
			return
		}
		web.OutputListEnter(w, "", totalPage, two, nil)
	}
}

//makeGetChannelOrgPeerInfoList 获取节点列表
func makeGetChannelOrgPeerInfoList(debug, db bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		queryForm, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			web.OutputEnter(w, "", nil, err)
			return
		}
		var channelName string
		if len(queryForm["channel"]) == 0 {
			channelName = ""
		} else {
			channelName = queryForm["channel"][0]
		}
		var orgName string
		if len(queryForm["orgname"]) == 0 {
			orgName = ""
		} else {
			orgName = queryForm["orgname"][0]
		}

		var pageSize, pageNum int
		if len(queryForm["pagenum"]) == 0 {
			pageNum = 0
		} else {
			pageNum, _ = strconv.Atoi(queryForm["pagenum"][0])
		}
		if len(queryForm["pagesize"]) == 0 {
			pageSize = 10
		} else {
			pageSize, _ = strconv.Atoi(queryForm["pagesize"][0])
		}
		if pageSize == 0 {
			web.OutputEnter(w, "", nil, errors.New("pageSize cannot be zero "))
			return
		}

		two, totalPage, err := dmysql.GetChannelOrgPeerInfoListFromIndent(channelName, orgName, pageNum, pageSize)
		if err != nil || totalPage == -1 {
			web.OutputEnter(w, "", nil, err)
			return
		}
		web.OutputListEnter(w, "", totalPage, two, nil)

	}
}

//makeGetchainCodeInfoList 获取合约列表
func makeGetchainCodeInfoList(debug, db bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		queryForm, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			web.OutputEnter(w, "", nil, err)
		}
		var channelName string
		if len(queryForm["channel"]) == 0 {
			channelName = ""
		} else {
			channelName = queryForm["channel"][0]
		}
		var pageSize, pageNum int
		if len(queryForm["pagenum"]) == 0 {
			pageNum = 0
		} else {
			pageNum, _ = strconv.Atoi(queryForm["pagenum"][0])
		}
		if len(queryForm["pagesize"]) == 0 {
			pageSize = 10
		} else {
			pageSize, _ = strconv.Atoi(queryForm["pagesize"][0])
		}
		if pageSize == 0 {
			web.OutputEnter(w, "", nil, errors.New("pageSize cannot be zero "))
			return
		}
		two, totalPage, err := dmysql.GetChainCodeInfoListFromIndent(channelName, pageNum, pageSize)
		if err != nil || totalPage == -1 {
			web.OutputEnter(w, "", nil, err)
			return
		}
		web.OutputListEnter(w, "", totalPage, two, nil)
	}
}

func makeGeneralTaskBase(kind, model string, bindent, bid, bbuild, del bool) http.HandlerFunc {
	core := dtask.MakePushTaskHandle(kind)
	return func(w http.ResponseWriter, r *http.Request) {
		switch model {
		case "newChaincodeUpload":
			err := dmysql.ReceiveChainCodeUploadFile(w, r)
			if err != nil {
				web.OutputEnter(w, "", nil, errors.WithMessage(err, "receive chaincode upload file error"))
				return
			}
			web.OutputEnter(w, "", nil, nil)
			return
		case "upgradeChaincodeUpload":
			err := dmysql.ReceiveUpgradeChainCodeUploadFile(w, r)
			if err != nil {
				web.OutputEnter(w, "", nil, errors.WithMessage(err, "receive chaincode upload file error"))
				return
			}
			web.OutputEnter(w, "", nil, nil)
			return
		default:
		}
		fmt.Println("general get  body info:")
		_, body := web.GetParamsBody(r)
		if len(body) == 0 {
			web.OutputEnter(w, "", nil, errors.New("Read body empty"))
			return
		}

		fmt.Println("request body info:", string(body))
		//首先检查基础变量
		one, err := dcache.GetGeneralOrgFromBuff(body)
		if err != nil {
			web.OutputEnter(w, "", nil, errors.WithMessage(err, "Read body error"))
			return
		}
		if one == nil {
			web.OutputEnter(w, "", nil, errors.WithMessage(err, "Read body error"))
			return
		}
		//检查任务id是否存在
		err = CheckID(one.ID, true)
		if err != nil {
			web.OutputEnter(w, "", nil, errors.WithMessage(err, "general add  union ID is not matching"))
			return
		}
		//订单补全
		if bindent {
			completion, err := daction.CheckCompletionIndent(one, model)
			if err != nil {
				web.OutputEnter(w, "", nil, errors.WithMessage(err, "Check indent error"))
				return
			}
			completionByte, _ := json.Marshal(completion)
			fmt.Println("completion", string(completionByte))
			one = completion
		}
		general, err := core(one)
		if err != nil {
			web.OutputEnter(w, "", nil, errors.WithMessage(err, "Task is rejected"))
			return
		}
		web.OutputEnter(w, "", general, nil)
		return
	}
}
