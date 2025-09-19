package main

import (
	"fmt"

	"imgop/db"
	"imgop/routers"
)

func main() {
	db.ConnectMongoDB()
	defer db.DisconnectMongoDB()

	router := routers.Setup()
	router.Run(":8080")
	fmt.Println("Sever is runnig in 8080")
}
