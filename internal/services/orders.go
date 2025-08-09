package services

import (
	"errors"
	"fmt"
	"github.com/Nikolay961996/goferma/internal/models"
	"github.com/Nikolay961996/goferma/internal/storage"
	"github.com/Nikolay961996/goferma/internal/utils"
	"io"
	"strconv"
	"strings"
)

func GetOrderNumber(contentType string, body io.ReadCloser) (string, error) {
	if contentType != "text/plain" {
		utils.Log.Error(errors.New("not text/plain"))
		return "", &models.FormatError{Err: errors.New("not text/plain")}
	}

	number, err := utils.ReadPlainTextBody(body)
	if err != nil {
		utils.Log.Error(errors.New("error reading text"))
		return "", &models.FormatError{Err: errors.New("error reading text")}
	}
	number = strings.ReplaceAll(number, " ", "")
	if !isCorrectOrderNumber(number) {
		utils.Log.Error(errors.New(fmt.Sprintf("order number is incorrect. '%s'", number)))
		return "", &models.IncorrectInputError{Err: errors.New("order number is incorrect")}
	}

	return number, nil
}

func RegisterOrder(db *storage.DBContext, orderNumber string, userId int64) (bool, error) {
	orderUserId, err := db.GetUserForOrder(orderNumber)
	if err != nil {
		return false, err
	}

	if orderUserId == 0 {
		err = db.SetUserOrder(userId, orderNumber)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	if orderUserId != userId {
		return false, &models.AlreadyExistError{Err: errors.New("already used order number by other user")}
	}

	return false, nil
}

func isCorrectOrderNumber(number string) bool {
	utils.Log.Info("checking order number ", number)
	if number == "" {
		return false
	}

	sum := 0
	parity := len(number) % 2
	for i, c := range number {
		digit, err := strconv.Atoi(string(c))
		if err != nil {
			return false
		}
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}
