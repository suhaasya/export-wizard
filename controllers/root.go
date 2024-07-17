package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/suhaasya/export-wizard/helpers"
	"github.com/xuri/excelize/v2"
)

func ExportPdf(w http.ResponseWriter, r *http.Request) {
	fmt.Println("exportPdf")
	vars := mux.Vars(r)
	urlId := vars["id"]
	code := r.URL.Query().Get("code")

	helpers.SetHeaders("post", w, http.StatusCreated)
	json.NewEncoder(w).Encode(urlId + " " + code)
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
