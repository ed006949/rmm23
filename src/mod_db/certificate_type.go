package mod_db

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"math/big"
	"net"
	"net/url"
	"time"

	"rmm23/src/mod_crypto"
	"rmm23/src/mod_net"
)

// Cert is the struct that represents an LDAP userPKCS12 attribute.
//
// when updating @src/mod_db/entry_type.go don't forget to update:
//
//	@src/mod_db/certificate_*.go
//	@src/mod_db/redis_*.go
type Cert struct {
	// db data
	Key string    `redis:",key"`  //
	Ver int64     `redis:",ver"`  //
	Ext time.Time `redis:",exat"` //

	UUID attrUUID `json:"uuid,omitempty" msgpack:"uuid"` //

	SerialNumber   *big.Int         `json:"serialNumber"`   // redis:",key"
	Issuer         attrDN           `json:"issuer"`         //
	Subject        attrDN           `json:"subject"`        //
	NotBefore      attrTime         `json:"notBefore"`      //
	NotAfter       attrTime         `json:"notAfter"`       // redis:",exat"
	DNSNames       []string         `json:"dnsNames"`       //
	EmailAddresses []string         `json:"emailAddresses"` //
	IPAddresses    []*attrIPAddress `json:"ipAddresses"`    //
	URIs           []*mod_net.URL   `json:"uris"`           //
	IsCA           bool             `json:"isCA"`           //

	// // element specific meta data
	// Type   attrEntryType   `json:"type,omitempty"   msgpack:"type"`   // (?) Certificate's type
	// Status attrEntryStatus `json:"status,omitempty" msgpack:"status"` //
	// BaseDN attrDN          `json:"baseDN,omitempty" msgpack:"baseDN"` //

	// // element meta data
	// UUID            attrUUID `json:"uuid,omitempty"            msgpack:"uuid"`            //  must be unique
	// DN              attrDN   `json:"dn,omitempty"              msgpack:"dn"`              //  must be unique
	// CreatorsName    attrDN   `json:"creatorsName,omitempty"    msgpack:"creatorsName"`    //
	// CreateTimestamp attrTime `json:"createTimestamp,omitempty" msgpack:"createTimestamp"` //
	// ModifiersName   attrDN   `json:"modifiersName,omitempty"   msgpack:"modifiersName"`   //
	// ModifyTimestamp attrTime `json:"modifyTimestamp,omitempty" msgpack:"modifyTimestamp"` //

	// element data
	Certificate *mod_crypto.Certificate `json:"certificate,omitempty" msgpack:"certificate"` //

	// // element data
	// CN           attrString        `json:"cn,omitempty"           msgpack:"cn"`           //  RDN in group's context
	// DC           attrString        `json:"dc,omitempty"           msgpack:"dc"`           //
	// Description  attrString        `json:"description,omitempty"  msgpack:"description"`  //
	// IPHostNumber attrIPPrefixes `json:"ipHostNumber,omitempty" msgpack:"ipHostNumber"` //
	// Mail         attrMails         `json:"mail,omitempty"         msgpack:"mail"`         //
	// O            attrString        `json:"o,omitempty"            msgpack:"o"`            //
	// OU           attrString        `json:"ou,omitempty"           msgpack:"ou"`           //
	// SN           attrString        `json:"sn,omitempty"           msgpack:"sn"`           //

	// 			d.PKI, is_new = i_PKI.verify(nil, d.FQDN, &x509.Certificate{
	//				SignatureAlgorithm: x509.ECDSAWithSHA512,
	//				// SignatureAlgorithm: x509.PureEd25519,
	//				SerialNumber: pki_crt_sn(),
	//				Subject: pkix.Name{
	//					Organization: []string{d.FQDN.String()},
	//					CommonName:   d.FQDN.String(),
	//					Names:        nil,
	//					ExtraNames:   nil,
	//				},
	//				NotBefore:             time.Now(),
	//				NotAfter:              pki_crt_expiry(),
	//				IsCA:                  true,
	//				ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	//				KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	//				BasicConstraintsValid: true,
	//				CRLDistributionPoints: []string{join_string("", "http://", join_string(".", "ns", d.FQDN), "/crl.pem")},
	//				DNSNames:              []string{d.FQDN.String()},
	//				EmailAddresses:        []string{join_string("@", "ns", d.FQDN)},
	//				IPAddresses:           nil,
	//			})

	Crt struct {
		Raw                         []byte                  `json:"-"` //
		RawTBSCertificate           []byte                  `json:"-"` //
		RawSubjectPublicKeyInfo     []byte                  `json:"-"` //
		RawSubject                  []byte                  `json:"-"` //
		RawIssuer                   []byte                  `json:"-"` //
		Signature                   []byte                  `json:"-"` //
		SignatureAlgorithm          x509.SignatureAlgorithm `json:"-"` //
		PublicKeyAlgorithm          x509.PublicKeyAlgorithm `json:"-"` //
		PublicKey                   any                     `json:"-"` //
		Version                     int                     `json:"-"` //
		SerialNumber                *big.Int                `json:"-"` //
		Issuer                      pkix.Name               `json:"-"` //
		Subject                     pkix.Name               `json:"-"` //
		NotBefore, NotAfter         time.Time               `json:"-"` //
		KeyUsage                    x509.KeyUsage           `json:"-"` //
		Extensions                  []pkix.Extension        `json:"-"` //
		ExtraExtensions             []pkix.Extension        `json:"-"` //
		UnhandledCriticalExtensions []asn1.ObjectIdentifier `json:"-"` //
		ExtKeyUsage                 []x509.ExtKeyUsage      `json:"-"` //
		UnknownExtKeyUsage          []asn1.ObjectIdentifier `json:"-"` //
		BasicConstraintsValid       bool                    `json:"-"` //
		IsCA                        bool                    `json:"-"` //
		MaxPathLen                  int                     `json:"-"` //
		MaxPathLenZero              bool                    `json:"-"` //
		SubjectKeyId                []byte                  `json:"-"` //
		AuthorityKeyId              []byte                  `json:"-"` //
		OCSPServer                  []string                `json:"-"` //
		IssuingCertificateURL       []string                `json:"-"` //
		DNSNames                    []string                `json:"-"` //
		EmailAddresses              []string                `json:"-"` //
		IPAddresses                 []net.IP                `json:"-"` //
		URIs                        []*url.URL              `json:"-"` //
		PermittedDNSDomainsCritical bool                    `json:"-"` //
		PermittedDNSDomains         []string                `json:"-"` //
		ExcludedDNSDomains          []string                `json:"-"` //
		PermittedIPRanges           []*net.IPNet            `json:"-"` //
		ExcludedIPRanges            []*net.IPNet            `json:"-"` //
		PermittedEmailAddresses     []string                `json:"-"` //
		ExcludedEmailAddresses      []string                `json:"-"` //
		PermittedURIDomains         []string                `json:"-"` //
		ExcludedURIDomains          []string                `json:"-"` //
		CRLDistributionPoints       []string                `json:"-"` //
		PolicyIdentifiers           []asn1.ObjectIdentifier `json:"-"` //
		Policies                    []x509.OID              `json:"-"` //
		InhibitAnyPolicy            int                     `json:"-"` //
		InhibitAnyPolicyZero        bool                    `json:"-"` //
		InhibitPolicyMapping        int                     `json:"-"` //
		InhibitPolicyMappingZero    bool                    `json:"-"` //
		RequireExplicitPolicy       int                     `json:"-"` //
		RequireExplicitPolicyZero   bool                    `json:"-"` //
		PolicyMappings              []x509.PolicyMapping    `json:"-"` //
	} `json:"-"` //
}
