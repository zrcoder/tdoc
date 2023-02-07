package view

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/zrcoder/tdoc/model"
)

const (
	MenuWidth = 40
	SepSpaces = 5
	Sep       = "  |  \n"
)

type docMsg *model.Doc

func Run(docs []*model.Doc) error {
	m := NewModel(docs)
	_, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion()).Run()
	return err
}

type Model struct {
	menue  *Menu
	doc    *Doc
	ready  bool
	height int
}

func NewModel(docs []*model.Doc) *Model {
	model := &Model{}
	model.menue = NewMenu(docs)
	model.doc = NewDoc(docs[0])
	return model
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			m.menue.Model = viewport.New(MenuWidth, msg.Height)
			m.doc.Model = viewport.New(msg.Width-MenuWidth-SepSpaces, msg.Height)
			m.ready = true
		} else {
			m.menue.Model.Width = MenuWidth
			m.menue.Model.Height = msg.Height
			m.doc.Model.Width = msg.Width - MenuWidth - SepSpaces
			m.doc.Model.Height = msg.Height
		}
		m.height = msg.Height
	}
	_, cmd1 := m.menue.Update(msg)
	_, cmd2 := m.doc.Update(msg)
	return m, tea.Batch(cmd1, cmd2)
}

func (m *Model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.menue.View(), m.sepView(), m.doc.View())
}

func (m *Model) sepView() string {
	return strings.Repeat(Sep, m.height)
}
