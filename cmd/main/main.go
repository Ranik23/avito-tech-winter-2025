package main

import (
	"avito/internal/app"
)




func main() {
	
	configPath := "config.yaml"
	app, err := app.NewApp(configPath)
	if err != nil {
		panic(err)
	}


	if err := app.Run(); err != nil {
		panic(err)
	}
}