package view

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/zrcoder/tdoc/model"
)

type Menu struct {
	list    list.Model
	docs    []*model.DocInfo
	width   int
	height  int
	current int
}

func NewMenu(title string, docs []*model.DocInfo) *Menu {
	items := make([]list.Item, len(docs))
	for i, d := range docs {
		items[i] = item(d.Title)
	}
	l := list.New(items, itemDelegate{}, 0, 0)
	if title != "" {
		l.Title = title
		l.Styles.Title = titleStyle
	} else {
		l.SetShowTitle(false)
	}
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.PaginationStyle = paginationStyle
	l.SetShowHelp(false)

	return &Menu{
		docs: docs,
		list: l,
	}
}

func (m *Menu) Init() tea.Cmd {
	return nil
}

func (m *Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var updateDocCmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			updateDocCmd = m.UpdateDoc()
		}
	case menuSizeMsg:
		m.current = 0
		m.width = msg.Width
		m.height = msg.Height
	}
	var updateListCmd tea.Cmd
	m.list, updateListCmd = m.list.Update(msg)
	m.list.SetSize(m.width, m.height)
	return m, tea.Batch(updateDocCmd, updateListCmd)
}

func (m *Menu) View() string {
	return "\n" + m.list.View()
}

func (m *Menu) UpdateDoc() tea.Cmd {
	return func() tea.Msg {
		return docMsg(m.docs[m.list.Index()])
	}
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
	if len(res) > MenuWidth {
		res = res[:MenuWidth]
	}
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(res))
}
