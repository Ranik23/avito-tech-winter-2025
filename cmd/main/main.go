package main

import (
	"avito/internal/app"
	"log"
)




func main() {
	
	configPath := "config.yaml"
	app, err := app.NewApp(configPath)
	if err != nil {
		panic(err)
	}

	
	if err := app.Run(); err != nil {
		log.Fatalf(err.Error())
	}
}