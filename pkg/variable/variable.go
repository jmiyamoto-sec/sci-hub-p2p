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

package variable

import (
	"os"
	"path/filepath"
)

var appBaseDir string

func GetAppBaseDir() string {
	if appBaseDir != "" {
		return appBaseDir
	}
	appBaseDir = filepath.Clean(os.ExpandEnv("$HOME/.sci-hub-p2p"))

	return appBaseDir
}

func GetAppTmpDir() string {
	return filepath.Join(GetAppBaseDir(), "tmp")
}

func GetPaperBoltPath() string {
	return filepath.Join(GetAppBaseDir(), "papers.bolt")
}

func GetTorrentStoragePath() string {
	return filepath.Join(GetAppBaseDir(), "torrents")
}

func NodeBucketName() []byte {
	return []byte("node-v0")
}

func BlockBucketName() []byte {
	return []byte("block-v0")
}
