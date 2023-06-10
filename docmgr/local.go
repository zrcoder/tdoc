package docmgr

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/zrcoder/tdoc/model"
)

type local struct {
	dir  string
	docs []*model.DocInfo
}

func NewWithDirectory(dir string) (Manager, error) {
	getter := func(filename string) ([]byte, error) {
		return os.ReadFile(filepath.Join(dir, filename))
	}
	mgr := &local{dir: dir}
	err := filepath.WalkDir(dir, func(path string, de fs.DirEntry, errin error) error {
		if errin != nil {
			return errin
		}
		if path == dir {
			return nil
		}
		if de.IsDir() {
			return filepath.SkipDir
		}
		name := de.Name()
		if !strings.HasSuffix(name, mdExtension) {
			return nil
		}
		fi, err := de.Info()
		if err != nil {
			return err
		}
		mgr.docs = append(mgr.docs, &model.DocInfo{
			Name:    name,
			Title:   name[:len(name)-len(mdExtension)],
			ModTime: fi.ModTime(),
			Getter:  getter,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(mgr.docs) == 0 {
		return nil, fmt.Errorf("no markdown files found in directory: %s", dir)
	}
	return mgr, nil
}

func (l *local) Docs() []*model.DocInfo {
	return l.docs
}

func (l *local) Sort(less Less) {
	sortDocs(l.docs, less)
}
