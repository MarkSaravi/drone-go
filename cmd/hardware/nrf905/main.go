package main

import (
	"github.com/MarkSaravi/drone-go/connectors/gpio"
	"github.com/MarkSaravi/drone-go/hardware/nrf905"
	"github.com/MarkSaravi/drone-go/types"
)

func main() {
	gpio.Open()
	nrf905 := nrf905.CreateNRF905(types.RadioLinkConfig{
		GPIO: types.RadioLinkGPIOPins{
			TXEN: 6,
			CE:   26,
			PWR:  5,
			CD:   25,
			AM:   23,
			DR:   24,
		},
		BusNumber:  1,
		ChipSelect: 2,
		RxAddress:  "",
		TxAddress:  "",
	})
	nrf905.IsDataReady()
	nrf905.Close()
	gpio.Close()
}
