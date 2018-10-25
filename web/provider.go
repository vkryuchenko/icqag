package web

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"icqag/icq"
	"icqag/webhook"
	"icqag/webhook/alertmanager"
	"icqag/webhook/grafana"
	"icqag/webhook/jenkins"
	"icqag/webhook/raw"
)

var payloadSourceMap = map[string]webhook.Payload{
	"raw":              raw.Message{},
	"grafana":          grafana.Message{},
	"jenkins-outbound": jenkins.OutboundMessage{},
	"jenkins-suspect":  jenkins.SuspectMessage{},
	"alertmanager":     alertmanager.Message{},
}

// Provider represent single instances of bot and echo
type Provider struct {
	Bot      *icq.Bot
	instance *echo.Echo
}

func (Provider) payloadBySourceName(sourceName string) (webhook.Payload, error) {
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
