-- MySQL dump 10.13  Distrib 5.7.19, for Linux (x86_64)
--
-- Host: localhost    Database: rbac
-- ------------------------------------------------------
-- Server version	5.7.19-0ubuntu0.16.04.1

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
-- Table structure for table `depart`
--

DROP TABLE IF EXISTS `depart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `depart` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '0',
  `pid` int(11) DEFAULT '0',
  `code` varchar(50) NOT NULL DEFAULT '0',
  `level` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `pids` varchar(100) DEFAULT NULL,
  `gmt_create` datetime DEFAULT NULL,
  `gmt_modified` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8 COMMENT='分组';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `depart`
--

LOCK TABLES `depart` WRITE;
/*!40000 ALTER TABLE `depart` DISABLE KEYS */;
INSERT INTO `depart` VALUES (1,'root',0,'0',0,NULL,NULL,NULL),(2,'china',1,'001',1,'1',NULL,NULL),(5,'shanghai',2,'002',2,'1,2',NULL,NULL),(11,'shandong',2,'003',2,'1,2',NULL,NULL),(12,'zibo',11,'001',3,'1,2,11',NULL,NULL),(13,'jinan',11,'002',3,'1,2,11',NULL,NULL);
/*!40000 ALTER TABLE `depart` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `depart_res`
--

DROP TABLE IF EXISTS `depart_res`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `depart_res` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `depart_id` int(11) DEFAULT NULL,
  `res_id` int(11) NOT NULL,
  `type` int(11) DEFAULT '0' COMMENT '0allow,1unallow',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `depart_res`
--

LOCK TABLES `depart_res` WRITE;
/*!40000 ALTER TABLE `depart_res` DISABLE KEYS */;
INSERT INTO `depart_res` VALUES (2,12,29,0);
/*!40000 ALTER TABLE `depart_res` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `depart_role`
--

DROP TABLE IF EXISTS `depart_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `depart_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `depart_id` int(11) DEFAULT NULL,
  `role_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `depart_role`
--

LOCK TABLES `depart_role` WRITE;
/*!40000 ALTER TABLE `depart_role` DISABLE KEYS */;
INSERT INTO `depart_role` VALUES (1,12,9),(2,2,7),(3,5,8),(4,13,9);
/*!40000 ALTER TABLE `depart_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `resource`
--

DROP TABLE IF EXISTS `resource`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `resource` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `res_name` char(50) NOT NULL,
  `res_code` char(50) NOT NULL,
  `pid` int(11) NOT NULL,
  `gmt_create` datetime NOT NULL,
  `gmt_modified` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8 COMMENT='资源';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `resource`
--

LOCK TABLES `resource` WRITE;
/*!40000 ALTER TABLE `resource` DISABLE KEYS */;
INSERT INTO `resource` VALUES (1,'root','root',0,'2017-10-10 00:00:00','2017-10-10 00:00:00'),(15,'2','2',14,'2017-10-12 08:48:41','2017-10-12 08:48:41'),(20,'1','1',1,'2017-10-12 00:54:44','2017-10-12 09:18:32'),(24,'2','2',20,'2017-10-12 09:19:09','2017-10-12 09:19:09'),(25,'3','3',20,'2017-10-19 06:02:21','2017-10-19 06:02:21'),(26,'user','user',1,'2017-10-18 23:16:07','2017-10-19 07:16:27'),(27,'user:add','user:add',26,'2017-10-18 23:16:40','2017-10-25 05:59:16'),(28,'user:edit','user:edit',26,'2017-10-19 07:16:54','2017-10-19 07:16:54'),(29,'user:del','user:del',26,'2017-10-19 07:17:03','2017-10-19 07:17:03'),(30,'user:detail','user:detail',26,'2017-10-25 05:59:28','2017-10-25 05:59:28'),(31,'depart','depart',1,'2017-10-25 06:00:11','2017-10-25 06:00:11'),(32,'depart:add','depart:add',31,'2017-10-25 06:00:26','2017-10-25 06:00:26'),(33,'depart:edit','depart:edit',31,'2017-10-25 06:00:36','2017-10-25 06:00:36'),(34,'depart:del','depart:del',31,'2017-10-25 06:00:48','2017-10-25 06:00:48'),(35,'depart:detail','depart:detail',31,'2017-10-25 06:01:00','2017-10-25 06:01:00'),(36,'depart:allotrole','depart:allotrole',31,'2017-10-25 06:30:55','2017-10-25 06:30:55'),(37,'depart:allotres','depart:allotres',31,'2017-10-25 06:31:10','2017-10-25 06:31:10'),(38,'user:allotdepart','user:allotdepart',26,'2017-10-25 06:34:11','2017-10-25 06:34:11'),(39,'user:allotrole','user:allotrole',26,'2017-10-25 06:34:22','2017-10-25 06:34:22'),(40,'user:allotres','user:allotres',26,'2017-10-25 06:34:34','2017-10-25 06:34:34');
/*!40000 ALTER TABLE `resource` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role`
--

DROP TABLE IF EXISTS `role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_name` char(30) NOT NULL DEFAULT '0',
  `description` varchar(100) DEFAULT '0',
  `pid` int(11) DEFAULT '0',
  `gmt_create` datetime DEFAULT NULL,
  `gmt_modified` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8 COMMENT='角色';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role`
--

LOCK TABLES `role` WRITE;
/*!40000 ALTER TABLE `role` DISABLE KEYS */;
INSERT INTO `role` VALUES (1,'root','root',0,NULL,NULL),(7,'country','undefined',1,'2017-10-18 00:48:41','2017-10-19 02:02:35'),(8,'province','province',7,'2017-10-19 02:48:35','2017-10-19 02:48:35'),(9,'city','city',8,'2017-10-20 17:46:55','2017-10-21 01:55:48');
/*!40000 ALTER TABLE `role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role_res`
--

DROP TABLE IF EXISTS `role_res`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `role_res` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL DEFAULT '0',
  `res_id` int(11) NOT NULL DEFAULT '0',
  `gmt_create` datetime NOT NULL,
  `gmt_modified` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=68 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role_res`
--

LOCK TABLES `role_res` WRITE;
/*!40000 ALTER TABLE `role_res` DISABLE KEYS */;
INSERT INTO `role_res` VALUES (56,7,24,'2017-10-19 10:02:35','2017-10-19 10:02:35'),(57,7,25,'2017-10-19 10:02:35','2017-10-19 10:02:35'),(58,7,27,'2017-10-19 10:02:35','2017-10-19 10:02:35'),(59,7,28,'2017-10-19 10:02:35','2017-10-19 10:02:35'),(60,7,29,'2017-10-19 10:02:35','2017-10-19 10:02:35'),(61,8,24,'2017-10-19 10:48:35','2017-10-19 10:48:35'),(62,8,25,'2017-10-19 10:48:35','2017-10-19 10:48:35'),(66,9,27,'2017-10-21 01:55:48','2017-10-21 01:55:48'),(67,9,28,'2017-10-21 01:55:48','2017-10-21 01:55:48');
/*!40000 ALTER TABLE `role_res` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` char(20) DEFAULT NULL,
  `birth` datetime DEFAULT NULL,
  `gender` tinyint(2) DEFAULT NULL,
  `addr` varchar(100) DEFAULT NULL,
  `pwd` char(32) DEFAULT NULL,
  `mobile` char(30) DEFAULT NULL,
  `email` char(30) DEFAULT NULL,
  `real_name` char(10) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  `token` char(32) DEFAULT NULL,
  `client_id` char(32) DEFAULT NULL,
  `last_login_time` datetime DEFAULT NULL,
  `gmt_create` datetime DEFAULT NULL,
  `gmt_modified` datetime DEFAULT NULL,
  `identity_card` char(18) DEFAULT NULL,
  `avatar` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,'gsy',NULL,0,'','e10adc3949ba59abbe56e057f20f883e','110','','gsy1',0,'','',NULL,NULL,NULL,'','static/avatar/1.png'),(2,'gsy2',NULL,0,'','e10adc3949ba59abbe56e057f20f883e','220','','gsy0',0,'','',NULL,NULL,NULL,'',NULL),(3,'gsy3',NULL,0,'','e10adc3949ba59abbe56e057f20f883e','330','','gsy3',0,'','',NULL,NULL,NULL,'330',NULL),(4,'gsy4',NULL,0,'','e10adc3949ba59abbe56e057f20f883e','440','','gsy4',0,'','',NULL,NULL,NULL,'440',NULL),(5,'gsy5',NULL,0,'','e10adc3949ba59abbe56e057f20f883e','550','','gsy5',0,'','',NULL,NULL,NULL,'550',NULL),(7,'gsy7',NULL,0,'','e10adc3949ba59abbe56e057f20f883e','770','','gsy7',0,'','',NULL,NULL,NULL,'770',NULL),(8,'gsy8',NULL,0,'','e10adc3949ba59abbe56e057f20f883e','110','','gsy',0,'','',NULL,NULL,NULL,'',NULL),(9,'root',NULL,0,'','e10adc3949ba59abbe56e057f20f883e','000','','root',0,'','',NULL,NULL,NULL,'000',NULL);
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_depart`
--

DROP TABLE IF EXISTS `user_depart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_depart` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `depart_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_depart`
--

LOCK TABLES `user_depart` WRITE;
/*!40000 ALTER TABLE `user_depart` DISABLE KEYS */;
INSERT INTO `user_depart` VALUES (15,1,12),(16,9,1);
/*!40000 ALTER TABLE `user_depart` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_res`
--

DROP TABLE IF EXISTS `user_res`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_res` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0',
  `res_id` int(11) NOT NULL DEFAULT '0',
  `type` int(11) DEFAULT '0' COMMENT '0allow,1unallow',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=93 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_res`
--

LOCK TABLES `user_res` WRITE;
/*!40000 ALTER TABLE `user_res` DISABLE KEYS */;
INSERT INTO `user_res` VALUES (92,1,35,0);
/*!40000 ALTER TABLE `user_res` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_role`
--

DROP TABLE IF EXISTS `user_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `role_id` int(11) NOT NULL,
  `type` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_role`
--

LOCK TABLES `user_role` WRITE;
/*!40000 ALTER TABLE `user_role` DISABLE KEYS */;
INSERT INTO `user_role` VALUES (1,9,1,0),(2,9,0,1);
/*!40000 ALTER TABLE `user_role` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-11-02 10:44:19
