package safewriter

import (
	"io"
	"os"
)

func Create(cb func(w io.Writer) (n int, err error), path string, mode ...os.FileMode) (n int, err error) {
	var w *Writer
	if w, err = Open(path, mode...); err != nil {
		return
	}
	defer func() {
		if err == nil {
			err = w.Close()
		} else {
			w.Close()
		}
	}()
	return cb(w)
}
