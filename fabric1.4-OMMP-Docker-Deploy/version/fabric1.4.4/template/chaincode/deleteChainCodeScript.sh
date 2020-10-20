#!/bin/bash
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
  docker rm -f ${CC_CONTAINER_NAME[$i]}
  set -x
  docker exec ${PEER_CONTAINER_NAME[$i]} sh -c "cd ${CC_PATH} && rm -f ${CC_NAME} && rm -f ${CC_NAME}.bak"
  res=$?
  verifyResult $res "peer delete chaincode  ${CC_NAME} failed"
  set +x
done

exit 0