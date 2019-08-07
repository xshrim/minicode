package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func parseParam(r *http.Request) (Config, string) {
	var sqls string

	r.ParseForm()

	var cfg Config
	if config != nil {
		cfg = *config
	} else {
		cfg = Config{}
	}

	if rkind := r.Form["kind"]; len(rkind) > 0 {
		if rkind[0] != "" && rkind[0] != "null" {
			cfg.Kind = strings.TrimSpace(rkind[0])
		}
	}

	if rhost := r.Form["host"]; len(rhost) > 0 {
		if rhost[0] != "" && rhost[0] != "null" {
			cfg.Host = strings.TrimSpace(rhost[0])
		}
	}

	if rport := r.Form["port"]; len(rport) > 0 {
		if rport[0] != "" && rport[0] != "null" {
			cfg.Port = strings.TrimSpace(rport[0])
		}
	}

	if rdbname := r.Form["dbname"]; len(rdbname) > 0 {
		if rdbname[0] != "" && rdbname[0] != "null" {
			cfg.Dbname = strings.TrimSpace(rdbname[0])
		}
	}

	if rcharset := r.Form["charset"]; len(rcharset) > 0 {
		if rcharset[0] != "" && rcharset[0] != "null" {
			cfg.Charset = strings.TrimSpace(rcharset[0])
		}
	}

	if ruser := r.Form["user"]; len(ruser) > 0 {
		if ruser[0] != "" && ruser[0] != "null" {
			cfg.User = strings.TrimSpace(ruser[0])
		}
	}

	if rpasswd := r.Form["passwd"]; len(rpasswd) > 0 {
		if rpasswd[0] != "" && rpasswd[0] != "null" {
			cfg.Passwd = strings.TrimSpace(rpasswd[0])
		}
	}

	if rsqls := r.Form["sql"]; len(rsqls) > 0 {
		sqls = strings.TrimSpace(rsqls[0])
	}

	sqls = strings.Replace(sqls, "\n", " ", -1)

	return cfg, sqls
}

func Exe(w http.ResponseWriter, r *http.Request) {
	// var kind, host, port, dbname, charset, user, passwd, sql, info string
	cfg, sqls := parseParam(r)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("content-type", "application/json")

	// vars := mux.Vars(r)
	// if v, ok := vars["info"]; ok {
	// 	info = v
	// }
	fmt.Println(cfg, sqls)

	client, err := New(cfg.Kind, cfg.Host, cfg.Port, cfg.Dbname, cfg.Charset, cfg.User, cfg.Passwd)
	if err != nil {
		fmt.Fprintf(w, "[{\"ERROR\": \""+err.Error()+"\"}]")
		return
	}

	var datas [][2]string
	for _, sql := range strings.Split(sqls, ";") {
		var res [2]string
		sql = strings.TrimSpace(sql)
		if sql == "" {
			continue
		}
		res[0] = sql
		stateprefix := strings.ToLower(strings.Split(sql, " ")[0])
		if stateprefix == "insert" || stateprefix == "update" || stateprefix == "delete" || stateprefix == "drop" || stateprefix == "alter" || stateprefix == "create" { // 写库
			data, err := client.Execute(sql)
			if err != nil {
				res[1] = "{\"ERROR\": \"" + err.Error() + "\"}"
			} else {
				res[1] = string(data)
			}
		} else { // 读库
			data, err := client.Query(sql)
			if err != nil {
				res[1] = "{\"ERROR\": \"" + err.Error() + "\"}"
			} else {
				res[1] = string(data)
			}
		}
		datas = append(datas, res)
	}

	dbytes, err := json.Marshal(datas)
	if err != nil {
		fmt.Fprintf(w, "[{\"ERROR\": \""+err.Error()+"\"}]")
		//client.Db.Close()
		return
	}

	fmt.Fprintf(w, string(dbytes))
	//fmt.Println(string(dbytes))
	//client.Db.Close()
}

func Info(w http.ResponseWriter, r *http.Request) {
	// var kind, host, port, dbname, charset, user, passwd, sql, info string

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("content-type", "application/json")

	cfg, sqls := parseParam(r)

	// vars := mux.Vars(r)
	// if v, ok := vars["info"]; ok {
	// 	info = v
	// }
	fmt.Println(sqls)

	client, err := New(cfg.Kind, cfg.Host, cfg.Port, cfg.Dbname, cfg.Charset, cfg.User, cfg.Passwd)
	if err != nil {
		fmt.Fprintf(w, "ERROR:"+err.Error())
		//client.Db.Close()
		return
	}

	data, err := client.Info(sqls)
	if err != nil {
		fmt.Fprintf(w, "ERROR:"+err.Error())
	} else {
		fmt.Fprintf(w, string(data))
	}

	//client.Db.Close()
}
