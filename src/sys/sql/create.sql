CREATE DATABASE `fcds`;

#用户表
CREATE TABLE `f_users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(100) NOT NULL,
  `password` varchar(100) NOT NULL DEFAULT '',
  `type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '0：普通用户，1：系统管理员',
  `totp` varchar(200) NOT NULL DEFAULT '' COMMENT 'totp短地址',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `last_login_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `token` varchar(500) DEFAULT '',
  `totp_secret` varchar(100) DEFAULT NULL COMMENT 'totp秘钥',
  `totp_url` varchar(300) DEFAULT NULL COMMENT 'totp完整url',
  `mail` varchar(100) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  KEY `(username,password)` (`password`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


#机群表
CREATE TABLE `f_sets_servers` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `ip` varchar(100) DEFAULT '' COMMENT 'outter ip',
  `port` int(11) DEFAULT NULL COMMENT 'ssh port',
  `type` tinyint(2) DEFAULT '0' COMMENT '类型，0：无效，1：生产，2：任务',
  `images` varchar(500) DEFAULT NULL COMMENT '基础镜像部署',
  `micro_service_num` int(11) DEFAULT NULL COMMENT '微服务数量',
  `username` varchar(100) DEFAULT NULL COMMENT 'ssh username',
  `password` varchar(100) DEFAULT NULL COMMENT 'ssh password',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip` (`ip`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

#二级代理表
CREATE TABLE `f_sets_proxys` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '0:secNginx，1：springCloud，2：go-kit',
  `ip` varchar(100) NOT NULL DEFAULT '',
  `port` varchar(100) NOT NULL DEFAULT '',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '0：断开，1：在线',
  `status_msg` varchar(100) NOT NULL DEFAULT '断开',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

#微服务名称表
CREATE TABLE `f_sets_ms` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

#杂项配置表
CREATE TABLE `f_sets_configs` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(100) NOT NULL DEFAULT '',
  `value` varchar(500) DEFAULT NULL,
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `comment` varchar(100) NULL DEFAULT '' COMMENT '说明',
  PRIMARY KEY (`id`),
  UNIQUE KEY `key` (`key`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

#微服务部署表
CREATE TABLE `f_ms` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `service_id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL DEFAULT '',
  `version` varchar(100) NOT NULL DEFAULT '',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '0:失效，1：部署中，2：在线',
  `statusMsg` varchar(100) DEFAULT NULL,
  `use_time` int(11) DEFAULT '0',
  `deploy_type` tinyint(2) NOT NULL DEFAULT '1' COMMENT '1：正式版本，2：灰度版本',
  `finish_time` datetime NOT NULL COMMENT '结束时间',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `(name,version)` (`service_id`),
  KEY `service_id` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

#用户操作日志
CREATE TABLE `f_logs` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `username` varchar(100) NOT NULL DEFAULT '',
  `route` varchar(100) DEFAULT NULL COMMENT '请求的路由',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `content` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

#任务机任务表
CREATE TABLE `f_job_tasks` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `tid` int(11) NOT NULL DEFAULT '0' COMMENT '任务机id',
  `time` varchar(200) NOT NULL DEFAULT '' COMMENT 'crontab时间格式，精确到秒',
  `value` varchar(500) NOT NULL DEFAULT '' COMMENT 'shell或url',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

#任务机代码仓库
CREATE TABLE `f_job_srcs` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `tid` int(11) DEFAULT NULL,
  `src` varchar(500) DEFAULT NULL COMMENT '仓库地址',
  `name` varchar(100) DEFAULT NULL COMMENT '存放名称',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `status` tinyint(2) DEFAULT NULL COMMENT '0:未部署或失败，1：部署成功',
  PRIMARY KEY (`id`),
  UNIQUE KEY `(tid,name)` (`tid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

#任务机命令管理
CREATE TABLE `f_job_cmds` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `tid` int(11) NOT NULL COMMENT '任务机id',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '命令名称',
  `value` varchar(500) NOT NULL DEFAULT '' COMMENT '命令实现',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `(tid,name)` (`tid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

#微服务镜像构建表
CREATE TABLE `f_builds` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `unique_id` varchar(200) NOT NULL DEFAULT '' COMMENT '唯一随机id',
  `service_id` int(11) DEFAULT NULL COMMENT '微服务id',
  `name` varchar(100) DEFAULT NULL COMMENT '微服务名称',
  `version` varchar(100) DEFAULT NULL COMMENT '版本号，如：v1.1.0',
  `src` varchar(500) DEFAULT NULL COMMENT 'svn地址',
  `path` varchar(500) DEFAULT NULL COMMENT '容器中项目代码目录',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '0:构建中，1：构建成功，2：构建失败',
  `status_msg` varchar(100) DEFAULT '' COMMENT '构建状态文本',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `log` text COMMENT '构建日志',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_id` (`unique_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

#核心配置
INSERT INTO `f_sets_configs` (`id`, `key`, `value`, `create_time`, `update_time`, `comment`)
VALUES
	(1, 'SshUsername', 'test', '2017-12-29 15:28:04', '2018-01-06 09:59:06', '核心依赖机ssh用户名称'),
	(2, 'SshPassword', 'test', '2017-12-29 15:28:13', '2018-01-06 09:59:01', '核心依赖机ssh用户密码'),
	(3, 'FacadeHost', '127.0.0.1 index.fuyoukache.com', '2017-12-29 15:27:44', '2018-01-06 10:08:53', '入口机host'),
	(4, 'FacadeAddress', '127.0.0.1:22', '2017-12-29 15:27:54', '2018-01-06 10:07:05', '入口机ssh地址'),
	(5, 'RegistryUsername', 'test', '2018-01-06 17:07:30', '2018-01-06 17:07:46', '镜像仓库统一账号'),
    (6, 'RegistryPassword', 'test', '2018-01-06 17:07:42', '2018-01-06 17:07:46', '镜像仓库统一密码'),
	(7, 'RegistryBaseAddress', '127.0.0.1:22', '2018-01-06 09:59:28', '2018-01-06 09:59:56', '基础镜像仓库机ssh地址'),
	(8, 'RegistryBaseDomain', 'registry-base.fuyoukache.com', '2017-12-29 15:28:30', '2018-01-03 14:32:10', '基础镜像仓库域'),
	(9, 'RegistryBaseHost', '127.0.0.1 registry-base.fuyoukache.com', '2017-12-29 15:28:37', '2018-01-03 14:32:10', '基础镜像仓库host'),
	(10, 'RegistryBaseRunName', 'registry-base', '2017-12-29 15:28:45', '2018-01-06 10:07:33', '基础镜像仓库docker名称'),
	(11, 'RegistryServiceAddress', '127.0.0.1:22', '2018-01-06 10:00:19', '2018-01-03 14:32:10', '微服务镜像仓库机ssh地址'),
	(12, 'RegistryServiceDomain', 'registry-service.fuyoukache.com', '2017-12-29 15:29:28', '2018-01-06 10:07:41', '微服务镜像仓库域'),
	(13, 'RegistryServiceHost', '127.0.0.1 registry-service.fuyoukache.com', '2017-12-29 15:29:35', '2018-01-06 10:07:51', '微服务镜像仓库host'),
	(14, 'RegistryServiceRunName', 'registry-service', '2017-12-29 15:29:44', '2018-01-06 10:07:57', '微服务镜像仓库docker名称'),
	(15, 'SvnUsername', 'test', '2017-12-29 15:41:49', '2018-01-06 10:08:00', 'svn用户名'),
	(16, 'SvnPassword', 'test', '2017-12-29 15:41:54', '2018-01-06 10:08:01', 'svn密码'),
	(17, 'SvnHost', '127.0.0.1 svn.fuyoukache.com', '2017-12-29 15:41:59', '2018-01-06 10:08:10', 'svn host'),
	(18, 'GitUsername', 'test', '2017-12-29 15:41:59', '2018-01-06 10:08:19', 'git用户名'),
	(19, 'GitPassword', 'test', '2017-12-29 15:42:14', '2018-01-06 10:08:22', 'git密码'),
	(20, 'GitHost', '127.0.0.1 git.fuyoukache.com', '2017-12-29 15:42:18', '2018-01-06 10:08:29', 'git host'),
	(21, 'DingDing', 'http://dingding.com/post/webhook', '2017-12-29 15:42:22', '2018-01-06 10:08:39', '钉钉报警webhook机器人地址'),
	(22, 'NodeTaskDefaultImages', 'java,python,php', '2017-12-29 16:01:17', '2018-01-03 14:32:11', '机群任务机默认镜像'),
	(23, 'AppMode', 'production', '2017-12-29 15:41:19', '2018-01-06 12:00:26', 'fcds运行环境production或develop');

