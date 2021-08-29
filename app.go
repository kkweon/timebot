package timebot

import (
	"log"

	"github.com/nlopes/slack"
)

// Main runs app
func Main(token string, debug bool) {
	api := slack.New(token, slack.OptionDebug(debug))

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:

		case *slack.MessageEvent:
			if zone, ok := IsTargetMessage(ev.Msg.Text); ok {
				t, err := ParseTime(ev.Msg.Text)

				if err != nil {
					continue
				}

				var targetTime string

				switch zone {
				case KST:
					targetTime = ToCaliforniaTime(t)
				case PST, PDT:
					targetTime = ToKoreaTime(t)
				default:
					// do nothing
				}

				if targetTime != "" {
					rtm.SendMessage(rtm.NewOutgoingMessage(targetTime, ev.Msg.Channel))
				}
			}

		case *slack.PresenceChangeEvent:

		case *slack.LatencyReport:

		case *slack.RTMError:
			log.Println("Error:", ev.Error())

		case *slack.InvalidAuthEvent:
			log.Fatalln("Invalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}
