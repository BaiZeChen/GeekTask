package middleware

import "GeekTask/web/server/route"

type MiddlewareBuild func(next Middleware) Middleware

type Middleware func(c *route.Context)
