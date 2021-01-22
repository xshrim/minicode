package signalscan

import (
	"alertsyslog/src/alertLogHandle"
	"os"
	"os/signal"

	"github.com/xshrim/gol"
)

func ScanSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	gol.Info("Got signal:", s)
	alertLogHandle.DBdataUpdate()
}
