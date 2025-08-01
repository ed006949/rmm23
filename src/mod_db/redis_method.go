package mod_db

import (
	"context"
	"strconv"
	"strings"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"

	"rmm23/src/l"
	"rmm23/src/mod_errors"
)

// CreateIndex creates the RediSearch index for the Entry struct.
func (r *RedisRepository) CreateIndex(ctx context.Context) (err error) {
	return r.repo.CreateIndex(ctx, func(schema om.FtCreateSchema) rueidis.Completed {
		return schema.
			FieldName("$." + string(_type)).As(string(_type)).Numeric().
			FieldName("$." + string(_status)).As(string(_status)).Numeric().
			FieldName("$." + string(_baseDN)).As(string(_baseDN)).Tag().Separator(sliceSeparator).
			FieldName("$." + string(_uuid)).As(string(_uuid)).Tag().Separator(sliceSeparator).
			FieldName("$." + string(_dn)).As(string(_dn)).Tag().Separator(sliceSeparator).
			FieldName("$." + string(_objectClass) + "[*]").As(string(_objectClass)).Tag().Separator(sliceSeparator).
			FieldName("$." + string(_creatorsName)).As(string(_creatorsName)).Tag().Separator(sliceSeparator).
			// FieldName("$.createTimestamp").As("createTimestamp").Numeric().
			FieldName("$." + string(_modifiersName)).As(string(_modifiersName)).Tag().Separator(sliceSeparator).
			// FieldName("$.modifyTimestamp").As("modifyTimestamp").Numeric().
			FieldName("$." + string(_cn)).As(string(_cn)).Tag().Separator(sliceSeparator).
			FieldName("$." + string(_dc)).As(string(_dc)).Tag().Separator(sliceSeparator).
			// FieldName("$.description").As("description").Tag().Separator(sliceSeparator).
			FieldName("$." + string(_destinationIndicator) + "[*]").As(string(_destinationIndicator)).Tag().Separator(sliceSeparator).
			// FieldName("$.displayName").As("displayName").Tag().Separator(sliceSeparator).
			FieldName("$." + string(_gidNumber)).As(string(_gidNumber)).Numeric().
			// FieldName("$.homeDirectory").As("homeDirectory").Tag().Separator(sliceSeparator).
			FieldName("$." + string(_ipHostNumber) + "[*]").As(string(_ipHostNumber)).Tag().Separator(sliceSeparator).
			FieldName("$." + string(_mail) + "[*]").As(string(_mail)).Tag().Separator(sliceSeparator).
			FieldName("$." + string(_member) + "[*]").As(string(_member)).Tag().Separator(sliceSeparator).
			// FieldName("$.o").As("o").Tag().Separator(sliceSeparator).
			// FieldName("$.ou").As("ou").Tag().Separator(sliceSeparator).
			FieldName("$." + string(_owner) + "[*]").As(string(_owner)).Tag().Separator(sliceSeparator).
			// FieldName("$.sn").As("sn").Tag().Separator(sliceSeparator).
			FieldName("$." + string(_sshPublicKey) + "[*]").As(string(_sshPublicKey)).Tag().Separator(sliceSeparator).
			FieldName("$." + string(_telephoneNumber) + "[*]").As(string(_telephoneNumber)).Tag().Separator(sliceSeparator).
			FieldName("$." + string(_telexNumber) + "[*]").As(string(_telexNumber)).Tag().Separator(sliceSeparator).
			FieldName("$." + string(_uid)).As(string(_uid)).Tag().Separator(sliceSeparator).
			FieldName("$." + string(_uidNumber)).As(string(_uidNumber)).Numeric().
			FieldName("$." + string(_userPKCS12) + "[*]").As(string(_userPKCS12)).Tag().Separator(sliceSeparator).
			// FieldName("$.userPassword").As("userPassword").Tag().Separator(sliceSeparator).

			// FieldName("$.host_aaa").As("host_aaa").Tag().Separator(sliceSeparator).
			// FieldName("$.host_acl").As("host_acl").Tag().Separator(sliceSeparator).
			// FieldName("$.host_type").As("host_type").Tag().Separator(sliceSeparator).
			// FieldName("$.host_asn").As("host_asn").Tag().Separator(sliceSeparator).
			// FieldName("$.host_upstream_asn").As("host_upstream_asn").Tag().Separator(sliceSeparator).
			// FieldName("$.host_hosting_uuid").As("host_hosting_uuid").Tag().Separator(sliceSeparator).
			// FieldName("$.host_url").As("host_url").Tag().Separator(sliceSeparator).
			// FieldName("$.host_listen").As("host_listen").Tag().Separator(sliceSeparator).

			FieldName("$." + string(_labeledURI) + "[*]").As(string(_labeledURI)).Tag().Separator(sliceSeparator).
			// FieldName("$.labeledURI[*].key").As("labeledURI_key").Tag().Separator(sliceSeparator).
			// FieldName("$.labeledURI[*].value").As("labeledURI_value").Tag().Separator(sliceSeparator).

			Build()
	})
}

func (r *Conf) Dial(ctx context.Context) (err error) {
	switch r.client, err = rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{r.URL.Host},
	}); {
	case err != nil:
		return
	}

	r.repo = NewRedisRepository(r.client)

	switch {
	case l.CLEAR:
		_ = r.repo.DropIndex(ctx)

		switch err = r.repo.CreateIndex(ctx); {
		case err != nil:
			return
		}
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
		l.Z{l.M: "redis resp", l.E: err}.Error()

		return
	}

	var (
		info map[string]string
	)

	switch info, err = resp.AsStrMap(); {
	case err != nil:
		l.Z{l.M: "redis info", l.E: err}.Error()

		return
	}

	for a, b := range info {
		switch a {
		case "hash_indexing_failures":
			switch c, d := strconv.ParseInt(b, 10, 64); {
			case d == nil && c == 0:
				l.Z{l.M: "redis", a: b}.Informational()
			case d == nil:
				l.Z{l.M: "redis", a: b}.Warning()
			default:
				l.Z{l.M: "redis info", l.E: err, a: b}.Error()
			}
		}
	}

	return
}

func (r *RedisRepository) SaveEntry(ctx context.Context, e *Entry) error {
	return r.repo.Save(ctx, e)
}

func (r *RedisRepository) FindEntry(ctx context.Context, id string) (*Entry, error) {
	return r.repo.Fetch(ctx, id)
}

func (r *RedisRepository) DeleteEntry(ctx context.Context, id string) error {
	return r.repo.Remove(ctx, id)
}

func (r *RedisRepository) DropIndex(ctx context.Context) (err error) { return r.repo.DropIndex(ctx) }

func (r *RedisRepository) SearchQ(ctx context.Context, query string) (count int64, entries []*Entry, err error) {
	return r.repo.Search(ctx, func(search om.FtSearchIndex) rueidis.Completed {
		return search.Query(query).
			Limit().OffsetNum(0, connMaxPaging).
			Build()
	})
}

func (r *RedisRepository) SearchFV(ctx context.Context, field entryFieldName, value string) (count int64, entries []*Entry, err error) {
	return r.SearchMFV(ctx, _MFV{{field, value}})
}

func (r *RedisRepository) SearchMFV(ctx context.Context, mfv _MFV) (count int64, entries []*Entry, err error) {
	return r.SearchQ(ctx, mfv.buildMFVQuery())
}

// SearchMFVField is not working:
//
// err is `unexpected end of JSON input`
//
// JSONRepository receives empty JSON stream.
func (r *RedisRepository) SearchMFVField(ctx context.Context, mfv _MFV, field entryFieldName) (count int64, entries []*Entry, err error) {
	return r.repo.Search(ctx, func(search om.FtSearchIndex) rueidis.Completed {
		var (
			command = search.Query(mfv.buildMFVQuery()).
				Return(strconv.FormatInt(1, 10)).
				Identifier(field.String()).
				Limit().OffsetNum(0, connMaxPaging).
				Build()
		)
		l.Z{l.M: "redis", "command": strings.Join(command.Commands(), " ")}.Informational()

		return command
	})
}

func (r *_MFV) buildMFVQuery() (outbound string) {
	var (
		interim = make([]string, len(*r), len(*r))
	)

	for i, fv := range *r {
		interim[i] = buildFVQuery(fv._F, fv._V)
	}

	return strings.Join(interim, " ")
}
