package main

import (
	"log"
	"net/http"

	"./sdk"
)

/* 实现连接不断开的情况下切换数据库
tx := db.MustBegin() // start transaction
tx.MustExec("use " + userDB) // switch to tenant db
tx.MustExec("insert into ....") // do some work
tx.MustExec("use `no-op-db`") // switch away from tenant db (there is no unuse, so I just use a dummy)
tx.Commit() // end transaction
*/

func main() {

	router := sdk.NewRouter()

	// Handle API routes
	// api := router.PathPrefix("/api/").Subrouter()
	// api.HandleFunc("/student", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintln(w, "From the API")
	// })

	// Serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Serve index page on all unhandled routes
	// router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./index.html")
	// })

	log.Fatal(http.ListenAndServe(":8080", router))

	// sql := "insert into blockchain.user values(0, '创始', '123456', '1980-01-01 02:00:00', false, NULL)"

	// sql = strings.Replace(sql, "\"", "'", -1)

	// res, err := db.Exec(sql)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// liid, _ := res.LastInsertId()
	// rnum, _ := res.RowsAffected()
	// fmt.Println(liid, rnum)
}
