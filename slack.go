package main

import (
	"log"
	"os"

	"github.com/nlopes/slack"
)

func StartSlackBot() {
	api := slack.New(conf.Slack.Token)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(false)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
		case *slack.MessageEvent:
			log.Println("Message: ", ev)

			user, err := api.GetUserInfo(ev.Msg.User)
			//message := ev.Msg.Text
			if err != nil {
				log.Println("Error get user name ", err)
			}
			if user != nil {
				if user.IsBot == false {
					rtm.SendMessage(rtm.NewOutgoingMessage("@"+user.Name+" "+ev.Msg.Text, ev.Msg.Channel))
				}
			}

		case *slack.PresenceChangeEvent:
			//fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			//fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			log.Println("Error:", ev.Error())

		case *slack.InvalidAuthEvent:
			log.Println("Invalid credentials")
			return

		default:
			log.Println("Unexpected: ", msg.Data)
		}
	}

}
