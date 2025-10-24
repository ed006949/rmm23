package mod_db

import (
	"context"

	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

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
				log.Error().Err(fnErr).Send()

				return
			}

			switch fnErr = baseDN.UnmarshalText([]byte(fnBaseDN)); {
			case fnErr != nil:
				log.Error().Err(fnErr).Send()

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
				case fnErr != nil:
					log.Error().Err(fnErr).Send()

					return
				}

				for _, entry := range fnEntry {
					entry.Type = entryType
					entry.BaseDN = baseDN
					entry.Status = entryStatusLoad

					entry.Key = uuid.NewSHA1(uuid.Nil, []byte(entry.DN.String())).String()

					_ = r.DeleteEntry(entry.Key)
				}

				switch e := r.SaveMultiEntry(fnEntry...); {
				case e != nil:
					for a, b := range e {
						switch {
						case b != nil:
							log.Warn().Str("key", fnEntry[a].Key).Str("DN", fnEntry[a].DN.String()).Msgf("save")
							log.Debug().Str("key", fnEntry[a].Key).Str("DN", fnEntry[a].DN.String()).Errs("errors", e).Msgf("save")
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
							log.Error().Err(fnErr).Send()

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
							log.Warn().Str("key", fnCerts[a].Key).Str("cert", fnCerts[a].Certificate.Certificate.Subject.String()).Msgf("save")
							log.Debug().Str("key", fnCerts[a].Key).Str("cert", fnCerts[a].Certificate.Certificate.Subject.String()).Errs("errors", e).Msgf("save")
						}
					}
				}
			}

			return
		}
	)
	switch err = inbound.SearchFn(ldap2doc); {
	case err != nil:
		log.Error().Err(err).Send()

		return
	}

	return
}
