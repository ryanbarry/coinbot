package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nlopes/slack"
	"github.com/ryanbarry/coinbot/btcaverage"
)

var avg *btcaverage.GlobalAvg

func main() {
	debugOn := flag.Bool("debug", false, "Enable debug logging?")
	slackToken := flag.String("slackToken", "", "API token for the Slack bot account to use")
	flag.Parse()

	if *debugOn {
		log.Println("Debug logging turned on.")
	}

	var err error
	if avg, err = btcaverage.GetCurrentGlobalAvg("USD"); err != nil {
		log.Fatal("Error initializing current BTC price: ", err.Error())
	}

	if *slackToken == "" {
		log.Println("Slack token not configured; not connecting to Slack!")
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
			case *slack.HelloEvent:
				log.Println("Got Hello from Slack.")
			case *slack.ConnectedEvent:
				log.Println("Connection count: ", ev.ConnectionCount)
			case *slack.MessageEvent:
				log.Printf("Got a message: %+v\n", ev)
				if strings.Contains(ev.Text, "$BTC") {
					avg = getCurrentBitcoinGlobalAvg()
					msg := slackRtm.NewOutgoingMessage(fmt.Sprintf("Bitcoin's current price is $%.2f USD.", avg.Last), ev.Channel)
					slackRtm.SendMessage(msg)
				}
			case *slack.PresenceChangeEvent:
				log.Printf("Presence change: %+v\n", ev)
			case *slack.LatencyReport:
				log.Printf("Current latency: %.0f\n", ev.Value)
			case *slack.RTMError:
				log.Printf("Slack Error: %s\n", ev.Error())
			case *slack.ChannelJoinedEvent:
				joined := ev.Channel
				log.Printf("Joined channel %q!\n", joined.Name)
				slackRtm.SendMessage(slackRtm.NewOutgoingMessage("Hi! I'm coinbot, and if I hear someone say `$BTC` I will post the USD price of the last BTC trade.", joined.ID))
			case *slack.InvalidAuthEvent:
				log.Fatalf("Invalid Slack credentials.")
				break Loop
			default:
				log.Printf("Got some other event from Slack: %+v\n", ev)
			}

		}
	}
}

func getCurrentBitcoinGlobalAvg() *btcaverage.GlobalAvg {
	currency := "USD"

	if avg == nil || time.Since(avg.Timestamp.Time) > time.Minute {
		newAvg, err := btcaverage.GetCurrentGlobalAvg(currency)
		if err != nil {
			log.Printf("Error getting global average!\n%v", err)
			return avg
		}
		avg = newAvg
	}

	return avg
}
