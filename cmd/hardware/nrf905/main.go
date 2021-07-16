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
	endChannel := createEndChannel()
	end := false
	for !end {
		isready := nrf905.IsDataReady()
		if isready {
			fmt.Println("Data Ready")
		}
		select {
		case end = <-endChannel:
		default:
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
