package config

type EnumLinkType uint16
type StopBits uint16
type Parity uint16

// Device types
const (
	LinkTypeNone EnumLinkType = iota
	LinkTypeSerial
	LinkTypeHidUsb
)

// Stop bits types
const (
	OneStopBit StopBits = iota
	OneHalfStopBits
	TwoStopBits
)

// Port parity types
const (
	NoParity Parity = iota
	OddParity
	EvenParity
	MarkParity
	SpaceParity
)

type SerialConfig struct {
	PortName string   `yaml:"port_name"`
	BaudRate uint32   `yaml:"baud_rate"`
	DataBits uint16   `yaml:"data_bits"`
	StopBits StopBits `yaml:"stop_bits"`
	Parity   Parity   `yaml:"parity"`
}

type HidUsbConfig struct {
	VendorID  uint16 `yaml:"vendor_id"`  // Device Vendor ID
	ProductID uint16 `yaml:"product_id"` // Device Product ID
	Serial    string `yaml:"serial"`     // Serial Number
}

type LinkerConfig struct {
	LinkType EnumLinkType  `yaml:"link_type"`
	Timeout  uint16        `yaml:"timeout"`
	Serial   *SerialConfig `yaml:"serial"`
	HidUsb   *HidUsbConfig `yaml:"hidusb"`
}

func GetDefaultLinkerConfig() *LinkerConfig {
	lnkCfg := &LinkerConfig{
		LinkType: LinkTypeNone,
		Timeout:  0,
		Serial: &SerialConfig{
			PortName: "",
			BaudRate: 9600,
			DataBits: 8,
			StopBits: OneStopBit,
			Parity:   OddParity,
		},
		HidUsb: &HidUsbConfig{},
	}
	return lnkCfg
}
