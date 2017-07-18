package main

import (
	"net/http"
	"os"
	"encoding/json"

	"github.com/joliva-ob/http-form-handler/configuration"
	"github.com/op/go-logging"
	"github.com/bluele/slack"
)

var (
	log *logging.Logger
	config configuration.ConfigType
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
	config = configuration.LoadConfiguration(filename)
	log = configuration.GetLog()

	server := http.Server{Addr: ":" + config.Server_port}
	http.HandleFunc("/formhandler", parseFormHandler)
	http.HandleFunc( "/info", infoHandler)
	http.HandleFunc( "/health", healthHandler)
	// TODO add health and info endpoints
	server.ListenAndServe()

}


func parseFormHandler(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()
	userParams := make(map[string]string)
	var chatMessage string

	// Collect post form params
	for key, _ := range request.Form {
		userParams[key] = request.Form.Get(key)
		chatMessage = chatMessage + key + ": " + request.Form.Get(key) + "\n"
	}
	log.Info("Post form received: %v", userParams)

	// Notice to a Slack channel
	api := slack.New(config.Slack_api_token)
//	options := new(slack.ChatPostMessageOpt)
//	options.AsUser = true
//	options.Username = "Developer Site"
	err := api.ChatPostMessage(config.Slack_channel_name, chatMessage, nil)
	if err != nil {
		log.Error("Slack ChatPostMessage: " + err.Error())

		writer.WriteHeader(500)
	} else {
		log.Info("Post form sent to Slack channel: %v", config.Slack_channel_name)

		//writer.WriteHeader(200)
		http.Redirect(writer, request, "http://developer.oneboxtm.com/thank_you_page.html", 301)
	}

	// Send form post params received to email
	// TODO Send form post params received to email

}


func infoHandler(w http.ResponseWriter, request *http.Request) {
	// Set json response struct
	var inforesponse InfoResponseType
	inforesponse.Version = "release/4.0.16"
	infojson, _ := json.Marshal(inforesponse)

	// Set response headers and body
	w.Header().Set("Content-Type", "application/json")
	w.Write(infojson)
}


func healthHandler(w http.ResponseWriter, request *http.Request) {
	// Set json response struct
	var healthresponse HealthResponseType
	healthresponse.Status = STATUS_UP
	// TODO fill the discovery and other resources statuses
	healthjson, _ := json.Marshal(healthresponse)

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Set response body
	w.Write(healthjson)
}

