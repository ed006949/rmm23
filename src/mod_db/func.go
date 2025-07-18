package mod_db

import (
	"context"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_ldap"
	"rmm23/src/mod_slices"
)

func CopyLDAP2DB(ctx context.Context, inbound *mod_ldap.Conf) (err error) {
	switch err = inbound.Fetch(); {
	case err != nil:
		return
	}

	var (
		rcAddress = "10.133.0.223:6379"
		// rcDB      = 0
		rcNetwork = "tcp"
		rcName    = "entryIdx"
	)

	var (
		// RediSearch requires DB 0 for index creation
		rsClient = redisearch.NewClientFromPool(&redis.Pool{
			DialContext: func(ctx context.Context) (redis.Conn, error) {
				return redis.DialContext(ctx, rcNetwork, rcAddress, redis.DialDatabase(0))
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) (tErr error) {
				_, tErr = c.Do("PING")
				return
			},
			MaxIdle:         4,
			MaxActive:       4,
			IdleTimeout:     240 * time.Second,
			Wait:            true,
			MaxConnLifetime: 0,
		}, rcName)
		entry  = Entry{}
		schema = entry.redisearchSchema()
		_      = rsClient.Drop() // test&dev, delete old entries
		// _        = rsClient.DropIndex(false) // prod, don't delete old entries
	)

	switch err = rsClient.CreateIndex(entry.redisearchSchema()); {
	case err != nil:
		return
	}

	for _, d := range inbound.Domain {
		var (
			doc redisearch.Document
		)

		switch doc, err = newRedisearchDocument(
			schema,
			mod_slices.Join([]string{"ldap", "entry", d.Domain.UUID.String()}, ":", mod_slices.FlagNone),
			1.0,
			d.Domain,
			false,
		); {
		case err != nil:
			return err
		default:
			doc.Set("type", entryTypeDomain)
			switch err = rsClient.Index([]redisearch.Document{doc}...); {
			case mod_errors.Contains(err, EDocExist):
				err = nil
			case err != nil:
				return err
			}
		}

		for _, f := range d.Groups {
			switch doc, err = newRedisearchDocument(
				schema,
				mod_slices.Join([]string{"ldap", "entry", f.UUID.String()}, ":", mod_slices.FlagNone),
				1.0,
				f,
				false,
			); {
			case err != nil:
				return err
			default:
				doc.Set("type", entryTypeGroup)
				switch err = rsClient.Index([]redisearch.Document{doc}...); {
				case mod_errors.Contains(err, EDocExist):
					err = nil
				case err != nil:
					return err
				}
			}
		}
		for _, f := range d.Users {
			switch doc, err = newRedisearchDocument(
				schema,
				mod_slices.Join([]string{"ldap", "entry", f.UUID.String()}, ":", mod_slices.FlagNone),
				1.0,
				f,
				false,
			); {
			case err != nil:
				return err
			default:
				doc.Set("type", entryTypeUser)
				switch err = rsClient.Index([]redisearch.Document{doc}...); {
				case mod_errors.Contains(err, EDocExist):
					err = nil
				case err != nil:
					return err
				}
			}
		}
		for _, f := range d.Hosts {
			switch doc, err = newRedisearchDocument(
				schema,
				mod_slices.Join([]string{"ldap", "entry", f.UUID.String()}, ":", mod_slices.FlagNone),
				1.0,
				f,
				false,
			); {
			case err != nil:
				return err
			default:
				doc.Set("type", entryTypeHost)
				switch err = rsClient.Index([]redisearch.Document{doc}...); {
				case mod_errors.Contains(err, EDocExist):
					err = nil
				case err != nil:
					return err
				}
			}
		}
	}

	return nil
}
