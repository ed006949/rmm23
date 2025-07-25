package mod_ldap

import (
	"errors"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"

	"rmm23/src/l"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_slices"
)

func (r *LDAPConfig) Fetch() (err error) {
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

	switch err = r.parse(); {
	case err != nil:
		return
	}

	return
}

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

func (r *LDAPConfig) Unmarshal(inbound any) (err error) {
	return
}

func (r *LDAPConfig) search() (err error) {
	for _, b := range r.Domains {
		switch b.SearchResults {
		case nil:
			b.SearchResults = make(map[string]*ldap.SearchResult)
		}

		for _, d := range r.Settings {
			switch d.Type {
			case "domain", "host", "user", "group":
				var (
					newErr        error
					searchResult  *ldap.SearchResult
				)

				switch searchResult, newErr = r.conn.Search(ldap.NewSearchRequest(
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
				case newErr != nil:
					err = errors.Join(err, newErr)
				}

				b.SearchResults[d.Type] = searchResult
			}
		}
	}

	return
}
func (r *LDAPConfig) parse() (err error) {
	for _, b := range r.Domains {
		switch newErr := b.unmarshal(); {
		case newErr != nil:
			err = errors.Join(err, newErr)
		}
	}

	return
}

func (r *LDAPDomain) unmarshal() (err error) {
	r.Domain = &Element{}

	switch newErr := r.Domain.unmarshal(r.SearchResults["domain"]); {
	case newErr != nil:
		err = errors.Join(err, newErr)
		l.Z{l.E: err, l.M: "LDAP Unmarshal Domain", "DN": r.DN}.Warning()
	}

	r.Hosts = make(Elements)

	switch newErr := r.Hosts.unmarshal(r.SearchResults["hosts"]); {
	case newErr != nil:
		err = errors.Join(err, newErr)
		l.Z{l.E: err, l.M: "LDAP Unmarshal Hosts", "DN": r.DN}.Warning()
	}

	r.Users = make(Elements)

	switch newErr := r.Users.unmarshal(r.SearchResults["users"]); {
	case newErr != nil:
		err = errors.Join(err, newErr)
		l.Z{l.E: err, l.M: "LDAP Unmarshal Users", "DN": r.DN}.Warning()
	}

	r.Groups = make(Elements)

	switch newErr := r.Groups.unmarshal(r.SearchResults["groups"]); {
	case newErr != nil:
		err = errors.Join(err, newErr)
		l.Z{l.E: err, l.M: "LDAP Unmarshal Groups", "DN": r.DN}.Warning()
	}

	return
}
func (r *Element) unmarshal(inbound *ldap.SearchResult) (err error) {
	for _, entry := range inbound.Entries {
		var (
			interim Element
		)

		switch newErr := UnmarshalEntry(entry, &interim); {
		case newErr != nil:
			err = errors.Join(err, newErr)
			l.Z{l.E: err, l.M: "LDAP Unmarshal", "DN": entry.DN}.Warning()
		}

		*r = interim
	}

	return
}
func (r Elements) unmarshal(inbound *ldap.SearchResult) (err error) {
	for _, entry := range inbound.Entries {
		var (
			interim Element
		)

		switch newErr := UnmarshalEntry(entry, &interim); {
		case newErr != nil:
			err = errors.Join(err, newErr)
			l.Z{l.E: err, l.M: "LDAP Unmarshal", "DN": entry.DN}.Warning()
		}

		r[interim.DN] = &interim
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

func (d *scopeIDType) UnmarshalJSON(data []byte) (err error) {
	switch value, ok := scopeIDMap[strings.TrimSpace(string(data))]; {
	case !ok:
		return mod_errors.EINVAL
	default:
		*d = value

		return
	}
}
func (d *scopeIDType) Int() (outbound int) { return int(*d) }
