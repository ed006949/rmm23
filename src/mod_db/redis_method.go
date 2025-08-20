package mod_db

import (
	"context"
	"encoding/json/v2"
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

	r.Repo = NewRedisRepository(ctx, client)

	switch {
	case !l.Run.DryRunValue():
		_ = r.Repo.DropEntryIndex()
		_ = r.Repo.DropCertIndex()

		switch err = r.Repo.CreateEntryIndex(); {
		case err != nil:
			return
		}

		switch err = r.Repo.CreateCertIndex(); {
		case err != nil:
			return
		}
	}

	switch err = r.Repo.getInfo(); {
	case err != nil:
		return
	}

	return
}

func (r *RedisRepository) getInfo() (err error) {
	mod_reflect.MakeMapIfNil(&r.info)

	var (
		repos []string
	)
	switch repos, err = r.client.Do(r.ctx, r.client.B().FtList().Build()).AsStrSlice(); {
	case err != nil:
		return
	}

	for _, b := range repos {
		var (
			redisResult    = r.client.Do(r.ctx, r.client.B().FtInfo().Index(b).Build())
			redisResultMap map[string]rueidis.RedisMessage
			redisResultAny map[string]any
			bytes          []byte
			interim        = new(ftInfo)
		)
		switch redisResultMap, err = redisResult.AsMap(); {
		case err != nil:
			return
		}

		switch redisResultAny, err = parseRedisMessages(redisResultMap); {
		case err != nil:
			return
		}

		switch bytes, err = json.Marshal(redisResultAny); {
		case err != nil:
			return
		}

		switch err = json.Unmarshal(bytes, interim); {
		case err != nil:
			return
		}

		r.info[b] = interim
	}

	for a, b := range r.info {
		switch value := b.HashIndexingFailures; {
		case value != 0:
			l.Z{l.M: redisearchTagName, "index": a, "failures": value}.Warning()
		}
	}

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

func (r *RedisRepository) waitIndexing(indexName string) (err error) {
	switch _, ok := r.info[indexName]; {
	case !ok:
		return mod_errors.ENODATA
	}

	for err = r.getInfo(); r.info[indexName].Indexing != 0; err = r.getInfo() {
		switch {
		case err != nil:
			return
		}

		switch err = mod_reflect.WaitCtx(r.ctx, l.RetryInterval); {
		case err != nil:
			return
		}
	}

	return
}

func (r *RedisRepository) SaveEntry(e *Entry) (err error) {
	switch {
	case l.Run.DryRunValue():
		return
	}

	err = r.entry.Save(r.ctx, e)
	_ = r.getInfo()

	return
}

func (r *RedisRepository) SaveCert(e *Cert) (err error) {
	switch {
	case l.Run.DryRunValue():
		return
	}

	err = r.cert.Save(r.ctx, e)
	_ = r.getInfo()

	return
}

//

func (r *RedisRepository) SaveMultiEntry(e ...*Entry) (err []error) {
	switch {
	case l.Run.DryRunValue():
		return
	}

	err = r.entry.SaveMulti(r.ctx, e...)
	_ = r.getInfo()

	return
}

func (r *RedisRepository) SaveMultiCert(e ...*Cert) (err []error) {
	switch {
	case l.Run.DryRunValue():
		return
	}

	err = r.cert.SaveMulti(r.ctx, e...)
	_ = r.getInfo()

	return
}

//

func (r *RedisRepository) FindEntry(id string) (entry *Entry, err error) {
	return r.entry.Fetch(r.ctx, id)
}

func (r *RedisRepository) FindCert(id string) (cert *Cert, err error) {
	return r.cert.Fetch(r.ctx, id)
}

//

func (r *RedisRepository) DeleteEntry(id string) (err error) {
	switch {
	case l.Run.DryRunValue():
		return
	}

	err = r.entry.Remove(r.ctx, id)

	return
}

func (r *RedisRepository) DeleteCert(id string) (err error) {
	switch {
	case l.Run.DryRunValue():
		return
	}

	err = r.cert.Remove(r.ctx, id)

	return
}

//

func (r *RedisRepository) DropEntryIndex() (err error) {
	switch {
	case l.Run.DryRunValue():
		return
	}

	err = r.entry.DropIndex(r.ctx)

	return
}

func (r *RedisRepository) DropCertIndex() (err error) {
	switch {
	case l.Run.DryRunValue():
		return
	}

	err = r.cert.DropIndex(r.ctx)

	return
}

//

func (r *RedisRepository) SearchEntryQ(query string) (count int64, entries []*Entry, err error) {
	_ = r.waitIndexing(r.entry.IndexName())

	return r.entry.Search(r.ctx, func(search om.FtSearchIndex) rueidis.Completed {
		return search.Query(query).
			Limit().OffsetNum(0, connMaxPaging).
			Build()
	})
}

func (r *RedisRepository) SearchCertQ(query string) (count int64, entries []*Cert, err error) {
	_ = r.waitIndexing(r.cert.IndexName())

	return r.cert.Search(r.ctx, func(search om.FtSearchIndex) rueidis.Completed {
		return search.Query(query).
			Limit().OffsetNum(0, connMaxPaging).
			Build()
	})
}

func (r *RedisRepository) SearchEntryFV(field mod_strings.EntryFieldName, value string) (count int64, entries []*Entry, err error) {
	return r.SearchEntryFVs(&mod_strings.FVs{{field, value}})
}

func (r *RedisRepository) SearchCertFV(field mod_strings.EntryFieldName, value string) (count int64, entries []*Cert, err error) {
	return r.SearchCertFVs(&mod_strings.FVs{{field, value}})
}

func (r *RedisRepository) SearchEntryFVs(fvs *mod_strings.FVs) (count int64, entries []*Entry, err error) {
	return r.SearchEntryQ(r.info[_entry].Attributes.buildQuery(fvs))
}

func (r *RedisRepository) SearchCertFVs(fvs *mod_strings.FVs) (count int64, entries []*Cert, err error) {
	return r.SearchCertQ(r.info[_certificate].Attributes.buildQuery(fvs))
}

// SearchEntryFVsField is not working:
//
// err is `unexpected end of JSON input`
//
// JSONRepository receives empty JSON stream.
func (r *RedisRepository) SearchEntryFVsField(ctx context.Context, fvs *mod_strings.FVs, field mod_strings.EntryFieldName) (count int64, entries []*Entry, err error) {
	_ = r.waitIndexing(r.entry.IndexName())

	return r.entry.Search(ctx, func(search om.FtSearchIndex) rueidis.Completed {
		var (
			command = search.Query(r.info[_entry].Attributes.buildQuery(fvs)).
				Return(strconv.FormatInt(1, 10)).
				Identifier(field.String()).
				Limit().OffsetNum(0, connMaxPaging).
				Build()
		)
		l.Z{l.M: "redis", "command": strings.Join(command.Commands(), " ")}.Informational()

		return command
	})
}
