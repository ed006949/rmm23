package mod_db

import (
	"context"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"

	"rmm23/src/mod_errors"
)

func (r *Conf) Dial(ctx context.Context) (err error) {
	switch r.client, err = rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{r.URL.Host},
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

// CreateIndex creates the RediSearch index for the Entry struct.
func (r *RedisRepository) DropIndex(ctx context.Context) (err error) { return r.repo.DropIndex(ctx) }

// CreateIndex creates the RediSearch index for the Entry struct.
func (r *RedisRepository) CreateIndex(ctx context.Context) (err error) {
	return r.repo.CreateIndex(ctx, func(schema om.FtCreateSchema) rueidis.Completed {
		return schema.
			FieldName("$.Key").As("key").Tag().
			FieldName("$.Ver").As("ver").Numeric().Sortable().
			FieldName("$.Type").As("type").Numeric().Sortable().
			FieldName("$.Status").As("status").Numeric().Sortable().
			FieldName("$.BaseDN").As("baseDN").Tag().
			FieldName("$.UUID").As("uuid").Tag().
			FieldName("$.DN").As("dn").Tag().
			FieldName("$.ObjectClass").As("objectClass").Tag().
			FieldName("$.CreatorsName").As("creatorsName").Tag().
			FieldName("$.CreateTimestamp").As("createTimestamp").Tag().
			FieldName("$.ModifiersName").As("modifiersName").Tag().
			FieldName("$.ModifyTimestamp").As("modifyTimestamp").Tag().
			FieldName("$.CN").As("cn").Tag().
			FieldName("$.DC").As("dc").Tag().
			FieldName("$.Description").As("description").Tag().
			FieldName("$.DestinationIndicator").As("destinationIndicator").Tag().
			FieldName("$.DisplayName").As("displayName").Tag().
			FieldName("$.GIDNumber").As("gidNumber").Numeric().Sortable().
			FieldName("$.HomeDirectory").As("homeDirectory").Tag().
			FieldName("$.IPHostNumber").As("ipHostNumber").Tag().
			FieldName("$.Mail").As("mail").Tag().
			FieldName("$.Member").As("member").Tag().
			FieldName("$.O").As("o").Tag().
			FieldName("$.OU").As("ou").Tag().
			FieldName("$.Owner").As("owner").Tag().
			FieldName("$.SN").As("sn").Tag().
			FieldName("$.SSHPublicKey").As("sshPublicKey").Tag().
			FieldName("$.TelephoneNumber").As("telephoneNumber").Tag().
			FieldName("$.TelexNumber").As("telexNumber").Tag().
			FieldName("$.UID").As("uid").Tag().
			FieldName("$.UIDNumber").As("uidNumber").Numeric().Sortable().
			FieldName("$.UserPKCS12").As("userPKCS12").Tag().
			FieldName("$.UserPassword").As("userPassword").Tag().
			FieldName("$.AAA").As("host_aaa").Tag().
			FieldName("$.ACL").As("host_acl").Tag().
			FieldName("$.HostType").As("host_type").Tag().
			FieldName("$.HostASN").As("host_asn").Numeric().Sortable().
			FieldName("$.HostUpstreamASN").As("host_upstream_asn").Numeric().
			FieldName("$.HostHostingUUID").As("host_hosting_uuid").Numeric().
			FieldName("$.HostURL").As("host_url").Tag().
			FieldName("$.HostListen").As("host_listen").Tag().
			FieldName("$.LabeledURI").As("labeledURI").Tag().
			Build()
	})
}
