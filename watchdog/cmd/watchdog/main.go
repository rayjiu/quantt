package main

import (
	"os"

	"github.com/rayjiu/quantt/watchdog/internal/chinaA/wogoo"
	"github.com/rayjiu/quantt/watchdog/internal/chinaA/wogoo/constants"
	"github.com/rayjiu/quantt/watchdog/internal/chinaA/wogoo/model"
	"github.com/rayjiu/quantt/watchdog/internal/logging"
	log "github.com/sirupsen/logrus"
)

func main() {
	logging.InitLogger()
	// test
	var tcpHost = os.Getenv("TCP_SRV_HOST")
	var tcpUId = os.Getenv("TCP_SRV_U_ID")
	var tcpDeviceId = os.Getenv("TCP_SRV_U_DEVICE_ID")
	var tcpToken = os.Getenv("TCP_SRV_TOKEN")
	wogoo.WogooHQServer.SetupClient(tcpHost, tcpUId, tcpDeviceId, tcpToken)
	wogoo.WogooHQServer.DoSub("601318", constants.MarketTypeSH, constants.PushBizTypeSnapshot, func(ss *model.StockSnapshot) {
		log.Infof("Reveived:%v \n", ss)
	})
	log.Info("watch dog started.")
	for {
	}
}
