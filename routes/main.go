package routes

import (
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func setupRoutes(router *mux.Router) {
	RootRoutes(router)
}

func CreateRouter() *http.Handler {
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
