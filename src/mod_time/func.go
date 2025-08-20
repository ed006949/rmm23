package mod_time

func UnmarshalText(inbound []byte) (outbound *Time, err error) {
	var (
		interim = new(Time)
	)
	switch err = interim.UnmarshalText(inbound); {
	case err != nil:
		return
	}

	return interim, nil
}
