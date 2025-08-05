package main

import (
	"os"

	"github.com/avfs/avfs"
	"github.com/avfs/avfs/vfs/memfs"

	"rmm23/src/l"
	"rmm23/src/mod_db"
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

	switch {
	case !l.Run.DryRunValue():
		switch err = mod_db.GetLDAPDocs(ctx, config.Conf.LDAP, config.Conf.DB); {
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
		switch err = mod_db.GetFSCerts(ctx, vfsDB, config.Conf.DB); {
		case err != nil:
			l.Z{l.E: err}.Critical()
		}
	}

	switch err = config.Conf.DB.Dial(ctx); {
	case err != nil:
		return
	}

	defer func() {
		_ = config.Conf.DB.Close()
	}()

	var (
		count   int64
		entries []*mod_db.Entry
		certs   []*mod_db.Cert
	)

	count, entries, err = config.Conf.DB.Repo.SearchEntryMFV(
		ctx,
		[]mod_db.FV{
			{
				mod_db.F_type,
				mod_db.EntryTypeHost.Number() + " " + mod_db.EntryTypeHost.Number(),
			},
		},
	)
	l.Z{l.M: count, l.E: err, "entries": len(entries)}.Warning()

	count, entries, err = config.Conf.DB.Repo.SearchEntryMFV(
		ctx,
		mod_db.MFV{
			{
				mod_db.F_baseDN,
				"dc=fabric,dc=domain,dc=tld",
			},
			{
				mod_db.F_objectClass,
				"posixAccount",
			},
		},
	)
	l.Z{l.M: count, l.E: err, "entries": len(entries)}.Warning()

	count, entries, err = config.Conf.DB.Repo.SearchEntryQ(ctx, "*")
	l.Z{l.M: count, l.E: err, "entries": len(entries)}.Warning()

	count, certs, err = config.Conf.DB.Repo.SearchCertMFV(
		ctx,
		[]mod_db.FV{
			{
				mod_db.F_isCA,
				"1 1",
			},
		},
	)
	l.Z{l.M: count, l.E: err, "entries": len(certs)}.Warning()

	os.Exit(1)
}
