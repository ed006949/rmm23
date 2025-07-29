package mod_db

import (
	"context"
	"fmt"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"

	"rmm23/src/mod_errors"
)

func (r *Conf) Dial(ctx context.Context) (err error) {
	switch r.client, err = rueidis.NewClient(rueidis.ClientOption{
		AlwaysRESP2:  true,
		DisableCache: true,
		InitAddress:  []string{r.URL.Host},
	}); {
	case err != nil:
		return
	}

	r.repo = NewRedisRepository(r.client)

	_ = r.repo.DropIndex(ctx)

	switch err = r.repo.CreateIndex(ctx); {
	case err != nil:
		return
	}

	_ = r.monitorIndexingFailures(ctx)

	return
}

func (r *Conf) Close() (err error) {
	switch {
	case r.client == nil:
		return mod_errors.ENoConn
	}

	r.client.Close()

	return
}

func (r *Conf) monitorIndexingFailures(ctx context.Context) (err error) {
	var (
		resp = r.client.Do(ctx, r.client.B().FtInfo().Index(r.repo.repo.IndexName()).Build())
	)

	switch err = resp.Error(); {
	case err != nil:
		return
	}

	var (
		info map[string]string
	)

	switch info, err = resp.AsStrMap(); {
	case err != nil:
		return
	}

	fmt.Printf("hash_indexing_failures: %s\n", info["hash_indexing_failures"])

	// for a, b := range info {
	// 	switch a {
	// 	case "hash_indexing_failures":
	// 		fmt.Printf("%s: %s\n", a, b)
	// 	}
	// }

	return
}

// SaveEntry saves an Entry to Redis.
func (r *RedisRepository) SaveEntry(ctx context.Context, e *Entry) error {
	return r.repo.Save(ctx, e)
}

// FindEntry finds an Entry by its ID.
func (r *RedisRepository) FindEntry(ctx context.Context, id string) (*Entry, error) {
	return r.repo.Fetch(ctx, id)
}

// DeleteEntry deletes an Entry by its ID.
func (r *RedisRepository) DeleteEntry(ctx context.Context, id string) error {
	return r.repo.Remove(ctx, id)
}

// func (r *RedisRepository) SearchEntries() {}

// DropIndex drops the RediSearch index for the Entry struct.
func (r *RedisRepository) DropIndex(ctx context.Context) (err error) { return r.repo.DropIndex(ctx) }

// CreateIndex creates the RediSearch index for the Entry struct.
func (r *RedisRepository) CreateIndex(ctx context.Context) (err error) {
	return r.repo.CreateIndex(ctx, func(schema om.FtCreateSchema) rueidis.Completed {
		return schema.
			FieldName("$.type").As("type").Numeric().
			FieldName("$.status").As("status").Numeric().
			FieldName("$.baseDN").As("baseDN").Tag().
			FieldName("$.uuid").As("uuid").Tag().
			FieldName("$.dn").As("dn").Tag().
			// FieldName("$.objectClass").As("objectClass").Tag().
			FieldName("$.creatorsName").As("creatorsName").Tag().
			FieldName("$.createTimestamp").As("createTimestamp").Numeric().
			FieldName("$.modifiersName").As("modifiersName").Tag().
			FieldName("$.modifyTimestamp").As("modifyTimestamp").Numeric().
			FieldName("$.cn").As("cn").Tag().
			FieldName("$.dc").As("dc").Tag().
			FieldName("$.description").As("description").Tag().
			// FieldName("$.destinationIndicator").As("destinationIndicator").Tag().
			FieldName("$.displayName").As("displayName").Tag().
			FieldName("$.gidNumber").As("gidNumber").Numeric().
			FieldName("$.homeDirectory").As("homeDirectory").Tag().
			// FieldName("$.ipHostNumber").As("ipHostNumber").Tag().
			// FieldName("$.mail").As("mail").Tag().
			FieldName("$.member").As("member").Tag().
			FieldName("$.o").As("o").Tag().
			FieldName("$.ou").As("ou").Tag().
			// FieldName("$.owner").As("owner").Tag().
			FieldName("$.sn").As("sn").Tag().
			// FieldName("$.sshPublicKey").As("sshPublicKey").Tag().
			// FieldName("$.telephoneNumber").As("telephoneNumber").Tag().
			// FieldName("$.telexNumber").As("telexNumber").Tag().
			FieldName("$.uid").As("uid").Tag().
			FieldName("$.uidNumber").As("uidNumber").Numeric().
			FieldName("$.userPKCS12").As("userPKCS12").Tag().
			// FieldName("$.userPassword").As("userPassword").Tag().
			// FieldName("$.host_aaa").As("host_aaa").Tag().
			// FieldName("$.host_acl").As("host_acl").Tag().
			// FieldName("$.host_type").As("host_type").Tag().
			// FieldName("$.host_asn").As("host_asn").Tag().
			// FieldName("$.host_upstream_asn").As("host_upstream_asn").Tag().
			// FieldName("$.host_hosting_uuid").As("host_hosting_uuid").Tag().
			// FieldName("$.host_url").As("host_url").Tag().
			// FieldName("$.host_listen").As("host_listen").Tag().

			FieldName("$.labeledURI").As("labeledURI").Tag().
			// FieldName("$.labeledURI[*].key").As("labeledURI_key").Tag().
			// FieldName("$.labeledURI[*].value").As("labeledURI_value").Tag().

			Build()
	})
}
