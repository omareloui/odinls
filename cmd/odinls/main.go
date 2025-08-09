package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/omareloui/formmap"
	"github.com/omareloui/odinls/config"
	"github.com/omareloui/odinls/internal/api/handler"
	"github.com/omareloui/odinls/internal/api/router"
	application "github.com/omareloui/odinls/internal/application/core"
	"github.com/omareloui/odinls/internal/logger"
	"github.com/omareloui/odinls/internal/repositories/mongo"
	"github.com/omareloui/odinls/internal/sanitizer/conformadaptor"
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

	validator := formmap.NewValidator()
	sanitizer := conformadaptor.NewSanitizer()

	app := application.NewApplication(repo, validator, sanitizer)

	h := handler.New(app)

	port := config.GetApplicationPort()

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router.New(h),
	}

	log.Printf("Starting to listen on http://localhost:%d\n", port)
	log.Fatalln(srv.ListenAndServe())
}
