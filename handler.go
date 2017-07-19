package main

import (
	"net/http"
	"encoding/json"
	"os"

	"github.com/bluele/slack"
	"time"
)

var (
	firstCountDay = time.Now().UTC()
	dailyCounter = 0
)

func parseFormHandler(writer http.ResponseWriter, request *http.Request) {

	request.ParseForm()
	userParams := make(map[string]string)
	var chatMessage string

	if dailyCounter < config.Security_daily_quota {
		// Collect post form params
		for key, _ := range request.Form {
			userParams[key] = request.Form.Get(key)
			chatMessage = chatMessage + key + ": " + request.Form.Get(key) + "\n"
		}
		log.Infof("Post form received: %v", userParams)

		// Notice to a Slack channel
		api := slack.New(config.Slack_api_token)
		//	options := new(slack.ChatPostMessageOpt)
		//	options.AsUser = true
		//	options.Username = "Developer Site"
		var err error
		if os.Getenv("ENV") != "dev" {
			err = api.ChatPostMessage(config.Slack_channel_name, chatMessage, nil)
		}
		if err != nil {
			log.Error("Slack ChatPostMessage: " + err.Error())
			writer.WriteHeader(500)
		} else {
			log.Infof("Post form [%v] sent to Slack channel: %s", dailyCounter, config.Slack_channel_name)
			http.Redirect(writer, request, config.Thankyou_page, 301)
		}

		// TODO Send form post params received to email
	} else {
		log.Errorf("Daily quota [%v of %v] reached out! You must wait until tomorrow morning to renew the quota.", config.Security_daily_quota, dailyCounter)
		writer.WriteHeader(500)
	}

	checkDailyQuota()

}


func checkDailyQuota() {

	if dailyCounter == 0 {
		firstCountDay = time.Now().UTC()
	} else if dailyCounter > config.Security_daily_quota {
		today := time.Now().UTC()
		if today.Year() > firstCountDay.Year() || today.YearDay() > firstCountDay.YearDay() {
			dailyCounter = 0
		}
	}
	dailyCounter++

}



func infoHandler(w http.ResponseWriter, request *http.Request) {
	// Set json response struct
	var inforesponse InfoResponseType
	inforesponse.Version = "release/4.0.16.1"
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



