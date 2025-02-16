package httpserver

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"avito-winter-test/internal/config"
	"avito-winter-test/internal/http-server/handlers"
	merch_shop_middleware "avito-winter-test/internal/http-server/middleware"
	"avito-winter-test/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Services struct {
	MerchShop service.MerchShopService
	Auth      service.AuthService
}

func NewServer(ctx context.Context, log *slog.Logger, cfg *config.Config, services *Services) *http.Server {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/api", func(r chi.Router) {
		r.Post("/auth", handlers.Authenticate(ctx, services.Auth, log))
		r.With(merch_shop_middleware.AuthMiddleware).Get("/info", handlers.GetUserInfo(ctx, services.MerchShop, log))
		r.With(merch_shop_middleware.AuthMiddleware).Post("/send-coin", handlers.SendCoin(ctx, services.MerchShop, log))
		r.With(merch_shop_middleware.AuthMiddleware).Get("/buy/{item}", handlers.BuyItem(ctx, services.MerchShop, log))
	})

	return &http.Server{
		Addr:              cfg.Address,
		BaseContext:       func(net.Listener) context.Context { return ctx },
		Handler:           router,
		ReadHeaderTimeout: 1 * time.Second,
	}
}
