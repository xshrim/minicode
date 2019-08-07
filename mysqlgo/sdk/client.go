package sdk

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var config *Config
var clients map[string]*Client

type Client struct {
	Db *sql.DB
}

func New(kind, host, port, dbname, charset, user, passwd string) (*Client, error) { // 不连接具体库
	hashkey := Hash(strings.ToLower(kind), host, port, dbname, charset, user, passwd)
	if client, ok := clients[hashkey]; ok {
		if err := client.Db.Ping(); err == nil {
			return client, nil
		}
	}

	db, err := sql.Open(strings.ToLower(kind), user+":"+passwd+"@tcp("+host+":"+port+")/"+dbname+"?charset="+charset)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	clients[hashkey] = &Client{Db: db}

	// db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})

	return clients[hashkey], nil
}

func (c *Client) Info(sql string) ([]byte, error) {
	if err := c.Db.Ping(); err != nil {
		return nil, err
	}

	query, err := c.Db.Query(sql)
	if err != nil {
		return nil, err
	}

	return parseResult(query)
}

func (c *Client) Query(sql string) ([]byte, error) {
	sql = strings.Replace(sql, "\"", "'", -1)
	//fmt.Println(sql)
	query, err := c.Db.Query(sql)
	if err != nil {
		return nil, err
	}

	return parseResult(query)
}

func (c *Client) Execute(sql string) ([]byte, error) {
	sql = strings.Replace(sql, "\"", "'", -1)
	//fmt.Println(sql)
	res, err := c.Db.Exec(sql)
	if err != nil {
		return nil, err
	}

	var lidstr, ratstr string

	lid, lerr := res.LastInsertId()
	if lerr != nil {
		lidstr = lerr.Error()
	} else {
		lidstr = strconv.FormatUint(uint64(lid), 10)
	}

	rat, rerr := res.RowsAffected()
	if rerr != nil {
		ratstr = rerr.Error()
	} else {
		ratstr = strconv.FormatUint(uint64(rat), 10)
	}

	return []byte("{\"LastInsertId\": \"" + lidstr + "\", \"RowsAffected\": \"" + ratstr + "\"}"), nil
}

// 查询库列表: show databases
// 查询库中所有表:
//   show tables
//   SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = '{database_name}'
// 查询表中所有字段:
//   desc tablename
//    SELECT TABLE_NAME, COLUMN_NAME, DATA_TYPE, COLUMN_COMMENT FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = '{database_name}' AND TABLE_NAME = '{tablename}'

func parseResult(query *sql.Rows) ([]byte, error) {
	column, _ := query.Columns()
	values := make([]sql.RawBytes, len(column))
	scans := make([]interface{}, len(column))

	for i := range values {
		scans[i] = &values[i]
	}

	var results []map[string]string

	i := 0

	for query.Next() {
		if err := query.Scan(scans...); err != nil {
			return nil, err
		}

		row := make(map[string]string)
		for k, v := range values {
			//fmt.Println(v)
			data := string(v)
			//bits.LeadingZeros8(v[0]) == 7
			if len(v) == 0 {
				data = "NULL"
			}
			if len(v) == 1 { // 判断布尔类型
				if v[0] == '\x00' {
					data = "false"
				} else if v[0] == '\x01' {
					data = "true"
				}
			}
			//fmt.Println(b, err)
			row[column[k]] = data
		}

		// fmt.Println(row)

		results = append(results, row)
		i++
	}

	// for k, v := range results {
	// 	fmt.Println(k, v)
	// }

	return json.Marshal(results)
}

func printValue(pval interface{}) string {
	var s_txt string
	switch v := (pval).(type) {
	case nil:
		s_txt = "NULL"
	case time.Time:
		s_txt = "'" + v.Format("2006-01-02 15:04:05.999") + "'"
	case int, int8, int16, int32, int64, float32, float64, byte:
		s_txt = fmt.Sprint(v)
	case []byte:
		s_txt = string(v)
	case bool:
		if v {
			s_txt = "'1'"
		} else {
			s_txt = "'0'"
		}
	default:
		s_txt = "'" + fmt.Sprint(v) + "'"
	}
	return s_txt
}
