package routes

import (
	"github.com/gorilla/mux"
	"github.com/suhaasya/export-wizard/controllers"
)

func RootRoutes(r *mux.Router) {
	protectedR := r.NewRoute().Subrouter()
	protectedR.HandleFunc("/pdf", controllers.ExportPdf).Methods("POST")
	protectedR.HandleFunc("/excel", controllers.ExportExcel).Methods("POST")
	protectedR.HandleFunc("/csv", controllers.ExportCSV).Methods("POST")
}
