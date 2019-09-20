package config

type EnumLinkType int16

// Device types
const (
	LinkTypeNone EnumLinkType = iota
	LinkTypeSerial
	LinkTypeHidUsb
)

type SerialConfig struct {
	PortNumber  int32 `yaml:"portNumber"`
	BaudRate    int32 `yaml:"baudRate"`
	DataBits    int32 `yaml:"dataBits"`
	StopBits    int32 `yaml:"stopBits"`
	Parity      int32 `yaml:"parity"`
	FlowControl int32 `yaml:"flowControl"`
	SendTimeout int32 `yaml:"sendTimeout"`
	RecvTimeout int32 `yaml:"recvTimeout"`
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
