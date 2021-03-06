/*
TO increase udp packet size in macOS use the following command
sudo sysctl -w net.inet.udp.maxdgram=65535
*/
package udplogger

import (
	"fmt"
	"math"
	"net"
	"strings"

	"github.com/MarkSaravi/drone-go/types"
)

type udpLogger struct {
	conn                 *net.UDPConn
	address              *net.UDPAddr
	enabled              bool
	buffer               []string
	imuDataPerSecond     int
	dataPerPacket        int
	dataPerPacketCounter int
	skipOffset           int
	maxDataPerPacket     int
	bufferCounter        int
}

func CreateUdpLogger(udpConfig types.UdpLoggerConfig, imuDataPerSecond int) types.UdpLogger {
	if !udpConfig.Enabled {
		return &udpLogger{
			enabled: false,
		}
	}
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		fmt.Println("UDP initialization error: ", err)
		return &udpLogger{
			enabled: false,
		}
	}
	address, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", udpConfig.IP, udpConfig.Port))
	if err != nil {
		fmt.Println("UDP initialization error: ", err)
		return &udpLogger{
			enabled: false,
		}
	}

	if udpConfig.PacketsPerSecond <= 0 {
		return &udpLogger{
			enabled: false,
		}
	}
	dataPerPacket := imuDataPerSecond / udpConfig.PacketsPerSecond
	var skipOffset int = 1
	var maxDataPerPacket = dataPerPacket
	if dataPerPacket > udpConfig.MaxDataPerPacket {
		maxDataPerPacket = udpConfig.MaxDataPerPacket
		skipOffset = int(math.Ceil(float64(dataPerPacket) / float64(udpConfig.MaxDataPerPacket)))
	}
	fmt.Println("DPP: ", imuDataPerSecond, dataPerPacket, maxDataPerPacket, skipOffset)

	return &udpLogger{
		conn:                 conn,
		address:              address,
		enabled:              true,
		imuDataPerSecond:     imuDataPerSecond,
		dataPerPacket:        dataPerPacket,
		dataPerPacketCounter: 0,
		bufferCounter:        0,
		skipOffset:           skipOffset,
		maxDataPerPacket:     maxDataPerPacket,
		buffer:               make([]string, maxDataPerPacket),
	}
}

func (l *udpLogger) appendData(imuRotations types.ImuRotations) {
	if !l.enabled {
		return
	}
	l.dataPerPacketCounter++
	if l.dataPerPacketCounter%l.skipOffset == 0 {
		l.buffer[l.bufferCounter] = imuDataToJson(imuRotations)
		l.bufferCounter++
	}
}

func (l *udpLogger) Send(imuRotations types.ImuRotations) {
	if !l.enabled {
		return
	}
	l.appendData(imuRotations)
	l.sendData()
}

func (l *udpLogger) sendData() {
	if l.enabled && l.dataPerPacketCounter == l.dataPerPacket {
		jsonPayload := fmt.Sprintf("{\"d\":[%s],\"dps\":%d}",
			strings.Join(l.buffer[0:l.bufferCounter], ","),
			l.imuDataPerSecond,
		)
		l.dataPerPacketCounter = 0
		l.bufferCounter = 0
		// data len should be less than sysctl net.inet.udp.maxdgram for macOS
		go func() {
			l.conn.WriteToUDP([]byte(jsonPayload), l.address)
		}()
	}
}

func imuDataToJson(imuRotations types.ImuRotations) string {
	return fmt.Sprintf(`{"a":{"r":%0.2f,"p":%0.2f,"y":%0.2f},"g":{"r":%0.2f,"p":%0.2f,"y":%0.2f},"r":{"r":%0.2f,"p":%0.2f,"y":%0.2f},"t":%d,"dt":%d}`,
		imuRotations.Accelerometer.Roll,
		imuRotations.Accelerometer.Pitch,
		imuRotations.Accelerometer.Yaw,
		imuRotations.Gyroscope.Roll,
		imuRotations.Gyroscope.Pitch,
		imuRotations.Gyroscope.Yaw,
		imuRotations.Rotations.Roll,
		imuRotations.Rotations.Pitch,
		imuRotations.Rotations.Yaw,
		imuRotations.ReadTime,
		imuRotations.ReadInterval,
	)
}
