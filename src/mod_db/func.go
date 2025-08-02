package mod_db

import (
	"context"

	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"

	"rmm23/src/l"
	"rmm23/src/mod_crypto"
	"rmm23/src/mod_ldap"
)

func CopyLDAP2DB(ctx context.Context, inbound *mod_ldap.Conf, outbound *Conf) (err error) {
	l.CLEAR = true

	l.Z{l.M: "indexing", l.E: err}.Warning()

	switch err = outbound.Dial(ctx); {
	case err != nil:
		return
	}

	defer func() {
		_ = outbound.Close()
	}()

	switch err = getLDAPDocs(ctx, inbound, outbound.repo, outbound); {
	case err != nil:
		return
	}

	var (
		count   int64
		entries []*Entry
	)

	l.Z{l.M: "indexed", l.E: err}.Warning()

	count, entries, err = outbound.repo.SearchEntryMFV(
		ctx,
		[]_FV{
			{
				_type,
				entryTypeHost.Number() + " " + entryTypeHost.Number(),
			},
		},
	)
	l.Z{l.M: count, l.E: err, "entries": len(entries)}.Warning()

	count, entries, err = outbound.repo.SearchEntryMFV(
		ctx,
		_MFV{
			{
				_baseDN,
				"dc=fabric,dc=domain,dc=tld",
			},
			{
				_objectClass,
				"posixAccount",
			},
		},
	)
	l.Z{l.M: count, l.E: err, "entries": len(entries)}.Warning()

	return
}

// getLDAPDocs fetches entries from LDAP and saves them to Redis using the provided RedisRepository.
//
// copy entry one-by-one to save memory.
func getLDAPDocs(ctx context.Context, inbound *mod_ldap.Conf, repo *RedisRepository, c *Conf) (err error) {
	switch l.CLEAR {
	case false:
		return
	}

	type certificate struct {
		UserPKCS12 mod_crypto.Certificates `ldap:"userPKCS12"`
	}

	var (
		ldap2doc = func(fnBaseDN string, fnSearchResultType string, fnSearchResult *ldap.SearchResult) (fnErr error) {
			for _, fnB := range fnSearchResult.Entries {
				var (
					fnEntry = new(Entry)
				)

				switch fnErr = mod_ldap.UnmarshalEntry(fnB, fnEntry); {
				case fnErr != nil:
					return
				}

				switch fnErr = fnEntry.Type.Parse(fnSearchResultType); {
				case fnErr != nil:
					return
				}

				fnEntry.BaseDN = attrDN(fnBaseDN)
				fnEntry.Status = entryStatusLoaded
				fnEntry.UUID = attrUUID(uuid.NewSHA1(uuid.Nil, []byte(fnEntry.DN.String())))

				fnEntry.Key = fnEntry.UUID.String()

				// switch l.CLEAR {
				// case true:
				// 	_ = entry.DeleteEntry(ctx, fnEntry.Key)
				// }

				_ = repo.DeleteEntry(ctx, fnEntry.Key)

				// Save the Entry using the RedisRepository
				switch fnErr = repo.SaveEntry(ctx, fnEntry); {
				case fnErr != nil:
					return
				}

				_ = c.monitorIndexingFailures(ctx)

				var (
					cert    = new(certificate)
					fnCerts []*Certificate
				)

				switch e := mod_ldap.UnmarshalEntry(fnB, cert); {
				case e != nil:
					l.Z{l.M: "parse LDAP", "DN": fnEntry.Key, "cert": "all", l.E: e}.Warning()

					continue
				}

				for a, e := range cert.UserPKCS12 {
					var (
						fnCert = new(Certificate)
					)
					// fnCert.BaseDN = attrDN(fnBaseDN)
					fnCert.Status = entryStatusLoaded
					// fnCert.DN = attrDN(a)

					fnCert.Key = a
					fnCert.Certificate = e

					_ = repo.DeleteCert(ctx, fnCert.Key)

					fnCerts = append(fnCerts, fnCert)
				}

				for i, b := range repo.SaveMultiCert(ctx, fnCerts...) {
					switch {
					case b != nil:
						l.Z{l.M: "parse LDAP", "DN": fnEntry.Key, "cert": fnCerts[i].Key, l.E: b}.Warning()
					}
				}

				_ = c.monitorIndexingFailures(ctx)
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
