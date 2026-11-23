CREATE DATABASE IF NOT EXISTS care_nurse DEFAULT CHARSET utf8mb4;
USE care_nurse;

CREATE TABLE `nurse_profile` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `real_name` VARCHAR(64) NOT NULL,
  `mobile` VARCHAR(20) NOT NULL,
  `license_no` VARCHAR(64) NOT NULL,
  `license_expire` DATE DEFAULT NULL,
  `skill_tags` JSON DEFAULT NULL COMMENT '擅长护理项目标签',
  `service_area` VARCHAR(128) DEFAULT NULL COMMENT '虚拟服务区域，用于本地调度匹配',
  `current_load` INT NOT NULL DEFAULT 0,
  `score_avg` DECIMAL(3,1) NOT NULL DEFAULT 5.0,
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '1可接单 0冻结',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_mobile` (`mobile`),
  KEY `idx_status` (`status`),
  KEY `idx_expire` (`license_expire`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='护士执业档案';

CREATE TABLE `training_record` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `nurse_id` BIGINT UNSIGNED NOT NULL,
  `course_name` VARCHAR(128) NOT NULL,
  `exam_month` CHAR(7) NOT NULL COMMENT 'YYYY-MM',
  `score` TINYINT DEFAULT NULL,
  `attempt_count` TINYINT NOT NULL DEFAULT 1,
  `passed` TINYINT NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_nurse_month` (`nurse_id`, `exam_month`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='培训考核记录';