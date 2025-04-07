package httpserver

import (
	"avito-winter-test/internal/config"
	"avito-winter-test/internal/http-server/handlers/authenticate"
	"avito-winter-test/internal/http-server/handlers/buy_item"
	"avito-winter-test/internal/http-server/handlers/get_user_info"
	"avito-winter-test/internal/http-server/handlers/send_coin"
	merch_shop_middleware "avito-winter-test/internal/http-server/middleware"
	"avito-winter-test/internal/storage"
	"avito-winter-test/internal/usecases/authenticate_user"
	buy_item_usecase "avito-winter-test/internal/usecases/buy_item"
	get_user_info_usecase "avito-winter-test/internal/usecases/get_user_info"
	"avito-winter-test/internal/usecases/send_coins"

	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServer(ctx context.Context, log *slog.Logger, cfg *config.Config, repo storage.Storage) *http.Server {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/api", func(r chi.Router) {
		r.Post("/auth", authenticate.New(authenticate_user.New(repo), log).Authenticate(ctx))
		r.With(merch_shop_middleware.AuthMiddleware).Get("/info", get_user_info.New(get_user_info_usecase.New(repo), log).GetUserInfo(ctx))
		r.With(merch_shop_middleware.AuthMiddleware).Post("/send-coin", send_coin.New(send_coins.New(repo), log).SendCoin(ctx))
		r.With(merch_shop_middleware.AuthMiddleware).Get("/buy/{item}", buy_item.New(buy_item_usecase.New(repo), log).BuyItem(ctx))
	})

	return &http.Server{
		Addr:              cfg.Address,
		BaseContext:       func(net.Listener) context.Context { return ctx },
		Handler:           router,
		ReadHeaderTimeout: 1 * time.Second,
	}
}
