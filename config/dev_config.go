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
	AutoLoad bool   `yaml:"autoLoad"`
}

type PrinterConfig struct {
	PrintName   string `yaml:"printName"`
	ImageFile   string `yaml:"imageFile"`
	PaperPath   int32  `yaml:"paperPath"`
	Landscape   bool   `yaml:"landscape"`
	ShowImage   int32  `yaml:"showImage"`
	EjectLength int32  `yaml:"ejectLength"`
}

type ReaderConfig struct {
	SkipPrefix int32 `yaml:"skipPrefix"`
	CardAccept int16 `yaml:"cardAccept"`
}

type PinPadConfig struct {
	NeedEnter bool  `yaml:"needEnter"`
	PinDigits int32 `yaml:"pinDigits"`
}

type ValidatorConfig struct {
	NotesMask  int64 `yaml:"notesMask"`
	NoteAlert  int32 `yaml:"noteAlert"`
	NoteLimit  int32 `yaml:"noteLimit"`
	ActDefault int32 `yaml:"actDefault"`
	StoreWait  int32 `yaml:"storeWait"`
	CurrCode   int32 `yaml:"currCode"`
}

type DispenserConfig struct {
	OutputDir int32 `yaml:"outputDir"`
	UseDivert int32 `yaml:"useDivert"`
	UseEscrow int32 `yaml:"useEscrow"`
}

type VendorConfig struct {
	UnitIndex int32 `yaml:"unitIndex"`
	ItemAlert int32 `yaml:"itemAlert"`
}

type DeviceConfig struct {
	Common    *CommonConfig    `yaml:"common"`
	Printer   *PrinterConfig   `yaml:"printer"`
	Reader    *ReaderConfig    `yaml:"reader"`
	Pinpad    *PinPadConfig    `yaml:"pinpad"`
	Validator *ValidatorConfig `yaml:"validator"`
	Dispenser *DispenserConfig `yaml:"dispenser"`
	Vendor    *VendorConfig    `yaml:"vendor"`
	Linker    *LinkerConfig    `yaml:"linker"`
}

func GetDefaultDeviceConfig() *DeviceConfig {
	devCfg := &DeviceConfig{
		Common:    &CommonConfig{},
		Printer:   &PrinterConfig{},
		Reader:    &ReaderConfig{},
		Pinpad:    &PinPadConfig{},
		Validator: &ValidatorConfig{},
		Dispenser: &DispenserConfig{},
		Vendor:    &VendorConfig{},
		Linker:    GetDefaultLinkerConfig(),
	}
	return devCfg
}
