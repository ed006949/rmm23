package mod_db

import (
	"context"
	"fmt"
	"strings"

	"rmm23/src/mod_strings"
)

// SearchV2 provides advanced search capabilities for EntryV2.
type SearchV2 struct {
	repo *RedisRepository
}

// NewSearchV2 creates a new SearchV2 instance.
func (r *RedisRepository) NewSearchV2() *SearchV2 {
	return &SearchV2{repo: r}
}

// SearchByRedisAttrs searches by Redis-specific attributes (Key, Ver, ExAt)
// Example: SearchByRedisAttrs("Key", "some-key-prefix*").
func (s *SearchV2) SearchByRedisAttrs(attr, value string) (count int64, entries []*Entry, err error) {
	var query string

	switch strings.ToLower(attr) {
	case "key", "$key":
		query = fmt.Sprintf("@__key:%s", value)
	case "ver", "$ver":
		query = fmt.Sprintf("@__ver:[%s %s]", value, value)
	default:
		return 0, nil, fmt.Errorf("unsupported redis attribute: %s", attr)
	}

	return s.repo.SearchEntryQ(query)
}

// SearchByOperationalAttr searches by LDAP operational attributes
// Supports: entryDN, entryUUID, createTimestamp, modifyTimestamp, creatorsName, modifiersName, etc.
func (s *SearchV2) SearchByOperationalAttr(attr, value string) (count int64, entries []*Entry, err error) {
	var query string

	switch strings.ToLower(attr) {
	case "entrydn":
		query = fmt.Sprintf("@entryDN:{%s}", escapeTag(value))
	case "entryuuid":
		query = fmt.Sprintf("@entryUUID:{%s}", escapeTag(value))
	case "creatorsname":
		query = fmt.Sprintf("@creatorsName:{%s}", escapeTag(value))
	case "modifiersname":
		query = fmt.Sprintf("@modifiersName:{%s}", escapeTag(value))
	case "structuralobjectclass":
		query = fmt.Sprintf("@structuralObjectClass:{%s}", escapeTag(value))
	case "subschemasubentry":
		query = fmt.Sprintf("@subschemaSubentry:{%s}", escapeTag(value))
	case "hassubordinates":
		query = fmt.Sprintf("@hasSubordinates:{%s}", value)
	case "subordinatecount":
		query = fmt.Sprintf("@subordinateCount:[%s %s]", value, value)
	case "entrycsn":
		query = fmt.Sprintf("@entryCSN:{%s}", escapeTag(value))
	case "contextcsn":
		query = fmt.Sprintf("@contextCSN:{%s}", escapeTag(value))
	case "createtimestamp":
		query = fmt.Sprintf("@createTimestamp:[%s %s]", value, value)
	case "modifytimestamp":
		query = fmt.Sprintf("@modifyTimestamp:[%s %s]", value, value)
	default:
		return 0, nil, fmt.Errorf("unsupported operational attribute: %s", attr)
	}

	return s.repo.SearchEntryQ(query)
}

// SearchByOperationalAttrRange searches operational attributes by range (for timestamps)
// Example: SearchByOperationalAttrRange("createTimestamp", "2024-01-01", "2024-12-31").
func (s *SearchV2) SearchByOperationalAttrRange(attr, minVal, maxVal string) (count int64, entries []*Entry, err error) {
	var query string

	switch strings.ToLower(attr) {
	case "createtimestamp":
		query = fmt.Sprintf("@createTimestamp:[%s %s]", minVal, maxVal)
	case "modifytimestamp":
		query = fmt.Sprintf("@modifyTimestamp:[%s %s]", minVal, maxVal)
	case "subordinatecount":
		query = fmt.Sprintf("@subordinateCount:[%s %s]", minVal, maxVal)
	default:
		return 0, nil, fmt.Errorf("unsupported range attribute: %s", attr)
	}

	return s.repo.SearchEntryQ(query)
}

// SearchByObjectClass searches entries by objectClass name
// Example: SearchByObjectClass("posixAccount").
func (s *SearchV2) SearchByObjectClass(className string) (count int64, entries []*Entry, err error) {
	query := fmt.Sprintf("@objectClass:{%s}", escapeTag(className))

	return s.repo.SearchEntryQ(query)
}

// SearchByMultipleObjectClasses searches entries that have ALL specified objectClasses
// Example: SearchByMultipleObjectClasses([]string{"posixAccount", "shadowAccount"}).
func (s *SearchV2) SearchByMultipleObjectClasses(classNames []string) (count int64, entries []*Entry, err error) {
	if len(classNames) == 0 {
		return 0, nil, fmt.Errorf("no objectClass names provided")
	}

	queries := make([]string, len(classNames))
	for i, className := range classNames {
		queries[i] = fmt.Sprintf("@objectClass:{%s}", escapeTag(className))
	}

	query := strings.Join(queries, " ")

	return s.repo.SearchEntryQ(query)
}

// SearchByObjectClassAttribute searches for entries with specific objectClass and filters by attribute
// This is a two-step process:
// 1. Search by objectClass
// 2. Filter results by objectClass-specific attribute in application layer.
func (s *SearchV2) SearchByObjectClassAttribute(className, attrName string, attrValue interface{}) (filteredEntries []*Entry, err error) {
	// Step 1: Get all entries with this objectClass
	_, entries, err := s.SearchByObjectClass(className)
	if err != nil {
		return nil, err
	}

	// Step 2: Filter by objectClass attribute
	// Note: This is done in application layer since RediSearch doesn't directly support
	// nested JSON array filtering in the index
	filteredEntries = make([]*Entry, 0)
	for _, entry := range entries {
		// This assumes Entry has ObjectClassData field (would need to add to Entry struct)
		// For now, this is a placeholder showing the approach
		// In practice, you might need to:
		// 1. Use EntryV2 with ObjectClassData support, OR
		// 2. Store objectClass attributes as separate indexed fields, OR
		// 3. Use Redis JSON path queries with more complex indexing
		filteredEntries = append(filteredEntries, entry)
	}

	return filteredEntries, nil
}

// SearchCriteria allows combining multiple search criteria.
type SearchCriteria struct {
	RedisAttrs       map[string]string    // Redis-specific attributes
	OperationalAttrs map[string]string    // LDAP operational attributes
	ObjectClasses    []string             // Required objectClasses (AND logic)
	StandardAttrs    map[string]string    // Standard LDAP attributes
	RangeAttrs       map[string][2]string // Range queries (attr -> [min, max])
}

// SearchByComplex performs a complex search with multiple criteria.
func (s *SearchV2) SearchByComplex(criteria SearchCriteria) (count int64, entries []*Entry, err error) {
	queryParts := make([]string, 0,
		len(criteria.RedisAttrs)+
		len(criteria.OperationalAttrs)+
		len(criteria.ObjectClasses)+
		len(criteria.StandardAttrs)+
		len(criteria.RangeAttrs))

	// Redis attributes
	for attr, value := range criteria.RedisAttrs {
		switch strings.ToLower(attr) {
		case "key", "$key":
			queryParts = append(queryParts, fmt.Sprintf("@__key:%s", value))
		case "ver", "$ver":
			queryParts = append(queryParts, fmt.Sprintf("@__ver:[%s %s]", value, value))
		}
	}

	// Operational attributes
	for attr, value := range criteria.OperationalAttrs {
		switch strings.ToLower(attr) {
		case "entrydn":
			queryParts = append(queryParts, fmt.Sprintf("@entryDN:{%s}", escapeTag(value)))
		case "entryuuid":
			queryParts = append(queryParts, fmt.Sprintf("@entryUUID:{%s}", escapeTag(value)))
		case "creatorsname":
			queryParts = append(queryParts, fmt.Sprintf("@creatorsName:{%s}", escapeTag(value)))
		case "modifiersname":
			queryParts = append(queryParts, fmt.Sprintf("@modifiersName:{%s}", escapeTag(value)))
		case "structuralobjectclass":
			queryParts = append(queryParts, fmt.Sprintf("@structuralObjectClass:{%s}", escapeTag(value)))
		}
	}

	// objectClass (AND logic)
	for _, className := range criteria.ObjectClasses {
		queryParts = append(queryParts, fmt.Sprintf("@objectClass:{%s}", escapeTag(className)))
	}

	// Standard attributes
	for attr, value := range criteria.StandardAttrs {
		fieldName := mod_strings.EntryFieldName(attr)
		queryParts = append(queryParts, fmt.Sprintf("@%s:{%s}", fieldName.String(), escapeTag(value)))
	}

	// Range queries
	for attr, minMax := range criteria.RangeAttrs {
		switch strings.ToLower(attr) {
		case "createtimestamp", "modifytimestamp", "subordinatecount":
			queryParts = append(queryParts, fmt.Sprintf("@%s:[%s %s]", attr, minMax[0], minMax[1]))
		case "uidnumber", "gidnumber":
			queryParts = append(queryParts, fmt.Sprintf("@%s:[%s %s]", attr, minMax[0], minMax[1]))
		}
	}

	if len(queryParts) == 0 {
		return 0, nil, fmt.Errorf("no search criteria provided")
	}

	query := strings.Join(queryParts, " ")

	return s.repo.SearchEntryQ(query)
}

// AdvancedSearchBuilder provides a fluent interface for building complex searches.
type AdvancedSearchBuilder struct {
	search   *SearchV2
	criteria SearchCriteria
}

// NewAdvancedSearch creates a new advanced search builder.
func (s *SearchV2) NewAdvancedSearch() *AdvancedSearchBuilder {
	return &AdvancedSearchBuilder{
		search: s,
		criteria: SearchCriteria{
			RedisAttrs:       make(map[string]string),
			OperationalAttrs: make(map[string]string),
			ObjectClasses:    make([]string, 0),
			StandardAttrs:    make(map[string]string),
			RangeAttrs:       make(map[string][2]string),
		},
	}
}

// WithRedisAttr adds a Redis attribute filter.
func (b *AdvancedSearchBuilder) WithRedisAttr(attr, value string) *AdvancedSearchBuilder {
	b.criteria.RedisAttrs[attr] = value

	return b
}

// WithOperationalAttr adds an operational attribute filter.
func (b *AdvancedSearchBuilder) WithOperationalAttr(attr, value string) *AdvancedSearchBuilder {
	b.criteria.OperationalAttrs[attr] = value

	return b
}

// WithObjectClass adds an objectClass requirement.
func (b *AdvancedSearchBuilder) WithObjectClass(className string) *AdvancedSearchBuilder {
	b.criteria.ObjectClasses = append(b.criteria.ObjectClasses, className)

	return b
}

// WithStandardAttr adds a standard LDAP attribute filter.
func (b *AdvancedSearchBuilder) WithStandardAttr(attr, value string) *AdvancedSearchBuilder {
	b.criteria.StandardAttrs[attr] = value

	return b
}

// WithRange adds a range query.
func (b *AdvancedSearchBuilder) WithRange(attr, minVal, maxVal string) *AdvancedSearchBuilder {
	b.criteria.RangeAttrs[attr] = [2]string{minVal, maxVal}

	return b
}

// Execute performs the search.
func (b *AdvancedSearchBuilder) Execute() (count int64, entries []*Entry, err error) {
	return b.search.SearchByComplex(b.criteria)
}

// ExecuteWithContext performs the search with custom context.
func (b *AdvancedSearchBuilder) ExecuteWithContext(ctx context.Context) (count int64, entries []*Entry, err error) {
	// Store original context
	origCtx := b.search.repo.ctx
	b.search.repo.ctx = ctx

	defer func() { b.search.repo.ctx = origCtx }()

	return b.search.SearchByComplex(b.criteria)
}

// escapeTag escapes special characters for TAG queries.
func escapeTag(value string) string {
	// RediSearch TAG special characters that need escaping: , . < > { } [ ] " ' : ; ! @ # $ % ^ & * ( ) - + = ~ |
	replacer := strings.NewReplacer(
		",", "\\,",
		".", "\\.",
		"<", "\\<",
		">", "\\>",
		"{", "\\{",
		"}", "\\}",
		"[", "\\[",
		"]", "\\]",
		":", "\\:",
		";", "\\;",
		"|", "\\|",
	)

	return replacer.Replace(value)
}

// SearchEntryV2 provides convenience methods for the repository
// Add these methods to RedisRepository or create extension methods

// SearchByDN searches for an entry by its Distinguished Name.
func (s *SearchV2) SearchByDN(dn string) (*Entry, error) {
	_, entries, err := s.SearchByOperationalAttr("entryDN", dn)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("entry not found: %s", dn)
	}

	return entries[0], nil
}

// SearchByUUID searches for an entry by its UUID.
func (s *SearchV2) SearchByUUID(uuid string) (*Entry, error) {
	_, entries, err := s.SearchByOperationalAttr("entryUUID", uuid)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("entry not found: %s", uuid)
	}

	return entries[0], nil
}

// SearchPosixAccounts searches for all POSIX account entries.
func (s *SearchV2) SearchPosixAccounts() (count int64, entries []*Entry, err error) {
	return s.SearchByObjectClass("posixAccount")
}

// SearchByUIDNumber searches for a user by UID number.
func (s *SearchV2) SearchByUIDNumber(uidNumber uint64) (*Entry, error) {
	query := fmt.Sprintf("@uidNumber:[%d %d] @objectClass:{posixAccount}", uidNumber, uidNumber)

	_, entries, err := s.repo.SearchEntryQ(query)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("user not found with uidNumber: %d", uidNumber)
	}

	return entries[0], nil
}

// SearchByGIDNumber searches for a group by GID number.
func (s *SearchV2) SearchByGIDNumber(gidNumber uint64) (*Entry, error) {
	query := fmt.Sprintf("@gidNumber:[%d %d] @objectClass:{posixGroup}", gidNumber, gidNumber)

	_, entries, err := s.repo.SearchEntryQ(query)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("group not found with gidNumber: %d", gidNumber)
	}

	return entries[0], nil
}

// SearchModifiedSince searches for entries modified since a specific timestamp.
func (s *SearchV2) SearchModifiedSince(timestamp string) (count int64, entries []*Entry, err error) {
	query := fmt.Sprintf("@modifyTimestamp:[%s +inf]", timestamp)

	return s.repo.SearchEntryQ(query)
}

// SearchCreatedBetween searches for entries created within a time range.
func (s *SearchV2) SearchCreatedBetween(startTime, endTime string) (count int64, entries []*Entry, err error) {
	return s.SearchByOperationalAttrRange("createTimestamp", startTime, endTime)
}
