package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Nikolay961996/goferma/internal/utils"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBContext struct {
	databaseDSN string
	db          *sql.DB
	tx          *sql.Tx

	sqlInsertNewUser *sql.Stmt
}

func NewDBStorage(databaseDSN string) *DBContext {
	s := DBContext{}
	s.open(databaseDSN)
	s.migrate()
	s.prepareSQL()

	return &s
}

func (db *DBContext) open(databaseDSN string) {
	dbContext, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		utils.Log.Fatal("Db connection error: ", err)
	}
	db.db = dbContext
	db.databaseDSN = databaseDSN
}

func (db *DBContext) migrate() {
	migrateFunc := func() error {
		driver, err := postgres.WithInstance(db.db, &postgres.Config{})
		if err != nil {
			utils.Log.Fatal(fmt.Sprintf("migration driver creation error: %s", err.Error()))
			return err
		}

		instance, err := migrate.NewWithDatabaseInstance("file://internal/storage/migrations", db.databaseDSN, driver)
		if err != nil {
			utils.Log.Fatal(fmt.Sprintf("migration instance creation error: %s", err.Error()))
			return err
		}

		if err := instance.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			utils.Log.Fatal(fmt.Sprintf("migration instance up error: %s", err.Error()))
			return err
		}
		return nil
	}

	err := migrateFunc()
	if err != nil {
		utils.Log.Fatal(err.Error())
	}
}

func (db *DBContext) prepareSQL() {
	sqlInsertNewUser, err := db.db.Prepare(
		`
		INSERT INTO users (login, password_hash)
		VALUES ($1, $2);`)
	if err != nil {
		utils.Log.Fatal(err.Error())
	}
	db.sqlInsertNewUser = sqlInsertNewUser
}
