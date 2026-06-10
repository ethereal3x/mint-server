CREATE TABLE IF NOT EXISTS `tbl_base_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户唯一id',
  `username` varchar(64) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户登录名',
  `nickname` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '平台昵称',
  `avatar_url` varchar(512) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户头像地址',
  `realname` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '真实名字',
  `password` char(60) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码',
  `password2` varchar(60) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '二级密码',
  `reg_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '注册时间',
  `mobile` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `last_login_ip` varchar(45) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '最后登录的ip',
  `last_login_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '最后登录的时间戳',
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_user_id` (`user_id`),
  UNIQUE KEY `uni_username` (`username`),
  KEY `idx_mobile` (`mobile`),
  KEY `idx_nickname` (`nickname`),
  KEY `idx_realname` (`realname`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户表';

CREATE TABLE IF NOT EXISTS `tbl_third_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户唯一id',
  `channel_code` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '渠道标识',
  `union_id` varchar(256) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `open_id` varchar(256) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_union_id` (`user_id`,`union_id`,`channel_code`),
  KEY `idx_user_open_id` (`user_id`,`open_id`,`channel_code`),
  KEY `idx_union_id_channel` (`union_id`,`channel_code`,`user_id`),
  KEY `idx_open_id_channel` (`open_id`,`channel_code`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='第三方用户关联表';
