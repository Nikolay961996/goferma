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

	sqlInsertNewUser             *sql.Stmt
	sqlSelectUserByLogin         *sql.Stmt
	sqlSelectUserIDByOrderNumber *sql.Stmt
	sqlInsertNewOrder            *sql.Stmt
	sqlSelectOrderByUserID       *sql.Stmt
	sqlSelectAccrualSum          *sql.Stmt
	sqlSelectWithdrawnSum        *sql.Stmt
	sqlSelectWithdrawnOrders     *sql.Stmt
	sqlSelectOrdersByNotInStatus *sql.Stmt
	sqlUpdateOrder               *sql.Stmt
}

func NewDBStorage(databaseDSN string) *DBContext {
	s := DBContext{}
	s.open(databaseDSN)
	s.migrate()

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

	sqlSelectUserByLogin, err := db.db.Prepare(
		`
		SELECT id, login, password_hash
		FROM users	
		WHERE login = $1;`)
	if err != nil {
		utils.Log.Fatal(err.Error())
	}

	sqlSelectUserIDByOrderNumber, err := db.db.Prepare(
		`
		SELECT user_id
		FROM orders	
		WHERE order_number = $1;`)
	if err != nil {
		utils.Log.Fatal(err.Error())
	}

	sqlInsertNewOrder, err := db.db.Prepare(
		`
		INSERT INTO orders (user_id, order_number, accrual, status, uploaded_at)
		VALUES ($1, $2, $3, $4, $5);`)
	if err != nil {
		utils.Log.Fatal(err.Error())
	}

	sqlSelectOrderByUserID, err := db.db.Prepare(
		`
		SELECT order_number, status, accrual, uploaded_at
		FROM orders
		WHERE user_id = $1;`)
	if err != nil {
		utils.Log.Fatal(err.Error())
	}

	sqlSelectAccrualSum, err := db.db.Prepare(
		`
		SELECT sum(accrual)
		FROM orders
		WHERE user_id = $1 and status = $2;`)
	if err != nil {
		utils.Log.Fatal(err.Error())
	}

	sqlSelectWithdrawnSum, err := db.db.Prepare(
		`
		SELECT sum(accrual)
		FROM orders
		WHERE user_id = $1 and status = $2 and accrual < 0;`)
	if err != nil {
		utils.Log.Fatal(err.Error())
	}

	sqlSelectWithdrawnOrders, err := db.db.Prepare(
		`
		SELECT order_number, accrual, uploaded_at
		FROM orders
		WHERE user_id = $1 and status = $2 and accrual < 0;`)
	if err != nil {
		utils.Log.Fatal(err.Error())
	}

	sqlSelectOrdersByNotInStatus, err := db.db.Prepare(
		`
		SELECT id, order_number, status
		FROM orders
		WHERE NOT (status = ANY($1))`)
	if err != nil {
		utils.Log.Fatal(err.Error())
	}

	sqlUpdateOrder, err := db.db.Prepare(
		`
		UPDATE orders SET
			status = $1,
			accrual = $2
		WHERE id = $3;`)
	if err != nil {
		utils.Log.Fatal(err.Error())
	}

	db.sqlInsertNewUser = sqlInsertNewUser
	db.sqlSelectUserByLogin = sqlSelectUserByLogin
	db.sqlSelectUserIDByOrderNumber = sqlSelectUserIDByOrderNumber
	db.sqlInsertNewOrder = sqlInsertNewOrder
	db.sqlSelectOrderByUserID = sqlSelectOrderByUserID
	db.sqlSelectAccrualSum = sqlSelectAccrualSum
	db.sqlSelectWithdrawnSum = sqlSelectWithdrawnSum
	db.sqlSelectWithdrawnOrders = sqlSelectWithdrawnOrders
	db.sqlSelectOrdersByNotInStatus = sqlSelectOrdersByNotInStatus
	db.sqlUpdateOrder = sqlUpdateOrder
}
