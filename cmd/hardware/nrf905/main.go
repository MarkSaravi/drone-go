package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MarkSaravi/drone-go/hardware/nrf905"
	"github.com/MarkSaravi/drone-go/types"
	"periph.io/x/periph/host"
)

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	config := types.RadioLinkConfig{
		GPIO: types.RadioLinkGPIOPins{
			TXEN: "GPIO6",
			CE:   "GPIO26",
			PWR:  "GPIO5",
			CD:   "GPIO25",
			AM:   "GPIO23",
			DR:   "GPIO24",
		},
		BusNumber:  1,
		ChipSelect: 2,
		RxAddress:  "",
		TxAddress:  "",
	}

	nrf905 := nrf905.CreateNRF905(config)
	nrf905.ReadData()
	endChannel := createEndChannel()
	end := false
	nrf905.PowerUp()
	fmt.Println(nrf905.ReadData())
	for !end {
		select {
		case end = <-endChannel:
		default:
			if nrf905.IsDataReady() {
				fmt.Println(string(nrf905.ReadData()))
			}
		}
	}
}

func createEndChannel() chan (bool) {
	end := make(chan (bool), 1)
	go func() {
		var b []byte = make([]byte, 1)
		os.Stdin.Read(b)
		if b[0] == '\n' {
			end <- true
			return
		}
	}()
	return end
}
