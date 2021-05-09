package appchaininteract

import (
	"encoding/json"
	"regexp"
	"strings"
)

var perfix = regexp.MustCompile(`^https?://`)

func removeURLParams(URL string) string {
	rootURL := strings.Split(URL, "?")[0]
	transedURL := strings.Replace(rootURL, `\/`, "/", -1)
	if !perfix.MatchString(transedURL) {
		transedURL = "http://" + transedURL
	}
	return transedURL
}

func jsonLoader(jsonData string) (app, error) {
	decoder := json.NewDecoder(strings.NewReader(jsonData))

	var a app
	err := decoder.Decode(&a)

	return a, err
}
