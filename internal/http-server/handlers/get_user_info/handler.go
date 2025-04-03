package get_user_info

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"avito-winter-test/internal/models/dto"
	"avito-winter-test/internal/tools"
)

type Handler struct {
	merchShopService merchService
	log              *slog.Logger
}

func New(merchShopService merchService, log *slog.Logger) *Handler {
	return &Handler{
		merchShopService: merchShopService,
		log:              log,
	}
}

func (h *Handler) GetUserInfo(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetUserInfo"
		w.Header().Set("Content-Type", "application/json")
		log := h.log.With(
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

		infoResponse, err := h.merchShopService.GetUserInfo(ctx, id)
		if err != nil {
			log.Error("error get user info:", tools.ErrAttr(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrorResponse{Error: "error get user info"})
			return
		}
		render.JSON(w, r, infoResponse)
	}
}
