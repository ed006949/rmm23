package mod_db

import (
	"context"
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"

	"rmm23/src/l"
	"rmm23/src/mod_ldap"
)

func CopyLDAP2DB(ctx context.Context, inbound *mod_ldap.Conf, outbound *Conf) (err error) {
	l.CLEAR = false

	l.Z{l.M: "indexing", l.E: err}.Warning()

	switch err = outbound.Dial(ctx); {
	case err != nil:
		return
	}

	defer func() {
		_ = outbound.Close()
	}()

	switch err = getLDAPDocs(ctx, inbound, outbound.repo, outbound); {
	case err != nil:
		return
	}

	var (
		count   int64
		entries []*Entry
	)

	l.Z{l.M: "indexed", l.E: err}.Warning()

	// count, entries, err = outbound.repo.repo.Search(ctx, func(search om.FtSearchIndex) rueidis.Completed {
	// 	return search.Query("*").Limit().OffsetNum(0, connMaxPaging).Build()
	// })
	// l.Z{l.M: count, l.E: err, "entries": len(entries)}.Warning()
	//
	// count, entries, err = outbound.repo.repo.Search(ctx, func(search om.FtSearchIndex) rueidis.Completed {
	// 	return search.Query("@objectClass:{posixGroup}").Limit().OffsetNum(0, connMaxPaging).Build()
	// })
	// l.Z{l.M: count, l.E: err, "entries": len(entries)}.Warning()

	count, entries, err = outbound.repo.repo.Search(ctx, func(search om.FtSearchIndex) rueidis.Completed {
		return search.Query("@baseDN:{dc\\=fabric\\,dc\\=domain\\,dc\\=tld} @objectClass:{posixGroup}").Limit().OffsetNum(0, connMaxPaging).Build()
	})
	l.Z{l.M: count, l.E: err, "entries": len(entries)}.Warning()

	count, entries, err = outbound.repo.SearchFV(ctx, _member, `uid=aaa\-auth,ou=People,dc=fabric,dc=domain,dc=tld`)
	l.Z{l.M: count, l.E: err, "entries": len(entries)}.Warning()

	count, entries, err = outbound.repo.SearchMFV(
		ctx,
		[]_FV{
			{
				_F: _baseDN,
				_V: "dc=fabric,dc=domain,dc=tld",
			},
			{
				_F: _objectClass,
				_V: "posixGroup",
			},
		},
	)
	l.Z{l.M: count, l.E: err, "entries": len(entries)}.Warning()

	// cmd := outbound.client.B().FtSearch().Index(outbound.repo.repo.IndexName()).Query(`*`).Build()
	// // count, entries, err = outbound.repo.repo.Search(ctx, func(search om.FtSearchIndex) rueidis.Completed {
	// // 	return search.Query(`@dn:dc`).Build()
	// // })
	// count, entries, err = outbound.client.Do(ctx, cmd).AsFtSearch()

	return
}

// getLDAPDocs fetches entries from LDAP and saves them to Redis using the provided RedisRepository.
//
// copy entry one-by-one to save memory.
func getLDAPDocs(ctx context.Context, inbound *mod_ldap.Conf, repo *RedisRepository, c *Conf) (err error) {
	switch l.CLEAR {
	case false:
		return
	}

	var (
		ldap2doc = func(fnBaseDN string, fnSearchResultType string, fnSearchResult *ldap.SearchResult) (fnErr error) {
			for _, fnB := range fnSearchResult.Entries {
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

				// switch l.CLEAR {
				// case true:
				// 	_ = repo.DeleteEntry(ctx, fnEntry.Key)
				// }

				_ = repo.DeleteEntry(ctx, fnEntry.Key)

				// Save the Entry using the RedisRepository
				switch fnErr = repo.SaveEntry(ctx, fnEntry); {
				case fnErr != nil:
					return
				}

				fmt.Printf("\nDN: %s\n", fnEntry.DN.String())

				_ = c.monitorIndexingFailures(ctx)
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
