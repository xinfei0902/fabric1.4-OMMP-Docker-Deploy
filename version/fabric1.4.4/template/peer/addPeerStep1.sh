#!/bin/bash
## 总存放目录 相关文件脚本存放文件夹名称
TARGET=$[target]
NEW_ORG_CLI_NAME=$[newOrgCliName]
## peer节点host
PEER_NAME=$[peerName]

## 是为了考虑新增节点时可以同时安装组织所有链码
IS_EXEC_INSTALL_CHAINCODE=$[isExecInstallCC]



##新节点启动
echo "========================================="
echo "启动新组织节点"
echo "========================================="
cd ${TARGET}
COMPOSE_PROJECT_NAME=baiyun docker-compose -f base.yaml up -d $PEER_NAME
sleep 3

COMPOSE_PROJECT_NAME=baiyun docker-compose -f base.yaml up -d ${NEW_ORG_CLI_NAME}
sleep 3
cd -

##检测容器是否存在
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



##新节点加入通道
docker cp ./${TARGET}  ${NEW_ORG_CLI_NAME}:/opt/gopath/src/github.com/hyperledger/fabric/peer/
sleep 5
docker cp ./${TARGET}/crypto-config/peerOrganizations  ${NEW_ORG_CLI_NAME}:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/
sleep 3
docker exec ${NEW_ORG_CLI_NAME} sh -c "cd ${TARGET} && chmod +x addPeerStep2.sh && ./addPeerStep2.sh"
ret=$?
if [ $ret -ne 0 ];then
   echo "Error exec script addPeerStep2 error!!!!!!" >&2
   exit 1
fi

## 安装组织内参入背书所有链码
if [ ${IS_EXEC_INSTALL_CHAINCODE} = "1" ];then
	
	docker exec ${NEW_ORG_CLI_NAME} sh -c "cd ${TARGET} && chmod +x addPeerStep3.sh && ./addPeerStep3.sh"
    ret=$?
	if [ $ret -ne 0 ];then
	   echo "Error exec script  addPeerStep3 error!!!!!!" >&2
	   exit 1
	fi
	
fi
exit 0