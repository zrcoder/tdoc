package model

import "time"

type Getter func(string) ([]byte, error)

type DocInfo struct {
	Name        string
	Title       string
	Description string
	ModTime     time.Time
	Getter      Getter
	cached      []byte
}

func (di *DocInfo) Get() ([]byte, error) {
	if len(di.cached) == 0 {
		return di.Getter(di.Name)
	}
	return di.cached, nil
}
