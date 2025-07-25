package config

import "go.mongodb.org/mongo-driver/mongo/options"

func GetApplicationPort() int {
	return getEnvironmentInt("PORT")
}

func GetDataSource() string {
	return getEnvironmentValue("DATA_SOURCE")
}

func GetJwtSecret() []byte {
	return []byte(getEnvironmentValue("TOKEN_SECRET"))
}

func GetMongoDatabaseName() string {
	return getEnvironmentValue("MONGO_DB_NAME")
}

func GetMongoQueryTimeout() int {
	return getEnvironmentInt("MONGO_DB_QUERY_TIMEOUT")
}

func GetMongoCred() options.Credential {
	return options.Credential{
		Username: getEnvironmentValue("MONGO_DB_USER"),
		Password: getEnvironmentValue("MONGO_DB_PASSWORD"),
	}
}

func GetLogLevel() int {
	return getEnvironmentIntWithDefault("LOG_LEVEL", 0)
}
