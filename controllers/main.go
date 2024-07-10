package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/suhaasya/export-wizard/helpers"
)

func ExportPdf(w http.ResponseWriter, r *http.Request) {
	fmt.Println("exportPdf")
	data := []int{1, 2, 3, 4}

	helpers.SetHeaders("post", w, http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}
