package types

type RadioLinkGPIOPins struct {
	TXEN uint8 `yaml:"txen"`
	CE   uint8 `yaml:"ce"`
	PWR  uint8 `yaml:"pwr"`
	CD   uint8 `yaml:"cd"`
	AM   uint8 `yaml:"am"`
	DR   uint8 `yaml:"dr"`
}

type RadioLinkConfig struct {
	GPIO       RadioLinkGPIOPins `yaml:"gpio"`
	RxAddress  string            `yaml:"rx_address"`
	TxAddress  string            `yaml:"tx_address"`
	BusNumber  int               `yaml:"bus_number"`
	ChipSelect int               `yaml:"chip_select"`
}
