package jenkins

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

//{
//
//"buildName": "#5",
//"buildUrl": "https://jenkins.example.com/job/Android/15/",
//"error": {
//"m.rublev": [
//"Assets/Scripts/Game/Tournament/TournamentPositionReward.cs(8,0): error CS1525: Unexpected symbol '}'"
//]
//},
//"event": "DONE",
//"projectName": "Android"
//}

var (
	emailRegexp = regexp.MustCompile(`<(?P<email>[\w.]+@(corp\.)?(mail|list|bk|inbox)\.ru)>$`)
)

// SuspectMessage represent data struct by Jenkins Outbound WebHook plugin
type SuspectMessage struct {
	BuildName   string              `json:"buildName"`
	BuildURL    string              `json:"buildUrl"`
	Event       string              `json:"event"`
	ProjectName string              `json:"projectName"`
	Branch      string              `json:"branch"`
	Error       map[string][]string `json:"error"`
}

func (*SuspectMessage) transform(data io.ReadCloser) (string, error) {
	messageBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return "", err
	}
	sm := SuspectMessage{}
	err = json.Unmarshal(messageBytes, &sm)
	if err != nil {
		return "", err
	}
	lines := []string{
		"Status: " + strings.ToUpper(sm.Event),
		"Name: " + sm.ProjectName,
		"Branch: " + sm.Branch,
		"URL: " + sm.BuildURL,
	}
	for author, errors := range sm.Error {
		suspect := author
		if emailRegexp.MatchString(author) {
			matchResult := emailRegexp.FindStringSubmatch(author)
			currentMention := matchResult[1]
			suspect = "@[" + currentMention + "]"
		}
		lines = append(lines, "\nSUSPECT: "+suspect+"\nERRORS:")
		for _, errLine := range errors {
			lines = append(lines, "### "+errLine)
		}
	}
	return strings.Join(lines, "\n"), nil
}

// Parse implement Payload.Parse()
func (m SuspectMessage) Parse(req *http.Request) (string, error) {
	return m.transform(req.Body)
}
