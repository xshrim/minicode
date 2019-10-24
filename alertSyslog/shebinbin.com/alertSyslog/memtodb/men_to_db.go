package memtodb

import (
	"time"

	"shebinbin.com/alertSyslog/alertLogHandle"
	"shebinbin.com/alertSyslog/config"
	"shebinbin.com/alertSyslog/heartbeat"
	"shebinbin.com/alertSyslog/mysqlconn"
	"shebinbin.com/alertSyslog/zapLogger"
)

var logger = zapLogger.LoggerFactory()
var count = 0

func countPlus() int {
	if count > config.CountTime {
		count = 0
	}
	count++
	return count
}

func init() {
	go func() {
		cron := time.NewTicker(time.Second * time.Duration(config.CronTime))
		for {
			// 任务 // cron.Stop()
			select {
			case <-cron.C:
				men2db()
			}
		}
	}()
}

// 如果当前数据库不可用，执行一次数据库存活检查，如果检查到数据库健康，那么开始将执行队列中的数据依次执行一遍。
func men2db() {
	if !config.Dbping {
		if countPlus() == 1 {
			heartbeat.AlertDbErr(config.MonitorIP, config.MonitorPort, config.NodeIP)
		}
		/*if config.ShowMemValue {
			alertLogHandle.PrintMemData()
		}*/
		logger.Info("数据库已断开一段时间，正在重连...")
		if mysqlconn.Dbping() {
			//mysqlconn.ReInit()
			config.Dbping = true
			// 那么开始将执行队列中的数据依次执行一遍
			alertLogHandle.DBdataUpdate()
			count = 0

		} else {
			logger.Error("数据库重连失败！")
		}
	} /*else{
		//logger.Info("定时任务...")
		//alertLogHandle.PrintMemData()
	}*/
}
