package icq

import "errors"

type apiResponse struct {
	Response struct {
		StatusCode       uint   `json:"statusCode"`
		StatusDetailCode uint   `json:"statusDetailCode"`
		StatusText       string `json:"statusText"`
		RequestID        string `json:"requestId"`
		Data             struct {
			SubCode struct {
				Error  uint   `json:"error"`
				Reason string `json:"reason"`
			} `json:"subCode"`
			MsgID string `json:"msgId"`
			State string `json:"state"`
		} `json:"data"`
	} `json:"response"`
}

func (r *apiResponse) msgError() error {
	if r.Response.StatusCode != 200 {
		return errors.New(r.Response.StatusText)
	}
	return nil
}

func (r *apiResponse) msgID() string {
	return r.Response.Data.MsgID
}
