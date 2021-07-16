package nrf905

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MarkSaravi/drone-go/types"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/host/sysfs"
)

const (
	READ_CONFIG      byte = 0b00010000
	WRITE_CONFIG     byte = 0b00000000
	READ_RX_PAYLOAD  byte = 0b00100100
	RX_PAYLOAD_WIDTH byte = 32
	TX_PAYLOAD_WIDTH byte = 32
)

type nRF905 struct {
	txe  gpio.PinOut
	pwr  gpio.PinOut
	ce   gpio.PinOut
	cd   gpio.PinIn
	am   gpio.PinIn
	dr   gpio.PinIn
	conn spi.Conn
}

func CreateNRF905(config types.RadioLinkConfig) *nRF905 {
	txe := initPin(config.GPIO.TXEN)
	pwr := initPin(config.GPIO.PWR)
	ce := initPin(config.GPIO.CE)
	cd := initPin(config.GPIO.CD)
	am := initPin(config.GPIO.AM)
	dr := initPin(config.GPIO.DR)
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
	r.initReceiver(config.RxAddress)
	return &r
}

func initPin(pinName string) gpio.PinIO {
	pin := gpioreg.ByName(pinName)
	if pin == nil {
		log.Fatal("Failed to find ", pinName)
	}
	return pin
}

func (rl *nRF905) standBy() {
	fmt.Println("Standby")
	rl.pwr.Out(gpio.High)
	rl.ce.Out(gpio.Low)
	time.Sleep(10 * time.Millisecond)
}

func (rl *nRF905) PowerUp() {
	fmt.Println("Power Up")
	rl.ce.Out(gpio.High)
	time.Sleep(10 * time.Millisecond)
	rl.pwr.Out(gpio.High)
	time.Sleep(10 * time.Millisecond)
}

func (rl *nRF905) initReceiver(address string) {
	rl.cd.In(gpio.PullDown, gpio.RisingEdge)
	rl.dr.In(gpio.PullDown, gpio.RisingEdge)
	rl.am.In(gpio.PullDown, gpio.RisingEdge)

	rl.txe.Out(gpio.Low)

	rl.standBy()

	rxadd := parseAddress(address)
	w := []byte{WRITE_CONFIG, 0x6C, 0xC, 0x44, RX_PAYLOAD_WIDTH, TX_PAYLOAD_WIDTH, rxadd[3], rxadd[2], rxadd[1], rxadd[0], 0xD8}
	fmt.Println("Writing: ", w)
	rl.conn.Tx(w, nil)
	time.Sleep(20 * time.Millisecond)

	r := make([]byte, 11)
	w[0] = READ_CONFIG
	rl.conn.Tx(w, r)
	fmt.Println("Config: ", r[1:])
	time.Sleep(20 * time.Millisecond)
}

func (rl *nRF905) IsDataReady() bool {
	return rl.dr.Read() == gpio.High
}

func (rl *nRF905) ReadData() []byte {
	w := make([]byte, RX_PAYLOAD_WIDTH+1)
	r := make([]byte, RX_PAYLOAD_WIDTH+1)
	w[0] = READ_RX_PAYLOAD
	rl.conn.Tx(w, r)
	return r[1:]
}
func (rl *nRF905) ReadConfig() []byte {
	r := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	w := make([]byte, 11)
	w[0] = READ_CONFIG
	rl.conn.Tx(w, r)
	return r[1:]
}

func (rl *nRF905) Close() {
	rl.ce.Out(gpio.Low)
	rl.pwr.Out(gpio.Low)
}

func parseAddress(address string) []byte {
	s := strings.Split(address, ":")
	b := make([]byte, 4)
	for i := 0; i < 4; i++ {
		n, _ := strconv.ParseInt(s[i], 16, 16)
		b[i] = byte(n)
	}
	return b
}
