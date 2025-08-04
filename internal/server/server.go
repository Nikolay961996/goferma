package server

import (
	"github.com/Nikolay961996/goferma/internal/server/router"
	"github.com/Nikolay961996/goferma/internal/storage"
	"github.com/Nikolay961996/goferma/internal/utils"
	"net/http"
)

func Run(config *Config) {
	utils.Log.Info("Server run on ", config.runAddress)

	dbContext := storage.NewDBStorage(config.databaseUri)

	err := http.ListenAndServe(config.runAddress, router.GofermaRouter(dbContext))
	if err != nil {
		utils.Log.Fatal("Can't start server: ", err)
	}
}
