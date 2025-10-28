package mod_db

import (
	"encoding/json"
	"strings"
)

// ObjectClassEntry represents an objectClass with its schema-specific attributes
// Supports RFC 2307bis and other objectClass-specific attributes.
type ObjectClassEntry struct {
	Name       string                 `json:"name"`       // objectClass name (e.g., "posixAccount", "inetOrgPerson")
	Attributes map[string]interface{} `json:"attributes"` // objectClass-specific attributes
}

// ObjectClassList is a slice of ObjectClassEntry that can be stored in Redis.
type ObjectClassList []ObjectClassEntry

// Names returns all objectClass names.
func (ocl ObjectClassList) Names() []string {
	names := make([]string, len(ocl))
	for i, oc := range ocl {
		names[i] = oc.Name
	}

	return names
}

// HasClass checks if a specific objectClass exists.
func (ocl ObjectClassList) HasClass(name string) bool {
	for _, oc := range ocl {
		if strings.EqualFold(oc.Name, name) {
			return true
		}
	}

	return false
}

// GetClass returns the objectClass entry by name.
func (ocl ObjectClassList) GetClass(name string) *ObjectClassEntry {
	for _, oc := range ocl {
		if strings.EqualFold(oc.Name, name) {
			return &oc
		}
	}

	return nil
}

// GetAttribute retrieves an attribute from any objectClass.
func (ocl ObjectClassList) GetAttribute(className, attrName string) (interface{}, bool) {
	if oc := ocl.GetClass(className); oc != nil {
		val, ok := oc.Attributes[attrName]

		return val, ok
	}

	return nil, false
}

// MarshalJSON implements json.Marshaler.
func (ocl ObjectClassList) MarshalJSON() ([]byte, error) {
	return json.Marshal([]ObjectClassEntry(ocl))
}

// UnmarshalJSON implements json.Unmarshaler.
func (ocl *ObjectClassList) UnmarshalJSON(data []byte) error {
	var entries []ObjectClassEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}

	*ocl = entries

	return nil
}

// ToRedisTagValue converts objectClass list to Redis TAG-searchable format
// Format: "posixAccount|inetOrgPerson|shadowAccount".
func (ocl ObjectClassList) ToRedisTagValue() string {
	return strings.Join(ocl.Names(), "|")
}

// FromLDAPObjectClass creates ObjectClassList from LDAP objectClass attribute.
func FromLDAPObjectClass(classes []string, attributeMap map[string]map[string]interface{}) ObjectClassList {
	list := make(ObjectClassList, 0, len(classes))
	for _, className := range classes {
		entry := ObjectClassEntry{
			Name:       className,
			Attributes: make(map[string]interface{}),
		}
		// Copy class-specific attributes if provided
		if attrs, ok := attributeMap[className]; ok {
			for k, v := range attrs {
				entry.Attributes[k] = v
			}
		}

		list = append(list, entry)
	}

	return list
}
