package mod_db

import (
	"context"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"

	"rmm23/src/mod_errors"
	"rmm23/src/mod_ldap"
	"rmm23/src/mod_slices"
	"rmm23/src/mod_strings"
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
		rsClient = redisearch.NewClientFromPool(&redis.Pool{
			DialContext: func(ctx context.Context) (redis.Conn, error) {
				return redis.DialContext(ctx, rcNetwork, rcAddress, redis.DialDatabase(0 /* only zero for indexing */))
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
		doc.Set("Type", entryTypeDomain)

		doc.Set("UUID", d.Domain.UUID)
		doc.Set("DN", d.Domain.DN)
		doc.Set("ObjectClass", mod_slices.Join(d.Domain.ObjectClass, mod_strings.SliceDelimiter, mod_slices.FlagNormalize))
		doc.Set("CreatorsName", d.Domain.CreatorsName)
		doc.Set("CreateTimestamp", d.Domain.CreateTimestamp)
		doc.Set("ModifiersName", d.Domain.ModifiersName)
		doc.Set("ModifyTimestamp", d.Domain.ModifyTimestamp)

		doc.Set("DC", d.Domain.DC)
		doc.Set("O", d.Domain.O)

		doc.Set("Legacy", d.Domain.LabeledURI)

		switch err = rsClient.Index([]redisearch.Document{doc}...); {
		case err != nil && mod_errors.Contains(err, EDocExist):
			err = nil
		case err != nil:
			return
		}

		for _, f := range d.Groups {
			doc = redisearch.NewDocument("ldap:entry:"+f.UUID.String(), 1.0)
			doc.Set("Type", entryTypeGroup)

			doc.Set("UUID", f.UUID)
			doc.Set("DN", f.DN)
			doc.Set("ObjectClass", mod_slices.Join(f.ObjectClass, mod_strings.SliceDelimiter, mod_slices.FlagNormalize))
			doc.Set("CreatorsName", f.CreatorsName)
			doc.Set("CreateTimestamp", f.CreateTimestamp)
			doc.Set("ModifiersName", f.ModifiersName)
			doc.Set("ModifyTimestamp", f.ModifyTimestamp)

			doc.Set("CN", f.CN)
			doc.Set("GIDNumber", f.GIDNumber)
			doc.Set("Member", mod_slices.Join(f.Member, mod_strings.SliceDelimiter, mod_slices.FlagNormalize))
			doc.Set("Owner", mod_slices.Join(f.Owner, mod_strings.SliceDelimiter, mod_slices.FlagNormalize))

			doc.Set("Legacy", f.LabeledURI)

			switch err = rsClient.Index([]redisearch.Document{doc}...); {
			case err != nil && mod_errors.Contains(err, EDocExist):
				err = nil
			case err != nil:
				return
			}
		}

		for _, f := range d.Users {
			doc = redisearch.NewDocument("ldap:entry:"+f.UUID.String(), 1.0)
			doc.Set("Type", entryTypeUser)

			doc.Set("UUID", f.UUID)
			doc.Set("DN", f.DN)
			doc.Set("ObjectClass", mod_slices.Join(f.ObjectClass, mod_strings.SliceDelimiter, mod_slices.FlagNormalize))
			doc.Set("CreatorsName", f.CreatorsName)
			doc.Set("CreateTimestamp", f.CreateTimestamp)
			doc.Set("ModifiersName", f.ModifiersName)
			doc.Set("ModifyTimestamp", f.ModifyTimestamp)

			doc.Set("CN", f.CN)
			doc.Set("Description", f.Description)
			doc.Set("DestinationIndicator", f.DestinationIndicator)
			doc.Set("DisplayName", f.DisplayName)
			doc.Set("GIDNumber", f.GIDNumber)
			doc.Set("HomeDirectory", f.HomeDirectory)
			doc.Set("IPHostNumber", f.IPHostNumber)
			doc.Set("Mail", f.Mail)
			// doc.Set("MemberOf", f.MemberOf)
			doc.Set("O", f.O)
			doc.Set("OU", f.OU)
			doc.Set("SN", f.SN)
			doc.Set("SSHPublicKey", f.SSHPublicKey)
			doc.Set("TelephoneNumber", f.TelephoneNumber)
			doc.Set("TelexNumber", f.TelexNumber)
			doc.Set("UID", f.UID)
			doc.Set("UIDNumber", f.UIDNumber)
			doc.Set("UserPKCS12", f.UserPKCS12)
			doc.Set("UserPassword", f.UserPassword)

			doc.Set("Legacy", f.LabeledURI)

			switch err = rsClient.Index([]redisearch.Document{doc}...); {
			case err != nil && mod_errors.Contains(err, EDocExist):
				err = nil
			case err != nil:
				return
			}
		}

		for _, f := range d.Hosts {
			doc = redisearch.NewDocument("ldap:entry:"+f.UUID.String(), 1.0)
			doc.Set("Type", entryTypeHost)

			doc.Set("UUID", f.UUID)
			doc.Set("DN", f.DN)
			doc.Set("ObjectClass", mod_slices.Join(f.ObjectClass, mod_strings.SliceDelimiter, mod_slices.FlagNormalize))
			doc.Set("CreatorsName", f.CreatorsName)
			doc.Set("CreateTimestamp", f.CreateTimestamp)
			doc.Set("ModifiersName", f.ModifiersName)
			doc.Set("ModifyTimestamp", f.ModifyTimestamp)

			doc.Set("CN", f.CN)
			doc.Set("GIDNumber", f.GIDNumber)
			doc.Set("HomeDirectory", f.HomeDirectory)
			// doc.Set("MemberOf", f.MemberOf)
			doc.Set("SN", f.SN)
			doc.Set("UID", f.UID)
			doc.Set("UIDNumber", f.UIDNumber)
			doc.Set("UserPKCS12", f.UserPKCS12)

			doc.Set("Legacy", f.LabeledURI)

			switch err = rsClient.Index([]redisearch.Document{doc}...); {
			case err != nil && mod_errors.Contains(err, EDocExist):
				err = nil
			case err != nil:
				return
			}
		}

	}

	return nil
}
