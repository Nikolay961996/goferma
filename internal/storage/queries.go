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

func (db *DBContext) SetUserOrder(userId int64, orderNumber string) error {
	query := `
		INSERT INTO orders (user_id, order_number, accrual, status, uploaded_at)
		VALUES ($1, $2, $3, $4);`
	_, err := db.db.Exec(query, userId, orderNumber, 0, models.NEW, time.Now())
	if err != nil {
		utils.Log.Error("error insert new order for user: ", err.Error())
		return err
	}

	return nil
}
