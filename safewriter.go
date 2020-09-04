package safewriter

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Writer struct {
	path string
	mode os.FileMode
	*os.File
}

func (this Writer) Path() string {
	return this.path
}

func Open(path string, mode ...os.FileMode) (w *Writer, err error) {
	var mode_ os.FileMode
	for _, mode_ = range mode {
	}

	if mode_ == 0 {
		if mode_, err = ModeOf(path); err != nil {
			return
		}
	}

	var f *os.File
	if f, err = ioutil.TempFile(filepath.Dir(path), filepath.Base(path)); err != nil {
		return
	}

	return &Writer{path: path, mode: mode_, File: f}, nil
}

func (this Writer) Close() (err error) {
	if this.File == nil {
		return io.ErrClosedPipe
	}
	defer func() {
		if err == nil {
			err = errors.Wrapf(this.File.Close(), "close %q", this.File.Name())
		} else {
			this.File.Close()
		}
		if err == nil {
			err = errors.Wrapf(os.Remove(this.File.Name()), "remove %q", this.File.Name())
		} else {
			os.Remove(this.File.Name())
		}
		this.File = nil
	}()
	if _, err = os.Stat(this.path); err != nil {
		if os.IsNotExist(err) {
			err = nil
		} else {
			return
		}
	}

	var f *os.File
	if f, err = os.OpenFile(this.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, this.mode); err != nil {
		err = errors.Wrapf(err, "open %q for create or truncate", this.path)
		return
	}

	if _, err = this.File.Seek(0, io.SeekStart); err != nil {
		err = errors.Wrapf(err, "seek %q to 0", this.File.Name())
		return
	}

	if _, err = io.Copy(f, this.File); err != nil {
		err = errors.Wrapf(err, "copy %q to %q", this.File.Name(), this.path)
	}
	return
}
