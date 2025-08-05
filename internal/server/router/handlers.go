package router

import (
	"errors"
	"github.com/Nikolay961996/goferma/internal/models"
	"github.com/Nikolay961996/goferma/internal/services"
	"github.com/Nikolay961996/goferma/internal/storage"
	"github.com/Nikolay961996/goferma/internal/utils"
	"net/http"
)

/*
200 — пользователь успешно зарегистрирован и аутентифицирован;
400 — неверный формат запроса;
409 — логин уже занят;
500 — внутренняя ошибка сервера.
*/
func registerHandler(db *storage.DBContext, secretKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.Log.Info("registerHandler")

		loginModel, err := services.ReadAuthModel(r)
		if err != nil {
			errorHandler(err, w)
			return
		}

		err = services.CreateUser(db, loginModel.Login, loginModel.Password)
		if err != nil {
			errorHandler(err, w)
			return
		}

		token, err := services.AuthUser(db, secretKey, loginModel.Login, loginModel.Password)
		if err != nil {
			errorHandler(err, w)
			return
		}
		w.Header().Set("Authorization", token)

	}
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

/*
400 FormatError
401 LoginPasswordError
409 AlreadyExistError
500 Internal
*/
func errorHandler(err error, w http.ResponseWriter) {
	var formatError *models.FormatError
	if errors.As(err, &formatError) {
		utils.Log.Error("error format request model:", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var alreadyExistError *models.AlreadyExistError
	if errors.As(err, &alreadyExistError) {
		utils.Log.Error("error already exist:", err.Error())
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	var loginPasswordError *models.LoginPasswordError
	if errors.As(err, &loginPasswordError) {
		utils.Log.Error("error login/password pair:", err.Error())
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	utils.Log.Error("internal error: ", err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	return
}
