package mod_db

import (
	"context"

	"github.com/redis/rueidis"

	"rmm23/src/l"
	"rmm23/src/mod_errors"
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

	switch err = r.Repo.getInfo(_entry, _certificate); {
	case err != nil:
		return
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
