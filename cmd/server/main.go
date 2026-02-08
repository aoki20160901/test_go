package main

import (
	"net/http"

	"myapi/internal/router"
)

func main() {
	r := router.Setup()

	http.ListenAndServe(":8080", r)
}
