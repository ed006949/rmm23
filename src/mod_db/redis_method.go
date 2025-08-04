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
			FieldName(F_type.FieldName()).As(F_type.String()).Numeric().
			FieldName(F_status.FieldName()).As(F_status.String()).Numeric().
			FieldName(F_baseDN.FieldName()).As(F_baseDN.String()).Tag().Separator(sliceSeparator).

			//
			FieldName(F_uuid.FieldName()).As(F_uuid.String()).Tag().Separator(sliceSeparator).
			FieldName(F_dn.FieldName()).As(F_dn.String()).Tag().Separator(sliceSeparator).
			FieldName(F_objectClass.FieldNameSlice()).As(F_objectClass.String()).Tag().Separator(sliceSeparator).
			FieldName(F_creatorsName.FieldName()).As(F_creatorsName.String()).Tag().Separator(sliceSeparator).
			// FieldName(	F_createTimestamp.FieldName()).As(	F_createTimestamp.String()).Numeric().
			FieldName(F_modifiersName.FieldName()).As(F_modifiersName.String()).Tag().Separator(sliceSeparator).
			// FieldName(	F_modifyTimestamp.FieldName()).As(	F_modifyTimestamp.String()).Numeric().

			//
			FieldName(F_cn.FieldName()).As(F_cn.String()).Tag().Separator(sliceSeparator).
			FieldName(F_dc.FieldName()).As(F_dc.String()).Tag().Separator(sliceSeparator).
			// FieldName(	F_description.FieldName()).As(	F_description.String()).Tag().Separator(sliceSeparator).
			FieldName(F_destinationIndicator.FieldNameSlice()).As(F_destinationIndicator.String()).Tag().Separator(sliceSeparator).
			// FieldName(	F_displayName.FieldName()).As(	F_displayName.String()).Tag().Separator(sliceSeparator).
			FieldName(F_gidNumber.FieldName()).As(F_gidNumber.String()).Numeric().
			// FieldName(	F_homeDirectory.FieldName()).As(	F_homeDirectory.String()).Tag().Separator(sliceSeparator).
			FieldName(F_ipHostNumber.FieldNameSlice()).As(F_ipHostNumber.String()).Tag().Separator(sliceSeparator).
			FieldName(F_mail.FieldNameSlice()).As(F_mail.String()).Tag().Separator(sliceSeparator).
			FieldName(F_member.FieldNameSlice()).As(F_member.String()).Tag().Separator(sliceSeparator).
			// FieldName(	F_o.FieldName()).As(	F_o.String()).Tag().Separator(sliceSeparator).
			// FieldName(	F_ou.FieldName()).As(	F_ou.String()).Tag().Separator(sliceSeparator).
			FieldName(F_owner.FieldNameSlice()).As(F_owner.String()).Tag().Separator(sliceSeparator).
			// FieldName(	F_sn.FieldName()).As(	F_sn.String()).Tag().Separator(sliceSeparator).
			FieldName(F_sshPublicKey.FieldNameSlice()).As(F_sshPublicKey.String()).Tag().Separator(sliceSeparator).
			FieldName(F_telephoneNumber.FieldNameSlice()).As(F_telephoneNumber.String()).Tag().Separator(sliceSeparator).
			FieldName(F_telexNumber.FieldNameSlice()).As(F_telexNumber.String()).Tag().Separator(sliceSeparator).
			FieldName(F_uid.FieldName()).As(F_uid.String()).Tag().Separator(sliceSeparator).
			FieldName(F_uidNumber.FieldName()).As(F_uidNumber.String()).Numeric().
			FieldName(F_userPKCS12.FieldNameSlice()).As(F_userPKCS12.String()).Tag().Separator(sliceSeparator).
			// FieldName(	F_userPassword.FieldName()).As(	F_userPassword.String()).Tag().Separator(sliceSeparator).

			//
			// FieldName(	F_host_aaa.FieldName()).As(	F_host_aaa.String()).Tag().Separator(sliceSeparator).
			// FieldName(	F_host_acl.FieldName()).As(	F_host_acl.String()).Tag().Separator(sliceSeparator).
			// FieldName(	F_host_type.FieldName()).As(	F_host_type.String()).Tag().Separator(sliceSeparator).
			// FieldName(	F_host_asn.FieldName()).As(	F_host_asn.String()).Tag().Separator(sliceSeparator).
			// FieldName(	F_host_upstream_asn.FieldName()).As(	F_host_upstream_asn.String()).Tag().Separator(sliceSeparator).
			// FieldName(	F_host_hosting_uuid.FieldName()).As(	F_host_hosting_uuid.String()).Tag().Separator(sliceSeparator).
			// FieldName(	F_host_url.FieldName()).As(	F_host_url.String()).Tag().Separator(sliceSeparator).
			// FieldName(	F_host_listen.FieldName()).As(	F_host_listen.String()).Tag().Separator(sliceSeparator).

			FieldName(F_labeledURI.FieldNameSlice()).As(F_labeledURI.String()).Tag().Separator(sliceSeparator).
			// FieldName( _labeledURI.FieldNameSlice() + ".key").As(	F_labeledURI.String() + "_key").Tag().Separator(sliceSeparator).
			// FieldName( _labeledURI.FieldNameSlice() + ".value").As(	F_labeledURI.String() + "_value").Tag().Separator(sliceSeparator).
			Build()
	})
}

// CreateCertIndex creates the RediSearch index for the Cert struct.
func (r *RedisRepository) CreateCertIndex(ctx context.Context) (err error) {
	return r.cert.CreateIndex(ctx, func(schema om.FtCreateSchema) rueidis.Completed {
		return schema.
			FieldName(F_type.FieldName()).As(F_type.String()).Numeric().
			FieldName(F_status.FieldName()).As(F_status.String()).Numeric().
			FieldName(F_baseDN.FieldName()).As(F_baseDN.String()).Tag().Separator(sliceSeparator).

			//
			FieldName(F_uuid.FieldName()).As(F_uuid.String()).Tag().Separator(sliceSeparator).
			FieldName(F_dn.FieldName()).As(F_dn.String()).Tag().Separator(sliceSeparator).

			//
			FieldName(F_creatorsName.FieldName()).As(F_creatorsName.String()).Tag().Separator(sliceSeparator).
			// FieldName(	F_createTimestamp.FieldName()).As(	F_createTimestamp.String()).Numeric().
			FieldName(F_modifiersName.FieldName()).As(F_modifiersName.String()).Tag().Separator(sliceSeparator).
			// FieldName(	F_modifyTimestamp.FieldName()).As(	F_modifyTimestamp.String()).Numeric().

			Build()
	})
}

func (r *Conf) Dial(ctx context.Context) (err error) {
	var (
		client rueidis.Client
	)
	switch client, err = rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{r.URL.Host},
	}); {
	case err != nil:
		return
	}

	r.Repo = NewRedisRepository(client)

	switch {
	case l.CLEAR:
		_ = r.Repo.DropEntryIndex(ctx)

		_ = r.Repo.DropCertIndex(ctx)
		switch err = r.Repo.CreateEntryIndex(ctx); {
		case err != nil:
			return
		}

		switch err = r.Repo.CreateCertIndex(ctx); {
		case err != nil:
			return
		}
	}

	_ = r.Repo.monitorIndexingFailures(ctx)

	return
}

func (r *Conf) Close() (err error) {
	switch {
	case r.Repo.client == nil:
		return mod_errors.ENoConn
	}

	r.Repo.client.Close()

	return
}

func (r *RedisRepository) monitorIndexingFailures(ctx context.Context) (err error) {
	for indexName, resp := range map[string]rueidis.RedisResult{
		r.entry.IndexName(): r.client.Do(ctx, r.client.B().FtInfo().Index(r.entry.IndexName()).Build()),
		r.cert.IndexName():  r.client.Do(ctx, r.client.B().FtInfo().Index(r.cert.IndexName()).Build()),
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

func (r *MFV) buildMFVQuery() (outbound string) {
	var (
		interim = make([]string, len(*r), len(*r))
	)

	for i, fv := range *r {
		interim[i] = buildFVQuery(fv.Field, fv.Value)
	}

	return strings.Join(interim, " ")
}

//

func (r *RedisRepository) SaveEntry(ctx context.Context, e *Entry) (err error) {
	err = r.entry.Save(ctx, e)
	_ = r.monitorIndexingFailures(ctx)

	return
}

func (r *RedisRepository) SaveCert(ctx context.Context, e *Cert) (err error) {
	err = r.cert.Save(ctx, e)
	_ = r.monitorIndexingFailures(ctx)

	return
}

//

func (r *RedisRepository) SaveMultiEntry(ctx context.Context, e ...*Entry) (err []error) {
	err = r.entry.SaveMulti(ctx, e...)
	_ = r.monitorIndexingFailures(ctx)

	return
}

func (r *RedisRepository) SaveMultiCert(ctx context.Context, e ...*Cert) (err []error) {
	err = r.cert.SaveMulti(ctx, e...)
	_ = r.monitorIndexingFailures(ctx)

	return
}

//

func (r *RedisRepository) FindEntry(ctx context.Context, id string) (entry *Entry, err error) {
	return r.entry.Fetch(ctx, id)
}

func (r *RedisRepository) FindCert(ctx context.Context, id string) (cert *Cert, err error) {
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

func (r *RedisRepository) SearchCertQ(ctx context.Context, query string) (count int64, entries []*Cert, err error) {
	return r.cert.Search(ctx, func(search om.FtSearchIndex) rueidis.Completed {
		return search.Query(query).
			Limit().OffsetNum(0, connMaxPaging).
			Build()
	})
}

func (r *RedisRepository) SearchEntryFV(ctx context.Context, field entryFieldName, value string) (count int64, entries []*Entry, err error) {
	return r.SearchEntryMFV(ctx, MFV{{field, value}})
}

func (r *RedisRepository) SearchEntryMFV(ctx context.Context, mfv MFV) (count int64, entries []*Entry, err error) {
	return r.SearchEntryQ(ctx, mfv.buildMFVQuery())
}

// SearchEntryMFVField is not working:
//
// err is `unexpected end of JSON input`
//
// JSONRepository receives empty JSON stream.
func (r *RedisRepository) SearchEntryMFVField(ctx context.Context, mfv MFV, field entryFieldName) (count int64, entries []*Entry, err error) {
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
