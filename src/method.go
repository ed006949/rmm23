package main

import (
	"encoding/xml"
	"path/filepath"

	"github.com/avfs/avfs"
	"github.com/avfs/avfs/vfs/memfs"

	"rmm23/src/l"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_vfs"
)

func (r *xmlConf) load() (err error) {
	var (
		vfsDB = &mod_vfs.VFSDB{
			List: make(map[string]string),
			VFS: memfs.NewWithOptions(&memfs.Options{
				Idm:        avfs.NewDummyIdm(),
				User:       nil,
				Name:       "",
				OSType:     avfs.CurrentOSType(),
				SystemDirs: nil,
			}),
		}
	)
	l.InitCLI()

	switch {
	case len(l.Config.String()) != 0: //
		l.CLI.Set()

		var (
			cliConfigFile string
			data          []byte
		)

		switch cliConfigFile, err = filepath.Abs(l.Config.String()); {
		case err != nil:
			return
		}
		switch err = vfsDB.CopyFromFS(cliConfigFile); {
		case err != nil:
			return
		}
		switch data, err = vfsDB.VFS.ReadFile(cliConfigFile); {
		case err != nil:
			return
		}
		switch err = xml.Unmarshal(data, r); {
		case err != nil:
			return
		}

	default:
		return mod_errors.ENOCONF
	}

	return
}
