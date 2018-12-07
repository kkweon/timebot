package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"encoding/json"
	"github.com/kkweon/timebot"
)

var clientId string
var clientSecret string

func init() {
	var ok bool
	clientId, ok = os.LookupEnv("SLACK_CLIENT_ID")
	if !ok {
		log.Fatalln("$SLACK_CLIENT_ID is not available")
	}
	clientSecret, ok = os.LookupEnv("SLACK_CLIENT_SECRET")
	if !ok {
		log.Fatalln("$SLACK_CLIENT_SECRET is not available")
	}
}

func main() {
	token := flag.String("token", "", "Slack Token")
	debug := flag.Bool("debug", false, "Turn on debug")
	flag.Parse()

	if *token == "" {
		tok, ok := os.LookupEnv("SLACK_BOT_TOKEN")
		if !ok {
			log.Fatalln("--token or env SLACK_BOT_TOKEN variable is required")
		}
		*token = tok
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Println("Unalbe to retrieve $PORT environment variable")
		port = "8080"
	}

	go timebot.Main(*token, *debug)

	http.HandleFunc("/healthcheck", healthCheck)
	http.HandleFunc("/oauth/slack", slackHandler)
	log.Println("The server is running at :" + port)
	http.ListenAndServe(":"+port, nil)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("true"))
}

type OAuthGenerated struct {
	AccessToken     string `json:"access_token"`
	Scope           string `json:"scope"`
	TeamName        string `json:"team_name"`
	TeamID          string `json:"team_id"`
	IncomingWebhook struct {
		URL              string `json:"url"`
		Channel          string `json:"channel"`
		ConfigurationURL string `json:"configuration_url"`
	} `json:"incoming_webhook"`
	Bot struct {
		BotUserID      string `json:"bot_user_id"`
		BotAccessToken string `json:"bot_access_token"`
	} `json:"bot"`
}

func slackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing ?code", http.StatusBadRequest)
		return
	}

	data := getClientSecret()
	data["code"] = []string{code}
	resp, err := http.PostForm("https://slack.com/api/oauth.access", data)

	if err != nil {
		http.Error(w, "Fail while sending a request https://slack.com/api/oauth.access", http.StatusInternalServerError)
		return
	}

	var oauth OAuthGenerated
	err = json.NewDecoder(resp.Body).Decode(&oauth)

	if err != nil {
		http.Error(w, "Fail while decoding JSON", http.StatusInternalServerError)
		return
	}

	token := oauth.Bot.BotAccessToken

	if token != "" {
		go timebot.Main(token, false)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("true"))
		return
	}

	http.Error(w, "Fail while decoding JSON", http.StatusInternalServerError)
	return
}

func getClientSecret() map[string][]string {
	args := make(map[string][]string)
	args["client_id"] = []string{clientId}
	args["client_secret"] = []string{clientSecret}

	return args
}
