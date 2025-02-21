package delete_user

import (
	"log"
	"net/http"
	"route256/cart/internal/usecases/cart/dto"
	"strconv"
)

type cartClient interface {
	DeleteUser(userID dto.UserID) error
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

	if err := h.cartClient.DeleteUser(dto.UserID(userID)); err != nil {
		log.Printf("delete user: %v", err)
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
