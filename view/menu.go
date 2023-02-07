package view

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"

	"github.com/zrcoder/tdoc/model"
)

type Menu struct {
	viewport.Model

	docs    []*model.Doc
	current int
}

func NewMenu(docs []*model.Doc) *Menu {
	return &Menu{docs: docs}
}

func (m *Menu) Init() tea.Cmd {
	return nil
}

func (m *Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if ok {
		switch keyMsg.String() {
		case "ctrl+u":
			m.Model.HalfViewUp()
		case "ctrl+d":
			m.Model.HalfViewDown()
		case "k":
			m.current = (m.current + len(m.docs) - 1) % len(m.docs)
			return m, m.UpdateDoc()
		case "j":
			m.current = (m.current + 1) % len(m.docs)
			return m, m.UpdateDoc()
		}
	}
	return m, nil
}

func (m *Menu) View() string {
	s, err := m.renderedContent()
	if err != nil {
		return err.Error()
	}
	m.Model.SetContent(s)
	return m.Model.View()
}

func (m *Menu) UpdateDoc() tea.Cmd {
	return func() tea.Msg {
		return docMsg(m.docs[m.current])
	}
}

func (m *Menu) renderedContent() (string, error) {
	buf := strings.Builder{}
	for i, v := range m.docs {
		prefix := "- [ ] ["
		if i == m.current {
			prefix = "- [x] ["
		}
		buf.WriteString(prefix)
		buf.WriteString(v.Title)
		buf.WriteString("]()\n")
	}
	render, err := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(m.Model.Width))
	if err != nil {
		return "", err
	}
	return render.Render(buf.String())
}
