package main

import (
	"os"

	"rmm23/src/l"
)

func main() {
	l.Z{l.M: "main", "daemon": l.Run.Name(), "commit": l.Run.CommitHash()}.Informational()
	defer l.Z{l.M: "exit", "daemon": l.Run.Name()}.Informational()

	switch err := l.Run.UnmarshalConfig(config); {
	case err != nil:
		panic(err)
	}

	// var (
	// 	err       error
	// 	xmlConfig = new(xmlConf)
	// 	vfsDB     = &mod_vfs.VFSDB{
	// 		List: make(map[string]string),
	// 		VFS: memfs.NewWithOptions(&memfs.Options{
	// 			Idm:        avfs.NotImplementedIdm,
	// 			User:       nil,
	// 			Name:       "",
	// 			OSType:     avfs.CurrentOSType(),
	// 			SystemDirs: nil,
	// 		}),
	// 	}
	// )

	// switch err = xmlConfig.load(); {
	// case errors.Is(err, mod_errors.ENOCONF):
	// 	flag.PrintDefaults()
	// 	l.Z{l.E: err}.Critical()
	// case err != nil:
	// 	flag.PrintDefaults()
	// 	l.Z{l.E: err}.Critical()
	// }
	//
	// switch err = vfsDB.CopyFromFS("./etc/legacy/"); {
	// case err != nil:
	// 	l.Z{l.E: err}.Critical()
	// }
	//
	// switch err = mod_db.CopyLDAP2DB(ctx, xmlConfig.LDAP); {
	// case err != nil:
	// 	l.Z{l.E: err}.Critical()
	// }
	os.Exit(1)
}
