package dto

type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "new"
	OrderStatusAwaitingPayment OrderStatus = "awaiting payment"
	OrderStatusFailed          OrderStatus = "failed"
	OrderStatusPayed           OrderStatus = "payed"
	OrderStatusCancelled       OrderStatus = "cancelled"
)

type Order struct {
	OrderID int64 `json:"order_id"`
	Status  OrderStatus
	User    int64
	Items   []Item
}

type Item struct {
	SKU   int64
	Count uint32
}
