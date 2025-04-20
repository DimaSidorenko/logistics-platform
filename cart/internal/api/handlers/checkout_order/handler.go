package checkout_order

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"route256/cart/internal/logger"
	"route256/cart/internal/tracing"
	cartDto "route256/cart/internal/usecases/cart/dto"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type cartClient interface {
	Checkout(ctx context.Context, userID cartDto.UserID) (orderID int64, err error)
}

type Handler struct {
	cartClient cartClient
}

func NewHandler(cartClient cartClient) *Handler {
	return &Handler{
		cartClient: cartClient,
	}
}

//nolint:unused
func injectTraceContext(ctx context.Context, req *http.Request) {
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	userID, err := strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, "not valid userID", http.StatusBadRequest)
		return
	}

	ctx, span := tracing.StartFromContext(req.Context(), "handler /checkout/order")
	defer span.End()

	orderID, err := h.cartClient.Checkout(ctx, cartDto.UserID(userID))
	if err != nil {
		logger.Warnw(req.Context(), "checkout order failed: %v", err)

		if errors.Is(err, cartDto.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if errors.Is(err, cartDto.ErrFailedToReserveStocks) {
			http.Error(w, "failed to reserve stocks", http.StatusPreconditionFailed)
			return
		}

		http.Error(w, "checkout failed", http.StatusInternalServerError)
		return
	}

	response := CheckoutResponse{
		OrderID: orderID,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Errorw(req.Context(), "checkout order failed: %v", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
