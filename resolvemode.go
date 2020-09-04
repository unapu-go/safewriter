package safewriter

import (
	"os"

	path_helpers "github.com/moisespsena-go/path-helpers"
	"github.com/pkg/errors"
)

func ModeOf(pth string) (mode os.FileMode, err error) {
	if mode, err = path_helpers.ResolveFileMode(pth); err != nil {
		err = errors.Wrapf(err, "resolve mode of %q", pth)
	}
	return
}
