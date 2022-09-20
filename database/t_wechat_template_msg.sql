CREATE TABLE `t_wechat_template_msg`
(
    `bind_id`     int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `wechat_id`   int(11) NOT NULL DEFAULT '0' COMMENT '微信id',
    `wechat_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '微信类型 1.公众号 2.小程序',
    `status`      tinyint(1) NOT NULL DEFAULT '0' COMMENT '推送状态 0.待推送 1.正在推送 2.已推送 3.推送失败',
    `push_data`   json                                    NOT NULL COMMENT '推送数据，json字符串',
    `errmsg`      varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '微信接口errmsg',
    `add_time`    datetime                                NOT NULL DEFAULT '1000-01-01 00:00:00' COMMENT '新增时间',
    `update_time` timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`bind_id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 79
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='微信模板消息表';