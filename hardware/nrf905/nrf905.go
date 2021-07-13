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
	conn *spi.Conn
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
		conn: &conn,
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

func (r *nRF905) initReceiver() {
	r.txen.SetAsOutput()
	r.txen.SetLow()
	r.pwr.SetAsOutput()
	r.pwr.SetHigh()
	r.ce.SetAsOutput()
	r.ce.SetHigh()

	r.dr.SetAsInput()
	r.am.SetAsInput()
	r.cd.SetAsInput()
}

func (r *nRF905) IsDataReady() bool {
	return false
}

func (r *nRF905) Close() {
}
