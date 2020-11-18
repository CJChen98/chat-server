SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for messages
-- ----------------------------
DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages`
(
    `id`              int(11) UNSIGNED                                              NOT NULL AUTO_INCREMENT,
    `user_id`         varchar(100)                                                  NOT NULL COMMENT '用户ID',
    `username`        varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '' COMMENT '发送者昵称',
    `conversation_id` int(11)                                                       NOT NULL COMMENT '会话ID',
    `receiver_id`     varchar(100)                                                  NOT NULL COMMENT '接收者ID',
    `content`         longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci     NULL COMMENT '聊天内容',
    `image_url`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL     DEFAULT '' COMMENT '图片URL',
    `created_at`      datetime(0)                                                   NULL     DEFAULT NULL,
    `updated_at`      datetime(0)                                                   NULL     DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(0),
    `deleted_at`      datetime(0)                                                   NULL     DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    INDEX `idx_id` (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci
  ROW_FORMAT = DYNAMIC;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `id`          int(11) UNSIGNED                                              NOT NULL AUTO_INCREMENT,
    `username`    varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '' COMMENT '昵称',
    `password`    varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL     DEFAULT '' COMMENT '密码',
    `avatar_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '/img?type=user&id=default' COMMENT '头像URL',
    `created_at`  datetime(0)                                                   NULL     DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(0),
    `updated_at`  datetime(0)                                                   NULL     DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(0),
    `deleted_at`  datetime(0)                                                   NULL     DEFAULT NULL,
    `snow_id`     varchar(100)                                                  not null comment '用户id',

    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `username` (`username`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci
  ROW_FORMAT = DYNAMIC;

SET FOREIGN_KEY_CHECKS = 1;

DROP TABLE if exists `rooms`;
create table `rooms`
(
    `id`           int(11) unsigned                                              not null auto_increment,
    `creator_id`   varchar(100)                                                  not null,
    `room_name`    varchar(50) char set utf8mb4 COLLATE utf8mb4_general_ci       not null default '' comment '群名',
    `member_size`  int(11)                                                       not null default 1 comment '群人数',
    `introduction` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci     NULL COMMENT '聊天介绍',
    `snow_id`      varchar(100)                                                  not null comment '群id',
    `created_at`   long                                                          null comment '群创建时间',
    `avatar_path`  varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '/img?type=room&id=default' COMMENT '头像URL',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `inx_id` (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci
  ROW_FORMAT = DYNAMIC;

DROP TABLE if exists `conversations`;
create table `conversations`
(
    `id`          int(11) unsigned not null auto_increment,
    `receiver_id` varchar(100)     not null comment '会话接收者',
    `private`     bool             not null default false comment '是否为私聊',
    `user_id`     varchar(100)     not null comment '会话发起者',
    `created_at`  datetime(0)      NULL     DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(0),
    `updated_at`  datetime(0)      NULL     DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(0),
    `deleted_at`  datetime(0)      NULL     DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `inx_id` (`id`) USING BTREE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci
  ROW_FORMAT = DYNAMIC;
