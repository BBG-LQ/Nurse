CREATE DATABASE IF NOT EXISTS care_elderly DEFAULT CHARSET utf8mb4;
USE care_elderly;

CREATE TABLE `elderly_archive` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `real_name` VARCHAR(64) NOT NULL,
  `gender` TINYINT NOT NULL DEFAULT 0 COMMENT '0未知 1男 2女',
  `birth_date` DATE DEFAULT NULL,
  `id_card_masked` VARCHAR(32) DEFAULT NULL,
  `mobile` VARCHAR(20) DEFAULT NULL,
  `address` VARCHAR(255) DEFAULT NULL COMMENT '本地验证兼作位置标签',
  `emergency_contact` VARCHAR(64) DEFAULT NULL,
  `emergency_mobile` VARCHAR(20) DEFAULT NULL,
  `chronic_disease` JSON DEFAULT NULL COMMENT '慢病史 JSON 数组',
  `care_level` TINYINT NOT NULL DEFAULT 1 COMMENT '护理等级 1-5',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '1正常 0停用',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_mobile` (`mobile`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='老人基础档案';

CREATE TABLE `elderly_medication` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `elderly_id` BIGINT UNSIGNED NOT NULL,
  `nurse_id` BIGINT UNSIGNED NOT NULL,
  `drug_name` VARCHAR(128) NOT NULL,
  `dosage` VARCHAR(64) DEFAULT NULL,
  `frequency` VARCHAR(32) DEFAULT NULL,
  `remind_times` JSON DEFAULT NULL COMMENT '["08:00","12:00"]',
  `start_date` DATE DEFAULT NULL,
  `end_date` DATE DEFAULT NULL,
  `status` TINYINT NOT NULL DEFAULT 1,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_elderly` (`elderly_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用药方案';