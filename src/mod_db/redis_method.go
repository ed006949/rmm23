package mod_db

import (
	"context"
	"strconv"
	"strings"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"

	"rmm23/src/l"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_strings"
)

// CreateEntryIndex creates the RediSearch index for the Entry struct.
func (r *RedisRepository) CreateEntryIndex(ctx context.Context) (err error) {
	return r.entry.CreateIndex(ctx, func(schema om.FtCreateSchema) rueidis.Completed {
		return schema.
			FieldName(mod_strings.F_type.FieldName()).As(mod_strings.F_type.String()).Numeric().
			FieldName(mod_strings.F_status.FieldName()).As(mod_strings.F_status.String()).Numeric().
			FieldName(mod_strings.F_baseDN.FieldName()).As(mod_strings.F_baseDN.String()).Tag().Separator(mod_strings.SliceSeparator).

			//
			FieldName(mod_strings.F_uuid.FieldName()).As(mod_strings.F_uuid.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_dn.FieldName()).As(mod_strings.F_dn.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_objectClass.FieldNameSlice()).As(mod_strings.F_objectClass.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_creatorsName.FieldName()).As(mod_strings.F_creatorsName.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_createTimestamp.FieldName()).As(	mod_strings.F_createTimestamp.String()).Numeric().
			FieldName(mod_strings.F_modifiersName.FieldName()).As(mod_strings.F_modifiersName.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_modifyTimestamp.FieldName()).As(	mod_strings.F_modifyTimestamp.String()).Numeric().

			//
			FieldName(mod_strings.F_cn.FieldName()).As(mod_strings.F_cn.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_dc.FieldName()).As(mod_strings.F_dc.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_description.FieldName()).As(	mod_strings.F_description.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_destinationIndicator.FieldNameSlice()).As(mod_strings.F_destinationIndicator.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_displayName.FieldName()).As(	mod_strings.F_displayName.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_gidNumber.FieldName()).As(mod_strings.F_gidNumber.String()).Numeric().
			// FieldName(	mod_strings.F_homeDirectory.FieldName()).As(	mod_strings.F_homeDirectory.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_ipHostNumber.FieldNameSlice()).As(mod_strings.F_ipHostNumber.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_mail.FieldNameSlice()).As(mod_strings.F_mail.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_member.FieldNameSlice()).As(mod_strings.F_member.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_o.FieldName()).As(	mod_strings.F_o.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_ou.FieldName()).As(	mod_strings.F_ou.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_owner.FieldNameSlice()).As(mod_strings.F_owner.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_sn.FieldName()).As(	mod_strings.F_sn.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_sshPublicKey.FieldNameSlice()).As(mod_strings.F_sshPublicKey.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_telephoneNumber.FieldNameSlice()).As(mod_strings.F_telephoneNumber.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_telexNumber.FieldNameSlice()).As(mod_strings.F_telexNumber.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_uid.FieldName()).As(mod_strings.F_uid.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_uidNumber.FieldName()).As(mod_strings.F_uidNumber.String()).Numeric().
			FieldName(mod_strings.F_userPKCS12.FieldNameSlice()).As(mod_strings.F_userPKCS12.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_userPassword.FieldName()).As(	mod_strings.F_userPassword.String()).Tag().Separator(mod_strings.SliceSeparator).

			//
			// FieldName(	mod_strings.F_host_aaa.FieldName()).As(	mod_strings.F_host_aaa.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_host_acl.FieldName()).As(	mod_strings.F_host_acl.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_host_type.FieldName()).As(	mod_strings.F_host_type.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_host_asn.FieldName()).As(	mod_strings.F_host_asn.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_host_upstream_asn.FieldName()).As(	mod_strings.F_host_upstream_asn.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_host_hosting_uuid.FieldName()).As(	mod_strings.F_host_hosting_uuid.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_host_url.FieldName()).As(	mod_strings.F_host_url.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName(	mod_strings.F_host_listen.FieldName()).As(	mod_strings.F_host_listen.String()).Tag().Separator(mod_strings.SliceSeparator).

			FieldName(mod_strings.F_labeledURI.FieldNameSlice()).As(mod_strings.F_labeledURI.String()).Tag().Separator(mod_strings.SliceSeparator).
			// FieldName( _labeledURI.FieldNameSlice() + ".key").As(	mod_strings.F_labeledURI.String() + "_key").Tag().Separator(mod_strings.SliceSeparator).
			// FieldName( _labeledURI.FieldNameSlice() + ".value").As(	mod_strings.F_labeledURI.String() + "_value").Tag().Separator(mod_strings.SliceSeparator).

			//
			Build()
	})
}

// CreateCertIndex creates the RediSearch index for the Cert struct.
func (r *RedisRepository) CreateCertIndex(ctx context.Context) (err error) {
	return r.cert.CreateIndex(ctx, func(schema om.FtCreateSchema) rueidis.Completed {
		return schema.
			// FieldName(mod_strings.F_type.FieldName()).As(mod_strings.F_type.String()).Numeric().
			// FieldName(mod_strings.F_status.FieldName()).As(mod_strings.F_status.String()).Numeric().
			// FieldName(mod_strings.F_baseDN.FieldName()).As(mod_strings.F_baseDN.String()).Tag().Separator(mod_strings.SliceSeparator).

			//
			FieldName(mod_strings.F_uuid.FieldName()).As(mod_strings.F_uuid.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_serialNumber.FieldName()).As(mod_strings.F_serialNumber.String()).Numeric().
			FieldName(mod_strings.F_issuer.FieldName()).As(mod_strings.F_issuer.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_subject.FieldName()).As(mod_strings.F_subject.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_notBefore.FieldName()).As(mod_strings.F_notBefore.String()).Numeric().
			FieldName(mod_strings.F_notAfter.FieldName()).As(mod_strings.F_notAfter.String()).Numeric().
			FieldName(mod_strings.F_dnsNames.FieldNameSlice()).As(mod_strings.F_dnsNames.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_emailAddresses.FieldNameSlice()).As(mod_strings.F_emailAddresses.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_ipAddresses.FieldNameSlice()).As(mod_strings.F_ipAddresses.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_uris.FieldNameSlice()).As(mod_strings.F_uris.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_isCA.FieldName()).As(mod_strings.F_isCA.String()).Numeric().

			//
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
	case !l.Run.DryRunValue():
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

func (r *RedisRepository) waitEntryIndexing(ctx context.Context) (err error) {
	return r.waitIndexing(ctx, r.entry.IndexName())
}

func (r *RedisRepository) waitCertIndexing(ctx context.Context) (err error) {
	return r.waitIndexing(ctx, r.cert.IndexName())
}

func (r *RedisRepository) waitIndexing(ctx context.Context, indexName string) (err error) {
	var (
		resp = r.client.Do(ctx, r.client.B().FtInfo().Index(indexName).Build())
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

	for info["percent_indexed"] != "1" {
		switch info, err = resp.AsStrMap(); {
		case err != nil:
			return
		}
	}

	return
}

func (r *RedisRepository) SaveEntry(ctx context.Context, e *Entry) (err error) {
	switch {
	case l.Run.DryRunValue():
		return
	}

	err = r.entry.Save(ctx, e)
	_ = r.monitorIndexingFailures(ctx)

	return
}

func (r *RedisRepository) SaveCert(ctx context.Context, e *Cert) (err error) {
	switch {
	case l.Run.DryRunValue():
		return
	}

	err = r.cert.Save(ctx, e)
	_ = r.monitorIndexingFailures(ctx)

	return
}

//

func (r *RedisRepository) SaveMultiEntry(ctx context.Context, e ...*Entry) (err []error) {
	switch {
	case l.Run.DryRunValue():
		return
	}

	err = r.entry.SaveMulti(ctx, e...)
	_ = r.monitorIndexingFailures(ctx)

	return
}

func (r *RedisRepository) SaveMultiCert(ctx context.Context, e ...*Cert) (err []error) {
	switch {
	case l.Run.DryRunValue():
		return
	}

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
	switch {
	case l.Run.DryRunValue():
		return
	}

	err = r.entry.Remove(ctx, id)

	return
}

func (r *RedisRepository) DeleteCert(ctx context.Context, id string) (err error) {
	switch {
	case l.Run.DryRunValue():
		return
	}

	err = r.cert.Remove(ctx, id)

	return
}

//

// func (r *RedisRepository) DropIndex(ctx context.Context) (err error) {
// 	l.Z{l.M: "redis", "drop index": "entry", l.E: r.DropEntryIndex(ctx)}.Informational()
// 	l.Z{l.M: "redis", "drop index": "cert", l.E: r.DropCertIndex(ctx)}.Informational()
// 	return
// }

func (r *RedisRepository) DropEntryIndex(ctx context.Context) (err error) {
	switch {
	case l.Run.DryRunValue():
		return
	}

	err = r.entry.DropIndex(ctx)

	return
}

func (r *RedisRepository) DropCertIndex(ctx context.Context) (err error) {
	switch {
	case l.Run.DryRunValue():
		return
	}

	err = r.cert.DropIndex(ctx)

	return
}

//

func (r *RedisRepository) SearchEntryQ(ctx context.Context, query string) (count int64, entries []*Entry, err error) {
	_ = r.waitEntryIndexing(ctx)

	return r.entry.Search(ctx, func(search om.FtSearchIndex) rueidis.Completed {
		return search.Query(query).
			Limit().OffsetNum(0, connMaxPaging).
			Build()
	})
}

func (r *RedisRepository) SearchCertQ(ctx context.Context, query string) (count int64, entries []*Cert, err error) {
	_ = r.waitCertIndexing(ctx)

	return r.cert.Search(ctx, func(search om.FtSearchIndex) rueidis.Completed {
		return search.Query(query).
			Limit().OffsetNum(0, connMaxPaging).
			Build()
	})
}

func (r *RedisRepository) SearchEntryFV(ctx context.Context, field mod_strings.EntryFieldName, value string) (count int64, entries []*Entry, err error) {
	return r.SearchEntryMFV(ctx, MFV{{field, value}})
}

func (r *RedisRepository) SearchCertFV(ctx context.Context, field mod_strings.EntryFieldName, value string) (count int64, entries []*Cert, err error) {
	return r.SearchCertMFV(ctx, MFV{{field, value}})
}

func (r *RedisRepository) SearchEntryMFV(ctx context.Context, mfv MFV) (count int64, entries []*Entry, err error) {
	return r.SearchEntryQ(ctx, mfv.buildMFVQuery())
}

func (r *RedisRepository) SearchCertMFV(ctx context.Context, mfv MFV) (count int64, entries []*Cert, err error) {
	return r.SearchCertQ(ctx, mfv.buildMFVQuery())
}

// SearchEntryMFVField is not working:
//
// err is `unexpected end of JSON input`
//
// JSONRepository receives empty JSON stream.
func (r *RedisRepository) SearchEntryMFVField(ctx context.Context, mfv MFV, field mod_strings.EntryFieldName) (count int64, entries []*Entry, err error) {
	_ = r.waitEntryIndexing(ctx)

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
