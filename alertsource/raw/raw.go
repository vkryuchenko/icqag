package raw

import (
	"io"
	"io/ioutil"
	"net/http"
)

// Message represent data from any data
type Message []byte

func transformMessage(data io.ReadCloser) (string, error) {
	messageBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return "", err
	}
	return string(messageBytes), nil
}

// Parse implement Payload.Parse()
func (m Message) Parse(req *http.Request) (string, error) {
	return transformMessage(req.Body)
}
