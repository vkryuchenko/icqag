package alertmanager

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

//https://prometheus.io/docs/alerting/configuration/#webhook_config

// Alert data
type Alert struct {
	Status       string      `json:"status"`
	Labels       interface{} `json:"labels"`
	Annotations  interface{} `json:"annotations"`
	StartsAt     string      `json:"startsAt"`
	EndsAt       string      `json:"endsAt"`
	GeneratorURL string      `json:"generatorURL"`
}

// Message data struct by Prometheus Alertmanager
type Message struct {
	Version      string      `json:"version"`
	GroupKey     string      `json:"groupKey"`
	Status       string      `json:"status"`
	Receiver     string      `json:"receiver"`
	GroupLabels  interface{} `json:"groupLabels"`
	CommonLabels interface{} `json:"commonLabels"`
	ExternalURL  string      `json:"externalURL"`
	Alerts       []Alert     `json:"alerts"`
}

func (*Message) transform(data io.ReadCloser) (string, []string, error) {
	messageBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return "", nil, err
	}
	msg := Message{}
	err = json.Unmarshal(messageBytes, &msg)
	if err != nil {
		return "", nil, err
	}
	lines := []string{
		msg.Status,
		msg.GroupKey,
		"Alerts:",
	}
	for _, alert := range msg.Alerts {
		lines = append(
			lines,
			"Started: "+alert.StartsAt,
			"Ends: "+alert.EndsAt,
			"Status: "+alert.Status,
			"URL: "+alert.GeneratorURL,
			"\n",
		)
	}
	return strings.Join(lines, "\n"), nil, nil
}

// Parse implement Payload.Parse()
func (gm Message) Parse(req *http.Request) (string, []string, error) {
	return gm.transform(req.Body)
}
