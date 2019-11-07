package mysqlconn

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"shebinbin.com/alertSyslog/config"
	"shebinbin.com/alertSyslog/zapLogger"
)

var logger = zapLogger.LoggerFactory()
var Syslogdb *sql.DB

func init() {
	var dberr error
	Syslogdb, dberr = sql.Open("mysql", config.MysqlConf)
	if dberr != nil || Syslogdb == nil || !Dbping() {
		logger.Error("mysql数据库连接失败!")
		config.Dbping = false
		return
	}
	smic, smoc := 1, 10
	Syslogdb.SetMaxIdleConns(smic)
	Syslogdb.SetMaxOpenConns(smoc)

	logger.Info("数据库连接成功，连接池完成初始化...")

	if CreateTable() != nil {
		logger.Info("数据库表创建失败")
	}

	logger.Info("数据库表创建成功")

}

func CreateTable() error {
	ssql := "create table if not exists syslogtb(keyword varchar(100), unixtime bigint)"
	psql := "create table if not exists projecttb(ename varchar(100), zhname varchar(100))"
	if err := Syslogdb.Ping(); err != nil {
		return err
	}
	_, err := Syslogdb.Query(ssql)
	if err != nil {
		return err
	}
	_, err = Syslogdb.Query(psql)
	if err != nil {
		return err
	}

	return nil
}

func FindProjectName(ename string) (string, bool) {
	findSql := "select zhname from projecttb where ename=?"
	ret, err := Syslogdb.Query(findSql, ename)
	if err != nil {
		logger.Error("查找项目（"+ename+"）失败 :", err)
		return err.Error(), false
	}
	defer ret.Close()
	if ret.Next() {
		var zhname string
		if err := ret.Scan(&zhname); err != nil {
			logger.Error("数据库该条数据："+ename+" 读入失败! err:", err)
			return err.Error(), false
		}
		logger.Info("查找记录（" + zhname + "）成功...")
		return zhname, true
	}
	return "查找失败！", false
}

func InsertProjectName(ename string, zhname string) (bool, string) {
	stmt, err := Syslogdb.Prepare("INSERT INTO projecttb (`ename`,`zhname`) VALUES (?,?)")
	if err != nil {
		logger.Error(err)
		return false, err.Error()
	}
	defer stmt.Close()
	_, err = stmt.Exec(ename, zhname)
	if err != nil {
		logger.Error("记录（"+ename+"）插入成功失败 :", err)
		return false, err.Error()
	}
	logger.Info("记录（" + ename + "）已保存至MySQL中！")
	return true, "数据插入成功！"
}

func UpdateProjectName(ename string, zhname string) (bool, string) {
	updateSql := "UPDATE projecttb SET ename=?,zhname=? WHERE ename=?"
	stmt, err := Syslogdb.Prepare(updateSql)
	if err != nil {
		logger.Error(err)
		return false, err.Error()
	}
	defer stmt.Close()
	res, err := stmt.Exec(ename, zhname, ename)
	if err != nil {
		logger.Error(err)
		return false, err.Error()
	}
	if res != nil {
		num, _ := res.RowsAffected()
		if num > 0 {
			logger.Info(ename + "项目信息更新成功！")
			return true, "更新成功!"
		}
	}
	return false, "更新失败！"
}

func DeleteProjectName(ename string) (bool, string) {
	deleteSql := "delete from projecttb where ename=?"
	_, err := Syslogdb.Exec(deleteSql, ename)
	//stmt, err := Syslogdb.Prepare(deleteSql)
	if err != nil {
		logger.Error("记录（"+ename+"）删除失败 :", err)
		return false, err.Error()
	}
	logger.Info("记录（" + ename + "）删除成功...")
	return true, "删除成功！"
}

//插入数据
func InsertData(keyword string, unixtime int64) bool {
	res := FindData(keyword)
	if res == 1 {
		insertSql := "INSERT INTO syslogtb (`keyword`,`unixtime`) VALUES (?,?)"
		_, err := Syslogdb.Exec(insertSql, keyword, unixtime)
		if err != nil {
			logger.Error("记录（"+keyword+"）插入失败 :", err)
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
	logger.Info("记录（" + keyword + "）已保存至MySQL中！")
	// 0 表示插入成功
	return true
}

func DeleteData(keyword string) int {
	deleteSql := "delete from syslogtb where keyword=?"
	_, err := Syslogdb.Exec(deleteSql, keyword)
	if err != nil {
		logger.Error("记录（"+keyword+"）删除失败 :", err)
		return 1
	}
	logger.Info("记录（" + keyword + "）删除成功...")
	return 0
}

func UpdateData(keyword string, unixtime int64) bool {
	updateSql := "UPDATE syslogtb SET unixtime=? WHERE keyword=?"
	stmt, err := Syslogdb.Prepare(updateSql)
	if err != nil {
		logger.Error(err)
		return false
	}
	defer stmt.Close()
	res, err := stmt.Exec(unixtime, keyword)
	if err != nil {
		logger.Error(err)
		return false
	}
	if res != nil {
		num, _ := res.RowsAffected()
		if num > 0 {
			logger.Debug("keyword时间戳更新成功！")
			return true
		}
	}
	logger.Error("keyword时间戳更新失败！")
	return false
}

func FindData(keyword string) int {
	findSql := "select unixtime from syslogtb where keyword=?"
	ret, err := Syslogdb.Query(findSql, keyword)
	if err != nil {
		logger.Error("查找记录（"+keyword+"）失败 :", err)
		// 2 表示数据库不可用
		return 2
	}
	defer ret.Close()
	if !ret.Next() {
		logger.Info("数据库中无该条（" + keyword + "）记录...")
		// 1 表示数据库中无该条数据
		return 1
	}
	logger.Info("查找记录（" + keyword + "）成功...")
	// 0 表示数据库中有该条数据库，查找成功！
	return 0
}

func MemDataInit() map[string]int64 {
	logger.Info("正在进行内存数据初始化...")
	selectSql := "select keyword,unixtime from syslogtb"
	rows, err := Syslogdb.Query(selectSql)
	if err != nil {
		logger.Error("内存数据初始化失败！无法访问数据库 ：", err)
	}
	defer rows.Close()
	dict := make(map[string]int64)
	for rows.Next() {
		var keyword string
		var unixtime int64
		if err := rows.Scan(&keyword, &unixtime); err != nil {
			logger.Error("数据库该条数据："+keyword+" 读入内存失败! err:", err)
			continue
		}
		dict[keyword] = unixtime
	}
	//logger.Info(time.Now().Unix())
	logger.Info("内存数据初始化成功！")
	return dict
}

func Dbping() bool {
	selectSql := "select keyword,unixtime from syslogtb where keyword = '1'"
	rows, err := Syslogdb.Query(selectSql)
	if err != nil {
		if CreateTable() != nil {
			logger.Error("无法访问数据库 ：", err)
			return false
		}
	}
	if rows != nil {
		rows.Close()
	}
	return true
}
