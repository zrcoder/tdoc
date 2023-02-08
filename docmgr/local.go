package docmgr

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/zrcoder/tdoc/model"
)

const (
	mdExtension = ".md"
)

type local struct {
	dir  string
	docs []*model.DocInfo
}

func New(dir string) (Manager, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	var getter model.Getter = func(filename string) ([]byte, error) {
		return os.ReadFile(filepath.Join(dir, filename))
	}
	mgr := &local{dir: dir}
	err = filepath.WalkDir(dir, func(path string, de fs.DirEntry, errin error) error {
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

func (d *local) Docs() []*model.DocInfo {
	return d.docs
}

func (d *local) Sort(less Less) {
	sortDocs(d.docs, less)
}
