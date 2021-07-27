package jenkins

import (
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// {
//   "buildName":"#1",
//   "buildUrl":"https://jenkins.example.com/job/Tools/job/icqWebhook/1/",
//   "event":"start",
//   "projectName":"icqWebhook"
// }

// OutboundMessage represent data struct by Jenkins Outbound WebHook plugin
type OutboundMessage struct {
	BuildName   string `json:"buildName"`
	BuildURL    string `json:"buildUrl"`
	Event       string `json:"event"`
	ProjectName string `json:"projectName"`
}

func (*OutboundMessage) transform(data io.ReadCloser, logger *zap.Logger) (string, error) {
	messageBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return "", err
	}
	logger.Debug(string(messageBytes))
	message := OutboundMessage{}
	err = json.Unmarshal(messageBytes, &message)
	if err != nil {
		return "", err
	}
	lines := []string{
		"Status: " + strings.ToUpper(message.Event),
		"Build: " + message.ProjectName + " :: " + message.BuildName,
		"URL: " + message.BuildURL,
	}
	return strings.Join(lines, "\n"), nil
}

// Parse implement Payload.Parse()
func (m OutboundMessage) Parse(req *http.Request, logger *zap.Logger) (string, error) {
	return m.transform(req.Body, logger)
}
