package docid

import "github.com/rs/xid"

func New() string {

	return xid.New().String()
}
