package router

import (
	"encoding/json"
	"errors"
	"fmt"
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
		authHandler(db, secretKey, w, r, true)
	}
}

/*
200 — пользователь успешно аутентифицирован;
400 — неверный формат запроса;
401 — неверная пара логин/пароль;
500 — внутренняя ошибка сервера.
*/
func loginHandler(db *storage.DBContext, secretKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHandler(db, secretKey, w, r, false)
	}
}

/*
200 — номер заказа уже был загружен этим пользователем;
202 — новый номер заказа принят в обработку;
400 — неверный формат запроса;
401 — пользователь не аутентифицирован;
409 — номер заказа уже был загружен другим пользователем;
422 — неверный формат номера заказа;
500 — внутренняя ошибка сервера.
*/
func setOrdersHandler(db *storage.DBContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(w, r)
		if userID == 0 {
			return
		}
		orderNumber, err := services.GetOrderNumber(r.Header.Get("content-type"), r.Body)
		if err != nil {
			errorHandler(err, w)
			return
		}

		createdNew, err := services.RegisterOrder(db, orderNumber, userID, models.New, 0)
		if err != nil {
			errorHandler(err, w)
			return
		}
		if createdNew {
			w.WriteHeader(http.StatusAccepted)
		}
	}
}

/*
200 — успешная обработка запроса.
204 — нет данных для ответа.
401 — пользователь не авторизован.
500 — внутренняя ошибка сервера.
*/
func getOrdersHandler(db *storage.DBContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		userID := getUserID(w, r)
		if userID == 0 {
			return
		}
		orders, err := db.GetUserOrders(userID)
		if err != nil {
			errorHandler(err, w)
			return
		}
		if orders == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		writeToResponse(w, orders)
	}
}

/*
200 — успешная обработка запроса.
401 — пользователь не авторизован.
500 — внутренняя ошибка сервера.
*/
func getBalanceHandler(db *storage.DBContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(w, r)
		if userID == 0 {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		data, err := services.GetUserBalance(db, userID)
		if err != nil {
			errorHandler(err, w)
			return
		}
		writeToResponse(w, data)
	}
}

/*
200 — успешная обработка запроса;
401 — пользователь не авторизован;
402 — на счету недостаточно средств;
422 — неверный номер заказа;
500 — внутренняя ошибка сервера.
*/
func withdrawBalanceHandler(db *storage.DBContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(w, r)
		if userID == 0 {
			return
		}
		model, err := services.ReadWithdrawnModel(r.Header.Get("content-type"), r.Body)
		if err != nil {
			errorHandler(err, w)
			return
		}
		err = services.Withdrawn(db, userID, model.Sum, model.Order)
		if err != nil {
			errorHandler(err, w)
			return
		}
	}
}

/*
200 — успешная обработка запроса.
204 — нет ни одного списания.
401 — пользователь не авторизован.
500 — внутренняя ошибка сервера.
*/
func showWithdrawalsBalanceHandler(db *storage.DBContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		userID := getUserID(w, r)
		if userID == 0 {
			return
		}
		withdrawnHistory, err := db.GerUserWithdrawnHistory(userID)
		if err != nil {
			errorHandler(err, w)
			return
		}
		if withdrawnHistory == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		writeToResponse(w, withdrawnHistory)
	}
}

func authHandler(db *storage.DBContext, secretKey string, w http.ResponseWriter, r *http.Request, isRegistration bool) {
	loginModel, err := services.ReadAuthModel(r.Header.Get("content-type"), r.Body)
	if err != nil {
		errorHandler(err, w)
		return
	}

	if isRegistration {
		err = services.CreateUser(db, loginModel.Login, loginModel.Password)
		if err != nil {
			errorHandler(err, w)
			return
		}
	}

	token, err := services.AuthUser(db, secretKey, loginModel.Login, loginModel.Password)
	if err != nil {
		errorHandler(err, w)
		return
	}
	w.Header().Set("Authorization", token)
}

func getUserID(w http.ResponseWriter, r *http.Request) int64 {
	userID, ok := r.Context().Value(models.UserIDKey).(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return 0
	}
	return userID
}

func writeToResponse(w http.ResponseWriter, data any) {
	resp, err := json.Marshal(data)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("Error marshalling body: %v", err))
		errorHandler(err, w)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("Error write body: %v", err))
		errorHandler(err, w)
		return
	}
}

/*
400 FormatError
402 NotEnoughError
409 AlreadyExistError
422 IncorrectInputError
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

	var notEnoughError *models.NotEnoughError
	if errors.As(err, &notEnoughError) {
		utils.Log.Error("error not enough:", err.Error())
		http.Error(w, http.StatusText(http.StatusPaymentRequired), http.StatusPaymentRequired)
		return
	}

	var incorrectInputError *models.IncorrectInputError
	if errors.As(err, &incorrectInputError) {
		utils.Log.Error("error incorrect input:", err.Error())
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	utils.Log.Error("internal error: ", err.Error())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
