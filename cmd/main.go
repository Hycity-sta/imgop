package main

import (
	"log"

	"imgop/internal/db"
	"imgop/internal/routers"
)

func main() {
	db.ConnectMongoDB()
	defer db.DisconnectMongoDB()

	router := routers.Setup()
	router.Run(":8080")
	log.Println("Sever is runnig in 8080")
}
