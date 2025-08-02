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

// CreateEntryIndex creates the RediSearch index for the Entry struct.
func (r *RedisRepository) CreateEntryIndex(ctx context.Context) (err error) {
	return r.entry.CreateIndex(ctx, func(schema om.FtCreateSchema) rueidis.Completed {
		return schema.
			FieldName(_type.FieldName()).As(_type.String()).Numeric().
			FieldName(_status.FieldName()).As(_status.String()).Numeric().
			FieldName(_baseDN.FieldName()).As(_baseDN.String()).Tag().Separator(sliceSeparator).

			//
			FieldName(_uuid.FieldName()).As(_uuid.String()).Tag().Separator(sliceSeparator).
			FieldName(_dn.FieldName()).As(_dn.String()).Tag().Separator(sliceSeparator).
			FieldName(_objectClass.FieldNameSlice()).As(_objectClass.String()).Tag().Separator(sliceSeparator).
			FieldName(_creatorsName.FieldName()).As(_creatorsName.String()).Tag().Separator(sliceSeparator).
			// FieldName(_createTimestamp.FieldName()).As(_createTimestamp.String()).Numeric().
			FieldName(_modifiersName.FieldName()).As(_modifiersName.String()).Tag().Separator(sliceSeparator).
			// FieldName(_modifyTimestamp.FieldName()).As(_modifyTimestamp.String()).Numeric().

			//
			FieldName(_cn.FieldName()).As(_cn.String()).Tag().Separator(sliceSeparator).
			FieldName(_dc.FieldName()).As(_dc.String()).Tag().Separator(sliceSeparator).
			// FieldName(_description.FieldName()).As(_description.String()).Tag().Separator(sliceSeparator).
			FieldName(_destinationIndicator.FieldNameSlice()).As(_destinationIndicator.String()).Tag().Separator(sliceSeparator).
			// FieldName(_displayName.FieldName()).As(_displayName.String()).Tag().Separator(sliceSeparator).
			FieldName(_gidNumber.FieldName()).As(_gidNumber.String()).Numeric().
			// FieldName(_homeDirectory.FieldName()).As(_homeDirectory.String()).Tag().Separator(sliceSeparator).
			FieldName(_ipHostNumber.FieldNameSlice()).As(_ipHostNumber.String()).Tag().Separator(sliceSeparator).
			FieldName(_mail.FieldNameSlice()).As(_mail.String()).Tag().Separator(sliceSeparator).
			FieldName(_member.FieldNameSlice()).As(_member.String()).Tag().Separator(sliceSeparator).
			// FieldName(_o.FieldName()).As(_o.String()).Tag().Separator(sliceSeparator).
			// FieldName(_ou.FieldName()).As(_ou.String()).Tag().Separator(sliceSeparator).
			FieldName(_owner.FieldNameSlice()).As(_owner.String()).Tag().Separator(sliceSeparator).
			// FieldName(_sn.FieldName()).As(_sn.String()).Tag().Separator(sliceSeparator).
			FieldName(_sshPublicKey.FieldNameSlice()).As(_sshPublicKey.String()).Tag().Separator(sliceSeparator).
			FieldName(_telephoneNumber.FieldNameSlice()).As(_telephoneNumber.String()).Tag().Separator(sliceSeparator).
			FieldName(_telexNumber.FieldNameSlice()).As(_telexNumber.String()).Tag().Separator(sliceSeparator).
			FieldName(_uid.FieldName()).As(_uid.String()).Tag().Separator(sliceSeparator).
			FieldName(_uidNumber.FieldName()).As(_uidNumber.String()).Numeric().
			FieldName(_userPKCS12.FieldNameSlice()).As(_userPKCS12.String()).Tag().Separator(sliceSeparator).
			// FieldName(_userPassword.FieldName()).As(_userPassword.String()).Tag().Separator(sliceSeparator).

			//
			// FieldName(_host_aaa.FieldName()).As(_host_aaa.String()).Tag().Separator(sliceSeparator).
			// FieldName(_host_acl.FieldName()).As(_host_acl.String()).Tag().Separator(sliceSeparator).
			// FieldName(_host_type.FieldName()).As(_host_type.String()).Tag().Separator(sliceSeparator).
			// FieldName(_host_asn.FieldName()).As(_host_asn.String()).Tag().Separator(sliceSeparator).
			// FieldName(_host_upstream_asn.FieldName()).As(_host_upstream_asn.String()).Tag().Separator(sliceSeparator).
			// FieldName(_host_hosting_uuid.FieldName()).As(_host_hosting_uuid.String()).Tag().Separator(sliceSeparator).
			// FieldName(_host_url.FieldName()).As(_host_url.String()).Tag().Separator(sliceSeparator).
			// FieldName(_host_listen.FieldName()).As(_host_listen.String()).Tag().Separator(sliceSeparator).

			FieldName(_labeledURI.FieldNameSlice()).As(_labeledURI.String()).Tag().Separator(sliceSeparator).
			// FieldName( _labeledURI.FieldNameSlice() + ".key").As(_labeledURI.String() + "_key").Tag().Separator(sliceSeparator).
			// FieldName( _labeledURI.FieldNameSlice() + ".value").As(_labeledURI.String() + "_value").Tag().Separator(sliceSeparator).
			Build()
	})
}

// CreateCertIndex creates the RediSearch index for the Certificate struct.
func (r *RedisRepository) CreateCertIndex(ctx context.Context) (err error) {
	return r.cert.CreateIndex(ctx, func(schema om.FtCreateSchema) rueidis.Completed {
		return schema.
			FieldName(_type.FieldName()).As(_type.String()).Numeric().
			FieldName(_status.FieldName()).As(_status.String()).Numeric().
			FieldName(_baseDN.FieldName()).As(_baseDN.String()).Tag().Separator(sliceSeparator).

			//
			FieldName(_uuid.FieldName()).As(_uuid.String()).Tag().Separator(sliceSeparator).
			FieldName(_dn.FieldName()).As(_dn.String()).Tag().Separator(sliceSeparator).

			//
			FieldName(_creatorsName.FieldName()).As(_creatorsName.String()).Tag().Separator(sliceSeparator).
			// FieldName(_createTimestamp.FieldName()).As(_createTimestamp.String()).Numeric().
			FieldName(_modifiersName.FieldName()).As(_modifiersName.String()).Tag().Separator(sliceSeparator).
			// FieldName(_modifyTimestamp.FieldName()).As(_modifyTimestamp.String()).Numeric().

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
		_ = r.repo.DropEntryIndex(ctx)
		_ = r.repo.DropCertIndex(ctx)

		switch err = r.repo.CreateEntryIndex(ctx); {
		case err != nil:
			return
		}

		switch err = r.repo.CreateCertIndex(ctx); {
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
	for indexName, resp := range map[string]rueidis.RedisResult{
		r.repo.entry.IndexName(): r.client.Do(ctx, r.client.B().FtInfo().Index(r.repo.entry.IndexName()).Build()),
		r.repo.cert.IndexName():  r.client.Do(ctx, r.client.B().FtInfo().Index(r.repo.cert.IndexName()).Build()),
	} {
		switch err = resp.Error(); {
		case err != nil:
			l.Z{l.M: "redis resp", "index": indexName, l.E: err}.Error()

			continue
		}

		var (
			info map[string]string
		)

		switch info, err = resp.AsStrMap(); {
		case err != nil:
			l.Z{l.M: "redis info", "index": indexName, l.E: err}.Error()

			continue
		}

		for a, b := range info {
			switch a {
			case "hash_indexing_failures":
				switch c, d := strconv.ParseInt(b, 10, 64); {
				case d == nil && c == 0:
					l.Z{l.M: "redis", "index": info["index_name"], a: b}.Debug()
				case d == nil:
					l.Z{l.M: "redis", "index": info["index_name"], a: b}.Warning()
				default:
					l.Z{l.M: "redis info", "index": info["index_name"], l.E: err, a: b}.Error()
				}
			}
		}
	}

	return nil
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

//

func (r *RedisRepository) SaveEntry(ctx context.Context, e *Entry) (err error) {
	return r.entry.Save(ctx, e)
}

func (r *RedisRepository) SaveCert(ctx context.Context, e *Certificate) (err error) {
	return r.cert.Save(ctx, e)
}

func (r *RedisRepository) SaveMultiEntry(ctx context.Context, e ...*Entry) (err []error) {
	return r.entry.SaveMulti(ctx, e...)
}

func (r *RedisRepository) SaveMultiCert(ctx context.Context, e ...*Certificate) (err []error) {
	return r.cert.SaveMulti(ctx, e...)
}

//

func (r *RedisRepository) FindEntry(ctx context.Context, id string) (entry *Entry, err error) {
	return r.entry.Fetch(ctx, id)
}

func (r *RedisRepository) FindCert(ctx context.Context, id string) (cert *Certificate, err error) {
	return r.cert.Fetch(ctx, id)
}

//

func (r *RedisRepository) DeleteEntry(ctx context.Context, id string) (err error) {
	return r.entry.Remove(ctx, id)
}

func (r *RedisRepository) DeleteCert(ctx context.Context, id string) (err error) {
	return r.cert.Remove(ctx, id)
}

//

// func (r *RedisRepository) DropIndex(ctx context.Context) (err error) {
// 	l.Z{l.M: "redis", "drop index": "entry", l.E: r.DropEntryIndex(ctx)}.Informational()
// 	l.Z{l.M: "redis", "drop index": "cert", l.E: r.DropCertIndex(ctx)}.Informational()
// 	return
// }

func (r *RedisRepository) DropEntryIndex(ctx context.Context) (err error) {
	return r.entry.DropIndex(ctx)
}

func (r *RedisRepository) DropCertIndex(ctx context.Context) (err error) {
	return r.cert.DropIndex(ctx)
}

//

func (r *RedisRepository) SearchEntryQ(ctx context.Context, query string) (count int64, entries []*Entry, err error) {
	return r.entry.Search(ctx, func(search om.FtSearchIndex) rueidis.Completed {
		return search.Query(query).
			Limit().OffsetNum(0, connMaxPaging).
			Build()
	})
}

func (r *RedisRepository) SearchEntryFV(ctx context.Context, field entryFieldName, value string) (count int64, entries []*Entry, err error) {
	return r.SearchEntryMFV(ctx, _MFV{{field, value}})
}

func (r *RedisRepository) SearchEntryMFV(ctx context.Context, mfv _MFV) (count int64, entries []*Entry, err error) {
	return r.SearchEntryQ(ctx, mfv.buildMFVQuery())
}

// SearchEntryMFVField is not working:
//
// err is `unexpected end of JSON input`
//
// JSONRepository receives empty JSON stream.
func (r *RedisRepository) SearchEntryMFVField(ctx context.Context, mfv _MFV, field entryFieldName) (count int64, entries []*Entry, err error) {
	return r.entry.Search(ctx, func(search om.FtSearchIndex) rueidis.Completed {
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
