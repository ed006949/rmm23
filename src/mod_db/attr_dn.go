package mod_db

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tsaarni/x500dn"

	"rmm23/src/mod_slices"
)

var oidMap = map[string]string{
	"2.5.4.15":                   "businesscategory",
	"2.5.4.6":                    "c",
	"2.5.4.3":                    "cn",
	"0.9.2342.19200300.100.1.25": "dc",
	"2.5.4.13":                   "description",
	"2.5.4.27":                   "destinationindicator",
	"2.5.4.49":                   "distinguishedName",
	"2.5.4.46":                   "dnqualifier",
	"1.2.840.113549.1.9.1":       "emailaddress",
	"2.5.4.47":                   "enhancedsearchguide",
	"2.5.4.23":                   "facsimiletelephonenumber",
	"2.5.4.44":                   "generationqualifier",
	"2.5.4.42":                   "givenname",
	"2.5.4.51":                   "houseidentifier",
	"2.5.4.43":                   "initials",
	"2.5.4.25":                   "internationalisdnnumber",
	"2.5.4.7":                    "l",
	"2.5.4.31":                   "member",
	"2.5.4.41":                   "name",
	"2.5.4.10":                   "o",
	"2.5.4.11":                   "ou",
	"2.5.4.32":                   "owner",
	"2.5.4.19":                   "physicaldeliveryofficename",
	"2.5.4.16":                   "postaladdress",
	"2.5.4.17":                   "postalcode",
	"2.5.4.18":                   "postOfficebox",
	"2.5.4.28":                   "preferreddeliverymethod",
	"2.5.4.26":                   "registeredaddress",
	"2.5.4.33":                   "roleoccupant",
	"2.5.4.14":                   "searchguide",
	"2.5.4.34":                   "seealso",
	"2.5.4.5":                    "serialnumber",
	"2.5.4.4":                    "sn",
	"2.5.4.8":                    "st",
	"2.5.4.9":                    "street",
	"2.5.4.20":                   "telephonenumber",
	"2.5.4.22":                   "teletexterminalidentifier",
	"2.5.4.21":                   "telexnumber",
	"2.5.4.12":                   "title",
	"0.9.2342.19200300.100.1.1":  "uid",
	"2.5.4.50":                   "uniquemember",
	"2.5.4.35":                   "userpassword",
	"2.5.4.24":                   "x121address",
}

func oidString(oid asn1.ObjectIdentifier) (outbound string) {
	var (
		ok bool
	)
	switch outbound, ok = oidMap[oid.String()]; {
	case ok:
		return outbound
	}

	return oid.String()
}

func (r *attrDN) String() (outbound string) {
	switch {
	case r == nil:
		return
	}

	var (
		seq     = r.Name.ToRDNSequence()
		dnParts = make([]string, len(seq), len(seq))
	)
	for a, b := range r.Name.ToRDNSequence() {
		var (
			attrs []string
		)
		for _, d := range b {
			var (
				valStr string
			)
			switch value := d.Value.(type) {
			case asn1.RawValue:
				valStr = string(value.Bytes)
			case string:
				valStr = value
			default:
				panic("unhandled type")
			}

			attrs = append(attrs, fmt.Sprintf("%s=%s", oidString(d.Type), valStr))
		}

		dnParts[a] = strings.Join(attrs, "+")
	}

	return strings.Join(dnParts, ",")
}
func (r *attrDN) LegacyString() (outbound string) { return r.Name.String() }

func (r *attrDN) Parse(inbound string) (err error) {
	var (
		interim *pkix.Name
	)
	switch interim, err = x500dn.ParseDN(inbound); {
	case err != nil:
		return
	}

	r.Name = *interim

	return
}
func (r *attrDN) MarshalJSON() (outbound []byte, err error) { return json.Marshal(r.String()) }

func (r *attrDN) UnmarshalJSON(inbound []byte) (err error) {
	var (
		interim string
	)
	switch err = json.Unmarshal(inbound, &interim); {
	case err != nil:
		return
	}

	switch err = r.Parse(interim); {
	case err != nil:
		return
	}

	return
}

func (r *attrDN) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		switch err = r.Parse(value); {
		case err != nil:
			continue
		}

		return // return only first value
	}

	return
}
