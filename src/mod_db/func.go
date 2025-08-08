package mod_db

import (
	"context"
	"io/fs"
	"slices"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"

	"rmm23/src/l"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_ldap"
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
	var (
		ldap2doc = func(fnBaseDN string, fnSearchResultType string, fnSearchResult *ldap.SearchResult) (fnErr error) {
			var (
				entryType attrEntryType
				baseDN    attrDN
			)
			switch fnErr = entryType.Parse(fnSearchResultType); {
			case fnErr != nil:
				return
			}

			switch baseDN, fnErr = parseDN(fnBaseDN); {
			case fnErr != nil:
				return
			}

			for en := 0; en < len(fnSearchResult.Entries); en += l.BulkOpsSize {
				var (
					fnEntry     []*Entry
					fnCerts     []*Cert
					end         = min(en+l.BulkOpsSize, len(fnSearchResult.Entries))
					bulkEntries = fnSearchResult.Entries[en:end]
				)

				// Parse LDAP Entries
				switch fnErr = mod_ldap.UnmarshalLDAPEntries(bulkEntries, &fnEntry); {
				case err != nil:
					return
				}

				for _, entry := range fnEntry {
					entry.Type = entryType
					entry.BaseDN = baseDN
					entry.Status = entryStatusLoaded

					entry.Key = uuid.NewSHA1(uuid.Nil, []byte(entry.DN.String())).String()

					_ = repo.DeleteEntry(ctx, entry.Key)
				}

				switch e := repo.SaveMultiEntry(ctx, fnEntry...); {
				case e != nil:
					for a, b := range e {
						switch {
						case b != nil:
							l.Z{l.M: "repo.SaveMultiEntry", "key": fnEntry[a].Key, "DN": fnEntry[a].DN.String(), l.E: e}.Warning()
						}
					}
				}

				// Parse LDAP Entries's Certificates
				switch fnErr = mod_ldap.UnmarshalLDAPEntries(bulkEntries, &fnCerts); {
				case err != nil:
					return
				}

				fnCerts = slices.DeleteFunc(fnCerts, func(cert *Cert) bool {
					return cert == nil || cert.Certificate == nil
				})

				for _, cert := range fnCerts {
					// cert.Type = entryType
					// cert.BaseDN = baseDN
					// cert.Status = entryStatusLoaded
					cert.Key = uuid.NewSHA1(uuid.Nil, cert.Certificate.Certificate.Raw).String()

					_ = repo.DeleteEntry(ctx, cert.Key)
				}

				switch e := repo.SaveMultiCert(ctx, fnCerts...); {
				case e != nil:
					for a, b := range e {
						switch {
						case b != nil:
							l.Z{l.M: "repo.SaveMultiCert", "key": fnCerts[a].Key, "cert": fnCerts[a].Certificate.Certificate.Subject.String(), l.E: e}.Warning()
						}
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
