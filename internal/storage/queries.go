package storage

import (
	"database/sql"
	"errors"
	"github.com/Nikolay961996/goferma/internal/models"
	"github.com/Nikolay961996/goferma/internal/utils"
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

func (db *DBContext) AddNewUser(login string, pswHash string) error {
	query := `
		INSERT INTO users (login, password_hash)
		VALUES ($1, $2);`
	_, err := db.db.Exec(query, login, pswHash)
	if err != nil {
		utils.Log.Error("error insert new user: ", err.Error())
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
	var userId int64
	err := db.db.QueryRow(query, orderNumber).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		utils.Log.Error("error get user order: ", err.Error())
		return 0, err
	}
	return userId, nil
}

func (db *DBContext) SetUserOrder(userId int64, orderNumber string, status models.OrderStatus, accrual float64) error {
	query := `
		INSERT INTO orders (user_id, order_number, accrual, status, uploaded_at)
		VALUES ($1, $2, $3, $4, $5);`
	_, err := db.db.Exec(query, userId, orderNumber, int64(accrual*100), status, time.Now())
	if err != nil {
		utils.Log.Error("error insert new order for user: ", err.Error())
		return err
	}

	return nil
}

func (db *DBContext) GetUserOrders(userId int64) ([]models.OrdersResponse, error) {
	query := `
		SELECT order_number, status, accrual, uploaded_at
		FROM orders
		WHERE user_id = $1;`

	var orders []models.OrdersResponse
	rows, err := db.db.Query(query, userId)
	defer rows.Close()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		utils.Log.Error("error get rows userOrders: ", err.Error())
		return nil, err
	}

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

func (db *DBContext) GerUserCurrentAccrual(userId int64) (float64, error) {
	query := `
		SELECT sum(accrual)
		FROM orders
		WHERE user_id = $1 and status = $2;`
	var accrualSum sql.NullInt64
	err := db.db.QueryRow(query, userId, models.PROCESSED).Scan(&accrualSum)
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

func (db *DBContext) GerUserWithdrawn(userId int64) (float64, error) {
	query := `
		SELECT sum(accrual)
		FROM orders
		WHERE user_id = $1 and status = $2 and accrual < 0;`
	var withdrawnSum sql.NullInt64
	err := db.db.QueryRow(query, userId, models.PROCESSED).Scan(&withdrawnSum)
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

func (db *DBContext) GerUserWithdrawnHistory(userId int64) ([]models.WithdrawHistoryResponse, error) {
	query := `
		SELECT order_number, accrual, uploaded_at
		FROM orders
		WHERE user_id = $1 and status = $2 and accrual < 0;`

	var withdrawHistory []models.WithdrawHistoryResponse
	rows, err := db.db.Query(query, userId, models.PROCESSED)
	defer rows.Close()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		utils.Log.Error("error get rows withdrawn history: ", err.Error())
		return nil, err
	}

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
