package server

import (
	"context"
	"github.com/Nikolay961996/goferma/internal/server/router"
	"github.com/Nikolay961996/goferma/internal/storage"
	"github.com/Nikolay961996/goferma/internal/utils"
	"github.com/Nikolay961996/goferma/internal/workers"
	"net/http"
)

func Run(config *Config) {
	utils.Log.Info("Server run on ", config.runAddress)

	done, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbContext := storage.NewDBStorage(config.databaseUri)
	runWorkers(dbContext, done, config.accrualSystemAddress)

	err := http.ListenAndServe(config.runAddress, router.GofermaRouter(dbContext, config.secretKey))
	if err != nil {
		utils.Log.Fatal("Can't start server: ", err)
	}
}

func runWorkers(db *storage.DBContext, done context.Context, loyaltyAddress string) {
	out := workers.CreateWorkerDistributor(db, done)
	for i := 0; i < 5; i++ {
		go workers.RunWorker(i, db, done, out, loyaltyAddress)
	}
}
