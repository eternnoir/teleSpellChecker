package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/eternnoir/gotelebot"
	"github.com/eternnoir/teleSpellChecker/handlers"
	"github.com/evalphobia/logrus_fluent"
	"os"
)

func initLog(fluentAddr string, fluentPort int) {
	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)
	if fluentAddr != "" {
		hook := logrus_fluent.NewHook(fluentAddr, fluentPort)
		hook.SetLevels([]log.Level{
			log.PanicLevel,
			log.ErrorLevel,
			log.InfoLevel,
			log.DebugLevel,
		})
		log.AddHook(hook)
	}
	log.SetLevel(log.DebugLevel)
}

func main() {
	tokenPtr := flag.String("t", "", "Telegram Bot Token")
	fluentAddrPtr := flag.String("fluent", "", "fluent addr")
	fluentportPtr := flag.Int("fport", 24224, "fluent addr")
	flag.Parse()
	if len(*tokenPtr) < 1 {
		fmt.Println("Token can not be empty.Use -h to get some help.")
		return
	}
	initLog(*fluentAddrPtr, *fluentportPtr)
	bot := gotelebot.InitTeleBot(*tokenPtr)
	log.WithFields(log.Fields{
		"tag": "fluent",
	}).Info("Start")
	log.Info("BotStart Use Toke :" + *tokenPtr)
	go func() {
		for {
			err := bot.StartPolling(false)
			if err != nil {
				log.WithFields(log.Fields{
					"tag":     "fluent",
					"botname": "spellchecker",
				}).Error(err)
			}
		}
	}()
	newMsgChan := bot.Messages
	for {
		m := <-newMsgChan
		handlers.MainHandler(bot, m)
	}
}
