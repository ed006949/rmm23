package mod_ldap

import "fmt"

// GetAttributeType retrieves an attribute type by OID or name.
func (s *Schema) GetAttributeType(nameOrOID string) (*AttributeType, bool) {
	attr, ok := s.AttributeTypes[nameOrOID]

	return attr, ok
}

// GetObjectClass retrieves an object class by OID or name.
func (s *Schema) GetObjectClass(nameOrOID string) (*ObjectClass, bool) {
	class, ok := s.ObjectClasses[nameOrOID]

	return class, ok
}

// GetSyntax retrieves a syntax by OID.
func (s *Schema) GetSyntax(oid string) (*Syntax, bool) {
	syntax, ok := s.Syntaxes[oid]

	return syntax, ok
}

// GetMatchingRule retrieves a matching rule by OID or name.
func (s *Schema) GetMatchingRule(nameOrOID string) (*MatchingRule, bool) {
	rule, ok := s.MatchingRules[nameOrOID]

	return rule, ok
}

// RegisterAttributeType adds an attribute type to the schema.
func (s *Schema) RegisterAttributeType(attr *AttributeType) {
	s.AttributeTypes[attr.OID] = attr
	for _, name := range attr.Names {
		s.AttributeTypes[name] = attr
	}
}

// RegisterObjectClass adds an object class to the schema.
func (s *Schema) RegisterObjectClass(class *ObjectClass) {
	s.ObjectClasses[class.OID] = class
	for _, name := range class.Names {
		s.ObjectClasses[name] = class
	}
}

// RegisterSyntax adds a syntax to the schema.
func (s *Schema) RegisterSyntax(syntax *Syntax) {
	s.Syntaxes[syntax.OID] = syntax
}

// RegisterMatchingRule adds a matching rule to the schema.
func (s *Schema) RegisterMatchingRule(rule *MatchingRule) {
	s.MatchingRules[rule.OID] = rule
	for _, name := range rule.Names {
		s.MatchingRules[name] = rule
	}
}

// ValidateEntry validates an LDAP entry against the schema.
func (s *Schema) ValidateEntry(objectClasses []string, attributes map[string][]string) error {
	if len(objectClasses) == 0 {
		return fmt.Errorf("entry must have at least one objectClass")
	}

	// Collect all required and allowed attributes from object classes
	requiredAttrs := make(map[string]bool)
	allowedAttrs := make(map[string]bool)

	for _, className := range objectClasses {
		class, ok := s.GetObjectClass(className)
		if !ok {
			return fmt.Errorf("unknown object class: %s", className)
		}

		// Add required attributes
		for _, attr := range class.Must {
			requiredAttrs[attr] = true
			allowedAttrs[attr] = true
		}

		// Add optional attributes
		for _, attr := range class.May {
			allowedAttrs[attr] = true
		}

		// Recursively add attributes from superior classes
		if err := s.addSuperClassAttributes(class, requiredAttrs, allowedAttrs); err != nil {
			return err
		}
	}

	// Check that all required attributes are present
	for requiredAttr := range requiredAttrs {
		found := false

		for attrName := range attributes {
			if s.attributeMatches(attrName, requiredAttr) {
				found = true

				break
			}
		}

		if !found {
			return fmt.Errorf("required attribute missing: %s", requiredAttr)
		}
	}

	// Check that all present attributes are allowed
	for attrName := range attributes {
		allowed := false

		for allowedAttr := range allowedAttrs {
			if s.attributeMatches(attrName, allowedAttr) {
				allowed = true

				break
			}
		}

		if !allowed {
			return fmt.Errorf("attribute not allowed by object classes: %s", attrName)
		}

		// Validate attribute type exists and check single-value constraint
		attr, ok := s.GetAttributeType(attrName)
		if !ok {
			return fmt.Errorf("unknown attribute type: %s", attrName)
		}

		if attr.SingleValue && len(attributes[attrName]) > 1 {
			return fmt.Errorf("attribute %s is single-valued but has multiple values", attrName)
		}
	}

	return nil
}

// addSuperClassAttributes recursively adds attributes from superior classes.
func (s *Schema) addSuperClassAttributes(class *ObjectClass, required, allowed map[string]bool) error {
	for _, superClassName := range class.SuperClasses {
		superClass, ok := s.GetObjectClass(superClassName)
		if !ok {
			return fmt.Errorf("unknown superior object class: %s", superClassName)
		}

		for _, attr := range superClass.Must {
			required[attr] = true
			allowed[attr] = true
		}

		for _, attr := range superClass.May {
			allowed[attr] = true
		}

		// Recursively process superior classes
		if err := s.addSuperClassAttributes(superClass, required, allowed); err != nil {
			return err
		}
	}

	return nil
}

// attributeMatches checks if two attribute names refer to the same attribute
// (considering aliases and OIDs).
func (s *Schema) attributeMatches(name1, name2 string) bool {
	if name1 == name2 {
		return true
	}

	attr1, ok1 := s.GetAttributeType(name1)
	attr2, ok2 := s.GetAttributeType(name2)

	if ok1 && ok2 {
		return attr1.OID == attr2.OID
	}

	return false
}

// GetAllowedAttributes returns all attributes allowed for a set of object classes.
func (s *Schema) GetAllowedAttributes(objectClasses []string) (required, optional []string, err error) {
	requiredMap := make(map[string]bool)
	optionalMap := make(map[string]bool)

	for _, className := range objectClasses {
		class, ok := s.GetObjectClass(className)
		if !ok {
			return nil, nil, fmt.Errorf("unknown object class: %s", className)
		}

		for _, attr := range class.Must {
			requiredMap[attr] = true
		}

		for _, attr := range class.May {
			if !requiredMap[attr] {
				optionalMap[attr] = true
			}
		}

		// Add attributes from superior classes
		tempRequired := make(map[string]bool)

		tempOptional := make(map[string]bool)
		if err := s.addSuperClassAttributes(class, tempRequired, tempOptional); err != nil {
			return nil, nil, err
		}

		for attr := range tempRequired {
			requiredMap[attr] = true
		}

		for attr := range tempOptional {
			if !requiredMap[attr] {
				optionalMap[attr] = true
			}
		}
	}

	required = make([]string, 0, len(requiredMap))
	for attr := range requiredMap {
		required = append(required, attr)
	}

	optional = make([]string, 0, len(optionalMap))
	for attr := range optionalMap {
		optional = append(optional, attr)
	}

	return required, optional, nil
}

// Merge combines another schema into this schema.
func (s *Schema) Merge(other *Schema) {
	// Merge attribute types
	for oid, attr := range other.AttributeTypes {
		if _, exists := s.AttributeTypes[oid]; !exists {
			s.AttributeTypes[oid] = attr
		}
	}

	// Merge object classes
	for oid, class := range other.ObjectClasses {
		if _, exists := s.ObjectClasses[oid]; !exists {
			s.ObjectClasses[oid] = class
		}
	}

	// Merge syntaxes
	for oid, syntax := range other.Syntaxes {
		if _, exists := s.Syntaxes[oid]; !exists {
			s.Syntaxes[oid] = syntax
		}
	}

	// Merge matching rules
	for oid, rule := range other.MatchingRules {
		if _, exists := s.MatchingRules[oid]; !exists {
			s.MatchingRules[oid] = rule
		}
	}
}
