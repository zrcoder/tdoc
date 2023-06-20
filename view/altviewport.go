package view

import (
	"github.com/charmbracelet/bubbles/viewport"
)

type altViewport struct {
	viewport.Model

	ready bool
}

func (av *altViewport) setSize(width, height int) {
	if av.ready {
		av.Width = width
		av.Height = height
	} else {
		av.Model = viewport.New(width, height)
		av.ready = true
	}
}
