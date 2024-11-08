package server

import "net/http"

type Router struct {
	V1 *http.ServeMux
}

func NewRouter() *Router {
	var (
		router = http.NewServeMux()
	)

	router.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	return &Router{V1: router}
}
