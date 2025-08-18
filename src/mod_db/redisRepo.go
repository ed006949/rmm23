package mod_db

import (
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"

	"rmm23/src/mod_strings"
)

var (
	elementFieldMap = map[mod_strings.EntryFieldName]string{
		mod_strings.F_type:   redisearchTagTypeNumeric,
		mod_strings.F_status: redisearchTagTypeNumeric,
		mod_strings.F_baseDN: redisearchTagTypeTag,

		mod_strings.F_uuid:            redisearchTagTypeTag,
		mod_strings.F_dn:              redisearchTagTypeTag,
		mod_strings.F_objectClass:     redisearchTagTypeTag,
		mod_strings.F_creatorsName:    redisearchTagTypeTag,
		mod_strings.F_createTimestamp: redisearchTagTypeTag,
		mod_strings.F_modifiersName:   redisearchTagTypeTag,
		mod_strings.F_modifyTimestamp: redisearchTagTypeTag,

		mod_strings.F_cn:                   redisearchTagTypeTag,
		mod_strings.F_dc:                   redisearchTagTypeTag,
		mod_strings.F_description:          redisearchTagTypeTag,
		mod_strings.F_destinationIndicator: redisearchTagTypeTag,
		mod_strings.F_displayName:          redisearchTagTypeTag,
		mod_strings.F_gidNumber:            redisearchTagTypeNumeric,
		mod_strings.F_homeDirectory:        redisearchTagTypeTag,
		mod_strings.F_ipHostNumber:         redisearchTagTypeTag,
		mod_strings.F_mail:                 redisearchTagTypeTag,
		mod_strings.F_member:               redisearchTagTypeTag,
		mod_strings.F_memberOf:             redisearchTagTypeTag,
		mod_strings.F_o:                    redisearchTagTypeTag,
		mod_strings.F_ou:                   redisearchTagTypeTag,
		mod_strings.F_owner:                redisearchTagTypeTag,
		mod_strings.F_sn:                   redisearchTagTypeTag,
		mod_strings.F_sshPublicKey:         redisearchTagTypeTag,
		mod_strings.F_telephoneNumber:      redisearchTagTypeTag,
		mod_strings.F_telexNumber:          redisearchTagTypeTag,
		mod_strings.F_uid:                  redisearchTagTypeTag,
		mod_strings.F_uidNumber:            redisearchTagTypeNumeric,
		mod_strings.F_userPKCS12:           redisearchTagTypeTag,
		mod_strings.F_userPassword:         redisearchTagTypeTag,

		mod_strings.F_labeledURI: redisearchTagTypeTag,

		mod_strings.F_serialNumber:   redisearchTagTypeNumeric,
		mod_strings.F_issuer:         redisearchTagTypeTag,
		mod_strings.F_subject:        redisearchTagTypeTag,
		mod_strings.F_notBefore:      redisearchTagTypeTag,
		mod_strings.F_notAfter:       redisearchTagTypeTag,
		mod_strings.F_dnsNames:       redisearchTagTypeTag,
		mod_strings.F_emailAddresses: redisearchTagTypeTag,
		mod_strings.F_ipAddresses:    redisearchTagTypeTag,
		mod_strings.F_uris:           redisearchTagTypeTag,
		mod_strings.F_isCA:           redisearchTagTypeTag,
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
