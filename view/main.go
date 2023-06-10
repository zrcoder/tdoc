package view

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	help       help.Model
	menu       *Menu
	doc        *Doc
	mainHeight int
}

func NewModel(docs []*model.DocInfo) *Model {
	model := &Model{}
	model.menu = NewMenu(docs)
	model.doc = NewDoc(docs[0])

	model.help = help.New()
	return model
}

func (m *Model) Init() tea.Cmd {
	m.help.ShowAll = true
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		msg.Height -= 5 // 5 lines for key help view
		cmds = append(cmds, menuSizeCmd(msg), docSizeCmd(msg))
		m.mainHeight = msg.Height
	}
	_, cmd1 := m.menu.Update(msg)
	_, cmd2 := m.doc.Update(msg)
	cmds = append(cmds, cmd1, cmd2)
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	main := lipgloss.JoinHorizontal(lipgloss.Top, m.menu.View(), m.sepView(), m.doc.View())
	return lipgloss.JoinVertical(lipgloss.Left, "\n", main, m.help.View(m))
}

func (m *Model) FullHelp() [][]key.Binding {
	mk := m.menu.list.KeyMap
	dk := m.doc.KeyMap
	return [][]key.Binding{
		{mk.CursorUp, mk.CursorDown},
		{mk.PrevPage, mk.NextPage},
		{dk.HalfPageUp, dk.HalfPageDown},
		{mk.Quit},
	}
}

func (m *Model) ShortHelp() []key.Binding {
	return nil
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
	return strings.Repeat(Sep, m.mainHeight)
}
