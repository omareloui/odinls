package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/omareloui/odinls/config"
	"github.com/omareloui/odinls/internal/adapters/db"
	"github.com/omareloui/odinls/internal/adapters/rest"
	"github.com/omareloui/odinls/internal/application/core/api"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

	dbAdapter, _ := db.NewAdapter(config.GetDataSource(), config.GetDataSourceCred())

	application := api.NewApplication(dbAdapter)

	restAdapter := rest.NewAdapter(application, config.GetApplicationPort())
	restAdapter.Run()
}
