package main

import (
	"net/netip"
	"os"

	"github.com/avfs/avfs"
	"github.com/avfs/avfs/vfs/memfs"

	"rmm23/src/l"
	"rmm23/src/mod_db"
	"rmm23/src/mod_net"
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
		entries []*mod_db.Entry
		// certs       []*mod_db.Cert
		usersSubnet = netip.MustParsePrefix("172.16.0.0/12")
		userBits    = mod_net.MaxIPv4Bits - mod_net.UserSubnetBits
	)

	switch entries, err = config.Conf.DB.Repo.CheckIPHostNumber(usersSubnet, userBits); {
	case err != nil:
		l.Z{l.E: err}.Critical()
	}

	for _, b := range entries {
		switch {
		case b.Status == mod_db.EntryStatusUpdated:
			l.Z{l.M: "updated entry", "DN": b.DN.String()}.Informational()

			b.Ver++
			switch err = config.Conf.DB.Repo.SaveEntry(b); {
			case err != nil:
				l.Z{l.E: err, "DN": b.DN.String()}.Warning()
			}
		}
	}

	os.Exit(1)
}
