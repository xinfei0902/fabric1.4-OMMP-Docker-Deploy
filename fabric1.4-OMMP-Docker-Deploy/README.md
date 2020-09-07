# fabric1.4-OMMP-Deploy

## 项目功能

```sh
运维管理平台服务 - 新增通道，组织，节点，合约 
代码基本都已经汉语注释
```



## 项目配置 

### 1. config.json

```sh
{
    "address": ":7000",    ##服务端端口
    "cert": null,
    "config": "D:/GoPath/src/deploy-server/local.json",  ##配置文件路径 也可以是linux路径 例如： /root/deploy-server/local.json
    "local": "D:/GoPath/src/deploy-server/local.json",
    "log": "console",
    "loglevel": "debug",
    "mysqlUser":"root",   ## 下面如果使用数据库配置
    "mysqlPass":"123456",
    "mysqlIP":"localhost",
    "mysqlPort":"3306",
    "mysqlDBName":"fabric",
    "toolsPort":"8080",   ## 请求依赖工具端口 依赖工具辅助执行脚本主要针对ssh登录权限(比如银行项目)
    "debug": null,
    "help": null,
    "key": null
}

```

### 2. local.json

```sh
{
    "consensus": {  ## 采用那种共识方式
        "consensus": [
            "kafka",
            "solo"
        ]
    },
    "version": [
        {
            "version": "fabric1.4.4",  ##指定版本
            "versionroot": "D:/GoPath/src/deploy-server/version/fabric1.4.4", ## 指定版本目录下的模板脚本
            "binroot": "D:/GoPath/src/deploy-server/version/fabric1.4.4/bin", ##指定证书等工具路径
            "chaincode": {
                "mycc": {
                    "version": "1.0",
                    "desc": "chaincode"
                }
            },
            "build": {
                "mycc": {  ## build目录下有mycc合约 没有用到 
                    "type": "chaincode",
                    "version": "1.0",
                    "bin": "mycc"
                }
            },
            "sub": { 
                "kafka": "2.12",
                "zookeeper": "3.4.10"
            }
        }
    ],
    "bin": {
        "linux": {
            "crypto": "cryptogen",
            "configtx": "configtxgen"
        },
        "windows": {
            "crypto": "cryptogen.exe",
            "configtx": "configtxgen.exe"
        }
    },
    "storepath": "D:/GoPath/src/deploy-server/output" ## 输出目录
}
```

## 项目编译

### 1. 命令

```sh
./build ./app
```


## 项目启动
```sh
./app.bin service
```
