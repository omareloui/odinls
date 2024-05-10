package config

import "go.mongodb.org/mongo-driver/mongo/options"

func GetApplicationPort() int {
	return getEnvironmentInt("PORT")
}

func GetDataSource() string {
	return getEnvironmentValue("DATA_SOURCE")
}

func GetDataCred() options.Credential {
	return options.Credential{
		Username: getEnvironmentValue("DB_USER"),
		Password: getEnvironmentValue("DB_PASSWORD"),
	}
}
