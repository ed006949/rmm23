package mod_db

import (
	"reflect"
	"strings"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"

	"rmm23/src/l"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_ldap"
	"rmm23/src/mod_reflect"
	"rmm23/src/mod_slices"
)

func buildRedisearchSchema(inbound interface{}) (schema *redisearch.Schema, schemaMap schemaMapType, err error) {
	var (
		rt reflect.Type
	)

	switch rt, err = mod_reflect.GetStructRT(inbound); {
	case err != nil:
		return
	}

	schema = redisearch.NewSchema(redisearch.DefaultOptions)

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
			return nil, nil, mod_errors.ETagNoType
		case len(types) > 1:
			return nil, nil, mod_errors.ETagMultiType
		case len(unknowns) > 0:
			return nil, nil, mod_errors.ETagUnknown
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
			return nil, nil, mod_errors.EUnwilling
		}
	}

	schemaMap = make(schemaMapType)
	for _, b := range schema.Fields {
		schemaMap[entryFieldName(b.Name)] = b.Type
	}

	return
}

func getLDAPDocs(inbound *mod_ldap.Conf, schema *redisearch.Schema) (outbound []*redisearch.Document, err error) {
	switch l.CLEAR {
	case false:
		return
	}

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

func escapeQueryValue(inbound string) string {
	replacer := strings.NewReplacer(
		"=", "\\=",
		",", "\\,",
		"(", "\\(",
		")", "\\)",
		"{", "\\{",
		"}", "\\}",
		"[", "\\[",
		"]", "\\]",
		"\"", "\\\"",
		"'", "\\'",
		"~", "\\~",
	)

	return replacer.Replace(inbound)
}
