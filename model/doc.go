package model

import "time"

type Doc struct {
	Title   string
	ModTime time.Time
	Content []byte
}
