package pdfexport

import (
	"fmt"
	"strings"
	"time"
)

type SingleLogbookEntry struct {
	Date  string
	MDate string

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

		CrossCountry string
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
	LogbookRows          int           `json:"logbook_rows"`
	Fill                 int           `json:"fill"`
	LeftMargin           float64       `json:"left_margin"`
	LeftMarginA          float64       `json:"left_margin_a"`
	LeftMarginB          float64       `json:"left_margin_b"`
	TopMargin            float64       `json:"top_margin"`
	BodyRow              float64       `json:"body_row_height"`
	FooterRow            float64       `json:"footer_row_height"`
	PageBreaks           string        `json:"page_breaks"`
	Columns              ColumnsWidth  `json:"columns"`
	Headers              ColumnsHeader `json:"headers"`
	ReplaceSPTime        bool          `json:"replace_sp_time"`
	IncludeSignature     bool          `json:"include_signature"`
	IsExtended           bool          `json:"is_extended"`
	TimeFieldsAutoFormat byte          `json:"time_fields_auto_format"`
	CustomTitle          string        `json:"custom_title"`
	CustomTitleBlob      []byte
}

type ColumnsWidth struct {
	Col1  float64 `json:"col1"`
	Col2  float64 `json:"col2"`
	Col3  float64 `json:"col3"`
	Col4  float64 `json:"col4"`
	Col5  float64 `json:"col5"`
	Col6  float64 `json:"col6"`
	Col7  float64 `json:"col7"`
	Col8  float64 `json:"col8"`
	Col9  float64 `json:"col9"`
	Col10 float64 `json:"col10"`
	Col11 float64 `json:"col11"`
	Col12 float64 `json:"col12"`
	Col13 float64 `json:"col13"`
	Col14 float64 `json:"col14"`
	Col15 float64 `json:"col15"`
	Col16 float64 `json:"col16"`
	Col17 float64 `json:"col17"`
	Col18 float64 `json:"col18"`
	Col19 float64 `json:"col19"`
	Col20 float64 `json:"col20"`
	Col21 float64 `json:"col21"`
	Col22 float64 `json:"col22"`
	Col23 float64 `json:"col23"`
}

type ColumnsHeader struct {
	Date      string `json:"date"`
	Departure string `json:"departure"`
	Arrival   string `json:"arrival"`
	Aircraft  string `json:"aircraft"`
	SPT       string `json:"spt"`
	MCC       string `json:"mcc"`
	Total     string `json:"total"`
	PICName   string `json:"pic_name"`
	Landings  string `json:"landings"`
	OCT       string `json:"oct"`
	PFT       string `json:"pft"`
	FSTD      string `json:"fstd"`
	Remarks   string `json:"remarks"`
	DepPlace  string `json:"dep_place"`
	DepTime   string `json:"dep_time"`
	ArrPlace  string `json:"arr_place"`
	ArrTime   string `json:"arr_time"`
	Model     string `json:"model"`
	Reg       string `json:"reg"`
	SE        string `json:"se"`
	ME        string `json:"me"`
	LandDay   string `json:"land_day"`
	LandNight string `json:"land_night"`
	Night     string `json:"night"`
	IFR       string `json:"ifr"`
	PIC       string `json:"pic"`
	COP       string `json:"cop"`
	Dual      string `json:"dual"`
	Instr     string `json:"instr"`
	SimType   string `json:"sim_type"`
	SimTime   string `json:"sim_time"`
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

	totals.Time.CrossCountry = dtoa(atod(totals.Time.CrossCountry) + atod(record.Time.CrossCountry))

	return totals
}

// atod converts formatted string to time.Duration
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

// exported dtoa function

// dtoa converts time.Duration to formatted string
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
