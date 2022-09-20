CREATE TABLE `t_wechat_conf`
(
    `wechat_id`    mediumint(9) NOT NULL AUTO_INCREMENT COMMENT '公众号/小程序配置id',
    `access_token` varchar(512) NOT NULL DEFAULT '' COMMENT '公众号接口调用token',
    PRIMARY KEY (`wechat_id`) USING BTREE
) ENGINE = InnoDB
  AUTO_INCREMENT = 16
  DEFAULT CHARSET = utf8
  ROW_FORMAT = DYNAMIC COMMENT ='微信公众号配置表';

