package buy_item

import (
	"avito-winter-test/internal/models/dto"
	"avito-winter-test/internal/tools"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Handler struct {
	buyItemUsecase buyItemUsecase
	log            *slog.Logger
}

func New(buyItemUsecase buyItemUsecase, log *slog.Logger) *Handler {
	return &Handler{
		buyItemUsecase: buyItemUsecase,
		log:            log,
	}
}

func (h *Handler) BuyItem(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.BuyItem"
		w.Header().Set("Content-Type", "application/json")
		log := h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("userID not found in request context")
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, dto.ErrorResponse{Error: "userID not found in request context"})
			return
		}
		itemName := chi.URLParam(r, "item")

		err := h.buyItemUsecase.BuyItem(ctx, userID, itemName)
		if err != nil {
			log.Error("error buy item", tools.ErrAttr(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrorResponse{Error: "error buy item"})
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
