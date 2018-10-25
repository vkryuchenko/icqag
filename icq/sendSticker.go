package icq

import (
	"net/url"
)

// SendSticker method to send any sticker by ID
func (bot *Bot) SendSticker(to, stickerID string) error {
	params := url.Values{
		"t":         []string{to},
		"stickerId": []string{stickerID},
	}
	data := apiResponse{}
	err := bot.postText("/im/sendSticker", params, &data)
	if err != nil {
		return err
	}
	return data.msgError()
}
