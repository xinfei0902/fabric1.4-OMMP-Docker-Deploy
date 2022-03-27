-- MySQL dump 10.13  Distrib 5.7.34, for Linux (x86_64)
--
-- Host: 127.0.0.1    Database: fabric
-- ------------------------------------------------------
-- Server version	5.7.34

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `ommp_chaincode`
--

DROP TABLE IF EXISTS `ommp_chaincode`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ommp_chaincode` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `task_id` varchar(128) NOT NULL COMMENT '每个任务唯一id',
  `cc_name` varchar(45) NOT NULL COMMENT '合约名称',
  `cc_version` varchar(45) NOT NULL COMMENT '合约版本',
  `cc_policy` varchar(512) DEFAULT NULL COMMENT '合约背书策略',
  `cc_org` varchar(45) DEFAULT NULL COMMENT '合约背书参入的组织',
  `channelname` varchar(45) DEFAULT NULL COMMENT '合约安装于那个通道\n',
  `createtime` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '创建时间',
  `detail` varchar(512) DEFAULT NULL COMMENT '合约描述 ',
  `is_install` int(11) NOT NULL COMMENT '合约是否安装  0--没有安装  1--安装',
  `status` int(11) NOT NULL COMMENT '合约状态   0---停用  1---启用  2---失败   3--正在删除   4---已删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ommp_chaincode`
--

LOCK TABLES `ommp_chaincode` WRITE;
/*!40000 ALTER TABLE `ommp_chaincode` DISABLE KEYS */;
/*!40000 ALTER TABLE `ommp_chaincode` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ommp_indent`
--

DROP TABLE IF EXISTS `ommp_indent`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ommp_indent` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `task_id` varchar(128) NOT NULL COMMENT '每个任务唯一标识',
  `source_id` varchar(128) DEFAULT NULL COMMENT '最开始一键部署任务id',
  `channelname` varchar(45) NOT NULL COMMENT '通道名称',
  `create_cife_type` varchar(45) DEFAULT NULL COMMENT ' 创建证书类型',
  `consensus` varchar(45) DEFAULT NULL COMMENT '共识机制 solo,etcdraft等',
  `org_name` varchar(45) NOT NULL COMMENT '组织名称',
  `org_domain` varchar(45) NOT NULL COMMENT '组织域名',
  `peer_ip` varchar(45) NOT NULL COMMENT '节点IP',
  `peer_name` varchar(45) NOT NULL COMMENT '节点名称',
  `peer_user` varchar(45) NOT NULL COMMENT '节点身份',
  `peer_port` int(11) NOT NULL COMMENT '节点端口',
  `peer_id` int(11) NOT NULL COMMENT '节点顺序 是第几个节点为 上报这点状态指定的id',
  `couchdb_port` int(11) NOT NULL COMMENT '节点数据库端口',
  `cc_port` int(11) NOT NULL COMMENT '节点合约端口',
  `peer_domain` varchar(45) NOT NULL COMMENT '节点域名',
  `ca_name` varchar(45) NOT NULL COMMENT 'CA名称',
  `ca_ip` varchar(45) NOT NULL COMMENT 'CAIP',
  `ca_port` int(11) NOT NULL COMMENT 'CA端口',
  `cli_name` varchar(45) NOT NULL COMMENT 'cli容器名称',
  `version` varchar(45) NOT NULL COMMENT ' fabric 版本',
  `nick_name` varchar(45) NOT NULL COMMENT '节点昵称',
  `accesskey` varchar(45) NOT NULL COMMENT '数据上链标识来源',
  `peer_status` int(11) NOT NULL DEFAULT '0' COMMENT '节点创建 失败 状态    0---正在创建  1---再用  2---失败   3--正在删除   4---已删除',
  `channel_status` int(11) NOT NULL DEFAULT '0' COMMENT '通道创建 失败状态   0---正在创建  1---再用  2---失败  3--正在删除  4---已删除',
  `org_status` int(11) NOT NULL DEFAULT '0' COMMENT '组织创建 失败状态   0---正在创建  1---再用  2---失败  3--正在删除   4---已删除',
  `peer_run_status` int(11) NOT NULL DEFAULT '0' COMMENT '节点启用 停用状态   0---停用  1---启用  ',
  PRIMARY KEY (`id`),
  UNIQUE KEY `create_cife_type_UNIQUE` (`create_cife_type`)
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ommp_indent`
--

LOCK TABLES `ommp_indent` WRITE;
/*!40000 ALTER TABLE `ommp_indent` DISABLE KEYS */;
/*!40000 ALTER TABLE `ommp_indent` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ommp_orderer`
--

DROP TABLE IF EXISTS `ommp_orderer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ommp_orderer` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `orderer_name` varchar(45) NOT NULL COMMENT 'orderer名称',
  `orderer_domain` varchar(45) NOT NULL COMMENT '域名  组合： 名称+域名后缀',
  `orderer_ip` varchar(45) NOT NULL COMMENT ' ip地址',
  `orderer_port` int(11) NOT NULL COMMENT ' 端口',
  `channelname` varchar(45) NOT NULL COMMENT '通道名称',
  `orderer_orgdomain` varchar(45) NOT NULL COMMENT '域名后缀',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ommp_orderer`
--

LOCK TABLES `ommp_orderer` WRITE;
/*!40000 ALTER TABLE `ommp_orderer` DISABLE KEYS */;
/*!40000 ALTER TABLE `ommp_orderer` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ommp_server`
--

DROP TABLE IF EXISTS `ommp_server`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ommp_server` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增排序',
  `task_id` varchar(128) NOT NULL COMMENT '每一个任务唯一id',
  `source_id` varchar(128) DEFAULT NULL COMMENT ' 原任务id',
  `server_name` varchar(45) NOT NULL COMMENT '主机服务名称',
  `server_des` varchar(128) DEFAULT NULL COMMENT '主机服务描述',
  `server_extip` varchar(45) NOT NULL COMMENT '主机服务外网ip',
  `server_intip` varchar(45) NOT NULL COMMENT '主机服务内网ip',
  `server_user` varchar(45) NOT NULL COMMENT '主机服务远程ssh连接用户名',
  `server_password` varchar(45) NOT NULL COMMENT '主机服务远程ssh连接密码',
  `server_num` varchar(10) NOT NULL COMMENT ' 主机服务上使用节点个数',
  `server_status` varchar(10) NOT NULL COMMENT '主机服务状态  0---正在创建  1---再用  2---失败  3--正在删除  4---已删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ommp_server`
--

LOCK TABLES `ommp_server` WRITE;
/*!40000 ALTER TABLE `ommp_server` DISABLE KEYS */;
/*!40000 ALTER TABLE `ommp_server` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-07-11 11:24:40
