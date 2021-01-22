package main

import (
	"auditlog/database"
	"auditlog/gc"
	"auditlog/server"
	"flag"
	"fmt"
	"time"
)

var expireTimeStr = flag.String("expire", "7d", "expire time of audit logs since created, default is 7 days")
var gcIntervalStr = flag.String("gc", "3h", "interval seconds for cleaning expired audit logs, default is 3 hours")
var serverIP = flag.String("ip", "0.0.0.0", "server ip address, default is 0.0.0.0")
var serverPort = flag.Int("port", 9090, "server port, default is 9090")
var noheader = flag.Bool("no-header", false, "discard request and response headers in audit logs")
var nobody = flag.Bool("no-body", false, "discard request and response bodys in audit logs")
var dbtype = flag.String("dbtype", "sqlite3", "database type, default is sqlite3")
var dburl = flag.String("dburl", "./data.db", "database connect url, default is ./data.db")

func main() {
	flag.Parse()

	expireTime, _ := time.ParseDuration(*expireTimeStr)
	gcInterval, _ := time.ParseDuration(*gcIntervalStr)

	db, err := database.Connect(*dbtype, *dburl)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	go gc.Gc(expireTime, gcInterval, db)
	fmt.Printf("Serving at %s:%d ...\n", *serverIP, *serverPort)
	server.New(*serverIP, *serverPort, *noheader, *nobody, db).Serve()
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// <-c
	// db.Close()
	// return
	// data := `{
	//   "requestBody": {
	//     "publicEndpoints": [
	//       {
	//           "addresses": [
	//               "10.64.3.58"
	//           ],
	//           "allNodes": true,
	//           "ingressId": {},
	//           "nodeId": [],
	//           "podId": null
	//       }
	//     ]
	//   }
	// }
	//   `

	// v := make(map[string]interface{})
	// _ = json.Unmarshal([]byte(data), &v)
	// fmt.Println(utils.Jsonify(v))
}
