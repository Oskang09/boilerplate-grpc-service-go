package config

import (
	"fmt"
	"os"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
)

var (
	MalaysiaTimezone    = new(time.Location)
	ValidatorTranslator ut.Translator

	Env = getEnv("ENV")

	RedisHostPath = getEnv("REDIS_HOST")
	RedisPassword = getEnv("REDIS_PASSWORD")

	DBHost     = getEnv("MYSQL_HOST")
	DBPort     = getEnv("MYSQL_PORT")
	DBUser     = getEnv("MYSQL_USER")
	DBPassword = getEnv("MYSQL_PASSWORD")
	DBName     = getEnv("MYSQL_DATABASE")
)

func init() {
	MalaysiaTimezone, _ = time.LoadLocation("Asia/Kuala_Lumpur")

	english := en.New()
	uni := ut.New(english, english)
	ValidatorTranslator, _ = uni.GetTranslator("en")
}

func getEnv(i string) string {
	l, isExist := os.LookupEnv(i)
	if !isExist {
		panic(fmt.Sprintf("%s cannot find the environment", i))
	}

	return l
}

// IsProduction : Return true if the environment is production
var IsProduction = func() bool {
	return Env == "production"
}

var IsLocal = func() bool {
	return Env == "local"
}
