package config

import (
	"io"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)





func LoadConfig(path string) (*Config, error) {

	file, err := os.Open(path)
	if err != nil {
		log.Println("error opening the config file:", err)
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		log.Println("failed to read the file:", err)
		return nil, err
	}


	var cfg Config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Println("error unmarshalling the data:", err)
		return nil, err
	}

	return &cfg, nil
}


type Config struct {
	HTTPServerConfig	    `yaml:"httpServerConfig"`
	PostgresConfig 			`yaml:"postgresConfig"`
	RedisConfig	   			`yaml:"redisConfig"`
}


type HTTPServerConfig struct {
	Host 		string		`yaml:"host"`
	Port		string		`yaml:"port"`
}

type PostgresConfig struct {
	Host 		string 		`yaml:"host"`
	Port 		string		`yaml:"port"`
	User 		string		`yaml:"user"`
	Password 	string		`yaml:"password"`
	DBname 		string		`yaml:"dbname"`
}

type RedisConfig struct {
	Host		string		`yaml:"host"`
	Port		string		`yaml:"port"`
	Password	string		`yaml:"password"`
}

type JWTConfig struct {
	Duration 			time.Time 	`yaml:"duration"`
}