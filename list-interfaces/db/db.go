package db

import (
	"database/sql"

	"../xlog"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func New() (*DB, error) {
	d, err := sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared&loc=auto")
	if err != nil {
		xlog.Errorln("Create database error: ", err)
		return nil, err
	}

	return &DB{
		db: d,
	}, nil
}

func (d *DB) Close() {
	d.db.Close()
}

func (d *DB) CreateTable(sql string) bool {
	_, err := d.db.Exec(sql)
	if err != nil {
		xlog.Errorln("Create table error: ", err)
		return false
	}
	return true
}

func (d *DB) CreateIndex(sql string) bool {
	_, err := d.db.Exec(sql)
	if err != nil {
		xlog.Errorln("Create index error: ", err)
		return false
	}

	return true
}

func (d *DB) Insert(sql string, args ...interface{}) (int64, error) {
	res, err := d.db.Exec(sql, args)
	if err != nil {
		xlog.Errorln("Insert row error: ", err)
		return 0, err
	}

	return res.LastInsertId()
}

func (d *DB) Update(sql string, args ...interface{}) (int64, error) {
	res, err := d.db.Exec(sql, args)
	if err != nil {
		xlog.Errorln("Update row error: ", err)
		return 0, err
	}

	return res.RowsAffected()
}

func (d *DB) Delete(sql string, args ...interface{}) (int64, error) {
	res, err := d.db.Exec(sql, args)
	if err != nil {
		xlog.Errorln("Delete rows error: ", err)
		return 0, err
	}

	return res.RowsAffected()
}

func (d *DB) Query(sql string, args ...interface{}) []interface{} {
	res, err := d.db.Query(sql, args)
	if err != nil {
		xlog.Errorln("Query rows error: ", err)
	}

	var rows []interface{}
	for res.Next() {
		var row interface{}
		if err := res.Scan(row); err != nil {
			xlog.Errorln("Fetch row error: ", err)
		} else {
			rows = append(rows, row)
		}
	}

	return rows
}
