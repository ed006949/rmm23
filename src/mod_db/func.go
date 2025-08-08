package mod_db

import (
	"context"
	"io/fs"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"

	"rmm23/src/l"
	"rmm23/src/mod_bools"
	"rmm23/src/mod_crypto"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_ldap"
	"rmm23/src/mod_net"
	"rmm23/src/mod_vfs"
)

func GetLDAPDocs(ctx context.Context, inbound *mod_ldap.Conf, outbound *Conf) (err error) {
	switch err = outbound.Dial(ctx); {
	case err != nil:
		return
	}

	defer func() {
		_ = outbound.Close()
	}()

	switch err = getLDAPDocs(ctx, inbound, outbound.Repo); {
	case err != nil:
		return
	}

	return
}

func GetFSCerts(ctx context.Context, inbound *mod_vfs.VFSDB, outbound *Conf) (err error) {
	switch err = outbound.Dial(ctx); {
	case err != nil:
		return
	}

	defer func() {
		_ = outbound.Close()
	}()

	switch err = getFSCerts(ctx, inbound, outbound.Repo); {
	case err != nil:
		return
	}

	return
}

func getLDAPDocs(ctx context.Context, inbound *mod_ldap.Conf, repo *RedisRepository) (err error) {
	type entryCerts struct {
		UserPKCS12 []*mod_crypto.Certificate `json:"userPKCS12,omitempty" ldap:"userPKCS12"`
	}

	var (
		ldap2doc = func(fnBaseDN string, fnSearchResultType string, fnSearchResult *ldap.SearchResult) (fnErr error) {
			var (
				entryType attrEntryType
			)
			switch fnErr = entryType.Parse(fnSearchResultType); {
			case fnErr != nil:
				return
			}

			for _, fnB := range fnSearchResult.Entries {
				var (
					fnEntry = new(Entry)
				)
				switch e := mod_ldap.UnmarshalLDAP(fnB, fnEntry); {
				case e != nil:
					l.Z{l.M: "mod_ldap.UnmarshalEntry", "DN": fnEntry.DN.String(), l.E: e}.Warning()

					return
				}

				fnEntry.Type = entryType
				fnEntry.BaseDN = mod_errors.StripErr1(parseDN(fnBaseDN))
				fnEntry.Status = entryStatusLoaded
				// tUUID := uuid.NewSHA1(uuid.Nil, []byte(fnEntry.DN.String()))
				// fnEntry.UUID = tUUID

				fnEntry.Key = uuid.NewSHA1(uuid.Nil, []byte(fnEntry.DN.String())).String()

				_ = repo.DeleteEntry(ctx, fnEntry.Key)
				switch e := repo.SaveEntry(ctx, fnEntry); {
				case e != nil:
					l.Z{l.M: "repo.SaveEntry", "DN": fnEntry.DN.String(), l.E: e}.Warning()

					continue
				}

				var (
					cert    = new(entryCerts)
					fnCerts []*Cert
				)
				switch e := mod_ldap.UnmarshalLDAP(fnB, cert); {
				case e != nil:
					l.Z{l.M: "mod_ldap.UnmarshalEntry", "DN": fnEntry.DN.String(), "cert": "all", l.E: e}.Warning()

					continue
				}

				for _, e := range cert.UserPKCS12 {
					var (
						fnCert = &Cert{
							// Key:            "",
							// Ver:            0,
							Ext:            e.Certificate.NotAfter,
							UUID:           uuid.NewSHA1(uuid.NameSpaceOID, e.Certificate.Raw),
							SerialNumber:   e.Certificate.SerialNumber,
							Issuer:         mod_errors.StripErr1(parseDN(e.Certificate.Issuer.String())),
							Subject:        mod_errors.StripErr1(parseDN(e.Certificate.Subject.String())),
							NotBefore:      e.Certificate.NotBefore,
							NotAfter:       e.Certificate.NotAfter,
							DNSNames:       e.Certificate.DNSNames,
							EmailAddresses: e.Certificate.EmailAddresses,
							IPAddresses:    mod_errors.StripErr1(mod_net.ParseNetIPs(e.Certificate.IPAddresses)),
							URIs:           e.Certificate.URIs,
							IsCA:           mod_bools.AttrBool(e.Certificate.IsCA),
							Certificate:    e,
						}
					)

					fnCert.Key = fnCert.UUID.String()

					// fnCert.Status = entryStatusLoaded

					_ = repo.DeleteCert(ctx, fnCert.Key)

					fnCerts = append(fnCerts, fnCert)
				}

				for a, e := range repo.SaveMultiCert(ctx, fnCerts...) {
					switch {
					case e != nil:
						l.Z{l.M: "repo.SaveMultiCert", "DN": fnEntry.DN.String(), "cert": fnCerts[a].Key, l.E: e}.Warning()
					}
				}
			}

			return
		}
	)
	switch err = inbound.SearchFn(ldap2doc); {
	case err != nil:
		return
	}

	return
}

func getFSCerts(ctx context.Context, vfsDB *mod_vfs.VFSDB, repo *RedisRepository) (err error) {
	var (
		c          = make(map[string][][]byte)
		fileExts   = 2
		totalFiles = 6

		fn = func(name string, dirEntry fs.DirEntry, err error) (fnErr error) {
			switch {
			case err != nil:
				return err
			}

			var (
				s = strings.Split(name, ".")
			)
			switch {
			case len(s) < fileExts:
				return
			}

			var (
				n = strings.Join(s[:len(s)-fileExts], ".")
			)
			switch _, ok := c[n]; {
			case !ok:
				c[n] = make([][]byte, totalFiles)
			}

			switch s[len(s)-1] {
			case "der":
				switch s[len(s)-2] {
				case "key":
					c[n][0], _ = vfsDB.VFS.ReadFile(name)
				case "crt":
					c[n][1], _ = vfsDB.VFS.ReadFile(name)
				case "ca":
					c[n][2], _ = vfsDB.VFS.ReadFile(name)
				case "csr":
					c[n][3], _ = vfsDB.VFS.ReadFile(name)
				case "crl":
					c[n][4], _ = vfsDB.VFS.ReadFile(name)
				}
			case "pem":
				c[n][5] = append(c[n][5], mod_errors.StripErr1(vfsDB.VFS.ReadFile(name))...)
			}

			return
		}
	)
	switch err = vfsDB.VFS.WalkDir("/", fn); {
	case err != nil:
		l.Z{l.E: err}.Error()
	}

	for a, b := range c {
		var (
			forErr  error
			forCert = new(Cert)
		)

		switch forErr = forCert.parseRaw(b[0], b[1], b[2], b[3], b[4], b[5]); {
		case forErr != nil:
			continue
		}

		_ = repo.DeleteCert(ctx, forCert.Key)

		switch forErr = repo.SaveCert(ctx, forCert); {
		case forErr != nil:
			l.Z{l.M: "repo.SaveCert", "cert": a, l.E: forErr}.Warning()
		}
	}

	return
}
