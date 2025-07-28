package mod_db

import (
	"context"

	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"

	"rmm23/src/l"
	"rmm23/src/mod_ldap"
)

func CopyLDAP2DB(ctx context.Context, inbound *mod_ldap.Conf, outbound *Conf) (err error) {
	l.CLEAR = true

	switch err = outbound.Dial(); {
	case err != nil:
		return
	}

	defer func() {
		_ = outbound.Close()
	}()

	switch err = getLDAPDocs(ctx, inbound, outbound.repo); {
	case err != nil:
		return
	}

	return
}

// getLDAPDocs fetches entries from LDAP and saves them to Redis using the provided RedisRepository.
//
// copy entry one-by-one to save memory.
func getLDAPDocs(ctx context.Context, inbound *mod_ldap.Conf, repo *RedisRepository) (err error) {
	switch l.CLEAR {
	case false:
		return
	}

	var (
		ldap2doc = func(fnBaseDN string, fnSearchResultType string, fnSearchResult *ldap.SearchResult) (fnErr error) {
			for counter, fnB := range fnSearchResult.Entries {
				switch {
				case l.CLEAR && counter > 3:
					return nil
				}

				var (
					fnEntry = new(Entry)
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

				fnEntry.Key = fnEntry.UUID.String()

				switch l.CLEAR {
				case true:
					_ = repo.repo.Remove(ctx, fnEntry.Key)
				}

				// Save the Entry using the RedisRepository
				switch fnErr = repo.SaveEntry(ctx, fnEntry); {
				case fnErr != nil:
					return
				}
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
