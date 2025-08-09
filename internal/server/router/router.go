package router

import (
	"github.com/Nikolay961996/goferma/internal/storage"
	"github.com/go-chi/chi/v5"
)

func GofermaRouter(dbContext *storage.DBContext, secretKey string) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/api/user/register", registerHandler(dbContext, secretKey))
	router.Post("/api/user/login", loginHandler(dbContext, secretKey))
	router.Post("/api/user/orders", setOrdersHandler(dbContext, secretKey))
	router.Get("/api/user/balance", getBalanceHandler(dbContext, secretKey))
	router.Post("/api/user/balance/withdraw", withdrawBalanceHandler)
	router.Get("/api/user/withdrawals", showWithdrawalsBalanceHandler)

	return router
}
