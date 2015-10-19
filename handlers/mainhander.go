package handlers

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/eternnoir/gotelebot"
	"github.com/eternnoir/gotelebot/types"
	"github.com/trustmaster/go-aspell"
	"strings"
)

func MainHandler(bot *gotelebot.TeleBot, newMsg *types.Message) {
	logInfo("GetMsg:" + newMsg.Text)
	if strings.Contains(newMsg.Text, "/help") || strings.Contains(newMsg.Text, "/start") {
		showHelp(bot, newMsg)
		return
	}
	processMessage(bot, newMsg)
}

func showHelp(bot *gotelebot.TeleBot, msg *types.Message) {
	returnMessage := `
	Send me a word. I will check and suggest for you.
	`
	bot.SendMessage(int(msg.Chat.Id), returnMessage, nil)
}

func processMessage(bot *gotelebot.TeleBot, msg *types.Message) {
	word := strings.TrimSpace(msg.Text)
	logInfo(fmt.Sprintf("Get Word: %s", word))
	if len(word) > 0 {
		speller, err := getSpeller("en_US")
		if err != nil {
			logError("Error: %s", err.Error())
			bot.SendMessage(int(msg.Chat.Id), "Ooops, something wrong", nil)
			return
		}
		defer speller.Delete()
		if speller.Check(word) {
			bot.SendMessage(int(msg.Chat.Id), "OK", nil)
		} else {
			retmsg := fmt.Sprintf("Incorrect word, suggestions:\n* %s\n", strings.Join(speller.Suggest(word), "\n* "))
			bot.SendMessage(int(msg.Chat.Id), retmsg, nil)
		}

	} else {
		returnMessage := `
		Please send text message.
		`
		bot.SendMessage(int(msg.Chat.Id), returnMessage, nil)
	}
}

func getSpeller(lang string) (*aspell.Speller, error) {
	speller, err := aspell.NewSpeller(map[string]string{
		"lang": lang,
	})
	if err != nil {
		logError("Error: %s", err.Error())
		return nil, err
	}
	return &speller, nil
}

func logInfo(args ...interface{}) {
	log.WithFields(log.Fields{
		"tag":     "fluent",
		"botname": "spellchecker",
	}).Info(args)
}
func logDebug(args ...interface{}) {
	log.WithFields(log.Fields{
		"tag":     "fluent",
		"botname": "spellchecker",
	}).Debug(args)
}
func logError(args ...interface{}) {
	log.WithFields(log.Fields{
		"tag":     "fluent",
		"botname": "spellchecker",
	}).Error(args)
}
