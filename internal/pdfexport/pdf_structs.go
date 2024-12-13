package pdfexport

import (
	"fmt"
	"strings"
	"time"
)

type SingleLogbookEntry struct {
	Date string

	Departure struct {
		Place string
		Time  string
	}

	Arrival struct {
		Place string
		Time  string
	}

	Aircraft struct {
		Model string
		Reg   string
	}

	Time struct {
		SE         string
		ME         string
		MCC        string
		Total      string
		Night      string
		IFR        string
		PIC        string
		CoPilot    string
		Dual       string
		Instructor string
	}

	Landings struct {
		Day   int
		Night int
	}

	SIM struct {
		Type string 
		Time string 
	}
	PIC     string 
	Remarks string 
}

type ExportPDF struct {
	LogbookRows          int
	Fill                 int
	LeftMargin           float64
	LeftMarginA          float64
	LeftMarginB          float64
	TopMargin            float64
	BodyRow              float64
	FooterRow            float64
	PageBreaks           string
	Columns              ColumnsWidth
	Headers              ColumnsHeader
	ReplaceSPTime        bool
	IncludeSignature     bool
	TimeFieldsAutoFormat byte
	CustomTitle          string
	CustomTitleBlob      []byte
}

type ColumnsWidth struct {
	Col1  float64
	Col2  float64
	Col3  float64
	Col4  float64
	Col5  float64
	Col6  float64
	Col7  float64
	Col8  float64
	Col9  float64
	Col10 float64
	Col11 float64
	Col12 float64
	Col13 float64
	Col14 float64
	Col15 float64
	Col16 float64
	Col17 float64
	Col18 float64
	Col19 float64
	Col20 float64
	Col21 float64
	Col22 float64
	Col23 float64
}

type ColumnsHeader struct {
	Date      string
	Departure string
	Arrival   string
	Aircraft  string
	SPT       string
	MCC       string
	Total     string
	PICName   string
	Landings  string
	OCT       string
	PFT       string
	FSTD      string
	Remarks   string
	DepPlace  string
	DepTime   string
	ArrPlace  string
	ArrTime   string
	Model     string
	Reg       string
	SE        string
	ME        string
	LandDay   string
	LandNight string
	Night     string
	IFR       string
	PIC       string
	COP       string
	Dual      string
	Instr     string
	SimType   string
	SimTime   string
}

func CalculateTotals(totals SingleLogbookEntry, record SingleLogbookEntry) SingleLogbookEntry {

	totals.Time.SE = dtoa(atod(totals.Time.SE) + atod(record.Time.SE))
	totals.Time.ME = dtoa(atod(totals.Time.ME) + atod(record.Time.ME))
	totals.Time.MCC = dtoa(atod(totals.Time.MCC) + atod(record.Time.MCC))
	totals.Time.Night = dtoa(atod(totals.Time.Night) + atod(record.Time.Night))
	totals.Time.IFR = dtoa(atod(totals.Time.IFR) + atod(record.Time.IFR))
	totals.Time.PIC = dtoa(atod(totals.Time.PIC) + atod(record.Time.PIC))
	totals.Time.CoPilot = dtoa(atod(totals.Time.CoPilot) + atod(record.Time.CoPilot))
	totals.Time.Dual = dtoa(atod(totals.Time.Dual) + atod(record.Time.Dual))
	totals.Time.Instructor = dtoa(atod(totals.Time.Instructor) + atod(record.Time.Instructor))
	totals.Time.Total = dtoa(atod(totals.Time.Total) + atod(record.Time.Total))
	totals.SIM.Time = dtoa(atod(totals.SIM.Time) + atod(record.SIM.Time))
	totals.Landings.Day += record.Landings.Day
	totals.Landings.Night += record.Landings.Night

	return totals
}

// string -> time.Duration
func atod(value string) time.Duration {
	if value == "" {
		value = "0:0"
	}

	strTime := fmt.Sprintf("%sm", strings.ReplaceAll(value, ":", "h"))

	duration, err := time.ParseDuration(strTime)
	if err != nil {
		fmt.Printf("Error parsing time %s\n", strTime)
		return 0
	}

	return duration
}

// time.Duration -> string
func dtoa(value time.Duration) string {

	d := value.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute

	if h == 0 && m == 0 {
		return "0:00"
	}
	return fmt.Sprintf("%01d:%02d", h, m)
}
