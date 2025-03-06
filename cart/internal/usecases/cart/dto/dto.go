package dto

type UserID int64
type SkuID int64

type Item struct {
	Sku   SkuID
	Name  string
	Count uint32
	Price uint32
}

type GetItemsResponse struct {
	Items      []Item
	TotalPrice uint32
}
