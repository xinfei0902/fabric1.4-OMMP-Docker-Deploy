#!/bin/bash

CHANNEL_NAME=$[channelName]
## orderer名称加端口 例如：orderer0.example.com:7050
ORDERER_NAME=$[ordererAddress]
##往cli容器添加hosts
ORDERER_HOSTS=$[ordererHosts]
#ORDERER_TLS_PATH=$[ordererTlsCa]
#ORDERER_MSP_PATH=$[ordererUsersMsp]
ORDERER_CA=$[ordererTlsCa]
##节点配置
ORG_NAMEMSP=$[orgMSPID]
PEER_TLS_PATH=$[peerTlsCa]
PEER_MSP_PATH=$[peerUsersMsp]
PEER_ADDRESS=$[peerAddress]
PEER_HOST=$[peerHosts]
verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "Error !!! ERROR  FAILE : $2" >&2
    exit 1
  fi
}

##hosts
echo "$ORDERER_HOSTS" >>/etc/hosts

echo "===================================="
echo "hosts"
echo "===================================="

peerTN=${#PEER_HOST[@]}
echo "peerTN=$peerTN"
for ((i=0;i<peerTN;i++))
do
  echo ${PEER_HOST[$i]} >>/etc/hosts
done


##开始执行
export CORE_PEER_LOCALMSPID=${ORG_NAMEMSP}
export CORE_PEER_TLS_ROOTCERT_FILE=${PEER_TLS_PATH}
export CORE_PEER_MSPCONFIGPATH=${PEER_MSP_PATH}
export CORE_PEER_ADDRESS=${PEER_ADDRESS}

set -x
peer channel fetch 0 $CHANNEL_NAME.block -o ${ORDERER_NAME} -c $CHANNEL_NAME --tls --cafile $ORDERER_CA >>peerFetch1.log 2>&1
ret=$?
set +x
verifyResult $ret "Fetching config block from orderer has Failed"
sleep 3

result=$(cat peerFetch1.log)
expr index "$result" "Error"
res=$?
if [ $res -ne 0 ];then
  echo "Error peer channel fetch faile !!!" >&2
  exit 1
fi

set -x
peer channel join -b $CHANNEL_NAME.block 
ret=$?
set +x

verifyResult $ret "peer join channel ${CHANNEL_NAME} failed"
sleep 2

exit 0