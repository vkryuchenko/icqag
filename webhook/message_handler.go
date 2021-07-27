package webhook

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func (p *Provider) handleMessage(w http.ResponseWriter, r *http.Request) {
	payload, err := p.payloadBySourceName(chi.URLParam(r, "source"))
	if err != nil {
		p.Logger.Error(err.Error(), zap.String("url", r.RequestURI))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	messageString, err := payload.Parse(r, p.Logger)
	if err != nil {
		p.Logger.Error(err.Error(), zap.String("url", r.RequestURI))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	for _, msg := range wrapMessage(messageString, 4000) {
		message := p.Bot.NewTextMessage(chi.URLParam(r, "target"), msg)
		err = message.Send()
		if err != nil {
			p.Logger.Error(err.Error(), zap.String("url", r.RequestURI))
			http.Error(w, "", http.StatusBadRequest)
			return
		}
	}
}

func wrapMessage(message string, maxLen int) []string {
	runes := bytes.Runes([]byte(message))
	runesCount := len(runes)
	if runesCount <= maxLen {
		return []string{message}
	}
	var result []string
	for i := 0; i < runesCount; i += maxLen {
		var s string
		start := i
		stop := i + maxLen
		if stop > runesCount {
			stop = runesCount
		}
		for _, r := range runes[start:stop] {
			s += string(r)
		}
		result = append(result, s)
	}
	return result
}
