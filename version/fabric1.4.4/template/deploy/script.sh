#!D:\Tools-Install-Package\git\Git\bin/bash

##orderer环境变量
ORDERER_ADDRESS=$[ordererAddress]
ORDERER_TLSCA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/$[ordererTlsCa]
##通道数组
CHANNEL_NAME=$[channelNameArray]
ORG_NAME=$[orgNameArray]
PEER_CONFIG=$[peerConfigArray]
# 每一个通道对应一个" "需要安装所有合约
CHAINCODE_NAME=$[ccNameArray]
# 针对每一个合约需要参与安装的组织
CHAINCODE_ORG=$[ccOrgArray]
## 每一个合约安装实例化信息
CHAINCODE_INSTALL_CONFIG=$[ccInstallConfigArray]

verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
    echo "========= ERROR !!! FAILED to execute End-2-End Scenario ==========="
    echo
    exit 1
  fi
}

## 创建通道
channelNum=${#CHANNEL_NAME[@]}
for ((i=0;i<channelNum;i++))
{   
    channelName=${CHANNEL_NAME[$i]}
    orgArray=("${ORG_NAME[$i]}")
    for orgName in ${orgArray[@]}; do
      for peerArray in "${PEER_CONFIG[@]}"; do
        peerConfig=($peerArray)
        if [ "$orgName" == "${peerConfig[0]}" ];then
            CORE_PEER_LOCALMSPID=${peerConfig[1]}
            CORE_PEER_TLS_ROOTCERT_FILE=${peerConfig[2]}
            CORE_PEER_MSPCONFIGPATH=${peerConfig[3]}
            CORE_PEER_ADDRESS=${peerConfig[4]}
            set -x
            peer channel create -o $ORDERER_ADDRESS -c $channelName -f ./channel-artifacts/$channelName.tx --tls true --cafile $ORDERER_TLSCA >&log.txt
            res=$?
            set +x
            verifyResult $res "Channel creation failed"
            sleep 1
        fi   
      done
    done
}

## 所有节点加入通道
channelNum=${#CHANNEL_NAME[@]}
for ((i=0;i<channelNum;i++))
{
  channelName=${CHANNEL_NAME[$i]}
  orgArray=(${ORG_NAME[$i]})
  for orgName in "${orgArray[@]}"; do
    for peerArray in "${PEER_CONFIG[@]}"; do
          peerConfig=($peerArray)
          if [ "$orgName" == "${peerConfig[0]}" ];then
            peerCN=${#peerConfig[@]}
            peerVarialN=$((($peerCN-2)/3))
            CORE_PEER_LOCALMSPID=${peerConfig[1]}
            for ((j=0;j<peerVarialN;j++))
            {
              CORE_PEER_TLS_ROOTCERT_FILE=${peerConfig[2+$j*3]}
              CORE_PEER_MSPCONFIGPATH=${peerConfig[3+$j*3]}
              CORE_PEER_ADDRESS=${peerConfig[4+$j*3]}
              set -x
              peer channel join -b $CHANNEL_NAME.block >&log.txt
              res=$?
              set +x
              verifyResult $res "peer=${CORE_PEER_ADDRESS} join channel=${channelName} failed"
              sleep 1
            }
          fi
    done
  done  
}

## 设置锚节点
channelNum=${#CHANNEL_NAME[@]}
for ((i=0;i<channelNum;i++))
{   
    channelName=${CHANNEL_NAME[$i]}
    orgArray=("${ORG_NAME[$i]}")
    for orgName in ${orgArray[@]}; do
      for peerArray in "${PEER_CONFIG[@]}"; do
        peerConfig=($peerArray)
        if [ "$orgName" == "${peerConfig[0]}" ];then
            CORE_PEER_LOCALMSPID=${peerConfig[1]}
            CORE_PEER_TLS_ROOTCERT_FILE=${peerConfig[2]}
            CORE_PEER_MSPCONFIGPATH=${peerConfig[3]}
            CORE_PEER_ADDRESS=${peerConfig[4]}
            set -x
            peer channel update -o $ORDERER_ADDRESS -c $channelName -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}Panchors@$channelName.tx --tls true --cafile $ORDERER_TLSCA >&log.txt
            res=$?
            set +x
            verifyResult $res "Channel creation failed"
            sleep 1
        fi   
      done
    done
}

## 安装合约
channelNum=${#CHANNEL_NAME[@]}
ccNum=0
for ((i=0;i<channelNum;i++))
{
  channelName=${CHANNEL_NAME[$i]}
  ccArray=("${CHAINCODE_NAME[$i]}")
  echo "ccArray" $ccArray
  for ccName in ${ccArray[@]}; do
      echo "合约名称" $ccName
      ccConfig=(${CHAINCODE_INSTALL_CONFIG[$i]})
      echo "ccConfig" "$ccConfig"
      ccCN=${#ccConfig[@]}
      echo "ccCN" $ccCN
      ccVarialN=$(($ccCN/3))
      for ((j=0;j<ccVarialN;j++))
      {
        if [ "$ccName" == "${ccConfig[$j*3]}" ];then
              CC_NAME=${ccConfig[$j*3]}
              CC_VERSION=${ccConfig[1+$j*3]}
              CC_POLICY=${ccConfig[2+$j*3]}
              echo "1" $CC_NAME
              echo "2" $CC_VERSION
              echo "3" $CC_POLICY
              ## 设置peer环境变量 所有peer安装合约
              orgArray=("${CHAINCODE_ORG[$ccNum]}")
              for orgName in ${orgArray[@]}; do
                  for peerArray in "${PEER_CONFIG[@]}"; do
                        peerConfig=($peerArray)
                        if [ "$orgName" == "${peerConfig[0]}" ];then
                          peerCN=${#peerConfig[@]}
                          peerVarialN=$((($peerCN-2)/3))
                          CORE_PEER_LOCALMSPID=${peerConfig[1]}
                          for ((j=0;j<peerVarialN;j++))
                          {
                            CORE_PEER_TLS_ROOTCERT_FILE=${peerConfig[2+$j*3]}
                            CORE_PEER_MSPCONFIGPATH=${peerConfig[3+$j*3]}
                            CORE_PEER_ADDRESS=${peerConfig[4+$j*3]}
                            echo "11" $CORE_PEER_LOCALMSPID
                            echo "11" $CORE_PEER_TLS_ROOTCERT_FILE
                            echo "11" $CORE_PEER_MSPCONFIGPATH
                            echo "11" $CORE_PEER_ADDRESS 
                            set -x
                            peer chaincode install -n ${CC_NAME} -v ${CC_VERSION} -l ${LANGUAGE} -p ${CC_PATH} >&log.txt
                            res=$?
                            set +x
                            verifyResult $res "peer=${CORE_PEER_ADDRESS} install chaincode=${} failed"
                            sleep 1
                          }
                        fi
                  done
              done
              echo "instantiate" $CC_NAME
              echo "instantiate" $CC_VERSION
              echo "instantiate" $CC_POLICY
              set -x
              peer chaincode instantiate -o $ORDERER_ADDRESS --tls true --cafile $ORDERER_TLSCA -C $channelName -n $CC_NAME -l "golang" -v $CC_VERSION -c '{"Args":["init"]}' -P "${CC_POLICY}" >&log.txt
              res=$?
              set +x
              verifyResult $res "Chaincode instantiation  channel '$channelName' failed"  
        fi    
      }
    ccNum=$(($ccNum+1)) 
  done  
}