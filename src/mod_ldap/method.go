package mod_ldap

import (
	"strings"

	"github.com/go-ldap/ldap/v3"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_slices"
	"rmm23/src/mod_strings"
)

func (r *Conf) SearchFn(fn func(fnBaseDN string, fnSearchResultType string, fnSearchResult *ldap.SearchResult) (fnErr error)) (err error) {
	switch err = r.dial(); {
	case err != nil:
		return
	}

	defer func() {
		_ = r.close()
	}()

	for _, b := range r.Domains {
		for _, d := range r.Settings {
			var (
				requestDN     = mod_strings.JoinStrings([]string{d.DN, b.DN}, ",", mod_slices.FlagFilterEmpty|mod_slices.FlagTrimSpace)
				searchRequest = ldap.NewSearchRequest(
					requestDN,          // Base DN
					d.Scope.Int(),      // Scope - search entire tree
					ldap.DerefAlways,   // Deref
					0,                  // Size limit (0 = no limit)
					0,                  // Time limit (0 = no limit)
					false,              // Types only
					d.Filter.String(),  // Filter - all objects
					[]string{"*", "+"}, // Attributes - all user and operational attributes
					nil,                // Controls
				)
				searchResult *ldap.SearchResult
			)
			switch searchResult, err = r.conn.Search(searchRequest); {
			case err != nil:
				return
			}

			switch err = fn(b.DN, d.Type, searchResult); {
			case err != nil:
				return
			}
		}
	}

	return
}

func (r *Conf) dial() (err error) {
	switch err = r.connect(); {
	case err != nil:
		return
	}

	switch err = r.bind(); {
	case err != nil:
		return
	}

	return
}
func (r *Conf) connect() (err error) {
	r.conn, err = ldap.DialURL(r.URL.String())

	return
}
func (r *Conf) bind() (err error) {
	switch {
	case r.conn == nil:
		return mod_errors.ENoConn
	}

	switch err = r.conn.Bind(r.URL.CleanUser()); {
	case err != nil:
		return
	case len(r.URL.CleanUsername()) == 0:
		return mod_errors.EAnonymousBind
	}

	return
}
func (r *Conf) close() (err error) {
	switch {
	case r.conn == nil:
		return mod_errors.ENoConn
	}

	return r.conn.Close()
}

func (d *attrSearchScope) UnmarshalJSON(data []byte) (err error) {
	switch value, ok := scopeIDMap[strings.Trim(string(data), " \"")]; {
	case !ok:
		return mod_errors.EINVAL
	default:
		*d = value

		return
	}
}
func (d *attrSearchScope) Int() (outbound int)        { return int(*d) }
func (d *attrSearchFilter) String() (outbound string) { return string(*d) }
