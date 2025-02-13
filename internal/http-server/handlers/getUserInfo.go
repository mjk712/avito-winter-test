package handlers

import (
	"avito-winter-test/internal/service"
	"avito-winter-test/internal/tools"
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

func GetUserInfo(ctx context.Context, merchShopService service.MerchShopService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetUserInfo"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := r.Context().Value("userID").(int)

		infoResponse, err := merchShopService.GetUserInfo(ctx, id)
		if err != nil {
			log.Error("error get user info:", tools.ErrAttr(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		render.JSON(w, r, infoResponse)
	}
}
