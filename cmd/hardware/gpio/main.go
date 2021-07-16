package main

import (
	"fmt"
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
	const pinOutName = "GPIO6"
	const pinInName = "GPIO24"
	testPinOut(setupPin(pinOutName))
	testPinIn(setupPin(pinInName), setupPin(pinOutName))
}

func setupPin(pinName string) gpio.PinIO {
	pin := gpioreg.ByName(pinName)
	if pin == nil {
		log.Fatal("Failed to find ", pinName)
	}
	return pin
}

func testPinOut(p gpio.PinOut) {
	// Use gpioreg GPIO pin registry to find a GPIO pin by name.
	for i := 0; i < 5; i++ {
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

func testPinIn(pin gpio.PinIn, pout gpio.PinOut) {
	end := make(chan (int), 1)
	go func() {
		fmt.Println("Setting level for ", pin)
		time.Sleep(time.Millisecond * 10)
		testPinOut(pout)
		end <- 1
	}()
	// Set it as input, with a pull down (defaults to Low when unconnected) and
	// enable rising edge triggering.
	if err := pin.In(gpio.PullDown, gpio.RisingEdge); err != nil {
		log.Fatal(err)
	}

	var run bool = true
	for run {
		select {
		case <-end:
			run = false
		default:
			if wh := pin.WaitForEdge(time.Millisecond); wh {
				fmt.Printf("%s went %s\n", pin, gpio.High)
			}
		}
	}
}
