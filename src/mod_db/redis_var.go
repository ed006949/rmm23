package mod_db

var (
	entryFieldMap = map[entryFieldName]string{
		_type:   rediSearchTagTypeNumeric,
		_status: rediSearchTagTypeNumeric,
		_baseDN: rediSearchTagTypeTag,

		_uuid:            rediSearchTagTypeTag,
		_dn:              rediSearchTagTypeTag,
		_objectClass:     rediSearchTagTypeTag,
		_creatorsName:    rediSearchTagTypeTag,
		_createTimestamp: rediSearchTagTypeTag,
		_modifiersName:   rediSearchTagTypeTag,
		_modifyTimestamp: rediSearchTagTypeTag,

		_cn:                   rediSearchTagTypeTag,
		_dc:                   rediSearchTagTypeTag,
		_description:          rediSearchTagTypeTag,
		_destinationIndicator: rediSearchTagTypeTag,
		_displayName:          rediSearchTagTypeTag,
		_gidNumber:            rediSearchTagTypeNumeric,
		_homeDirectory:        rediSearchTagTypeTag,
		_ipHostNumber:         rediSearchTagTypeTag,
		_mail:                 rediSearchTagTypeTag,
		_member:               rediSearchTagTypeTag,
		_o:                    rediSearchTagTypeTag,
		_ou:                   rediSearchTagTypeTag,
		_owner:                rediSearchTagTypeTag,
		_sn:                   rediSearchTagTypeTag,
		_sshPublicKey:         rediSearchTagTypeTag,
		_telephoneNumber:      rediSearchTagTypeTag,
		_telexNumber:          rediSearchTagTypeTag,
		_uid:                  rediSearchTagTypeTag,
		_uidNumber:            rediSearchTagTypeNumeric,
		_userPKCS12:           rediSearchTagTypeTag,
		_userPassword:         rediSearchTagTypeTag,

		_labeledURI: rediSearchTagTypeTag,
	}
	entryFieldValueEnclosure = map[string][2]string{
		rediSearchTagTypeText:    {enclosureEmpty0, enclosureEmpty1},
		rediSearchTagTypeTag:     {enclosureCurly0, enclosureCurly1},
		rediSearchTagTypeNumeric: {enclosureSquare0, enclosureSquare1},
		rediSearchTagTypeGeo:     {enclosureSquare0, enclosureSquare1},
	}
)
