package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Nikolay961996/goferma/internal/models"
	"github.com/Nikolay961996/goferma/internal/storage"
	"github.com/Nikolay961996/goferma/internal/utils"
	"io"
	"strings"
)

func GetUserBalance(db *storage.DBContext, userID int64) (*models.BalanceResponse, error) {
	accrual, err := db.GerUserCurrentAccrual(userID)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("Error get user balance for user id %d", userID))
		return nil, err
	}
	withdrawn, err := db.GerUserWithdrawn(userID)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("Error get user withdrawn for user id %d", userID))
		return nil, err
	}

	return &models.BalanceResponse{
		Accrual:   accrual,
		Withdrawn: withdrawn,
	}, nil
}

func ReadWithdrawnModel(contentType string, body io.ReadCloser) (*models.WithdrawnRequest, error) {
	if contentType != "application/json" {
		utils.Log.Error(errors.New("not json"))
		return nil, &models.FormatError{Err: errors.New("not text/plain")}
	}

	bytes, err := utils.ReadJSONBody(body)
	if err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}

	var model models.WithdrawnRequest
	err = json.Unmarshal(bytes, &model)
	if err != nil {
		utils.Log.Error(err.Error())
		return nil, &models.FormatError{Err: err}
	}

	return &model, nil
}

func Withdrawn(db *storage.DBContext, userID int64, withdrawn float64, orderNumber string) error {
	orderNumber = strings.ReplaceAll(orderNumber, " ", "")
	if !isCorrectOrderNumber(orderNumber) {
		utils.Log.Error(fmt.Errorf("order number is incorrect. '%s'", orderNumber))
		return &models.IncorrectInputError{Err: errors.New("order number is incorrect")}
	}
	accrual, err := db.GerUserCurrentAccrual(userID)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("Error get user balance for user id %d", userID))
		return err
	}

	if accrual < withdrawn {
		utils.Log.Error(fmt.Sprintf("Not enoth accrual. Current %f, userID %d", accrual, userID))
		return &models.NotEnoughError{Err: fmt.Errorf("not enoth accrual. Current %f, userID %d", accrual, userID)}
	}
	createdNew, err := RegisterOrder(db, orderNumber, userID, models.Processed, -withdrawn)
	if err != nil {
		utils.Log.Error(fmt.Sprintf("Error register order for user id %d", userID))
		return err
	}
	if !createdNew {
		utils.Log.Error(fmt.Sprintf("already used order for withdrawn. Order %s, user %d", orderNumber, userID))
		return errors.New("already used order for withdrawn")
	}

	return nil
}
