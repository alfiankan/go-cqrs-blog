package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// ApplicationConfig contain all application config loaded from env or (dot)env file
type ApplicationConfig struct {
	PostgreeHost          string
	PostgreeUser          string
	PostgreePass          string
	PostgreeDb            string
	PostgreePort          int
	PostgreeSsl           string
	ElasticSearchAdresses []string
	ElasticSearchUsername string
	ElasticSearchPassword string
	RedisHost             string
	RedisPass             string
	HTTPApiPort           string
}

// Load load config from env
func Load(configFile ...string) ApplicationConfig {

	if err := godotenv.Load(configFile...); err != nil {
		panic(err)
	}

	postgreeDbPort, err := strconv.Atoi(os.Getenv("PG_DATABASE_PORT"))
	if err != nil {
		panic(err)
	}

	return ApplicationConfig{
		PostgreeHost:          os.Getenv("PG_DATABASE_HOST"),
		PostgreeUser:          os.Getenv("PG_DATABASE_USERNAME"),
		PostgreePass:          os.Getenv("PG_DATABASE_PASSWORD"),
		PostgreeDb:            os.Getenv("PG_DATABASE_NAME"),
		PostgreePort:          postgreeDbPort,
		PostgreeSsl:           os.Getenv("PG_DATABASE_SSL_MODE"),
		ElasticSearchAdresses: strings.Split(os.Getenv("ELASTICSEARCH_ADDRESSES"), ";"),
		ElasticSearchUsername: os.Getenv("ELASTICSEARCH_USERNAME"),
		ElasticSearchPassword: os.Getenv("ELASTICSEARCH_PASSWORD"),
		RedisHost:             os.Getenv("REDIS_HOST"),
		RedisPass:             os.Getenv("REDIS_PASSWORD"),
		HTTPApiPort:           os.Getenv("HTTP_API_PORT"),
	}

}
