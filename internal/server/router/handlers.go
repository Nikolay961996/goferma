package router

import (
	"github.com/Nikolay961996/goferma/internal/utils"
	"net/http"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	utils.Log.Info("registerHandler")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	utils.Log.Info("loginHandler")

}

func setOrdersHandler(w http.ResponseWriter, r *http.Request) {
	utils.Log.Info("setOrdersHandler")

}

func getBalanceHandler(w http.ResponseWriter, r *http.Request) {
	utils.Log.Info("getBalanceHandler")

}

func withdrawBalanceHandler(w http.ResponseWriter, r *http.Request) {
	utils.Log.Info("withdrawBalanceHandler")

}

func showWithdrawalsBalanceHandler(w http.ResponseWriter, r *http.Request) {
	utils.Log.Info("showWithdrawalsBalanceHandler")

}
