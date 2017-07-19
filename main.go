package main

import (
	"net/http"
	"os"
)

const (
	STATUS_UP string = "UP"
)

type InfoResponseType struct {

	Version string `json:"version"`
}

// Health response struct
type HealthResponseType struct {

	Status string `json:"status"`
}




func main() {

	// Load configuration in order to start application
	var filename = os.Getenv("CONF_PATH") + "/" + os.Getenv("ENV") + ".yml"
	config = LoadConfiguration(filename)

	server := http.Server{Addr: ":" + config.Server_port}
	http.HandleFunc("/formhandler", parseFormHandler)
	http.HandleFunc( "/info", infoHandler)
	http.HandleFunc( "/health", healthHandler)
	server.ListenAndServe()

}

