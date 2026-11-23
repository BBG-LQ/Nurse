CREATE DATABASE IF NOT EXISTS care_order DEFAULT CHARSET utf8mb4;
USE care_order;

CREATE TABLE `service_order` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `order_no` VARCHAR(32) NOT NULL,
  `elderly_id` BIGINT UNSIGNED NOT NULL,
  `nurse_id` BIGINT UNSIGNED DEFAULT NULL,
  `service_type` VARCHAR(64) NOT NULL,
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0待派 1已接单 2服务中 3已完成 4已取消',
  `urgency` TINYINT NOT NULL DEFAULT 0 COMMENT '1紧急',
  `scheduled_at` DATETIME DEFAULT NULL,
  `started_at` DATETIME DEFAULT NULL,
  `finished_at` DATETIME DEFAULT NULL,
  `amount` DECIMAL(10,2) DEFAULT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_elderly` (`elderly_id`),
  KEY `idx_nurse_status` (`nurse_id`, `status`),
  KEY `idx_status_urgency` (`status`, `urgency`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='护理服务订单';

CREATE TABLE `satisfaction` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `order_id` BIGINT UNSIGNED NOT NULL,
  `nurse_id` BIGINT UNSIGNED NOT NULL,
  `score` TINYINT NOT NULL COMMENT '1-5',
  `comment` VARCHAR(512) DEFAULT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order` (`order_id`),
  KEY `idx_nurse` (`nurse_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='服务满意度';