package mod_ldap

import (
	"github.com/go-ldap/ldap/v3"
)

var (
	scopeIDMap = map[string]scopeIDType{
		// strconv.Itoa(ldap.ScopeBaseObject):   ldap.ScopeBaseObject,
		// strconv.Itoa(ldap.ScopeSingleLevel):  ldap.ScopeSingleLevel,
		// strconv.Itoa(ldap.ScopeWholeSubtree): ldap.ScopeWholeSubtree,
		// strconv.Itoa(ldap.ScopeChildren):     ldap.ScopeChildren,
		"base":                               ldap.ScopeBaseObject,
		"one":                                ldap.ScopeSingleLevel,
		"sub":                                ldap.ScopeWholeSubtree,
		"child":                              ldap.ScopeChildren,
	}
)
