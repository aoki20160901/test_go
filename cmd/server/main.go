package main

import (
	"net/http"

	"myapi/internal/router"
)

// main関数
func main() {
	r := router.Setup()

	http.ListenAndServe(":8080", r)
}
