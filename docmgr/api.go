package docmgr

import (
	"sort"

	"github.com/zrcoder/tdoc/model"
)

type Manager interface {
	Docs() []*model.DocInfo
	Sort(Less)
}

type Less func(*model.DocInfo, *model.DocInfo) bool

var (
	ByModTime Less = func(a, b *model.DocInfo) bool {
		return a.ModTime.Before(b.ModTime)
	}

	ByTitle Less = func(a, b *model.DocInfo) bool {
		return a.Name < b.Name
	}
)

func sortDocs(docs []*model.DocInfo, less Less) {
	sort.Slice(docs, func(i, j int) bool {
		return less(docs[i], docs[j])
	})
}
