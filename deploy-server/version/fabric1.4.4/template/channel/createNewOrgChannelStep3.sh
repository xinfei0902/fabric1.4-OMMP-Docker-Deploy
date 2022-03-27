#!/bin/bash
## 总存放目录 相关文件脚本存放文件夹名称
TARGET=$[target]
CHANNEL_NAME=$[channelName]
CHANNEL_CONFIG=./channel-artifacts/${CHANNEL_NAME}.tx
NEW_ORG_CLI_NAME=$[newOrgCliName]
## orderer名称加端口 例如：orderer0.example.com:7050
ORDERER_NAME=$[ordererAddress]
ORDERER_CA=$[ordererTlsCa]
##节点配置
ORG_NAMEMSP=$[orgMSPID]
ORG_PEER_NUM=$[peerNum]
PEER_TLS_PATH=$[peerTlsCa]
PEER_MSP_PATH=$[peerUsersMsp]
PEER_ADDRESS=$[peerAddress]
## 锚节点
PEER_ADMIN_ANCHOR=$[peerAnchor]
PEER_ADMIN_TLS_PATH=$[peerAnchorTls]
PEER_ADMIN_MSP_PATH=$[peerAnchorMsp]
PEER_ADMIN_ADDRESS=$[peerAnchorAddreess]

verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "Error !!! ERROR  FAILE : $2" >&2
    exit 1
  fi
}

##创建通道
set -x
peer channel create -o $ORDERER_NAME -c $CHANNEL_NAME -f $CHANNEL_CONFIG --tls true --cafile $ORDERER_CA
ret=$?
set +x
verifyResult $ret "peer create channel: CHANNEL_NAME failed"


## 所有组织加入通道
orgNum=${#ORG_NAMEMSP[@]}
for ((j=0;j<orgNum;j++))
 {
   peerN=${ORG_PEER_NUM[$j]}
   if [ $j -ne 0 ];then
      peerAf=$(($j-1))
      peerTotal=$(($peerTotal+${ORG_PEER_NUM[$peerAf]}))
      echo "peerTotal=$peerTotal"
   fi
   for ((k=0;k<peerN;k++))
    {
        num=$(($peerTotal+$k))
        echo "num=$num"
        export CORE_PEER_LOCALMSPID=${ORG_NAMEMSP[$j]}
		export CORE_PEER_TLS_ROOTCERT_FILE=${PEER_TLS_PATH[$num]}
		export CORE_PEER_MSPCONFIGPATH=${PEER_MSP_PATH[$num]}
		export CORE_PEER_ADDRESS=${PEER_ADDRESS[$num]}
		set -x
		peer channel join -b $CHANNEL_NAME.block
		ret=$?
		set +x
		verifyResult $ret "peer join channel  ${CHANNEL_NAME} failed"
    }
 }

##定义锚节点
orgNum=${#ORG_NAMEMSP[@]}
for ((i=0;i<orgNum;i++))
{
	export CORE_PEER_LOCALMSPID=${ORG_NAMEMSP[$i]}
	export CORE_PEER_TLS_ROOTCERT_FILE=${PEER_ADMIN_TLS_PATH[$i]}
	export CORE_PEER_MSPCONFIGPATH=${PEER_ADMIN_MSP_PATH[$i]}
	export CORE_PEER_ADDRESS=${PEER_ADMIN_ADDRESS[$i]}
	set -x
	peer channel update -o $ORDERER_NAME -c $CHANNEL_NAME -f ./channel-artifacts/${ORG_NAMEMSP[$i]}Panchors.tx --tls true --cafile $ORDERER_CA 
	ret=$?
	set +x
	verifyResult $ret "peer update channel  $CHANNEL_NAME failed"
}