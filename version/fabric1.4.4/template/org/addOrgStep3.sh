#!/bin/bash
## 总存放目录 相关文件脚本存放文件夹名称
TARGET=$[target]
##CLI_NAME=$[cliName]
NEW_ORG_CLI_NAME=$[newOrgCliName]
CA_CONTINUE_NAME=$[caContinueName]
## peer节点host
PEER_NAME=$[peerName]


##新组织启动
echo "========================================="
echo "启动新组织节点"
echo "========================================="
cd ${TARGET}
## COMPOSE_PROJECT_NAME=baiyun 是docker网络名称 是启动容器加入baiyun网络中
COMPOSE_PROJECT_NAME=baiyun docker-compose -f base.yaml up -d $PEER_NAME
sleep 3
COMPOSE_PROJECT_NAME=baiyun docker-compose -f base.yaml up -d $CA_CONTINUE_NAME
sleep 3
COMPOSE_PROJECT_NAME=baiyun docker-compose -f base.yaml up -d ${NEW_ORG_CLI_NAME}
sleep 3
cd -

##检测容器是否启动成功
pExist=`docker inspect --format '{{.State.Running}}' ${PEER_NAME}`
if [ "${pExist}" != "true" ];then
   echo "Error docker container ${PEER_NAME} is not exist" >&2
   exit 1
fi 

cExist=`docker inspect --format '{{.State.Running}}' ${NEW_ORG_CLI_NAME}`
if [ "${cExist}" != "true" ];then
   echo "Error docker container ${NEW_ORG_CLI_NAME} is not exist" >&2
   exit 1
fi 

##新组织加入通道
docker cp ./${TARGET}  ${NEW_ORG_CLI_NAME}:/opt/gopath/src/github.com/hyperledger/fabric/peer/
sleep 5
docker cp ./${TARGET}/crypto-config/peerOrganizations  ${NEW_ORG_CLI_NAME}:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/
sleep 3
docker exec ${NEW_ORG_CLI_NAME} sh -c "cd ${TARGET} && chmod +x addOrgStep4.sh && ./addOrgStep4.sh"
ret=$?
if [ $ret -ne 0 ];then
   echo "Error exec script addOrgStep4 error!!!!!!" >&2
   exit 1
fi

exit 0