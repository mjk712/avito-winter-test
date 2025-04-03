package authenticate

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
	authService authService
	log         *slog.Logger
}

func New(authService authService, log *slog.Logger) *Handler {
	return &Handler{
		authService: authService,
		log:         log,
	}
}

func (h *Handler) Authenticate(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Authenticate"

		w.Header().Set("Content-Type", "application/json")
		log := h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		// Читаем request
		var req dto.AuthRequest

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("error reading request", tools.ErrAttr(err))
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, dto.ErrorResponse{Error: "error reading request"})
			return
		}

		// Проводим аутентификацию
		token, err := h.authService.Authenticate(ctx, req)
		if err != nil {
			log.Error("error authenticating", tools.ErrAttr(err))
			w.WriteHeader(http.StatusUnauthorized)
			render.JSON(w, r, dto.ErrorResponse{Error: "error authenticating"})
			return
		}
		render.JSON(w, r, dto.AuthResponse{Token: token})
	}
}
