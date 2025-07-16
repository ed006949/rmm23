package mod_errors

func StripErr(err error)                                 {}
func StripErr1[E any](inbound E, err error) (outbound E) { return inbound }
