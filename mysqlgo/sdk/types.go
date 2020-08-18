package sdk

import "net/http"

type Config struct {
	Kind    string `json:"kind"`
	Host    string `json:"Host"`
	Port    string `json:"port"`
	Dbname  string `json:"dbname"`
	Charset string `json:"charset"`
	User    string `json:"user"`
	Passwd  string `json:"passwd"`
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
