package view

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/zrcoder/tdoc/model"
)

type Menu struct {
	list    list.Model
	help    help.Model
	docs    []*model.DocInfo
	current int
}

func NewMenu(title string, docs []*model.DocInfo) *Menu {
	items := make([]list.Item, len(docs))
	for i, d := range docs {
		items[i] = item(d.Title)
	}
	l := list.New(items, itemDelegate{}, 0, 0)
	l.Title = title // can be ""
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.PaginationStyle = paginationStyle

	return &Menu{
		docs: docs,
		list: l,
		help: help.New(),
	}
}

func (m *Menu) Init() tea.Cmd {
	m.help.ShowAll = true
	return nil
}

func (m *Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var updateDocCmd, updateListCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			updateDocCmd = func() tea.Msg { return docMsg(m.docs[m.list.Index()]) }
		default:
			m.list, updateListCmd = m.list.Update(msg)
		}
	case tea.WindowSizeMsg:
		m.current = 0
		m.list.SetSize(msg.Width, msg.Height)
	}
	m.help, _ = m.help.Update(msg)
	return m, tea.Batch(updateDocCmd, updateListCmd)
}

func (m *Menu) View() string {
	return m.list.View()
}

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
)

type item string

func (item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	itm, ok := listItem.(item)
	if !ok {
		return
	}
	res := string(itm)
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(res))
}
