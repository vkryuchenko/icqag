package raw

import (
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
)

// Message represent data from any data
type Message []byte

func transformMessage(data io.ReadCloser, logger *zap.Logger) (string, error) {
	messageBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return "", err
	}
	logger.Debug(string(messageBytes))
	return string(messageBytes), nil
}

// Parse implement Payload.Parse()
func (m Message) Parse(req *http.Request, logger *zap.Logger) (string, error) {
	return transformMessage(req.Body, logger)
}
