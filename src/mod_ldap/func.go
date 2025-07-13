package mod_ldap

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-ldap/ldap/v3"
)

// SearchLDAP connects to an LDAP server, performs a search, and returns the entries.
func SearchLDAP(ldapURL, bindDN, bindPassword, baseDN, filter string, attributes []string) ([]*ldap.Entry, error) {
	l, err := ldap.DialURL(ldapURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to LDAP server: %w", err)
	}
	defer l.Close()

	// Bind to the LDAP server
	err = l.Bind(bindDN, bindPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to bind to LDAP server: %w", err)
	}

	// Create a search request
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0,
		false,
		filter,
		attributes,
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to perform LDAP search: %w", err)
	}

	return sr.Entries, nil
}

// ReadLDAPAndStoreInRedis reads LDAP entries, unmarshals them, and stores them in Redis with RediSearch schemas.
func ReadLDAPAndStoreInRedis(redisAddr, redisIndexName, ldapURL, bindDN, bindPassword, baseDN, filter string, attributes []string) error {
	// Initialize RediSearch client
	client := redisearch.NewClient(redisAddr, redisIndexName)

	// Create LDAP schema in RediSearch
	err := createLDAPSchema(client)
	if err != nil {
		return fmt.Errorf("failed to create RediSearch schema: %w", err)
	}

	// Search LDAP entries
	entries, err := SearchLDAP(ldapURL, bindDN, bindPassword, baseDN, filter, attributes)
	if err != nil {
		return fmt.Errorf("failed to search LDAP: %w", err)
	}

	var docs []redisearch.Document
	for _, entry := range entries {
		// Determine the type of LDAP entry based on objectClass and UnmarshalEntry accordingly
		objectClasses := entry.GetAttributeValues("objectClass")
		switch {
		case contains(objectClasses, "person"):
			var user ElementUser
			err = UnmarshalEntry(entry, &user)
			if err != nil {
				log.Printf("Warning: failed to UnmarshalEntry user entry %s: %v", entry.DN, err)
				continue
			}
			doc, err := StructToDocument(entry.DN, 1.0, user)
			if err != nil {
				log.Printf("Warning: failed to convert user to document %s: %v", entry.DN, err)
				continue
			}
			docs = append(docs, *doc)
		case contains(objectClasses, "groupOfNames"), contains(objectClasses, "groupOfUniqueNames"):
			var group ElementGroup
			err = UnmarshalEntry(entry, &group)
			if err != nil {
				log.Printf("Warning: failed to UnmarshalEntry group entry %s: %v", entry.DN, err)
				continue
			}
			doc, err := StructToDocument(entry.DN, 1.0, group)
			if err != nil {
				log.Printf("Warning: failed to convert group to document %s: %v", entry.DN, err)
				continue
			}
			docs = append(docs, *doc)
		case contains(objectClasses, "device"):
			var host ElementHost
			err = UnmarshalEntry(entry, &host)
			if err != nil {
				log.Printf("Warning: failed to UnmarshalEntry host entry %s: %v", entry.DN, err)
				continue
			}
			doc, err := StructToDocument(entry.DN, 1.0, host)
			if err != nil {
				log.Printf("Warning: failed to convert host to document %s: %v", entry.DN, err)
				continue
			}
			docs = append(docs, *doc)
		default:
			log.Printf("Info: Skipping unsupported LDAP entry type for DN: %s, ObjectClasses: %v", entry.DN, objectClasses)
		}
	}

	// Index documents in RediSearch
	if len(docs) > 0 {
		if err := client.Index(docs...); err != nil {
			return fmt.Errorf("failed to index documents in RediSearch: %w", err)
		}
	}

	log.Printf("Successfully read %d LDAP entries and stored them in Redis with RediSearch.", len(docs))
	return nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func createLDAPSchema(client *redisearch.Client) (err error) {
	// Drop existing index if any
	_ = client.Drop()

	// Define schema for entire LDAP tree: Users, Groups, devices, etc.
	// Common LDAP attributes and indexes for fast search

	var (
		schema = redisearch.NewSchema(redisearch.DefaultOptions).
			// User attributes
			AddField(redisearch.NewTextFieldOptions("uid", redisearch.TextFieldOptions{Sortable: true})).
			AddField(redisearch.NewNumericFieldOptions("uidNumber", redisearch.NumericFieldOptions{Sortable: true})).
			AddField(redisearch.NewNumericField("gidNumber")).
			AddField(redisearch.NewTextField("cn")).
			AddField(redisearch.NewTextField("mail")).
			AddField(redisearch.NewTextField("sn")). // surname
			AddField(redisearch.NewTextField("givenName")).
			AddField(redisearch.NewTextField("displayName")).
			AddField(redisearch.NewTextField("description")).
			AddField(redisearch.NewTextField("homeDirectory")).
			AddField(redisearch.NewTextField("loginShell")).
			AddField(redisearch.NewTextField("memberOf")). // Groups user belongs to

			// Group attributes
			AddField(redisearch.NewTextFieldOptions("cn", redisearch.TextFieldOptions{Sortable: true})).
			AddField(redisearch.NewTextField("gidNumber")).
			AddField(redisearch.NewTextField("member")). // Users in group

			// Device attributes (example)
			AddField(redisearch.NewTextField("deviceID")).
			AddField(redisearch.NewTextField("deviceType")).
			AddField(redisearch.NewTextField("deviceOwner")).
			AddField(redisearch.NewTextField("deviceDescription"))
	)

	// Create the index
	return client.CreateIndex(schema)
}

func StructToDocument(id string, score float32, s interface{}) (*redisearch.Document, error) {
	doc := redisearch.NewDocument(id, score)

	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct")
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip unexported fields
		if !value.CanInterface() {
			continue
		}

		// Get redisearch tag or use field name
		fieldName := field.Tag.Get("redisearch")
		if fieldName == "" {
			fieldName = strings.ToLower(field.Name)
		}

		// Set field in document
		doc.Set(fieldName, value.Interface())
	}

	return &doc, nil
}

func main() {
	// Example usage:
	// client := redisearch.NewClient("localhost:6379", "ldapIndex")
	// err := createLDAPSchema(client)
	// if err != nil {
	//     log.Fatalf("Failed to create LDAP schema: %v", err)
	// }
}
