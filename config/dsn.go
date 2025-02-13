package config

import "fmt"



func CreatePostgresDSN(cfg *Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.PostgresConfig.Password, cfg.PostgresConfig.Host, cfg.PostgresConfig.Port, cfg.DBname, "disable",
	)
}