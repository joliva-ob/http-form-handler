package main

import (
	"net/http"
	"os"

	"github.com/joliva-ob/http-form-handler/configuration"
	"github.com/op/go-logging"
	"github.com/bluele/slack"
)

var (
	log *logging.Logger
	config configuration.ConfigType
)

func main() {

	// Load configuration in order to start application
	var filename = os.Getenv("CONF_PATH") + "/" + os.Getenv("ENV") + ".yml"
	config = configuration.LoadConfiguration(filename)
	log = configuration.GetLog()

	server := http.Server{Addr: ":" + config.Server_port}
	http.HandleFunc("/formhandler", parseFormHandler)
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
	err := api.ChatPostMessage(config.Slack_channel_name, chatMessage, nil)
	if err != nil {
		log.Error("Slack ChatPostMessage: " + err.Error())
		writer.WriteHeader(500)
	} else {
		log.Info("Post form sent to Slack channel: %v", config.Slack_channel_name)
		writer.WriteHeader(200)
	}

	// Send form post params received to email
	// TODO Send form post params received to email

}


