package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
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

	_ = validator.RegisterValidation("not_blank", validators.NotBlank)
	_ = validator.RegisterValidation("alphanum_with_underscore", IsAlphaNumWithUnderScore)

	h := handler.New(app)

	port := config.GetApplicationPort()

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router.New(h),
	}

	log.Printf("Starting to listen on http://localhost:%d\n", port)
	log.Fatalln(srv.ListenAndServe())
}

func IsAlphaNumWithUnderScore(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[A-Za-z0-9_]+$`)
	field := fl.Field()
	return re.Match([]byte(field.String()))
}
