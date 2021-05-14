package inandoutinteract

import (
	"time"

	"gorm.io/gorm")


type User struct {
	gorm.Model

	QQId uint64 `gorm:"not null"`

	Records []PaymentRecord
}

type PaymentRecord struct {
	gorm.Model

	User User 
	UserId uint

	Date time.Time

	PaymentFloat uint16
	PaymentInt int32

	GoodType string `gorm:"size:2;default:￥"`

	Message string `gorm:"size:256"`
}