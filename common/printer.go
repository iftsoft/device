package common

const (
	CmdInitPrinter     = "InitPrinter"
	CmdPrintText       = "PrintText"
	CmdPrinterProgress = "PrinterProgress"
)

type PrinterQuery struct {
	Text string `json:"text"`
}

type PrinterSetup struct {
	PaperPath int32 `json:"paper_path"`
	Landscape bool  `json:"landscape"`
	ShowImage int32 `json:"show_image"`
}

type PrinterProgress struct {
	DocName  string `json:"doc_name"`
	PageDone int32  `json:"page_done"`
	PagesAll int32  `json:"pages_all"`
}

type PrinterCallback interface {
	PrinterProgress(name string, reply *PrinterProgress) error
}

type PrinterManager interface {
	InitPrinter(name string, query *PrinterSetup) error
	PrintText(name string, query *PrinterQuery) error
}
