package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var sql_table = `
CREATE TABLE IF NOT EXISTS auditlog(
  id INTEGER PRIMARY KEY autoincrement,
  auditid VARCHAR(64) PRIMARY KEY NOT NULL,
  scope INTEGER NULL,
  requri VARCHAR(128) NULL,
  remoteaddr VARCHAR(24) NULL,
  uname VARCHAR(32) NULL,
  ugroup VARCHAR(256) NULL,
  cluster VARCHAR(32) NULL,
  project VARCHAR(32) NULL,
  method VARCHAR(16) NULL,
  respcode INTEGER NULL,
  reqts TIMESTAMP NULL,
  respts TIMESTAMP NULL,
  reqheader TEXT NULL,
  respheader TEXT NULL,
  reqbody TEXT NULL,
  respbody TEXT NULL
);
`

var sql_indices = []string{
	`
  CREATE INDEX scope ON auditlog (scope);
  `,
	`
  CREATE INDEX uname ON auditlog (uname);
  `,
	`
  CREATE INDEX cluster ON auditlog (cluster);
  `,
	`
  CREATE INDEX project ON auditlog (project);
  `,
	`
  CREATE INDEX method ON auditlog (method);
  `,
	`
  CREATE INDEX respcode ON auditlog (respcode);
  `,
	`
  CREATE INDEX reqts ON auditlog (reqts);
  `,
}

type Database struct {
	kind string
	db   *sql.DB
}

func Connect(kind, dburl string) (*Database, error) {
	var database *Database
	switch kind {
	case "sqlite3":
		db, err := sql.Open("sqlite3", dburl)
		if err != nil {
			return nil, err
		}
		database = &Database{kind: kind, db: db}
	}

	if database.db != nil {
		_, err := database.CreateTable(sql_table)
		if err != nil {
			return nil, err
		}
		for _, sql_index := range sql_indices {
			_, _ = database.CreateIndex(sql_index)
		}
	}
	return database, nil
}

func (db *Database) CreateTable(sqlstr string, args ...interface{}) (sql.Result, error) {
	return db.Execute(sqlstr)
}

func (db *Database) CreateIndex(sqlstr string, args ...interface{}) (sql.Result, error) {
	return db.Execute(sqlstr)
}

func (db *Database) Query(sqlstr string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := db.db.Prepare(sqlstr)
	if err != nil {
		return nil, err
	}
	return stmt.Query(args...)
}

func (db *Database) Execute(sqlstr string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.db.Prepare(sqlstr)
	if err != nil {
		return nil, err
	}
	return stmt.Exec(args...)
}

func (db *Database) Insert(sqlstr string, args ...interface{}) (int64, error) {
	res, err := db.Execute(sqlstr, args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (db *Database) Update(sqlstr string, args ...interface{}) (int64, error) {
	res, err := db.Execute(sqlstr, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (db *Database) Delete(sqlstr string, args ...interface{}) (int64, error) {
	res, err := db.Execute(sqlstr, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (db *Database) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}
