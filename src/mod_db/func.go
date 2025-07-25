package mod_db

import (
	"context"
	"net/url"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_ldap"
	"rmm23/src/mod_net"
	"rmm23/src/mod_slices"
)

func redisNetwork(inbound *url.URL) (outbound string, err error) {
	switch outbound = inbound.Scheme; outbound {
	case "redis", "redis-sentinel":
		return "tcp", nil
	case "file":
		return "unix", nil
	default:
		return outbound, mod_errors.EUnknownScheme
	}
}

// func (r Elements) unmarshal(inbound *ldap.SearchResult) (err error) {
// 	for _, entry := range inbound.Entries {
// 		var (
// 			interim Element
// 		)
//
// 		switch newErr := UnmarshalEntry(entry, &interim); {
// 		case newErr != nil:
// 			err = errors.Join(err, newErr)
// 			l.Z{l.E: err, l.M: "LDAP Unmarshal", "DN": entry.DN}.Warning()
// 		}
//
// 		r[interim.DN] = &interim
// 	}
//
// 	return
// }

func CopyLDAP2DB(ctx context.Context, inbound *mod_ldap.LDAPConfig, outbound *mod_net.URL) (err error) {
	switch err = inbound.Search(); {
	case err != nil:
		return
	}

	var (
		// RediSearch requires DB 0 for index creation
		// rcDB      = 0

		rcNetwork string
		rcName    = "entryIdx"
	)

	// var (
	// 	searchRequest = ldap.NewSearchRequest(
	// 		mod_slices.JoinStrings([]string{d.DN.String(), b.DN.String()}, ",", mod_slices.FlagFilterEmpty), // Base DN
	// 		ldap.ScopeBaseObject, // Scope - search entire tree
	// 		ldap.DerefAlways,     // Deref
	// 		0,                    // Size limit (0 = no limit)
	// 		0,                    // Time limit (0 = no limit)
	// 		false,                // Types only
	// 		d.Filter,             // Filter - all objects
	// 		[]string{"*", "+"},   // Attributes - all user and operational attributes
	// 		nil,                  // Controls
	// 	)
	// )

	switch rcNetwork, err = redisNetwork(outbound.URL); {
	case err != nil:
		return
	}

	var (
		rsClient = redisearch.NewClientFromPool(&redis.Pool{
			DialContext: func(ctx context.Context) (redis.Conn, error) {
				return redis.DialContext(ctx, rcNetwork, outbound.Host, redis.DialDatabase(0))
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) (tErr error) {
				_, tErr = c.Do(_PING)

				return
			},
			MaxIdle:         connMaxIdle,
			MaxActive:       connMaxActive,
			IdleTimeout:     connIdleTimeout,
			Wait:            connWait,
			MaxConnLifetime: connMaxConnLifetime,
		}, rcName)
		entry  = Entry{}
		schema = entry.redisearchSchema()
		_      = rsClient.Drop() // test&dev, delete old entries
		// _        = rsClient.DropIndex(false) // prod, don't delete old entries
	)

	switch err = rsClient.CreateIndex(schema); {
	case err != nil:
		return
	}

	for _, d := range inbound.Domains {
		var (
			doc redisearch.Document
			// t = mod_ldap.UnmarshalResult(d)
		)

		switch doc, err = newRedisearchDocument(
			schema,
			mod_slices.JoinStrings([]string{"ldap", "entry", d.Domain.UUID.String()}, ":", mod_slices.FlagNone),
			1.0,
			d.Domain,
			false,
		); {
		case err != nil:
			return
		default:
			doc.Set("type", entryTypeDomain)

			switch err = rsClient.Index([]redisearch.Document{doc}...); {
			case mod_errors.Contains(err, mod_errors.EDocExist):
			case err != nil:
				return
			}
		}

		for _, f := range d.Groups {
			switch doc, err = newRedisearchDocument(
				schema,
				mod_slices.JoinStrings([]string{"ldap", "entry", f.UUID.String()}, ":", mod_slices.FlagNone),
				1.0,
				f,
				false,
			); {
			case err != nil:
				return
			default:
				doc.Set("type", entryTypeGroup)

				switch err = rsClient.Index([]redisearch.Document{doc}...); {
				case mod_errors.Contains(err, mod_errors.EDocExist):
				case err != nil:
					return
				}
			}
		}

		for _, f := range d.Users {
			switch doc, err = newRedisearchDocument(
				schema,
				mod_slices.JoinStrings([]string{"ldap", "entry", f.UUID.String()}, ":", mod_slices.FlagNone),
				1.0,
				f,
				false,
			); {
			case err != nil:
				return
			default:
				doc.Set("type", entryTypeUser)

				switch err = rsClient.Index([]redisearch.Document{doc}...); {
				case mod_errors.Contains(err, mod_errors.EDocExist):
				case err != nil:
					return
				}
			}
		}

		for _, f := range d.Hosts {
			switch doc, err = newRedisearchDocument(
				schema,
				mod_slices.JoinStrings([]string{"ldap", "entry", f.UUID.String()}, ":", mod_slices.FlagNone),
				1.0,
				f,
				false,
			); {
			case err != nil:
				return
			default:
				doc.Set("type", entryTypeHost)

				switch err = rsClient.Index([]redisearch.Document{doc}...); {
				case mod_errors.Contains(err, mod_errors.EDocExist):
				case err != nil:
					return
				}
			}
		}
	}

	return
}
