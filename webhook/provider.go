package webhook

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/mail-ru-im/bot-golang"
	"icqag/alertsource/alertmanager"
	"icqag/alertsource/grafana"
	"icqag/alertsource/jenkins"
	"icqag/alertsource/json"
	"icqag/alertsource/raw"
	"icqag/alertsource/teamcity"
	"net/http"
	"os"
	"strings"
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
	Parse(req *http.Request, logger echo.Logger) (string, error)
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
	if strings.EqualFold(os.Getenv("DEBUG"), "true") {
		e.Debug = true
		e.Logger.SetLevel(log.DEBUG)
		e.Logger.Debug("debug mode enabled")
	}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//
	e.GET("/health", func(echo.Context) error { return nil })
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
