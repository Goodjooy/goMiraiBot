package xmlchaininteract

import (
	"encoding/xml"
	"net/url"
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
	//rootURL:=URL
	URL=idCapture(URL)
	transedURL:=strings.Replace(URL,`\/`,"/",-1)
	return transedURL
}

func idCapture(URL string)string{
	var reValue url.Values=make(url.Values)
	data,_:=url.Parse(URL)

	params:=data.RawQuery

	values,_:=url.ParseQuery(params)

	if v,ok:=values["id"];ok{
		reValue["id"]=v
	}
	data.RawQuery=reValue.Encode()

return data.String()
}