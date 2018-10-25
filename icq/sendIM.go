package icq

import (
	"net/url"
	"strings"
)

// SendIm method to send any text message
func (bot *Bot) SendIm(to, message string, mentions []string) (string, error) {
	params := url.Values{
		"t":       []string{to},
		"message": []string{message},
		//"mentions": []string{},
		//"parse": []string{}
	}
	if mentions != nil {
		params["mentions"] = []string{strings.Join(mentions, ",")}
	}
	data := apiResponse{}
	err := bot.postText("/im/sendIM", params, &data)
	if err != nil {
		return "", err
	}
	return data.msgID(), data.msgError()
}
