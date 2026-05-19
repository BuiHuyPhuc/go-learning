-- +goose Up
CREATE TABLE IF NOT EXISTS `pre_go_acc_user_info_9999` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_account` varchar(255) NOT NULL,
  `user_nickname` varchar(255),
  `user_avatar` varchar(255),
  `user_state` tinyint unsigned NOT NULL COMMENT 'User state: 0-Locked, 1-Activated, 2-Not Activated',
  `user_mobile` varchar(20),
  `user_gender` tinyint unsigned,
  `user_birthday` date,
  `user_email` varchar(255),
  `user_is_authentication` tinyint unsigned NOT NULL COMMENT 'Authentication status: 0-Not Authentication, 1-Pending, 2-Authentication',

  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `unique_user_account` (`user_account`),
  INDEX `idx_user_mobile` (`user_mobile`),
  INDEX `idx_user_email` (`user_email`),
  INDEX `idx_user_state` (`user_state`),
  INDEX `idx_user_is_authentication` (`user_is_authentication`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='pre_go_acc_user_info_9999';

-- +goose Down
DROP TABLE IF EXISTS `pre_go_acc_user_base_9999`;
