package view

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/zrcoder/tdoc/model"
)

var (
	ErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f00"))
)

const (
	MenuWidth = 34
	SepSpaces = 5
	Sep       = "  |  \n"
)

type (
	docMsg      *model.DocInfo
	menuSizeMsg tea.WindowSizeMsg
	docSizeMsg  tea.WindowSizeMsg
)

type Getter func(string) ([]byte, error)

func Run(docs []*model.DocInfo) error {
	m := NewModel(docs)
	_, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion()).Run()
	return err
}

type Model struct {
	menue  *Menu
	doc    *Doc
	height int
}

func NewModel(docs []*model.DocInfo) *Model {
	model := &Model{}
	model.menue = NewMenu(docs)
	model.doc = NewDoc(docs[0])
	return model
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		cmds = append(cmds, menuSizeCmd(msg), docSizeCmd(msg))
		m.height = msg.Height
	}
	_, cmd1 := m.menue.Update(msg)
	_, cmd2 := m.doc.Update(msg)
	cmds = append(cmds, cmd1, cmd2)
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.menue.View(), m.sepView(), m.doc.View())
}

func menuSizeCmd(msg tea.WindowSizeMsg) tea.Cmd {
	return func() tea.Msg {
		return menuSizeMsg{Width: MenuWidth, Height: msg.Height}
	}
}

func docSizeCmd(msg tea.WindowSizeMsg) tea.Cmd {
	return func() tea.Msg {
		return docSizeMsg{Width: msg.Width - MenuWidth - SepSpaces, Height: msg.Height}
	}
}

func (m *Model) sepView() string {
	return strings.Repeat(Sep, m.height)
}
