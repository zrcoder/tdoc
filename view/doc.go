package view

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"

	"github.com/zrcoder/tdoc/model"
)

type Doc struct {
	viewport.Model

	Doc *model.Doc
}

func NewDoc(doc *model.Doc) *Doc {
	return &Doc{Doc: doc}
}

func (d *Doc) Init() tea.Cmd {
	return nil
}

func (d *Doc) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		if key == "u" {
			d.Model.HalfViewUp()
		} else if key == "d" {
			d.Model.HalfViewDown()
		}
	case docMsg:
		d.Doc = msg
		d.GotoTop()
	}
	return d, nil
}

func (d *Doc) View() string {
	content, err := d.renderedContent()
	if err != nil {
		return err.Error()
	}
	d.Model.SetContent(content)
	return d.Model.View()
}

func (d *Doc) renderedContent() (string, error) {
	render, err := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(d.Model.Width))
	if err != nil {
		return "", err
	}
	return render.Render(string(d.Doc.Content))
}
