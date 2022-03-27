#!/bin/bash

CHANNEL_NAME=$[channelName]
## orderer名称加端口 例如：orderer0.example.com:7050
ORDERER_NAME=$[ordererAddress]
ORDERER_CA=$[ordererTlsCa]
##节点配置
ORG_NAMEMSP=$[orgMSPID]
PEER_TLS_PATH=$[peerTlsCa]
PEER_MSP_PATH=$[peerUsersMsp]
PEER_ADDRESS=$[peerHosts]

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

set -x
peer channel fetch 0 $CHANNEL_NAME.block -o ${ORDERER_NAME} -c $CHANNEL_NAME --tls --cafile $ORDERER_CA >&log.txt
res=$?
set +x

verifyResult $res "Fetching config block from orderer has Failed"
sleep 3
set -x
peer channel join -b $CHANNEL_NAME.block >&log.txt
res=$?
set +x

verifyResult $res "peer join channel ${CHANNEL_NAME} failed"
sleep 3

exit 0