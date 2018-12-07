package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/kkweon/timebot"
)

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

	http.HandleFunc("/healthcheck", root)
	log.Println("The server is running at :" + port)
	http.ListenAndServe(":"+port, nil)
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("true"))
}
