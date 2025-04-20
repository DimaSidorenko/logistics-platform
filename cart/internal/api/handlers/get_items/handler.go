package get_items

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"route256/cart/internal/logger"
	cartDto "route256/cart/internal/usecases/cart/dto"
)

type cartClient interface {
	GetItems(ctx context.Context, userID cartDto.UserID) (cartDto.GetItemsResponse, error)
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

	items, err := h.cartClient.GetItems(req.Context(), cartDto.UserID(userID))
	if err != nil {
		if errors.Is(err, cartDto.ErrUserNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		logger.Errorw(req.Context(), "get items: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(convertResponse(items))
	if err != nil {
		http.Error(w, fmt.Sprintf("encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func convertResponse(response cartDto.GetItemsResponse) (out GetItemsResponse) {
	out.Items = make([]Item, len(response.Items))
	for i := range response.Items {
		out.Items[i] = Item{
			Sku:   int64(response.Items[i].Sku),
			Name:  response.Items[i].Name,
			Count: response.Items[i].Count,
			Price: response.Items[i].Price,
		}
	}
	out.TotalPrice = response.TotalPrice

	return out
}
