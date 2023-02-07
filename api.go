package tdoc

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/zrcoder/tdoc/model"
	"github.com/zrcoder/tdoc/view"
)

func ParseFromDir(dir string) ([]*model.Doc, error) {
	root, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	var res []*model.Doc
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, errin error) error {
		if errin != nil {
			return errin
		}
		if path == root {
			return nil
		}
		if d.IsDir() {
			return filepath.SkipDir
		}
		fi, err := d.Info()
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		res = append(res, &model.Doc{
			Title:   getTitle(d.Name()), // todo
			ModTime: fi.ModTime(),
			Content: data,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("invalid or empty directory: %s", dir)
	}
	return res, nil
}

func Sort(docs []*model.Doc, less func(*model.Doc, *model.Doc) bool) {
	sort.Slice(docs, func(i, j int) bool {
		return less(docs[i], docs[j])
	})
}

func ByModTime(a, b *model.Doc) bool {
	return a.ModTime.Before(b.ModTime)
}

func ByTitle(a, b *model.Doc) bool {
	return a.Title < b.Title
}

func Run(docs []*model.Doc) error {
	return view.Run(docs)
}

func getTitle(fileName string) string {
	i := strings.LastIndex(fileName, ".")
	if i < 0 {
		return fileName
	}
	return fileName[:i]
}
