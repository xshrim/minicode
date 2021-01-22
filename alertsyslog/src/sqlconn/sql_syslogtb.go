package sqlconn

import (
	"database/sql"

	"alertsyslog/src/config"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xshrim/gol"
)

var Syslogdb *sql.DB

func init() {
	var dberr error
	Syslogdb, dberr = sql.Open("sqlite3", config.SqliteConfig)
	if dberr != nil || Syslogdb == nil || !Dbping() {
		gol.Error("sqlite3数据库创建失败!")
		return
	}
	smic, smoc := 1, 10
	Syslogdb.SetMaxIdleConns(smic)
	Syslogdb.SetMaxOpenConns(smoc)

	gol.Info("数据库连接成功，连接池完成初始化...")

	if CreateTable() != nil {
		gol.Info("数据库表创建失败")
	}

	gol.Info("数据库表创建成功")

}

func CreateTable() error {
	ssql := "create table if not exists syslogtb(keyword varchar(100), unixtime bigint)"
	psql := "create table if not exists projecttb(ename varchar(100), zhname varchar(100))"
	if err := Syslogdb.Ping(); err != nil {
		return err
	}
	_, err := Syslogdb.Exec(ssql)
	if err != nil {
		return err
	}
	_, err = Syslogdb.Exec(psql)
	if err != nil {
		return err
	}

	return nil
}

func FindProjectName(ename string) (string, bool) {
	findSQL := "select zhname from projecttb where ename=?"
	ret, err := Syslogdb.Query(findSQL, ename)
	if err != nil {
		gol.Error("查找项目（"+ename+"）失败 :", err)
		return err.Error(), false
	}
	defer ret.Close()
	if ret.Next() {
		var zhname string
		if err := ret.Scan(&zhname); err != nil {
			gol.Error("数据库该条数据："+ename+" 读入失败! err:", err)
			return err.Error(), false
		}
		gol.Info("查找记录（" + zhname + "）成功...")
		return zhname, true
	}
	return "查找失败！", false
}

func InsertProjectName(ename string, zhname string) (bool, string) {
	stmt, err := Syslogdb.Prepare("INSERT INTO projecttb (`ename`,`zhname`) VALUES (?,?)")
	if err != nil {
		gol.Error(err)
		return false, err.Error()
	}
	defer stmt.Close()
	_, err = stmt.Exec(ename, zhname)
	if err != nil {
		gol.Error("记录（"+ename+"）插入成功失败 :", err)
		return false, err.Error()
	}
	gol.Info("记录（" + ename + "）已保存至MySQL中！")
	return true, "数据插入成功！"
}

func UpdateProjectName(ename string, zhname string) (bool, string) {
	updateSQL := "UPDATE projecttb SET ename=?,zhname=? WHERE ename=?"
	stmt, err := Syslogdb.Prepare(updateSQL)
	if err != nil {
		gol.Error(err)
		return false, err.Error()
	}
	defer stmt.Close()
	res, err := stmt.Exec(ename, zhname, ename)
	if err != nil {
		gol.Error(err)
		return false, err.Error()
	}
	if res != nil {
		num, _ := res.RowsAffected()
		if num > 0 {
			gol.Info(ename + "项目信息更新成功！")
			return true, "更新成功!"
		}
	}
	return false, "更新失败！"
}

func DeleteProjectName(ename string) (bool, string) {
	deleteSQL := "delete from projecttb where ename=?"
	_, err := Syslogdb.Exec(deleteSQL, ename)
	//stmt, err := Syslogdb.Prepare(deleteSQL)
	if err != nil {
		gol.Error("记录（"+ename+"）删除失败 :", err)
		return false, err.Error()
	}
	gol.Info("记录（" + ename + "）删除成功...")
	return true, "删除成功！"
}

//插入数据
func InsertData(keyword string, unixtime int64) bool {
	res := FindData(keyword)
	if res == 1 {
		insertSQL := "INSERT INTO syslogtb (`keyword`,`unixtime`) VALUES (?,?)"
		_, err := Syslogdb.Exec(insertSQL, keyword, unixtime)
		if err != nil {
			gol.Error("记录（"+keyword+"）插入失败 :", err)
			// 2 表示数据库断开连接
			return false
		}
	} else if res == 0 {
		// 1 表示数据库数据重复，不用再次插入
		return true
	} else {
		// 2  表示数据库无法使用
		return false
	}
	gol.Info("记录（" + keyword + "）已保存至MySQL中！")
	// 0 表示插入成功
	return true
}

func DeleteData(keyword string) int {
	deleteSQL := "delete from syslogtb where keyword=?"
	_, err := Syslogdb.Exec(deleteSQL, keyword)
	if err != nil {
		gol.Error("记录（"+keyword+"）删除失败 :", err)
		return 1
	}
	gol.Info("记录（" + keyword + "）删除成功...")
	return 0
}

func UpdateData(keyword string, unixtime int64) bool {
	updateSQL := "UPDATE syslogtb SET unixtime=? WHERE keyword=?"
	stmt, err := Syslogdb.Prepare(updateSQL)
	if err != nil {
		gol.Error(err)
		return false
	}
	defer stmt.Close()
	res, err := stmt.Exec(unixtime, keyword)
	if err != nil {
		gol.Error(err)
		return false
	}
	if res != nil {
		num, _ := res.RowsAffected()
		if num > 0 {
			gol.Debug("keyword时间戳更新成功！")
			return true
		}
	}
	gol.Error("keyword时间戳更新失败！")
	return false
}

func FindData(keyword string) int {
	findSQL := "select unixtime from syslogtb where keyword=?"
	ret, err := Syslogdb.Query(findSQL, keyword)
	if err != nil {
		gol.Error("查找记录（"+keyword+"）失败 :", err)
		// 2 表示数据库不可用
		return 2
	}
	defer ret.Close()
	if !ret.Next() {
		gol.Info("数据库中无该条（" + keyword + "）记录...")
		// 1 表示数据库中无该条数据
		return 1
	}
	gol.Info("查找记录（" + keyword + "）成功...")
	// 0 表示数据库中有该条数据库，查找成功！
	return 0
}

func MemDataInit() map[string]int64 {
	gol.Info("正在进行内存数据初始化...")
	selectSQL := "select keyword,unixtime from syslogtb"
	rows, err := Syslogdb.Query(selectSQL)
	if err != nil {
		gol.Error("内存数据初始化失败！无法访问数据库 ：", err)
	}
	defer rows.Close()
	dict := make(map[string]int64)
	for rows.Next() {
		var keyword string
		var unixtime int64
		if err := rows.Scan(&keyword, &unixtime); err != nil {
			gol.Error("数据库该条数据："+keyword+" 读入内存失败! err:", err)
			continue
		}
		dict[keyword] = unixtime
	}
	//gol.Info(time.Now().Unix())
	gol.Info("内存数据初始化成功！")
	return dict
}

func Dbping() bool {
	selectSQL := "select keyword,unixtime from syslogtb where keyword = '1'"
	rows, err := Syslogdb.Query(selectSQL)
	if err != nil && CreateTable() != nil {
		gol.Error("无法访问数据库：", err)
		return false
	}
	if rows != nil {
		rows.Close()
	}
	return true
}
