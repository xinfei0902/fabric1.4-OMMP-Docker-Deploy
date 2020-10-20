#!/bin/bash

#脚本后台运行命令
#nohup ./peer-status-monitor-baiyun-admin.sh >/dev/null 2>&1 &

#CHECK_PORT 是要检测的端口 可以检测多个端口 例如CHECK_PORT="7050 7051"
CHECK_PORT="$[checkPeerPort]"
#ADDRESS 是调用上报的ip 提供检测报告
ADDRESS=172.16.50.203
#PPORT 是调用上报信息的端口
PORT=9050
#URL 调用地址
URL=/webapi/blockchainSync/updateNode
#INIT_STATUS 启动peer初始状态 0=停止 1=运行 2=删除
INIT_STATUS=2

#CHECK_SLEEP  设定多长时间检测一次
CHECK_SLEEP=5

######################################################################
# 组装 参数说明                                                      #
# token: 调用地址标识                                                #
# name: 节点名称                                                     #
# accessKey sdk中节点accesskey的key值                                #
# id: 节点id                                                         #
# ip: 节点所在机器外网IP                                             #
# port: 节点所在机器的端口                                           #
# cityId：所在城市id                                                 #
# localion: 所在地址的经纬度                                         #
# status: 节点状态 0=停止 1=运行                                     #
######################################################################

#发送application/json 格式组装数据
#peerSuccessInfo='{"name":"peer0.orgrw.com","ip":"59.110.150.162","port":"7051","status":"1","cityId":"001","location":"201-118-33","deleted":"0"}'
#peerErrorInfo='{"name":"peer0.orgrw.com","ip":"59.110.150.162","port":"7051","status":"0","cityId":"001","location":"201-118-33","deleted":"0"}'

#发送application/x-www-form-urlencoded格式组装数据
peerDeleteInfo='token=fb3f261e82a744249d8ef3a7173dd32d&name=$[peerNick]&accessKey=$[accessKey]&id=$[peerID]&ip=$[peerIP]&port=$[peerID]&cityId=101260102&localion=106.65°E,26.68°N&status=3'
peerErrorInfo='token=fb3f261e82a744249d8ef3a7173dd32d&name=$[peerNick]&accessKey=$[accessKey]&id=$[peerID]&ip=$[peerIP]&port=$[peerID]&cityId=101260102&localion=106.65°E,26.68°N&status=0'

putCheckDeleteMessage(){
 url="http://${ADDRESS}:${PORT}${URL}"
 #code=$(curl -H "Content-Type:application/json" -X POST --data ${peerSuccessInfo} $url |jq '.code')
 code=$(curl -X POST --data ${peerDeleteInfo}  $url |jq '.code')
 if [ $code -eq 200 ]; then
    echo success $code
    exit 0
 fi
}


putCheckErrorMessage(){
 url="http://${ADDRESS}:${PORT}${URL}"
 #code=$(curl -H "Content-Type:application/json" -X POST --data ${peerErrorInfo} $url |jq '.code')
  code=$(curl -X POST --data ${peerErrorInfo}  $url |jq '.code')
  if [ $code -eq 200 ]; then
     echo success $code
  fi
}

## 安装工具
which jq
ret=$?
if [ $ret -ne 0 ];then
 echo "install jq"
 apt-get -y update && apt-get -y install jq
fi

while :
 do
  for i in $CHECK_PORT
   do  
    netstat -anput | grep "LISTEN" |grep $i 1>/dev/null
    if [ $? -ne 0 ]; then 
      putCheckDeleteMessage
      sleep ${CHECK_SLEEP}
    else
      putCheckErrorMessage
     sleep ${CHECK_SLEEP}
    fi
  done
done  

