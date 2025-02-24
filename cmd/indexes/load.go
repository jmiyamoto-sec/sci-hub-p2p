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

package indexes

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/itchio/lzma"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"

	"sci_hub_p2p/internal/utils"
	"sci_hub_p2p/pkg/consts"
	"sci_hub_p2p/pkg/logger"
	"sci_hub_p2p/pkg/vars"
)

var loadCmd = &cobra.Command{
	Use:           "load",
	Short:         "Load indexes into database.",
	Example:       "indexes load /path/to/*.jsonlines.lzma [--glob '/path/to/data/*.jsonlines.lzma']",
	SilenceErrors: false,
	PreRunE:       utils.EnsureDir(vars.GetAppBaseDir()),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		args, err = utils.MergeGlob(args, glob)
		if err != nil {
			return errors.Wrap(err, "can't load any index files")
		}
		db, err := bbolt.Open(vars.IndexesBoltPath(), consts.DefaultFilePerm, bbolt.DefaultOptions)
		if err != nil {
			return errors.Wrap(err, "cant' open database file, maybe another process is running")
		}
		defer func(db *bbolt.DB) {
			if e := db.Close(); e != nil {
				e = errors.Wrap(e, "can't save data to disk")
				if err == nil {
					err = e
				} else {
					logger.Error("", zap.Error(e))
				}
			}
		}(db)

		var count int
		err = db.Batch(func(tx *bbolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists(consts.IndexBucketName())
			if err != nil {
				return errors.Wrap(err, "can't create bucket in database")
			}

			for _, file := range args {
				c, err := loadIndexFile(b, file)
				if err != nil {
					return errors.Wrap(err, "can't load indexes file "+file)
				}
				count += c

			}

			return nil
		})
		if err != nil {
			return errors.Wrap(err, "can't save torrent data to database")
		}
		fmt.Printf("successfully load %d index file into database\n", len(args))
		fmt.Printf("%d records\n", count)

		return nil
	},
}

var glob string

func init() {
	loadCmd.Flags().StringVar(&glob, "glob", "",
		"glob pattern to search indexes to avoid 'Argument list too long' error")
}

func loadIndexFile(b *bbolt.Bucket, name string) (success int, err error) {
	f, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	reader := lzma.NewReader(f)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		var s []string

		err = json.Unmarshal(scanner.Bytes(), &s)
		if err != nil || len(s) != 2 {
			return 0, errors.Wrap(err, "can't parse json "+scanner.Text())
		}

		value, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			return 0, errors.Wrap(err, "can't decode base64")
		}

		key, err := url.QueryUnescape(strings.TrimSuffix(s[0], ".pdf"))
		if err != nil {
			return 0, err
		}

		err = b.Put([]byte(key), value)
		if err != nil {
			return 0, errors.Wrap(err, "can't save record to database")
		}

		success++
	}

	err = scanner.Err()
	if err != nil {
		return 0, errors.Wrap(err, "can't scan file")
	}

	return success, errors.Wrap(scanner.Err(), "can't scan file")
}
