#!/bin/bash
## 总存放目录 相关文件脚本存放文件夹名称
TARGET=$[target]
SOURCE_ORG_CLI_NAME=$[sourceOrgCliName]

## 检测容器
cExist=`docker inspect --format '{{.State.Running}}' ${SOURCE_ORG_CLI_NAME}`
if [ "${cExist}" != "true" ];then
   echo "Error docker container ${SOURCE_ORG_CLI_NAME} is not exist" >&2
   exit 1
fi 

docker cp ./${TARGET}  ${SOURCE_ORG_CLI_NAME}:/opt/gopath/src/github.com/hyperledger/fabric/peer/  
sleep 5
docker exec ${SOURCE_ORG_CLI_NAME} sh -c "cd ${TARGET} && chmod +x createChannelStep2.sh && ./createChannelStep2.sh"
ret=$?
if [ $? -ne 0 ];then
   echo "Error exec script createChannelStep2 error!!!!!!" >&2
   exit 1
fi
sleep 3
exit 0