package storage

import (
	"database/sql"
	"errors"
	"github.com/Nikolay961996/goferma/internal/models"
	"github.com/Nikolay961996/goferma/internal/utils"
	"github.com/lib/pq"
	"time"
)

func (db *DBContext) CreateNewUser(login string, pswHash string) error {
	query := `
		INSERT INTO users (login, password_hash)
		VALUES ($1, $2);`
	_, err := db.db.Exec(query, login, pswHash)
	if err != nil {
		utils.Log.Error("create user error ", err.Error())
		return err
	}

	return nil
}

func (db *DBContext) GetUser(login string) (*models.User, error) {
	query := `
		SELECT id, login, password_hash
		FROM users	
		WHERE login = $1;`
	var user models.User
	err := db.db.QueryRow(query, login).Scan(&user.ID, &user.Login, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		utils.Log.Error("error get user: ", err.Error())
		return nil, err
	}

	return &user, nil
}

func (db *DBContext) GetUserForOrder(orderNumber string) (int64, error) {
	query := `
		SELECT user_id
		FROM orders	
		WHERE order_number = $1;`
	var userID int64
	err := db.db.QueryRow(query, orderNumber).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		utils.Log.Error("error get user order: ", err.Error())
		return 0, err
	}
	return userID, nil
}

func (db *DBContext) SetUserOrder(userID int64, orderNumber string, status models.OrderStatus, accrual float64) error {
	query := `
		INSERT INTO orders (user_id, order_number, accrual, status, uploaded_at)
		VALUES ($1, $2, $3, $4, $5);`
	_, err := db.db.Exec(query, userID, orderNumber, int64(accrual*100), status, time.Now())
	utils.Log.Warn("SetUserOrder: ", orderNumber, " = ", int64(accrual*100), ", status ", status)
	if err != nil {
		utils.Log.Error("error insert new order for user: ", err.Error())
		return err
	}

	return nil
}

func (db *DBContext) GetUserOrders(userID int64) ([]models.OrdersResponse, error) {
	query := `
		SELECT order_number, status, accrual, uploaded_at
		FROM orders
		WHERE user_id = $1;`

	var orders []models.OrdersResponse
	rows, err := db.db.Query(query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		utils.Log.Error("error get rows userOrders: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.OrdersResponse
		var accrual int64
		var status models.OrderStatus
		var uploadedAt time.Time
		err := rows.Scan(&order.Number, &status, &accrual, &uploadedAt)
		if err != nil {
			utils.Log.Error("error get userOrders: ", err.Error())
			return nil, err
		}
		order.Accrual = float64(accrual) / 100
		order.Status = status.String()
		order.UploadedAt = uploadedAt.Format(time.RFC3339)
		orders = append(orders, order)

	}

	err = rows.Err()
	if err != nil {
		utils.Log.Error(err)
		return nil, err
	}

	return orders, nil
}

func (db *DBContext) GerUserCurrentAccrual(userID int64) (float64, error) {
	query := `
		SELECT sum(accrual)
		FROM orders
		WHERE user_id = $1 and status = $2;`
	var accrualSum sql.NullInt64
	err := db.db.QueryRow(query, userID, models.Processed).Scan(&accrualSum)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		utils.Log.Error("error get user balance: ", err.Error())
		return 0, err
	}
	if !accrualSum.Valid {
		return 0, nil
	}

	return float64(accrualSum.Int64) / 100, nil
}

func (db *DBContext) GerUserWithdrawn(userID int64) (float64, error) {
	query := `
		SELECT sum(accrual)
		FROM orders
		WHERE user_id = $1 and status = $2 and accrual < 0;`
	var withdrawnSum sql.NullInt64
	err := db.db.QueryRow(query, userID, models.Processed).Scan(&withdrawnSum)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		utils.Log.Error("error get user balance: ", err.Error())
		return 0, err
	}
	if !withdrawnSum.Valid {
		return 0, nil
	}

	return float64(-withdrawnSum.Int64) / 100, nil
}

func (db *DBContext) GerUserWithdrawnHistory(userID int64) ([]models.WithdrawHistoryResponse, error) {
	query := `
		SELECT order_number, accrual, uploaded_at
		FROM orders
		WHERE user_id = $1 and status = $2 and accrual < 0;`

	var withdrawHistory []models.WithdrawHistoryResponse
	rows, err := db.db.Query(query, userID, models.Processed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		utils.Log.Error("error get rows withdrawn history: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var withdraw models.WithdrawHistoryResponse
		var sum int64
		var uploadedAt time.Time
		err := rows.Scan(&withdraw.Order, &sum, &uploadedAt)
		if err != nil {
			utils.Log.Error("error get withdrawn history: ", err.Error())
			return nil, err
		}
		withdraw.Sum = float64(-sum) / 100
		withdraw.ProcessedAt = uploadedAt.Format(time.RFC3339)
		withdrawHistory = append(withdrawHistory, withdraw)
	}

	err = rows.Err()
	if err != nil {
		utils.Log.Error(err)
		return nil, err
	}

	return withdrawHistory, nil
}

func (db *DBContext) GerUnprocessedOrders() ([]models.Order, error) {
	query := `
		SELECT id, order_number, status
		FROM orders
		WHERE NOT (status = ANY($1));`

	var orders []models.Order
	rows, err := db.db.Query(query, pq.Array([]models.OrderStatus{models.Processed, models.Invalid}))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		utils.Log.Error("error get rows unprocessed orders: ", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.Number, &order.CurrentStatus)
		if err != nil {
			utils.Log.Error("error get unprocessed orders: ", err.Error())
			return nil, err
		}
		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		utils.Log.Error(err)
		return nil, err
	}

	return orders, nil
}

func (db *DBContext) UpdateOrder(orderId int64, newStatus models.OrderStatus, accrual float64) error {
	query := `
		UPDATE orders SET
			status = $1,
			accrual = $2
		WHERE id = $3;`
	_, err := db.db.Exec(query, newStatus, int64(accrual*100), orderId)
	if err != nil {
		utils.Log.Error("error update order: ", err.Error())
		return err
	}

	return nil
}
