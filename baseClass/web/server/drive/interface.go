package drive

import (
	"GeekTask/baseClass/web/server/route"
)

type Server interface {
	route.Routable

	Start(address string) error
}
