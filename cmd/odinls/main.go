package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/omareloui/odinls/config"
	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/api/fibersrv"
	"github.com/omareloui/odinls/internal/api/resthandlers"
	application "github.com/omareloui/odinls/internal/application/core"
	"github.com/omareloui/odinls/internal/repositories/mongo"
	"github.com/omareloui/odinls/internal/validator/playgroundvalidator"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}
}

func main() {
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
	jwtAdapter := jwtadapter.NewJWTV5Adapter(config.GetJwtSecret())
	handler := resthandlers.NewHandler(app, jwtAdapter)

	api := fibersrv.NewAdapter(handler, config.GetApplicationPort())
	api.Run()
}
