package config

import "fmt"

type EnumLinkType uint16
type EnumStopBits uint16
type EnumParity uint16

// Device types
const (
	LinkTypeNone EnumLinkType = iota
	LinkTypeSerial
	LinkTypeHidUsb
)

// Stop bits types
const (
	OneStopBit EnumStopBits = iota
	OneHalfStopBits
	TwoStopBits
)

// Port parity types
const (
	NoParity EnumParity = iota
	OddParity
	EvenParity
	MarkParity
	SpaceParity
)

func (e EnumLinkType) String() string {
	switch e {
	case LinkTypeNone:			return "Off-line"
	case LinkTypeSerial:		return "Serial"
	case LinkTypeHidUsb:		return "HID/USB"
	default:					return "Undefined"
	}
}

func (e EnumStopBits) String() string {
	switch e {
	case OneStopBit:			return "One stop bit"
	case OneHalfStopBits:		return "One and half"
	case TwoStopBits:			return "Two stop bits"
	default:					return "Undefined"
	}
}

func (e EnumParity) String() string {
	switch e {
	case NoParity:				return "No parity"
	case OddParity:				return "Odd parity"
	case EvenParity:			return "Even parity"
	case MarkParity:			return "Mark parity"
	case SpaceParity:			return "Space parity"
	default:					return "Undefined"
	}
}


type SerialConfig struct {
	PortName string       `yaml:"port_name"`
	BaudRate uint32       `yaml:"baud_rate"`
	DataBits uint16       `yaml:"data_bits"`
	StopBits EnumStopBits `yaml:"stop_bits"`
	Parity   EnumParity   `yaml:"parity"`
}

func (cfg *SerialConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\tSerial config: " +
		"PortName = %s, BaudRate = %d, DataBits = %d, StopBits = %s, Parity = %s.",
		cfg.PortName, cfg.BaudRate, cfg.DataBits, cfg.StopBits.String(), cfg.Parity.String())
	return str
}

type HidUsbConfig struct {
	VendorID  uint16 `yaml:"vendor_id"`  // Device Vendor ID
	ProductID uint16 `yaml:"product_id"` // Device Product ID
	Serial    string `yaml:"serial"`     // Serial Number
}

func (cfg *HidUsbConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\tHID/USB config: " +
		"VendorID = %X, ProductID = %X, Serial = %s.",
		cfg.VendorID, cfg.ProductID, cfg.Serial)
	return str
}

type LinkerConfig struct {
	LinkType EnumLinkType  `yaml:"link_type"`
	Timeout  uint16        `yaml:"timeout"`
	Serial   *SerialConfig `yaml:"serial"`
	HidUsb   *HidUsbConfig `yaml:"hid_usb"`
}

func (cfg *LinkerConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\tLinker config: " +
		"LinkType = %s, Timeout = %d, %s %s",
		cfg.LinkType.String(), cfg.Timeout, cfg.Serial.String(), cfg.HidUsb.String())
	return str
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
			Parity:   NoParity,
		},
		HidUsb: &HidUsbConfig{},
	}
	return lnkCfg
}
