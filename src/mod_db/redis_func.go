package mod_db

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_ldap"
	"rmm23/src/mod_reflect"
	"rmm23/src/mod_slices"
)

func buildRedisearchSchema(inbound interface{}) (outbound *redisearch.Schema, err error) {
	var (
		schema = redisearch.NewSchema(redisearch.DefaultOptions)
		rt     reflect.Type
	)

	switch rt, err = mod_reflect.GetStructRT(inbound); {
	case err != nil:
		return
	}

	for i := 0; i < rt.NumField(); i++ {
		var (
			field                    = rt.Field(i)
			redisTag, redisearchTag  = field.Tag.Get(redisTagName), field.Tag.Get(rediSearchTagName)
			parts                    = mod_slices.SplitString(redisearchTag, tagSeparator, mod_slices.FlagNormalize)
			types, options, unknowns = make(map[string]bool), make(map[string]bool), make(map[string]bool)
		)

		switch {
		case len(redisTag) == 0 || len(redisearchTag) == 0 || len(parts) == 0:
			continue
		}

		for _, opt := range parts {
			switch opt {
			case rediSearchTagTypeIgnore, rediSearchTagTypeText, rediSearchTagTypeNumeric, rediSearchTagTypeTag, rediSearchTagTypeGeo:
				types[opt] = true
			case rediSearchTagOptionSortable:
				options[opt] = true
			default:
				unknowns[opt] = true
			}
		}

		switch {
		case len(types) == 0:
			return nil, mod_errors.ETagNoType
		case len(types) > 1:
			return nil, mod_errors.ETagMultiType
		case len(unknowns) > 0:
			return nil, mod_errors.ETagUnknown
		}

		switch {
		case types[rediSearchTagTypeIgnore]:
		case types[rediSearchTagTypeText]:
			schema.AddField(redisearch.NewTextFieldOptions(redisTag, redisearch.TextFieldOptions{
				Sortable: options[rediSearchTagOptionSortable],
			}))
		case types[rediSearchTagTypeNumeric]:
			schema.AddField(redisearch.NewNumericFieldOptions(redisTag, redisearch.NumericFieldOptions{
				Sortable: options[rediSearchTagOptionSortable],
			}))
		case types[rediSearchTagTypeTag]:
			schema.AddField(redisearch.NewTagFieldOptions(redisTag, redisearch.TagFieldOptions{
				Sortable:  options[rediSearchTagOptionSortable],
				Separator: sliceSeparator,
			}))
		case types[rediSearchTagTypeGeo]:
			schema.AddField(redisearch.NewGeoFieldOptions(redisTag, redisearch.GeoFieldOptions{}))
		default:
			return nil, mod_errors.EUnwilling
		}
	}

	return schema, nil
}

func getLDAPDocs(inbound *mod_ldap.Conf, schema *redisearch.Schema) (outbound []*redisearch.Document, err error) {
	var (
		ldap2doc = func(fnBaseDN string, fnSearchResultType string, fnSearchResult *ldap.SearchResult) (fnErr error) {
			for _, fnB := range fnSearchResult.Entries {
				var (
					fnDoc   *redisearch.Document
					fnEntry = new(entry)
				)

				switch fnErr = mod_ldap.UnmarshalEntry(fnB, fnEntry); {
				case fnErr != nil:
					return
				}

				switch fnErr = fnEntry.Type.Parse(fnSearchResultType); {
				case fnErr != nil:
					return
				}

				fnEntry.BaseDN = attrDN(fnBaseDN)
				fnEntry.Status = entryStatusLoaded
				fnEntry.UUID = attrUUID(uuid.NewSHA1(uuid.Nil, []byte(fnEntry.DN.String())))

				switch fnDoc, fnErr = marshalRedisearchDoc(
					schema,
					fnEntry.UUID.Entry(),
					1.0,
					fnEntry,
					false,
				); {
				case fnErr != nil:
					return
				}

				outbound = append(outbound, fnDoc)
			}

			return
		}
	)

	switch err = inbound.SearchFn(ldap2doc); {
	case err != nil:
		return
	}

	return
}

func createQuery(inbound ...any) (outbound string) {
	return escapeQuery(mod_slices.Join(inbound, ":", mod_slices.FlagNone))
}

// escapeQuery escapes special characters in a string for Redisearch queries.
// It adds a backslash before any character that has special meaning in Redisearch query syntax.
func escapeQuery(inbound any) (outbound string) {
	var (
		interim strings.Builder
	)

	for _, b := range fmt.Sprint(inbound) {
		switch b {
		case ',', '.', '<', '>', '{', '}', '[', ']', '"', '\'', ':', ';', '!', '@', '#', '$', '%', '^', '&', '*', '-', '+', '=', '~', '|':
			interim.WriteRune('\\') // Add escape character
		}

		interim.WriteRune(b) // Add the character itself
	}

	return interim.String()
}
