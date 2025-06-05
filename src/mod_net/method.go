package mod_net

import (
	"encoding/xml"
	"net/url"
	"strings"

	"rmm23/src/l"
)

func (r *URL) UnmarshalXMLAttr(attr xml.Attr) (err error) {
	var (
		interim *url.URL
	)

	switch interim, err = url.Parse(attr.Value); {
	case err != nil:
		return
	}

	r.URL = interim
	return
}

func (r *URL) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name:  name,
		Value: r.String(),
	}, nil
}

func (r *URL) CleanPath() (outbound string) { return strings.TrimPrefix(r.URL.Path, "/") }
func (r *URL) CleanUser() (username string, password string) {
	return r.CleanUsername(), r.CleanPassword()
}
func (r *URL) CleanUsername() (outbound string) { return r.URL.User.Username() }
func (r *URL) CleanPassword() (outbound string) { return l.StripIfBool1(r.URL.User.Password()) }
