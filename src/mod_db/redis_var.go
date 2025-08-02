package mod_db

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
	}
	entryFieldValueEnclosure = map[string][2]string{
		redisearchTagTypeText:    {enclosureEmpty0, enclosureEmpty1},
		redisearchTagTypeTag:     {enclosureCurly0, enclosureCurly1},
		redisearchTagTypeNumeric: {enclosureSquare0, enclosureSquare1},
		redisearchTagTypeGeo:     {enclosureSquare0, enclosureSquare1},
	}
)
