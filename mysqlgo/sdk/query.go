package sdk

import (
	"fmt"
	"net/http"
	"strings"
)

func Query(w http.ResponseWriter, r *http.Request) {
	// var kind, host, port, dbname, charset, user, passwd, sql, info string
	var sql string

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("content-type", "application/json")

	r.ParseForm()

	var cfg Config
	if config != nil {
		cfg = *config
	} else {
		cfg = Config{}
	}

	if rkind := r.Form["kind"]; len(rkind) > 0 {
		cfg.Kind = strings.TrimSpace(rkind[0])
	}

	if rhost := r.Form["host"]; len(rhost) > 0 {
		cfg.Host = strings.TrimSpace(rhost[0])
	}

	if rport := r.Form["port"]; len(rport) > 0 {
		cfg.Port = strings.TrimSpace(rport[0])
	}

	if rdbname := r.Form["dbname"]; len(rdbname) > 0 {
		cfg.Dbname = strings.TrimSpace(rdbname[0])
	}

	if rcharset := r.Form["charset"]; len(rcharset) > 0 {
		cfg.Charset = strings.TrimSpace(rcharset[0])
	}

	if ruser := r.Form["user"]; len(ruser) > 0 {
		cfg.User = strings.TrimSpace(ruser[0])
	}

	if rpasswd := r.Form["passwd"]; len(rpasswd) > 0 {
		cfg.Passwd = strings.TrimSpace(rpasswd[0])
	}

	if rsql := r.Form["sql"]; len(rsql) > 0 {
		sql = strings.TrimSpace(rsql[0])
	}

	// vars := mux.Vars(r)
	// if v, ok := vars["info"]; ok {
	// 	info = v
	// }
	fmt.Println(sql)

	client, err := New(cfg.Kind, cfg.Host, cfg.Port, cfg.Dbname, cfg.Charset, cfg.User, cfg.Passwd)
	if err != nil {
		fmt.Fprintf(w, "{\"ERROR\": \""+err.Error()+"\"}")
		client.Db.Close()
		return
	}

	stateprefix := strings.ToLower(strings.Split(sql, " ")[0])
	if stateprefix == "insert" || stateprefix == "update" || stateprefix == "delete" || stateprefix == "drop" || stateprefix == "alter" || stateprefix == "create" { // 写库
		data, err := client.Execute(sql)
		if err != nil {
			fmt.Fprintf(w, "{\"ERROR\": \""+err.Error()+"\"}")
		} else {
			fmt.Fprintf(w, string(data))
		}
	} else { // 读库
		data, err := client.Query(sql)
		if err != nil {
			fmt.Fprintf(w, "{\"ERROR\": \""+err.Error()+"\"}")
		} else {
			fmt.Fprintf(w, string(data))
		}
	}
	client.Db.Close()
}
