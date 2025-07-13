package main

import (
	"errors"
	"flag"

	"github.com/avfs/avfs"
	"github.com/avfs/avfs/vfs/memfs"

	"rmm23/src/l"
	"rmm23/src/mod_db"
	"rmm23/src/mod_vfs"
)

func main() {
	l.Name.Set("rmm23")
	l.CLI.Set()

	l.Z{l.M: "main", "daemon": l.Name.String(), "commit": l.GitCommit.String()}.Informational()
	defer l.Z{l.M: "exit", "daemon": l.Name.String()}.Debug()

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
		a = mod_db.Entry{}
		b = a.RedisearchSchema()
	)
	b = b

	switch err = xmlConfig.load(); {
	case errors.Is(err, l.ENOCONF):
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

	switch err = xmlConfig.LDAP.Fetch(); {
	case err != nil:
		l.Z{l.E: err}.Critical()
	}
}

// 	var (
//		vfsDB = &mod_vfs.VFSDB{
//			List: make(map[string]string),
//			VFS: memfs.NewWithOptions(&memfs.Options{
//				Idm:        avfs.NewDummyIdm(),
//				User:       nil,
//				Name:       "",
//				OSType:     avfs.CurrentOSType(),
//				SystemDirs: nil,
//			}),
//		}
//	)
//	l.InitCLI()
//
//	switch {
//	case len(l.Config.String()) != 0: //
//		l.CLI.Set()
//
//		var (
//			cliConfigFile string
//			data          []byte
//		)
//
//		switch cliConfigFile, err = filepath.Abs(l.Config.String()); {
//		case err != nil:
//			return
//		}
//		switch err = vfsDB.CopyFromFS(cliConfigFile); {
//		case err != nil:
//			return
//		}
//		switch data, err = vfsDB.VFS.ReadFile(cliConfigFile); {
//		case err != nil:
//			return
//		}
//		switch err = xml.Unmarshal(data, r); {
//		case err != nil:
//			return
//		}
//
//	default:
//		return l.ENOCONF
//	}
//
//	return
