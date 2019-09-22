package config

type EnumLinkType int16

// Device types
const (
	LinkTypeNone EnumLinkType = iota
	LinkTypeSerial
	LinkTypeHidUsb
)

type SerialConfig struct {
	PortName    string `yaml:"portName"`
	BaudRate    uint32 `yaml:"baudRate"`
	DataBits    uint16 `yaml:"dataBits"`
	StopBits    uint16 `yaml:"stopBits"`
	Parity      uint16 `yaml:"parity"`
	FlowControl uint16 `yaml:"flowControl"`
	SendTimeout uint16 `yaml:"sendTimeout"`
	RecvTimeout uint16 `yaml:"recvTimeout"`
}

type HidUsbConfig struct {
	VendorID  uint16 `yaml:"vendorId"`  // Device Vendor ID
	ProductID uint16 `yaml:"productId"` // Device Product ID
	Serial    string `yaml:"serial"`    // Serial Number
}

type LinkerConfig struct {
	LinkType EnumLinkType  `yaml:"linkType"`
	Serial   *SerialConfig `yaml:"serial"`
	HidUsb   *HidUsbConfig `yaml:"hidUsb"`
}

func GetDefaultLinkerConfig() *LinkerConfig {
	lnkCfg := &LinkerConfig{
		LinkType: LinkTypeNone,
		Serial:   &SerialConfig{},
		HidUsb:   &HidUsbConfig{},
	}
	return lnkCfg
}
