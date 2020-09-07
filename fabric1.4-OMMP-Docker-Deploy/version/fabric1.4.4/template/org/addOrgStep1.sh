#!/bin/bash
## 总存放目录 相关文件脚本存放文件夹名称
TARGET=$[target]
CLI_NAME=$[cliName]
NEW_ORG_CLI_NAME=$[newOrgCliName]
CA_CONTINUE_NAME=$[caContinueName]
## peer节点host
PEER_NAME=$[peerName]

##此变量此版本目前不考虑 是为了考虑新增组织时可以同时升级链码更新背书策略
IS_EXEC_INSTALL_CHAINCODE=1

##检测容器是否运行
cExist=`docker inspect --format '{{.State.Running}}' ${CLI_NAME}`
if [ "${cExist}" != "true" ];then
   echo "Error docker container ${CLI_NAME} is not exist" >&2
   exit 1
fi 


##执行相应文件拷贝到cli容器
docker cp ./${TARGET}  ${CLI_NAME}:/opt/gopath/src/github.com/hyperledger/fabric/peer/
sleep 5
docker cp ./${TARGET}/crypto-config/peerOrganizations  ${CLI_NAME}:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/
sleep 3

docker exec ${CLI_NAME} sh -c "cd ${TARGET} && chmod +x addOrgStep2.sh && ./addOrgStep2.sh"
ret=$?
echo "exec script Step1 end"
if [ $ret -ne 0 ];then
   echo "Error exec new add Org script addOrgStep2.sh error!!!!!!" >&2
   exit 1
fi

exit 0