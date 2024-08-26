CREATE TABLE `request_history` (
                                   `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
                                   `link_id` varchar(16) COLLATE utf8mb4_general_ci NOT NULL COMMENT '短链',
                                   `bid` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
                                   `url` varchar(512) COLLATE utf8mb4_general_ci NOT NULL COMMENT 'url',
                                   `ip` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'IP',
                                   `ua` varchar(265) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '访问者',
                                   `time` datetime NOT NULL COMMENT '创建时间',
                                   PRIMARY KEY (`id`),
                                   KEY `idx_link_id_bid` (`link_id`,`bid`) USING BTREE,
                                   KEY `idx_link_id` (`link_id`,`time` DESC) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- test.short_link definition

CREATE TABLE `short_link` (
                              `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
                              `link_id` varchar(16) NOT NULL COMMENT '短链',
                              `url` varchar(512) NOT NULL COMMENT '长连接',
                              `description` varchar(128) DEFAULT NULL COMMENT '描述',
                              `is_enable` tinyint NOT NULL DEFAULT '1' COMMENT '是否启用',
                              `created_by` varchar(32) DEFAULT NULL COMMENT '创建人',
                              `create_time` datetime NOT NULL COMMENT '创建时间',
                              `update_time` datetime NOT NULL COMMENT '更新时间',
                              PRIMARY KEY (`id`),
                              UNIQUE KEY `bid_UNIQUE` (`link_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- test.`user` definition

CREATE TABLE `user` (
                        `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                        `username` varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名',
                        `role` int NOT NULL COMMENT '角色',
                        `password` varchar(128) COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码',
                        `salt` varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '盐',
                        `create_time` datetime NOT NULL COMMENT '创建时间',
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `username_UNIQUE` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;