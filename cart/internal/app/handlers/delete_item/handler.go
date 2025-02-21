package delete_item

import (
	"log"
	"net/http"
	"route256/cart/internal/usecases/cart/dto"
	"strconv"
)

type cartClient interface {
	DeleteItem(userID dto.UserID, skuID dto.SkuID) error
}

type Handler struct {
	cartClient cartClient
}

func NewHandler(cartClient cartClient) *Handler {
	return &Handler{
		cartClient: cartClient,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	userID, err := strconv.ParseInt(req.PathValue("user_id"), 10, 64)
	if err != nil || userID == 0 {
		http.Error(w, "not valid userID", http.StatusBadRequest)
		return
	}
	skuID, err := strconv.ParseInt(req.PathValue("sku_id"), 10, 64)
	if err != nil || skuID == 0 {
		http.Error(w, "not valid skuID", http.StatusBadRequest)
		return
	}

	if err := h.cartClient.DeleteItem(dto.UserID(userID), dto.SkuID(skuID)); err != nil {
		log.Printf("delete item: %v", err)
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
