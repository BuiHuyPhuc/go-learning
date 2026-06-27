-- +goose Up
CREATE TABLE IF NOT EXISTS `pre_go_ticket_item_9999` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `description` text,
  `stock_initial` int(11) NOT NULL,
  `stock_available` int(11) NOT NULL,
  `is_stock_prepared` boolean NOT NULL DEFAULT 0,
  `price_original` bigint(20) NOT NULL,
  `price_flash` bigint(20) NOT NULL,
  `sale_start_time` datetime NOT NULL,
  `sale_end_time` datetime NOT NULL,
  `status` int(11) NOT NULL DEFAULT 0,
  `acctivity_id` bigint(20) NOT NULL,

  `created_at` datetime NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_start_time` (`sale_start_time`),
  INDEX `idx_end_time` (`sale_end_time`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='pre_go_ticket_item_9999';

-- +goose Down
DROP TABLE IF EXISTS `pre_go_ticket_item_9999`;
