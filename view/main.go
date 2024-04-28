package view

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/tdoc/model"
)

var ErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f00"))

const (
	HelpHeight = 4
)

type (
	docMsg  *model.DocInfo
	menuMsg struct{}
)

func Run(docs []*model.DocInfo, cfg ...model.Config) error {
	m := NewModel(docs, cfg...)
	_, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion()).Run()
	return err
}

type Model struct {
	isMenu bool
	help   help.Model
	menu   *Menu
	doc    *Doc
}

func NewModel(docs []*model.DocInfo, cfg ...model.Config) *Model {
	model := &Model{}
	title := ""
	if len(cfg) > 0 {
		title = cfg[0].Title
	}
	model.isMenu = len(docs) > 1
	model.menu = NewMenu(title, docs)
	model.doc = NewDoc(docs[0], model.menu.list.KeyMap.Quit)
	model.help = help.New()
	return model
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case docMsg:
		m.isMenu = false
	case menuMsg:
		m.isMenu = true
	}
	_, cmd1 := m.menu.Update(msg)
	_, cmd2 := m.doc.Update(msg)
	return m, tea.Batch(cmd1, cmd2)
}

func (m *Model) View() string {
	if m.isMenu {
		return m.menu.View()
	}
	return m.doc.View()
}
