package mod_net

import (
	"net/netip"
	"net/url"

	"rmm23/src/mod_crypto"
)

type URL struct{ *url.URL }

type OpenVPN struct {
	Listen   []netip.AddrPort `xml:"listen,omitempty"`
	Hostname []string         `xml:"hostname,omitempty"`
}
type CiscoVPN struct {
	Listen   []netip.AddrPort `xml:"listen,omitempty"`
	Hostname []string         `xml:"hostname,omitempty"`
}
type InterimHost struct {
	ASN      uint32 `xml:"asn,omitempty"`
	Hostname string `xml:"hostname,omitempty"`
	Auth     mod_crypto.AuthDB
}
