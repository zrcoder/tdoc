package tdoc

import (
	"github.com/zrcoder/tdoc/model"
	"github.com/zrcoder/tdoc/view"
)

func Run(docs []*model.DocInfo) error {
	return view.Run(docs)
}
