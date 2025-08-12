package workers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Nikolay961996/goferma/internal/models"
	"github.com/Nikolay961996/goferma/internal/storage"
	"github.com/Nikolay961996/goferma/internal/utils"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

func CreateWorkerDistributor(db *storage.DBContext, done context.Context) <-chan job {
	ticker := time.NewTicker(100 * time.Millisecond)
	out := make(chan job, 10)

	go func() {
		defer ticker.Stop()
		defer close(out)
		for {
			select {
			case <-ticker.C:
				orders, err := db.GerUnprocessedOrders()
				if err != nil {
					utils.Log.Error("error getting unprocessed orders: ", err.Error())
				}
				utils.Log.Info("Find ", len(orders), " unprocessed orders")
				for _, order := range orders {
					out <- job{Order: order}
				}
			case <-done.Done():
				utils.Log.Info("WorkerDistributor done")
				return
			}
		}
	}()

	return out
}

func RunWorker(workerID int, db *storage.DBContext, done context.Context, jobs <-chan job, loyaltyAddress string) {
	utils.Log.Info("Starting worker", workerID)

	for {
		select {
		case job := <-jobs:
			loyalty := sendToLoyalty(loyaltyAddress, job.Order.Number)
			loyaltyStatus := loyaltyStatus(loyalty.Status)
			newStatus := loyaltyStatusToOrderStatus(loyaltyStatus)
			if job.Order.CurrentStatus == newStatus {
				utils.Log.Info("Order", job.Order.Number, "status is not changed. New status=", loyaltyStatus)
				continue
			}
			updateOrder(db, job.Order.ID, loyalty.Accrual, newStatus)
			utils.Log.Warn("Worker ", workerID, " is update ", job.Order.Number, " status ", newStatus, ", sum= ", loyalty.Accrual)

		case <-done.Done():
			utils.Log.Info("WorkerDistributor done")
			return
		}
	}
}

func sendToLoyalty(loyaltyAddress string, orderNumber string) *loyaltyResponse {
	client := resty.New()
	url := fmt.Sprintf("%s/api/orders/%s", loyaltyAddress, orderNumber)

	var body []byte
	err := utils.RetryerCon(func() error {
		b, err := sendRequest(client, url)
		if err == nil {
			body = b
		}
		return err
	}, func(err error) bool {
		var tooManyRequestsError *models.TooManyRequestsError
		return errors.As(err, &tooManyRequestsError)
	})
	if err != nil {
		utils.Log.Error(err.Error())
		return nil
	}

	var data loyaltyResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		utils.Log.Error(err.Error())
		return nil
	}

	return &data
}

func sendRequest(client *resty.Client, url string) ([]byte, error) {
	request := client.R()
	response, err := request.Get(url)
	if err != nil {
		utils.Log.Error(err.Error())
		return nil, err
	}
	if response.StatusCode() == http.StatusTooManyRequests {
		return nil, &models.TooManyRequestsError{Err: errors.New("to many requests")}
	}
	if response.StatusCode() != http.StatusOK {
		utils.Log.Error(fmt.Errorf("unexpected status code: %d", response.StatusCode()))
		return nil, nil
	}

	return response.Body(), nil
}

func updateOrder(db *storage.DBContext, orderID int64, accrual float64, newStatus models.OrderStatus) {
	err := db.UpdateOrder(orderID, newStatus, accrual)
	if err != nil {
		utils.Log.Error(err.Error())
	}
}

func loyaltyStatusToOrderStatus(loyaltyStatus loyaltyStatus) models.OrderStatus {
	switch loyaltyStatus {
	case registered:
		return models.New
	case processing:
		return models.Processing
	case processed:
		return models.Processed
	case invalid:
		return models.Invalid
	}
	utils.Log.Error("Unknown loyalty status. ", loyaltyStatus, " Set as invalid")
	return models.Invalid
}
