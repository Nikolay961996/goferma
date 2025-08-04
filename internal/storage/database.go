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

	sqlInsertOrUpdateGauge   *sql.Stmt
	sqlInsertOrUpdateCounter *sql.Stmt
	sqlGetGauge              *sql.Stmt
	sqlGetCounter            *sql.Stmt
	sqlGetAll                *sql.Stmt
}

func NewDBStorage(databaseDSN string) *DBContext {
	s := DBContext{}
	s.open(databaseDSN)
	s.migrate()

	return &s
}

func (m *DBContext) open(databaseDSN string) {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		utils.Log.Fatal("Db connection error: ", err)
	}
	m.db = db
	m.databaseDSN = databaseDSN
}

func (m *DBContext) migrate() {
	migrateFunc := func() error {
		driver, err := postgres.WithInstance(m.db, &postgres.Config{})
		if err != nil {
			utils.Log.Fatal(fmt.Sprintf("migration driver creation error: %s", err.Error()))
			return err
		}

		instance, err := migrate.NewWithDatabaseInstance("file://internal/storage/migrations", m.databaseDSN, driver)
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
