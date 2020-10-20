#!/bin/bash
# 文件目录
TARGET=$[target]
# 新增组织 MSPID名称
ORG_NAMEMSP=$[orgMSPID]
##通道名称
CHANNEL_NAME=$[channelName]

##节点配置
ORG_NAMEMSP=$[orgMSPID]
PEER_TLS_PATH=$[peerTlsCa]
PEER_MSP_PATH=$[peerUsersMsp]
PEER_ADDRESS=$[peerHosts]
#是否开启tls模式
CORE_PEER_TLS_ENABLED=true
#合约配置
##INSTALL_CC_NAME 包含所有需要节点安装所有链码 是数组形式 例如INSTALL_CC_NAME=(bysms bypms) 
INSTALL_CC_NAME=$[installCCName]
##链码版本 与上面链码名称一一对应
INSTALL_CC_VERSION=$[installCCVersion]
INSTALL_CC_PATH=$[installCCPath]

verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "Error !!! ERROR  FAILE : $2" >&2
    exit 1
  fi
}

export CORE_PEER_LOCALMSPID=${ORG_NAMEMSP}
export CORE_PEER_TLS_ROOTCERT_FILE=${PEER_TLS_PATH}
export CORE_PEER_MSPCONFIGPATH=${PEER_MSP_PATH}
export CORE_PEER_ADDRESS=${PEER_ADDRESS}

##安装
echo "===================================="
echo "安装"
echo "===================================="
num=${#INSTALL_CC_NAME[@]}
echo "========num=${num}========="
for ((i=0;i<num;i++))
{
	set -x
    peer chaincode install -n ${INSTALL_CC_NAME[$i]} -v ${INSTALL_CC_VERSION[$i]} -p ${INSTALL_CC_PATH[$i]}
	res=$?
	set +x
	verifyResult $res "peer install chaincode  ${INSTALL_CC_NAME[$i]} failed"
}
exit 0

