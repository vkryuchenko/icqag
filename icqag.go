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
	bot, err := botgolang.NewBot(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	webProvider := webhook.Provider{Bot: bot}
	log.Fatal(webProvider.Start())
}
