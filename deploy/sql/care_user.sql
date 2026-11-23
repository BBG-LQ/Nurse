CREATE DATABASE IF NOT EXISTS care_user DEFAULT CHARSET utf8mb4;
USE care_user;

CREATE TABLE `sys_user` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `mobile` VARCHAR(20) NOT NULL,
  `password_hash` VARCHAR(128) NOT NULL COMMENT '本地测试可用明文/固定hash',
  `real_name` VARCHAR(64) DEFAULT NULL,
  `role` VARCHAR(16) NOT NULL COMMENT 'elderly/family/nurse/admin',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '1正常 0停用',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_mobile` (`mobile`),
  KEY `idx_role` (`role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统账号';

CREATE TABLE `family_binding` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `family_user_id` BIGINT UNSIGNED NOT NULL COMMENT '家属 sys_user.id',
  `elderly_id` BIGINT UNSIGNED NOT NULL COMMENT 'elderly_archive.id',
  `relation` VARCHAR(16) DEFAULT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_family_elderly` (`family_user_id`, `elderly_id`),
  KEY `idx_elderly` (`elderly_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='家属绑定';