package add_item

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	dto2 "route256/cart/internal/usecases/cart/dto"
	"strconv"
)

type cartClient interface {
	AddItem(userID dto2.UserID, skuID dto2.SkuID, quantity uint32) error
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

	data, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "cannot read request body", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	var request AddItemRequest
	if err = json.Unmarshal(data, &request); err != nil {
		http.Error(w, "cannot parse request body", http.StatusInternalServerError)
		return
	}
	if request.Count == 0 {
		http.Error(w, "count should be positive", http.StatusBadRequest)
		return
	}

	if err := h.cartClient.AddItem(dto2.UserID(userID), dto2.SkuID(skuID), request.Count); err != nil {
		log.Printf("add item: %v", err)
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	w.WriteHeader(http.StatusOK)
}
