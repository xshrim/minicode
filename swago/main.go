// Package main Testing go-swagger generation
//
// The purpose of this application is to test go-swagger in a simple GET requests.
//
// Schemes: http
// Host: 192.168.43.103:9090
// Version: 0.0.1
// License: MIT http://opensource.org/licenses/MIT
// Contract: xshrim<xshrim@gmail.com>
// Terms Of Service:
// there are no TOS at this moment.
//
// Consumes:
// - text/plain
// - application/json
//
// Produces:
// - text/plain
// - application/json
//
// swagger:meta
package main

import (
	"log"
	"net/http"

	"./rest"
)

func main() {

	handler := rest.New()
	log.Fatal(http.ListenAndServe(":9090", handler))

}
