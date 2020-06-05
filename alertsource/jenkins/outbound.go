package jenkins

import (
	"encoding/json"
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

func (*OutboundMessage) transform(data io.ReadCloser) (string, error) {
	messageBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return "", err
	}
	om := OutboundMessage{}
	err = json.Unmarshal(messageBytes, &om)
	if err != nil {
		return "", err
	}
	lines := []string{
		"Status: " + strings.ToUpper(om.Event),
		"Build: " + om.ProjectName + " :: " + om.BuildName,
		"URL: " + om.BuildURL,
	}
	return strings.Join(lines, "\n"), nil
}

// Parse implement Payload.Parse()
func (m OutboundMessage) Parse(req *http.Request) (string, error) {
	return m.transform(req.Body)
}
