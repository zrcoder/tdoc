package docmgr

import (
	"github.com/zrcoder/tdoc/model"
)

type Manager interface {
	Docs() []*model.DocInfo
	Sort(Less)
}

type Less func(*model.DocInfo, *model.DocInfo) bool

var ByModTime Less = func(a, b *model.DocInfo) bool {
	return a.ModTime.Before(b.ModTime)
}

var ByTitle Less = func(a, b *model.DocInfo) bool {
	return a.Name < b.Name
}
