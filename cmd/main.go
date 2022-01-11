package main

import (
	"ctmgr/containermgr"
	"ctmgr/server"
	"os"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func initLog(logPath string) error {
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(file)
	log.SetLevel(logrus.DebugLevel)
	log.SetReportCaller(true)
	return nil
}
func main() {
	if err := initLog("/home/ubuntu/log/container_mgr.log"); err != nil {
		return
	}
	if err := containermgr.CreateInstance(); err != nil {
		return
	}
	server.Init()
	grpcSvr, err := server.NewServer(":5555")
	if err != nil {
		return
	}
	grpcSvr.Run()
}
