package controllers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"

	"github.com/johnfercher/maroto/v2"

	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"

	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/suhaasya/export-wizard/helpers"
	"github.com/xuri/excelize/v2"
)

func ExportPdf(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data []map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	m := GetMaroto(data)

	document, err := m.Generate()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = document.Save("billingv2.pdf")
	if err != nil {
		log.Fatal(err.Error())
	}

	err = document.GetReport().Save("billingv2.txt")
	if err != nil {
		log.Fatal(err.Error())
	}

	helpers.SetHeaders("post", w, http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

func ExportExcel(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data []map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Set headers
	headers := make([]string, 0, len(data[0]))
	for key := range data[0] {
		headers = append(headers, key)
	}

	// Write headers
	for i, header := range headers {
		col := string(rune('A' + i))
		cell := col + "1"
		f.SetCellValue("Sheet1", cell, header)
		f.SetColWidth("Sheet1", col, col, 40) // Set column width to 20
	}

	// Write data rows
	for rowIndex, item := range data {
		for colIndex, header := range headers {
			col := string(rune('A' + colIndex))
			cell := col + fmt.Sprint(rowIndex+2)
			value := item[header]
			f.SetCellValue("Sheet1", cell, value)
		}
	}
	// Save spreadsheet by the given path.
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}

	helpers.SetHeaders("post", w, http.StatusCreated)
	json.NewEncoder(w).Encode("created")
}

func ExportCSV(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data []map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// Create a new file
	file, err := os.Create("output.csv")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush() // Ensure all data is written

	// Set headers
	headers := make([]string, 0, len(data[0]))
	for key := range data[0] {
		headers = append(headers, key)
	}
	if err := writer.Write(headers); err != nil {
		fmt.Println("Error writing record to file:", err)
		return
	}

	// Write data rows
	for _, item := range data {
		record := make([]string, 0, len(data[0]))
		for _, header := range headers {
			value := item[header]
			if strValue, ok := value.(string); ok {
				record = append(record, strValue)
			} else {
				record = append(record, fmt.Sprintf("%v", value)) // Convert non-string values to string
			}
		}
		if err := writer.Write(record); err != nil {
			fmt.Println("Error writing record to file:", err)
			return
		}
	}

	fmt.Println("CSV file created successfully")
}

func GetMaroto(data []map[string]interface{}) core.Maroto {
	cfg := config.NewBuilder().
		WithPageNumber().
		WithLeftMargin(10).
		WithTopMargin(15).
		WithRightMargin(10).
		Build()

	mrt := maroto.New(cfg)
	m := maroto.NewMetricsDecorator(mrt)

	headers := make([]string, 0, len(data[0]))
	for key := range data[0] {
		headers = append(headers, key)
	}

	columns := make([]core.Col, len(headers))

	for i, header := range headers {
		columns[i] = text.NewCol(3, header, props.Text{Size: 9, Align: align.Left, Style: fontstyle.Bold})
	}

	rows := []core.Row{
		row.New(5).Add(
			columns...,
		),
	}

	var contentsRow []core.Row

	for _, item := range data {
		record := make([]string, 0, len(data[0]))
		for _, header := range headers {
			value := item[header]
			if strValue, ok := value.(string); ok {
				record = append(record, strValue)
			} else {
				record = append(record, fmt.Sprintf("%v", value)) // Convert non-string values to string
			}
		}

		columns := make([]core.Col, len(record))

		for i, value := range record {
			columns[i] = text.NewCol(3, value, props.Text{Size: 9, Align: align.Left, Style: fontstyle.Bold})
		}

		rows := []core.Row{
			row.New(4).Add(
				columns...,
			),
		}

		contentsRow = append(contentsRow, rows...)
	}

	rows = append(rows, contentsRow...)

	m.AddRows(rows...)

	return m
}
