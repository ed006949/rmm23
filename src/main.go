package main

import (
	"errors"
	"fmt"
	"net/netip"
	"os"
	"strconv"

	"github.com/avfs/avfs"
	"github.com/avfs/avfs/vfs/memfs"

	"rmm23/src/l"
	"rmm23/src/mod_db"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_net"
	"rmm23/src/mod_strings"
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

	// switch err = mod_net.Subnets.Generate(netip.MustParsePrefix("100.64.0.0/10"), mod_net.MaxIPv4Bits); {
	// case err != nil:
	// 	l.Z{l.E: err}.Critical()
	// }

	var (
		vlans        = []int{0, 1, 55, 66, 2001, 4094, 4095}
		vlansSubnets []netip.Prefix
	)

	switch vlansSubnets, err = mod_net.Subnets.SubnetList(netip.MustParsePrefix("10.240.192.0/18"), mod_net.MaxIPv4Bits-mod_net.HostSubnetBits, vlans...); {
	case err != nil:
		l.Z{l.E: err}.Critical()
	}

	for a, b := range vlans {
		fmt.Printf("VLAN%04d: %18s\n", b, vlansSubnets[a])
	}

	switch vlansSubnets, err = mod_net.Subnets.SubnetList(netip.MustParsePrefix("10.240.192.0/30"), mod_net.MaxIPv4Bits-mod_net.HostSubnetBits, 0, 0, 0); {
	case err != nil:
		l.Z{l.E: err}.Critical()
	}

	for a, b := range []int{0, 0, 0} {
		fmt.Printf("VLAN%04d: %18s\n", b, vlansSubnets[a])
	}

	switch err = mod_net.Subnets.Generate(netip.MustParsePrefix("172.16.0.0/12"), mod_net.MaxIPv4Bits-mod_net.UserSubnetBits); {
	case err != nil:
		l.Z{l.E: err}.Critical()
	}

	// switch vlansSubnets, err = mod_net.Subnets.Subnets(netip.MustParsePrefix("10.240.192.0/16"), mod_net.MaxIPv4Bits-mod_net.HostSubnetBits-2); {
	// case err != nil:
	// 	l.Z{l.E: err}.Critical()
	// }
	//
	// for a, b := range vlansSubnets {
	// 	fmt.Printf("ID%010d: %18s\n", a, b)
	// }

	// switch vlansSubnets, err = mod_net.Subnets.Subnets(netip.MustParsePrefix("10.92.0.0/16"), mod_net.MaxIPv4Bits-mod_net.HostSubnetBits); {
	// case err != nil:
	// 	l.Z{l.E: err}.Critical()
	// }
	//
	// for a, b := range vlansSubnets {
	// 	fmt.Printf("TI%05d: %18s\n", a, b)
	// }

	switch _, entries, err = config.Conf.DB.Repo.SearchEntryFVs(
		&mod_strings.FVs{
			{
				mod_strings.F_type,
				mod_db.EntryTypeUser.Number() + " " + mod_db.EntryTypeUser.Number(),
			},
		},
	); {
	case err != nil:
		l.Z{l.E: err}.Critical()
	}

	for _, b := range entries {
		switch len(b.IPHostNumber) {
		case 0:
			l.Z{l.E: fmt.Errorf("no prefix in '%v'", b.DN.String())}.Warning()
		case 1:
			switch err = mod_net.Subnets.PrefixUse(netip.MustParsePrefix("172.16.0.0/12"), mod_net.MaxIPv4Bits-mod_net.UserSubnetBits, b.IPHostNumber[0]); {
			case errors.Is(err, mod_errors.EEXIST):
				l.Z{l.E: fmt.Errorf("prefix '%v' in '%v' is already used", b.IPHostNumber[0].String(), b.DN.String())}.Warning()
			case err != nil:
				l.Z{l.E: fmt.Errorf("invalid prefix '%v' in '%v'", b.IPHostNumber[0].String(), b.DN.String())}.Warning()
			}
		default:
			l.Z{l.E: fmt.Errorf("too many prefixes '%v' in '%v'", b.IPHostNumber, b.DN.String())}.Warning()
		}
	}

	os.Exit(1)
}
