package mod_db

// type attrTime struct{ time.Time }
//
// func (r *attrTime) String() (outbound string) { return r.Time.String() }
//
// // func (r *attrTime) MarshalText() (outbound []byte, err error) { return r.Time.MarshalText() }
// //
// // // UnmarshalText and all `attrTime` is for LDAP "specific" behavior
// // func (r *attrTime) UnmarshalText(inbound []byte) (err error) {
// // 	switch swInterim, swErr := ber.ParseGeneralizedTime(inbound); {
// // 	case swErr == nil:
// // 		r.Time = swInterim
// // 		return
// // 	}
// // 	var (
// // 		interim time.Time
// // 	)
// // 	switch err = interim.UnmarshalText(inbound); {
// // 	case err != nil:
// // 		return
// // 	}
// //
// // 	r.Time = interim
// //
// // 	return
// // }
//
// func (r *attrTime) MarshalJSON() (outbound []byte, err error) { return r.Time.MarshalJSON() }
//
// func (r *attrTime) UnmarshalJSON(inbound []byte) (err error) {
// 	var (
// 		interim string
// 	)
// 	switch err = json.Unmarshal(inbound, &interim); {
// 	case err != nil:
// 		return
// 	}
//
// 	switch swInterim, swErr := ber.ParseGeneralizedTime([]byte(interim)); {
// 	case swErr == nil:
// 		r.Time = swInterim
//
// 		return
// 	}
//
// 	switch swInterim, swErr := time.Parse(time.RFC3339, interim); {
// 	case swErr == nil:
// 		r.Time = swInterim
//
// 		return
// 	}
//
// 	return mod_errors.EParse
// }
