package storage

import (
	"database/sql"
	"errors"
	"github.com/Nikolay961996/goferma/internal/models"
	"github.com/Nikolay961996/goferma/internal/utils"
)

func (db *DBContext) CreateNewUser(login string, pswHash string) error {
	query := `
		INSERT INTO users (login, passwordHash)
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
		INSERT INTO users (login, passwordHash)
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
		SELECT id, login, passwordHash
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
