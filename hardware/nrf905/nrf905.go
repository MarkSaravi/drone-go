package nrf905

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/connectors/gpio"
	"github.com/MarkSaravi/drone-go/types"
)

type nRF905 struct {
	txen *gpio.Pin
	pwr  *gpio.Pin
	ce   *gpio.Pin
	cd   *gpio.Pin
	am   *gpio.Pin
	dr   *gpio.Pin
	miso *gpio.Pin
	mosi *gpio.Pin
	sck  *gpio.Pin
	csn  *gpio.Pin
}

func CreateNRF905(config types.RadioLinkConfig) *nRF905 {
	txen := createPin(gpio.GetGPIO(config.TXENGpioPin))
	pwr := createPin(gpio.GetGPIO(config.PWRGpioPin))
	ce := createPin(gpio.GetGPIO(config.CEGpioPin))
	cd := createPin(gpio.GetGPIO(config.CDGpioPin))
	am := createPin(gpio.GetGPIO(config.AMGpioPin))
	dr := createPin(gpio.GetGPIO(config.DRGpioPin))
	miso := createPin(gpio.GetGPIO(config.MISOGpioPin))
	mosi := createPin(gpio.GetGPIO(config.MOSIGpioPin))
	sck := createPin(gpio.GetGPIO(config.SCKGpioPin))
	csn := createPin(gpio.GetGPIO(config.CSNGpioPin))
	r := nRF905{
		txen: txen,
		pwr:  pwr,
		ce:   ce,
		cd:   cd,
		dr:   dr,
		am:   am,
		miso: miso,
		mosi: mosi,
		sck:  sck,
		csn:  csn,
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
	// set as receiver
	r.txen.SetAsOutput()
	r.txen.SetLow()
	// enable receiver
	r.ce.SetAsOutput()
	r.ce.SetHigh()
	r.pwr.SetAsOutput()
	r.pwr.SetHigh()
	r.dr.SetAsInput()
	r.am.SetAsInput()
	r.cd.SetAsInput()
}

func (r *nRF905) IsDataReady() bool {
	return false
}

func (r *nRF905) Close() {
}
