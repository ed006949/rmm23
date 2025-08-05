package mod_net

import (
	"net/url"
)

type URL struct{ *url.URL }
type URLs []*URL
