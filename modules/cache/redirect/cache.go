package redirect

import (
	"proxy/modules/config"
)

type Cache interface {
	Add(hash string, server config.Server)
	Find(hash string) (config.Server, bool)
	Clear()
}
