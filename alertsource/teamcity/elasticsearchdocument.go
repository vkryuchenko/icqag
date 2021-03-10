package teamcity

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// {
//   "branchName": "<default>",
//   "branchDisplayName": "refs/heads/develop",
//   "branchIsDefault": "true",
//   "buildStartTime": "2021-03-10T09:39:35.799Z",
//   "timestamp": "2021-03-10T09:40:41.733Z",
//   "buildFinishTime": "2021-03-10T09:40:41.733Z",
//   "buildEvent": "buildFinished",
//   "buildName": "CodegenCheck",
//   "buildStatusUrl": "https://ci.shared.ittvrn.dev/viewLog.html?buildTypeId=ITT_Projects_Random_CodegenCheck&buildId=553",
//   "buildNumber": "16",
//   "triggeredBy": "Snapshot dependency; Michail Rublev; ITT / Random / Build / Server",
//   "buildResult": "failure",
//   "buildResultPrevious": "success",
//   "buildResultDelta": "broken"
// }

type ElasticsearchDocumentMessage struct {
	AgentName         string `json:"agentName"`
	BranchDisplayName string `json:"branchDisplayName"`
	BuildName         string `json:"buildName"`
	BuildStartTime    string `json:"buildStartTime"`
	BuildFinishTime   string `json:"buildFinishTime"`
	BuildStatusUrl    string `json:"buildStatusUrl"`
	BuildResult       string `json:"buildResult"`
	BuildResultDelta  string `json:"buildResultDelta"`
	TriggeredBy       string `json:"triggeredBy"`
}

func (*ElasticsearchDocumentMessage) transform(data io.ReadCloser) (string, error) {
	messageBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return "", err
	}
	message := ElasticsearchDocumentMessage{}
	err = json.Unmarshal(messageBytes, &message)
	if err != nil {
		return "", err
	}
	status := strings.ToUpper(message.BuildResult)
	if message.BuildResultDelta != "unchanged" {
		status += strings.ToUpper(" / " + message.BuildResultDelta)
	}
	lines := []string{
		"Status: " + status,
		"Build: " + message.BuildName,
		"Branch: " + strings.TrimLeft(message.BranchDisplayName, "refs/heads/"),
		"Agent: " + message.AgentName,
		"Triggered: " + message.TriggeredBy,
		"URL: " + message.BuildStatusUrl,
	}
	return strings.Join(lines, "\n"), nil
}

// Parse implement Payload.Parse()
func (m ElasticsearchDocumentMessage) Parse(req *http.Request) (string, error) {
	return m.transform(req.Body)
}
