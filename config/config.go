package config

func GetApplicationPort() int {
	return getEnvironmentInt("PORT")
}
