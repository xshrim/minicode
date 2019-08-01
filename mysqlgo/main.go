package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

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

func printResult(query *sql.Rows) {
	column, _ := query.Columns()
	values := make([]sql.RawBytes, len(column))
	scans := make([]interface{}, len(column))

	for i := range values {
		scans[i] = &values[i]
	}

	results := make(map[int]map[string]string)

	i := 0

	for query.Next() {
		if err := query.Scan(scans...); err != nil {
			fmt.Println(err)
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

		results[i] = row
		i++
	}

	for k, v := range results {
		fmt.Println(k, v)
	}
}

/* 实现连接不断开的情况下切换数据库
tx := db.MustBegin() // start transaction
tx.MustExec("use " + userDB) // switch to tenant db
tx.MustExec("insert into ....") // do some work
tx.MustExec("use `no-op-db`") // switch away from tenant db (there is no unuse, so I just use a dummy)
tx.Commit() // end transaction
*/

func New(dbtype, host, port, user, passwd, dbname, charset string) (*sql.DB, error) {
	db, err := sql.Open(dbtype, user+":"+passwd+"@tcp("+host+":"+port+")/"+dbname+"?charset="+charset)
	if err != nil {
		return nil, err
	}

	if db.Ping() != nil {
		return nil, err
	}

	return db, nil
}

func query(sql string) {

}
func execute(sql string) {

}

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/?charset=utf8")
	if err != nil {
		fmt.Println(err)
	}

	// q, err := db.Query("desc blockchain.user")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// printResult(q)

	query, _ := db.Query("select * from blockchain.user")

	printResult(query)

	return

	sql := "insert into blockchain.user values(0, '创始', '123456', '1980-01-01 02:00:00', false, NULL)"

	sql = strings.Replace(sql, "\"", "'", -1)

	res, err := db.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
	liid, _ := res.LastInsertId()
	rnum, _ := res.RowsAffected()
	fmt.Println(liid, rnum)
}
