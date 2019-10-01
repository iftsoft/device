package config

type EnumDevRole int32

// Device types
const (
	DevRolePrinter EnumDevRole = 1 << iota
	DevRoleBarScanner
	DevRoleCardReader
	DevRoleValidator
	DevRoleItemVendor
	DevRoleDispenser
	DevRoleTerminal
	DevRoleCustom
)

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

type PrinterConfig struct {
	PrintName   string `yaml:"print_name"`
	ImageFile   string `yaml:"image_file"`
	PaperPath   int32  `yaml:"paper_path"`
	Landscape   bool   `yaml:"landscape"`
	ShowImage   int32  `yaml:"show_image"`
	EjectLength int32  `yaml:"eject_length"`
}

type ReaderConfig struct {
	SkipPrefix int32 `yaml:"skip_prefix"`
	CardAccept int16 `yaml:"card_accept"`
}

type PinPadConfig struct {
	NeedEnter bool  `yaml:"need_enter"`
	PinDigits int32 `yaml:"pin_digits"`
}

type ValidatorConfig struct {
	NotesMask  int64 `yaml:"notes_mask"`
	NoteAlert  int32 `yaml:"note_alert"`
	NoteLimit  int32 `yaml:"note_limit"`
	ActDefault int32 `yaml:"act_default"`
	StoreWait  int32 `yaml:"store_wait"`
	CurrCode   int32 `yaml:"curr_code"`
}

type DispenserConfig struct {
	OutputDir int32 `yaml:"output_dir"`
	UseDivert int32 `yaml:"use_divert"`
	UseEscrow int32 `yaml:"use_escrow"`
}

type VendorConfig struct {
	UnitIndex int32 `yaml:"unit_index"`
	ItemAlert int32 `yaml:"item_alert"`
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

func GetDefaultDeviceConfig() *DeviceConfig {
	devCfg := &DeviceConfig{
		Linker:    GetDefaultLinkerConfig(),
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
