package webhook

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/mail-ru-im/bot-golang"
	"go.uber.org/zap"
	"icqag/alertsource/alertmanager"
	"icqag/alertsource/grafana"
	"icqag/alertsource/jenkins"
	"icqag/alertsource/json"
	"icqag/alertsource/raw"
	"icqag/alertsource/teamcity"
	"net/http"
)

var payloadSourceMap = map[string]Payload{
	"json":                json.Message{},
	"raw":                 raw.Message{},
	"grafana":             grafana.Message{},
	"jenkins-outbound":    jenkins.OutboundMessage{},
	"jenkins-suspect":     jenkins.SuspectMessage{},
	"alertmanager":        alertmanager.Message{},
	"teamcity-elasticdoc": teamcity.ElasticsearchDocumentMessage{},
}

// Payload interface for any data from any alert systems
type Payload interface {
	Parse(r *http.Request, logger *zap.Logger) (string, error)
}

// Provider represent single instances of bot and echo
type Provider struct {
	Bot    *botgolang.Bot
	Logger *zap.Logger
}

func (Provider) payloadBySourceName(sourceName string) (Payload, error) {
	payload, ok := payloadSourceMap[sourceName]
	if !ok {
		return nil, errors.New("unknown alert source")
	}
	return payload, nil
}

// Start prepare routes and serve them
func (p *Provider) Start() error {
	router := chi.NewRouter()
	router.Put("/{source}/{target}", p.handleMessage)
	router.Post("/{source}/{target}", p.handleMessage)
	router.Get(
		"/health",
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusOK)

		},
	)
	return http.ListenAndServe(":8888", router)
}
