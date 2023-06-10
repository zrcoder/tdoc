package view

import (
	"bytes"

	tea "github.com/charmbracelet/bubbletea"
	mr "github.com/charmbracelet/glamour"

	"github.com/zrcoder/tdoc/model"
)

type Doc struct {
	altViewport
	Doc *model.DocInfo
}

func NewDoc(doc *model.DocInfo) *Doc {
	return &Doc{Doc: doc}
}

func (d *Doc) Init() tea.Cmd {
	return nil
}

func (d *Doc) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case docMsg:
		d.Doc = msg
		d.altViewport.GotoTop()
	case docSizeMsg:
		d.altViewport.Update(msg)
	}
	var cmd tea.Cmd
	d.Model, cmd = d.Model.Update(msg)
	d.disableUnusedKeys()
	return d, cmd
}

func (d *Doc) View() string {
	content, err := d.renderedContent()
	if err != nil {
		return ErrStyle.Copy().Render(err.Error())
	}
	d.altViewport.SetContent(string(content))
	return d.altViewport.View()
}

func (d *Doc) disableUnusedKeys() {
	d.altViewport.KeyMap.Down.SetEnabled(false)
	d.altViewport.KeyMap.Up.SetEnabled(false)
	d.altViewport.KeyMap.PageDown.SetEnabled(false)
	d.altViewport.KeyMap.PageUp.SetEnabled(false)
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
