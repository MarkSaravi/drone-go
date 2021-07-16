package types

type RadioLinkGPIOPins struct {
	TXEN string `yaml:"txen"`
	CE   string `yaml:"ce"`
	PWR  string `yaml:"pwr"`
	CD   string `yaml:"cd"`
	AM   string `yaml:"am"`
	DR   string `yaml:"dr"`
}

type RadioLinkConfig struct {
	GPIO       RadioLinkGPIOPins `yaml:"gpio"`
	RxAddress  string            `yaml:"rx_address"`
	TxAddress  string            `yaml:"tx_address"`
	BusNumber  int               `yaml:"bus_number"`
	ChipSelect int               `yaml:"chip_select"`
}
