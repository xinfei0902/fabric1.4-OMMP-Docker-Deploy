#!/bin/bash

CHANNEL_NAME=$[channelName]
## orderer名称加端口 例如：orderer0.example.com:7050
ORDERER_NAME=$[ordererAddress]
ORDERER_CA=$[ordererTlsCa]
##合约配置
INSTALL_CC_NAME=$[installCCName]
##链码版本 与上面链码名称一一对应
INSTALL_CC_VERSION=$[installCCVersion]
#INSTALL_CC_INIT='$[installCCInit]'
INSTALL_CC_POLICY="$[installCCPolicy]"
INSTALL_CC_PATH=$[installCCPath]
##节点配置
ORG_NAMEMSP=$[orgMSPID]
ORG_PEER_NUM=$[orgPeerNum]
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

echo "===================================="
echo "hosts"
echo "===================================="

peerTN=${#PEER_HOST[@]}
echo "peerTN=$peerTN"
for ((i=0;i<peerTN;i++))
do
  echo ${PEER_HOST[$i]} >>/etc/hosts
done

echo "===================================="
echo "安装"
echo "===================================="
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
		peer chaincode install -n ${INSTALL_CC_NAME} -v ${INSTALL_CC_VERSION} -p ${INSTALL_CC_PATH}
		res=$?
		set +x
		verifyResult $res "peer install chaincode  ${INSTALL_CC_NAME[$i]} failed"
    }
 }
 
sleep 3

echo "===================================="
echo "实例化"
echo "===================================="
export CORE_PEER_LOCALMSPID=${ORG_NAMEMSP[0]}
export CORE_PEER_TLS_ROOTCERT_FILE=${PEER_TLS_PATH[0]}
export CORE_PEER_MSPCONFIGPATH=${PEER_MSP_PATH[0]}
export CORE_PEER_ADDRESS=${PEER_ADDRESS[0]}
set -x
peer chaincode upgrade -o $ORDERER_NAME --tls true --cafile $ORDERER_CA -C $CHANNEL_NAME -n ${INSTALL_CC_NAME} -v ${INSTALL_CC_VERSION} -c '{"Args":["init","a", "100", "b","200"]}' -P ${INSTALL_CC_POLICY}
ret=$?
set +x  
verifyResult $ret "peer instantiate chaincode  ${INSTALL_CC_NAME[$i]} failed"	

exit 0