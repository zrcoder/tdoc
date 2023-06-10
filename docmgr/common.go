package docmgr

import (
	"sort"

	"github.com/zrcoder/tdoc/model"
)

const (
	mdExtension = ".md"
)

func sortDocs(docs []*model.DocInfo, less Less) {
	sort.Slice(docs, func(i, j int) bool {
		return less(docs[i], docs[j])
	})
}

type common struct {
	docs []*model.DocInfo
}

func NewWithDocs(docs []*model.DocInfo) Manager {
	return &common{
		docs: docs,
	}
}

func (c *common) Docs() []*model.DocInfo {
	return c.docs
}

func (c *common) Sort(less Less) {
	sortDocs(c.docs, less)
}
