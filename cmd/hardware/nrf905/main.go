package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"unsafe"

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
		RxAddress:  "39:B5:3C:90",
		TxAddress:  "",
	}

	nrf905 := nrf905.CreateNRF905(config)
	nrf905.ReadData()
	endChannel := createEndChannel()
	end := false
	nrf905.PowerUp()
	fmt.Println(nrf905.ReadData())
	for !end {
		time.Sleep(time.Millisecond * 5)
		select {
		case end = <-endChannel:
		default:
			if nrf905.IsDataReady() {
				byteData := nrf905.ReadData()
				intData := byteArrayToInt16Array(byteData, 32)
				fmt.Println(intData)
			}
		}
	}
}

func byteArrayToInt16Array(ba []byte, size int) []int16 {
	type pInt16Array = *([]int16)
	var pi16 pInt16Array = pInt16Array(unsafe.Pointer(&ba))
	var ia []int16 = make([]int16, size/2)
	for i := 0; i < size/2; i++ {
		ia[i] = (*pi16)[i]
	}
	return ia
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

// func Int16ToUint8(n int16) []uint8 {
// 	var ui uint16 = uint16(n)
// 	return []uint8{uint8(ui & 0b0000000011111111), uint8(ui >> 8)}
// }

// func Uint8ToInt16(ui8 []uint8) int16 {
// 	var ui16 uint16 = uint16(ui8[1])<<8 | uint16(ui8[0])
// 	return int16(ui16)
// }
