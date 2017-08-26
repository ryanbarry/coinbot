package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
	"github.com/ryanbarry/coinbot/btcaverage"
)

func main() {
	debugOn := flag.Bool("debug", false, "Enable debug logging?")
	slackToken := flag.String("slackToken", "", "API token for the Slack bot account to use")
	flag.Parse()

	if *debugOn {
		log.Println("Debug logging turned on.")
	}

	btcusdTracker, err := btcaverage.NewGlobalTracker()
	if err != nil {
		log.Fatal("Could not initialize the Global BTC Tracker! Error: ", err.Error())
	}

	if *slackToken == "" {
		log.Fatalln("Error: Slack token not configured; not connecting to Slack!")
	} else {
		slackApi := slack.New(*slackToken)
		slackApi.SetDebug(*debugOn)
		if *debugOn {
			slack.SetLogger(log.New(os.Stderr, "[Slack] ", log.LstdFlags|log.LUTC))
		}

		slackRtm := slackApi.NewRTM()
		go slackRtm.ManageConnection()

		// slackRtmInfo := slackRtm.GetInfo()

	Loop:
		for {
			msg := <-slackRtm.IncomingEvents

			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				if ev.Type == "message" && ev.SubMessage == nil {

					log.Printf("Message from %s/%s in channel %s: %q\n", ev.User, ev.Team, ev.Channel, ev.Text)
					if strings.Contains(ev.Text, "$BTC") {
						ticker := btcusdTracker.GetAvg("BTCUSD")
						text := fmt.Sprintf("Bitcoin's current price is $%.2f USD.", ticker.Last)
						msg := slackRtm.NewOutgoingMessage(text, ev.Channel)
						slackRtm.SendMessage(msg)
					}
				} else {
					if ev.SubMessage != nil {
						log.Printf("Got message: %+v and submessage: %+v", ev, ev.SubMessage)
					} else {
						log.Printf("Got message: %+v\n", ev)
					}
				}
			case *slack.ChannelJoinedEvent:
				joined := ev.Channel
				log.Printf("Joined channel %q!\n", joined.Name)
				slackRtm.SendMessage(slackRtm.NewOutgoingMessage("Hi! I'm coinbot, and if I hear someone say `$BTC` I will post the USD price of the last BTC trade.", joined.ID))
			case *slack.InvalidAuthEvent:
				log.Fatalf("Invalid Slack credentials.")
				break Loop
			}

		}
	}
}
