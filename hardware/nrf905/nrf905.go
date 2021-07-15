package nrf905

import (
	"fmt"
	"os"
	"time"

	"github.com/MarkSaravi/drone-go/connectors/gpio"
	"github.com/MarkSaravi/drone-go/types"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host/sysfs"
)

type nRF905 struct {
	txe  *gpio.Pin
	pwr  *gpio.Pin
	ce   *gpio.Pin
	cd   *gpio.Pin
	am   *gpio.Pin
	dr   *gpio.Pin
	conn spi.Conn
}

func CreateNRF905(config types.RadioLinkConfig) *nRF905 {
	txe := createPin(gpio.GetGPIO(config.GPIO.TXEN))
	pwr := createPin(gpio.GetGPIO(config.GPIO.PWR))
	ce := createPin(gpio.GetGPIO(config.GPIO.CE))
	cd := createPin(gpio.GetGPIO(config.GPIO.CD))
	am := createPin(gpio.GetGPIO(config.GPIO.AM))
	dr := createPin(gpio.GetGPIO(config.GPIO.DR))
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
		txe:  txe, // High to enable Tx
		pwr:  pwr, // Power up
		ce:   ce,  // Tx and Rx enable
		cd:   cd,
		dr:   dr,
		am:   am,
		conn: conn,
	}
	r.initReceiver()
	return &r
}

func createPin(gpioPinNum uint8) *gpio.Pin {
	pin, err := gpio.NewPin(gpioPinNum)
	if err != nil {
		panic(fmt.Sprintf("Can't create pin %d", gpioPinNum))
	}
	return pin
}

func (rl *nRF905) standBy() {
	fmt.Println("Standby")
	rl.pwr.SetHigh()
	rl.txe.SetLow()
	rl.ce.SetLow()
}

func (rl *nRF905) powerUp() {
	fmt.Println("Power Up")
	rl.pwr.SetHigh()
	rl.txe.SetHigh()
	rl.ce.SetHigh()
}

func (rl *nRF905) initReceiver() {
	rl.ce.SetAsOutput()
	rl.txe.SetAsOutput()
	rl.pwr.SetAsOutput()

	rl.cd.SetAsInput()
	rl.dr.SetAsInput()
	rl.am.SetAsInput()

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
	w := []uint8{0xFF}
	r := []uint8{0}
	err := rl.conn.Tx(w, r)

	if err != nil {
		fmt.Println(err)
	}
	var dataReady uint8 = r[0] & 0b00100000
	return dataReady != 0
}

func (rl *nRF905) Close() {
	rl.ce.SetLow()
	rl.pwr.SetLow()
}
