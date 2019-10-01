package config

type EnumLinkType int16

// Device types
const (
	LinkTypeNone EnumLinkType = iota
	LinkTypeSerial
	LinkTypeHidUsb
)

type SerialConfig struct {
	PortName string `yaml:"port_name"`
	BaudRate uint32 `yaml:"baud_rate"`
	DataBits uint16 `yaml:"data_bits"`
	StopBits uint16 `yaml:"stop_bits"`
	Parity   uint16 `yaml:"parity"`
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
		Serial:   &SerialConfig{},
		HidUsb:   &HidUsbConfig{},
	}
	return lnkCfg
}
