package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/omareloui/odinls/config"
	"github.com/omareloui/odinls/internal/api/restfiber"
	application "github.com/omareloui/odinls/internal/application/core"
	"github.com/omareloui/odinls/internal/repositories/mongo"
	"github.com/omareloui/odinls/internal/validator/playgroundvalidator"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

	repo, err := mongo.NewRepository(
		config.GetDataSource(),
		"ODINLS_DEV", // TODO: get from env var
		14,
	)
	if err != nil {
		log.Fatal(err)
	}

	validator := playgroundvalidator.NewValidator()
	app := application.NewApplication(repo, validator)
	handler := restfiber.NewHandler(app)

	api := restfiber.NewAdapter(handler, config.GetApplicationPort())
	api.Run()
}
