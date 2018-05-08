package ibredirect

import (
	"sync"
	"time"

	"github.com/skudriavtsev/caddy"
	"github.com/skudriavtsev/caddy/caddyhttp/httpserver"
)

const (
	templateDir = "templates"
	authTpl = templateDir + "/auth.tpl"

	authTokenCookieName = "ib-auth-token"

	tokenParamName = "ib-auth-token"
	originUrlParamName = "ib-origin-url"
)

type globalState struct {
	sync.RWMutex
}

var gState globalState

// Redirect is middleware to respond with HTTP redirects
type Redirect struct {
	Next   httpserver.Handler
	Suffix string
}

func init() {
	caddy.RegisterPlugin("ibredir", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})

	fakeAuthTokens["7b89ab36-524b-11e8-bc8b-3b395b206d1f"] = &authItem{
		Domains: []string{"domain.a", "domain.b", "test.pvt"},
		Expires: time.Now().Add(120 * time.Second),
	}
	fakeAuthTokens["7c20cd0e-524b-11e8-bcf5-2fb7fb8b1b6a"] = &authItem{
		Domains: []string{"domain.c", "domain.d", "test2.pvt"},
		Expires: time.Now().Add(120 * time.Second),
	}
}

// setup configures a new Redirect middleware instance.
func setup(c *caddy.Controller) error {
	var suffix string

	for c.Next() {
		if !c.NextArg() {
			// need an argument
			return c.ArgErr()
		}
		suffix = c.Val()
		if c.NextArg() {
			// only one argument allowed
			return c.ArgErr()
		}
	}

	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return Redirect{Next: next, Suffix: suffix}
	})

	return nil
}
