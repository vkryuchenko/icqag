package alertmanager

import (
	"encoding/json"
	"go.uber.org/zap"
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

func (*Message) transform(data io.ReadCloser, logger *zap.Logger) (string, error) {
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
	lines := []string{
		message.Status,
		message.GroupKey,
		"Alerts:",
	}
	for _, alert := range message.Alerts {
		lines = append(
			lines,
			"Started: "+alert.StartsAt,
			"Ends: "+alert.EndsAt,
			"Status: "+alert.Status,
			"URL: "+alert.GeneratorURL,
			"\n",
		)
	}
	return strings.Join(lines, "\n"), nil
}

// Parse implement Payload.Parse()
func (gm Message) Parse(req *http.Request, logger *zap.Logger) (string, error) {
	return gm.transform(req.Body, logger)
}
