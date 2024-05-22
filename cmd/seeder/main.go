package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/omareloui/odinls/config"
	"github.com/omareloui/odinls/internal/application/core/role"
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

	roleService := role.NewRoleService(repo, validator)

	err = roleService.SeedRoles()
	if err != nil {
		log.Fatalln("err seeding the roles:", err)
		return
	}
}
