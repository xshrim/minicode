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
		"Home",
		"GET",
		"/home",
		Home,
	},
	Route{
		"Exe",
		"POST",
		"/exe",
		Exe,
	},
	Route{
		"Info",
		"POST",
		"/info",
		Info,
	},
}

func NewRouter() *mux.Router {

	config, _ = LoadConfig("./config.json")

	clients = make(map[string]*Client)

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
