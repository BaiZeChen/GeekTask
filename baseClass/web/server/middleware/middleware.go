package middleware

import (
	"GeekTask/baseClass/web/server/route"
)

type MiddlewareBuild func(next Middleware) Middleware

type Middleware func(c *route.Context)
