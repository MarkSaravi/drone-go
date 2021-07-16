package main

import (
	"fmt"
	"os"
	"time"

	"github.com/MarkSaravi/drone-go/hardware/nrf905"
	"github.com/MarkSaravi/drone-go/types"
	"github.com/stianeikeland/go-rpio"
)

func main() {
	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rpio.Close()

	var run bool = true
	go func() {
		var b []byte = make([]byte, 1)
		os.Stdin.Read(b)
		if b[0] == '\n' {
			run = false
			return
		}
	}()

	config := types.RadioLinkConfig{
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
	}

	nrf905 := nrf905.CreateNRF905(config)
	var isready bool = true
	for run {
		drstate := nrf905.IsDataReady()
		if isready != drstate {
			isready = drstate
			if isready {
				fmt.Println("Data Ready")
			} else {
				fmt.Println("Waiting")
			}
		}
		time.Sleep(time.Second / 10)
	}
}
