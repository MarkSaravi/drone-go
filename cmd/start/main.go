package main

import (
	"log"

	"github.com/MarkSaravi/drone-go/flightcontrol"
	"github.com/MarkSaravi/drone-go/modules/esc"
	"github.com/MarkSaravi/drone-go/modules/imu"
	"github.com/MarkSaravi/drone-go/modules/udplogger"
	"github.com/MarkSaravi/drone-go/utils"
	"periph.io/x/periph/host"
)

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	appConfig := utils.ReadConfigs()
	udpLogger := udplogger.CreateUdpLogger(appConfig.UDP, appConfig.Flight.Imu.ImuDataPerSecond)
	imu := imu.CreateIM(appConfig)
	pid := flightcontrol.CreatePidController()
	esc := esc.NewESCsHandler(appConfig.Flight.Esc)
	flightControl := flightcontrol.CreateFlightControl(imu, pid, esc, udpLogger)

	flightControl.Start()
}
