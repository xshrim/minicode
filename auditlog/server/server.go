package server

import (
	"auditlog/database"
	"auditlog/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type User struct {
	Name  string   `json:"name"`
	Group []string `json:"group"`
}

type AuditLog struct {
	AuditID           string      `json:"auditID"`
	RequestURI        string      `json:"requestURI"`
	RemoteAddr        string      `json:"remoteAddr"`
	User              User        `json:"user"`
	Method            string      `json:"method"`
	ResponseCode      int         `json:"responseCode"`
	RequestTimestamp  string      `json:"requestTimestamp"`
	ResponseTimestamp string      `json:"responseTimestamp"`
	RequestHeader     interface{} `json:"requestHeader"`
	ResponseHeader    interface{} `json:"responseHeader"`
	RequestBody       interface{} `json:"requestBody"`
	ResponseBody      interface{} `json:"responseBody"`
}

type Request struct {
	Offset         int64  `json:"offset"`
	Limit          int    `json:"limit"`
	AuditID        string `json:"auditID"`
	User           string `json:"user"`
	Scope          int    `json:"scope"`
	RemoteAddr     string `json:"remoteAddr"`
	Cluster        string `json:"cluster"`
	Project        string `json:"project"`
	Method         string `json:"method"`
	ResponseCode   int    `json:"responseCode"`
	StartTimestamp int64  `json:"startTimestamp"`
	EndTimestamp   int64  `json:"endTimestamp"`
}

type Row struct {
	AuditID    string    `json:"auditID"`
	Scope      int       `json:"scope"`
	Requri     string    `json:"requri"`
	RemoteAddr string    `json:"remoteAddr"`
	Uname      string    `json:"uname"`
	Ugroup     []string  `json:"ugroup"`
	Cluster    string    `json:"cluster"`
	Project    string    `json:"project"`
	Method     string    `json:"method"`
	RespCode   int       `json:"responseCode"`
	Reqts      time.Time `json:"requestTimestamp"`
	Respts     time.Time `json:"responseTimestamp"`
	Reqheader  string    `json:"requestHeader"`
	Respheader string    `json:"responseHeader"`
	Reqbody    string    `json:"requestBody"`
	Respbody   string    `json:"responseBody"`
}

type Response struct {
	Offset int64 `json:"offset"`
	Rows   []Row `json:"rows"`
}

type Server struct {
	host     string
	port     int
	noheader bool
	nobody   bool
	db       *database.Database
}

func New(host string, port int, nh, nb bool, db *database.Database) *Server {
	if host == "" {
		host = "0.0.0.0"
	}
	return &Server{
		host:     host,
		port:     port,
		noheader: nh,
		nobody:   nb,
		db:       db,
	}
}

func write(db *database.Database, log *AuditLog, noheader, nobody bool) {
	scope := 1
	cluster := ""
	project := ""
	uname := log.User.Name
	ugroup := strings.Join(log.User.Group, ",")
	projectExp := regexp.MustCompile(`^/v3/projects?/([\w\-]*):([\w\-]*)`)
	clusterExp := regexp.MustCompile(`^/v3/clusters?/([\w\-]*)`)

	if projectExp.Match([]byte(log.RequestURI)) {
		scope = 3
		m := projectExp.FindStringSubmatch(log.RequestURI)
		cluster = m[1]
		project = m[2]
	} else if clusterExp.Match([]byte(log.RequestURI)) {
		scope = 2
		m := clusterExp.FindStringSubmatch(log.RequestURI)
		cluster = m[1]
	}
	remoteAddr := log.RemoteAddr

	reqtm, _ := time.ParseInLocation("2006-01-02T15:04:05+08:00", log.RequestTimestamp, time.Local)
	requestTimestamp := reqtm.Unix()
	resptm, _ := time.ParseInLocation("2006-01-02T15:04:05+08:00", log.ResponseTimestamp, time.Local)
	responseTimestamp := resptm.Unix()

	reqheader := ""
	respheader := ""
	if !noheader {
		if log.RequestHeader != nil {
			reqheader = utils.Jsonify(log.RequestHeader.(map[string]interface{}))
		}
		if log.ResponseHeader != nil {
			respheader = utils.Jsonify(log.ResponseHeader.(map[string]interface{}))
		}
	}

	reqbody := ""
	respbody := ""
	if !nobody {
		if log.RequestBody != nil {
			reqbody = utils.Jsonify(log.RequestBody.(map[string]interface{}))
		}
		if log.ResponseBody != nil {
			respbody = utils.Jsonify(log.ResponseBody.(map[string]interface{}))
		}
	}

	sqlstr := "INSERT INTO auditlog(auditid, scope, requri, remoteaddr, uname, ugroup, cluster, project, method, respcode, reqts, respts, reqheader, respheader, reqbody, respbody) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	args := []interface{}{log.AuditID, scope, log.RequestURI, remoteAddr, uname, ugroup, cluster, project, log.Method, log.ResponseCode, requestTimestamp, responseTimestamp, reqheader, respheader, reqbody, respbody}
	if n, err := db.Insert(sqlstr, args...); err != nil {
		fmt.Printf("Inserted %d rows, error: %s\n", n, err.Error())
	}
}

func (svr *Server) writeHandle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		payload, _ := ioutil.ReadAll(r.Body)
		var auditLog AuditLog
		err := json.Unmarshal(payload, &auditLog)
		if err != nil {
			fmt.Println(err)
		}
		write(svr.db, &auditLog, svr.noheader, svr.nobody)
	}
}

func read(db *database.Database, req Request) ([]byte, error) {
	sqlstr := "SELECT * FROM auditlog WHERE"
	if req.AuditID != "" {
		sqlstr += fmt.Sprintf(" auditid='%s' AND", req.AuditID)
	}
	if req.User != "" {
		sqlstr += fmt.Sprintf(" uname='%s' AND", req.User)
	}
	if req.Scope != 0 {
		sqlstr += fmt.Sprintf(" scope=%d AND", req.Scope)
	}
	if req.Cluster != "" {
		sqlstr += fmt.Sprintf(" cluster='%s' AND", req.Cluster)
	}
	if req.Project != "" {
		sqlstr += fmt.Sprintf(" project='%s' AND", req.Project)
	}
	if req.RemoteAddr != "" {
		sqlstr += fmt.Sprintf(" remoteaddr='%s' AND", req.RemoteAddr)
	}
	if req.Method != "" {
		sqlstr += fmt.Sprintf(" method='%s' AND", req.Method)
	}
	if req.ResponseCode != 0 {
		sqlstr += fmt.Sprintf(" respcode='%d' AND", req.ResponseCode)
	}
	if req.StartTimestamp != 0 && req.EndTimestamp != 0 {
		sqlstr += fmt.Sprintf(" reqts BETWEEN %d AND %d AND", req.StartTimestamp, req.EndTimestamp)
	} else {
		if req.StartTimestamp != 0 {
			sqlstr += fmt.Sprintf(" reqts>=%d AND", req.StartTimestamp)
		} else {
			sqlstr += fmt.Sprintf(" reqts<=%d AND", req.EndTimestamp)
		}
	}
	if req.Offset != 0 {
		sqlstr += fmt.Sprintf(" id>=%d AND", req.Offset)
	}

	sqlstr = strings.TrimRight(sqlstr, " AND")
	sqlstr = strings.TrimRight(sqlstr, " WHERE")

	if req.Limit != 0 {
		sqlstr += fmt.Sprintf(" LIMIT %d ORDER BY id", req.Limit)
	}

	rs, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}

	var auditID string
	var scope int
	var requri string
	var remoteAddr string
	var uname string
	var ugroup string
	var cluster string
	var project string
	var method string
	var responseCode int
	var requestTimestamp time.Time
	var responseTimestamp time.Time
	var reqbody string
	var respbody string
	var reqheader string
	var respheader string

	var rows []Row
	var offset int64

	for rs.Next() {
		err := rs.Scan(&offset, &auditID, &scope, &requri, &remoteAddr, &uname, &ugroup, &cluster, &project, &method, &responseCode, &requestTimestamp, &responseTimestamp, &reqheader, &respheader, &reqbody, &respbody)
		if err != nil {
			fmt.Println(err)
		}
		row := Row{
			AuditID:    auditID,
			Scope:      scope,
			Requri:     requri,
			RemoteAddr: remoteAddr,
			Uname:      uname,
			Ugroup:     strings.Split(ugroup, ","),
			Cluster:    cluster,
			Project:    project,
			Method:     method,
			RespCode:   responseCode,
			Reqts:      requestTimestamp.Local(),
			Respts:     responseTimestamp.Local(),
			Reqheader:  reqheader,
			Respheader: respheader,
			Reqbody:    reqbody,
			Respbody:   respbody,
		}
		rows = append(rows, row)
		// fmt.Println(id, scope, requri, remoteAddr, uname, user, cluster, project, method, stage, timestamp, reqbody, respbody)
	}
	resp := Response{
		Offset: offset + 1,
		Rows:   rows,
	}
	return json.Marshal(resp)
}

func (svr *Server) readHandle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == http.MethodGet {
		payload, _ := ioutil.ReadAll(r.Body)
		var req Request
		err := json.Unmarshal(payload, &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		data, err := read(svr.db, req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		} else {
			_, _ = w.Write(data)
		}
	}
}

func (svr *Server) pingHandle(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Pong"))
}

func (svr *Server) Serve() error {
	http.HandleFunc("/set", svr.writeHandle)
	http.HandleFunc("/get", svr.readHandle)
	http.HandleFunc("/ping", svr.pingHandle)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", svr.host, svr.port), nil)
}
