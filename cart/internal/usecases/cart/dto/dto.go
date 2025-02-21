package dto

type Product struct {
	Name  string `json:"name"`
	Price int32  `json:"price"`
	Sku   int64  `json:"sku"`
}

type UserID int64
type SkuID int64

type Item struct {
	Sku   int64
	Name  string
	Count uint32
	Price uint32
}

type GetItemsResponse struct {
	Items      []Item
	TotalPrice uint32
}
