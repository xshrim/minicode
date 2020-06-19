package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"label.admission/webhook"
	"label.admission/xlog"
)

func main() {
	var port int
	var certfile, keyfile string

	xlog.Level = xlog.DEBUG

	// get command line parameters
	flag.IntVar(&port, "port", 443, "Webhook server port.")
	flag.StringVar(&certfile, "tlsCertFile", "/etc/webhook/certs/cert.pem", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&keyfile, "tlsKeyFile", "/etc/webhook/certs/key.pem", "File containing the x509 private key to --tlsCertFile.")
	flag.Parse()

	fmt.Print(xlog.Sprint("abc%v", 12))
	fmt.Printf("abc")

	whsvr := webhook.New(webhook.NewParameters(port, certfile, keyfile))

	whsvr.Server()
	xlog.Info("Server started")

	// listening OS shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	xlog.Info("Got OS shutdown signal, shutting down webhook server gracefully...")
	whsvr.Shutdown(context.Background())
}
