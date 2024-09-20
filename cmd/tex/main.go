package main

import (
	"fmt"
	"texApi/config"
	"texApi/database"
	app "texApi/internal"
)

func main() {
	config.InitConfig()
	database.InitDB()
	server := app.InitApp()
	address := fmt.Sprintf("%v:%v", config.ENV.API_HOST, config.ENV.API_PORT)
	server.Run(address)
}
