package routes

import (
	"github.com/gorilla/mux"
	"github.com/suhaasya/export-wizard/controllers"
)

func RootRoutes(r *mux.Router) {
	protectedR := r.NewRoute().Subrouter()
	protectedR.HandleFunc("/", controllers.ExportPdf).Methods("GET")
}
