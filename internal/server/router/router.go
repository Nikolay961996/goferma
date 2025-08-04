package router

import (
	"github.com/Nikolay961996/goferma/internal/storage"
	"github.com/go-chi/chi/v5"
)

func GofermaRouter(dbContext *storage.DBContext) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/api/user/register", registerHandler(dbContext))
	router.Post("/api/user/login", loginHandler)
	router.Post("/api/user/orders", setOrdersHandler)
	router.Get("/api/user/balance", getBalanceHandler)
	router.Post("/api/user/balance/withdraw", withdrawBalanceHandler)
	router.Get("/api/user/withdrawals", showWithdrawalsBalanceHandler)

	return router
}
