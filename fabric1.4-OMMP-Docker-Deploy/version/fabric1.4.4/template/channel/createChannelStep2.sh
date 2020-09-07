#!/bin/bash

CHANNEL_NAME=$[channelName]
## orderer名称加端口 例如：orderer0.prisons.guizhou:7050
ORDERER_NAME=$[ordererAddress]
ORDERER_CA=$[ordererTlsCa]

ORG_NAMEMSP=$[orgMSPID]

CHANNEL_CONFIG=./channel-artifacts/${CHANNEL_NAME}.tx
ANCHORS_CONFIG=./channel-artifacts/${ORG_NAMEMSP}anchors.tx.tx
##节点配置

PEER_TLS_PATH=$[peerTlsCa]
PEER_MSP_PATH=$[peerUsersMsp]
PEER_ADDRESS=$[peerAddress]

verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "Error !!! ERROR  FAILE : $2" >&2
    exit 1
  fi
}
## 创建通道
set -x
export CORE_PEER_LOCALMSPID=$ORG_NAMEMSP
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER_TLS_PATH
export CORE_PEER_MSPCONFIGPATH=$PEER_MSP_PATH
export CORE_PEER_ADDRESS=$PEER_ADDRESS
sleep 2

peer channel create -o $ORDERER_NAME -c $CHANNEL_NAME -f $CHANNEL_CONFIG --tls true --cafile $ORDERER_CA
set +x
## 所有节点加入通道

export CORE_PEER_LOCALMSPID=$ORG_NAMEMSP
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER_TLS_PATH
export CORE_PEER_MSPCONFIGPATH=$PEER_MSP_PATH
export CORE_PEER_ADDRESS=$PEER_ADDRESS
set -x
peer channel join -b $CHANNEL_NAME.block
res=$?
set +x
verifyResult $res "peer join channel ${CHANNEL_NAME} failed"


##定义锚节点

export CORE_PEER_LOCALMSPID=$ORG_NAMEMSP
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER_TLS_PATH
export CORE_PEER_MSPCONFIGPATH=$PEER_MSP_PATH
export CORE_PEER_ADDRESS=$PEER_ADDRESS
set -x
peer channel update -o $ORDERER_NAME -c $CHANNEL_NAME -f $ANCHORS_CONFIG --tls true --cafile $ORDERER_CA 
res=$?
set +x
verifyResult $res "peer update anchors  failed"

exit 0