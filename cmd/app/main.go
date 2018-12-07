package main

import (
	"flag"
	"log"
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

	timebot.Main(*token, *debug)
}
