package mod_db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"rmm23/src/l"
	"rmm23/src/mod_crypto"
)

//
// attrDN

func (r *attrDN) String() (outbound string) { return string(*r) }

//
// attrUUID

func (r *attrUUID) String() (outbound string) { return uuid.UUID(*r).String() }
func (r *attrUUID) Entry() (outbound string)  { return entryKeyHeader + ":" + uuid.UUID(*r).String() }

func (r *attrUUID) MarshalJSON() (outbound []byte, err error) {
	return []byte(fmt.Sprintf("%q", r.String())), nil
}

func (r *attrUUID) UnmarshalJSON(inbound []byte) (err error) {
	var (
		interim string
	)

	switch err = json.Unmarshal(inbound, &interim); {
	case err != nil:
		return
	}

	var (
		interimUUID uuid.UUID
	)
	switch interimUUID, err = uuid.Parse(interim); {
	case err != nil:
		return
	}

	*r = attrUUID(interimUUID)

	return
}

//
// attrTime
// store/retrieve time.Time as int64 to utilize redisearch NUMERIC search

func (r *attrTime) String() (outbound string) { return time.Time(*r).String() }

func (r *attrTime) MarshalJSON() (outbound []byte, err error) {
	return []byte(fmt.Sprintf("%d", time.Time(*r).Unix())), nil
}

func (r *attrTime) UnmarshalJSON(inbound []byte) (err error) {
	var (
		interim int64
		// interimTime time.Time
	)

	switch err = json.Unmarshal(inbound, &interim); {
	case err != nil:
		return
	}

	*r = attrTime(time.Unix(interim, 0))

	return
}

//
// attrUserPKCS12s
// store/retrieve as map[string]p12

func (r *attrUserPKCS12s) MarshalJSON() (outbound []byte, err error) {
	var (
		interim = make(map[string][]byte)
	)

	for a, b := range *r {
		var (
			pfxData []byte
			forErr  error
		)

		switch pfxData, forErr = b.EncodeP12(); {
		case forErr != nil:
			l.Z{l.M: "EncodeP12", l.E: forErr}.Warning()

			continue
		}

		interim[a.String()] = pfxData
	}

	return json.Marshal(interim)
}

func (r *attrUserPKCS12s) UnmarshalJSON(inbound []byte) (err error) {
	var (
		interim                = make(map[string][]byte)
		interimAttrUserPKCS12s = make(attrUserPKCS12s)
	)

	switch err = json.Unmarshal(inbound, &interim); {
	case err != nil:
		return
	}

	for a, b := range interim {
		var (
			forErr      error
			certificate = new(mod_crypto.Certificate)
		)

		switch forErr = certificate.DecodeP12(b); {
		case forErr != nil:
			l.Z{l.M: "DecodeP12", l.E: forErr}.Warning()

			continue
		}

		interimAttrUserPKCS12s[attrDN(a)] = certificate
	}

	*r = interimAttrUserPKCS12s

	return
}
