package main

import (
	"log"

	"imgop/internal/db"
	"imgop/internal/routers"
)

func main() {
	db.ConnectMongoDB()
	defer db.DisconnectMongoDB()

	app := routers.Setup()
	app.Run(":8080")
	log.Println("Sever is runnig in 8080")
}
