package nrf905

import (
	"fmt"
	"os"
	"time"

	"github.com/MarkSaravi/drone-go/types"
	"github.com/stianeikeland/go-rpio"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host/sysfs"
)

type nRF905 struct {
	txe  *rpio.Pin
	pwr  *rpio.Pin
	ce   *rpio.Pin
	cd   *rpio.Pin
	am   *rpio.Pin
	dr   *rpio.Pin
	conn spi.Conn
}

func CreateNRF905(config types.RadioLinkConfig) *nRF905 {
	txe := rpio.Pin(config.GPIO.TXEN)
	pwr := rpio.Pin(config.GPIO.PWR)
	ce := rpio.Pin(config.GPIO.CE)
	cd := rpio.Pin(config.GPIO.CD)
	am := rpio.Pin(config.GPIO.AM)
	dr := rpio.Pin(config.GPIO.DR)
	d, err := sysfs.NewSPI(config.BusNumber, config.ChipSelect)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// SPI1 only supports Mode0
	//to enable SPI1 in raspberry pi follow instructions here https://docs.rs/rppal/0.8.1/rppal/spi/index.html
	// or add "dtoverlay=spi1-3cs" to /boot/config.txt
	conn, err := d.Connect(physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	r := nRF905{
		txe:  &txe, // High to enable Tx
		pwr:  &pwr, // Power up
		ce:   &ce,  // Tx and Rx enable
		cd:   &cd,
		dr:   &dr,
		am:   &am,
		conn: conn,
	}
	r.initReceiver()
	return &r
}

func (rl *nRF905) standBy() {
	fmt.Println("Standby")
	rl.pwr.High()
	rl.txe.Low()
	rl.ce.Low()
}

func (rl *nRF905) powerUp() {
	fmt.Println("Power Up")
	rl.pwr.High()
	rl.txe.High()
	rl.ce.High()
}

func (rl *nRF905) initReceiver() {
	rl.ce.Output()
	rl.txe.Output()
	rl.pwr.Output()

	rl.cd.Input()
	rl.dr.Input()
	rl.am.Input()

	rl.dr.PullDown()
	rl.am.PullDown()
	rl.cd.PullDown()

	rl.standBy()

	r := make([]uint8, 11)
	const READ_RX_ADDRESS uint8 = 0b00010000
	const WRITE_RX_ADDRESS uint8 = 0b00000000

	w := []uint8{WRITE_RX_ADDRESS, 0x6C, 0xC, 0x44, 0x20, 0x20, 0x58, 0x6F, 0x2E, 0x10, 0xD8}
	fmt.Println("Writing: ", w)
	err := rl.conn.Tx(w, nil)
	time.Sleep(20 * time.Millisecond)

	w[0] = READ_RX_ADDRESS
	err = rl.conn.Tx(w, r)
	fmt.Println(r[1:])
	fmt.Printf("0x%X, 0x%X, 0x%X, 0x%X, 0x%X, 0x%X, 0x%X, 0x%X, 0x%X, 0x%X\n %v\n", r[1], r[2], r[3], r[4], r[5], r[6], r[7], r[8], r[9], r[10], err)
	time.Sleep(20 * time.Millisecond)
	rl.powerUp()
}

func (rl *nRF905) IsDataReady() bool {
	return rl.dr.Read() == rpio.High
}

func (rl *nRF905) Close() {
	rl.ce.Low()
	rl.pwr.Low()
}
