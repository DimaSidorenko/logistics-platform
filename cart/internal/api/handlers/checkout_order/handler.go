package checkout_order

import (
	"encoding/json"
	"log"
	"net/http"
	cartDto "route256/cart/internal/usecases/cart/dto"
	"strconv"
)

type cartClient interface {
	Checkout(userID cartDto.UserID) (orderID int64, err error)
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

	orderID, err := h.cartClient.Checkout(cartDto.UserID(userID))
	if err != nil {
		log.Printf("checkout failed: %v", err)
		http.Error(w, "checkout failed", http.StatusInternalServerError)
		return
	}

	response := CheckoutResponse{
		OrderID: orderID,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode response: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
