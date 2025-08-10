package router

import (
	"github.com/Nikolay961996/goferma/internal/storage"
	"github.com/go-chi/chi/v5"
)

func GofermaRouter(dbContext *storage.DBContext, secretKey string) *chi.Mux {
	router := chi.NewRouter()
	router.Use(WithLogger)

	router.Group(func(router chi.Router) {
		router.Post("/api/user/register", registerHandler(dbContext, secretKey))
		router.Post("/api/user/login", loginHandler(dbContext, secretKey))
	})

	router.Group(func(router chi.Router) {
		router.Use(WithAuth(secretKey))

		router.Post("/api/user/orders", setOrdersHandler(dbContext))
		router.Get("/api/user/orders", getOrdersHandler(dbContext))
		router.Get("/api/user/balance", getBalanceHandler(dbContext))
		router.Post("/api/user/balance/withdraw", withdrawBalanceHandler(dbContext))
		router.Get("/api/user/withdrawals", showWithdrawalsBalanceHandler(dbContext))
	})

	return router
}
