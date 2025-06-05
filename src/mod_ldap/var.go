package io_ldap

var (
	fields = map[string]string{
		"CN":                   "cn",                   //
		"DN":                   "dn",                   //
		"Description":          "description",          //
		"DestinationIndicator": "destinationIndicator", // remote interim host
		"DisplayName":          "displayName",          //
		"GIDNumber":            "gidNumber",            //
		"HomeDirectory":        "homeDirectory",        //
		"IPHostNumber":         "ipHostNumber",         // user's subnet
		"LabeledURI":           "labeledURI",           //
		"Mail":                 "mail",                 //
		"Member":               "member",               //
		"MemberOf":             "memberOf",             //
		"O":                    "o",                    //
		"OU":                   "ou",                   //
		"ObjectClass":          "objectClass",          //
		"Owner":                "owner",                //
		"SN":                   "sn",                   //
		"SSHPublicKey":         "sshPublicKey",         //
		"TelephoneNumber":      "telephoneNumber",      //
		"TelexNumber":          "telexNumber",          // Signaling
		"UID":                  "uid",                  //
		"UIDNumber":            "uidNumber",            //
		"UserPKCS12":           "userPKCS12",           //
		"UserPassword":         "userPassword",         //

		"CreateTimestamp": "createTimestamp", // service
		"CreatorsName":    "creatorsName",    // service
		"EntryUUID":       "entryUUID",       // service
		"ModifiersName":   "modifiersName",   // service
		"ModifyTimestamp": "modifyTimestamp", // service
	}
)
