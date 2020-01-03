package model

import (
	"github.com/AdrianOrlow/links-api/app/utils"
	"time"
)

type Link struct {
	Model
	Url        string      `json:"url"`
	Permalink  string      `json:"permalink"`
	Visits     int         `json:"visits"`
	ValidUntil *time.Time  `json:"validUntil"`
}

func (l *Link) WithHashId() *Link {
	l.HashID = utils.EncodeId(int(l.ID))
	return l
}

func (l *Link) IsValid() bool {
	isNil := l.ValidUntil == nil
	dateExceeded := !isNil && l.ValidUntil.After(time.Now())
	return isNil || dateExceeded
}