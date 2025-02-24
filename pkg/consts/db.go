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

package consts

import "os"

func IndexBucketName() []byte { return []byte("index-v0") }
func TorrentBucket() []byte   { return []byte("torrent-v0") }
func NodeBucketName() []byte  { return []byte("node-v0") }
func BlockBucketName() []byte { return []byte("block-v0") }

const (
	DefaultFilePerm  os.FileMode = 0640
	DefaultDirPerm               = os.ModeDir | DefaultFilePerm
	SecurityPerm     os.FileMode = 0600
	PrivateKeyLength             = 4096
	IPFSBlockDB                  = "ipfs.block.bolt"
)
