package mod_ldap

// AttributeType represents an LDAP attribute type definition.
type AttributeType struct {
	OID          string   // Object Identifier
	Names        []string // Attribute names (primary and aliases)
	Description  string   // Human-readable description
	Obsolete     bool     // Whether this attribute is obsolete
	SuperType    string   // Parent attribute type (inheritance)
	Equality     string   // Equality matching rule OID
	Ordering     string   // Ordering matching rule OID
	Substring    string   // Substring matching rule OID
	Syntax       string   // Syntax OID
	SingleValue  bool     // Whether only one value is allowed
	Collective   bool     // Whether this is a collective attribute
	NoUserModify bool     // Whether users can modify this attribute
	Usage        string   // Attribute usage (userApplications, directoryOperation, etc.)
	MinLength    int      // Minimum length constraint
	MaxLength    int      // Maximum length constraint
}

// ObjectClass represents an LDAP object class definition.
type ObjectClass struct {
	OID          string   // Object Identifier
	Names        []string // Object class names (primary and aliases)
	Description  string   // Human-readable description
	Obsolete     bool     // Whether this object class is obsolete
	SuperClasses []string // Parent object classes
	Kind         string   // STRUCTURAL, AUXILIARY, or ABSTRACT
	Must         []string // Required attribute types
	May          []string // Optional attribute types
}

// Syntax represents an LDAP syntax definition.
type Syntax struct {
	OID         string // Object Identifier
	Description string // Human-readable description
	Obsolete    bool   // Whether this syntax is obsolete
}

// MatchingRule represents an LDAP matching rule definition.
type MatchingRule struct {
	OID         string   // Object Identifier
	Names       []string // Matching rule names
	Description string   // Human-readable description
	Obsolete    bool     // Whether this matching rule is obsolete
	Syntax      string   // Syntax OID that this rule applies to
}

// Schema represents a collection of LDAP schema definitions.
type Schema struct {
	AttributeTypes map[string]*AttributeType // Key: OID or name
	ObjectClasses  map[string]*ObjectClass   // Key: OID or name
	Syntaxes       map[string]*Syntax        // Key: OID
	MatchingRules  map[string]*MatchingRule  // Key: OID or name
}

// NewSchema creates a new empty schema.
func NewSchema() *Schema {
	return &Schema{
		AttributeTypes: make(map[string]*AttributeType),
		ObjectClasses:  make(map[string]*ObjectClass),
		Syntaxes:       make(map[string]*Syntax),
		MatchingRules:  make(map[string]*MatchingRule),
	}
}
