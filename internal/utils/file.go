// Copyright 2021 Trim21<trim21.me@gmail.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		_ = out.Close()

		return err
	}

	return out.Close()
}

func GlobWithExpand(glob string) ([]string, error) {
	if strings.Contains(glob, "~") {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return nil, errors.Wrap(err, "can't determine homedir to expand ~")
		}

		glob = strings.ReplaceAll(glob, "~", homedir)
	}

	s, err := filepath.Glob(glob)

	return s, err
}

var ErrNotAFile = errors.New("not a file")
var ErrNotADir = errors.New("not a dir")

func FileExist(name string) (bool, error) {
	s, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}

	if s.IsDir() {
		return false, ErrNotAFile
	}

	return true, err
}

func DirExist(name string) (bool, error) {
	s, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}

	if !s.IsDir() {
		return false, ErrNotADir
	}

	return true, err
}

func ReadFileAt(name string, offset, length int64) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file %s", name)
	}
	defer f.Close()

	var p = make([]byte, length)

	if _, err = f.ReadAt(p, offset); err != nil {
		return nil, errors.Wrapf(err, "failed to read file %s", name)
	}

	return p, nil
}
