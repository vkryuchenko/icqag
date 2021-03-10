package webhook

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mail-ru-im/bot-golang"
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
	Parse(req *http.Request) (string, error)
}

// Provider represent single instances of bot and echo
type Provider struct {
	Bot      *botgolang.Bot
	instance *echo.Echo
}

func (Provider) payloadBySourceName(sourceName string) (Payload, error) {
	payload, ok := payloadSourceMap[sourceName]
	if !ok {
		return nil, errors.New("unknown alert source")
	}
	return payload, nil
}

func (p *Provider) initEcho() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//
	e.PUT("/:source/:target", p.handleMessage)
	e.POST("/:source/:target", p.handleMessage)
	//
	p.instance = e

}

// Start prepare echo instance and start it
func (p *Provider) Start() error {
	p.initEcho()
	return p.instance.Start(":8888")
}