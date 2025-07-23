package mod_log

import (
	"github.com/rs/zerolog"
)

func (r Object) MarshalZerologObject(e *zerolog.Event) {
	for a, b := range r {
		// switch a {
		// case E:
		// 	a = zerolog.ErrorFieldName
		// case M:
		// 	a = zerolog.MessageFieldName
		// 	// case T:
		// 	// 	a = zerolog.TypeFieldName
		// }

		switch value := b.(type) {
		// case nameValue:
		// e.Str(a, value.String())
		// case configValue:
		// e.Str(a, value.String())
		// case dryRunFlag:
		// e.Bool(a, value.Flag())
		// case modeValue:
		// e.Str(a, value.String())
		// case verbosityLevel:
		// e.Str(a, value.String())
		case error:
			e.AnErr(a, value)
		case []error:
			e.Errs(a, value)
		default:
			switch a {
			case T:
				e.Type(a, b)
			default:
				e.Interface(a, value)
			}
		}
	}

	switch {
	// case l.DryRunValue.Flag():
	// 	e.Bool(l.DryRunValue.Name(), l.DryRunValue.Flag())
	}
}
