CREATE DATABASE IF NOT EXISTS care_vital DEFAULT CHARSET utf8mb4;
USE care_vital;

CREATE TABLE `vital_record` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `elderly_id` BIGINT UNSIGNED NOT NULL,
  `heart_rate` SMALLINT DEFAULT NULL,
  `systolic` SMALLINT DEFAULT NULL,
  `diastolic` SMALLINT DEFAULT NULL,
  `spo2` DECIMAL(4,1) DEFAULT NULL,
  `posture` TINYINT DEFAULT NULL COMMENT '1站立 2卧床 3跌倒',
  `source` VARCHAR(16) NOT NULL DEFAULT 'keyboard',
  `record_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_elderly_time` (`elderly_id`, `record_time`),
  KEY `idx_posture` (`posture`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='体征记录（键盘模拟）';