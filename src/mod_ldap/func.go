package mod_ldap

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/RediSearch/redisearch-go/redisearch"
)

func createLDAPSchema(client *redisearch.Client) (err error) {
	// Drop existing index if any
	_ = client.Drop()

	// Define schema for entire LDAP tree: users, groups, devices, etc.
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
			AddField(redisearch.NewTextField("memberOf")). // groups user belongs to

			// Group attributes
			AddField(redisearch.NewTextFieldOptions("cn", redisearch.TextFieldOptions{Sortable: true})).
			AddField(redisearch.NewTextField("gidNumber")).
			AddField(redisearch.NewTextField("member")). // users in group

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
