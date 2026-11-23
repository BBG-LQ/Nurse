CREATE DATABASE IF NOT EXISTS care_alert DEFAULT CHARSET utf8mb4;
USE care_alert;

CREATE TABLE `alert_event` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `elderly_id` BIGINT UNSIGNED DEFAULT NULL,
  `nurse_id` BIGINT UNSIGNED DEFAULT NULL,
  `order_id` BIGINT UNSIGNED DEFAULT NULL,
  `alert_type` VARCHAR(32) NOT NULL COMMENT 'vital_abnormal/sos_elderly/sos_nurse',
  `severity` TINYINT NOT NULL DEFAULT 1 COMMENT '1提示 2高危',
  `content` VARCHAR(512) DEFAULT NULL,
  `handled` TINYINT NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_elderly_type` (`elderly_id`, `alert_type`),
  KEY `idx_handled` (`handled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='安全警报';

CREATE TABLE `adverse_event` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `event_type` VARCHAR(32) NOT NULL COMMENT 'fall/pressure_ulcer/medication_error/dispute',
  `description` TEXT,
  `anonymous` TINYINT NOT NULL DEFAULT 1,
  `reporter_hash` VARCHAR(64) DEFAULT NULL,
  `nurse_id` BIGINT UNSIGNED DEFAULT NULL,
  `elderly_id` BIGINT UNSIGNED DEFAULT NULL,
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0待复盘 1已复盘',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_type` (`event_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='不良事件上报（匿名）';