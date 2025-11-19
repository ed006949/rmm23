package main

import (
	"github.com/avfs/avfs"
	"github.com/avfs/avfs/vfs/memfs"
	"github.com/rs/zerolog/log"

	"rmm23/src/l"
	"rmm23/src/mod_vfs"
)

func main() {
	var (
		err error
	)

	defer func() {
		ctxCancel()

		switch err {
		case nil:
			log.Info().
				Msg("exit")
		default:
			log.Fatal().
				Err(err).
				Msg("exited with error")
		}
	}()

	l.Initialize(ctx)

	log.Info().
		Str("commit", l.Run.CommitHashValue()).
		Str("built", l.Run.BuildTimeValue()).
		Bool(l.Run.DryRunName(), l.Run.DryRunValue()).
		Msg("main")

	var (
		config = new(ConfigRoot)
		vfsDB  = &mod_vfs.VFSDB{
			List: make(map[string]string),
			VFS: memfs.NewWithOptions(
				&memfs.Options{
					Idm:        avfs.NotImplementedIdm,
					User:       nil,
					Name:       "",
					OSType:     avfs.CurrentOSType(),
					SystemDirs: nil,
				},
			),
		}
	)
	switch err = l.Run.ConfigUnmarshal(&config); {
	case err != nil:
		return
	}

	switch err = config.Conf.DB.Dial(ctx); {
	case err != nil:
		return
	}

	defer func() {
		_ = config.Conf.DB.Close()
	}()

	switch {
	case !l.Run.DryRunValue():
		switch err = config.Conf.DB.Repo.GetLDAPDocs(ctx, config.Conf.LDAP); {
		case err != nil:
			return
		}
	}

	switch err = vfsDB.CopyFromFS(config.Conf.Legacy.PKI); {
	case err != nil:
		return
	}

	switch {
	case !l.Run.DryRunValue():
		switch err = config.Conf.DB.Repo.GetFSCerts(ctx, vfsDB); {
		case err != nil:
			return
		}
	}

	switch err = config.Conf.DB.Repo.CheckIPHostNumber(config.Conf.Networking.User.Subnet, config.Conf.Networking.User.Bits); {
	case err != nil:
		return
	}

	// a, _ := mod_net.Subnets.Subnets(netip.MustParsePrefix("10.92.0.0/16"), config.Conf.Networking.User.Bits)
	// for c, d := range a {
	// 	switch {
	// 	case c >= (8192-64) && c < (8192+64):
	// 		log.Info().Msgf("%05d %18s", c, d)
	// 	}

	// }

	return
}
