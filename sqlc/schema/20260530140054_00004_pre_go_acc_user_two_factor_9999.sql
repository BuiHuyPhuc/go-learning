-- +goose Up
CREATE TABLE IF NOT EXISTS `pre_go_acc_user_two_factor_9999` (
  `two_factor_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `two_factor_auth_type` tinyint unsigned NOT NULL COMMENT 'Auth type: 0-SMS, 1-Email, 2-App',
  `two_factor_auth_secret` varchar(255),
  `two_factor_phone` varchar(20),
  `two_factor_email` varchar(255),
  `two_factor_is_active` bit NOT NULL DEFAULT 1,
  `two_factor_created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `two_factor_updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  PRIMARY KEY (`two_factor_id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_auth_type` (`two_factor_auth_type`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='pre_go_acc_user_two_factor_9999';

-- +goose Down
DROP TABLE IF EXISTS `pre_go_acc_user_two_factor_9999`;
