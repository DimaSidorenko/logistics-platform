package delete_item

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"route256/cart/internal/logger"
	"route256/cart/internal/usecases/cart/dto"
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
	vars := mux.Vars(req)

	userID, err := strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, "not valid userID", http.StatusBadRequest)
		return
	}

	skuID, err := strconv.ParseInt(vars["sku_id"], 10, 64)
	if err != nil || skuID == 0 {
		http.Error(w, "not valid skuID", http.StatusBadRequest)
		return
	}

	if err := h.cartClient.DeleteItem(dto.UserID(userID), dto.SkuID(skuID)); err != nil {
		logger.Warnw(req.Context(), "delete item: %v", err)
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
