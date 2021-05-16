package xmlchaininteract

type msg struct {
	TemplateID int    `xml:"templateID,attr"`
	URL        string `xml:"url,attr"`

	ServiceID int `xml:"serviceID,attr"`

	Action       string `xml:"action,attr"`
	ActionData   string `xml:"actionData,attr"`
	A_actionData string `xml:"a_actionData,attr"`
	I_actionData string `xml:"i_actionData,attr"`

	Brief string `xml:"brief,attr"`
	Flag  int    `xml:"flag,attr"`

	Item   item   `xml:"item"`
	Source source `xml:"source"`
}

type item struct {
	Layout int `xml:"layout,attr"`

	Picture picture `xml:"picture"`

	Title   string `xml:"title"`
	Summary string `xml:"summary"`
}

type picture struct {
	Cover string `xml:"cover,attr"`
}

type source struct {
	URL  string `xml:"url,attr"`
	Icon string `xml:"icon,attr"`
	Name string `xml:"name,attr"`

	Appid string `xml:"appid,attr"`

	Action       string `xml:"action,attr"`
	ActionData   string `xml:"actionData,attr"`
	A_actionData string `xml:"a_actionData,attr"`
	I_actionData string `xml:"i_actionData,attr"`
}