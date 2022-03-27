#!/bin/bash

CC_CONTAINER_NAME=$[ccContainerName]
PEER_CONTAINER_NAME=$[peerContainerName]
CC_NAME=$[ccName]
CC_PATH="/var/hyperledger/production/chaincodes"


verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "Error !!! ERROR  FAILE : $2" >&2
    exit 1
  fi
}


##进入镜像执行命令
PEERNUM=${#PEER_CONTAINER_NAME[@]}
for ((i=0;i<PEERNUM;i++))
do
  ##清理合约镜像
  ##检测容器
  cExist=`docker inspect --format '{{.State.Running}}'  ${CC_CONTAINER_NAME[$i]}`
  if [ "${cExist}" != "true" ];then
    echo "Error docker container ${ORG_CLI_NAME} is not exist" >&2
  else
        docker rm -f ${CC_CONTAINER_NAME[$i]}
  fi 
  
  set -x
  docker exec ${PEER_CONTAINER_NAME[$i]} sh -c "cd ${CC_PATH} && mv ${CC_NAME} ${CC_NAME}.bak"
  res=$?
  verifyResult $res "peer disable chaincode  ${CC_NAME} failed"
  set +x
done

exit 0