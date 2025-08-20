package mod_url

import (
	"encoding/json/v2"
	"net/url"
	"strings"

	"rmm23/src/mod_errors"
)

func (r *URL) UnmarshalText(inbound []byte) (err error) {
	var (
		interim *url.URL
	)

	switch interim, err = url.Parse(string(inbound)); {
	case err != nil:
		return
	}

	r.URL = interim

	return
}

func (r *URL) MarshalText() (outbound []byte, err error) { return json.Marshal(r.String()) }

func (r *URL) CleanPath() (outbound string) { return strings.TrimPrefix(r.Path, "/") }
func (r *URL) CleanUser() (username string, password string) {
	return r.CleanUsername(), r.CleanPassword()
}
func (r *URL) CleanUsername() (outbound string) { return r.User.Username() }
func (r *URL) CleanPassword() (outbound string) {
	outbound, _ = r.User.Password()

	return
}

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
