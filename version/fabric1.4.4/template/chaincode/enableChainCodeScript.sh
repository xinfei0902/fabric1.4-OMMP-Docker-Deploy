#!/bin/bash

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
  set -x
  docker exec ${PEER_CONTAINER_NAME[$i]} sh -c "cd ${CC_PATH} && mv ${CC_NAME}.bak ${CC_NAME}"
  res=$?
  verifyResult $res "peer disable chaincode  ${CC_NAME} failed"
  set +x
done

exit 0