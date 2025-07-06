package config

import (
	"log"
	"os"
	"strconv"
)

func getEnvironmentValue(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s environment variable is missing", key)
	}
	return v
}

func getEnvironmentInt(key string) int {
	envStr := getEnvironmentValue(key)
	prasedInt, err := strconv.Atoi(envStr)
	if err != nil {
		log.Fatalf("%s: %s is invalid integer", key, envStr)
	}
	return prasedInt
}

func getEnvironmentIntWithDefault(key string, dflt int) int {
	v := os.Getenv(key)
	i := dflt
	if v != "" {
		parsedInt, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("%s: %s is invalid integer", key, v)
		}
		i = parsedInt
	}
	return i
}
