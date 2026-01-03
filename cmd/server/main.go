package main

import (
	"myapi/internal/router"
)

func main() {
	r := router.Setup()
	r.Run(":8080")
}
