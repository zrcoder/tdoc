package view

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type altViewport struct {
	viewport.Model

	ready bool
}

func (av *altViewport) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case menuSizeMsg:
		av.setSize(msg.Width, msg.Height)
	case docSizeMsg:
		av.setSize(msg.Width, msg.Height)
	case tea.WindowSizeMsg:
		av.setSize(msg.Width, msg.Height)
	}
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
