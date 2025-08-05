package mod_db

import (
	"context"

	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"

	"rmm23/src/l"
	"rmm23/src/mod_bools"
	"rmm23/src/mod_crypto"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_ldap"
	"rmm23/src/mod_net"
)

func CopyLDAP2DB(ctx context.Context, inbound *mod_ldap.Conf, outbound *Conf) (err error) {
	switch err = outbound.Dial(ctx); {
	case err != nil:
		return
	}

	defer func() {
		_ = outbound.Close()
	}()

	switch {
	case !l.Run.DryRunValue():
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
				switch e := mod_ldap.UnmarshalEntry(fnB, fnEntry); {
				case e != nil:
					l.Z{l.M: "mod_ldap.UnmarshalEntry", "DN": fnEntry.DN.String(), l.E: e}.Warning()

					return
				}

				fnEntry.Type = entryType
				_ = fnEntry.BaseDN.Parse(fnBaseDN)
				fnEntry.Status = entryStatusLoaded
				fnEntry.UUID.Generate(uuid.Nil, []byte(fnEntry.DN.String()))

				fnEntry.Key = fnEntry.UUID.String()

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
				switch e := mod_ldap.UnmarshalEntry(fnB, cert); {
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
							UUID:           attrUUID{},
							SerialNumber:   e.Certificate.SerialNumber,
							Issuer:         attrDN{},
							Subject:        attrDN{},
							NotBefore:      attrTime{e.Certificate.NotBefore},
							NotAfter:       attrTime{e.Certificate.NotAfter},
							DNSNames:       e.Certificate.DNSNames,
							EmailAddresses: e.Certificate.EmailAddresses,
							IPAddresses:    mod_errors.StripErr1(mod_net.ParseNetIPs(e.Certificate.IPAddresses)),
							URIs:           e.Certificate.URIs,
							IsCA:           mod_bools.AttrBool(e.Certificate.IsCA),
							Certificate:    e,
						}
					)

					_ = fnCert.Issuer.Parse(e.Certificate.Issuer.String())
					_ = fnCert.Subject.Parse(e.Certificate.Subject.String())
					fnCert.UUID.Generate(uuid.NameSpaceOID, fnCert.Certificate.Certificate.Raw)
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
