package mod_crypto

import (
	"rmm23/src/l"
	"rmm23/src/mod_slices"
	"rmm23/src/mod_strings"
)

func (r *Certificates) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_strings.Normalize(values, mod_slices.FlagNormalize) {
		var (
			forErr  error
			interim = new(Certificate)
		)

		switch forErr = interim.ParsePEM([]byte(value)); {
		case forErr != nil:
			l.Z{l.M: "ParsePEM", l.E: forErr}.Warning()

			continue
		}

		(*r)[interim.Certificate.Subject.String()] = interim
	}

	return
}
