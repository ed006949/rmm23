package mod_db

import (
	"context"
	"crypto"
	"crypto/x509"
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
				_ = fnEntry.BaseDN.parse(fnBaseDN)
				fnEntry.Status = entryStatusLoaded
				fnEntry.UUID.generate(uuid.Nil, []byte(fnEntry.DN.String()))

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
							UUID:           generateUUID(uuid.NameSpaceOID, e.Certificate.Raw),
							SerialNumber:   e.Certificate.SerialNumber,
							Issuer:         mod_errors.StripErr1(parseDN(e.Certificate.Issuer.String())),
							Subject:        mod_errors.StripErr1(parseDN(e.Certificate.Subject.String())),
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
		totalFiles = 4

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
			case s[len(s)-1] != "der":
				return
			}

			var (
				n = strings.Join(s[:len(s)-fileExts], ".")
			)

			switch _, ok := c[n]; {
			case !ok:
				c[n] = make([][]byte, totalFiles)
			}

			switch s[len(s)-2] {
			case "key":
				c[n][0], _ = vfsDB.VFS.ReadFile(name)
			case "crt":
				c[n][1], _ = vfsDB.VFS.ReadFile(name)
			case "crl":
				c[n][2], _ = vfsDB.VFS.ReadFile(name)
			case "csr":
				c[n][3], _ = vfsDB.VFS.ReadFile(name)
			default:
				return
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
			forErr error
			key    crypto.PrivateKey
			cert   *x509.Certificate
			crl    *x509.RevocationList
			csr    *x509.CertificateRequest
		)

		switch key, forErr = mod_crypto.ParsePrivateKey(b[0]); {
		case forErr != nil:
			continue
		}

		switch cert, forErr = x509.ParseCertificate(b[1]); {
		case forErr != nil:
			continue
		}

		switch crl, forErr = x509.ParseRevocationList(b[2]); {
		case forErr != nil:
			// continue
		}

		switch csr, forErr = x509.ParseCertificateRequest(b[3]); {
		case forErr != nil:
			// continue
		}

		var (
			fnCert = &Cert{
				// Key:            "",
				// Ver:            0,
				Ext:            cert.NotAfter,
				UUID:           generateUUID(uuid.NameSpaceOID, cert.Raw),
				SerialNumber:   cert.SerialNumber,
				Issuer:         mod_errors.StripErr1(parseDN(cert.Issuer.String())),
				Subject:        mod_errors.StripErr1(parseDN(cert.Subject.String())),
				NotBefore:      attrTime{cert.NotBefore},
				NotAfter:       attrTime{cert.NotAfter},
				DNSNames:       cert.DNSNames,
				EmailAddresses: cert.EmailAddresses,
				IPAddresses:    mod_errors.StripErr1(mod_net.ParseNetIPs(cert.IPAddresses)),
				URIs:           cert.URIs,
				IsCA:           mod_bools.AttrBool(cert.IsCA),
				Certificate: &mod_crypto.Certificate{
					P12: nil,
					DER: nil,
					PEM: nil,
					CRL: func() (outbound []byte) {
						switch {
						case crl != nil:
							return crl.Raw
						}

						return
					}(),
					CSR: func() (outbound []byte) {
						switch {
						case csr != nil:
							return csr.Raw
						}

						return
					}(),
					PrivateKeyDER:         nil,
					CertificateRequestDER: nil,
					CertificateDER:        nil,
					CertificateCAChainDER: nil,
					RevocationListDER:     nil,
					PrivateKeyPEM:         nil,
					CertificateRequestPEM: nil,
					CertificatePEM:        nil,
					CertificateCAChainPEM: nil,
					RevocationListPEM:     nil,
					PrivateKey:            key,
					CertificateRequest:    csr,
					Certificate:           cert,
					CertificateCAChain:    nil,
					RevocationList:        crl,
				},
			}
		)

		fnCert.Key = fnCert.UUID.String()

		switch forErr = fnCert.Certificate.EncodeP12(); {
		case forErr != nil:
			continue
		}

		_ = repo.DeleteCert(ctx, fnCert.Key)

		switch forErr = repo.SaveCert(ctx, fnCert); {
		case forErr != nil:
			l.Z{l.M: "repo.SaveCert", "cert": a, l.E: forErr}.Warning()
		}
	}

	return
}
