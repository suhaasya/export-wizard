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

	// You can now iterate over the data and access the keys/values dynamically
	singleItem := data[0]

	idx := 0
	for key := range singleItem {
		f.SetCellValue("Sheet1", string(rune(66+idx))+"2", key)
		idx++
	}

	for i, item := range data {
		idx := 0
		for _, value := range item {
			f.SetCellValue("Sheet1", string(rune(66+idx))+""+fmt.Sprint(3+i), value)
			idx++
		}
	}

	// Save spreadsheet by the given path.
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}

	helpers.SetHeaders("post", w, http.StatusCreated)
	json.NewEncoder(w).Encode("created")
}
