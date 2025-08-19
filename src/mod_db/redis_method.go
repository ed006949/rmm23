package mod_db

import (
	"context"
	"encoding/json/v2"
	"os"
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

	_ = r.Repo.getInfo(ctx)

	os.Exit(1)

	return
}

func (r *RedisRepository) getInfo(ctx context.Context) (err error) {
	var (
		repos []string
	)
	switch repos, err = r.client.Do(ctx, r.client.B().FtList().Build()).AsStrSlice(); {
	case err != nil:
		return
	}

	for _, b := range repos {
		var (
			redisResult    = r.client.Do(ctx, r.client.B().FtInfo().Index(b).Build())
			redisResultMap map[string]rueidis.RedisMessage
			redisResultAny map[string]any
			bytes          []byte
			ftInfo         = new(FTInfo)
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

		switch err = json.Unmarshal(bytes, ftInfo); {
		case err != nil:
			return
		}
	}

	return
}

func parseRedisMessages(messages map[string]rueidis.RedisMessage) (outbound map[string]any, err error) {
	outbound = make(map[string]any)

	for c, message := range messages {
		switch outbound[c], err = parseRedisMessage(message); {
		case err != nil:
			return nil, err
		}
	}

	return
}

func parseRedisMessage(message rueidis.RedisMessage) (outbound any, err error) {
	switch {
	// case message.IsCacheHit():

	case message.IsArray():
		switch messages, swErr := message.AsMap(); {
		case swErr == nil:
			return parseRedisMessages(messages)
		}

		switch messages, swErr := message.ToArray(); {
		case swErr == nil:
			var (
				interim []any
			)

			for _, b := range messages {
				switch message2, swErr2 := parseRedisMessage(b); {
				case swErr2 != nil:
					return nil, swErr2
				default:
					interim = append(interim, message2)
				}
			}

			return interim, nil
		}

		return nil, mod_errors.EParse

	case message.IsBool():
		return message.AsBool()

	case message.IsFloat64():
		return message.AsFloat64()

	case message.IsInt64():
		return message.AsInt64()

	case message.IsNil():
		return nil, nil

	case message.IsMap():
		return message.ToMap()

	case message.IsString():
		return message.ToString()

	default:
		return message.ToAny()
	}
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
		r.entry.IndexName():  r.client.Do(ctx, r.client.B().FtInfo().Index(r.entry.IndexName()).Build()),
		r.cert.IndexName():   r.client.Do(ctx, r.client.B().FtInfo().Index(r.cert.IndexName()).Build()),
		r.issued.IndexName(): r.client.Do(ctx, r.client.B().FtInfo().Index(r.cert.IndexName()).Build()),
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
			case mod_strings.FTInfo_hash_indexing_failures:
				switch c, d := strconv.ParseInt(b, 10, 64); {
				case d == nil && c == 0:
					l.Z{l.M: "redis", "index": info[mod_strings.FTInfo_index_name], a: b}.Debug()
				case d == nil:
					l.Z{l.M: "redis", "index": info[mod_strings.FTInfo_index_name], a: b}.Warning()
				default:
					l.Z{l.M: "redis info", "index": info[mod_strings.FTInfo_index_name], l.E: err, a: b}.Error()
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

	for info, forErr := resp.AsStrMap(); info[mod_strings.FTInfo_percent_indexed] != "1"; info, forErr = resp.AsStrMap() {
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
	return r.SearchEntryMFV(ctx, &mod_strings.FVs{{field, value}})
}

func (r *RedisRepository) SearchCertFV(ctx context.Context, field mod_strings.EntryFieldName, value string) (count int64, entries []*Cert, err error) {
	return r.SearchCertMFV(ctx, &mod_strings.FVs{{field, value}})
}

func (r *RedisRepository) SearchEntryMFV(ctx context.Context, mfv *mod_strings.FVs) (count int64, entries []*Entry, err error) {
	return r.SearchEntryQ(ctx, elementFieldMap.BuildQuery(mfv))
}

func (r *RedisRepository) SearchCertMFV(ctx context.Context, mfv *mod_strings.FVs) (count int64, entries []*Cert, err error) {
	return r.SearchCertQ(ctx, elementFieldMap.BuildQuery(mfv))
}

// SearchEntryMFVField is not working:
//
// err is `unexpected end of JSON input`
//
// JSONRepository receives empty JSON stream.
func (r *RedisRepository) SearchEntryMFVField(ctx context.Context, mfv *mod_strings.FVs, field mod_strings.EntryFieldName) (count int64, entries []*Entry, err error) {
	_ = r.waitEntryIndexing(ctx)

	return r.entry.Search(ctx, func(search om.FtSearchIndex) rueidis.Completed {
		var (
			command = search.Query(elementFieldMap.BuildQuery(mfv)).
				Return(strconv.FormatInt(1, 10)).
				Identifier(field.String()).
				Limit().OffsetNum(0, connMaxPaging).
				Build()
		)
		l.Z{l.M: "redis", "command": strings.Join(command.Commands(), " ")}.Informational()

		return command
	})
}
