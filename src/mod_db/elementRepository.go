package mod_db

import (
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"

	"rmm23/src/mod_strings"
)

var (
	elementFieldMap = mod_strings.EntryFieldMap{
		mod_strings.F_type:   mod_strings.RedisearchTagTypeNumeric,
		mod_strings.F_status: mod_strings.RedisearchTagTypeNumeric,
		mod_strings.F_baseDN: mod_strings.RedisearchTagTypeTag,

		mod_strings.F_uuid:            mod_strings.RedisearchTagTypeTag,
		mod_strings.F_dn:              mod_strings.RedisearchTagTypeTag,
		mod_strings.F_objectClass:     mod_strings.RedisearchTagTypeTag,
		mod_strings.F_creatorsName:    mod_strings.RedisearchTagTypeTag,
		mod_strings.F_createTimestamp: mod_strings.RedisearchTagTypeTag,
		mod_strings.F_modifiersName:   mod_strings.RedisearchTagTypeTag,
		mod_strings.F_modifyTimestamp: mod_strings.RedisearchTagTypeTag,

		mod_strings.F_cn:                   mod_strings.RedisearchTagTypeTag,
		mod_strings.F_dc:                   mod_strings.RedisearchTagTypeTag,
		mod_strings.F_description:          mod_strings.RedisearchTagTypeTag,
		mod_strings.F_destinationIndicator: mod_strings.RedisearchTagTypeTag,
		mod_strings.F_displayName:          mod_strings.RedisearchTagTypeTag,
		mod_strings.F_gidNumber:            mod_strings.RedisearchTagTypeNumeric,
		mod_strings.F_homeDirectory:        mod_strings.RedisearchTagTypeTag,
		mod_strings.F_ipHostNumber:         mod_strings.RedisearchTagTypeTag,
		mod_strings.F_mail:                 mod_strings.RedisearchTagTypeTag,
		mod_strings.F_member:               mod_strings.RedisearchTagTypeTag,
		mod_strings.F_memberOf:             mod_strings.RedisearchTagTypeTag,
		mod_strings.F_o:                    mod_strings.RedisearchTagTypeTag,
		mod_strings.F_ou:                   mod_strings.RedisearchTagTypeTag,
		mod_strings.F_owner:                mod_strings.RedisearchTagTypeTag,
		mod_strings.F_sn:                   mod_strings.RedisearchTagTypeTag,
		mod_strings.F_sshPublicKey:         mod_strings.RedisearchTagTypeTag,
		mod_strings.F_telephoneNumber:      mod_strings.RedisearchTagTypeTag,
		mod_strings.F_telexNumber:          mod_strings.RedisearchTagTypeTag,
		mod_strings.F_uid:                  mod_strings.RedisearchTagTypeTag,
		mod_strings.F_uidNumber:            mod_strings.RedisearchTagTypeNumeric,
		mod_strings.F_userPKCS12:           mod_strings.RedisearchTagTypeTag,
		mod_strings.F_userPassword:         mod_strings.RedisearchTagTypeTag,

		mod_strings.F_labeledURI: mod_strings.RedisearchTagTypeTag,

		mod_strings.F_serialNumber:   mod_strings.RedisearchTagTypeNumeric,
		mod_strings.F_issuer:         mod_strings.RedisearchTagTypeTag,
		mod_strings.F_subject:        mod_strings.RedisearchTagTypeTag,
		mod_strings.F_notBefore:      mod_strings.RedisearchTagTypeTag,
		mod_strings.F_notAfter:       mod_strings.RedisearchTagTypeTag,
		mod_strings.F_dnsNames:       mod_strings.RedisearchTagTypeTag,
		mod_strings.F_emailAddresses: mod_strings.RedisearchTagTypeTag,
		mod_strings.F_ipAddresses:    mod_strings.RedisearchTagTypeTag,
		mod_strings.F_uris:           mod_strings.RedisearchTagTypeTag,
		mod_strings.F_isCA:           mod_strings.RedisearchTagTypeTag,
	}
)

// RedisRepository provides methods for interacting with Redis using rueidis.
type RedisRepository struct {
	client rueidis.Client
	entry  om.Repository[Entry]
	cert   om.Repository[Cert]
	issued om.Repository[Cert]
}

// type REntries om.Repository[Entry]
// type RCerts om.Repository[Cert]
// type RIssued om.Repository[Cert]
