package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/yisaer/sidecar-inject-server/pkg/config"
	"github.com/yisaer/sidecar-inject-server/pkg/webhook"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var parameters webhook.WhSvrParameters

	// get command line parameters
	flag.IntVar(&parameters.Port, "port", 443, "Webhook server port.")
	flag.StringVar(&parameters.CertFile, "tlsCertFile", "/etc/webhook/certs/cert.pem", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&parameters.KeyFile, "tlsKeyFile", "/etc/webhook/certs/key.pem", "File containing the x509 private key to --tlsCertFile.")
	flag.StringVar(&parameters.Token, "authToken", "/var/run/secrets/kubernetes.io/serviceaccount/token", "Token to communicate with apiServer")
	flag.StringVar(&parameters.Crt, "ca", "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt", "crt to verify api server")
	flag.Parse()

	url := "https://kubernetes.default.svc.cluster.local/apis/yisaer.github.io/v1alpha1/sidecars"

	// load apiServer ca.crt
	certPool, err := config.LoadCA(parameters.Crt)
	if err != nil {
		glog.Errorf("Filed to load crt: %v", err)
	}

	// load apiServer token
	token, err := config.LoadToken(parameters.Token)
	if err != nil {
		glog.Errorf("failed to load token")
	}

	// load tls key pair
	pair, err := tls.LoadX509KeyPair(parameters.CertFile, parameters.KeyFile)
	if err != nil {
		glog.Errorf("Filed to load key pair: %v", err)
	}

	// build server
	whsvr := &webhook.WebhookServer{
		Client: &config.WebClient{
			Url:   url,
			Token: token,
			Client: &http.Client{Transport: &http.Transport{
				TLSClientConfig: &tls.Config{RootCAs: certPool},
			}},
		},
		Server: &http.Server{
			Addr:      fmt.Sprintf(":%v", parameters.Port),
			TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}},
		},
	}

	// define http server and server handler
	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", whsvr.Serve)
	whsvr.Server.Handler = mux

	// start webhook server in new rountine
	go func() {
		if err := whsvr.Server.ListenAndServeTLS("", ""); err != nil {
			glog.Errorf("Filed to listen and serve webhook server: %v", err)
		}
	}()

	// listening OS shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	glog.Infof("Got OS shutdown signal, shutting down wenhook server gracefully...")
	whsvr.Server.Shutdown(context.Background())
}
