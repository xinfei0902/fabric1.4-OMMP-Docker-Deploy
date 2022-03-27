#!/bin/bash
## 总存放目录 相关文件脚本存放文件夹名称
TARGET=$[target]
CLI_NAME=$[cliName]

##检测容器是否运行
cExist=`docker inspect --format '{{.State.Running}}' ${CLI_NAME}`
if [ "${cExist}" != "true" ];then
   echo "Error docker container ${CLI_NAME} is not exist" >&2
   exit 1
fi 


##执行相应文件拷贝到cli容器
docker cp ./ ${CLI_NAME}:/opt/gopath/src/github.com/hyperledger/fabric/peer/
sleep 5
docker cp ./crypto-config/peerOrganizations  ${CLI_NAME}:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/
sleep 3

docker exec ${CLI_NAME} sh -c "cd scripts && chmod +x script.sh && ./script.sh"
ret=$?
echo "exec script Step1 end"
if [ $ret -ne 0 ];then
   echo "Error exec ./script.sh error!!!!!!" >&2
   exit 1
fi

exit 0