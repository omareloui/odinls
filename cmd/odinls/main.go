package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/omareloui/odinls/config"
	"github.com/omareloui/odinls/internal/api/fibersrv"
	"github.com/omareloui/odinls/internal/api/resthandlers"
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
		config.GetMongoDatabaseName(),
		config.GetMongoQueryTimeout(),
	)
	if err != nil {
		log.Fatal(err)
	}

	validator := playgroundvalidator.NewValidator()
	app := application.NewApplication(repo, validator)
	handler := resthandlers.NewHandler(app)

	api := fibersrv.NewAdapter(handler, config.GetApplicationPort())
	api.Run()
}
