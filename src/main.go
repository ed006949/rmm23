package main

import (
	"os"

	"github.com/avfs/avfs"
	"github.com/avfs/avfs/vfs/memfs"

	"rmm23/src/l"
	"rmm23/src/mod_vfs"
)

func main() {
	l.Initialize()

	l.Z{l.M: "main", "commit": l.Run.CommitHashValue(), "built": l.Run.BuildTimeValue()}.Informational()
	defer l.Z{l.M: "exit"}.Informational()

	var (
		config = new(ConfigRoot)
		err    error
		vfsDB  = &mod_vfs.VFSDB{
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
	switch err = l.Run.ConfigUnmarshal(&config); {
	case err != nil:
		os.Exit(1)
	}

	switch err = config.Conf.DB.Dial(ctx); {
	case err != nil:
		return
	}

	defer func() {
		_ = config.Conf.DB.Close()
	}()

	switch {
	case !l.Run.DryRunValue():
		switch err = config.Conf.DB.Repo.GetLDAPDocs(ctx, config.Conf.LDAP); {
		case err != nil:
			l.Z{l.E: err}.Critical()
		}
	}

	switch err = vfsDB.CopyFromFS("./etc/legacy/"); {
	case err != nil:
		l.Z{l.E: err}.Critical()
	}

	switch {
	case !l.Run.DryRunValue():
		switch err = config.Conf.DB.Repo.GetFSCerts(ctx, vfsDB); {
		case err != nil:
			l.Z{l.E: err}.Critical()
		}
	}

	var (
	// count   int64
	// entries []*mod_db.Entry
	// certs       []*mod_db.Cert
	)

	switch _, err = config.Conf.DB.Repo.CheckIPHostNumber(config.Conf.Networking.User.Subnet, config.Conf.Networking.User.Bits); {
	case err != nil:
		l.Z{l.E: err}.Critical()
	}

	// switch errs := config.Conf.DB.Repo.UpdateMultiEntry(entries...); {
	// case errs != nil:
	// 	l.Z{l.E: err}.Critical()
	// }

	os.Exit(1)
}
