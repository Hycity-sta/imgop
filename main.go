package main

import (
	"imgop/routers"
)

func main() {
	router := routers.Setup()
	router.Run(":8080")
}
