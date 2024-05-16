package config

import "go.mongodb.org/mongo-driver/mongo/options"

func GetApplicationPort() int {
	return getEnvironmentInt("PORT")
}

func GetDataSource() string {
	return getEnvironmentValue("DATA_SOURCE")
}

func GetMongoCred() options.Credential {
	return options.Credential{
		Username: getEnvironmentValue("MONGO_DB_USER"),
		Password: getEnvironmentValue("MONGO_DB_PASSWORD"),
	}
}
