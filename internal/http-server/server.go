package httpServer

import (
	"avito-winter-test/internal/config"
	"avito-winter-test/internal/http-server/handlers"
	merch_shop_middleware "avito-winter-test/internal/http-server/middleware"
	"avito-winter-test/internal/service"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net"
	"net/http"
	"time"
)

func NewServer(ctx context.Context, log *slog.Logger, cfg *config.Config, merchShopService service.MerchShopService) *http.Server {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/api", func(r chi.Router) {
		r.Post("/auth", handlers.Authenticate(ctx, merchShopService, log))
		r.With(merch_shop_middleware.AuthMiddleware).Get("/info", handlers.GetUserInfo(ctx, merchShopService, log))
		r.With(merch_shop_middleware.AuthMiddleware).Post("/send-coin", handlers.SendCoin(ctx, merchShopService, log))
		r.With(merch_shop_middleware.AuthMiddleware).Get("/buy/{item}", handlers.BuyItem(ctx, merchShopService, log))
	})

	return &http.Server{
		Addr:        cfg.Address,
		BaseContext: func(net.Listener) context.Context { return ctx },
		Handler:     router,
	}
}
