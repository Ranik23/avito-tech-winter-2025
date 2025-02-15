package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/spf13/viper"
)

func LoadConfig() (*Config, error) {

	viper.SetConfigFile(".env") 
	viper.AutomaticEnv()  

	err := viper.ReadInConfig()
	if err != nil && os.IsNotExist(err) {
		log.Println(".env file not found, using environment variables.")
	} else if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	} else {
		log.Println(".env file loaded.")
	}


	cfg := &Config{
		HTTPServerConfig: HTTPServerConfig{
			Host: viper.GetString("HTTP_SERVER_HOST"),
			Port: viper.GetString("HTTP_SERVER_PORT"),
		},
		PostgresConfig: PostgresConfig{
			Host:     viper.GetString("POSTGRES_HOST"),
			Port:     viper.GetString("POSTGRES_PORT"),
			User:     viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			DBname:   viper.GetString("POSTGRES_DBNAME"),
		},
		LogConfig: LogConfig{
			Level: viper.GetString("LOG_LEVEL"),
		},
		JWTConfig: JWTConfig{
			Duration: viper.GetString("JWT_DURATION"),
			Secret:   viper.GetString("JWT_SECRET"),
		},
	}

	log.Println("Checking if variables came from environment or file:")
	if viper.IsSet("HTTP_SERVER_HOST") {
		if viper.InConfig("HTTP_SERVER_HOST") {
			log.Println("HTTP_SERVER_HOST was loaded from the .env file.")
		} else {
			log.Println("HTTP_SERVER_HOST was loaded from the environment.")
		}
	}

	return cfg, nil
}

func CreatePostgresDSN(cfg *Config) string {
	user := strings.TrimSpace(cfg.PostgresConfig.User)
	password := strings.TrimSpace(cfg.PostgresConfig.Password)
	host := strings.TrimSpace(cfg.PostgresConfig.Host)
	port := strings.TrimSpace(cfg.PostgresConfig.Port)
	dbname := strings.TrimSpace(cfg.PostgresConfig.DBname)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname,
	)
	log.Println("Postgres DSN:", dsn)
	return dsn
}


type Config struct {
	HTTPServerConfig HTTPServerConfig
	PostgresConfig   PostgresConfig
	LogConfig        LogConfig
	JWTConfig        JWTConfig
}

type HTTPServerConfig struct {
	Host string
	Port string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
}

type LogConfig struct {
	Level string
}

type JWTConfig struct {
	Secret   string
	Duration string
}
