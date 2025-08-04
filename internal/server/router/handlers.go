package router

import (
	"errors"
	"github.com/Nikolay961996/goferma/internal/models"
	"github.com/Nikolay961996/goferma/internal/services"
	"github.com/Nikolay961996/goferma/internal/utils"
	"net/http"
)

/*
200 — пользователь успешно зарегистрирован и аутентифицирован;
400 — неверный формат запроса;
409 — логин уже занят;
500 — внутренняя ошибка сервера.
*/
func registerHandler(w http.ResponseWriter, r *http.Request) {
	utils.Log.Info("registerHandler")

	loginModel, err := services.ReadAuthModel(r)
	utils.Log.Info("User ", loginModel.Login)
	if err != nil {
		var formatError *models.FormatError
		if errors.As(err, &formatError) {
			utils.Log.Error("error reading request model:", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		utils.Log.Error("error reading request model:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := services.BuildJWTToken(1, "")
	if err != nil {
		utils.Log.Error("error building JWT token")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", token)

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
