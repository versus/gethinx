package main

import (
	"log"
	"os"

	"fmt"

	"strings"

	"sync/atomic"

	"github.com/nlopes/slack"
)

func statusMsg(api *slack.Client, channel string) {

	channels := []string{channel}
	params := slack.FileUploadParameters{
		Title:    "Status Gethinx",
		Filetype: "txt",
		Content:  GetStatusTable(),
		Channels: channels,
	}
	_, err := api.UploadFile(params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

}

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
			log.Println("Message: ", ev)
			// Ignore hello

		case *slack.ConnectedEvent:
			log.Println("Message: ", ev)
		case *slack.MessageEvent:
			log.Println("Message: ", ev)

			user, err := api.GetUserInfo(ev.Msg.User)

			if err != nil {
				log.Println("Error get user name ", err)
			}
			if user != nil {
				if user.IsBot == false {
					message := ev.Msg.Text
					//TODO: add reload and update commands
					//rtm.SendMessage(rtm.NewOutgoingMessage("@"+user.Name+" "+ev.Msg.Text, ev.Msg.Channel))
					if strings.Contains(message, "status") {
						statusMsg(api, ev.Msg.Channel)
					}
					if strings.Contains(message, "last") {
						msg := fmt.Sprintf("@ %s last block is %d", user.Name, atomic.LoadInt64(&LastBlock.Dig))
						rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Msg.Channel))
					}
				}

			}

		case *slack.RTMError:
			log.Println("Error:", ev.Error())

		case *slack.InvalidAuthEvent:
			log.Println("Invalid credentials")
			return

		default:
			//log.Println("Unexpected: ", msg.Data)
		}
	}

}
