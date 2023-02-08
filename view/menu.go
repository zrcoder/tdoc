package view

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/zrcoder/tdoc/model"
)

type Menu struct {
	altViewport
	docs    []*model.DocInfo
	current int
}

func NewMenu(docs []*model.DocInfo) *Menu {
	return &Menu{docs: docs, current: 0}
}

func (m *Menu) Init() tea.Cmd {
	return nil
}

func (m *Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+u":
			m.altViewport.HalfViewUp()
		case "ctrl+d":
			m.altViewport.HalfViewDown()
		case "k", "up":
			m.current = (m.current + len(m.docs) - 1) % len(m.docs)
			return m, m.UpdateDoc()
		case "j", "down":
			m.current = (m.current + 1) % len(m.docs)
			return m, m.UpdateDoc()
		}
	case menuSizeMsg:
		m.current = 0
		m.altViewport.Update(msg)
		return m, m.UpdateDoc()
	}
	return m, nil
}

func (m *Menu) View() string {
	content := m.renderedContent()
	m.altViewport.SetContent(content)
	return m.altViewport.View()
}

func (m *Menu) UpdateDoc() tea.Cmd {
	return func() tea.Msg {
		return docMsg(m.docs[m.current])
	}
}

func (m *Menu) renderedContent() string {
	const (
		selectedPrefix = "> "
		normalPrefix   = "  "
		prefixLen      = 2
		dotsSuffix     = " ..."
	)
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	buf := strings.Builder{}
	buf.WriteString("\n")
	for i, v := range m.docs {
		title := v.Name
		renderedWidth := lipgloss.Width(title) + prefixLen
		exraLen := prefixLen + len(dotsSuffix)
		if renderedWidth > m.altViewport.Width && m.altViewport.Width > exraLen {
			title = title[:m.altViewport.Width-exraLen] + dotsSuffix
		}
		if i == m.current {
			buf.WriteString(selectedStyle.Width(m.altViewport.Width).Render(selectedPrefix + title))
		} else {
			buf.WriteString(normalPrefix + title)
		}
		buf.WriteString("\n")
	}
	return buf.String()
}
