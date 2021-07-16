package main

import (
	"fmt"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio"
)

func main() {
	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rpio.Close()

	// config := types.RadioLinkConfig{
	// 	GPIO: types.RadioLinkGPIOPins{
	// 		TXEN: 6,
	// 		CE:   26,
	// 		PWR:  5,
	// 		CD:   25,
	// 		AM:   23,
	// 		DR:   24,
	// 	},
	// 	BusNumber:  1,
	// 	ChipSelect: 2,
	// 	RxAddress:  "",
	// 	TxAddress:  "",
	// }

	inputTest()

}

func inputTest() {
	var counter int = 0
	var run bool = true
	go func() {
		var b []byte = make([]byte, 1)
		os.Stdin.Read(b)
		if b[0] == '\n' {
			run = false
			return
		}
	}()

	go func() {
		outpin := rpio.Pin(25)
		fmt.Println("Output started")
		outpin.Output()
		for run {
			outpin.High()
			time.Sleep(time.Second)
			outpin.Low()
			time.Sleep(time.Second)
		}
		fmt.Println("Output stopped")
	}()

	inpin := rpio.Pin(24)
	inpin.Input()
	inpin.PullDown()
	pState := rpio.Low
	for run {
		nState := inpin.Read()
		if nState != pState {
			pState = nState
			fmt.Println(nState, counter)
			counter++
		}
	}
}

func outputTest(pin *rpio.Pin) {
	pin.Output()
	start := time.Now()
	for time.Since(start) < time.Second*5 {
		pin.High()
		time.Sleep(time.Millisecond * 250)
		pin.Low()
		time.Sleep(time.Millisecond * 250)
	}
	pin.Input()
}
