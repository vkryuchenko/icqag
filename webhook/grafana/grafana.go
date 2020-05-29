package grafana

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

//{
//  "evalMatches": [
//    {
//      "value": 100,
//      "metric": "High value",
//      "tags": null
//    },
//    {
//      "value": 200,
//      "metric": "Higher Value",
//      "tags": null
//    }
//  ],
//  "imageUrl": "http://grafana.org/assets/img/blog/mixed_styles.png",
//  "message": "Someone is testing the alert notification within grafana.",
//  "ruleId": 0,
//  "ruleName": "Test notification",
//  "ruleUrl": "http://localhost:3000/",
//  "state": "alerting",
//  "title": "[Alerting] Test notification"
//}

// MetricValue represent
type MetricValue struct {
	Metric string  `json:"metric"`
	Value  float64 `json:"value"`
	//Tags interface{} `json:"tags"`
}

// Message represent data struct by Grafana
type Message struct {
	RuleID      uint          `json:"ruleId"`
	ImageURL    string        `json:"imageUrl"`
	Message     string        `json:"message"`
	RuleName    string        `json:"ruleName"`
	RuleURL     string        `json:"ruleUrl"`
	State       string        `json:"state"`
	Title       string        `json:"title"`
	EvalMatches []MetricValue `json:"evalMatches"`
}

func (*Message) transform(data io.ReadCloser) (string, error) {
	messageBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return "", err
	}
	gm := Message{}
	err = json.Unmarshal(messageBytes, &gm)
	if err != nil {
		return "", err
	}
	lines := []string{
		gm.Title,
		gm.Message,
	}
	for _, metric := range gm.EvalMatches {
		lines = append(lines, metric.Metric+": "+fmt.Sprint(metric.Value))
	}
	return strings.Join(lines, "\n"), nil
}

// Parse implement Payload.Parse()
func (gm Message) Parse(req *http.Request) (string, error) {
	return gm.transform(req.Body)
}
