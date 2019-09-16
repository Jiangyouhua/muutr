# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 114.115.210.121 (MySQL 5.5.56-MariaDB)
# Database: jiang_ido
# Generation Time: 2018-05-08 11:37:51 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table jyh_ad
# ------------------------------------------------------------

DROP TABLE IF EXISTS `jyh_ad`;

CREATE TABLE `jyh_ad` (
  `aid` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `image` varchar(255) NOT NULL DEFAULT '' COMMENT '显示的图片',
  `link` varchar(255) NOT NULL DEFAULT '' COMMENT '打开的连接',
  `start` datetime DEFAULT NULL COMMENT '投放开始时间',
  `end` datetime DEFAULT NULL COMMENT '设放结束时间',
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '发布时间',
  `admin` varchar(255) NOT NULL DEFAULT '' COMMENT '发布的管理员',
  `status` tinyint(4) NOT NULL DEFAULT '1',
  PRIMARY KEY (`aid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table jyh_chat
# ------------------------------------------------------------

DROP TABLE IF EXISTS `jyh_chat`;

CREATE TABLE `jyh_chat` (
  `cid` varchar(255) NOT NULL,
  `sid` varchar(255) DEFAULT '' COMMENT '分类ID',
  `user` varchar(255) DEFAULT '' COMMENT '发布者',
  `date` timestamp(6) NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '时间',
  `content` varchar(2047) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '内容',
  `username` varchar(255) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '发布者别名',
  `at` varchar(2047) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '@什么人',
  `analysis` varchar(4) DEFAULT '0' COMMENT '内容的后缀',
  `status` tinyint(4) DEFAULT '1',
  PRIMARY KEY (`cid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table jyh_config
# ------------------------------------------------------------

DROP TABLE IF EXISTS `jyh_config`;

CREATE TABLE `jyh_config` (
  `uid` int(11) NOT NULL AUTO_INCREMENT,
  `server` varchar(255) DEFAULT '' COMMENT '主服务器',
  `back` varchar(255) DEFAULT '' COMMENT '后台服务器',
  `sync` varchar(255) DEFAULT '' COMMENT '同步服务器',
  `file` varchar(255) DEFAULT '' COMMENT '文件服务器',
  `upload` varchar(255) DEFAULT '' COMMENT '上传服务器',
  `socket` varchar(255) DEFAULT '' COMMENT 'WebSocket服务器',
  `html` varchar(255) DEFAULT '' COMMENT '静态服务器',
  `date` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table jyh_holiday
# ------------------------------------------------------------

DROP TABLE IF EXISTS `jyh_holiday`;

CREATE TABLE `jyh_holiday` (
  `uid` int(11) NOT NULL,
  `year` int(11) DEFAULT '0' COMMENT '年',
  `month` int(11) DEFAULT '0' COMMENT '月',
  `day` int(11) DEFAULT '0' COMMENT '日',
  `date` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `status` tinyint(4) DEFAULT '1' COMMENT '是否有节假日',
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table jyh_item
# ------------------------------------------------------------

DROP TABLE IF EXISTS `jyh_item`;

CREATE TABLE `jyh_item` (
  `iid` varchar(255) NOT NULL,
  `sid` varchar(255) DEFAULT '' COMMENT '分类ID',
  `explain` varchar(255) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '任务说明',
  `cycle` tinyint(4) DEFAULT '0' COMMENT '周期',
  `option` varchar(255) DEFAULT '' COMMENT '周期的选项',
  `start` date DEFAULT NULL COMMENT '开始日期',
  `complete` varchar(255) DEFAULT '' COMMENT '完成时间',
  `alert` varchar(255) DEFAULT '' COMMENT '提醒时间',
  `end` date DEFAULT NULL COMMENT '结束日期',
  `first` timestamp NULL DEFAULT NULL COMMENT '发布时间',
  `rank` TINYINT(4) NOT NULL  DEFAULT '0',
  `date` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `editor` int(11) NULL DEFAULT NULL COMMENT '最后的编辑者ID',
  `status` tinyint(4) DEFAULT '1',
  PRIMARY KEY (`iid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table jyh_join
# ------------------------------------------------------------

DROP TABLE IF EXISTS `jyh_join`;

CREATE TABLE `jyh_join` (
  `uid` varchar(255) NOT NULL COMMENT 'user,sid',
  `sid` varchar(255) DEFAULT '' COMMENT '分类ID',
  `name` varchar(255) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '自定义的分类名称',
  `user` varchar(255) DEFAULT '' COMMENT '用户',
  `sequence` int(11) DEFAULT '0' COMMENT '排序',
  `first` timestamp NULL DEFAULT NULL COMMENT '发布时间',
  `date` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `editor` int(11) NULL DEFAULT NULL COMMENT '最后的编辑者ID',
  `status` tinyint(4) DEFAULT '1',
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table jyh_logs
# ------------------------------------------------------------

DROP TABLE IF EXISTS `jyh_logs`;

CREATE TABLE `jyh_logs` (
  `uid` varchar(255) NOT NULL,
  `user` varchar(255) DEFAULT '' COMMENT '用户',
  `iid` varchar(255) DEFAULT '' COMMENT '任务ID',
  `explain` varchar(1024) DEFAULT '' COMMENT '生成的内容',
  `date` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '生成的时间',
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table jyh_mate
# ------------------------------------------------------------

DROP TABLE IF EXISTS `jyh_mate`;

CREATE TABLE `jyh_mate` (
  `uid` varchar(255) NOT NULL COMMENT 'sid,user',
  `sid` varchar(255) DEFAULT '' COMMENT '分类ID',
  `user` varchar(255) DEFAULT '' COMMENT '用户',
  `username` varchar(255) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '用户名称',
  `role` tinyint(4) DEFAULT '1' COMMENT '权限，成员为1, 创建者为2',
  `first` timestamp NULL DEFAULT NULL COMMENT '发布时间',
  `date` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `editor` int(11) NULL DEFAULT NULL COMMENT '最后的编辑者ID',
  `status` tinyint(4) DEFAULT '1',
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table jyh_money
# ------------------------------------------------------------

DROP TABLE IF EXISTS `jyh_money`;

CREATE TABLE `jyh_money` (
  `mid` varchar(255) NOT NULL DEFAULT '' COMMENT 'user+iid+ym',
  `user` varchar(255) DEFAULT '' COMMENT '用户',
  `year` int(11) DEFAULT '0' COMMENT '年',
  `month` int(11) DEFAULT '0' COMMENT '月',
  `sid` varchar(255) DEFAULT '' COMMENT '类ID',
  `iid` varchar(255) DEFAULT '' COMMENT '任务ID',
  `income` decimal(15,2) DEFAULT '0.00' COMMENT '收入',
  `expenses` decimal(15,2) DEFAULT '0.00' COMMENT '支出',
  `date` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '生成的时间',
  `status` tinyint(4) DEFAULT '1',
  PRIMARY KEY (`mid`),
  UNIQUE KEY `moneyIndex` (`user`,`year`,`month`,`iid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table jyh_node
# ------------------------------------------------------------

DROP TABLE IF EXISTS `jyh_node`;

CREATE TABLE `jyh_node` (
  `uid` varchar(255) NOT NULL COMMENT 'timline.uid,user',
  `tid` varchar(255) DEFAULT '' COMMENT '周期ID',
  `iid` varchar(255) DEFAULT '' COMMENT '任务ID',
  `ymd` int(11) DEFAULT '0' COMMENT '周期所有的日期',
  `user` varchar(255) DEFAULT '' COMMENT '用户',
  `username` varchar(255) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '用户名称',
  `first` timestamp NULL DEFAULT NULL COMMENT '发布的时间',
  `date` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改的时间',
  `editor` int(11) NULL DEFAULT NULL COMMENT '最后的编辑者ID',
  `status` tinyint(4) DEFAULT '0',
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table jyh_note
# ------------------------------------------------------------

DROP TABLE IF EXISTS `jyh_note`;

CREATE TABLE `jyh_note` (
  `nid` varchar(255) NOT NULL COMMENT 'user,timestamp',
  `iid` varchar(255) DEFAULT '' COMMENT '任务ID',
  `first` timestamp NULL DEFAULT NULL COMMENT '发布的时间',
  `date` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改的时间',
  `user` varchar(255) DEFAULT '' COMMENT '用户',
  `content` varchar(2045) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '内容',
  `username` varchar(255) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '用户名称',
  `money` decimal(15,2) DEFAULT '0.00' COMMENT '金额',
  `analysis` varchar(4) DEFAULT '0' COMMENT '是否需要解析',
  `editor` int(11) NULL DEFAULT NULL COMMENT '最后的编辑者ID',
  `status` tinyint(4) DEFAULT '1',
  PRIMARY KEY (`nid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table jyh_pusher
# ------------------------------------------------------------

DROP TABLE IF EXISTS `jyh_pusher`;

CREATE TABLE `jyh_pusher` (
  `pid` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) DEFAULT '' COMMENT '标题',
  `body` varchar(1024) DEFAULT '' COMMENT '内容',
  `image` varchar(255) DEFAULT '' COMMENT '图片',
  `category` varchar(255) DEFAULT '' COMMENT '分布的类型',
  `push` datetime DEFAULT NULL COMMENT '发布的时间',
  `admin` varchar(255) DEFAULT '' COMMENT '发布的管理员',
  `status` tinyint(4) DEFAULT '1',
  PRIMARY KEY (`pid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table jyh_record
# ------------------------------------------------------------

DROP TABLE IF EXISTS `jyh_record`;

CREATE TABLE `jyh_record` (
  `rid` varchar(255) NOT NULL DEFAULT '' COMMENT '聊天内容ID',
  `cid` varchar(255) NOT NULL DEFAULT '' COMMENT '聊天内容ID',
  `user` varchar(255) NOT NULL DEFAULT '' COMMENT '用户',
  `editor` int(11) NULL DEFAULT NULL COMMENT '最后的编辑者ID',
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '生成的时间',
  `status` tinyint(4) NOT NULL DEFAULT '1',
  KEY `recordIndex` (`user`,`cid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table jyh_sort
# ------------------------------------------------------------

DROP TABLE IF EXISTS `jyh_sort`;

CREATE TABLE `jyh_sort` (
  `sid` varchar(255) NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '名称',
  `explain` varchar(255) DEFAULT '' COMMENT '口号，说明',
  `share` int(11) DEFAULT '0' COMMENT '发享的次数',
  `first` timestamp NULL DEFAULT NULL COMMENT '发布的时间',
  `date` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改的时间',
  `editor` int(11) NULL DEFAULT NULL COMMENT '最后的编辑者ID',
  `status` tinyint(4) DEFAULT '1',
  PRIMARY KEY (`sid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table jyh_task
# ------------------------------------------------------------

DROP TABLE IF EXISTS `jyh_claim`;

CREATE TABLE `jyh_claim` (
  `tid` varchar(255) NOT NULL DEFAULT '' COMMENT '周期的ID',
  `uid` varchar(255) DEFAULT '',
  `sid` varchar(255) DEFAULT '' COMMENT '分类ID',
  `user` varchar(255) DEFAULT '' COMMENT '用户',
  `iid` varchar(255) DEFAULT '' COMMENT '任务ID',
  `alert` varchar(25) DEFAULT '' COMMENT '提醒的日期，未用',
  `date` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '生成的时间',
  `editor` int(11) NULL DEFAULT NULL COMMENT '最后的编辑者ID',
  `status` tinyint(4) DEFAULT '1',
  PRIMARY KEY (`tid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table jyh_user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `jyh_user`;

CREATE TABLE `jyh_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` varchar(255) NULL DEFAULT NULL,
  `account` varchar(255) DEFAULT '',
  `code` varchar(64) DEFAULT '',
  `username` varchar(255) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '用户名',
  `score` int(11) DEFAULT '0' COMMENT '完成任务周期数',
  `first` timestamp NULL DEFAULT NULL COMMENT '注册时间',
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后登录时间',
  `alert` tinyint(4) DEFAULT '1' COMMENT '提醒',
  `claim` tinyint(4) DEFAULT '0' COMMENT '显示领取任务',
  `cellular` tinyint(4) DEFAULT '1' COMMENT '使用蜂窝数据',
  `role` tinyint(4) DEFAULT '1' COMMENT '权限',
  `token` varchar(255) DEFAULT '' COMMENT '用户的苹果Tokan',
  `status` tinyint(4) DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `index` (`user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `jyh_user` WRITE;
/*!40000 ALTER TABLE `jyh_user` DISABLE KEYS */;

INSERT INTO `jyh_user` (`id`, `user`, `account`, `code`, `username`, `score`, `first`, `date`, `alert`, `claim`, `cellular`, `role`, `token`, `status`)
VALUES
	(1,NULL ,'admin@muutr.com','3ff3b966da4dad7ba1eca026d8167e28','管理员',0,'2017-01-01 00:00:00','2018-04-17 19:34:03',1,0,1,9,'',1),
	(2,NULL ,'18086306467','382a96aa2ef012b8c8036b9c7a07d2ae','芷兰',0,'2018-04-16 10:57:42','2018-05-05 11:53:44',1,1,1,1,'496b41f1cafe372313cdecdf52e200d9a3747726c073a049f2d09eb99a00a20d',1),
	(3,NULL ,'13797780504','18b9679efde0a1e990a0397793546b35','灰白',10,'2018-04-16 11:03:40','2018-05-08 19:11:38',1,0,0,1,'',1);

/*!40000 ALTER TABLE `jyh_user` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
