package main

import (
	"github.com/labstack/gommon/log"
	"github.com/mail-ru-im/bot-golang"
	"icqag/webhook"
	"os"
)

func init() {
	log.SetLevel(log.INFO)
	//http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func main() {
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
		log.Fatal(err)
	}
	webProvider := webhook.Provider{Bot: bot}
	log.Fatal(webProvider.Start())
}
