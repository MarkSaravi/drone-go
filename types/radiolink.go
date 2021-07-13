package types

type RadioLinkConfig struct {
	TXENGpioPin        int
	CEGpioPin          int
	PWRGpioPin         int
	CDGpioPin          int
	AMGpioPin          int
	DRGpioPin          int
	MISOGpioPin        int
	MOSIGpioPin        int
	SCKGpioPin         int
	CSNGpioPin         int
	ReceiverAddress    string
	TransmitterAddress string
}
