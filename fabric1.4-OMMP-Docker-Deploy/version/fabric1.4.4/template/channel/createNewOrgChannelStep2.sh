#!/bin/bash
## 总存放目录 相关文件脚本存放文件夹名称
TARGET=$[target]
NEW_ORG_CLI_NAME=$[newOrgCliName]

## 检测容器
cExist=`docker inspect --format '{{.State.Running}}' ${NEW_ORG_CLI_NAME}`
if [ "${cExist}" != "true" ];then
   echo "Error docker container ${NEW_ORG_CLI_NAME} is not exist" > &2
   exit 1
fi 

##拷贝证书到cli
docker cp ./${TARGET}  ${NEW_ORG_CLI_NAME}:/opt/gopath/src/github.com/hyperledger/fabric/peer/  
sleep 5
docker cp ./${TARGET}/crypto-config/peerOrganizations  ${NEW_ORG_CLI_NAME}:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/
sleep 3

##执行容器内脚本
docker exec ${NEW_ORG_CLI_NAME} sh -c "cd ${TARGET} && chmod +x createNewOrgChannelStep3.sh && ./createNewOrgChannelStep3.sh"
ret=$?
if [ $ret -ne 0 ];then
   echo "Error exec script createNewOrgChannelStep3 error!!!!!!" >&2
   exit 1
fi
exit 0