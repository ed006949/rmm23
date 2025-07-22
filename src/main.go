package main

import (
	"errors"
	"flag"

	"github.com/avfs/avfs"
	"github.com/avfs/avfs/vfs/memfs"

	"rmm23/src/l"
	"rmm23/src/mod_db"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_vfs"
)

func main() {
	l.Z{l.M: "main", "daemon": l.Name.String(), "commit": l.GitCommit.String()}.Informational()
	defer l.Z{l.M: "exit", "daemon": l.Name.String()}.Informational()

	var (
		err       error
		xmlConfig = new(xmlConf)
		vfsDB     = &mod_vfs.VFSDB{
			List: make(map[string]string),
			VFS: memfs.NewWithOptions(&memfs.Options{
				Idm:        avfs.NotImplementedIdm,
				User:       nil,
				Name:       "",
				OSType:     avfs.CurrentOSType(),
				SystemDirs: nil,
			}),
		}
	)

	switch err = xmlConfig.load(); {
	case errors.Is(err, mod_errors.ENOCONF):
		flag.PrintDefaults()
		l.Z{l.E: err}.Critical()
	case err != nil:
		flag.PrintDefaults()
		l.Z{l.E: err}.Critical()
	}

	switch err = vfsDB.CopyFromFS("./etc/legacy/"); {
	case err != nil:
		l.Z{l.E: err}.Critical()
	}

	switch err = mod_db.CopyLDAP2DB(ctx, xmlConfig.LDAP); {
	case err != nil:
		l.Z{l.E: err}.Critical()
	}
}
