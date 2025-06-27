package mod_db

import (
	"github.com/RediSearch/redisearch-go/redisearch"
)

func (e *ElementUser) RedisearchSchema() *redisearch.Schema {
	return redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextFieldOptions("uuid", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextFieldOptions("dn", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextField("objectClass")).
		AddField(redisearch.NewTextField("cn")).
		AddField(redisearch.NewTextField("description")).
		AddField(redisearch.NewTextField("destinationIndicator")).
		AddField(redisearch.NewTextFieldOptions("displayName", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericField("gidNumber")).
		AddField(redisearch.NewTextField("homeDirectory")).
		AddField(redisearch.NewTextFieldOptions("ipHostNumber", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextField("labeledURI")).
		AddField(redisearch.NewTextField("mail")).
		AddField(redisearch.NewTextField("memberOf")).
		AddField(redisearch.NewTextField("o")).
		AddField(redisearch.NewTextField("ou")).
		AddField(redisearch.NewTextField("sn")).
		AddField(redisearch.NewTextField("sshPublicKey")).
		AddField(redisearch.NewTextField("telephoneNumber")).
		AddField(redisearch.NewTextField("telexNumber")).
		AddField(redisearch.NewTextFieldOptions("uid", redisearch.TextFieldOptions{Sortable: true})).
		AddField(redisearch.NewNumericFieldOptions("uidNumber", redisearch.NumericFieldOptions{Sortable: true})).
		AddField(redisearch.NewTextField("userPKCS12")).
		AddField(redisearch.NewTextField("userPassword")).
		AddField(redisearch.NewTextField("creatorsName")).
		AddField(redisearch.NewTextField("createTimestamp")).
		AddField(redisearch.NewTextField("modifiersName")).
		AddField(redisearch.NewTextField("modifyTimestamp"))
}
