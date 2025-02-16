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

func GetUserInfo(ctx context.Context, merchShopService service.MerchShopService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetUserInfo"
		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("userId not found in request context")
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, dto.ErrorResponse{Error: "userId not found in request context"})
			return
		}

		infoResponse, err := merchShopService.GetUserInfo(ctx, id)
		if err != nil {
			log.Error("error get user info:", tools.ErrAttr(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrorResponse{Error: "error get user info"})
			return
		}
		render.JSON(w, r, infoResponse)
	}
}
