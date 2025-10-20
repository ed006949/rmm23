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

func (r *RedisRepository) getInfo(indexNames ...string) (err error) {
	mod_reflect.MakeMapIfNil(&r.info)

	switch {
	case len(indexNames) == 0:
		switch indexNames, err = r.client.Do(r.ctx, r.client.B().FtList().Build()).AsStrSlice(); {
		case err != nil:
			return
		}
	}

	for _, indexName := range indexNames {
		var (
			redisResult    = r.client.Do(r.ctx, r.client.B().FtInfo().Index(indexName).Build())
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

		r.info[indexName] = interim
	}

	_ = r.checkIndexFailure()

	return
	// return r.checkIndexExist(indexNames...)
}

func (r *RedisRepository) checkIndexFailure() (err error) {
	for indexName, indexInfo := range r.info {
		switch value := indexInfo.HashIndexingFailures; {
		case value != 0:
			err = mod_errors.EINVAL

			l.Z{l.M: redisearchTagName, "index": indexName, "failures": value}.Warning()
		}
	}

	return
}

func (r *RedisRepository) checkIndexExist(indexNames ...string) (err error) {
	for _, indexName := range indexNames {
		switch _, ok := r.info[indexName]; {
		case !ok:
			err = mod_errors.EUnwilling
			l.Z{l.M: redisearchTagName, "index": indexName, l.E: mod_errors.ENOTFOUND}.Error()
		}
	}

	return
}

func (r *RedisRepository) waitIndex(indexName string) (err error) {
	// switch err = r.checkIndexExist(indexName); {
	// case err != nil:
	// 	return
	// }
	switch err = r.getInfo(indexName); {
	case err != nil:
		return
	}

	// for err = r.getInfo(indexName); r.info[indexName].Indexing != 0 || r.info[indexName].PercentIndexed != 1; err = r.getInfo(indexName) {
	for r.info[indexName].Indexing != 0 || r.info[indexName].PercentIndexed != 1 {
		switch err = r.getInfo(indexName); {
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
	_ = r.waitIndex(r.entry.IndexName())

	return r.entry.Search(r.ctx, func(search om.FtSearchIndex) rueidis.Completed {
		return search.Query(query).
			Limit().OffsetNum(0, connMaxPaging).
			Build()
	})
}

func (r *RedisRepository) SearchCertQ(query string) (count int64, entries []*Cert, err error) {
	_ = r.waitIndex(r.cert.IndexName())

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
	_ = r.waitIndex(r.entry.IndexName())

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

//

func (r *RedisRepository) UpdateEntry(e *Entry) (err error) {
	switch e.Status {
	case EntryStatusUpdated:
		e.Ver++
		err = r.SaveEntry(e)
	case EntryStatusDeleted:
		err = r.DeleteEntry(e.Key)
	default:
	}

	switch err {
	case nil:
		l.Z{l.M: "entry", "flag": e.Status.String(), "DN": e.DN.String()}.Informational()
	default:
		l.Z{l.M: "entry", "flag": e.Status.String(), l.E: err, "DN": e.DN.String()}.Warning()
	}

	return
}

// func (r *RedisRepository) UpdateMultiEntry(e ...*Entry) (err []error) {
// 	var (
// 		fErr  = make([]error, len(e), len(e))
// 		isErr bool
// 	)
// 	mod_reflect.MakeSliceIfNil(&err, len(e))
//
// 	for a, b := range e {
// 		var (
// 			forErr error
// 		)
//
// 		switch b.Status {
// 		case EntryStatusUpdated:
// 			l.Z{l.M: "updated entry", "DN": b.DN.String()}.Informational()
// 			b.Ver++
// 			forErr = r.SaveEntry(b)
// 			// switch forErr = r.SaveEntry(b); {
// 			// case swErr != nil:
// 			// 	fErr[a] = swErr
// 			// 	l.Z{l.E: swErr, "DN": b.DN.String()}.Warning()
// 			// }
// 		case EntryStatusDeleted:
// 			l.Z{l.M: "deleted entry", "DN": b.DN.String()}.Informational()
// 			forErr = r.DeleteEntry(b.Key)
// 			// switch swErr := r.DeleteEntry(b.Key); {
// 			// case swErr != nil:
// 			// 	fErr[a] = swErr
// 			// 	l.Z{l.E: swErr, "DN": b.DN.String()}.Warning()
// 			// }
// 		default:
// 		}
//
// 		switch {
// 		case forErr != nil:
// 			fErr[a] = forErr
// 			isErr = true
//
// 			l.Z{l.M: b.Status.String(), l.E: forErr, "DN": b.DN.String()}.Warning()
// 		}
// 	}
//
// 	switch {
// 	case isErr:
// 		return fErr
// 	}
//
// 	return
// }

// func (r *RedisRepository) UpdateMultiCert(e ...*Cert) (err []error) {
// 	err = r.cert.SaveMulti(r.ctx, e...)
// 	_ = r.getInfo()
//
// 	return
// }

//
