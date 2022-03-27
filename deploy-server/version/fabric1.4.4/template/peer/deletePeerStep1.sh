#!/bin/bash
## 脚本主要功能 删除节点容器

## 节点容器名称
PEER_CONTAINER_ARRAY=$[peerContainerArray]
## 安装在组织链码镜像
DEL_CC_IMAGES=$[deleteCCImages]

##组织节点停掉 节点cli
echo "========================================="
echo "停掉节点 cli"
echo "========================================="

peerContainerNum=${#PEER_CONTAINER_ARRAY[@]}
for ((i=0;i<peerContainerNum;i++))
{
    docker stop ${PEER_CONTAINER_ARRAY[$i]}##检测容器是否停用成功
    sleep 2
    pExist=`docker inspect --format '{{.State.Running}}' ${PEER_CONTAINER_ARRAY[$i]}`
    if [ "${pExist}" != "false" ];then
        echo "Error docker container ${PEER_CONTAINER_ARRAY[$i]} stop fail" >&2
        exit 1
    fi 
}
sleep 3


## 删除所有退出状态容器
docker rm $(docker ps --all -q -f status=exited)
sleep 3

## 删除链码镜像

ccImages=${#DEL_CC_IMAGES[@]}
for ((j=0;j<ccImages;j++))
{
  result=$(docker images ${DEL_CC_IMAGES[$j]})
  var=${result// /$'\n'}
  for str in $var;
  do
    if [[ $str == ${DEL_IMAGES} ]];then
      docker rmi $str
    fi
  done
  sleep 2
}
exit 0