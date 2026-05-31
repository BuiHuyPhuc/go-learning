-- +goose Up
ALTER TABLE `pre_go_acc_user_base_9999`
ADD COLUMN `is_two_factor_enabled` bit NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE `pre_go_acc_user_base_9999`
DROP COLUMN `is_two_factor_enabled`;
