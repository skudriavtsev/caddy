// Copyright 2015 Light Code Labs, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ibproxy

import (
	"github.com/skudriavtsev/caddy"
	"github.com/skudriavtsev/caddy/caddyhttp/httpserver"
)

func init() {
	caddy.RegisterPlugin("ibproxy", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// setup configures a new Proxy middleware instance.
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
		return Proxy{Next: next, Suffix: suffix}
	})

	return nil
}
