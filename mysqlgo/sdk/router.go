package sdk

import (
	"net/http"

	"github.com/gorilla/mux"
)

var routes = []Route{
	Route{
		"Index",
		"GET",
		"/index",
		Index,
	},
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Help",
		"GET",
		"/help",
		Help,
	},
	Route{
		"Query",
		"POST",
		"/query",
		Query,
	},
}

func NewRouter() *mux.Router {

	config, _ = LoadConfig("./config.json")

	//fmt.Println(config)

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
