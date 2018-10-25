package icq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const (
	apiRootURL string = "https://botapi.icq.net"
	botName           = "icqAlertGateBot"
	botNick           = "icqAlertGateBot"
	botUin            = "999999"
	botVersion        = "0.0.1"
	libVersion string = "0.0.1"
)

// Bot is singleton with bot API params
type Bot struct {
	RootURL    string
	Token      string
	Name       string
	Nick       string
	Uin        string
	Version    string
	LibVersion string
	userAgent  string
}

func (bot *Bot) Init() error {
	if bot.Token == "" {
		return errors.New("token shod't be empty")
	}
	if bot.RootURL == "" {
		bot.RootURL = apiRootURL
	}
	if bot.Name == "" {
		bot.Name = botName
	}
	if bot.Nick == "" {
		bot.Nick = botNick
	}
	if bot.Uin == "" {
		bot.Uin = botUin
	}
	if bot.Version == "" {
		bot.Version = botVersion
	}
	if bot.LibVersion == "" {
		bot.LibVersion = libVersion
	}
	bot.userAgent = fmt.Sprintf("%s/%s (uin=%s; nick=%s) gocq/%s",
		bot.Name, bot.Version, bot.Uin, bot.Nick, bot.LibVersion)
	return nil
}

func (bot *Bot) postText(path string, params url.Values, result interface{}) error {
	uuid, err := newUUID()
	if err != nil {
		return err
	}
	uri := bot.RootURL + path
	params["r"] = []string{uuid}
	params["aimsid"] = []string{bot.Token}
	request, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader([]byte(params.Encode())))
	if err != nil {
		return err
	}
	request.Header.Set("User-Agent", bot.userAgent)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(responseData, result)
	if err != nil {
		return err
	}
	return nil
}

func (bot *Bot) postFile(filePath string, result interface{}) error {
	var fileData []byte
	params := url.Values{
		"aimsid":   []string{bot.Token},
		"filename": []string{filepath.Base(filePath)},
	}
	uri := bot.RootURL + "/im/sendFile?" + params.Encode()
	//
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Read(fileData)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(fileData))
	if err != nil {
		return err
	}
	request.Header.Set("User-Agent", bot.userAgent)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(responseData, result)
	if err != nil {
		return err
	}
	return nil
}
