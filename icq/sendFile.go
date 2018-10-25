package icq

// SendFile method upload file to server and send link to receiver
func (bot *Bot) SendFile(to, filePath string) (string, error) {
	fileURL, err := bot.UploadFile(to, filePath)
	if err != nil {
		return "", err
	}
	return bot.SendIm(to, fileURL, nil)
}
