package rest

import "net/http"

// swagger:ignore
type Route struct {
	Name   string
	Method string
	Path   string
	Handle http.HandlerFunc
}

// swagger:response
type BookResponse struct {
	Resp Book `json:"resp"`
}

// A MyError is an error.
// swagger:response
type MyError struct {
	// The error message
	// in: body
	Message string `json:"message"`
}

// A book entity
//
// swagger:model Book
type Book struct {
	// the id for this book
	//
	// required: true
	// min length: 1
	Id string `json:"id"`

	// the authors of this book
	Authors []string `json:"authors"`

	// the title of this book
	//
	// min length: 3
	// example: The Golang Programming
	Title string `json:"title"`

	// the price of this book
	//
	// required: true
	// min: 1
	Price uint `json:"price"`
}
