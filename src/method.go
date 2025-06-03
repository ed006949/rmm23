package main

import (
	"encoding/xml"
	"path/filepath"

	"github.com/avfs/avfs"
	"github.com/avfs/avfs/vfs/memfs"

	"rmm23/src/io_vfs"
	"rmm23/src/l"
)

func (r *xmlConf) load() (err error) {
	var (
		vfsDB = &io_vfs.VFSDB{
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
		return l.ENOCONF
	}

	return
}
