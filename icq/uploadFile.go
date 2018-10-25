package icq

// UploadFile method to upload any file to server.
func (bot *Bot) UploadFile(to, filePath string) (string, error) {
	response := struct {
		Data struct {
			StaticURL string `json:"static_url"`
		} `json:"data"`
	}{}
	err := bot.postFile(filePath, &response)
	if err != nil {
		return "", err
	}
	return response.Data.StaticURL, nil
}
