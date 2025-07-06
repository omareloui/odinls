package main

import (
	"github.com/joho/godotenv"
	"github.com/omareloui/odinls/config"
	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/api/chisrv"
	"github.com/omareloui/odinls/internal/api/resthandlers"
	application "github.com/omareloui/odinls/internal/application/core"
	"github.com/omareloui/odinls/internal/logger"
	"github.com/omareloui/odinls/internal/repositories/mongo"
	"github.com/omareloui/odinls/internal/sanitizer/conformadaptor"
	"github.com/omareloui/odinls/internal/validator/playgroundvalidator"
	"go.uber.org/zap"
)

func init() {
	if err := godotenv.Load(); err != nil {
		l := logger.Get()
		l.Fatal("Error loading .env file", zap.Error(err))
	}
}

func main() {
	repo, err := mongo.NewRepository(
		config.GetDataSource(),
		config.GetMongoDatabaseName(),
		config.GetMongoQueryTimeout(),
	)
	if err != nil {
		l := logger.Get()
		l.Fatal("Error creating repository", zap.Error(err))
	}

	validator := playgroundvalidator.NewValidator()
	sanitizer := conformadaptor.NewSanitizer()

	app := application.NewApplication(repo, validator, sanitizer)
	jwtAdapter := jwtadapter.NewJWTV5Adapter(config.GetJwtSecret())
	handler := resthandlers.NewHandler(app, jwtAdapter)

	api := chisrv.NewAdapter(handler, config.GetApplicationPort())
	l := logger.Get()
	l.Info("Starting OdinLS API", zap.Int("port", config.GetApplicationPort()))
	api.Run()
}
