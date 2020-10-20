#!/bin/bash

CHANNEL_NAME=$[channelName]
## orderer名称加端口 例如：orderer0.example.com:7050
ORDERER_NAME=$[ordererAddress]
ORDERER_CA=$[ordererTlsCa]
##节点配置
ORG_NAMEMSP=$[orgMSPID]
PEER_TLS_PATH=$[peerTlsCa]
PEER_MSP_PATH=$[peerUsersMsp]
PEER_ADDRESS=$[peerAddress]
##示例（"192.168.8.1 peer0.baiyun.example.com" "192.168.8.2 peer2.baiyun.example.com"）
PEER_HOST=$[peerHosts]
verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "Error !!! ERROR  FAILE : $2" >&2
    exit 1
  fi
}

##开始执行
export CORE_PEER_LOCALMSPID=${ORG_NAMEMSP}
export CORE_PEER_TLS_ROOTCERT_FILE=${PEER_TLS_PATH}
export CORE_PEER_MSPCONFIGPATH=${PEER_MSP_PATH}
export CORE_PEER_ADDRESS=${PEER_ADDRESS}

echo "===================================="
echo "hosts"
echo "===================================="

peerTN=${#PEER_HOST[@]}
echo "peerTN=$peerTN"
for ((i=0;i<peerTN;i++))
do
  echo ${PEER_HOST[$i]} >>/etc/hosts
done

set -x
peer channel fetch 0 $CHANNEL_NAME.block -o ${ORDERER_NAME} -c $CHANNEL_NAME --tls --cafile $ORDERER_CA >&log.txt
ret=$?
set +x

verifyResult $ret "Fetching config block from orderer has Failed"
sleep 3
set -x
peer channel join -b $CHANNEL_NAME.block >&log.txt
ret=$?
set +x
verifyResult $ret "peer join channel ${CHANNEL_NAME} failed"

exit 0