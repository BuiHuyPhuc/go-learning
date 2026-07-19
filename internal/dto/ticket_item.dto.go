package dto

type TicketItemRequest struct {
	TicketId string `json:"ticket_id"`
}

type TicketItemResponse struct {
	TicketId       int    `json:"ticket_id"`
	TicketName     string `json:"ticket_name"`
	StockAvailable int    `json:"stock_available"`
	SotckInitial   int    `json:"stock_initial"`
}

type OrderRequest struct {
	TicketId int    `json:"ticket_id" validate:"required"`
	UserId   int    `json:"user_id"`
	Quantity int    `json:"quantity" validate:"gte=1"`
	Notes    string `json:"notes"`
}
