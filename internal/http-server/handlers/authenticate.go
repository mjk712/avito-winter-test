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

// Authenticate Авторизация пользователя
// @Summary Авторизация пользователя

func Authenticate(ctx context.Context, authService service.AuthService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Authenticate"

		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		// Читаем request
		var req dto.AuthRequest

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("error reading request", tools.ErrAttr(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Проводим аутентификацию
		token, err := authService.Authenticate(ctx, req)
		if err != nil {
			log.Error("error authenticating", tools.ErrAttr(err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		render.JSON(w, r, dto.AuthResponse{Token: token})
	}
}
