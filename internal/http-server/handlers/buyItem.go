package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"avito-winter-test/internal/service"
	"avito-winter-test/internal/tools"
)

func BuyItem(ctx context.Context, merchShopService service.MerchShopService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.BuyItem"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("userID not found in request context")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		itemName := chi.URLParam(r, "item")

		err := merchShopService.BuyItem(ctx, userID, itemName)
		if err != nil {
			log.Error("error buy item", tools.ErrAttr(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
