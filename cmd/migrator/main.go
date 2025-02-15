package main

import (
	"avito/config"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}

	migrationsPath := flag.String("path", "./internal/migrations", "Путь к файлам миграции")
	flag.Parse()

	dsn := config.CreatePostgresDSN(cfg)
	fmt.Println("Используем DSN:", dsn)

	m, err := migrate.New(
		"file://"+*migrationsPath,
		dsn,
	)
	if err != nil {
		log.Fatalf("Ошибка при создании мигратора: %v", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Ошибка при применении миграций: %v", err)
	}

	fmt.Println("Миграции успешно применены или изменений нет")
}
