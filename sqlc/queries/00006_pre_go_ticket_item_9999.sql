-- name: GetTicketItemById :one
SELECT id, name, stock_initial, stock_available
FROM `pre_go_ticket_item_9999`
WHERE id = ?;
