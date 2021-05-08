package datautil

import "strings"

type TargetValues []string

func NewTargetValues(sign string, targets ...string) TargetValues {
	v := TargetValues{}
	v = append(v, sign)
	v = append(v, targets...)

	return v
}

func (t TargetValues) Match(value string) bool {
	value = strings.ToLower(value)
	for _, v := range t {
		if strings.ToLower(v) == value {
			return true
		}
	}
	return false
}

func (t TargetValues) GetSign() string {
	return t[0]
}
