package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func New() http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.Name(route.Name).Methods(route.Method).Path(route.Path).Handler(Logger(route.Handle, route.Name))
	}

	handler := cors.Default().Handler(router)

	return handler
}
