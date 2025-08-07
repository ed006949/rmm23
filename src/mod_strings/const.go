package mod_strings

const (
	SliceSeparator  = "\x1f"
	JSONPathHeader  = "$."
	TagSeparator    = ","
	HeaderSeparator = ":"
)

const (
	F_key EntryFieldName = "key"
	F_ver EntryFieldName = "ver"

	F_type   EntryFieldName = "type"
	F_status EntryFieldName = "status"
	F_baseDN EntryFieldName = "baseDN"

	F_uuid            EntryFieldName = "uuid"
	F_dn              EntryFieldName = "dn"
	F_objectClass     EntryFieldName = "objectClass"
	F_creatorsName    EntryFieldName = "creatorsName"
	F_createTimestamp EntryFieldName = "createTimestamp"
	F_modifiersName   EntryFieldName = "modifiersName"
	F_modifyTimestamp EntryFieldName = "modifyTimestamp"

	F_cn                   EntryFieldName = "cn"
	F_dc                   EntryFieldName = "dc"
	F_description          EntryFieldName = "description"
	F_destinationIndicator EntryFieldName = "destinationIndicator"
	F_displayName          EntryFieldName = "displayName"
	F_gidNumber            EntryFieldName = "gidNumber"
	F_homeDirectory        EntryFieldName = "homeDirectory"
	F_ipHostNumber         EntryFieldName = "ipHostNumber"
	F_mail                 EntryFieldName = "mail"
	F_member               EntryFieldName = "member"
	F_memberOf             EntryFieldName = "memberOf"
	F_o                    EntryFieldName = "o"
	F_ou                   EntryFieldName = "ou"
	F_owner                EntryFieldName = "owner"
	F_sn                   EntryFieldName = "sn"
	F_sshPublicKey         EntryFieldName = "sshPublicKey"
	F_telephoneNumber      EntryFieldName = "telephoneNumber"
	F_telexNumber          EntryFieldName = "telexNumber"
	F_uid                  EntryFieldName = "uid"
	F_uidNumber            EntryFieldName = "uidNumber"
	F_userPKCS12           EntryFieldName = "userPKCS12"
	F_userPassword         EntryFieldName = "userPassword"

	F_labeledURI EntryFieldName = "labeledURI"

	F_serialNumber   EntryFieldName = "serialNumber"
	F_issuer         EntryFieldName = "issuer"
	F_subject        EntryFieldName = "subject"
	F_notBefore      EntryFieldName = "notBefore"
	F_notAfter       EntryFieldName = "notAfter"
	F_dnsNames       EntryFieldName = "dnsNames"
	F_emailAddresses EntryFieldName = "emailAddresses"
	F_ipAddresses    EntryFieldName = "ipAddresses"
	F_uris           EntryFieldName = "uris"
	F_isCA           EntryFieldName = "isCA"
)
