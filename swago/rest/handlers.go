package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET / index Index
	//
	// 欢迎信息
	// ---
	// produces:
	// - text/plain
	// responses:
	//   200:
	//     description: 欢迎信息
	//     type: string
	fmt.Fprintln(w, "Welcome")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /hello/{name} hello Hello
	//
	// Returns a simple Hello message
	// ---
	// consumes:
	// - text/plain
	// produces:
	// - text/plain
	// parameters:
	// - name: name
	//   in: path
	//   description: Name to be returned
	//   required: true
	//   type: string
	// responses:
	//   200:
	//     description: The hello message
	//     type: string
	log.Println("Responsing to /hello request")
	log.Println(r.UserAgent())

	vars := mux.Vars(r)
	name := vars["name"]

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello:", name)
}

func List(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /list list List
	//
	// Return a list of books
	// ---
	// consumes:
	// - text/plain
	// produces:
	// - application/json
	// parameters:
	// - name: min
	//   in: query
	//   type: string
	//   required: true
	// - name: max
	//   in: query
	//   type: string
	//   required: true
	// responses:
	//   200:
	//     description: book list
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/Book"
	r.ParseForm()

	minp, _ := strconv.Atoi(r.Form["min"][0])
	maxp, _ := strconv.Atoi(r.Form["max"][0])

	var sbooks []Book

	for _, v := range books {
		if v.Price > uint(minp) && v.Price < uint(maxp) {
			sbooks = append(sbooks, v)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // WriteHeader必须在Header().Set()之后

	data, _ := json.Marshal(sbooks)
	fmt.Fprintln(w, string(data))

}

func View(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /view/{id} view View
	//
	// Return a book entity
	// ---
	// consumes:
	// - text/plain
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: The book id
	//   type: string
	//   required: true
	// responses:
	//   200:
	//     descript: The book entity
	//     schema:
	//       "$ref": "#/definitions/Book"
	//   404: MyError
	vars := mux.Vars(r)
	id := vars["id"]

	book, ok := books[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, MyError{
			Message: "item not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // WriteHeader必须在Header().Set()之后

	fmt.Fprintln(w, book)

	// data, err := json.Marshal(book)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintln(w)
	// 	return
	// }

	// w.WriteHeader(http.StatusOK)
	// fmt.Fprintln(w, string(data))
}

func Swagger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, "swagger.json")
}
