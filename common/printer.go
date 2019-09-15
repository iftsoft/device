package common

const (
	CmdInitPrinter     = "InitPrinter"
	CmdPrintText       = "PrintText"
	CmdPrinterProgress = "PrinterProgress"
)

type PrinterQuery struct {
	Text string
}

type PrinterSetup struct {
	PaperPath int32
	Landscape bool
	ShowImage int32
}

type PrinterProgress struct {
	DocName  string
	PageDone int32
	PagesAll int32
}

type PrinterCallback interface {
	PrinterProgress(name string, reply *PrinterProgress) error
}

type PrinterManager interface {
	InitPrinter(name string, query *PrinterSetup) error
	PrintText(name string, query *PrinterQuery) error
}
