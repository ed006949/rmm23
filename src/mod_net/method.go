package mod_net

import (
	"encoding/json"
	"encoding/xml"
	"net/url"
	"strings"

	"rmm23/src/mod_bools"
	"rmm23/src/mod_errors"
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
	return xml.Attr{Name: name, Value: r.URL.String()}, nil
}

func (r *URL) UnmarshalJSON(inbound []byte) (err error) {
	var (
		interim    string
		interimURL *url.URL
	)
	switch err = json.Unmarshal(inbound, &interim); {
	case err != nil:
		return
	}

	switch interimURL, err = url.Parse(interim); {
	case err != nil:
		return
	}

	r.URL = interimURL

	return
}
func (r *URL) MarshalJSON() ([]byte, error) { return json.Marshal(r.URL.String()) }

func (r *URL) CleanPath() (outbound string) { return strings.TrimPrefix(r.URL.Path, "/") }
func (r *URL) CleanUser() (username string, password string) {
	return r.CleanUsername(), r.CleanPassword()
}
func (r *URL) CleanUsername() (outbound string) { return r.User.Username() }
func (r *URL) CleanPassword() (outbound string) { return mod_bools.StripIfBool1(r.User.Password()) }

func (r *URL) RedisNetwork() (outbound string, err error) {
	switch outbound = r.Scheme; outbound {
	case "redis", "redis-sentinel":
		return "tcp", nil
	case "file":
		return "unix", nil
	default:
		return outbound, mod_errors.EUnknownScheme
	}
}
