package main

import (
	"github.com/MarkSaravi/drone-go/hardware/nrf905"
	"github.com/MarkSaravi/drone-go/types"
)

func main() {
	nrf905 := nrf905.CreateNRF905(types.RadioLinkConfig{
		TXENGpioPin: 6,
		CEGpioPin:   26,
		PWRGpioPin:  5,
		CDGpioPin:   25,
		AMGpioPin:   23,
		DRGpioPin:   24,
		MISOGpioPin: 19,
		MOSIGpioPin: 20,
		SCKGpioPin:  21,
		CSNGpioPin:  16,
	})
	nrf905.IsDataReady()
	nrf905.Close()
}
