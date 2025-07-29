package mod_crypto

import (
	"rmm23/src/l"
	"rmm23/src/mod_slices"
)

func (r *Certificates) UnmarshalLDAPAttr(values []string) (err error) {
	for _, value := range mod_slices.StringsNormalize(values, mod_slices.FlagNormalize) {
		var (
			forErr  error
			interim *Certificate
		)

		switch interim, forErr = ParsePEM([]byte(value)); {
		case forErr != nil:
			l.Z{l.M: "ParsePEM", l.E: forErr}.Warning()

			continue
		}

		(*r)[interim.Certificates[0].Subject.String()] = interim
	}

	return
}
