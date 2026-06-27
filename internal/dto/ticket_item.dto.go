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
