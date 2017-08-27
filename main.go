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

type Options struct {
	debugOn    bool
	slackToken string
}

func readOptions() Options {
	var envDebugOn string
	var envDebugOnSet bool
	var flagDebugOn *bool
	var envSlackToken string
	var flagSlackToken *string

	envDebugOn, envDebugOnSet = os.LookupEnv("DEBUG_ON")
	envSlackToken = os.Getenv("SLACK_TOKEN")

	flagDebugOn = flag.Bool("debug", false, "Enable debug logging?")
	flagSlackToken = flag.String("slackToken", "", "API token for the Slack bot account to use")
	flag.Parse()

	finalOpts := new(Options)

	if len(*flagSlackToken) > 0 {
		finalOpts.slackToken = *flagSlackToken
	} else {
		finalOpts.slackToken = envSlackToken
	}

	if *flagDebugOn {
		finalOpts.debugOn = true
	} else {
		if envDebugOnSet && strings.ToLower(envDebugOn) != "false" {
			finalOpts.debugOn = true
		}
	}

	return *finalOpts
}

func main() {
	opts := readOptions()

	if opts.debugOn {
		log.Println("Debug logging turned on.")
		log.Println("Got slackToken \"" + opts.slackToken + "\"")
	}

	btcusdTracker, err := btcaverage.NewGlobalTracker()
	if err != nil {
		log.Fatal("Could not initialize the Global BTC Tracker! Error: ", err.Error())
	}

	if opts.slackToken == "" {
		log.Fatalln("Error: Slack token not configured; not connecting to Slack!")
	} else {
		slackApi := slack.New(opts.slackToken)
		slackApi.SetDebug(opts.debugOn)
		if opts.debugOn {
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
