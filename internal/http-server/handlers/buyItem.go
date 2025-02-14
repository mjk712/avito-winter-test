package handlers

import (
	"avito-winter-test/internal/service"
	"avito-winter-test/internal/tools"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

func BuyItem(ctx context.Context, merchShopService service.MerchShopService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.BuyItem"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userId := r.Context().Value("userID").(int)
		itemName := chi.URLParam(r, "item")
		
		err := merchShopService.BuyItem(ctx, userId, itemName)
		if err != nil {
			log.Error("error buy item", tools.ErrAttr(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)

	}
}
