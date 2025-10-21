package mod_db

import (
	"rmm23/src/l"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_reflect"
)

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

		switch err = mod_reflect.WaitWithCtx(r.ctx, l.RetryInterval); {
		case err != nil:
			return
		}
	}

	return
}
