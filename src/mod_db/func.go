package mod_db

import (
	"context"
	"strings"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"

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
		rcDB      = 0 // only zero for indexing
		rcNetwork = "tcp"
		rcName    = "entryIdx"
	)

	var (
		rsClient = redisearch.NewClientFromPool(&redis.Pool{
			DialContext: func(ctx context.Context) (redis.Conn, error) {
				return redis.DialContext(ctx, rcNetwork, rcAddress, redis.DialDatabase(rcDB))
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
		entry = Entry{}
		_     = rsClient.Drop() // test&dev, delete old entries
		// _        = rsClient.DropIndex(false) // prod, don't delete old entries
	)

	switch err = rsClient.CreateIndex(entry.RedisearchSchema()); {
	case err != nil:
		return
	}

	for _, d := range inbound.Domain {
		var (
			doc = redisearch.NewDocument("ldap:entry:"+d.Domain.UUID.String(), 1.0)
		)
		doc.Set("Type", EntryTypeDomain)
		doc.Set("UUID", d.Domain.UUID)
		doc.Set("DN", d.Domain.DN)
		doc.Set("ObjectClass", mod_slices.Join(d.Domain.ObjectClass, sliceDelimiter, mod_slices.FlagNormalize))
		doc.Set("CreatorsName", d.Domain.CreatorsName)
		doc.Set("CreateTimestamp", d.Domain.CreateTimestamp)
		doc.Set("ModifiersName", d.Domain.ModifiersName)
		doc.Set("ModifyTimestamp", d.Domain.ModifyTimestamp)
		doc.Set("DC", d.Domain.DC)
		doc.Set("O", d.Domain.O)
		doc.Set("Legacy", d.Domain.LabeledURI)

		switch err = rsClient.Index([]redisearch.Document{doc}...); {
		case err != nil && strings.Contains(err.Error(), EDocExist.Error()):
			err = nil
		case err != nil:
			return
		}

		for _, g := range d.Groups {
			doc = redisearch.NewDocument("ldap:entry:"+g.UUID.String(), 1.0)
			doc.Set("Type", EntryTypeGroup)
			doc.Set("UUID", g.UUID)
			doc.Set("DN", g.DN)
			doc.Set("ObjectClass", mod_slices.Join(g.ObjectClass, sliceDelimiter, mod_slices.FlagNormalize))
			doc.Set("CreatorsName", g.CreatorsName)
			doc.Set("CreateTimestamp", g.CreateTimestamp)
			doc.Set("ModifiersName", g.ModifiersName)
			doc.Set("ModifyTimestamp", g.ModifyTimestamp)
			doc.Set("CN", g.CN)
			doc.Set("Owner", mod_slices.Join(g.Owner, sliceDelimiter, mod_slices.FlagNormalize))
			doc.Set("Member", mod_slices.Join(g.Member, sliceDelimiter, mod_slices.FlagNormalize))
			doc.Set("GIDNumber", g.GIDNumber)
			doc.Set("Legacy", g.LabeledURI)

			switch err = rsClient.Index([]redisearch.Document{doc}...); {
			case err != nil && strings.Contains(err.Error(), EDocExist.Error()):
				err = nil
			case err != nil:
				return
			}
		}

		for _, u := range d.Users {
			doc = redisearch.NewDocument("ldap:entry:"+u.UUID.String(), 1.0)
			doc.Set("Type", EntryTypeUser)
			doc.Set("UUID", u.UUID)
			doc.Set("DN", u.DN)
			doc.Set("ObjectClass", mod_slices.Join(u.ObjectClass, sliceDelimiter, mod_slices.FlagNormalize))
			doc.Set("CreatorsName", u.CreatorsName)
			doc.Set("CreateTimestamp", u.CreateTimestamp)
			doc.Set("ModifiersName", u.ModifiersName)
			doc.Set("ModifyTimestamp", u.ModifyTimestamp)
			doc.Set("CN", u.CN)
			doc.Set("GIDNumber", u.GIDNumber)
			doc.Set("Legacy", u.LabeledURI)

			switch err = rsClient.Index([]redisearch.Document{doc}...); {
			case err != nil && strings.Contains(err.Error(), EDocExist.Error()):
				err = nil
			case err != nil:
				return
			}
		}

		for _, h := range d.Hosts {
			doc = redisearch.NewDocument("ldap:entry:"+h.UUID.String(), 1.0)
			doc.Set("Type", EntryTypeHost)
			doc.Set("UUID", h.UUID)
			doc.Set("DN", h.DN)
			doc.Set("ObjectClass", mod_slices.Join(h.ObjectClass, sliceDelimiter, mod_slices.FlagNormalize))
			doc.Set("CreatorsName", h.CreatorsName)
			doc.Set("CreateTimestamp", h.CreateTimestamp)
			doc.Set("ModifiersName", h.ModifiersName)
			doc.Set("ModifyTimestamp", h.ModifyTimestamp)
			doc.Set("CN", h.CN)
			doc.Set("GIDNumber", h.GIDNumber)
			doc.Set("Legacy", h.LabeledURI)

			switch err = rsClient.Index([]redisearch.Document{doc}...); {
			case err != nil && strings.Contains(err.Error(), EDocExist.Error()):
				err = nil
			case err != nil:
				return
			}
		}

	}

	return nil
}
