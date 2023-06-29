package view

import (
	"bytes"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	mr "github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/zrcoder/tdoc/model"
)

type Doc struct {
	altViewport
	quitKey key.Binding
	help    help.Model
	Doc     *model.DocInfo
}

func NewDoc(doc *model.DocInfo, quitKey key.Binding) *Doc {
	return &Doc{Doc: doc, help: help.New(), quitKey: quitKey}
}

func (d *Doc) Init() tea.Cmd {
	return nil
}

func (d *Doc) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case docMsg:
		d.Doc = msg
		d.altViewport.GotoTop()
	case tea.WindowSizeMsg:
		d.altViewport.setSize(msg.Width, msg.Height-HelpHeight)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, d.KeyMap.HalfPageUp):
			d.altViewport.HalfViewUp()
		case key.Matches(msg, d.KeyMap.HalfPageDown):
			d.altViewport.HalfViewDown()
		case key.Matches(msg, menuKey):
			return d, func() tea.Msg { return menuMsg{} }
		}
	}
	d.help, _ = d.help.Update(msg)
	return d, nil
}

func (d *Doc) View() string {
	content, err := d.renderedContent()
	if err != nil {
		return ErrStyle.Copy().Render(err.Error())
	}
	str := lipgloss.NewStyle().Width(d.altViewport.Width).Height(d.altViewport.Height).Render(string(content))
	d.altViewport.SetContent(str)
	return lipgloss.JoinVertical(lipgloss.Left, d.altViewport.View(), "\n", "  "+d.help.View(d))
}

func (d *Doc) renderedContent() ([]byte, error) {
	content, err := d.Doc.Get()
	if err != nil {
		return nil, err
	}
	// workaround for glamar's bug
	content = bytes.ReplaceAll(content, []byte{'\t'}, []byte("    "))
	render, err := mr.NewTermRenderer(mr.WithAutoStyle(), mr.WithWordWrap(d.altViewport.Width))
	if err != nil {
		return nil, err
	}
	return render.RenderBytes(content)
}

var menuKey = key.NewBinding(
	key.WithKeys("m"),
	key.WithHelp("m", "show menu"),
)

func (d *Doc) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{d.KeyMap.HalfPageUp, d.KeyMap.HalfPageDown},
		{menuKey, d.quitKey},
	}
}

func (d *Doc) ShortHelp() []key.Binding {
	return []key.Binding{d.KeyMap.HalfPageUp, d.KeyMap.HalfPageDown, menuKey, d.quitKey}
}
