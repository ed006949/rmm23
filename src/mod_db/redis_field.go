package mod_db

import (
	"fmt"
	"strings"
)

const (
	enclosureEmpty0  = ""
	enclosureEmpty1  = ""
	enclosureSquare0 = "["
	enclosureSquare1 = "]"
	enclosureCurly0  = "{"
	enclosureCurly1  = "}"
)

type MFV []FV

type FV struct {
	Field entryFieldName
	Value string
}

var (
	entryFieldMap = map[entryFieldName]string{
		F_type:   redisearchTagTypeNumeric,
		F_status: redisearchTagTypeNumeric,
		F_baseDN: redisearchTagTypeTag,

		F_uuid:            redisearchTagTypeTag,
		F_dn:              redisearchTagTypeTag,
		F_objectClass:     redisearchTagTypeTag,
		F_creatorsName:    redisearchTagTypeTag,
		F_createTimestamp: redisearchTagTypeTag,
		F_modifiersName:   redisearchTagTypeTag,
		F_modifyTimestamp: redisearchTagTypeTag,

		F_cn:                   redisearchTagTypeTag,
		F_dc:                   redisearchTagTypeTag,
		F_description:          redisearchTagTypeTag,
		F_destinationIndicator: redisearchTagTypeTag,
		F_displayName:          redisearchTagTypeTag,
		F_gidNumber:            redisearchTagTypeNumeric,
		F_homeDirectory:        redisearchTagTypeTag,
		F_ipHostNumber:         redisearchTagTypeTag,
		F_mail:                 redisearchTagTypeTag,
		F_member:               redisearchTagTypeTag,
		F_memberOf:             redisearchTagTypeTag,
		F_o:                    redisearchTagTypeTag,
		F_ou:                   redisearchTagTypeTag,
		F_owner:                redisearchTagTypeTag,
		F_sn:                   redisearchTagTypeTag,
		F_sshPublicKey:         redisearchTagTypeTag,
		F_telephoneNumber:      redisearchTagTypeTag,
		F_telexNumber:          redisearchTagTypeTag,
		F_uid:                  redisearchTagTypeTag,
		F_uidNumber:            redisearchTagTypeNumeric,
		F_userPKCS12:           redisearchTagTypeTag,
		F_userPassword:         redisearchTagTypeTag,

		F_labeledURI: redisearchTagTypeTag,

		F_serialNumber:   redisearchTagTypeNumeric,
		F_issuer:         redisearchTagTypeTag,
		F_subject:        redisearchTagTypeTag,
		F_notBefore:      redisearchTagTypeTag,
		F_notAfter:       redisearchTagTypeTag,
		F_dnsNames:       redisearchTagTypeTag,
		F_emailAddresses: redisearchTagTypeTag,
		F_ipAddresses:    redisearchTagTypeTag,
		F_uris:           redisearchTagTypeTag,
		F_isCA:           redisearchTagTypeNumeric,
	}
	entryFieldValueEnclosure = map[string][2]string{
		redisearchTagTypeText:    {enclosureEmpty0, enclosureEmpty1},
		redisearchTagTypeTag:     {enclosureCurly0, enclosureCurly1},
		redisearchTagTypeNumeric: {enclosureSquare0, enclosureSquare1},
		redisearchTagTypeGeo:     {enclosureSquare0, enclosureSquare1},
	}
)

const (
	F_key entryFieldName = "key"
	F_ver entryFieldName = "ver"

	F_type   entryFieldName = "type"
	F_status entryFieldName = "status"
	F_baseDN entryFieldName = "baseDN"

	F_uuid            entryFieldName = "uuid"
	F_dn              entryFieldName = "dn"
	F_objectClass     entryFieldName = "objectClass"
	F_creatorsName    entryFieldName = "creatorsName"
	F_createTimestamp entryFieldName = "createTimestamp"
	F_modifiersName   entryFieldName = "modifiersName"
	F_modifyTimestamp entryFieldName = "modifyTimestamp"

	F_cn                   entryFieldName = "cn"
	F_dc                   entryFieldName = "dc"
	F_description          entryFieldName = "description"
	F_destinationIndicator entryFieldName = "destinationIndicator"
	F_displayName          entryFieldName = "displayName"
	F_gidNumber            entryFieldName = "gidNumber"
	F_homeDirectory        entryFieldName = "homeDirectory"
	F_ipHostNumber         entryFieldName = "ipHostNumber"
	F_mail                 entryFieldName = "mail"
	F_member               entryFieldName = "member"
	F_memberOf             entryFieldName = "memberOf"
	F_o                    entryFieldName = "o"
	F_ou                   entryFieldName = "ou"
	F_owner                entryFieldName = "owner"
	F_sn                   entryFieldName = "sn"
	F_sshPublicKey         entryFieldName = "sshPublicKey"
	F_telephoneNumber      entryFieldName = "telephoneNumber"
	F_telexNumber          entryFieldName = "telexNumber"
	F_uid                  entryFieldName = "uid"
	F_uidNumber            entryFieldName = "uidNumber"
	F_userPKCS12           entryFieldName = "userPKCS12"
	F_userPassword         entryFieldName = "userPassword"

	F_labeledURI entryFieldName = "labeledURI"

	F_serialNumber   entryFieldName = "serialNumber"
	F_issuer         entryFieldName = "issuer"
	F_subject        entryFieldName = "subject"
	F_notBefore      entryFieldName = "notBefore"
	F_notAfter       entryFieldName = "notAfter"
	F_dnsNames       entryFieldName = "dnsNames"
	F_emailAddresses entryFieldName = "emailAddresses"
	F_ipAddresses    entryFieldName = "ipAddresses"
	F_uris           entryFieldName = "uris"
	F_isCA           entryFieldName = "isCA"
)

func (r *MFV) buildMFVQuery() (outbound string) {
	var (
		interim = make([]string, len(*r), len(*r))
	)

	for i, fv := range *r {
		interim[i] = buildFVQuery(fv.Field, fv.Value)
	}

	return strings.Join(interim, " ")
}

func buildFVQuery(field entryFieldName, value string) (outbound string) {
	return fmt.Sprintf(
		"@%s:%s%v%s",
		field.String(),
		entryFieldValueEnclosure[entryFieldMap[field]][0],
		escapeQueryValue(value),
		entryFieldValueEnclosure[entryFieldMap[field]][1],
	)
}

func escapeQueryValue(inbound string) (outbound string) {
	replacer := strings.NewReplacer(
		`=`, `\=`, //
		`,`, `\,`, //
		`(`, `\(`, //
		`)`, `\)`, //
		`{`, `\{`, //
		`}`, `\}`, //
		`[`, `\[`, //
		`]`, `\]`, //
		`"`, `\"`, //
		`'`, `\'`, //
		`~`, `\~`, //
		`-`, `\-`, // (?)
	)

	return replacer.Replace(inbound)
}
