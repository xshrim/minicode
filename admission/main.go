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

func mai() {
	var port int
	var certfile, keyfile string

	xlog.SetLevel(xlog.INFO)
	xlog.SetFlag(xlog.Ldate | xlog.Ltime | xlog.Lshortfile | xlog.Lcolor)

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

	x := xlog.New(os.Stdout, "YES", xlog.DEBUG, xlog.Ldate|xlog.Ltime|xlog.Lshortfile|xlog.Lcolor)

	x.Info("AMD %s", "yes")
	x.Trace("trace")

	// fmt.Println("aaa")
	xlog.Info("AMD %123", "YES")
	xlog.Error("Error")
	xlog.Warn("Warn")
	xlog.Info("Info")
	xlog.Debug("Debug")
	xlog.Trace("Trace")
	//	xlog.Fatal("Fatal")
	xlog.Panic("Panic")

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
