package rest

var books = map[string]Book{
	"1": Book{
		Id:      "1",
		Authors: []string{"Tom", "Lucy"},
		Title:   "The Golang Programming",
		Price:   100,
	},
	"2": Book{
		Id:      "2",
		Authors: []string{"Nick", "Susan"},
		Title:   "The Python Programming",
		Price:   80,
	},
	"3": Book{
		Id:      "3",
		Authors: []string{"Jack", "Lina"},
		Title:   "The Javascript Programming",
		Price:   90,
	},
}
