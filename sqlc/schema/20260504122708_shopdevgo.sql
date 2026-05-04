-- +goose Up
CREATE TABLE IF NOT EXISTS `pre_go_acc_user_9999` (
  `user_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'User  ID',
  `user_account` varchar(255) NOT NULL COMMENT 'User account',
  `user_nickname` varchar(255) DEFAULT NULL COMMENT 'User nickname',
  `user_avatar` varchar(255) DEFAULT NULL COMMENT 'User avatar',
  `user_state` tinyint unsigned NOT NULL COMMENT 'User state: 0-Locked, 1-Activated, 2-Not Activated',
  `user_mobile` varchar(20) DEFAULT NULL COMMENT 'Mobile phone number',
  `user_gender` tinyint unsigned DEFAULT NULL COMMENT 'User gender: 0-Secret, 1-Male, 2-Female',
  `user_birthday` date DEFAULT NULL COMMENT 'User birthday',
  `user_email` varchar(255) DEFAULT NULL COMMENT 'User email address',
  `user_is_authentication` tinyint unsigned NOT NULL COMMENT 'User gender: 0-Not Authentication, 1-Authentication',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record creation time',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Record update time',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `unique_user_account` (`user_account`),
  KEY `idx_user_mobile` (`user_mobile`),
  KEY `idx_user_email` (`user_email`),
  KEY `idx_user_state` (`user_state`),
  KEY `idx_user_is_authentication` (`user_is_authentication`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='pre_go_acc_user_9999';

CREATE TABLE IF NOT EXISTS `pre_go_acc_user_base_9999` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_account` varchar(255) NOT NULL,
  `user_password` varchar(255) NOT NULL,
  `user_salt` varchar(255) NOT NULL,
  `user_login_time` timestamp NULL DEFAULT NULL,
  `user_logout_time` timestamp NULL DEFAULT NULL,
  `user_login_ip` varchar(45) DEFAULT NULL,
  `user_created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `user_updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `unique_user_account` (`user_account`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='pre_go_acc_user_base_9999';

CREATE TABLE IF NOT EXISTS `pre_go_acc_user_verify_9999` (
  `verify_id` int(11) NOT NULL AUTO_INCREMENT,
  `verify_otp` varchar(6) NOT NULL,
  `verify_key` varchar(255) NOT NULL,
  `verify_key_hash` varchar(255) NOT NULL,
  `verify_type` int(11) DEFAULT '1',
  `is_verified` int(11) DEFAULT '0',
  `is_deleted` int(11) DEFAULT '0',
  `verify_created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `verify_updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`verify_id`),
  UNIQUE KEY `unique_verify_key` (`verify_key`),
  KEY `idx_verify_otp` (`verify_otp`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='pre_go_acc_user_verify_9999';

-- +goose Down
DROP TABLE IF EXISTS `pre_go_acc_user_9999`;
DROP TABLE IF EXISTS `pre_go_acc_user_base_9999`;
DROP TABLE IF EXISTS `pre_go_acc_user_verify_9999`;
