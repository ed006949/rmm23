package mod_time

// var (
// 	timeToMillis = json.MarshalToFunc(func(enc *jsontext.Encoder, t time.Time) error { return enc.WriteToken(jsontext.Int(t.Unix())) })
// 	millisToTime = json.UnmarshalFunc(
// 		func(data []byte, dst *time.Time) error {
// 			var (
// 				seconds int64
// 			)
// 			switch err := json.Unmarshal(data, &seconds); {
// 			case err != nil:
// 				return err
// 			}
// 			*dst = time.Unix(seconds, 0)
// 			return nil
// 		},
// 	)
// )
