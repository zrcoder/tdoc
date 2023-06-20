package view

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

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
		desc := d.Description
		if desc == "" {
			t := d.ModTime
			if t.IsZero() {
				t = time.Now()
			}
			desc = t.Format("2006-01-02 15:04")
		}
		items[i] = item{
			title: d.Title,
			desc:  desc,
		}
	}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = title // can be ""
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

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
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }
