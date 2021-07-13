package types

type RadioLinkGPIOPins struct {
	TXEN int `yaml:"txen"`
	CE   int `yaml:"ce"`
	PWR  int `yaml:"pwr"`
	CD   int `yaml:"cd"`
	AM   int `yaml:"am"`
	DR   int `yaml:"dr"`
}

type RadioLinkConfig struct {
	GPIO       RadioLinkGPIOPins `yaml:"gpio"`
	RxAddress  string            `yaml:"rx_address"`
	TxAddress  string            `yaml:"tx_address"`
	BusNumber  int               `yaml:"bus_number"`
	ChipSelect int               `yaml:"chip_select"`
}
