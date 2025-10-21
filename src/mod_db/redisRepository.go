package mod_db

import (
	"context"
	"strconv"
	"strings"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"

	"rmm23/src/l"
	"rmm23/src/mod_reflect"
	"rmm23/src/mod_strings"
)

func (r *RedisRepository) SaveEntry(e *Entry) (err error) {
	l.Z{l.M: "save", "DN": e.DN.String()}.Informational()

	switch {
	case l.Run.DryRunValue():
		return
	}

	var (
		fn = func() error { return r.entry.Save(r.ctx, e) }
	)

	switch err = mod_reflect.RetryWithCtx(r.ctx, 0, l.RetryInterval, fn); {
	case err == nil:
		e.Status = entryStatusReady
	}

	_ = r.getInfo()

	return
}

func (r *RedisRepository) SaveCert(e *Cert) (err error) {
	l.Z{l.M: "save", "DN": e.Subject.String()}.Informational()

	switch {
	case l.Run.DryRunValue():
		return
	}

	var (
		fn = func() error { return r.cert.Save(r.ctx, e) }
	)

	switch err = mod_reflect.RetryWithCtx(r.ctx, 0, l.RetryInterval, fn); {
	case err == nil:
		e.Status = entryStatusReady
	}

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
	l.Z{l.M: "find", "DN": id}.Informational()

	return r.entry.Fetch(r.ctx, id)
}

func (r *RedisRepository) FindCert(id string) (cert *Cert, err error) {
	l.Z{l.M: "find", "DN": id}.Informational()

	return r.cert.Fetch(r.ctx, id)
}

//

func (r *RedisRepository) DeleteEntry(id string) (err error) {
	l.Z{l.M: "delete", "DN": id}.Informational()

	switch {
	case l.Run.DryRunValue():
		return
	}

	var (
		fn = func() error { return r.entry.Remove(r.ctx, id) }
	)

	err = mod_reflect.RetryWithCtx(r.ctx, 0, l.RetryInterval, fn)

	return
}

func (r *RedisRepository) DeleteCert(id string) (err error) {
	l.Z{l.M: "delete", "DN": id}.Informational()

	switch {
	case l.Run.DryRunValue():
		return
	}

	var (
		fn = func() error { return r.cert.Remove(r.ctx, id) }
	)

	err = mod_reflect.RetryWithCtx(r.ctx, 0, l.RetryInterval, fn)

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
	l.Z{l.M: "update", "DN": e.DN.String()}.Informational()

	switch e.Status {
	case entryStatusUpdate:
		e.Ver++
		err = r.SaveEntry(e)
	case entryStatusDelete:
		err = r.DeleteEntry(e.Key)
	default:
	}

	switch {
	case err != nil:
		l.Z{l.M: "entry", "flag": e.Status.String(), l.E: err, "DN": e.DN.String()}.Warning()
	}

	return
}
