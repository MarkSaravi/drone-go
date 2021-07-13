|nRF905|Raspberry Pi|Arduino Uno|Description
| ----------- | ------------------- | ----------- | ----------- |
|VCC	      |3.3V                 | 3.3V	      |Power (3.3V)|
|TXEN	      |GPIO6                |9	          |TX or RX mode – High = TX, Low = RX|
|CE           |GPIO26               |7	          |Standby – High = TX/RX mode, Low = standby|
|PWR          |GPIO5                |8	          |Power-up – High = on, Low = off|
|CD	          |GPIO25               |4	          |Carrier detect – High when a signal is detected, for collision avoidance|
|AM           |GPIO23               |2            |Address Match – High when receiving a packet that has the same address as the one set for this device, optional since state is stored in register, if interrupts are used (default) then this pin must be connected|
|DR           |GPIO24               |3	          |Data Ready – High when finished transmitting/High when new data received, optional since state is stored in register, if interrupts are used (default) then this pin must be connected|
|MISO         |SPI1 MISO GPIO19 (35)|12	          |SPI MISO (Mega pin 50)|
|MOSI         |SPI1 MOSI GPIO20 (38)|11           |SPI MOSI (Mega pin 51)|
|SCK          |SPI1 SCLK GPIO21 (40)|13           |SPI SCK (Mega pin 52)|
|CSN          |SPI1 CE2 GPIO16 (36) |10           |SPI SS|
|GND          |GND                  |GND          |Ground|