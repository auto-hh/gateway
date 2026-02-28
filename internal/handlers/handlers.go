package handlers

import (
	"fmt"
	"net/http"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ProxyHandler is ready to work")
}
