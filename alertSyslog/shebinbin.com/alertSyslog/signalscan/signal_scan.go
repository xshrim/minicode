package signalscan

import (
	"os"
	"os/signal"
	"shebinbin.com/alertSyslog/alertLogHandle"
	"shebinbin.com/alertSyslog/zapLogger"
)

var logger = zapLogger.LoggerFactory()

func ScanSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	logger.Info("Got signal:", s)
	alertLogHandle.DBdataUpdate()
}
