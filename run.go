package tdoc

import (
	"github.com/zrcoder/tdoc/model"
	"github.com/zrcoder/tdoc/view"
)

func Run(docs []*model.DocInfo, cfg ...model.Config) error {
	return view.Run(docs, cfg...)
}
