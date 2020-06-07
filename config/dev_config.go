package config

import (
	"fmt"
	"github.com/iftsoft/device/common"
)

type CommonConfig struct {
	Model    string `yaml:"model"`
	Version  string `yaml:"version"`
	Timeout  int32  `yaml:"timeout"`
	AutoLoad bool   `yaml:"auto_load"`
}
func (cfg *CommonConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\tCommon config: " +
		"Model = %s, Version = %s, Timeout = %d, AutoLoad = %t.",
		cfg.Model, cfg.Version, cfg.Timeout, cfg.AutoLoad)
	return str
}
func GetDefaultCommonConfig() *CommonConfig {
	cfg := &CommonConfig{
		Model:    "",
		Version:  "",
		Timeout:  60,
		AutoLoad: false,
	}
	return cfg
}


type PrinterConfig struct {
	PrintName   string 			`yaml:"print_name"`
	Landscape   bool			`yaml:"landscape"`
	PaperPath   EnumPaperPath	`yaml:"paper_path"`
	ShowImage   EnumShowImage	`yaml:"show_image"`
	ImageFile   string			`yaml:"image_file"`
	EjectLength int32			`yaml:"eject_length"`
}
func (cfg *PrinterConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\tPrinter config: " +
		"PrintName = %s, Landscape = %t, PaperPath = %s, ShowImage = %s, ImageFile = %s.",
		cfg.PrintName, cfg.Landscape, cfg.PaperPath, cfg.ShowImage, cfg.ImageFile)
	return str
}
func GetDefaultPrinterConfig() *PrinterConfig {
	cfg := &PrinterConfig{
	}
	return cfg
}


type ReaderConfig struct {
	SkipPrefix	EnumSkipPrefix	`yaml:"skip_prefix"`
	CardAccept	EnumCardAccept	`yaml:"card_accept"`
}
func (cfg *ReaderConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\tReader config: " +
		"SkipPrefix = %s, CardAccept = %s.",
		cfg.SkipPrefix, cfg.CardAccept)
	return str
}
func GetDefaultReaderConfig() *ReaderConfig {
	cfg := &ReaderConfig{
	}
	return cfg
}


type PinPadConfig struct {
	NeedEnter bool   `yaml:"need_enter"`
	PinDigits uint16 `yaml:"pin_digits"`
}
func (cfg *PinPadConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\tPIN pad config: " +
		"NeedEnter = %t, PinDigits = %d.",
		cfg.NeedEnter, cfg.PinDigits)
	return str
}
func GetDefaultPinPadConfig() *PinPadConfig {
	cfg := &PinPadConfig{
	}
	return cfg
}


type ValidatorConfig struct {
	NotesMask	int64              `yaml:"notes_mask"`
	NoteAlert	int32              `yaml:"note_alert"`
	NoteLimit	int32              `yaml:"note_limit"`
	ActDefault	EnumBillAction	   `yaml:"act_default"`
	StoreWait	int32              `yaml:"store_wait"`
	CurrCode	common.DevCurrency `yaml:"curr_code"`
}
func (cfg *ValidatorConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\tValidator config: " +
		"NotesMask = %d, NoteAlert = %d, NoteLimit = %d, ActDefault = %s, StoreWait = %d, CurrCode = %d.",
		cfg.NotesMask, cfg.NoteAlert, cfg.NoteLimit, cfg.ActDefault, cfg.StoreWait, cfg.CurrCode)
	return str
}
func GetDefaultValidatorConfig() *ValidatorConfig {
	cfg := &ValidatorConfig{
	}
	return cfg
}


type DispenserConfig struct {
	OutputDir EnumOutputDir `yaml:"output_dir"`
	UseDivert EnumUnitUsage `yaml:"use_divert"`
	UseEscrow EnumUnitUsage `yaml:"use_escrow"`
}
func (cfg *DispenserConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\tDispenser config: " +
		"OutputDir = %s, UseDivert = %s, UseEscrow = %s.",
		cfg.OutputDir.String(), cfg.UseDivert, cfg.UseEscrow)
	return str
}
func GetDefaultDispenserConfig() *DispenserConfig {
	cfg := &DispenserConfig{
	}
	return cfg
}


type VendorConfig struct {
	UnitIndex int32 `yaml:"unit_index"`
	ItemAlert int32 `yaml:"item_alert"`
}
func (cfg *VendorConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\tVendor config: " +
		"UnitIndex = %d, ItemAlert = %d.",
		cfg.UnitIndex, cfg.ItemAlert)
	return str
}
func GetDefaultVendorConfig() *VendorConfig {
	cfg := &VendorConfig{
	}
	return cfg
}


type DeviceConfig struct {
	Linker    *LinkerConfig    `yaml:"linker"`
	Common    *CommonConfig    `yaml:"common"`
	Printer   *PrinterConfig   `yaml:"printer"`
	Reader    *ReaderConfig    `yaml:"reader"`
	Pinpad    *PinPadConfig    `yaml:"pinpad"`
	Validator *ValidatorConfig `yaml:"validator"`
	Dispenser *DispenserConfig `yaml:"dispenser"`
	Vendor    *VendorConfig    `yaml:"vendor"`
}
func (cfg *DeviceConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\nDevice config: %s %s %s %s %s %s %s %s",
		cfg.Linker, cfg.Common, cfg.Printer, cfg.Reader,
		cfg.Pinpad, cfg.Validator, cfg.Dispenser, cfg.Vendor)
	return str
}


func GetDefaultDeviceConfig() *DeviceConfig {
	devCfg := &DeviceConfig{
		Linker:    GetDefaultLinkerConfig(),
		Common:    GetDefaultCommonConfig(),
		Printer:   GetDefaultPrinterConfig(),
		Reader:    GetDefaultReaderConfig(),
		Pinpad:    GetDefaultPinPadConfig(),
		Validator: GetDefaultValidatorConfig(),
		Dispenser: GetDefaultDispenserConfig(),
		Vendor:    GetDefaultVendorConfig(),
	}
	return devCfg
}


type ConfigOverwrite struct {
	LinkType  EnumLinkType  `yaml:"link_type"`
	PortName  string        `yaml:"port_name"`
	VendorID  uint16        `yaml:"vendor_id"`  	// Device Vendor ID
	ProductID uint16        `yaml:"product_id"` 	// Device Product ID
}

func (cfg *ConfigOverwrite) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\t\tLinker overwrite: " +
		"LinkType = %s, PortName = %s, VendorID = %X, ProductID = %X.",
		cfg.LinkType, cfg.PortName, cfg.VendorID, cfg.ProductID)
	return str
}

