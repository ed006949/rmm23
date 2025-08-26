package main

import (
	"fmt"
	"net/netip"
	"os"
	"strconv"

	"github.com/avfs/avfs"
	"github.com/avfs/avfs/vfs/memfs"

	"rmm23/src/l"
	"rmm23/src/mod_db"
	"rmm23/src/mod_strings"
	"rmm23/src/mod_vfs"
	"rmm23/src/mod_vlan"
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
		vlanSubnets = mod_vlan.NewSubnets()
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
		switch err = mod_db.GetLDAPDocs(ctx, config.Conf.LDAP, config.Conf.DB.Repo); {
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
		switch err = mod_db.GetFSCerts(ctx, vfsDB, config.Conf.DB.Repo); {
		case err != nil:
			l.Z{l.E: err}.Critical()
		}
	}

	var (
		count   int64
		entries []*mod_db.Entry
		certs   []*mod_db.Cert
	)

	count, entries, err = config.Conf.DB.Repo.SearchEntryFVs(
		&mod_strings.FVs{
			{
				mod_strings.F_type,
				mod_db.EntryTypeHost.Number() + " " + mod_db.EntryTypeHost.Number(),
			},
		},
	)
	l.Z{l.M: count, l.E: err, "entries": len(entries)}.Warning()

	count, entries, err = config.Conf.DB.Repo.SearchEntryFVs(
		&mod_strings.FVs{
			{
				mod_strings.F_baseDN,
				"dc=fabric,dc=domain,dc=tld",
			},
			{
				mod_strings.F_objectClass,
				"posixAccount",
			},
		},
	)
	l.Z{l.M: count, l.E: err, "entries": len(entries)}.Warning()

	count, entries, err = config.Conf.DB.Repo.SearchEntryQ("*")
	l.Z{l.M: count, l.E: err, "entries": len(entries)}.Warning()

	count, certs, err = config.Conf.DB.Repo.SearchCertFVs(
		&mod_strings.FVs{
			{
				mod_strings.F_isCA,
				strconv.FormatBool(true),
			},
		},
	)
	l.Z{l.M: count, l.E: err, "entries": len(certs)}.Warning()

	var (
		vlans        = []int{0, 1, 2001, 4094, 4095}
		vlansSubnets []netip.Prefix
	)

	switch err = vlanSubnets.GenerateSubnets(netip.MustParseAddr("10.240.192.0"), mod_vlan.MaxIPv4Bits-mod_vlan.HostSubnetSize); {
	case err != nil:
		l.Z{l.E: err}.Critical()
	}

	switch vlansSubnets, err = vlanSubnets.Subnets(netip.MustParseAddr("10.240.192.0"), vlans...); {
	case err != nil:
		l.Z{l.E: err}.Critical()
	}

	for a, b := range vlansSubnets {
		fmt.Printf("VLAN%04d: %18s\n", vlans[a], b)
	}

	// switch vlansSubnets, err = vlanSubnets.Subnets(netip.MustParseAddr("10.240.192.0")); {
	// case err != nil:
	// 	l.Z{l.E: err}.Critical()
	// }
	// for a, b := range vlansSubnets {
	// 	fmt.Printf("VLAN%04d: %18s\n", a, b)
	// }

	os.Exit(1)
}
