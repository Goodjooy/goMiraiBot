package xmlchaininteract

import (
	"encoding/xml"
	"regexp"
	"strings"
)

var headPattern = regexp.MustCompile(`(?:\?>)(.+)`)

func loadXML(xmlSource string) (msg, error) {
	xmlBodys:= headPattern.FindStringSubmatch(xmlSource)
	xmlBody:=xmlBodys[1]

	var m msg
	decoder := xml.NewDecoder(strings.NewReader(xmlBody))
	err := decoder.Decode(&m)

	return m, err
}

func removeURLParams(URL string)string{
	rootURL:= strings.Split(URL,"?")[0]
	transedURL:=strings.Replace(rootURL,`\/`,"/",-1)
	return transedURL
}