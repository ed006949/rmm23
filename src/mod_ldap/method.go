package mod_ldap

import (
	"encoding/xml"
	"errors"

	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"

	"rmm23/src/l"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_slices"
)

func (r *Conf) Fetch() (err error) {
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

func (r *Conf) search() (err error) {
	for _, b := range r.Domain {
		switch {
		case b.searchResults == nil:
			b.searchResults = make(map[string]*ldap.SearchResult)
		}
		for _, d := range r.Settings {
			switch d.Type {
			case "domain":
				var (
					newErr        error
					searchResult  *ldap.SearchResult
					searchRequest = ldap.NewSearchRequest(
						mod_slices.Join([]string{d.DN.String(), b.DN.String()}, ",", mod_slices.FlagFilterEmpty), // Base DN
						ldap.ScopeBaseObject, // Scope - search entire tree
						ldap.DerefAlways,     // Deref
						0,                    // Size limit (0 = no limit)
						0,                    // Time limit (0 = no limit)
						false,                // Types only
						d.Filter,             // Filter - all objects
						[]string{"*", "+"},   // Attributes - all user and operational attributes
						nil,                  // Controls
					)
				)
				switch searchResult, newErr = r.conn.Search(searchRequest); {
				case newErr != nil:
					err = errors.Join(err, newErr)
					l.Z{l.E: err, l.M: "LDAP Search", "DN": searchRequest.BaseDN}.Warning()
				}
				b.searchResults[d.Type] = searchResult

			case "hosts", "users", "groups":
				var (
					newErr        error
					searchResult  *ldap.SearchResult
					searchRequest = ldap.NewSearchRequest(
						mod_slices.Join([]string{d.DN.String(), b.DN.String()}, ",", mod_slices.FlagFilterEmpty), // Base DN
						ldap.ScopeWholeSubtree, // Scope - search entire tree
						ldap.DerefAlways,       // Deref
						0,                      // Size limit (0 = no limit)
						0,                      // Time limit (0 = no limit)
						false,                  // Types only
						d.Filter,               // Filter - all objects
						[]string{"*", "+"},     // Attributes - all user and operational attributes
						nil,                    // Controls
					)
				)
				switch searchResult, newErr = r.conn.Search(searchRequest); {
				case newErr != nil:
					err = errors.Join(err, newErr)
					l.Z{l.E: err, l.M: "LDAP Search", "DN": searchRequest.BaseDN}.Warning()
				}
				b.searchResults[d.Type] = searchResult

			default:
				// 		WTF
			}
		}
	}
	return
}
func (r *Conf) parse() (err error) {
	for _, b := range r.Domain {
		switch newErr := b.unmarshal(); {
		case newErr != nil:
			err = errors.Join(err, newErr)
			l.Z{l.E: err, l.M: "LDAP Unmarshal", "DN": b.DN, "URL": r.URL.Redacted()}.Warning()
		}
	}

	return
}

// UnmarshalEntry
func (r *ConfDomain) unmarshal() (err error) {
	r.Domain = &ElementDomain{}
	switch newErr := r.Domain.unmarshal(r.searchResults["domain"]); {
	case newErr != nil:
		err = errors.Join(err, newErr)
		l.Z{l.E: err, l.M: "LDAP Unmarshal Domain", "DN": r.DN}.Warning()
	}
	r.Hosts = make(ElementHosts)
	switch newErr := r.Hosts.unmarshal(r.searchResults["hosts"]); {
	case newErr != nil:
		err = errors.Join(err, newErr)
		l.Z{l.E: err, l.M: "LDAP Unmarshal Hosts", "DN": r.DN}.Warning()
	}
	r.Users = make(ElementUsers)
	switch newErr := r.Users.unmarshal(r.searchResults["users"]); {
	case newErr != nil:
		err = errors.Join(err, newErr)
		l.Z{l.E: err, l.M: "LDAP Unmarshal Users", "DN": r.DN}.Warning()
	}
	r.Groups = make(ElementGroups)
	switch newErr := r.Groups.unmarshal(r.searchResults["groups"]); {
	case newErr != nil:
		err = errors.Join(err, newErr)
		l.Z{l.E: err, l.M: "LDAP Unmarshal Groups", "DN": r.DN}.Warning()
	}
	return
}
func (r *ElementDomain) unmarshal(inbound *ldap.SearchResult) (err error) {
	for _, entry := range inbound.Entries {
		var (
			interim ElementDomain
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
func (r ElementHosts) unmarshal(inbound *ldap.SearchResult) (err error) {
	for _, entry := range inbound.Entries {
		var (
			interim ElementHost
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
func (r ElementUsers) unmarshal(inbound *ldap.SearchResult) (err error) {
	for _, entry := range inbound.Entries {
		var (
			interim ElementUser
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
func (r ElementGroups) unmarshal(inbound *ldap.SearchResult) (err error) {
	for _, entry := range inbound.Entries {
		var (
			interim ElementGroup
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

// find if exist

func (r *Conf) FindUser(inbound AttrDN) (outbound *ElementUser) {
	for _, b := range r.Domain {
		switch value, ok := b.Users[inbound]; {
		case ok:
			return value
		}
	}
	return
}
func (r *Conf) FindGroup(inbound AttrDN) (outbound *ElementGroup) {
	for _, b := range r.Domain {
		switch value, ok := b.Groups[inbound]; {
		case ok:
			return value
		}
	}
	return
}
func (r *Conf) FindHost(inbound AttrDN) (outbound *ElementHost) {
	for _, b := range r.Domain {
		switch value, ok := b.Hosts[inbound]; {
		case ok:
			return value
		}
	}
	return
}

// check if exist

// func (r *Conf) IsUser(inbound AttrDN) (outbound bool)  { return r.FindUser(inbound) != nil }
// func (r *Conf) IsGroup(inbound AttrDN) (outbound bool) { return r.FindGroup(inbound) != nil }
// func (r *Conf) IsHost(inbound AttrDN) (outbound bool)  { return r.FindHost(inbound) != nil }

func (r *Conf) AddUser(inbound *ElementUser) (err error) {
	switch {
	case r.FindUser(inbound.DN) != nil:
		l.Z{l.E: mod_errors.EEXIST, "DN": inbound.DN}.Warning()
		return mod_errors.EEXIST
	}
	return
}

// conn
func (r *Conf) connect() (err error) {
	r.conn, err = ldap.DialURL(r.URL.String())
	return
}
func (r *Conf) bind() (err error) {
	switch {
	case r.conn == nil:
		return mod_errors.ENoConn
	}

	switch err = r.conn.Bind(r.URL.CleanUsername(), r.URL.CleanPassword()); {
	case err != nil:
		return
	case len(r.URL.CleanUsername()) == 0:
		// return EAnonymousBind
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

// XML
func (r *AttrDN) UnmarshalXMLAttr(attr xml.Attr) (err error) {
	switch _, err = ldap.ParseDN(attr.Value); {
	case err == nil:
		*r = AttrDN(attr.Value)
	}
	return
}
func (r *AttrDN) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: r.String()}, nil
}

// String
func (r *AttrDN) String() string { return string(*r) }
func (r *AttrDNs) String() (outbound []string) {
	for _, b := range *r {
		outbound = append(outbound, b.String())
	}
	return
}

func (r *AttrUUID) String() string { return uuid.UUID(*r).String() }
