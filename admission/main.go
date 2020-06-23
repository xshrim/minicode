package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"label.admission/webhook"
	"label.admission/xlog"
)

func main() {
	var port int
	var certfile, keyfile string

	xlog.Level = xlog.TRACE
	xlog.Color = true

	// get command line parameters
	// kubectl create secret tls tls-secret --cert=path/to/tls.cert --key=path/to/tls.key
	// volumeMounts:
	// - mountPath: /etc/webhook/certs
	//   name: tls-secret
	//   readOnly: true
	// volumes:
	// - name: tls-secret
	//   secret:
	//     defaultMode: 420
	//     secretName: tls-secret
	flag.IntVar(&port, "port", 443, "Webhook server port.")
	flag.StringVar(&certfile, "tlsCertFile", "/etc/webhook/certs/tls.crt", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&keyfile, "tlsKeyFile", "/etc/webhook/certs/tls.key", "File containing the x509 private key to --tlsCertFile.")
	flag.Parse()

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
