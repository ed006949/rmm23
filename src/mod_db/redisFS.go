package mod_db

import (
	"context"
	"io/fs"
	"strings"

	"rmm23/src/l"
	"rmm23/src/mod_errors"
	"rmm23/src/mod_vfs"
)

func (r *RedisRepository) GetFSCerts(ctx context.Context, vfsDB *mod_vfs.VFSDB) (err error) {
	var (
		content    = make(map[string][][]byte)
		fileExts   = 2
		totalFiles = 6

		fn = func(name string, dirEntry fs.DirEntry, err error) (fnErr error) {
			switch {
			case err != nil:
				return err
			}

			var (
				s = strings.Split(name, ".")
			)
			switch {
			case len(s) < fileExts:
				return
			}

			var (
				n = strings.Join(s[:len(s)-fileExts], ".")
			)
			switch _, ok := content[n]; {
			case !ok:
				content[n] = make([][]byte, totalFiles)
			}

			switch s[len(s)-1] {
			case "der":
				switch s[len(s)-2] {
				case "key":
					content[n][0], _ = vfsDB.VFS.ReadFile(name)
				case "crt":
					content[n][1], _ = vfsDB.VFS.ReadFile(name)
				case "ca":
					content[n][2], _ = vfsDB.VFS.ReadFile(name)
				case "csr":
					content[n][3], _ = vfsDB.VFS.ReadFile(name)
				case "crl":
					content[n][4], _ = vfsDB.VFS.ReadFile(name)
				}
			case "pem":
				content[n][5] = append(content[n][5], mod_errors.StripErr1(vfsDB.VFS.ReadFile(name))...)
			}

			return
		}
	)
	switch err = vfsDB.VFS.WalkDir("/", fn); {
	case err != nil:
		l.Z{l.E: err}.Error()
	}

	for a, b := range content {
		var (
			forErr  error
			forCert = new(Cert)
		)

		switch forErr = forCert.parseRaw(b...); {
		case forErr != nil:
			continue
		}

		_ = r.DeleteCert(forCert.Key)

		switch forErr = r.SaveCert(forCert); {
		case forErr != nil:
			l.Z{l.M: "r.SaveCert", "cert": a, l.E: forErr}.Warning()
		}
	}

	return
}
