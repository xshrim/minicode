package rest

var routes = []Route{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Hello",
		"GET",
		"/hello/{name}",
		Hello,
	},
	Route{
		"View",
		"GET",
		"/view/{id}",
		View,
	},
	Route{
		"List",
		"GET",
		"/list",
		List,
	},
	Route{
		"Swagger",
		"GET",
		"/swagger.json",
		Swagger,
	},
}
