package main

import (
	"os"

	"rmm23/src/l"
)

func main() {
	l.Initialize()

	l.Z{l.M: "main", "commit": l.Run.CommitHashValue(), "built": l.Run.BuildTimeValue()}.Informational()
	defer l.Z{l.M: "exit"}.Informational()

	var (
		config = new(ConfigRoot)
		err    error
	)

	switch err = l.Run.ConfigUnmarshal(config); {
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
