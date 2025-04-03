package httpserver

import (
	"avito-winter-test/internal/http-server/handlers/authenticate"
	"avito-winter-test/internal/http-server/handlers/buy_item"
	"avito-winter-test/internal/http-server/handlers/get_user_info"
	"avito-winter-test/internal/http-server/handlers/send_coin"
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"avito-winter-test/internal/config"
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
		r.Post("/auth", authenticate.New(services.Auth, log).Authenticate(ctx))
		r.With(merch_shop_middleware.AuthMiddleware).Get("/info", get_user_info.New(services.MerchShop, log).GetUserInfo(ctx))
		r.With(merch_shop_middleware.AuthMiddleware).Post("/send-coin", send_coin.New(services.MerchShop, log).SendCoin(ctx))
		r.With(merch_shop_middleware.AuthMiddleware).Get("/buy/{item}", buy_item.New(services.MerchShop, log).BuyItem(ctx))
	})

	return &http.Server{
		Addr:              cfg.Address,
		BaseContext:       func(net.Listener) context.Context { return ctx },
		Handler:           router,
		ReadHeaderTimeout: 1 * time.Second,
	}
}
