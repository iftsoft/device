package config

import "fmt"

// Device error codes
const (
	CardAcceptAnyCard int16 = iota
	CardAcceptMagnetic
	CardAcceptSmart
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


type PrinterConfig struct {
	PrintName   string `yaml:"print_name"`
	ImageFile   string `yaml:"image_file"`
	PaperPath   int32  `yaml:"paper_path"`
	Landscape   bool   `yaml:"landscape"`
	ShowImage   int32  `yaml:"show_image"`
	EjectLength int32  `yaml:"eject_length"`
}
func (cfg *PrinterConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\tPrinter config: " +
		"PrintName = %s, ImageFile = %s, PaperPath = %d, Landscape = %t.",
		cfg.PrintName, cfg.ImageFile, cfg.PaperPath, cfg.Landscape)
	return str
}


type ReaderConfig struct {
	SkipPrefix int32 `yaml:"skip_prefix"`
	CardAccept int16 `yaml:"card_accept"`
}
func (cfg *ReaderConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\tReader config: " +
		"SkipPrefix = %d, CardAccept = %d.",
		cfg.SkipPrefix, cfg.CardAccept)
	return str
}


type PinPadConfig struct {
	NeedEnter bool  `yaml:"need_enter"`
	PinDigits int32 `yaml:"pin_digits"`
}
func (cfg *PinPadConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\tPIN pad config: " +
		"NeedEnter = %t, PinDigits = %d.",
		cfg.NeedEnter, cfg.PinDigits)
	return str
}


type ValidatorConfig struct {
	NotesMask  int64 `yaml:"notes_mask"`
	NoteAlert  int32 `yaml:"note_alert"`
	NoteLimit  int32 `yaml:"note_limit"`
	ActDefault int32 `yaml:"act_default"`
	StoreWait  int32 `yaml:"store_wait"`
	CurrCode   int32 `yaml:"curr_code"`
}
func (cfg *ValidatorConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\tValidator config: " +
		"NotesMask = %d, NoteAlert = %d, NoteLimit = %d, ActDefault = %d, StoreWait = %d, CurrCode = %d.",
		cfg.NotesMask, cfg.NoteAlert, cfg.NoteLimit, cfg.ActDefault, cfg.StoreWait, cfg.CurrCode)
	return str
}


type DispenserConfig struct {
	OutputDir int32 `yaml:"output_dir"`
	UseDivert int32 `yaml:"use_divert"`
	UseEscrow int32 `yaml:"use_escrow"`
}
func (cfg *DispenserConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\n\tDispenser config: " +
		"OutputDir = %d, UseDivert = %d, UseEscrow = %d.",
		cfg.OutputDir, cfg.UseDivert, cfg.UseEscrow)
	return str
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


func GetDefaultDeviceConfig(linker *LinkerConfig) *DeviceConfig {
	devCfg := &DeviceConfig{
		Linker:    linker,
		Common:    &CommonConfig{},
		Printer:   &PrinterConfig{},
		Reader:    &ReaderConfig{},
		Pinpad:    &PinPadConfig{},
		Validator: &ValidatorConfig{},
		Dispenser: &DispenserConfig{},
		Vendor:    &VendorConfig{},
	}
	return devCfg
}
