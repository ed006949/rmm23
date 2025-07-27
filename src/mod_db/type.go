package mod_db

import (
	"github.com/RediSearch/redisearch-go/redisearch"

	"rmm23/src/mod_net"
)

type Conf struct {
	URL       *mod_net.URL `json:"url,omitempty"`
	Name      string       `json:"name,omitempty"`
	rsClient  *redisearch.Client
	rcNetwork string
}
