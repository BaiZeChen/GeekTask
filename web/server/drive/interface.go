package drive

import "GeekTask/web/server/route"

type Server interface {
	route.Routable

	Start(address string) error
}
