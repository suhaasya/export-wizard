package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/suhaasya/export-wizard/routes"
)

func setupRoutes(router *mux.Router) {
	routes.RootRoutes(router)
}

func createRouter() *http.Handler {
	allowed_origins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), " ")
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   allowed_origins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	r := mux.NewRouter()

	setupRoutes(r)
	routerProtected := corsHandler.Handler(r)
	return &routerProtected
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}

	PORT := os.Getenv("PORT")

	router := createRouter()
	fmt.Println("Starting the server on port " + PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", PORT), *router))
}
