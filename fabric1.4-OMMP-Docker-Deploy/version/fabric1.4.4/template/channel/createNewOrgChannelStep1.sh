#!/bin/bash
## 总存放目录 相关文件脚本存放文件夹名称
TARGET=$[target]
NEW_ORG_CLI_NAME=$[newOrgCliName]
NEW_ORG_CA_NAME=$[newOrgCAName]
## peer节点host
ORG_NAME=$[orgMSPName]
PEER_NAME=$[peerName]
#是否启动ca
START_CA=$[startCA]
echo "========================================="
echo "启动新组织节点"
echo "========================================="
cd ${TARGET}

if [ ${START_CA} -eq 1 ];then
   COMPOSE_PROJECT_NAME=baiyun docker-compose -f  base.yaml up -d ${NEW_ORG_CA_NAME}
   ret=$?
	if [ $ret -ne 0 ];then
		echo "Error start container CA:${NEW_ORG_CA_NAME} error!!!!!!" >&2
		exit 1
	fi
   sleep 3
fi

 COMPOSE_PROJECT_NAME=baiyun docker-compose -f base.yaml up -d ${PEER_NAME}
 ret=$?
 if [ $ret -ne 0 ];then
	echo "Error start container Peer:${PEER_NAME} error!!!!!!" >&2
	exit 1
 fi
 sleep 3
 COMPOSE_PROJECT_NAME=baiyun docker-compose -f base.yaml up -d ${NEW_ORG_CLI_NAME}
 ret=$?
 if [ $ret -ne 0 ];then
	echo "Error start container cli:${NEW_ORG_CLI_NAME} error!!!!!!" >&2
	exit 1
 fi
 
cd -
exit 0