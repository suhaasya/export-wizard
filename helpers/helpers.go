package helpers

import "net/http"

var methodChoices = map[string]string{
	"get":   "GET",
	"post":  "POST",
	"patch": "PATCH",
	"del":   "DELETE",
}

func SetHeaders(type_ string, w http.ResponseWriter, status int) {
	method := methodChoices[type_]
	if method == "" {
		method = "GET"
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if method != "GET" {
		w.Header().Set("Access-Control-Allow-Methods", method)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
}
