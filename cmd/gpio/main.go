package main

import (
	"log"
	"time"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	testPinOut("GPIO24")
}

// func inputTest() {
// 	var counter int = 0
// 	var run bool = true
// 	go func() {
// 		var b []byte = make([]byte, 1)
// 		os.Stdin.Read(b)
// 		if b[0] == '\n' {
// 			run = false
// 			return
// 		}
// 	}()

// 	go func() {
// 		outpin := rpio.Pin(25)
// 		fmt.Println("Output started")
// 		outpin.Output()
// 		for run {
// 			outpin.High()
// 			time.Sleep(time.Second)
// 			outpin.Low()
// 			time.Sleep(time.Second)
// 		}
// 		fmt.Println("Output stopped")
// 	}()

// 	inpin := rpio.Pin(24)
// 	inpin.Input()
// 	inpin.PullDown()
// 	pState := rpio.Low
// 	for run {
// 		nState := inpin.Read()
// 		if nState != pState {
// 			pState = nState
// 			fmt.Println(nState, counter)
// 			counter++
// 		}
// 	}
// }

func testPinOut(pin string) {
	// Use gpioreg GPIO pin registry to find a GPIO pin by name.
	p := gpioreg.ByName(pin)
	if p == nil {
		log.Fatal("Failed to find GPIO6")
	}
	if err := p.Out(gpio.High); err != nil {
		log.Fatal(err)
	}
	start := time.Now()
	for time.Since(start) < time.Second*5 {
		setPinOutLevel(p, gpio.High)
		time.Sleep(time.Second / 2)
		setPinOutLevel(p, gpio.Low)
		time.Sleep(time.Second / 2)
	}
}

func setPinOutLevel(p gpio.PinOut, level gpio.Level) {
	if err := p.Out(level); err != nil {
		log.Fatal(err)
	}
}
