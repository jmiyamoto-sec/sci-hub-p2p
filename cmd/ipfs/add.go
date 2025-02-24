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

package ipfs

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"

	"sci_hub_p2p/internal/utils"
	"sci_hub_p2p/pkg/consts"
	"sci_hub_p2p/pkg/dag"
	"sci_hub_p2p/pkg/logger"
	"sci_hub_p2p/pkg/vars"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "add all files in a zip files",
	PreRunE: utils.EnsureDir(vars.GetAppBaseDir()),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if recursive && glob != "" {
			return errors.New("can't use --glob with --recursive")
		}
		if recursive {
			var zipFiles []string
			for _, arg := range args {
				err := filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
					if strings.HasSuffix(strings.ToLower(path), ".zip") {
						zipFiles = append(zipFiles, path)
					}

					return nil
				})
				if err != nil {
					return errors.Wrap(err, "can't search dir in args")
				}
			}
			if len(zipFiles) != 0 {
				fmt.Printf("find %d zip files to add\n", len(zipFiles))
			}
			args = zipFiles
		} else {
			args, err = utils.MergeGlob(args, glob)
			if err != nil {
				return errors.Wrap(err, "failed to find zip files")
			}
		}

		logger.Info("open database", zap.String("db", vars.IpfsDBPath()))
		db, err := bbolt.Open(vars.IpfsDBPath(), consts.DefaultFilePerm, &bbolt.Options{NoSync: true})
		if err != nil {
			return errors.Wrap(err, "failed to open database")
		}
		defer func(db *bbolt.DB) {
			err := db.Close()
			if err != nil {
				logger.Error("failed to close DataBase", zap.Error(err))
			}
		}(db)
		err = dag.InitDB(db)
		if err != nil {
			return errors.Wrap(err, "failed to initialize database")
		}

		width := len(strconv.Itoa(len(args)))

		for i, file := range args {
			fmt.Printf("%0*d/%d processing %s\n", width, i+1, len(args), file)
			if err := dag.AddZip(db, file); err != nil {
				logger.Error("failed to add files from zip archive", zap.Error(err))
			}

			if i%10 == 0 {
				err := db.Sync()
				if err != nil {
					logger.Error("failed to sync database to DB", zap.Error(err))
				}
			}
		}

		return errors.Wrap(db.Sync(), "failed to flush data to disk")
	},
}

var glob string
var recursive bool

func init() {
	addCmd.Flags().StringVar(&glob, "glob", "", "glob pattern")
	addCmd.Flags().BoolVarP(&recursive, "", "r", false, "recursively search all sub directory")
}
