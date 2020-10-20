#!/bin/bash
# 文件目录
TARGET=$[target]
# 新增组织 MSPID名称
ORG_NAMEMSP=$[orgMSPID]
##通道名称
CHANNEL_NAME=$[channelName]
## orderer名称加端口 例如：orderer0.prisons.guizhou:7050
ORDERER_NAME=$[ordererAddress]
ORDERER_TLS_PATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/$[ordererTlsCa]
ORDERER_MSP_PATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/$[ordererUsersMsp]
##orderer-ca证书路径
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/$[ordererTlsCa]
#是否开启tls模式
CORE_PEER_TLS_ENABLED=true
#签名配置
##SIGN_ORG_NAME 包含所有组织名称(MSPID)且数组元素位置按序排列 例如（Org1MSP Org2MSP） 
SIGN_ORG_NAME=$[signOrgMSPID]
##证书路径 与上面组织名称一一对应
SIGN_ORG_TLS_PATH=$[signOrgTlsCa]
SIGN_ORG_MSP_PATH=$[signOrgUsersMsp]
SIGN_ORG_ADDRESS=$[signOrgPeerHosts]

##设置orderer环境变量
setOrdererGlobals() {
  CORE_PEER_LOCALMSPID="OrdererMSP"
  CORE_PEER_TLS_ROOTCERT_FILE=${ORDERER_TLS_PATH}
  CORE_PEER_MSPCONFIGPATH=${ORDERER_MSP_PATH}
}

##获取通道配置块
fetchChannelConfig() {
  CHANNEL=$1
  OUTPUT=$2

  setOrdererGlobals

  echo "Fetching the most recent configuration block for the channel"
  if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
    set -x
    peer channel fetch config config_block.pb -o $ORDERER_NAME -c $CHANNEL_NAME --cafile $ORDERER_CA
    set +x
  else
    set -x
    peer channel fetch config config_block.pb -o $ORDERER_NAME -c $CHANNEL_NAME --tls --cafile $ORDERER_CA >peerFetch.log 2>&1
    set +x
    res=$(cat peerFetch.log)
    expr index "$res" "Error"
    ret=$?
    if [ $ret -ne 0 ];then
          echo "Error peer channel fetch faile !!!" >&2
          exit 1
    fi
  fi

  echo "Decoding config block to JSON and isolating config to ${OUTPUT}"
  set -x
  ./bin/configtxlator proto_decode --input config_block.pb --type common.Block | jq .data.data[0].payload.data.config >"${OUTPUT}"
  set +x
}

##更新获取的配置块
createConfigUpdate() {
  CHANNEL=$1
  ORIGINAL=$2
  MODIFIED=$3
  OUTPUT=$4

  set -x
  ./bin/configtxlator proto_encode --input "${ORIGINAL}" --type common.Config >original_config.pb
  ./bin/configtxlator proto_encode --input "${MODIFIED}" --type common.Config >modified_config.pb
  ./bin/configtxlator compute_update --channel_id "${CHANNEL}" --original original_config.pb --updated modified_config.pb >config_update.pb
  ./bin/configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate >config_update.json
  echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . >config_update_in_envelope.json
  ./bin/configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope >"${OUTPUT}" 
}

## 对信封pb文件 签名
signConfigtxAsPeerOrg() { 
  TX=$1
  num=${#SIGN_ORG_NAME[@]}
  echo "========num=${num}========="
  for ((i=0;i<num;i++))
  {
    CORE_PEER_LOCALMSPID=${SIGN_ORG_NAME[$i]}
    CORE_PEER_TLS_ROOTCERT_FILE=${SIGN_ORG_TLS_PATH[$i]}
    CORE_PEER_MSPCONFIGPATH=${SIGN_ORG_MSP_PATH[$i]}
    CORE_PEER_ADDRESS=${SIGN_ORG_ADDRESS[$i]}
    set -x
    peer channel signconfigtx -f "${TX}" 
    set +x
  }
  
}

##进入目录
##cd ${TARGET}
## 安装工具
which jq
if [ $? -ne 0 ];then
 echo "install jq"
 apt-get -y update && apt-get -y install jq
fi
## 执行命令
export FABRIC_CFG_PATH=$PWD
./bin/configtxgen -printOrg ${ORG_NAMEMSP} > ./${ORG_NAMEMSP}.json

export FABRIC_CFG_PATH=/etc/hyperledger/fabric/

##
echo "===================================="
echo "获取块配置"
echo "===================================="
fetchChannelConfig ${CHANNEL_NAME} config.json 
ret=$?
if [ $ret -ne 0 ];then
   echo "Error peer channel update  fetch block file error!!!!!!" >&2
   exit 1
fi
set -x
jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups": {"'${ORG_NAMEMSP}'":.[1]}}}}}' config.json ./${ORG_NAMEMSP}.json > modified_config.json  
set +x
echo "===================================="
echo "配置更新"
echo "===================================="
createConfigUpdate ${CHANNEL_NAME} config.json modified_config.json ${ORG_NAMEMSP}_update_in_envelope.pb 
ret=$?
if [ $ret -ne 0 ];then
   echo "Error peer channel update  update .pb file error!!!!!!" >&2
   exit 1
fi
##签名
echo "===================================="
echo "签名"
echo "===================================="
export FABRIC_CFG_PATH=/etc/hyperledger/fabric/
signConfigtxAsPeerOrg  ${ORG_NAMEMSP}_update_in_envelope.pb
ret=$?
if [ $ret -ne 0 ];then
   echo "Error peer channel update sign .pb file error!!!!!!" >&2
   exit 1
fi
##设置环境变量
export CORE_PEER_LOCALMSPID=${SIGN_ORG_NAME[0]}
export CORE_PEER_TLS_ROOTCERT_FILE=${SIGN_ORG_TLS_PATH[0]}
export CORE_PEER_MSPCONFIGPATH=${SIGN_ORG_MSP_PATH[0]}
export CORE_PEER_ADDRESS=${SIGN_ORG_ADDRESS[0]}
##提交
echo "===================================="
echo "提交"
echo "===================================="
set -x
peer channel update -f ${ORG_NAMEMSP}_update_in_envelope.pb -c ${CHANNEL_NAME} -o ${ORDERER_NAME} --tls --cafile ${ORDERER_CA} >>addOrg.log 2>&1
set +x
result=$(cat addOrg.log)
expr index "$result" "Error"
res=$?
if [ $res -ne 0 ];then
  echo "Error peer channel update faile !!!" >&2
  exit 1
fi

exit 0