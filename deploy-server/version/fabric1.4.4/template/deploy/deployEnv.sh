#!/bin/bash

## 安装curl
curl --version &> /dev/null
if [ $? -ne 0 ]; then
   echo " != 0"
   if [ -f /etc/redhat-release ]; then
      echo "centos"
       yum update
       yum instll -y curl
   fi
   if [ -f /etc/lsb-release ]; then
        echo "ubunt"
        sudo apt-get update
        sudo apt install curl
   fi
fi
echo " curl 已安装"


## 安装git
git version &> /dev/null
if [ $? -ne 0 ]; then
   echo " != 0"
   if [ -f /etc/redhat-release ]; then
      echo "centos"
       yum instll -y git
   fi
   if [ -f /etc/lsb-release ]; then
        echo "ubunt"
        sudo apt install -y git
   fi
fi
echo "git 已安装"

## 安装golang 
go version &> /dev/null
if [ $? -ne 0 ]; then
    echo "开始安装golang"
    tar -C /usr/local -zxf  go1.15.13.linux-amd64.tar.gz
    echo "export GOROOT=/usr/local/go" >>/etc/profile
    echo "export GOPATH=$[DEPLOY_TOOLS_PATH]/gopath" >>/etc/profile
    echo 'export PATH=$PATH:$GOROOT/bin' >>/etc/profile
    source /etc/profile
    go version &> /dev/null
    if [ $? -ne 0 ]; then
       echo "安装golang失败"
       exit 1
    fi
    go env -w GO111MODULE="on"
    go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/
    mkdir -p $[DEPLOY_TOOLS_PATH]/gopath/src
    mkdir -p $[DEPLOY_TOOLS_PATH]/gopath/pkg
    mkdir -p $[DEPLOY_TOOLS_PATH]/gopath/bin
    echo "安装golang成功"
else
    echo "golang 环境已经安装"   
fi

source /etc/profile

## 安装docker
docker version &> /dev/null
if [ $? -ne 0 ]; then
 if [ -f /etc/redhat-release ]; then
    echo "centos install docker"
    yum update -y
    yum install -y yum-utils device-mapper-persistent-data lvm2
    yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
    yum install -y docker-ce-18.03.1.ce
    systemctl start docker
    docker version &> /dev/null
    if [ $? -ne 0 ]; then
        echo "安装docker失败"
        exit 1
    fi
    echo "docker 安装成功"
 fi
 if [ -f /etc/lsb-release ]; then
    echo "ubuntu install docker"
    sudo apt-get update
    sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
    sudo apt-get update
    sudo apt-get install docker-ce=18.03.1~ce-0~ubuntu
    sudo systemctl start docker
    docker version &> /dev/null
    if [ $? -ne 0 ]; then
        echo "安装docker失败"
        exit 1
    fi
    echo "docker 安装成功"
 fi
else
 echo "docker 环境已经安装"  
fi
## 安装docker-compose
docker-compose version &> /dev/null
if [ $? -ne 0 ]; then
   sudo curl -L "https://get.daocloud.io/docker/compose/releases/download/1.25.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
   chmod +x /usr/local/bin/docker-compose
   docker-compose --version &> /dev/null
   if [ $? -ne 0 ]; then
     echo "docker-compose 安装失败"
     exit 1
   fi
   echo "docker-compose 安装成功"
else
 echo "docker-compose 已经安装"
fi

## 
## 下载镜像
FABRIC_TAG=1.4.4
CA_TAG=1.4.4
THIRDPARTY_TAG=0.4.18
for IMAGES in peer orderer ccenv  tools; do
    result=`docker images | grep hyperledger/fabric-$IMAGES | awk '{FS=" "} {print $1}'`
    if [ ! $result ]; then
        echo "==> FABRIC IMAGE: $IMAGES"
        echo
        docker pull hyperledger/fabric-$IMAGES:$FABRIC_TAG
        docker tag hyperledger/fabric-$IMAGES:$FABRIC_TAG hyperledger/fabric-$IMAGES
    fi   
done

resultca=`docker images | grep hyperledger/fabric-ca | awk '{FS=" "} {print $1}'`
if [ ! $resultca ]; then
    docker pull hyperledger/fabric-ca:$CA_TAG
    docker tag hyperledger/fabric-ca:$CA_TAG hyperledger/fabric-ca
fi
for IMAGES in couchdb kafka zookeeper; do

    resultca=`docker images | grep hyperledger/fabric-$IMAGES | awk '{FS=" "} {print $1}'`
    if [ ! $resultca ]; then
        echo "==> THIRDPARTY DOCKER IMAGE: $IMAGES"
        docker pull hyperledger/fabric-$IMAGES:$THIRDPARTY_TAG
        docker tag hyperledger/fabric-$IMAGES:$THIRDPARTY_TAG hyperledger/fabric-$IMAGES
    fi
done
echo "pull docker images success"

exit 0
