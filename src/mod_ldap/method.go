package mod_ldap

import (
	"errors"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_slices"
)

func (r *LDAPConfig) Search() (err error) {
	switch err = r.connect(); {
	case err != nil:
		return
	}

	defer func() {
		_ = r.close()
	}()

	switch err = r.bind(); {
	case errors.Is(err, mod_errors.EAnonymousBind):
	case err != nil:
		return
	}

	switch err = r.search(); {
	case err != nil:
		return
	}

	return
}

func (r *LDAPConfig) search() (err error) {
	for _, b := range r.Domains {
		switch {
		case b.SearchResults == nil:
			b.SearchResults = make(map[string]*ldap.SearchResult)
		}

		for _, d := range r.Settings {
			var (
				searchResult *ldap.SearchResult
			)

			switch searchResult, err = r.conn.Search(ldap.NewSearchRequest(
				mod_slices.JoinStrings([]string{d.DN.String(), b.DN.String()}, ",", mod_slices.FlagFilterEmpty), // Base DN
				d.Scope.Int(),      // Scope - search entire tree
				ldap.DerefAlways,   // Deref
				0,                  // Size limit (0 = no limit)
				0,                  // Time limit (0 = no limit)
				false,              // Types only
				d.Filter,           // Filter - all objects
				[]string{"*", "+"}, // Attributes - all user and operational attributes
				nil,                // Controls
			)); {
			case err != nil:
				return
			}

			b.SearchResults[d.Type] = searchResult
		}
	}

	return
}

func (r *AttrDN) String() string   { return string(*r) }
func (r *AttrUUID) String() string { return uuid.UUID(*r).String() }

func (r *LDAPConfig) connect() (err error) {
	r.conn, err = ldap.DialURL(r.URL.String())

	return
}
func (r *LDAPConfig) bind() (err error) {
	switch {
	case r.conn == nil:
		return mod_errors.ENoConn
	}

	switch err = r.conn.Bind(r.URL.CleanUsername(), r.URL.CleanPassword()); {
	case err != nil:
		return
	case len(r.URL.CleanUsername()) == 0:
		return mod_errors.EAnonymousBind
	}

	return
}
func (r *LDAPConfig) close() (err error) {
	switch {
	case r.conn == nil:
		return mod_errors.ENoConn
	}

	return r.conn.Close()
}

func (d *AttrSearchScope) UnmarshalJSON(data []byte) (err error) {
	switch value, ok := scopeIDMap[strings.Trim(string(data), " \"")]; {
	case !ok:
		return mod_errors.EINVAL
	default:
		*d = value

		return
	}
}
func (d *AttrSearchScope) Int() (outbound int) { return int(*d) }
