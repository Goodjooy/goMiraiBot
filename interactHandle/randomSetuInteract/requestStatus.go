package randomsetuinteract

import "gorm.io/gorm"

type setuRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`

	Quota            uint16 `json:"quota"`
	QuotaRecoverTime uint   `json:"quota_min_ttl"`

	Count uint `json:"count"`

	Setus []setuMeta `json:"data"`
}

type setuMeta struct {
	Pid    uint64 `json:"pid"`
	PicNum uint   `json:"p"`

	AuthorId uint64 `json:"uid"`

	Title  string `json:"title"`
	Author string `json:"author"`

	URL string `json:"url"`

	R18 bool `json:"r18"`

	Width  uint16 `json:"width"`
	Height uint16 `json:"height"`

	Tags []string `json:"tags"`
}

type setuInfo struct {
	gorm.Model

	URL string `gorm:"size:512"`
	R18 bool 
}