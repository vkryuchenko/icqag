package main

import (
	"github.com/mail-ru-im/bot-golang"
	"go.uber.org/zap"
	"icqag/webhook"
	"os"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	var options []botgolang.BotOption
	apiUrl := os.Getenv("API_URL")
	if apiUrl != "" {
		options = append(options, botgolang.BotApiURL(apiUrl))
	}
	debug := os.Getenv("DEBUG")
	if debug == "true" {
		options = append(options, botgolang.BotDebug(true))
	}
	bot, err := botgolang.NewBot(os.Getenv("BOT_TOKEN"), options...)
	if err != nil {
		logger.Fatal(err.Error())
	}
	webProvider := webhook.Provider{Bot: bot, Logger: logger}
	logger.Fatal(webProvider.Start().Error())
}
