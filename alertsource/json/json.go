package json

import (
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
)

// Message represent data from any JSON
type Message map[string]interface{}

func transformMessage(data io.ReadCloser, logger *zap.Logger) (string, error) {
	messageBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return "", err
	}
	logger.Debug(string(messageBytes))
	message := Message{}
	err = json.Unmarshal(messageBytes, &message)
	if err != nil {
		return "", err
	}
	indentJSON, err := json.MarshalIndent(message, "", " ")
	if err != nil {
		return "", err
	}
	return string(indentJSON[:]), nil
}

// Parse implement Payload.Parse()
func (m Message) Parse(req *http.Request, logger *zap.Logger) (string, error) {
	return transformMessage(req.Body, logger)
}
