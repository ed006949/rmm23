package mod_db

import (
	"context"
	"crypto/x509/pkix"

	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"

	"rmm23/src/l"
	"rmm23/src/mod_crypto"
	"rmm23/src/mod_ldap"
)

func CopyLDAP2DB(ctx context.Context, inbound *mod_ldap.Conf, outbound *Conf) (err error) {
	l.CLEAR = false

	switch err = outbound.Dial(ctx); {
	case err != nil:
		return
	}

	defer func() {
		_ = outbound.Close()
	}()

	switch l.CLEAR {
	case true:
		switch err = getLDAPDocs(ctx, inbound, outbound.Repo); {
		case err != nil:
			return
		}
	}

	return
}

// getLDAPDocs fetches entries from LDAP and saves them to Redis using the provided RedisRepository.
//
// copy entry one-by-one to save memory.
func getLDAPDocs(ctx context.Context, inbound *mod_ldap.Conf, repo *RedisRepository) (err error) {
	type entryCerts struct {
		UserPKCS12 mod_crypto.Certificates `ldap:"userPKCS12"`
	}

	var (
		ldap2doc = func(fnBaseDN *pkix.Name, fnSearchResultType string, fnSearchResult *ldap.SearchResult) (fnErr error) {
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
				switch fnErr = mod_ldap.UnmarshalEntry(fnB, fnEntry); {
				case fnErr != nil:
					return
				}

				fnEntry.Type = entryType
				fnEntry.BaseDN = &attrDN{*fnBaseDN}
				fnEntry.Status = entryStatusLoaded
				fnEntry.UUID.Generate(uuid.Nil, []byte(fnEntry.DN.String()))

				fnEntry.Key = fnEntry.UUID.String()

				_ = repo.DeleteEntry(ctx, fnEntry.Key)
				switch fnErr = repo.SaveEntry(ctx, fnEntry); {
				case fnErr != nil:
					return
				}

				var (
					cert    = new(entryCerts)
					fnCerts []*Cert
				)
				switch e := mod_ldap.UnmarshalEntry(fnB, cert); {
				case e != nil:
					l.Z{l.M: "mod_ldap.UnmarshalEntry", "DN": fnEntry.DN.String(), "cert": "all", l.E: e}.Warning()

					continue
				}

				for a, e := range cert.UserPKCS12 {
					var (
						fnCert = new(Cert)
					)

					// fnCert.Status = entryStatusLoaded
					fnCert.Certificate = e

					fnCert.Key = a

					_ = repo.DeleteCert(ctx, fnCert.Key)

					fnCerts = append(fnCerts, fnCert)
				}

				for a, e := range repo.SaveMultiCert(ctx, fnCerts...) {
					switch {
					case e != nil:
						l.Z{l.M: "Repo.SaveMultiCert", "DN": fnEntry.DN.String(), "cert": fnCerts[a].Key, l.E: e}.Warning()
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
