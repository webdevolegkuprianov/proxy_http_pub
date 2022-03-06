package main

import (
	logger "github.com/webdevolegkuprianov/proxy_http/app/logger"
	"github.com/webdevolegkuprianov/proxy_http/app/model"
	"github.com/webdevolegkuprianov/proxy_http/app/proxyserver"
)

func main() {

	config, err := model.NewConfig()
	if err != nil {
		logger.ErrorLogger.Println(err)

	}

	if err := proxyserver.Start(config); err != nil {
		logger.ErrorLogger.Println(err)
	}

}
