package mod_dn

func UnmarshalText(inbound []byte) (outbound DN, err error) {
	var (
		interim = new(DN)
	)
	switch err = interim.UnmarshalText(inbound); {
	case err != nil:
		return
	}

	return *interim, err
}
