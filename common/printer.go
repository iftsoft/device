package common

const (
	CmdInitPrinter     = "InitPrinter"
	CmdPrintText       = "PrintText"
	CmdPrinterProgress = "PrinterProgress"
)

type PrinterQuery struct {
	Text string	`json:"text"`
}

type PrinterSetup struct {
	PaperPath int32	`json:"paperPath"`
	Landscape bool	`json:"landscape"`
	ShowImage int32	`json:"showImage"`
}

type PrinterProgress struct {
	DocName  string	`json:"docName"`
	PageDone int32	`json:"pageDone"`
	PagesAll int32	`json:"pagesAll"`
}

type PrinterCallback interface {
	PrinterProgress(name string, reply *PrinterProgress) error
}

type PrinterManager interface {
	InitPrinter(name string, query *PrinterSetup) error
	PrintText(name string, query *PrinterQuery) error
}
