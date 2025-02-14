package main

import (
	"avito/internal/app"
)




func main() {
	
	app, err := app.NewApp("config.yaml")
	if err != nil {
		panic(err)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}