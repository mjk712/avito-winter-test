package send_coin

import (
	"avito-winter-test/internal/models/dto"
	"avito-winter-test/internal/tools"
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Handler struct {
	sendCoinUsecase sendCoinUsecase
	log             *slog.Logger
}

func New(sendCoinUsecase sendCoinUsecase, log *slog.Logger) *Handler {
	return &Handler{
		sendCoinUsecase: sendCoinUsecase,
		log:             log,
	}
}

func (h *Handler) SendCoin(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.SendCoin"
		w.Header().Set("Content-Type", "application/json")
		log := h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		fromUserID, ok := r.Context().Value("userID").(int)
		if !ok {
			log.Error("userId not found in request context")
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, dto.ErrorResponse{Error: "userId not found in request context"})
			return
		}

		var req dto.SendCoinRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("error decode json", tools.ErrAttr(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrorResponse{Error: "error decode json"})
			return
		}

		err := h.sendCoinUsecase.SendCoins(ctx, fromUserID, req)
		if err != nil {
			log.Error("error sending coin", tools.ErrAttr(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, dto.ErrorResponse{Error: "error sending coin"})
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
