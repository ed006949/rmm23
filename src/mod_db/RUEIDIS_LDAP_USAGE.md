# RueiDIS LDAP Backend Usage Guide

This guide demonstrates how to use rueidis as a backend for LDAP entries storage with support for:
1. Redis-specific attributes (Key, Ver, ExAt)
2. LDAP operational attributes (entryDN, entryUUID, timestamps, etc.)
3. objectClass with schema-specific sub-attributes (RFC 2307bis support)
4. Advanced search capabilities across all attribute types

## Architecture Overview

```
Entry/EntryV2 (Go struct)
    ├── Redis-specific: Key, Ver, ExAt
    ├── LDAP Operational: entryDN, entryUUID, createTimestamp, etc.
    ├── objectClass with sub-attributes (ObjectClassList)
    └── Standard LDAP attributes: cn, uid, mail, etc.
         ↓
    rueidis/om (Object Mapper)
         ↓
    Redis + RediSearch
         ↓
    Searchable by: Redis attrs, operational attrs, objectClass, standard attrs
```

## Basic Setup

```go
import (
    "context"
    "time"

    "github.com/google/uuid"
    "github.com/redis/rueidis"

    "rmm23/src/mod_db"
    "rmm23/src/mod_dn"
    "rmm23/src/mod_time"
)

// Initialize Redis connection
ctx := context.Background()
client, err := rueidis.NewClient(rueidis.ClientOption{
    InitAddress: []string{"localhost:6379"},
})
if err != nil {
    panic(err)
}

// Create repository
repo := mod_db.NewRedisRepository(ctx, client)

// Create RediSearch index
err = repo.CreateEntryIndex()
if err != nil {
    panic(err)
}
```

## Creating Entries

### Example 1: POSIX User with RFC 2307bis

```go
// Create a POSIX user entry
entry := &mod_db.EntryV2{
    // Redis-specific attributes
    Key:  uuid.NewSHA1(uuid.Nil, []byte("uid=jdoe,ou=users,dc=example,dc=com")).String(),
    Ver:  1,
    ExAt: time.Now().Add(365 * 24 * time.Hour),

    // LDAP Operational attributes
    EntryDN:   mod_dn.MustParse("uid=jdoe,ou=users,dc=example,dc=com"),
    EntryUUID: uuid.New(),
    CreateTimestamp: mod_time.Now(),
    ModifyTimestamp: mod_time.Now(),
    CreatorsName: mod_dn.MustParse("cn=admin,dc=example,dc=com"),
    ModifiersName: mod_dn.MustParse("cn=admin,dc=example,dc=com"),
    StructuralObjectClass: "inetOrgPerson",

    // Standard LDAP attributes
    UID:       "jdoe",
    UIDNumber: 1000,
    GIDNumber: 1000,
    CN:        "John Doe",
    SN:        "Doe",
    Mail:      []string{"jdoe@example.com"},
    HomeDirectory: "/home/jdoe",

    // Custom attributes
    Type:   mod_db.EntryTypeUser,
    Status: mod_db.EntryStatusReady,
    BaseDN: mod_dn.MustParse("dc=example,dc=com"),
}

// Add objectClasses with sub-attributes
entry.AddObjectClass("inetOrgPerson", map[string]interface{}{
    "displayName": "John Doe",
    "givenName":   "John",
    "title":       "Software Engineer",
})

entry.AddObjectClass("posixAccount", map[string]interface{}{
    "loginShell":   "/bin/bash",
    "gecos":        "John Doe,Room 123,x1234",
    "homeDirectory": "/home/jdoe",
})

entry.AddObjectClass("shadowAccount", map[string]interface{}{
    "shadowLastChange": 19000,
    "shadowMax":        99999,
    "shadowWarning":    7,
})

// Sync objectClass list
entry.SyncObjectClass()

// Save to Redis
err = repo.SaveEntry(entry)
```

### Example 2: POSIX Group

```go
entry := &mod_db.EntryV2{
    Key:  uuid.NewSHA1(uuid.Nil, []byte("cn=developers,ou=groups,dc=example,dc=com")).String(),
    Ver:  1,
    ExAt: time.Now().Add(365 * 24 * time.Hour),

    EntryDN:   mod_dn.MustParse("cn=developers,ou=groups,dc=example,dc=com"),
    EntryUUID: uuid.New(),
    CreateTimestamp: mod_time.Now(),
    ModifyTimestamp: mod_time.Now(),
    StructuralObjectClass: "posixGroup",

    CN:        "developers",
    GIDNumber: 5000,
    Member: []mod_dn.DN{
        mod_dn.MustParse("uid=jdoe,ou=users,dc=example,dc=com"),
        mod_dn.MustParse("uid=asmith,ou=users,dc=example,dc=com"),
    },
    Description: "Developer group",

    Type:   mod_db.EntryTypeGroup,
    Status: mod_db.EntryStatusReady,
}

entry.AddObjectClass("posixGroup", map[string]interface{}{
    "memberUid": []string{"jdoe", "asmith"},
})

entry.AddObjectClass("groupOfNames", map[string]interface{}{
    "businessCategory": "Engineering",
})

entry.SyncObjectClass()
err = repo.SaveEntry(entry)
```

### Example 3: Host Entry with Network Information

```go
entry := &mod_db.EntryV2{
    Key:  uuid.NewSHA1(uuid.Nil, []byte("cn=server01,ou=hosts,dc=example,dc=com")).String(),
    Ver:  1,
    ExAt: time.Now().Add(365 * 24 * time.Hour),

    EntryDN:   mod_dn.MustParse("cn=server01,ou=hosts,dc=example,dc=com"),
    EntryUUID: uuid.New(),
    CreateTimestamp: mod_time.Now(),
    StructuralObjectClass: "ipHost",

    CN: "server01",
    IPHostNumber: []netip.Prefix{
        netip.MustParsePrefix("192.168.1.10/32"),
        netip.MustParsePrefix("2001:db8::10/128"),
    },
    Description: "Production web server",

    // Host-specific attributes
    HostType:   "provider",
    HostASN:    65000,
    HostListen: netip.MustParseAddr("0.0.0.0"),

    Type:   mod_db.EntryTypeHost,
    Status: mod_db.EntryStatusReady,
}

entry.AddObjectClass("ipHost", nil)
entry.AddObjectClass("device", map[string]interface{}{
    "serialNumber": "SN12345",
    "l":            "Data Center 1",
})

entry.SyncObjectClass()
err = repo.SaveEntry(entry)
```

## Searching Entries

### 1. Search by Redis-Specific Attributes

```go
search := repo.NewSearchV2()

// Search by key prefix
count, entries, err := search.SearchByRedisAttrs("key", "prefix*")

// Search by version
count, entries, err = search.SearchByRedisAttrs("ver", "1")
```

### 2. Search by LDAP Operational Attributes

```go
search := repo.NewSearchV2()

// Search by entryDN
entry, err := search.SearchByDN("uid=jdoe,ou=users,dc=example,dc=com")

// Search by entryUUID
entry, err = search.SearchByUUID("550e8400-e29b-41d4-a716-446655440000")

// Search by creator
count, entries, err := search.SearchByOperationalAttr("creatorsName", "cn=admin,dc=example,dc=com")

// Search by structural object class
count, entries, err = search.SearchByOperationalAttr("structuralObjectClass", "inetOrgPerson")

// Search entries modified since timestamp
count, entries, err = search.SearchModifiedSince("1704067200000") // Unix timestamp in ms

// Search entries created in a time range
count, entries, err = search.SearchCreatedBetween("1704067200000", "1735689600000")
```

### 3. Search by objectClass

```go
search := repo.NewSearchV2()

// Search by single objectClass
count, entries, err := search.SearchByObjectClass("posixAccount")

// Search by multiple objectClasses (AND logic)
count, entries, err = search.SearchByMultipleObjectClasses([]string{
    "posixAccount",
    "shadowAccount",
    "inetOrgPerson",
})

// Search for POSIX users (convenience method)
count, entries, err = search.SearchPosixAccounts()
```

### 4. Search by Standard LDAP Attributes

```go
search := repo.NewSearchV2()

// Search by UID
count, entries, err := repo.SearchEntryFV(mod_strings.F_uid, "jdoe")

// Search by UID number
entry, err := search.SearchByUIDNumber(1000)

// Search by GID number
entry, err = search.SearchByGIDNumber(5000)

// Search by mail
count, entries, err = repo.SearchEntryFV(mod_strings.F_mail, "jdoe@example.com")

// Search by CN
count, entries, err = repo.SearchEntryFV(mod_strings.F_cn, "John Doe")
```

### 5. Search by objectClass Sub-Attributes

Since objectClass sub-attributes are stored in nested JSON, you have two approaches:

**Approach A: Two-step search (application-level filtering)**

```go
search := repo.NewSearchV2()

// Get all posixAccount entries, then filter by sub-attribute
filteredEntries, err := search.SearchByObjectClassAttribute(
    "posixAccount",
    "loginShell",
    "/bin/bash",
)
```

**Approach B: Index specific attributes as top-level fields**

For frequently searched objectClass attributes, promote them to top-level Entry fields:
- `loginShell` → add as top-level field in Entry struct
- Add to RediSearch index
- Search directly: `@loginShell:{/bin/bash}`

### 6. Complex Searches with Multiple Criteria

```go
search := repo.NewSearchV2()

// Using SearchCriteria struct
criteria := mod_db.SearchCriteria{
    OperationalAttrs: map[string]string{
        "structuralObjectClass": "inetOrgPerson",
    },
    ObjectClasses: []string{"posixAccount", "shadowAccount"},
    StandardAttrs: map[string]string{
        "cn": "John*",
    },
    RangeAttrs: map[string][2]string{
        "uidNumber": {"1000", "2000"},
    },
}

count, entries, err := search.SearchByComplex(criteria)

// Using fluent builder API
count, entries, err = search.NewAdvancedSearch().
    WithOperationalAttr("structuralObjectClass", "inetOrgPerson").
    WithObjectClass("posixAccount").
    WithObjectClass("shadowAccount").
    WithStandardAttr("cn", "John*").
    WithRange("uidNumber", "1000", "2000").
    Execute()
```

### 7. Advanced Query Examples

```go
search := repo.NewSearchV2()

// Find all users with SSH keys
query := "@sshPublicKey:{*} @objectClass:{posixAccount}"
count, entries, err := repo.SearchEntryQ(query)

// Find groups with more than 10 members (if you track member count)
query = "@objectClass:{posixGroup} @subordinateCount:[10 +inf]"
count, entries, err = repo.SearchEntryQ(query)

// Find entries modified by a specific admin in the last 30 days
thirtyDaysAgo := time.Now().Add(-30 * 24 * time.Hour).UnixMilli()
query = fmt.Sprintf("@modifiersName:{cn=admin,dc=example,dc=com} @modifyTimestamp:[%d +inf]", thirtyDaysAgo)
count, entries, err = repo.SearchEntryQ(query)

// Find all IPv4 hosts in a specific subnet (requires custom indexing)
query = "@objectClass:{ipHost} @ipHostNumber:{192.168.1.*}"
count, entries, err = repo.SearchEntryQ(query)
```

## Working with objectClass Sub-Attributes

```go
// Set objectClass sub-attribute
entry.SetObjectClassAttribute("posixAccount", "loginShell", "/bin/zsh")

// Get objectClass sub-attribute
value, ok := entry.GetObjectClassAttribute("shadowAccount", "shadowLastChange")
if ok {
    lastChange := value.(int64)
    fmt.Printf("Last password change: %d days since epoch\n", lastChange)
}

// Check if entry has specific objectClass
if entry.ObjectClassData.HasClass("shadowAccount") {
    // Handle shadow account specific logic
}

// Get all attributes of a specific objectClass
if oc := entry.ObjectClassData.GetClass("posixAccount"); oc != nil {
    for attrName, attrValue := range oc.Attributes {
        fmt.Printf("%s: %v\n", attrName, attrValue)
    }
}
```

## Batch Operations

```go
// Save multiple entries
entries := []*mod_db.Entry{entry1, entry2, entry3}
errors := repo.SaveMultiEntry(entries...)

// Check for errors
for i, err := range errors {
    if err != nil {
        fmt.Printf("Failed to save entry %d: %v\n", i, err)
    }
}
```

## Index Management

```go
// Create index
err = repo.CreateEntryIndex()

// Drop index (careful!)
err = repo.DropEntryIndex()

// Wait for index to be ready
// (automatically handled by search methods via repo.waitIndex)
```

## Best Practices

### 1. objectClass Design

- Use `ObjectClassData` for schema-specific attributes that vary by objectClass
- Promote frequently-searched objectClass attributes to top-level Entry fields
- Keep objectClass names in the `ObjectClass` slice for TAG search performance

### 2. Indexing Strategy

- Index only searchable fields to minimize index size
- Use TAG for exact match fields (DN, UUID, objectClass names)
- Use NUMERIC for range queries (timestamps, UID/GID numbers)
- Use TEXT for full-text search (descriptions, names)

### 3. Key Management

- Use deterministic key generation (UUID v5 from DN)
- Include DN in key generation for idempotency
- Consider key prefixes for namespace separation

### 4. Search Performance

- Combine multiple criteria in a single query when possible
- Use `SearchByComplex` or fluent builder for multi-attribute searches
- Leverage RediSearch's query syntax for complex expressions
- For objectClass sub-attributes, consider application-level filtering

### 5. Data Consistency

- Always call `SyncObjectClass()` before saving
- Use `Ver` field for optimistic locking
- Set appropriate `ExAt` for entry expiration
- Update `ModifyTimestamp` and `ModifiersName` on changes

## Migration from Existing Entry Struct

To migrate from the current `Entry` struct to `EntryV2`:

1. Update struct definition to use `EntryV2`
2. Convert `ObjectClass []string` to `ObjectClassData ObjectClassList`
3. Add operational attribute fields
4. Update index creation to include new fields
5. Migrate existing data:

```go
func MigrateEntryToV2(old *Entry) *EntryV2 {
    v2 := &EntryV2{
        Key:  old.Key,
        Ver:  old.Ver,
        ExAt: old.Ext,

        // Map operational attributes
        EntryDN:   old.DN,
        EntryUUID: old.UUID,
        CreateTimestamp: old.CreateTimestamp,
        ModifyTimestamp: old.ModifyTimestamp,
        CreatorsName: old.CreatorsName,
        ModifiersName: old.ModifiersName,

        // Standard attributes
        CN:        old.CN,
        UID:       old.UID,
        UIDNumber: old.UIDNumber,
        GIDNumber: old.GIDNumber,
        // ... map other fields

        Type:   old.Type,
        Status: old.Status,
        BaseDN: old.BaseDN,
    }

    // Convert objectClass to ObjectClassData
    v2.ObjectClassData = mod_db.FromLDAPObjectClass(old.ObjectClass, nil)
    v2.SyncObjectClass()

    return v2
}
```

## Troubleshooting

### Search returns no results

1. Verify index exists: `FT.INFO idx:entry`
2. Check index is populated: `FT.SEARCH idx:entry * LIMIT 0 0`
3. Validate query syntax: Use Redis CLI to test queries
4. Ensure `SyncObjectClass()` was called before saving

### objectClass sub-attributes not searchable

- RediSearch doesn't index nested JSON arrays deeply by default
- Use application-level filtering or promote attributes to top-level fields

### Performance issues

- Reduce indexed fields if index is too large
- Use specific queries instead of wildcards
- Consider pagination for large result sets
- Use `LIMIT` clause in queries

## References

- [RediSearch Query Syntax](https://redis.io/docs/interact/search-and-query/query/)
- [rueidis Documentation](https://github.com/redis/rueidis)
- [RFC 4512 - LDAP Models](https://datatracker.ietf.org/doc/html/rfc4512)
- [RFC 2307bis - POSIX Schema](https://datatracker.ietf.org/doc/html/draft-howard-rfc2307bis)
