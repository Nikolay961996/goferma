package router

import (
	"github.com/go-chi/chi/v5"
)

func GofermaRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/api/user/register", registerHandler)
	router.Post("/api/user/login", loginHandler)
	router.Post("/api/user/orders", setOrdersHandler)
	router.Get("/api/user/balance", getBalanceHandler)
	router.Post("/api/user/balance/withdraw", withdrawBalanceHandler)
	router.Get("/api/user/withdrawals", showWithdrawalsBalanceHandler)

	return router
}
