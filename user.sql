CREATE TABLE `users` (
  `id` int(64) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL COMMENT '用户名, 用户名不允许重复的',
  `password` varchar(255) NOT NULL COMMENT '不能保持用户的明文密码',
  `label` varchar(255) NOT NULL COMMENT '用户标签',
  `role` tinyint(4) NOT NULL COMMENT '用户角色,0表示普通用户,1表示审核用户,2表示管理员',
  `create_time` int(64) NOT NULL COMMENT '创建时间',
  `update_time` int(64) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_user` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4;