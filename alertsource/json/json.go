package json

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// Message represent data from any JSON
type Message map[string]interface{}

func transformMessage(data io.ReadCloser) (string, error) {
	messageBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return "", err
	}
	rm := Message{}
	err = json.Unmarshal(messageBytes, &rm)
	if err != nil {
		return "", err
	}
	indentJSON, err := json.MarshalIndent(rm, "", " ")
	if err != nil {
		return "", err
	}
	return string(indentJSON[:]), nil
}

// Parse implement Payload.Parse()
func (m Message) Parse(req *http.Request) (string, error) {
	return transformMessage(req.Body)
}
