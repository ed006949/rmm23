package mod_db

import (
	"context"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"

	"rmm23/src/l"
	"rmm23/src/mod_errors"
)

func (r *Conf) dial() (err error) {
	switch r.rcNetwork, err = r.URL.RedisNetwork(); {
	case err != nil:
		return
	}

	r.rsClient = redisearch.NewClientFromPool(&redis.Pool{
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.DialContext(ctx, r.rcNetwork, r.URL.Host, redis.DialDatabase(0))
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
	}, r.Name)

	return
}

func (r *Conf) createIndex() (err error) {
	var (
		indexInfo *redisearch.IndexInfo
	)

	// define indexDefinition
	r.indexDefinition = redisearch.NewIndexDefinition().AddPrefix(entryDocIDHeader)

	switch l.CLEAR {
	case true:
		_ = r.rsClient.Drop()
	// _ = r.rsClient.DropIndex(true)
	default:
		_ = r.rsClient.DropIndex(false)
	}

	switch swErr := r.rsClient.CreateIndexWithIndexDefinition(r.schema, r.indexDefinition); {
	case mod_errors.Contains(swErr, mod_errors.EIndexExist):
	case swErr != nil:
		return swErr
	}

	// wait for index to complete
	switch indexInfo, err = r.rsClient.Info(); {
	case err != nil:
		return
	}

	for indexInfo.IsIndexing {
		switch indexInfo, err = r.rsClient.Info(); {
		case err != nil:
			return
		}
	}

	return
}

func (r *Conf) getDoc(inbound string) (outbound *redisearch.Document, err error) {
	return r.rsClient.Get(inbound)
}

func (r *Conf) getDocByUUID(inbound attrUUID) (outbound *redisearch.Document, err error) {
	return r.getDoc(inbound.Entry())
}

func (r *Conf) getDocsByKV(key entryFieldName, value string) (outbound []redisearch.Document, count int, err error) {
	var (
		interim string
	)

	switch r.schemaMap[key] {
	case redisearch.TextField:
		interim = "@" + key.String() + ":" + escapeQueryValue(value) + ""
	case redisearch.NumericField:
		interim = "@" + key.String() + ":[" + escapeQueryValue(value) + "]"
	case redisearch.GeoField:
		interim = "@" + key.String() + ":[" + escapeQueryValue(value) + "]"
	case redisearch.TagField:
		interim = "@" + key.String() + ":{" + escapeQueryValue(value) + "}"
	default:
		return nil, 0, mod_errors.EUnwilling
	}

	return r.rsClient.Search(redisearch.NewQuery(interim).SetInFields(key.String()).Limit(0, connMaxPaging))
}

// SaveEntry saves an entry to Redis.
func (r *RedisRepository) SaveEntry(ctx context.Context, e *entry) error {
	return r.repo.Save(ctx, e)
}

// FindEntry finds an entry by its ID.
func (r *RedisRepository) FindEntry(ctx context.Context, id string) (*entry, error) {
	return r.repo.Fetch(ctx, id)
}

// CreateIndex creates the RediSearch index for the entry struct.
func (r *RedisRepository) CreateIndex(ctx context.Context) error {
	return r.repo.CreateIndex(ctx, func(schema om.FtCreateSchema) rueidis.Completed {
		return schema.
			FieldName("$.Type").As("type").Numeric().Sortable().
			FieldName("$.Status").As("status").Numeric().Sortable().
			FieldName("$.BaseDN").As("baseDN").Tag().Sortable().
			FieldName("$.UUID").As("uuid").Tag().Sortable().
			FieldName("$.DN").As("dn").Tag().Sortable().
			FieldName("$.ObjectClass").As("objectClass").Tag().
			FieldName("$.CreatorsName").As("creatorsName").Tag().
			FieldName("$.CreateTimestamp").As("createTimestamp").Tag().
			FieldName("$.ModifiersName").As("modifiersName").Tag().
			FieldName("$.ModifyTimestamp").As("modifyTimestamp").Tag().
			FieldName("$.CN").As("cn").Tag().
			FieldName("$.DC").As("dc").Tag().Sortable().
			FieldName("$.Description").As("description").Tag().
			FieldName("$.DestinationIndicator").As("destinationIndicator").Tag().
			FieldName("$.DisplayName").As("displayName").Tag().Sortable().
			FieldName("$.GIDNumber").As("gidNumber").Numeric().Sortable().
			FieldName("$.HomeDirectory").As("homeDirectory").Tag().
			FieldName("$.IPHostNumber").As("ipHostNumber").Tag().Sortable().
			FieldName("$.Mail").As("mail").Tag().
			FieldName("$.Member").As("member").Tag().Sortable().
			FieldName("$.O").As("o").Tag().
			FieldName("$.OU").As("ou").Tag().
			FieldName("$.Owner").As("owner").Tag().
			FieldName("$.SN").As("sn").Tag().
			FieldName("$.SSHPublicKey").As("sshPublicKey").Tag().
			FieldName("$.TelephoneNumber").As("telephoneNumber").Tag().
			FieldName("$.TelexNumber").As("telexNumber").Tag().
			FieldName("$.UID").As("uid").Tag().Sortable().
			FieldName("$.UIDNumber").As("uidNumber").Numeric().Sortable().
			FieldName("$.UserPKCS12").As("userPKCS12").Tag().
			FieldName("$.UserPassword").As("userPassword").Tag().
			FieldName("$.AAA").As("host_aaa").Tag().
			FieldName("$.ACL").As("host_acl").Tag().
			FieldName("$.HostType").As("host_type").Tag().
			FieldName("$.HostASN").As("host_asn").Numeric().Sortable().
			FieldName("$.HostUpstreamASN").As("host_upstream_asn").Numeric().
			FieldName("$.HostHostingUUID").As("host_hosting_uuid").Numeric().
			FieldName("$.HostURL").As("host_url").Tag().Sortable().
			FieldName("$.HostListen").As("host_listen").Tag().Sortable().
			FieldName("$.LabeledURI").As("labeledURI").Tag().
			Build()
	})
}
