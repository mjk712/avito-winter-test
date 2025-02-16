package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"avito-winter-test/internal/models/dto"
	"avito-winter-test/internal/service"
	"avito-winter-test/internal/tools"
)

func SendCoin(ctx context.Context, merchShopService service.MerchShopService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.SendCoin"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		fromUserID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("userId not found in request context")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var req dto.SendCoinRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("error decode json", tools.ErrAttr(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := merchShopService.SendCoin(ctx, fromUserID, req)
		if err != nil {
			log.Error("error sending coin", tools.ErrAttr(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
