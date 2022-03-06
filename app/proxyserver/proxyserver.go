package proxyserver

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"path/filepath"

	logger "github.com/webdevolegkuprianov/proxy_http/app/logger"
	"github.com/webdevolegkuprianov/proxy_http/app/model"
)

func Start(config *model.Config) error {

	//cert, key files
	fcert, err := filepath.Abs("/root/cert/onsales.st.tech.crt")
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	fkey, err := filepath.Abs("/root/cert/onsales.st.tech.key")
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	//cert, key load
	cer, err := tls.LoadX509KeyPair(fcert, fkey)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	configCert := &tls.Config{Certificates: []tls.Certificate{cer}}

	caCert, err := ioutil.ReadFile(fcert)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	server := newServer(config)

	//setup HTTPS server
	srv := &http.Server{
		Addr:      config.Spec.Ports.Addr,
		TLSConfig: configCert,
		Handler:   server.router,
	}

	return srv.ListenAndServeTLS(fcert, fkey)

}
