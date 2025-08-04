package storage

import "github.com/Nikolay961996/goferma/internal/utils"

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
