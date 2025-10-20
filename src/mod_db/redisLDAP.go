package mod_db

import (
	"context"

	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"

	"rmm23/src/l"
	"rmm23/src/mod_dn"
	"rmm23/src/mod_ldap"
	"rmm23/src/mod_strings"
)

func (r *RedisRepository) GetLDAPDocs(ctx context.Context, inbound *mod_ldap.Conf) (err error) {
	var (
		ldap2doc = func(fnBaseDN string, fnSearchResultType string, fnSearchResult *ldap.SearchResult) (fnErr error) {
			var (
				entryType attrEntryType
				baseDN    mod_dn.DN
			)
			switch fnErr = entryType.Parse(fnSearchResultType); {
			case fnErr != nil:
				return
			}

			switch fnErr = baseDN.UnmarshalText([]byte(fnBaseDN)); {
			case fnErr != nil:
				return
			}

			for en := 0; en < len(fnSearchResult.Entries); en += l.BulkOpsSize {
				var (
					fnEntry     []*Entry
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
					entry.Status = EntryStatusLoad

					entry.Key = uuid.NewSHA1(uuid.Nil, []byte(entry.DN.String())).String()

					_ = r.DeleteEntry(entry.Key)
				}

				switch e := r.SaveMultiEntry(fnEntry...); {
				case e != nil:
					for a, b := range e {
						switch {
						case b != nil:
							l.Z{l.M: "save", "key": fnEntry[a].Key, "DN": fnEntry[a].DN.String()}.Warning()
							l.Z{l.M: "save", "key": fnEntry[a].Key, "DN": fnEntry[a].DN.String(), l.E: e}.Debug()
						}
					}
				}

				// Parse LDAP Entries's Certificates `userPKCS12`
				var (
					fnCerts []*Cert
				)

				for _, b := range bulkEntries {
					for _, d := range b.GetRawAttributeValues(mod_strings.F_userPKCS12.String()) {
						var (
							fnCert = new(Cert)
						)
						switch forErr := fnCert.parseRaw(d); {
						case forErr != nil:
							continue
						}

						fnCerts = append(fnCerts, fnCert)
						_ = r.DeleteCert(fnCert.Key)
					}
				}

				switch e := r.SaveMultiCert(fnCerts...); {
				case e != nil:
					for a, b := range e {
						switch {
						case b != nil:
							l.Z{l.M: "save", "key": fnCerts[a].Key, "cert": fnCerts[a].Certificate.Certificate.Subject.String()}.Warning()
							l.Z{l.M: "save", "key": fnCerts[a].Key, "cert": fnCerts[a].Certificate.Certificate.Subject.String(), l.E: e}.Debug()
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
