#!/bin/bash
## 总存放目录 相关文件脚本存放文件夹名称
TARGET=$[target]
ORG_CLI_NAME=$[cliName]
CC_NAME=$[installCCName]

##检测容器
cExist=`docker inspect --format '{{.State.Running}}' ${ORG_CLI_NAME}`
if [ "${cExist}" != "true" ];then
   echo "Error docker container ${ORG_CLI_NAME} is not exist" >&2
   exit 1
fi 

##安装链码
docker cp ./${TARGET}  ${ORG_CLI_NAME}:/opt/gopath/src/github.com/hyperledger/fabric/peer/
sleep 5
docker cp ./${TARGET}/chaincode/${CC_NAME}/  ${ORG_CLI_NAME}:/opt/gopath/src/github.com/chaincode
sleep 5
docker cp ./${TARGET}/crypto-config/peerOrganizations  ${ORG_CLI_NAME}:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/
sleep 3
docker exec ${ORG_CLI_NAME} sh -c "cd ${TARGET} && chmod +x addChainCodeStep2.sh && ./addChainCodeStep2.sh"
ret=$?
if [ $ret -ne 0 ];then
   echo "Error exec script addChainCodeStep2 error!!!!!!" >&2
   exit 1
fi
exit 0