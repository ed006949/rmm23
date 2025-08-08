package mod_db

import (
	"context"
	"strconv"
	"strings"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"

	"rmm23/src/l"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_reflect"
	"rmm23/src/mod_strings"
)

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
			case ftInfo_hash_indexing_failures:
				switch c, d := strconv.ParseInt(b, 10, 64); {
				case d == nil && c == 0:
					l.Z{l.M: "redis", "index": info[ftInfo_index_name], a: b}.Debug()
				case d == nil:
					l.Z{l.M: "redis", "index": info[ftInfo_index_name], a: b}.Warning()
				default:
					l.Z{l.M: "redis info", "index": info[ftInfo_index_name], l.E: err, a: b}.Error()
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

	for info, forErr := resp.AsStrMap(); info[ftInfo_percent_indexed] != "1"; info, forErr = resp.AsStrMap() {
		switch {
		case forErr != nil:
			return
		}

		_ = mod_reflect.WaitCtx(ctx, l.RetryInterval)
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
