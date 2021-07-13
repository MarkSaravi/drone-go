package nrf905

import (
	"fmt"
	"os"

	"github.com/MarkSaravi/drone-go/connectors/gpio"
	"github.com/MarkSaravi/drone-go/types"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host/sysfs"
)

type nRF905 struct {
	txen *gpio.Pin
	pwr  *gpio.Pin
	ce   *gpio.Pin
	cd   *gpio.Pin
	am   *gpio.Pin
	dr   *gpio.Pin
	conn spi.Conn
}

func CreateNRF905(config types.RadioLinkConfig) *nRF905 {
	txen := createPin(gpio.GetGPIO(config.GPIO.TXEN))
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
		txen: txen,
		pwr:  pwr,
		ce:   ce,
		cd:   cd,
		dr:   dr,
		am:   am,
		conn: conn,
	}
	r.initReceiver()
	return &r
}

func createPin(gpioPinNum int) *gpio.Pin {
	pin, err := gpio.NewPin(gpioPinNum)
	if err != nil {
		panic(fmt.Sprintf("Can't create pin %d", gpioPinNum))
	}
	return pin
}

func (rl *nRF905) initReceiver() {
	rl.dr.SetAsInput()
	rl.am.SetAsInput()
	rl.cd.SetAsInput()
	rl.txen.SetAsOutput()
	rl.pwr.SetAsOutput()
	rl.ce.SetAsOutput()

	rl.txen.SetLow()
	rl.pwr.SetLow()
	rl.ce.SetLow()
	w := make([]uint8, 5)
	r := make([]uint8, 5)
	const READ_RX_ADDRESS uint8 = 0b00010101
	const WRITE_RX_ADDRESS uint8 = 0b00000101
	w[0] = READ_RX_ADDRESS
	err := rl.conn.Tx(w, r)
	fmt.Println(err, r)

	w = []uint8{WRITE_RX_ADDRESS, 0x58, 0x6F, 0x2E, 0x10}
	fmt.Println("Writing: ", w)
	err = rl.conn.Tx(w, r)
	fmt.Println(err, r)

	w[0] = READ_RX_ADDRESS
	err = rl.conn.Tx(w, r)
	fmt.Println(err, r)

	rl.pwr.SetHigh()
	rl.ce.SetHigh()
}

func (rl *nRF905) IsDataReady() bool {
	if rl.dr.GetLevel() == gpio.High {
		return true
	}
	return false
}

func (rl *nRF905) Close() {
	rl.ce.SetLow()
	rl.pwr.SetLow()
}
