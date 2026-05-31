-- name: EnableTwoFactorTypeEmail :exec
INSERT INTO `pre_go_acc_user_two_factor_9999` (user_id, two_factor_auth_type, two_factor_email, two_factor_auth_secret, two_factor_is_active, two_factor_created_at, two_factor_updated_at)
VALUES (?, ?, ?, "OTP", 0, now(), now());

-- name: DisableTwoFactor :exec
UPDATE `pre_go_acc_user_two_factor_9999`
SET two_factor_is_active = 0,
    two_factor_updated_at = now()
WHERE user_id = ? AND two_factor_auth_type = ?;

-- name: UpdateTwoFactorStatus :exec
UPDATE `pre_go_acc_user_two_factor_9999`
SET two_factor_is_active = b'1',
    two_factor_updated_at = now()
WHERE user_id = ? AND two_factor_auth_type = ? AND two_factor_is_active = b'0';

-- name: VerifyTwoFactor :one
SELECT COUNT(*)
FROM `pre_go_acc_user_two_factor_9999`
WHERE user_id = ? AND two_factor_auth_type = ? AND two_factor_is_active = b'1';

-- name: GetTwoFactorStatus :one
SELECT two_factor_is_active
FROM `pre_go_acc_user_two_factor_9999`
WHERE user_id = ? AND two_factor_auth_type = ?;

-- name: IsTwoFactorEnabled :one
SELECT COUNT(*)
FROM `pre_go_acc_user_two_factor_9999`
WHERE user_id = ? AND two_factor_is_active = b'1';

-- name: AddOrUpdatePhone :exec
INSERT INTO `pre_go_acc_user_two_factor_9999` (user_id, two_factor_phone, two_factor_is_active)
VALUES (?, ?, b'1')
ON DUPLICATE KEY UPDATE two_factor_phone = ?, two_factor_updated_at = now();

-- name: AddOrUpdateEmail :exec
INSERT INTO `pre_go_acc_user_two_factor_9999` (user_id, two_factor_email, two_factor_is_active)
VALUES (?, ?, b'1')
ON DUPLICATE KEY UPDATE two_factor_email = ?, two_factor_updated_at = now();

-- name: GetTwoFactorMethods :many
SELECT two_factor_id, user_id, two_factor_auth_type, two_factor_auth_secret, two_factor_phone, two_factor_email, two_factor_is_active, two_factor_created_at, two_factor_updated_at
FROM `pre_go_acc_user_two_factor_9999`
WHERE user_id = ?;

-- name: ReactivateTwoFactor :exec
UPDATE `pre_go_acc_user_two_factor_9999`
SET two_factor_is_active = b'1',
    two_factor_updated_at = now()
WHERE user_id = ? AND two_factor_auth_type = ?;

-- name: RemoveTwoFactor :exec
DELETE FROM `pre_go_acc_user_two_factor_9999`
WHERE user_id = ? AND two_factor_auth_type = ?;

-- name: CountActiveTwoFactorMethods :one
SELECT COUNT(*)
FROM `pre_go_acc_user_two_factor_9999`
WHERE user_id = ? AND two_factor_is_active = b'1';

-- name: GetTwoFactorMethodByID :one
SELECT two_factor_id, user_id, two_factor_auth_type, two_factor_auth_secret, two_factor_phone, two_factor_email, two_factor_is_active, two_factor_created_at, two_factor_updated_at
FROM `pre_go_acc_user_two_factor_9999`
WHERE user_id = ?;

-- name: GetTwoFactorMethodByIDAndType :one
SELECT two_factor_id, user_id, two_factor_auth_type, two_factor_auth_secret, two_factor_phone, two_factor_email, two_factor_is_active, two_factor_created_at, two_factor_updated_at
FROM `pre_go_acc_user_two_factor_9999`
WHERE user_id = ? AND two_factor_auth_type = ?;
