package appchaininteract

type app struct {
	AppTitle string `json:"desc"`
	Meta meta `json:"meta"`

}

type meta struct {
	Detial detial `json:"detail_1"`
}

type detial struct {
	Title string `json:"desc"`
	Icon string `json:"icon"`
	Preview string `json:"preview"`
	QqDocURL string `json:"qqdocurl"`
 }