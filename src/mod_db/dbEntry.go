package mod_db

import (
	"time"

	"github.com/google/uuid"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/om"

	"rmm23/src/mod_dn"
	"rmm23/src/mod_strings"
	"rmm23/src/mod_time"
)

type DBEntry struct {
	// db mandatory
	Key string    `redis:",key"`  //
	Ver int64     `redis:",ver"`  //
	Ext time.Time `redis:",exat"` //

	// element  meta data
	Status attrEntryStatus `json:"status,omitempty"` //

	// db operational
	BaseDN          mod_dn.DN     `json:"baseDN,omitempty"`                                 //
	ObjectClass     []string      `json:"objectClass,omitempty"     ldap:"objectClass"`     //
	UUID            uuid.UUID     `json:"uuid,omitempty"            ldap:"entryUUID"`       //  must be unique
	DN              mod_dn.DN     `json:"dn,omitempty"              ldap:"entryDN"`         //  must be unique
	CreatorsName    mod_dn.DN     `json:"creatorsName,omitempty"    ldap:"creatorsName"`    //
	CreateTimestamp mod_time.Time `json:"createTimestamp,omitempty" ldap:"createTimestamp"` //
	ModifiersName   mod_dn.DN     `json:"modifiersName,omitempty"   ldap:"modifiersName"`   //
	ModifyTimestamp mod_time.Time `json:"modifyTimestamp,omitempty" ldap:"modifyTimestamp"` //
}

func (r *RedisRepository) CreateDBEntryIndex() (err error) {
	return r.entry.CreateIndex(r.ctx, func(schema om.FtCreateSchema) rueidis.Completed {
		return schema.
			FieldName(mod_strings.F_status.FieldName()).As(mod_strings.F_status.String()).Numeric().

			//
			FieldName(mod_strings.F_baseDN.FieldName()).As(mod_strings.F_baseDN.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_objectClass.FieldNameSlice()).As(mod_strings.F_objectClass.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_uuid.FieldName()).As(mod_strings.F_uuid.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_dn.FieldName()).As(mod_strings.F_dn.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_creatorsName.FieldName()).As(mod_strings.F_creatorsName.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_createTimestamp.FieldName()).As(mod_strings.F_createTimestamp.String()).Numeric().
			FieldName(mod_strings.F_modifiersName.FieldName()).As(mod_strings.F_modifiersName.String()).Tag().Separator(mod_strings.SliceSeparator).
			FieldName(mod_strings.F_modifyTimestamp.FieldName()).As(mod_strings.F_modifyTimestamp.String()).Numeric().

			//
			Build()
	})
}

type ObjectClassLabeledURI struct {
	LabeledURI []string `json:"labeledURI,omitempty" ldap:"labeledURI"`
}
