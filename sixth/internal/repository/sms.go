package repository

import (
	"GeekTask/sixth/internal/domain"
	"context"
)

/*
懒得的实现了，这里用接口代替一下;表结构如下：
CREATE TABLE `demo` (
	  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
	  `status` tinyint(10) unsigned NOT NULL COMMENT '状态 1 未重试 2 重试成功 3 重试失败（人工介入）',
	  `create_time` int(10) unsigned NOT NULL COMMENT '创建时间',
      `update_time` int(10) unsigned NOT NULL COMMENT '修改时间',
	  `extra` json NOT NULL COMMENT '额外的参数',
	  PRIMARY KEY (`id`),
	  KEY `create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
*/

type SMSRepository interface {
	Create(ctx context.Context, extra domain.SMSCallBackArgs) error // create_time建立时，不会立马重试，延迟10s
	GetRetryList(ctx context.Context) []domain.SMSCallBackArgs      // 获取小于等于当前时间的状态为1的列表，默认取10
	UpdateStatus(ctx context.Context, status int) error
}
